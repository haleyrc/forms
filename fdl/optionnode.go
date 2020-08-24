package fdl

import (
	"fmt"
	"strings"
)

type Option struct {
	Name     string
	Location Location
}

func (o Option) Valid() error {
	if o.Name == "" {
		return fmt.Errorf("invalid option: name must be provided")
	}
	return nil
}

type OptionNode struct {
	Name        string
	Description string
	Size        int
	Options     []Option
}

func (n OptionNode) Valid() error {
	if n.Name == "" {
		return fmt.Errorf("invalid option node: name must be provided")
	}
	if n.Description == "" {
		return fmt.Errorf("invalid option node: description must be provided")
	}
	if n.Size == 0 {
		return fmt.Errorf("invalid option node: size must be greater than zero")
	}
	for _, opt := range n.Options {
		if err := opt.Valid(); err != nil {
			return fmt.Errorf("invalid option node: %w", err)
		}
	}
	return nil
}

func (n OptionNode) Type() string {
	return "StyleOption"
}

func (n OptionNode) Types() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "type %sOption string\n", n.Name)
	fmt.Fprintln(&sb)

	fmt.Fprintf(&sb, "const (\n")
	for _, opt := range n.Options {
		fmt.Fprintf(&sb,
			"        %sOption%s = %q\n",
			n.Name,
			opt.Name,
			opt.Name,
		)
	}
	fmt.Fprintf(&sb, ")\n")
	fmt.Fprintln(&sb)

	fmt.Fprintf(&sb, "func (o %sOption) Location() fdl.Location {\n", n.Name)
	fmt.Fprintln(&sb, "        switch o {")
	for _, opt := range n.Options {
		fmt.Fprintf(&sb,
			"        case %sOption%s:\n",
			n.Name,
			opt.Name,
		)
		fmt.Fprintf(&sb,
			"                return fdl.Location{X: %d, Y: %d}\n",
			opt.Location.X,
			opt.Location.Y,
		)
	}
	fmt.Fprintln(&sb, "        }")
	fmt.Fprintf(&sb,
		"        panic(fmt.Errorf(\"generated: invalid %sOption: %%s\", o))\n",
		n.Name,
	)
	fmt.Fprintln(&sb, "}")
	fmt.Fprintln(&sb)

	return sb.String()
}

func (n OptionNode) Render() string {
	return fmt.Sprintf(
		"r.PrintCheck(\nfrl.Location{\nX: f.%s.Location().X,\nY: f.%s.Location().Y,\n}, \nfrl.Size(%d),\n)",
		n.Name,
		n.Name,
		n.Size,
	)
}
