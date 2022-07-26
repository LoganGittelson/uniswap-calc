package main

import (
	"context"
	"log"

	"github.com/shurcooL/graphql"
)

func main() {

	// Create the client
	client := graphql.NewClient("https://api.thegraph.com/subgraphs/name/ianlapham/uniswap-v3-subgraph", nil)

	var q struct {
		PoolDayDatas []struct {
			Id      graphql.String
			Date    graphql.Int
			TvlUSD  graphql.String `graphql:"tvlUSD"`
			FeesUSD graphql.String `graphql:"feesUSD"`
		} `graphql:"poolDayDatas(first: 1, where: {date: $day, id_gt: $lastID})"`
	}

	variables := map[string]interface{}{
		"day":    graphql.Int(1647388800),
		"lastID": "0x8ad599c3a0ff1de082011efddc58f1908eb6e6d7-19067",
	}

	for {

		err := client.Query(context.Background(), &q, variables)
		if err != nil {
			log.Fatal(err)
		}

		for _, p := range q.PoolDayDatas {
			log.Printf(
				`
			ID: %v
			Date: %v
			TVL: %v
			Fees: %v
			`, p.Id, p.Date, p.TvlUSD, p.FeesUSD,
			)
		}

		log.Printf("Records retrieved: %d", len(q.PoolDayDatas))

		if len(q.PoolDayDatas) == 0 {
			break
		}
		break

		// TODO: Add pagination
		// TODO: Iterate through dates
		// TODO: Store to cummlative sums
	}

}
