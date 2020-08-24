package fdl

import "fmt"

type TextNode struct {
	Name     string
	Location Location
	FontSize int
}

func (n TextNode) Type() string {
	return "string"
}

func (n TextNode) Types() string {
	return ""
}

func (n TextNode) Render() string {
	return fmt.Sprintf(
		"r.PrintText(frl.Location{X: %d, Y: %d}, frl.FontSize(%d), n.%s)",
		n.Location.X,
		n.Location.Y,
		n.FontSize,
		n.Name,
	)
}
