package goi

type Validatable interface {
	Validate(data *any) error
}

type Validator struct {
	defaultValue  any
	originalValue *any
	rules         []func(*any) error
	ruleNames     []string
	label         string
	Validator     Validatable

	required bool
}

func findIndex(arr []string, el string) int {
	for i, v := range arr {
		if el == v {
			return i
		}
	}
	return -1
}
