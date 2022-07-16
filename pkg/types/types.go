package types

import (
	"github.com/bykof/gostradamus"
)

func NewDateTime(y, m, d, h, mm, s int) gostradamus.DateTime {
	return gostradamus.NewDateTime(y, m, d, h, mm, s, 0, gostradamus.UTC)
}
