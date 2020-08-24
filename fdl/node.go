package fdl

type Node interface {
	// Type returns the Go type of the field for the form definition.
	Type() string

	// Types returns a rendered string of any supporting types required for the
	// form to work correctly, e.g. constants for option node variants.
	Types() string

	// Render returns a string suitable for including in a frl form's BuildPDF
	// method.
	Render() string

	// Valid validates the the node contents and returns an error if any.
	Valid() error
}
