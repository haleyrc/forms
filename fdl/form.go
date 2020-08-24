package fdl

type FormCode string

type Form struct {
	Code           FormCode
	BackgroundFile string
	Nodes          []Node
}
