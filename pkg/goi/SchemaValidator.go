package goi

import (
	"fmt"
	"reflect"
)

type SchemaValidator struct {
	schemaRuleNamesIn []string
	schemasIn         []any

	schemaValueNames []string
	values           []any
}

func Schema(m map[string]any) *SchemaValidator {
	s := &SchemaValidator{}
	for k, v := range m {
		s.schemaRuleNamesIn = append(s.schemaRuleNamesIn, k)
		s.schemasIn = append(s.schemasIn, v)
	}
	return s
}

func (s *SchemaValidator) Validate(data *any) error {
	dataVal := reflect.ValueOf(*data)
	if dataVal.Kind() != reflect.Map {
		return fmt.Errorf("not a map")
	}

	for _, k := range dataVal.MapKeys() {
		s.schemaValueNames = append(s.schemaValueNames, k.String())
		v := dataVal.MapIndex(k).Interface()
		s.values = append(s.values, v)
	}

	for i, v := range s.schemasIn {
		//Nombre de schema de entrada
		fieldName := s.schemaRuleNamesIn[i]
		//Indice en los schema que se reciben
		index := findIndex(s.schemaValueNames, fieldName)
		var var1 any
		label := fieldName
		if index == -1 {
			var1 = nil
		} else {
			var1 = (any)(s.values[index])
		}
		//Valor que recibe en validate
		labelMethod := reflect.ValueOf(v).MethodByName("Label")
		if labelMethod.IsValid() {
			labelMethod.Call([]reflect.Value{reflect.ValueOf(label)})
		}
		valMethod := reflect.ValueOf(v).MethodByName("Validate")
		if !valMethod.IsValid() {
			panic("Method validate not found")
		}
		result := valMethod.Call([]reflect.Value{reflect.ValueOf(&var1)})
		if len(result) > 0 && result[0].Interface() != nil {
			return result[0].Interface().(error)
		}
		reflect.ValueOf(*data).SetMapIndex(reflect.ValueOf(fieldName), reflect.ValueOf(var1))
	}

	return nil
}
