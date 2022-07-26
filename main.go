package main

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/shurcooL/graphql"
)

type ValueAndFee struct {
	value float64
	fee   float64
}

func main() {

	var pageSize int = 1000
	var dayIncrement int = 86400
	var rangeStart int = 1640995200
	var rangeEnd int = 1646006400

	var lastWinner string = "0x7845cfd7acb64e988988f0eeec47ec84c4fb0021"

	cummlatives := make(map[string]ValueAndFee)

	// Create the client
	client := graphql.NewClient("https://api.thegraph.com/subgraphs/name/ianlapham/uniswap-v3-subgraph", nil)

	var q struct {
		PoolDayDatas []struct {
			Id      graphql.String
			Date    graphql.Int
			TvlUSD  graphql.String `graphql:"tvlUSD"`
			FeesUSD graphql.String `graphql:"feesUSD"`
		} `graphql:"poolDayDatas(first: $pageSize, where: {date: $day, id_gt: $lastID})"`
	}

	variables := map[string]interface{}{
		"pageSize": graphql.Int(pageSize),
		"day":      graphql.Int(0),
		"lastID":   "",
	}

	// Iterate through dates
	for r := rangeStart; r <= rangeEnd; r += dayIncrement {
		variables["day"] = graphql.Int(r)
		variables["lastID"] = ""
		for {

			err := client.Query(context.Background(), &q, variables)
			if err != nil {
				log.Fatal(err)
			}

			for _, p := range q.PoolDayDatas {
				if (p.TvlUSD == "0" && p.FeesUSD != "0") || strings.HasPrefix(string(p.Id), lastWinner) {
					log.Printf(
						`
					ID: %v
					Date: %v
					TVL: %v
					Fees: %v
					`, p.Id, p.Date, p.TvlUSD, p.FeesUSD,
					)
				}
				// Store to cummlative sums

				cValue, err := strconv.ParseFloat(string(p.TvlUSD), 64)
				if err != nil {
					log.Fatal(err)
				}

				cFee, err := strconv.ParseFloat(string(p.FeesUSD), 64)
				if err != nil {
					log.Fatal(err)
				}

				cID := strings.Split(string(p.Id), "-")[0]

				vf, found := cummlatives[cID]
				if !found {
					vf = ValueAndFee{0, 0}
				}
				cummlatives[cID] = ValueAndFee{value: vf.value + cValue, fee: vf.fee + cFee}
			}

			log.Printf("Date: %d", r)
			log.Printf("Records retrieved: %d", len(q.PoolDayDatas))

			if len(q.PoolDayDatas) < pageSize {
				break
			}

			variables["lastID"] = q.PoolDayDatas[len(q.PoolDayDatas)-1].Id

		}
	}

	calcRatios(cummlatives)

}

func calcRatios(cummlatives map[string]ValueAndFee) {
	ratios := make(map[string]float64)
	bestRatio := float64(0)
	bestPool := ""
	for i, t := range cummlatives {
		cRatio := t.fee / t.value
		if t.value == 0 {
			cRatio = 0
		}
		log.Printf("%v: %v / %v = %v", i, t.fee, t.value, cRatio)
		if cRatio > bestRatio {
			bestRatio = cRatio
			bestPool = i
		}
		ratios[i] = cRatio
	}

	log.Printf("Address of pool: %v", bestPool)
	log.Printf("Earnings: $%f", bestRatio)
}
