package utils

import (
	"strings"

	"github.com/bykof/gostradamus"
)

func NewDateTime(y, m, d, h, mm, s int) gostradamus.DateTime {
	return gostradamus.NewDateTime(y, m, d, h, mm, s, 0, gostradamus.UTC)
}

func GetGivethDate(dt gostradamus.DateTime) string {
	ret := dt.Format(gostradamus.Iso8601)
	ret = strings.Replace(ret, "-", "%2F", -1)
	ret = strings.Replace(ret, ":", "%3A", -1)
	ret = strings.Replace(ret, "T", "-", -1)
	return strings.Split(ret, ".")[0]
}
