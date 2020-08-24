package fdl

import "fmt"

type TextNode struct {
	Name        string
	Description string
	Location    Location
	FontSize    int
}

func (n TextNode) Type() string {
	return "string"
}

func (n TextNode) Types() string {
	return ""
}

func (n TextNode) Render() string {
	return fmt.Sprintf(
		"r.PrintText(\nfrl.Location{\nX: %d,\nY: %d,\n},\nfrl.FontSize(%d),\nf.%s,\n)",
		n.Location.X,
		n.Location.Y,
		n.FontSize,
		n.Name,
	)
}

func (n TextNode) Valid() error {
	if n.Name == "" {
		return fmt.Errorf("invalid text node: name must be provided")
	}
	if n.Description == "" {
		return fmt.Errorf("invalid text node: description must be provided")
	}
	if n.FontSize == 0 {
		return fmt.Errorf("invalid text node: font size must be greater than zero")
	}
	return nil
}
