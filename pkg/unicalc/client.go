package unicalc

import (
	"github.com/shurcooL/graphql"
)

const API_ENDPOINT = "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v3"

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
		Day:      graphql.Int(0),
		LastID:   graphql.String(""),
	}
}

func (uc *UniswapClient) varsToStringMap() map[string]interface{} {
	return map[string]interface{}{
		"pageSize": uc.PageSize,
		"day":      uc.Day,
		"lastID":   uc.LastID,
	}
}

func (uc *UniswapClient) setDay(day int) {
	uc.Day = graphql.Int(day)
}

func (uc *UniswapClient) setLastID(id graphql.String) {
	uc.LastID = id
}
