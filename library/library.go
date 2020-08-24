package library

import "github.com/haleyrc/forms/fdl"

//go:generate go run generate.go

var Forms = fdl.NewRegistry()

func init() {
	panicOnError(Forms.Register(fdl.Form{
		Code:        "1234",
		Description: "The RIC for TX SI sales.",
		Nodes: []fdl.Node{
			fdl.TextNode{
				Name:        "FirstName",
				Description: "The customer's first name.",
				Location:    fdl.Location{X: 10, Y: 20},
				FontSize:    12,
			},
			fdl.OptionNode{
				Name:        "Style",
				Description: "The body style of the vehicle.",
				Size:        10,
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
