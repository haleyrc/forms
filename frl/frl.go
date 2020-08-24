package frl

type Renderer interface {
	PrintText(Location, FontSize, string) error
	PrintCheck(Location, Size) error
}

type Location struct {
	X int
	Y int
}

type FontSize int

type Size int
