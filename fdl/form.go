package fdl

import "fmt"

type FormCode string

type Form struct {
	Code           FormCode
	Description    string
	BackgroundFile string
	Nodes          []Node
}

func (f Form) Valid() error {
	if f.Code == "" {
		return fmt.Errorf("invalid form: code must be provided")
	}
	if f.Description == "" {
		return fmt.Errorf("invalid form: description must be provided")
	}
	for _, node := range f.Nodes {
		if err := node.Valid(); err != nil {
			return fmt.Errorf("invalid form: %w", err)
		}
	}
	return nil
}
