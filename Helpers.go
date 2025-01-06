package goi

import "fmt"

type Helper struct {
	tag string
}

func (h *Helper) Error(tag string) error {
	h.tag = tag
	if h.tag != "" {
		return fmt.Errorf("goi error: %v", h.tag)
	}
	return nil
}
