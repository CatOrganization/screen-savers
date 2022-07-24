package effect

import (
	"math"
	"time"
)

type fadeInOutState string

const (
	fadeIn               fadeInOutState = "fade-in"
	fadeInPendingFadeout fadeInOutState = "fade-in-pending-fade-out"
	fadeOut              fadeInOutState = "fade-out"
	visible              fadeInOutState = "visible"
	invisible            fadeInOutState = "invisible"
)

type FadeInOut struct {
	state fadeInOutState
	value float64

	fadeInDuration  time.Duration
	fadeOutDuration time.Duration
}

func NewFadeInOut(fadeInDuration, fadeOutDuration time.Duration) *FadeInOut {
	return &FadeInOut{
		state:           invisible,
		value:           0,
		fadeInDuration:  fadeInDuration,
		fadeOutDuration: fadeOutDuration,
	}
}

func (f *FadeInOut) Value() float64 {
	return f.value
}

func (f *FadeInOut) OnFadeInEvent() {
	if f.state == invisible || f.state == fadeOut {
		f.state = fadeIn
	}
}

func (f *FadeInOut) OnFadeOutEvent() {
	if f.state == visible {
		f.state = fadeOut
	}

	if f.state == fadeIn {
		f.state = fadeInPendingFadeout
	}
}

func (f *FadeInOut) Step(dt float32) {
	switch f.state {
	case fadeIn, fadeInPendingFadeout:
		f.value += float64(dt) / f.fadeInDuration.Seconds()
		f.value = math.Min(1, f.value)
		if f.value == 1 {
			if f.state == fadeInPendingFadeout {
				f.state = fadeOut
			} else {
				f.state = visible
			}
		}
	case fadeOut:
		f.value -= float64(dt) / f.fadeOutDuration.Seconds()
		f.value = math.Max(0, f.value)
		if f.value == 0 {
			f.state = "invisible"
		}
	}
}
