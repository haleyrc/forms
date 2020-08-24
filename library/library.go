package library

import "github.com/haleyrc/forms/fdl"

//go:generate go run generate.go

var Forms = fdl.NewRegistry()

func init() {
	panicOnError(Forms.Register(fdl.Form{
		Code: "1234",
		Nodes: []fdl.Node{
			fdl.TextNode{
				Name:     "FirstName",
				Location: fdl.Location{X: 10, Y: 20},
				FontSize: 12,
			},
			fdl.OptionNode{
				Name: "Style",
				Size: 10,
				Options: []fdl.Option{
					fdl.Option{
						Location: fdl.Location{X: 10, Y: 30},
						Name:     "Sedan",
					},
					fdl.Option{
						Location: fdl.Location{X: 30, Y: 30},
						Name:     "Truck",
					},
				},
			},
		},
	}))
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
