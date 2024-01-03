package fail

// gophenum:generate colors
type color string

const (
	Blue   color = "blue"
	Yellow color = "yellow"
	Red    color = "red"
)

// gophenum:generate cars
type car int

const (
	Sub car = iota + 1
	Sport
	Family
)
