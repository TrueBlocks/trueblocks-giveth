package data

import (
	"io"
	"log"

	"github.com/TrueBlocks/trueblocks-giveth/pkg/output"
)

func GetDataFromCmd(cmd, fn, fmt string) interface{} {
	var iFace interface{}
	switch cmd {
	case "purple-list":
		iFace, _ = NewPurpleList(fn)
	case "eligible":
		fallthrough
	case "not-eligible":
		fallthrough
	case "purple-verified":
		iFace, _ = NewDonations(fn, fmt)
	case "calc-givback":
		iFace, _ = NewGivback(fn, fmt)
	default:
		log.Fatal("Should not happen", cmd)
	}
	return iFace
}

func RenderQueries(fmt string, w io.Writer, queries []Query) {
	for i, q := range queries {
		iFace := GetDataFromCmd(q.Cmd, q.Fn, fmt)
		if iFace != nil {
			switch q.Cmd {
			case "purple-list":
				if i == 0 {
					output.Header(PurpleList{}, w, fmt)
					defer output.Footer(PurpleList{}, w, fmt)
				}
				if fmt == "txt" || fmt == "csv" {
					for j, p := range iFace.(PurpleList).List {
						output.Line(p, w, fmt, i == 0 && j == 0)
					}
				} else {
					output.Line(iFace, w, fmt, true)
				}
			case "eligible":
				fallthrough
			case "not-eligible":
				fallthrough
			case "purple-verified":
				if i == 0 {
					output.Header(Donation{}, w, fmt)
					defer output.Footer(Donation{}, w, fmt)
				}
				for j, d := range iFace.([]Donation) {
					output.Line(d, w, fmt, i == 0 && j == 0)
				}
			case "calc-givback":
				if i == 0 {
					output.Header(Givback{}, w, fmt)
					defer output.Footer(Givback{}, w, fmt)
				}
				for j, d := range iFace.([]Givback) {
					output.Line(d, w, fmt, i == 0 && j == 0)
				}
			}
		}
	}
}
