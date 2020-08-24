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
    fmt.Fprintf(sb, "type %sOption string\n", n.Name)
    fmt.Fprintln(sb)

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
    fmt.Fprintln(sb)

    fmt.Fprintf(sb, "func (o %sOption) Location() fdl.Location {\n", n.Name)
    fmt.Fprintln(sb, "        switch o {")
    for _, opt := range n.Options {
        fmt.Fprintf(sb,
            "        case %sOption%s:\n",
            n.Name,
            opt.Name,
        )
        fmt.Fprintf(sb,
            "                return fdl.Location{X: %d, Y: %d}\n",
            opt.Location.X,
            opt.Location.Y,
        )
    }
    fmt.Fprintln(sb, "        }")
    fmt.Fprintf(sb, "        panic(fmt.Errorf(\"generated: invalid %sOption: \%s\", o))\n")
    fmt.Fprintln(sb, "}")
    fmt.Fprintln(sb)

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

## `package form`

The `form` package is the end-result of all of the compilation steps and is the
import that would be used in a live production system. Here is our toy form in
its compiled state:

```go
// File: form1234.go
package form

// The following type and method definitions were generated from the
// OptionsNode's Types() method.

// This line derives from the following:
//
//   fmt.Fprintf(sb, "type %sOption string\n", n.Name)
/
// where in this case, n.Name == "Style".
type StyleOption string

// The following lines are generated from the snippet:
//
//    fmt.Fprintf(sb, "const (\n")
//    for _, opt := range n.Options {
//        fmt.Fprintf(sb,
//            "        %sOption%s = \"%s\"\n",
//            n.Name,
//            opt.Name,
//            opt.Name,
//        )
//        fmt.Fprintf(sb, ")\n")
//    }
//
const (
    StyleOptionSedan StyleOption = "Sedan"
    StyleOptionTruck StyleOption = "Truck"
)

// This method for determing where a check needs to be placed based on the
// selected style was generated from:
//
//    fmt.Fprintf(sb, "func (o %sOption) Location() fdl.Location {\n", n.Name)
//    fmt.Fprintln(sb, "        switch o {")
//    for _, opt := range n.Options {
//        fmt.Fprintf(sb,
//            "        case %sOption%s:\n",
//            n.Name,
//            opt.Name,
//        )
//        fmt.Fprintf(sb,
//            "                return fdl.Location{X: %d, Y: %d}\n",
//            opt.Location.X,
//            opt.Location.Y,
//        )
//    }
//    fmt.Fprintln(sb, "        }")
//    fmt.Fprintf(sb, "        panic(fmt.Errorf(\"generated: invalid %sOption: \%s\", o))\n")
//    fmt.Fprintln(sb, "}")
//
func (so StyleOption) Location() fdl.Location {
    switch so {
    case StyleOptionSedan:
        return fdl.Location{X: 10, Y: 30}
    case StyleOptionTruck:
        return fdl.Location{X: 30, Y: 30}
    }
    panic(fmt.Errorf("generated: invalid StyleOption: %s", so))
}

// The FormXYZ type is generated directly by concatenating the form code to
// "Form" and then enumerating the nodes one by one using their Name() and
// Type() values, e.g.:
//
//   fmt.Printf("type Form%s struct{\n", f.Code)
//   for _, n := f.Nodes {
//       fmt.Printf("        %s %s\n", n.Name(), n.Type())
//   }
//   fmt.Printf("}\n")
//
type Form1234 struct {
    FirstName string
    Style     StyleOption
}

// Finally, every Form has a BuildPDF method that takes a Renderer and returns
// an error. The contents of the method are derived using each Node's Render
// method in turn. For instance:
//
//   fmt.Printf("func (f Form%s) BuildPDF(r frl.Renderer) erorr {\n")
//   for _, node := range f.Nodes {
//       fmt.Println(node.Render())
//   }
//   fmt.Println("}")
//
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

To use this package to print form 1234, we can do something like the following:

```go
package main

import (
    // Other imports...

    "github.com/haleyrc/forms/form"
    "github.com/haleyrc/forms/render"
)

func main() {
    // This is a concrete implementation of the frl.Renderer interface. In our
    // fiction, this might output .py files that can be imported into a project
    // using the insert_correct_name PDF library.
    r := render.NewPythonRenderer()

    // We populate the form however is appropriate, but probably as a mix of
    // stored data and an incoming request:
    f := form.Form1234 {
        FirstName: "Joe",
        Style: form.StyleOptionSedan,
    }

    // Then we call BuildPDF, passing it the renderer we've chosen. For testing,
    // this could just be a mock that verifies a sequence of commands was
    // executed. It could also be an externally defined implementation that
    // creates Canvas elements for displaying fields in a browser. It really
    // doesn't matter as long as the implementation correctly translates the
    // inputs into what the client program is expecting.
    if err := f.BuildPDF(r); err != nil {
        panic(err)
    }
}
```
