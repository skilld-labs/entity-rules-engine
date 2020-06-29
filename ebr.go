// Package the rules logic
// The exported function ApplyOn check if requiered conditions are valid and execute actions from entityRules object

package entityrules

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/hashicorp/go-bexpr"
	"reflect"
	"regexp"
	"strings"
	"text/template"
)

type EntityRules struct {
	Conditions *MethodsExecution
	Actions    *MethodsExecution
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

var (
	ErrNotRegistered = errors.New("is not registered")
	ErrIsNotBool     = errors.New("do not return boolean")
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
		fmt.Println("Entity does not have rules")
		return nil
	} else {
		for _, r := range rr {
			if r.Name == "" {
				return fmt.Errorf("Rule missing name: %v", r)
			}
			if r.When == nil {
				return fmt.Errorf("rule: %s does not have When, if you do not want any when conditions set it to true", r.Name)
			}
			if r.If == "" {
				return fmt.Errorf("rule: %s does not have If, if you do not want any if conditions set it to true", r.Name)
			}
			if r.Do == nil {
				return fmt.Errorf("rule: %s does not have Do", r.Name)
			}
		}
	}
	return nil
}

func (mm MethodsExecution) validate() error {
	for _, m := range mm {
		if m.Name == "" {
			return fmt.Errorf("cannot use empty name: %v", m)
		}
		if m.Method == "" {
			return fmt.Errorf("cannot use empty method: %v", m)
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
	if when[0] == "true" {
		return true, nil
	}
	for _, name := range when {
		isTriggered, err := entityRules.executeMethodBool(entity, name)
		if err != nil {
			return false, fmt.Errorf("%w: %s", err, name)
		}
		if isTriggered {
			return true, nil
		}
	}
	return false, nil
}

func (entityRules *EntityRules) entityHasConditions(entity reflect.Value, on string) (hasConditions bool, err error) {
	//Check if conditions are valid using If attribute of the current rule
	if on == "true" { //If no conditions are requiered
		return true, nil
	}
	//Formating On for latter evaluation expression
	namedConditions := make(map[string]bool)
	condition := on
	var functionMatcher = regexp.MustCompile(`[\w\d]+`)
	matches := functionMatcher.FindAllString(condition, -1)
	for _, match := range matches {
		if isNotOperator(match) {
			condition = strings.Replace(condition, match, match+" == true", -1)
			if namedConditions[match], err = entityRules.executeMethodBool(entity, match); err != nil {
				return false, fmt.Errorf("%w: %s", err, match)
			}
		}
	}
	for strings.Contains(condition, "== true == true") {
		condition = strings.Replace(condition, "== true == true", "== true", -1)
	}
	//Evaluating expression with bexpr package
	eval, err := bexpr.CreateEvaluatorForType(condition, nil, (map[string](bool))(nil))
	if err != nil {
		return false, err
	}
	hasConditions, err = eval.Evaluate(namedConditions)
	if err != nil {
		return false, err
	}
	return hasConditions, nil
}

func isNotOperator(s string) bool {
	return s != "or" && s != "and" && s != "not"
}

func (entityRules *EntityRules) execute(entity reflect.Value, do []string) error {
	//Execute actions from Do attributes in the current rule
	for _, name := range do {
		if _, err := executeMethod(entity, *entityRules.Actions, name); err != nil {
			return fmt.Errorf("%w: %s", err, name)
		}
	}
	return nil
}

func (entityRules *EntityRules) executeMethodBool(entity reflect.Value, name string) (result bool, err error) {
	//Evaluating with entity the method while checking if it return a boolean
	/*	defer func() {
		if r := recover(); r != nil {
			result, err = false, r.(error)
		}
	}()*/
	resultMethod, err := executeMethod(entity, *entityRules.Conditions, name)
	if err != nil {
		return false, err
	}
	result, ok := resultMethod[0].Interface().(bool)
	if !ok {
		return false, ErrIsNotBool
	}
	return result, nil
}

func executeMethod(entity reflect.Value, methods MethodsExecution, name string) (result []reflect.Value, err error) {
	//Evaluating with entity the method identified by name and contained in the input methods map
	defer func() {
		if r := recover(); r != nil {
			result, err = nil, r.(error)
		}
	}()
	method := methods.GetByName(name)
	if method.Name == "" {
		return nil, ErrNotRegistered
	}
	rps, err := evaluateParams(entity, method)
	if err != nil {
		return nil, err
	}
	return entity.MethodByName(method.Method).Call(rps), nil
}

func evaluateParams(entity reflect.Value, methodExecution MethodExecution) (rps []reflect.Value, err error) {
	for _, param := range methodExecution.Arguments {
		param, err = evaluateParam(entity, param)
		if err != nil {
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

func (mm *MethodsExecution) GetByName(name string) MethodExecution {
	for _, m := range *mm {
		if m.Name == name {
			return m
		}
	}
	return MethodExecution{}
}
