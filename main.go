package main

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/shurcooL/graphql"
)

// Set DEBUG to true for addtional printouts
const DEBUG = false

type PoolVals struct {
	value      float64
	fee        float64
	earned     float64
	daysActive int
}

type PoolDayData struct {
	Id      graphql.String
	Date    graphql.Int
	TvlUSD  graphql.String `graphql:"tvlUSD"`
	FeesUSD graphql.String `graphql:"feesUSD"`
}

func main() {
	// Page size to request from GraphQL API, higher means fewer queries but could timeout
	const pageSize int = 1000
	// Difference between epoch days
	const dayIncrement int = 86400
	// Start and end of range to query (in epoch format)
	const rangeStart int = 1640995200
	const rangeEnd int = 1646006400

	cummlatives := make(map[string]PoolVals)

	// Create the client
	client := graphql.NewClient("https://api.thegraph.com/subgraphs/name/ianlapham/uniswap-v3-subgraph", nil)

	// Create the query
	var q struct {
		PoolDayDatas []PoolDayData `graphql:"poolDayDatas(first: $pageSize, where: {date: $day, id_gt: $lastID})"`
	}

	// Set the variables
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
				cID, c := processPoolDay(p)

				// Update cummulative values
				pv, found := cummlatives[cID]
				if !found {
					pv = PoolVals{0, 0, 0, 0}
				}
				cummlatives[cID] = PoolVals{value: pv.value + c.value, fee: pv.fee + c.fee, earned: pv.earned + c.earned, daysActive: pv.daysActive + 1}
			}

			// Print metdata about iteration
			if DEBUG {
				log.Printf("Date: %d", r)
				log.Printf("Records retrieved: %d", len(q.PoolDayDatas))
			}

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
func processPoolDay(p PoolDayData) (cID string, c PoolVals) {
	var lastWinner string = "0xd3ca35355106cb8bc5fd7c534275509673319d83"
	var err error

	// Parse the current values
	cID = strings.Split(string(p.Id), "-")[0]

	c.value, err = strconv.ParseFloat(string(p.TvlUSD), 64)
	if err != nil {
		log.Fatal(err)
	}

	c.fee, err = strconv.ParseFloat(string(p.FeesUSD), 64)
	if err != nil {
		log.Fatal(err)
	}

	c.earned = c.fee / c.value

	// Print select data
	if strings.HasPrefix(string(p.Id), lastWinner) && DEBUG {
		log.Printf(
			`
					ID: %v
					Date: %v
					TVL: %v
					Fees: %v
					Earned: %v
				`, cID, p.Date, c.value, c.fee, c.earned,
		)
	}

	if validatePoolVals(c) {
		return
	}

	// If parsed values are not considered valid, return all 0's
	return cID, PoolVals{0, 0, 0, 0}
}

// Given a set of PoolVals, return true if valid
func validatePoolVals(c PoolVals) bool {
	return c.value >= 1
}

// Takes the cummlative values and fees map and prints the best ratio found
func calcRatios(cummlatives map[string]PoolVals) {
	// Setup variables
	bestEarnings := float64(0)
	bestPool := ""

	// Iterate through cummlative values dictionary
	for i, t := range cummlatives {
		cRatio := t.fee / t.value
		if t.value == 0 {
			cRatio = 0
		}
		if DEBUG {
			log.Printf("%v earned a total of %v with an average ratio of %v over %v days", i, t.earned, cRatio, t.daysActive)
		}

		if t.earned > bestEarnings {
			bestEarnings = t.earned
			bestPool = i
		}
	}

	// Print best pool found
	if DEBUG {
		log.Printf("%v had a cummlative tlv of %v and a cummlative fee of %v over %v days", bestPool, cummlatives[bestPool].value, cummlatives[bestPool].fee, cummlatives[bestPool].daysActive)
	}

	log.Printf("Address of pool: %v", bestPool)
	log.Printf("Earnings: $%f", bestEarnings)
}

// TODO: Fill out readme
// TODO: Delete queries.txt

// Future Improvements:
// 	Make start and end date command-line args
//	Cache query results locally for faster subsequent runs
