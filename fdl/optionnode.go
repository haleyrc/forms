package fdl

import (
	"fmt"
	"strings"
)

type Option struct {
	Name     string
	Location Location
}

type OptionNode struct {
	Name    string
	Size    int
	Options []Option
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
		"r.PrintCheck(f.%s.Location(), frl.Size(%d))",
		n.Name,
		n.Size,
	)
}
