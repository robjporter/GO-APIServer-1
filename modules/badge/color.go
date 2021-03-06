package badge

// Color represents color of the badge.
type Color string

// ColorScheme contains named colors that could be used to render the badge.
var ColorScheme = map[string]string{
	"brightgreen": "#4c1",
	"green":       "#97ca00",
	"yellow":      "#dfb317",
	"yellowgreen": "#a4a61d",
	"orange":      "#fe7d37",
	"red":         "#e05d44",
	"blue":        "#007ec6",
	"grey":        "#555",
	"gray":        "#555",
	"lightgrey":   "#9f9f9f",
	"lightgray":   "#9f9f9f",
	"brown":       "#804000",
}

// Standard colors.
const (
	ColorBrightgreen = Color("brightgreen")
	ColorGreen       = Color("green")
	ColorYellow      = Color("yellow")
	ColorYellowgreen = Color("yellowgreen")
	ColorOrange      = Color("orange")
	ColorRed         = Color("red")
	ColorBlue        = Color("blue")
	ColorGrey        = Color("grey")
	ColorGray        = Color("gray")
	ColorLightgrey   = Color("lightgrey")
	ColorLightgray   = Color("lightgray")
	ColorBrown       = Color("brown")
)

func ColorString(color string) Color {
	switch color {
	case "brightgreen":
		return Color("brightgreen")
	case "green":
		return Color("green")
	case "yellow":
		return Color("yellow")
	case "yellowgreen":
		return Color("yellowgreen")
	case "orange":
		return Color("orange")
	case "red":
		return Color("red")
	case "blue":
		return Color("blue")
	case "grey":
		return Color("grey")
	case "gray":
		return Color("gray")
	case "lightgrey":
		return Color("lightgrey")
	case "lightgray":
		return Color("lightgray")
	case "brown":
		return Color("brown")
	default:
		return Color("red")
	}
}

func (c Color) String() string {
	color, ok := ColorScheme[string(c)]
	if ok {
		return color
	} else {
		return string(c)
	}
}
