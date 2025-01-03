package goi

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"unicode/utf8"
)

type stringValidator struct {
	Validator
	min   int
	max   int
	regex *regexp.Regexp
}

func String() *stringValidator {
	v := &stringValidator{}
	v.StringBase()
	return v
}

func (s *stringValidator) StringBase() {
	s.ruleNames = append(s.ruleNames, "base")
	s.required = false
	s.rules = append(s.rules, func(value *any) error {
		if s.originalValue == nil || *s.originalValue == nil {
			return nil
		}
		if reflect.TypeOf(*s.originalValue).Kind() != reflect.String {
			return fmt.Errorf("%s not a string", s.label)
		}
		return nil
	})
}

func (s *stringValidator) Required() *stringValidator {
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

func (s *stringValidator) LowerCase() *stringValidator {
	s.ruleNames = append(s.ruleNames, "lowercase")
	s.rules = append(s.rules, func(value *any) error {
		if *value != nil && *s.originalValue != nil {
			val := (any)(strings.ToLower((*value).(string)))
			*value = val
		}
		return nil
	})
	return s
}

func (s *stringValidator) Trim() *stringValidator {
	s.ruleNames = append(s.ruleNames, "trim")
	s.rules = append(s.rules, func(value *any) error {
		if *value != nil && *s.originalValue != nil {
			val := (any)(strings.TrimSpace((*value).(string)))
			*value = val
		}
		return nil
	})
	return s
}

func (s *stringValidator) Regex(regex string) *stringValidator {
	s.ruleNames = append(s.ruleNames, "regex")
	s.regex = regexp.MustCompile(regex)
	s.rules = append(s.rules, func(value *any) error {
		if *value != nil && !s.regex.MatchString((*s.originalValue).(string)) {
			return fmt.Errorf("%s does not match regex", s.label)
		}
		return nil
	})
	return s
}

func (s *stringValidator) Min(length int) *stringValidator {
	s.ruleNames = append(s.ruleNames, "min")
	s.min = length
	s.rules = append(s.rules, func(value *any) error {
		if *s.originalValue != nil {
			cast, _ := (*s.originalValue).(string)
			if utf8.RuneCountInString(cast) < s.min {
				return fmt.Errorf("%s must be at least %d length", s.label, s.min)
			}
		}
		return nil
	})
	return s
}

func (s *stringValidator) Max(length int) *stringValidator {
	s.ruleNames = append(s.ruleNames, "max")
	s.max = length
	s.rules = append(s.rules, func(value *any) error {
		if *s.originalValue != nil {
			cast, _ := (*s.originalValue).(string)
			if utf8.RuneCountInString(cast) > s.max {
				return fmt.Errorf("%s must be least than %d length", s.label, length)
			}
		}
		return nil
	})
	return s
}

func (s *stringValidator) Alphanum() *stringValidator {
	s.ruleNames = append(s.ruleNames, "alphanum")
	s.rules = append(s.rules, func(value *any) error {
		rgx := regexp.MustCompile(`^[a-zA-Z0-9]*$`)
		if *s.originalValue != nil && !rgx.MatchString((*s.originalValue).(string)) {
			return fmt.Errorf("%s not alphanum", s.label)
		}
		return nil
	})
	return s
}

// Sets a default value if the original value is undefined
func (s *stringValidator) Default(defaultValue any) *stringValidator {
	s.ruleNames = append(s.ruleNames, "default")
	s.defaultValue = defaultValue
	return s
}

func (s *stringValidator) Optional() *stringValidator {
	s.ruleNames = append(s.ruleNames, "optional")
	s.required = false
	return s
}

func (s *stringValidator) Label(label string) *stringValidator {
	s.label = label
	return s
}

func (s *stringValidator) Valid(value []any) *stringValidator {
	s.ruleNames = append(s.ruleNames, "valid")
	for _, v := range value {
		typeVar := reflect.TypeOf(v).Kind()
		if typeVar != reflect.String {
			panic("Can not coerce type " + typeVar.String() + " to string")
		}
	}
	s.rules = append(s.rules, func(v *any) error {
		//Allowed cause required validate this
		if *s.originalValue == nil {
			return nil
		}
		ok := false
		for _, v := range value {
			if v == *s.originalValue {
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

func (s *stringValidator) Custom(customFunction func(value *any, helpers *Helper) any) *stringValidator {
	s.ruleNames = append(s.ruleNames, "custom")
	s.rules = append(s.rules, func(v *any) error {
		h := &Helper{}
		result := customFunction(v, h)
		fmt.Println(reflect.TypeOf(result))
		fmt.Println(reflect.TypeOf(fmt.Errorf("")))
		if reflect.TypeOf(result) == reflect.TypeOf(fmt.Errorf("")) {
			return (result).(error)
		} else {
			*v = result
		}
		return nil
	})
	return s
}

// You must pass the pointer to the value you want to validate
func (s *stringValidator) Validate(data *any) error {
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
