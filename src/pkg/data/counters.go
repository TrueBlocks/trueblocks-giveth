package data

import (
	"bufio"
	"fmt"
	"os"
	"sort"

	"github.com/TrueBlocks/trueblocks-giveth/pkg/output"
)

type SortDirection int

const (
	Forward SortDirection = iota
	Reverse
)

type StringCounter struct {
	Key   string
	Count int
}

type TwoStringCounter struct {
	Key1  string
	Key2  string
	Count int
}

func Sort(iFace interface{}, dir SortDirection) error {
	twoKey, ok := iFace.([]TwoStringCounter)
	if ok {
		sort.Slice(twoKey, func(i, j int) bool {
			if twoKey[i].Count == twoKey[j].Count {
				if twoKey[i].Key1 == twoKey[j].Key1 {
					return twoKey[i].Key2 < twoKey[j].Key2
				}
				return twoKey[i].Key1 < twoKey[j].Key1
			}
			var test bool
			if dir == Forward {
				test = twoKey[i].Count < twoKey[j].Count
			} else {
				test = twoKey[i].Count > twoKey[j].Count
			}
			return test
		})
	} else {
		oneKey, ok := iFace.([]StringCounter)
		if !ok {
			return fmt.Errorf("unknown type")
		}
		sort.Slice(oneKey, func(i, j int) bool {
			if oneKey[i].Count == oneKey[j].Count {
				return oneKey[i].Key < oneKey[j].Key
			}
			var test bool
			if dir == Forward {
				test = oneKey[i].Count < oneKey[j].Count
			} else {
				test = oneKey[i].Count > oneKey[j].Count
			}
			return test
		})
	}

	return nil
}

func WriteSummary(fn string, results interface{}, dir SortDirection, format string) error {

	Sort(results, dir)

	out, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer out.Close()

	w := bufio.NewWriter(out)
	defer w.Flush() // order matters

	array, ok := results.([]TwoStringCounter)
	if ok {
		output.Header(TwoStringCounter{}, w, format)
		defer output.Footer(TwoStringCounter{}, w, format)
		for i, item := range array {
			output.Line(item, w, format, i == 0)
		}
	} else {
		array1, ok := results.([]StringCounter)
		if ok {
			output.Header(StringCounter{}, w, format)
			defer output.Footer(StringCounter{}, w, format)
			for i, item := range array1 {
				output.Line(item, w, format, i == 0)
			}
		} else {
			return fmt.Errorf("unknown type")
		}
	}

	return nil
}
