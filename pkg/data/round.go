package data

import (
	"encoding/json"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/paths"
	tbUtils "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/TrueBlocks/trueblocks-giveth/pkg/utils"
	"github.com/bykof/gostradamus"
)

type Round struct {
	Id           int
	StartDate    gostradamus.DateTime
	EndDate      gostradamus.DateTime
	GnosisRange  paths.FileRange
	MainnetRange paths.FileRange
	Available    int
	Price        float64
}

func (r Round) String() string {
	rr := roundInternal{
		Id:           r.Id,
		Start:        r.StartDate.Format("YYYY-MM-DDTHH:mm:ss"),
		End:          r.EndDate.Format("YYYY-MM-DDTHH:mm:ss"),
		GnosisRange:  r.GnosisRange,
		MainnetRange: r.MainnetRange,
		Available:    r.Available,
		Price:        r.Price,
	}
	b, _ := json.MarshalIndent(rr, "", "  ")
	return string(b)
}

type roundInternal struct {
	Id           int             `json:"id"`
	Start        string          `json:"start"`
	End          string          `json:"end"`
	GnosisRange  paths.FileRange `json:"gnosisRange"`
	MainnetRange paths.FileRange `json:"mainnetRange"`
	Available    int             `json:"available"`
	Price        float64         `json:"price"`
}

func GetRounds(filter, max int) (rounds []Round, err error) {
	max = tbUtils.Max(max, 30)
	for i := 1; i <= max; i++ {
		if filter < 1 || filter == i {
			round := Round{
				Id:           i,
				StartDate:    utils.NewDateTime(2021, 12, 10+(14*i), 16, 0, 0),
				EndDate:      utils.NewDateTime(2021, 12, 10+(14*(i+1)), 16, 0, -1),
				GnosisRange:  params[i].GnosisRange,
				MainnetRange: params[i].MainnetRange,
				Available:    params[i].Available,
				Price:        params[i].Price,
			}
			if filter == -1 || round.StartDate.Time().Before(gostradamus.Now().Time()) {
				rounds = append(rounds, round)
			}
		}
	}
	return rounds, nil
}

type Params struct {
	Id           int
	Available    int
	Price        float64
	GnosisRange  paths.FileRange
	MainnetRange paths.FileRange
}

var params = []Params{
	{Id: 0, Available: 1000000, Price: 0.0400000000000000000},
	{Id: 1, Available: 1000000, Price: 0.0400000000000000000, GnosisRange: paths.FileRange{First: 19747830, Last: 19983928}, MainnetRange: paths.FileRange{First: 13868853, Last: 13959253}},
	{Id: 2, Available: 1000000, Price: 0.0400000000000000000, GnosisRange: paths.FileRange{First: 19983929, Last: 20225243}, MainnetRange: paths.FileRange{First: 13959253, Last: 14049867}},
	{Id: 3, Available: 1000000, Price: 0.0400000000000000000, GnosisRange: paths.FileRange{First: 20225244, Last: 20466904}, MainnetRange: paths.FileRange{First: 14049867, Last: 14140481}},
	{Id: 4, Available: 1000000, Price: 0.0400000000000000000, GnosisRange: paths.FileRange{First: 20466905, Last: 20705042}, MainnetRange: paths.FileRange{First: 14140481, Last: 14231033}},
	{Id: 5, Available: 1000000, Price: 0.0400000000000000000, GnosisRange: paths.FileRange{First: 20705043, Last: 20936540}, MainnetRange: paths.FileRange{First: 14231033, Last: 14321469}},
	{Id: 6, Available: 1000000, Price: 0.3921273841685301600, GnosisRange: paths.FileRange{First: 20936541, Last: 21176019}, MainnetRange: paths.FileRange{First: 14321469, Last: 14411361}},
	{Id: 7, Available: 1000000, Price: 0.3654260658558120000, GnosisRange: paths.FileRange{First: 21176020, Last: 21413720}, MainnetRange: paths.FileRange{First: 14411361, Last: 14501315}},
	{Id: 8, Available: 1000000, Price: 0.2919862076284853000, GnosisRange: paths.FileRange{First: 21413721, Last: 21654153}, MainnetRange: paths.FileRange{First: 14501315, Last: 14590874}},
	{Id: 9, Available: 1000000, Price: 0.2313727411584708700, GnosisRange: paths.FileRange{First: 21654154, Last: 21887672}, MainnetRange: paths.FileRange{First: 14590874, Last: 14680127}},
	{Id: 10, Available: 1000000, Price: 0.116349312012078970, GnosisRange: paths.FileRange{First: 21887673, Last: 22126195}, MainnetRange: paths.FileRange{First: 14680127, Last: 14768318}},
	{Id: 11, Available: 1000000, Price: 0.098797375293326950, GnosisRange: paths.FileRange{First: 22126196, Last: 22361918}, MainnetRange: paths.FileRange{First: 14768318, Last: 14855026}},
	{Id: 12, Available: 1000000, Price: 0.058188181877412634, GnosisRange: paths.FileRange{First: 22361919, Last: 22589919}, MainnetRange: paths.FileRange{First: 14855026, Last: 14939330}},
	{Id: 13, Available: 1000000, Price: 0.049162283948361910, GnosisRange: paths.FileRange{First: 22589920, Last: 22827238}, MainnetRange: paths.FileRange{First: 14939330, Last: 15019094}},
	{Id: 14, Available: 1000000, Price: 0.058925060751612926, GnosisRange: paths.FileRange{First: 22827239, Last: 23063889}, MainnetRange: paths.FileRange{First: 15019094, Last: 15102864}},
	{Id: 15, Available: 1000000, Price: 0.040000000000000000, GnosisRange: paths.FileRange{First: 23063890, Last: 23302259}, MainnetRange: paths.FileRange{First: 15102864, Last: 15193339}},
	{Id: 16, Available: 1000000, Price: 0.053374358501172090, GnosisRange: paths.FileRange{First: 23302260, Last: 23542115}, MainnetRange: paths.FileRange{First: 15193339, Last: 15283271}},
	{Id: 17, Available: 1000000, Price: 0.046536434913078640, GnosisRange: paths.FileRange{First: 23542116, Last: 23780247}, MainnetRange: paths.FileRange{First: 15283271, Last: 15372231}},
	{Id: 18, Available: 1000000, Price: 0.045112760240151956, GnosisRange: paths.FileRange{First: 23780248, Last: 24006020}, MainnetRange: paths.FileRange{First: 15372231, Last: 15460088}},
	{Id: 19, Available: 1000000, Price: 0.042464801724041410, GnosisRange: paths.FileRange{First: 24006021, Last: 24210992}, MainnetRange: paths.FileRange{First: 15460089, Last: 15547222}},
	{Id: 20, Available: 1000000, Price: 0.040000000000000000, GnosisRange: paths.FileRange{First: 24210993, Last: 24393124}, MainnetRange: paths.FileRange{First: 15547223, Last: 15640048}},
	{Id: 21, Available: 1000000, Price: 0.040000000000000000, GnosisRange: paths.FileRange{First: 24393125, Last: 00000000}, MainnetRange: paths.FileRange{First: 15640049, Last: 00000000}},
	{Id: 22, Available: 1000000, Price: 0.040000000000000000},
	{Id: 23, Available: 1000000, Price: 0.040000000000000000},
	{Id: 24, Available: 1000000, Price: 0.040000000000000000},
	{Id: 25, Available: 1000000, Price: 0.040000000000000000},
	{Id: 26, Available: 1000000, Price: 0.040000000000000000},
	{Id: 27, Available: 1000000, Price: 0.040000000000000000},
	{Id: 28, Available: 1000000, Price: 0.040000000000000000},
	{Id: 29, Available: 1000000, Price: 0.040000000000000000},
	{Id: 30, Available: 1000000, Price: 0.040000000000000000},
}
