## `fdl`

`fdl` (pronounced "fiddle") contains primitives for describing the shape of form
fields. Forms are created as entries in a global form map. These are then
converted to an input type and an associated build method using the Form Render
Library (`frl`) and `go generate`.

The ultimate goal of `fdl` is to create a rich type system for describing forms
that can be used to describe all the inputs and variations of every Frazer form.

### Example

The following example demonstrates the basic document structure and a number of
basic primitives. Additional functionality can be added to the Node interface to
support additional configuration.

It's important to note that development of the primitives is related but largely
orthogonal to the development of form definitions. A developer modifying and
creating primitives (core) needs to know the requirements from the developers
writing form definitions (forms programmers) as well as the internals of the
rendering and generation processes. By contrast, the forms programmers need only
be familiar with the `fdl` library of primitives (which basically means knowing
all of the node types and their usages).

```go
package fdl

type Document struct {
    Code           string
    BackgroundFile string
    Nodes          []Node
}

type Node interface {
    // Type returns the Go type of the field for the form definition.
    Type() string

    // Types returns a rendered string of any supporting types required for the
    // form to work correctly, e.g. constants for option node variants.
    Types()  string

    // Render returns a string suitable for including in a frl form's BuildPDF
    // method.
    Render() string
}

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
        "f.renderer.PrintText(frl.Location{X: %d, Y: %d}, frl.FontSize(%d), n.%s)",
        n.Location.X,
        n.Location.Y,
        n.FontSize,
        n.Name,
    )
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
    fmt.Fprintf(sb, "type %sOption string\n\n", n.Name)
    fmt.Fprintf(sb, "const (\n")
    for _, opt := range n.Options {
        fmt.Fprintf(sb,
            "        %sOption%s = \"%s\"\n",
            n.Name,
            opt.Name,
            opt.Name,
        )
        fmt.Fprintf(sb, ")\n")
    }
    return sb.String()
}

func (n OptionNode) Render() string {
    return fmt.Sprintf(
        "f.renderer.PrintCheck(f.%s.Location(), frl.Size(%d))",
        n.Name,
        n.Size,
    )
}
```

Given these primitives, we can construct the toy form below. As you can see, we
are creating a declarative framework for forms much like modern web frameworks
provide a declarative definition of user interfaces. For some trial output from
generation of this toy example, see the section on `forms`.

```go
type FormRegistry map[string]Document

func (fr *FormRegistry) Register(d Document) error {
    if _, found := fr[d.Code]; found {
        return ErrDuplicateFormCode
    }
    fr[d.Code] = d
}

var globalForms = map[string]Document{}

func init() {
    panicOnError(globalForms.Register(Document {
        Code: "1234",
        BackgroundFile: "form_1234.pdf",
        Nodes: []Node{
            TextNode{
                Name:     "FirstName",
                Location: Location{X: 10, Y: 20},
                FontSize: 12,
            },
            OptionsNode{
                Name: "Style",
                Size: 10,
                Options: []Option{
                    Option{
                        Location: Location{X: 10, Y: 30},
                        Name: "Sedan",
                    },
                    Option{
                        Location: Location{ X: 30, Y: 30},
                        Name: "Truck",
                    },
                },
            },
        },
    }))
)

```

## `frl`

`frl` (pronounced "furl") contains primitives for rendering nodes onto a form.
Rendering primitives include interfaces for drawing text, checking boxes, etc.

The focus of `frl` is to provide a common interface for printing that the `fdl`
primitives can refer to under the hood for defining their rendering behavior.
In order for the process to work, we export the interface that all renderers
support as well as concrete implementations of renderers that we can then
develop in parallel with other efforts.

In this way, we separate the three parts of the forms printing process into
orthogonal concerns:

- Form description
  - The process of describing a form in terms of its properties and fields.
- Form compilation
  - The process of generating a standard form interface for each definition that
    can be used in production.
- Form printing
  - The process of converting form input data into a rendered form by utilizing
    the standard form interface and a specific rendering implementation.

This design also increases the long-term utility of the architecture by allowing
us to modify and update the generation and rendering processes independently of
the definitions. As long as our description library is rich enough to express
all field permutations, the definitions are never at risk of going stale, even
in the face of a different conversion process during generation, or an entirely
new rendering implementation.

### Example

The following minimal example contains enough code to support the generation of
our toy form. Additional functionality can be added to support new use-cases as
long as the exported concrete implementations are also updated to include the
new behaviors.

Note that external packages can also provide their own `Renderer` implementation
to support unexpected use-cases, since our generated forms will only contain a
reference to the `Renderer` interface and not a specific implementation.

```go
package frl

type Location struct {
    X int
    Y int
}

type FontSize int
type Size int

type Renderer interface {
    PrintCheck(Location, Size) error
    PrintText(Location, FontSize, string) error
}
```

## `form`

```go
type StyleOption string

const (
    StyleOptionSedan StyleOption = "Sedan"
    StyleOptionTruck StyleOption = "Truck"
)

func (so StyleOption) Location() fdl.Location {
    switch so {
    case StyleOptionSedan:
        return fdl.Location{X: 10, Y: 30}
    case StyleOptionTruck:
        return fdl.Location{X: 30, Y: 30}
    }
    panic(fmt.Errorf("generated: invalid StyleOption: %s", so))
}

type Form1234 struct {
    FirstName string
    Style     StyleOption
}

func (f Form1234) BuildPDF(r frl.Renderer) error {
    if err := r.PrintText(frl.Location{X: 10, Y: 20}, frl.FontSize(12), n.FirstName); err != nil {
        return err
    }
    if err := r.PrintCheck(f.Style.Location(), frl.Size(10)); err != nil {
        return err
    }
    return nil
}
```

TODO: Add the generation code for the Location method to the fdl section
TODO: Add comments in the generated code explaining how the lines were derived
