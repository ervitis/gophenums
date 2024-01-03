package simple

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
