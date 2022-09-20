package unicalc

import (
	"log"
	"strconv"
	"strings"
)

// Set DEBUG to true for addtional printouts
const DEBUG = true

func RunCalcs() error { // Page size to request from GraphQL API, higher means fewer queries but could timeout
	const pageSize int = 1000
	// Difference between epoch days
	const dayIncrement int = 86400
	// Start and end of range to query (in epoch format)
	const rangeStart int = 1640995200
	const rangeEnd int = 1646006400

	cummlatives := make(map[string]PoolVals)

	// Create the client
	if DEBUG {
		log.Print("Establishing connection...")
	}
	uclient := newUniswapClient(pageSize)

	// Create the query
	var q struct {
		PoolDayDatas []PoolDayData `graphql:"poolDayDatas(first: $pageSize, where: {date: $day, id_gt: $lastID})"`
	}

	if DEBUG {
		log.Print("Setup complete, starting iteration")
	}

	// Iterate through dates
	for r := rangeStart; r <= rangeEnd; r += dayIncrement {
		uclient.setDay(r)
		uclient.setLastID("")
		// Page through data for a particular date
		for {
			// Run query
			err := uclient.Query(q)
			if err != nil {
				return err
			}

			// Parse query results
			for _, p := range q.PoolDayDatas {
				cID, c, err := processPoolDay(p)
				if err != nil {
					return err
				}

				// Update cummulative values
				pv, found := cummlatives[cID]
				if !found {
					pv = PoolVals{0, 0, 0, 0}
				}
				cummlatives[cID] = PoolVals{
					value:      pv.value + c.value,
					fee:        pv.fee + c.fee,
					earned:     pv.earned + c.earned,
					daysActive: pv.daysActive + 1}
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
			uclient.setLastID(q.PoolDayDatas[len(q.PoolDayDatas)-1].Id)

		}
	}

	if DEBUG {
		log.Print("Iteration complete, calculating ratios")
	}

	var days = (rangeEnd - rangeStart) / dayIncrement
	calcRatios(cummlatives, days)

	return nil
}

// Takes one record of PoolDayData and returns the parsed ID, TLV, and Fees
func processPoolDay(p PoolDayData) (cID string, c PoolVals, err error) {
	var lastWinner string = "0xd3ca35355106cb8bc5fd7c534275509673319d83"

	// Parse the current values
	cID = strings.Split(string(p.Id), "-")[0]

	c.value, err = strconv.ParseFloat(string(p.TvlUSD), 64)
	if err != nil {
		return
	}

	c.fee, err = strconv.ParseFloat(string(p.FeesUSD), 64)
	if err != nil {
		return
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
	return cID, PoolVals{0, 0, 0, 0}, nil
}

// Given a set of PoolVals, return true if valid
func validatePoolVals(c PoolVals) bool {
	return c.value >= 1
}

// Takes the cummlative values and fees map and prints the best ratio found
func calcRatios(cummlatives map[string]PoolVals, days int) {
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

	var interest = bestEarnings
	var principle = 1.0

	var APR = (((interest / principle) / float64(days)) * 365) * 100
	log.Printf("Calulcated APR of: %v%%", APR)

}

// Future Improvements:
// 	Make start and end date command-line args
//	Cache query results locally for faster subsequent runs
