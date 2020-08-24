package fdl

import "fmt"

func NewIterator(r *Registry) *Iterator {
	forms := make([]*Form, 0, len(r.forms))
	for _, form := range r.forms {
		forms = append(forms, &form)
	}
	return &Iterator{
		curr:  0,
		forms: forms,
	}
}

type Iterator struct {
	curr  int
	forms []*Form
}

func (i *Iterator) Next() *Form {
	if i.curr >= len(i.forms) {
		return nil
	}
	f := i.forms[i.curr]
	i.curr++
	return f
}

func NewRegistry() *Registry {
	return &Registry{
		forms: make(map[FormCode]Form),
	}
}

type Registry struct {
	forms map[FormCode]Form
}

func (fr *Registry) Iter() *Iterator {
	return NewIterator(fr)
}

func (fr *Registry) Register(f Form) error {
	if _, found := fr.forms[f.Code]; found {
		return fmt.Errorf("duplicate form code: %s", f.Code)
	}
	fr.forms[f.Code] = f
	return nil
}
