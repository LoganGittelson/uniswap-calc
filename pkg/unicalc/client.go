package unicalc

import (
	"context"
	"log"

	"github.com/shurcooL/graphql"
)

const API_ENDPOINT = "https://api.thegraph.com/subgraphs/name/ianlapham/uniswap-v3-subgraph"

type UniswapClient struct {
	graph    *graphql.Client
	PageSize graphql.Int
	Day      graphql.Int
	LastID   graphql.String
}

func newUniswapClient(pageSize int) *UniswapClient {
	return &UniswapClient{
		graph:    graphql.NewClient(API_ENDPOINT, nil),
		PageSize: graphql.Int(pageSize),
		Day:      0,
		LastID:   "",
	}
}

func (uc *UniswapClient) varsToStringMap() map[string]interface{} {
	return map[string]interface{}{
		"pageSize": uc.PageSize,
		"day":      uc.Day,
		"lastID":   uc.LastID,
	}
}

func (uc *UniswapClient) Query(q interface{}) error {
	log.Print("Getting vars")
	vars := uc.varsToStringMap()
	log.Print("Running query")
	err := uc.graph.Query(context.Background(), &q, vars)
	log.Printf("Got query: " + err.Error())
	return err
}

// TODO: set functions for day and last id

func (uc *UniswapClient) setDay(day int) {
	uc.Day = graphql.Int(day)
}

func (uc *UniswapClient) setLastID(id graphql.String) {
	uc.LastID = id
}
