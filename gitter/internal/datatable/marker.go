package datatable

import "github.com/cwdot/go-stdlib/color"

func NewMarker() *Marker {
	lut := make(map[string]color.Color)
	return &Marker{
		lut: lut,
		def: color.Normal,
	}
}

type Marker struct {
	lut map[string]color.Color
	def color.Color // default value
}

func (m *Marker) Set(name string, c color.Color) {
	m.lut[name] = c
}

func (m *Marker) Mark(value string) string {
	if c, ok := m.lut[value]; ok {
		return c.It(value)
	}
	return m.def.It(value)
}
