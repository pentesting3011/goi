package goi

import (
	"fmt"
	"reflect"
)

type NumberValidator struct {
	Validator
	min float64
	max float64
}

var validNums = []string{"int8", "int16", "int32", "int64", "int", "float32", "float64", "uint8", "uint16", "uint32", "uint64", "uint"}

func Number() *NumberValidator {
	v := &NumberValidator{}
	v.NumberBase()
	return v
}

func (s *NumberValidator) NumberBase() {
	s.ruleNames = append(s.ruleNames, "base")
	s.required = false
	s.rules = append(s.rules, func(value *any) error {
		if s.originalValue == nil || *s.originalValue == nil {
			return nil
		}
		if reflect.TypeOf(*s.originalValue).Kind() != reflect.Float32 && reflect.TypeOf(*s.originalValue).Kind() != reflect.Float64 {
			return fmt.Errorf("%s not a number", s.label)
		}
		return nil
	})
}

func (s *NumberValidator) Required() *NumberValidator {
	s.ruleNames = append(s.ruleNames, "required")
	s.required = true
	s.rules = append(s.rules, func(value *any) error {
		if s.required && (s.originalValue == nil || *s.originalValue == nil) {
			return fmt.Errorf("%s must be defined", s.label)
		}
		return nil
	})
	return s
}

func (s *NumberValidator) Min(length float64) *NumberValidator {
	s.ruleNames = append(s.ruleNames, "min")
	s.min = length
	s.rules = append(s.rules, func(value *any) error {
		if *s.originalValue != nil {
			cast, _ := (*s.originalValue).(float64)
			if cast < float64(s.min) {
				return fmt.Errorf("%s must be greater than or equal to %.0f", s.label, s.min)
			}
		}
		return nil
	})
	return s
}

func (s *NumberValidator) Max(length float64) *NumberValidator {
	s.ruleNames = append(s.ruleNames, "max")
	s.max = length
	s.rules = append(s.rules, func(value *any) error {
		if *s.originalValue != nil {
			cast, _ := (*s.originalValue).(float64)
			if cast > float64(s.max) {
				return fmt.Errorf("%s must be less than %.0f", s.label, s.max)
			}
		}
		return nil
	})
	return s
}

// Sets a default value if the original value is undefined
func (s *NumberValidator) Default(defaultValue any) *NumberValidator {
	s.ruleNames = append(s.ruleNames, "default")
	s.defaultValue = defaultValue
	return s
}

func (s *NumberValidator) Optional() *NumberValidator {
	s.ruleNames = append(s.ruleNames, "optional")
	s.required = false
	return s
}

func (s *NumberValidator) Label(label string) *NumberValidator {
	s.label = label
	return s
}

func (s *NumberValidator) Valid(value []any) *NumberValidator {
	s.ruleNames = append(s.ruleNames, "valid")
	for _, v := range value {
		typeVar := reflect.TypeOf(v).Kind()
		if findIndex(validNums, typeVar.String()) == -1 {
			panic("Can not coerce type " + typeVar.String() + " to any kind of number")
		}
	}
	s.rules = append(s.rules, func(v *any) error {
		//Allowed cause required validate this
		if *s.originalValue == nil {
			return nil
		}
		ok := false
		for _, v := range value {
			vType := reflect.TypeOf(v)
			valueOfOriginalValue := reflect.ValueOf(*s.originalValue)
			canCast := valueOfOriginalValue.CanConvert(vType)
			if !canCast {
				continue
			}
			if v == valueOfOriginalValue.Convert(vType).Interface() {
				ok = true
			}
		}
		if !ok {
			return fmt.Errorf("%s not in valid array", s.label)
		}
		return nil
	})
	return s
}

// You must pass the pointer to the value you want to validate
func (s *NumberValidator) Validate(data *any) error {
	if s.label == "" {
		s.label = "value"
	}
	s.originalValue = new(any)
	*s.originalValue = nil
	if data != nil && *data != nil {
		*s.originalValue = *data
	} else {
		*data = s.defaultValue
	}
	for _, rule := range s.rules {
		if err := rule(data); err != nil {
			return err
		}
	}
	return nil
}
