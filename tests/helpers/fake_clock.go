package helpers

import "time"

type FakeClock struct {
	Current time.Time
}

func NewFakeClock(t time.Time) *FakeClock {

	return &FakeClock{
		Current: t.UTC(),
	}
}

func (f *FakeClock) Now() time.Time {

	return f.Current.UTC()
}

func (f *FakeClock) Set(t time.Time) {

	f.Current = t.UTC()
}

func (f *FakeClock) Advance(duration time.Duration) {

	f.Current = f.Current.Add(duration)
}
