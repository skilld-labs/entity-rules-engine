// Package the rules logic
// The exported function ApplyOn check if requiered conditions are valid and execute actions from entityRules object

package entityrules

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"text/template"

	"github.com/hashicorp/go-bexpr"
)

type EntityRules struct {
	Conditions MethodsExecution
	Actions    MethodsExecution
	Rules      Rules
}
type MethodsExecution []MethodExecution

type MethodExecution struct {
	Name      string
	Method    string
	Arguments Arguments
}
type Arguments []Argument
type Argument interface{}

type Rules []Rule

type Rule struct {
	Name        string
	Description string
	When        []string
	If          string
	Do          []string
}

var functionMatcher = regexp.MustCompile(`[\w\d]+`)

const (
	TrueValue = "true"
)

var (
	ErrEmptyRules                 = errors.New("entity does not have rules")
	ErrEmptyRuleName              = errors.New("rule misses name")
	ErrEmptyRuleWhen              = errors.New("rule does not have when, if you do not want any when conditions set it to true")
	ErrEmptyRuleIf                = errors.New("rule does not have if, if you do not want any if conditions set it to true")
	ErrEmptyRuleDo                = errors.New("rule does not have do, please add one")
	ErrEmptyMethodExecutionName   = errors.New("methodexecution misses name")
	ErrEmptyMethodExecutionMethod = errors.New("methodexecution misses method")
	ErrMethodNotRegistered        = errors.New("method is not registered")
	ErrMethodReturnNotBool        = errors.New("method does not return boolean")
)

func (entityRules *EntityRules) ApplyOn(entities ...interface{}) error {
	if err := entityRules.Validate(); err != nil {
		return err
	}
	for _, e := range entities {
		if err := entityRules.applyOnOne(e); err != nil {
			return err
		}
	}
	return nil
}

func (entityRules *EntityRules) Validate() error {
	if err := entityRules.Actions.validate(); err != nil {
		return err
	}
	if err := entityRules.Conditions.validate(); err != nil {
		return err
	}
	if rr := entityRules.Rules; rr == nil {
		return ErrEmptyRules
	} else {
		for _, r := range rr {
			if r.Name == "" {
				return fmt.Errorf("%w: %v", ErrEmptyRuleName, r)
			}
			if r.When == nil {
				return fmt.Errorf("%w: %s", ErrEmptyRuleWhen, r.Name)
			}
			if r.If == "" {
				return fmt.Errorf("%w: %s", ErrEmptyRuleIf, r.Name)
			}
			if r.Do == nil {
				return fmt.Errorf("%w: %s", ErrEmptyRuleDo, r.Name)
			}
		}
	}
	return nil
}

func (mm MethodsExecution) validate() error {
	for _, m := range mm {
		if m.Name == "" {
			return fmt.Errorf("%w: %v", ErrEmptyMethodExecutionName, m)
		}
		if m.Method == "" {
			return fmt.Errorf("%w: %v", ErrEmptyMethodExecutionMethod, m)
		}
	}
	return nil
}

func (entityRules *EntityRules) applyOnOne(entity interface{}) error {
	re := reflect.ValueOf(entity)
	for _, rule := range entityRules.Rules {
		isTriggered, err := entityRules.isRuleTriggered(re, rule.When)
		if err != nil {
			return fmt.Errorf("errors occured in triggers check: %w", err)
		}
		if isTriggered {
			hasConditions, err := entityRules.entityHasConditions(re, rule.If)
			if err != nil {
				return fmt.Errorf("errors occured in conditions check: %w", err)
			}
			if hasConditions {
				if err := entityRules.execute(re, rule.Do); err != nil {
					return fmt.Errorf("errors occured in actions execution: %w", err)
				}
			}
		}
	}
	return nil
}

func (entityRules *EntityRules) isRuleTriggered(entity reflect.Value, when []string) (bool, error) {
	//Check if at least one trigger is valid using When attribute of the current rule
	if when[0] == TrueValue {
		return true, nil
	}
	for _, name := range when {
		isTriggered, err := entityRules.executeMethodBool(entity, name)
		if err != nil {
			return false, err
		}
		if isTriggered {
			return true, nil
		}
	}
	return false, nil
}

func (entityRules *EntityRules) entityHasConditions(entity reflect.Value, on string) (hasConditions bool, err error) {
	//Check if conditions are valid using If attribute of the current rule
	if on == TrueValue { //If no conditions are requiered
		return true, nil
	}
	//Formating On for latter evaluation expression
	namedConditions := make(map[string]bool)
	matches := functionMatcher.FindAllString(on, -1)
	for _, match := range matches {
		if !isOperator(match) {
			on = strings.Replace(on, match, match+" == true", -1)
			if namedConditions[match], err = entityRules.executeMethodBool(entity, match); err != nil {
				return false, err
			}
		}
	}
	for strings.Contains(on, "== true == true") {
		on = strings.Replace(on, "== true == true", "== true", -1)
	}
	//Evaluating expression with bexpr package
	eval, err := bexpr.CreateEvaluator(on)
	if err != nil {
		return false, err
	}
	if hasConditions, err = eval.Evaluate(namedConditions); err != nil {
		return false, err
	}
	return hasConditions, nil
}

func isOperator(s string) bool {
	return s == "or" || s == "and" || s == "not"
}

func (entityRules *EntityRules) execute(entity reflect.Value, do []string) error {
	//Execute actions from Do attributes in the current rule
	for _, name := range do {
		if _, err := executeMethod(entity, entityRules.Actions, name); err != nil {
			return err
		}
	}
	return nil
}

func (entityRules *EntityRules) executeMethodBool(entity reflect.Value, name string) (bool, error) {
	//Evaluating with entity the method while checking if it return a boolean
	resultMethod, err := executeMethod(entity, entityRules.Conditions, name)
	if err != nil {
		return false, err
	}
	if len(resultMethod) > 0 && resultMethod[0].Interface() != nil {
		if result, ok := resultMethod[0].Interface().(bool); ok {
			return result, nil
		}
	}
	return false, ErrMethodReturnNotBool
}

func executeMethod(entity reflect.Value, methods MethodsExecution, name string) (result []reflect.Value, err error) {
	//Evaluating with entity the method identified by name and contained in the input methods map
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case string:
				result, err = nil, errors.New(r.(string)+": "+name)
			case *reflect.ValueError:
				re := r.(*reflect.ValueError)
				result, err = nil, errors.New(re.Error()+": "+name)
			}
		}
	}()
	method := methods.GetByName(name)
	if method.Name == "" {
		return nil, fmt.Errorf("%w: %s", ErrMethodNotRegistered, name)
	}
	rps, err := evaluateParams(entity, method)
	if err != nil {
		return nil, err
	}
	return entity.MethodByName(method.Method).Call(rps), nil

}

func evaluateParams(entity reflect.Value, methodExecution MethodExecution) (rps []reflect.Value, err error) {
	for _, param := range methodExecution.Arguments {
		if param, err = evaluateParam(entity, param); err != nil {
			return nil, err
		}
		rps = append(rps, reflect.ValueOf(param))
	}
	return rps, nil
}

func evaluateParam(entity interface{}, param interface{}) (interface{}, error) {
	if tmpl, ok := param.(*template.Template); ok {
		var evalParam bytes.Buffer
		if err := tmpl.Execute(&evalParam, entity); err != nil {
			return nil, err
		}
		param = evalParam.String()
	}
	return param, nil
}

func (mm MethodsExecution) GetByName(name string) MethodExecution {
	for _, m := range mm {
		if m.Name == name {
			return m
		}
	}
	return MethodExecution{}
}
