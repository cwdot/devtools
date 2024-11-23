package color

func NewRgb(r int, g int, b int) *Color {
	return &Color{
		rgb: []int{r, g, b},
	}
}
func NewRgbw(r int, g int, b int, white int) *Color {
	return &Color{
		rgbw: []int{r, g, b, white},
	}
}
func NewRgbww(r int, g int, b int, white int, warm int) *Color {
	return &Color{
		rgbww: []int{r, g, b, white, warm},
	}
}

type Color struct {
	rgb   []int
	rgbw  []int
	rgbww []int
}

func (c *Color) Values() (string, any) {
	switch {
	case c.rgb != nil:
		return "rgb_color", c.rgb
	case c.rgbw != nil:
		return "rgbw_color", c.rgbw
	case c.rgbww != nil:
		return "rgbww_color", c.rgbww
	}
	return "", ""
}
