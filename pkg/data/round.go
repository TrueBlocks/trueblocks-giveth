package data

import (
	"encoding/json"
	"log"

	"github.com/TrueBlocks/trueblocks-giveth/pkg/types"
	"github.com/bykof/gostradamus"
)

type Round struct {
	Id        int
	StartDate gostradamus.DateTime
	EndDate   gostradamus.DateTime
	Available int
	Price     string
}

func (r Round) String() string {
	rr := roundInternal{
		Id:        r.Id,
		Start:     r.StartDate.Format("YYYY-MM-DDTHH:mm:ss"),
		End:       r.EndDate.Format("YYYY-MM-DDTHH:mm:ss"),
		Available: r.Available,
		Price:     r.Price,
	}
	b, _ := json.MarshalIndent(rr, "", "  ")
	return string(b)
}

type roundInternal struct {
	Id        int    `json:"id"`
	Start     string `json:"start"`
	End       string `json:"end"`
	Available int    `json:"available"`
	Price     string `json:"price"` // need a string to preserve decimals
}

func GetRounds() (rounds []Round) {
	for i := 1; i <= 25; i++ {
		rounds = append(rounds, Round{
			Id:        i,
			StartDate: types.NewDateTime(2021, 12, 10+(14*i), 16, 0, 0),
			EndDate:   types.NewDateTime(2021, 12, 10+(14*(i+1)), 16, 0, -1),
			Available: getParams(i).Available,
			Price:     getParams(i).Price,
		})
	}
	return rounds
}

type Params struct {
	Id        int
	Available int
	Price     string
}

var params []Params

func getParams(i int) Params {
	if len(params) == 0 {
		params = append(params, Params{Id: 0, Available: 1000000, Price: "0.04"})
		params = append(params, Params{Id: 1, Available: 1000000, Price: "0.04"})
		params = append(params, Params{Id: 2, Available: 1000000, Price: "0.04"})
		params = append(params, Params{Id: 3, Available: 1000000, Price: "0.04"})
		params = append(params, Params{Id: 4, Available: 1000000, Price: "0.04"})
		params = append(params, Params{Id: 5, Available: 1000000, Price: "0.04"})
		params = append(params, Params{Id: 6, Available: 1000000, Price: "0.39212738416853016"})
		params = append(params, Params{Id: 7, Available: 1000000, Price: "0.365426065855812"})
		params = append(params, Params{Id: 8, Available: 1000000, Price: "0.2919862076284853"})
		params = append(params, Params{Id: 9, Available: 1000000, Price: "0.23137274115847087"})
		params = append(params, Params{Id: 10, Available: 1000000, Price: "0.11634931201207897"})
		params = append(params, Params{Id: 11, Available: 1000000, Price: "0.09879737529332695"})
		params = append(params, Params{Id: 12, Available: 1000000, Price: "0.058188181877412634"})
		params = append(params, Params{Id: 13, Available: 1000000, Price: "0.04916228394836191"})
		params = append(params, Params{Id: 14, Available: 1000000, Price: "0.04"})
		params = append(params, Params{Id: 15, Available: 1000000, Price: "0.04"})
		params = append(params, Params{Id: 16, Available: 1000000, Price: "0.04"})
		params = append(params, Params{Id: 17, Available: 1000000, Price: "0.04"})
		params = append(params, Params{Id: 18, Available: 1000000, Price: "0.04"})
		params = append(params, Params{Id: 19, Available: 1000000, Price: "0.04"})
		params = append(params, Params{Id: 20, Available: 1000000, Price: "0.04"})
		params = append(params, Params{Id: 21, Available: 1000000, Price: "0.04"})
		params = append(params, Params{Id: 22, Available: 1000000, Price: "0.04"})
		params = append(params, Params{Id: 23, Available: 1000000, Price: "0.04"})
		params = append(params, Params{Id: 24, Available: 1000000, Price: "0.04"})
		params = append(params, Params{Id: 25, Available: 1000000, Price: "0.04"})
		params = append(params, Params{Id: 26, Available: 1000000, Price: "0.04"})
		params = append(params, Params{Id: 27, Available: 1000000, Price: "0.04"})
		params = append(params, Params{Id: 28, Available: 1000000, Price: "0.04"})
		params = append(params, Params{Id: 29, Available: 1000000, Price: "0.04"})
		params = append(params, Params{Id: 30, Available: 1000000, Price: "0.04"})
	}

	for _, p := range params {
		if p.Id == i {
			return p
		}

	}

	log.Fatal("no such parameter", i)
	return Params{}
}
