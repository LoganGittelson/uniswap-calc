package main

import (
	"context"
	"log"

	"github.com/machinebox/graphql"
)

func main() {

	// create a client (safe to share across requests)
	client := graphql.NewClient("https://api.thegraph.com/subgraphs/name/ianlapham/uniswap-v3-subgraph")

	// make a request
	req := graphql.NewRequest(`
		query ($day: Int!) {
			poolDayDatas (
				where: {
					date: $day
				}
			) {
				id
				date
				tvlUSD
				feesUSD
			}
		}
	`)

	// set any variables
	req.Var("day", 1647388800)

	// set header fields
	req.Header.Set("Cache-Control", "no-cache")

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var respData map[string]interface{}
	err := client.Run(ctx, req, &respData)
	if err != nil {
		log.Fatal(err)
	}

	log.Print(respData)

}
