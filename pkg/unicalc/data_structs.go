package unicalc

import "github.com/shurcooL/graphql"

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
