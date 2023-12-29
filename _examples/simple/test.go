package main

// gophenum:generate color
type color string

const (
	Blue   color = "blue"
	Yellow color = "yellow"
	Red    color = "red"
)

// gophenum:generate car
type car int

const (
	Sub car = iota + 1
	Sport
	Family
)

/*

it generates:

type Color interface {
	__private()

	Value() string
}

func (c color) __private() {}

func (c color) Value() string {
	return string(c)
}
*/
