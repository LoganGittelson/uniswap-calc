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

type PoolDayData struct {
	Id      graphql.String
	Date    graphql.Int
	TvlUSD  graphql.String `graphql:"tvlUSD"`
	FeesUSD graphql.String `graphql:"feesUSD"`
}

func main() {

	var pageSize int = 1000
	var dayIncrement int = 86400
	var rangeStart int = 1640995200
	var rangeEnd int = 1646006400

	cummlatives := make(map[string]ValueAndFee)

	// Create the client
	client := graphql.NewClient("https://api.thegraph.com/subgraphs/name/ianlapham/uniswap-v3-subgraph", nil)

	var q struct {
		PoolDayDatas []PoolDayData `graphql:"poolDayDatas(first: $pageSize, where: {date: $day, id_gt: $lastID})"`
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
		// Page through data for a particular date
		for {
			// Run query
			err := client.Query(context.Background(), &q, variables)
			if err != nil {
				log.Fatal(err)
			}

			// Parse query results
			for _, p := range q.PoolDayDatas {
				cID, cValue, cFee := processPoolDay(p)

				// Update cummulative values
				vf, found := cummlatives[cID]
				if !found {
					vf = ValueAndFee{0, 0}
				}
				cummlatives[cID] = ValueAndFee{value: vf.value + cValue, fee: vf.fee + cFee}
			}

			// Print metdata about iteration
			log.Printf("Date: %d", r)
			log.Printf("Records retrieved: %d", len(q.PoolDayDatas))

			// If we have fetched an incomplete page, it must be the last one
			if len(q.PoolDayDatas) < pageSize {
				break
			}

			// Update highest seen ID for paging
			variables["lastID"] = q.PoolDayDatas[len(q.PoolDayDatas)-1].Id

		}
	}

	calcRatios(cummlatives)

}

// Takes one record of PoolDayData and returns the parsed ID, TLV, and Fees
func processPoolDay(p PoolDayData) (cID string, cValue float64, cFee float64) {
	var lastWinner string = "0xe1263d62e3961467ca88c084f0d0d14feb0846d6"
	var err error

	// Print select data
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

	// Parse the current values
	cValue, err = strconv.ParseFloat(string(p.TvlUSD), 64)
	if err != nil {
		log.Fatal(err)
	}

	cFee, err = strconv.ParseFloat(string(p.FeesUSD), 64)
	if err != nil {
		log.Fatal(err)
	}

	cID = strings.Split(string(p.Id), "-")[0]

	return
}

// Takes the cummlative values and fees map and prints the best ratio found
func calcRatios(cummlatives map[string]ValueAndFee) {
	// Setup variables
	bestRatio := float64(0)
	bestPool := ""

	// Iterate through cummlative values dictionary
	for i, t := range cummlatives {
		cRatio := t.fee / t.value
		if t.value == 0 {
			cRatio = 0
		}
		log.Printf("%v: %v / %v = %v", i, t.fee, t.value, cRatio)
		// Constraining best ratio to only pools that have a cummulative TLV greater than what we're considering investing
		if cRatio > bestRatio && t.value > 1 {
			bestRatio = cRatio
			bestPool = i
		}
	}

	// Print best pool found
	log.Printf("Address of pool: %v", bestPool)
	log.Printf("Earnings: $%f", bestRatio)
}

// TODO: Cleanup print statements
// TODO: Figure out what requirements to put on valid returns
// TODO: Integrate days active scaling factor
