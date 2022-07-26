# uniswap-calc

This is a script to fetch the per-day pool values from uniswap V3 (via The Graph) and calculate which pool had the best return over a hard-coded time-frame.

# To Run
```bash
go run  main.go
```

# Calculating Returns
Return is calculated for each pool each day by dividing the total fees earned (`feesUSD`) by the total value (`tvlUSD`). Return over time is the sum of these daily returns.

## Example
Liquidity pool with ID `0x8ad599c3a0ff1de082011efddc58f1908eb6e6d8` on 16th March 2022 (`1647388800`) had US$356,129 of fees collected and a total value of US$391,636,206. We would calculate this as `356129 / 391636206 = 0.00090933625` meaning that if you had invested $1 you would have expected US$0.00090933625 earned for that day.

## Assumptions
* **Assumption 1**: We only consider pools which had a tvl of at least $1 for the day being looked at. Otherwise that's calculated as 0 return.
* **Assumption 2**: No additional fees or rewards are being considered (e.g. gas, additional incentives provided by Uniswap, etc)
* **Assumption 3**: Returns are *not* being reinvested into the pool. This would essentially be as if the rewards were immediately cashed out in USD and set aside each day.