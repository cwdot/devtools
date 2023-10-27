package datatable

import "github.com/cwdot/go-stdlib/color"

func NewPen(trueC color.Color, falseC color.Color) *Pen {
	return &Pen{
		trueC:  trueC,
		falseC: falseC,
	}
}

type Pen struct {
	trueC  color.Color
	falseC color.Color
}

func (p *Pen) Ternary(value bool, trueT string, falseT string) string {
	if value {
		return it(p.trueC, trueT)
	}
	return it(p.falseC, falseT)
}

func (p *Pen) Mark(value bool, text string) string {
	if value {
		return it(p.trueC, text)
	}
	return it(p.falseC, text)
}

func it(value color.Color, text string) string {
	if value == "" {
		return text
	}
	return value.It(text)
}
