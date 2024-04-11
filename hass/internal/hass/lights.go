package hass

import (
	"time"

	"hass/internal/hass/color"
)

type LightOnOpts struct {
	Color      *color.Color
	Flash      string
	TurnOff    time.Duration
	Brightness int
}

func LongFlash() func(*LightOnOpts) {
	return func(s *LightOnOpts) {
		s.Flash = "long"
	}
}

func ShortFlash() func(*LightOnOpts) {
	return func(s *LightOnOpts) {
		s.Flash = "short"
	}
}

func TurnOff(secs int) func(*LightOnOpts) {
	return func(s *LightOnOpts) {
		s.TurnOff = time.Second * time.Duration(secs)
	}
}

func Brightness(b int) func(*LightOnOpts) {
	return func(s *LightOnOpts) {
		s.Brightness = b
	}
}

func Red() func(*LightOnOpts) {
	return func(s *LightOnOpts) {
		s.Color = color.NewRgb(255, 0, 0)
	}
}

func Green() func(*LightOnOpts) {
	return func(s *LightOnOpts) {
		s.Color = color.NewRgb(0, 255, 0)
	}
}

func Blue() func(*LightOnOpts) {
	return func(s *LightOnOpts) {
		s.Color = color.NewRgb(0, 0, 255)
	}
}

func Yellow() func(*LightOnOpts) {
	return func(s *LightOnOpts) {
		s.Color = color.NewRgb(255, 255, 0)
	}
}

func White() func(*LightOnOpts) {
	return func(s *LightOnOpts) {
		s.Color = color.NewRgb(255, 255, 255)
	}
}
