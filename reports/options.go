package reports

import (
	"time"

	"github.com/johnfercher/maroto/pkg/consts"
)

type SheetOptions struct {
	size         consts.PageSize
	orientation  consts.Orientation
	numLanes     int
	eventOrder   []OrderFunc
	sessionTimes []time.Time
	bySession    bool
}

func (s *SheetOptions) Size() consts.PageSize {
	if s == nil || s.orientation == "" {
		return consts.A4
	}
	return s.size
}

func (s *SheetOptions) Orientation() consts.Orientation {
	if s == nil || s.orientation == "" {
		return consts.Portrait
	}
	return s.orientation
}

const defaultNumLanes = 3

func (s *SheetOptions) Lanes() int {
	if s == nil || s.numLanes == 0 {
		return defaultNumLanes
	}
	return s.numLanes
}

func (s *SheetOptions) EventOrder() *Order {
	if s == nil || len(s.eventOrder) == 0 {
		return NewOrder(DefaultEventOrder...)
	}
	return NewOrder(s.eventOrder...)
}

func (s *SheetOptions) SessionTime(session int) time.Time {
	if s == nil || len(s.sessionTimes) < session {
		return time.Now()
	}
	return s.sessionTimes[session-1]
}

func (s *SheetOptions) SessionTimes() []time.Time {
	if s == nil || len(s.sessionTimes) == 0 {
		return []time.Time{time.Now()}
	}
	return s.sessionTimes
}

func (s *SheetOptions) BySession() bool {
	if s == nil {
		return false
	}
	return s.bySession
}

type SheetOption func(*SheetOptions)

func SizeOption(size consts.PageSize) SheetOption {
	return SheetOption(func(s *SheetOptions) {
		s.size = size
	})
}

func OrientationOption(orientation consts.Orientation) SheetOption {
	return SheetOption(func(s *SheetOptions) {
		s.orientation = orientation
	})
}

func NumLanesOption(lanes int) SheetOption {
	return SheetOption(func(s *SheetOptions) {
		s.numLanes = lanes
	})
}

func EventOrderOption(order ...OrderFunc) SheetOption {
	return SheetOption(func(s *SheetOptions) {
		s.eventOrder = order
	})
}

func SessionTimesOption(sessions []time.Time) SheetOption {
	return SheetOption(func(s *SheetOptions) {
		s.sessionTimes = sessions
	})
}

func BySessionOption(b bool) SheetOption {
	return SheetOption(func(s *SheetOptions) {
		s.bySession = b
	})
}

func applyOptions(opts []SheetOption) *SheetOptions {
	s := &SheetOptions{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
