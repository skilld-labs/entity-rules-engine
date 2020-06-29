package entityrules

import (
	"encoding/json"
	"fmt"
	"github.com/knadh/koanf"
	cjson "github.com/knadh/koanf/parsers/json"
	cyaml "github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"reflect"
	"text/template"
)

func LoadFromJSON(filePath string, opts ...LoadOption) (entityRules *EntityRules, err error) {
	if entityRules, err = loadFromFile(filePath, "json", opts...); err != nil {
		return nil, err
	}
	return entityRules, nil
}

func LoadFromYAML(filePath string, opts ...LoadOption) (entityRules *EntityRules, err error) {
	if entityRules, err = loadFromFile(filePath, "yaml", opts...); err != nil {
		return nil, err
	}
	return entityRules, nil
}

func LoadFromMap(m map[string]interface{}, opts ...LoadOption) (*EntityRules, error) {
	em, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	var entityRules EntityRules
	if err := json.Unmarshal(em, &entityRules); err != nil {
		return nil, err
	}
	if opts != nil {
		fm := make(template.FuncMap)
		optFm := &LoadOptions{FuncMap: fm}
		for _, opt := range opts {
			opt(optFm)
		}
		if err := (&entityRules).parse(fm); err != nil {
			return nil, err
		}
	}
	return &entityRules, nil
}

func LoadFromInterface(er interface{}, opts ...LoadOption) (*EntityRules, error) {
	m, ok := er.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("cannot load EntityRules, cannot parse interface into map[string]interface{}")
	}
	return LoadFromMap(m, opts...)
}

func loadFromFile(filePath string, fileType string, opts ...LoadOption) (*EntityRules, error) {
	k := koanf.New(".")
	switch fileType {
	case "yaml":
		if err := k.Load(file.Provider(filePath), cyaml.Parser()); err != nil {
			return nil, err
		}
	case "json":
		if err := k.Load(file.Provider(filePath), cjson.Parser()); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("cannot read file type: %s", fileType)
	}
	return LoadFromInterface(k.Get("entityRules"), opts...)
}

func (entityRules *EntityRules) parse(fm template.FuncMap) error {
	if err := entityRules.Actions.parse(fm); err != nil {
		return err
	}
	if err := entityRules.Conditions.parse(fm); err != nil {
		return err
	}
	return nil
}

func (me *MethodsExecution) parse(fm template.FuncMap) (err error) {
	for _, a := range *me {
		for argIdx, arg := range a.Arguments {
			if _, isString := arg.(string); isString {
				a.Arguments[argIdx], err = template.New("param").Funcs(fm).Parse(arg.(string))
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

type LoadOption func(*LoadOptions)

type LoadOptions struct {
	FuncMap template.FuncMap
}

func WithFuncMap(funcMap template.FuncMap) LoadOption {
	return func(opts *LoadOptions) {
		fm := opts.FuncMap
		for funcName, function := range funcMap {
			fm[funcName] = function
		}
	}
}

func WithEntityFuncs(entity interface{}) LoadOption {
	return func(opts *LoadOptions) {
		fm := opts.FuncMap
		loadFuncMap(entity, fm)
	}
}

func loadFuncMap(entity interface{}, fm template.FuncMap) {
	rte := reflect.TypeOf(entity)
	for methodIdx := 0; methodIdx < rte.NumMethod(); methodIdx++ {
		method := rte.Method(methodIdx)
		fm[method.Name] = func(entity interface{}, params ...interface{}) (interface{}, error) {
			var rps []reflect.Value
			for _, p := range params {
				rps = append(rps, reflect.ValueOf(p))
			}
			responses := reflect.ValueOf(entity).MethodByName(method.Name).Call(rps)
			var resp interface{}
			var err error
			if len(responses) > 1 {
				var lastParamIsError bool
				var isError bool
				v := responses[len(responses)-1].Interface()
				err, isError = responses[len(responses)-1].Interface().(error)
				if isError {
					lastParamIsError = true
				} else {
					lastParamIsError = v == nil
				}
				if lastParamIsError {
					responses = responses[:len(responses)-1]
				}
			}
			if len(responses) == 1 {
				resp = responses[0]
			} else if len(responses) > 1 {
				resp = responses
			}
			return resp, err
		}
	}
}
