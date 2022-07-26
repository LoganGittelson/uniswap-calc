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
		} `graphql:"poolDayDatas(first: 100, where: {date: $day, id_gt: $lastID})"`
	}

	variables := map[string]interface{}{
		"day":    graphql.Int(1647388800),
		"lastID": "",
	}

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

	// TODO: Add pagination
	// TODO: Store to cummlative sums

}
