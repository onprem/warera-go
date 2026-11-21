// Package production computes what items cost to produce and what a
// production point earns on them, from the recipes in [warera.GameConfig]
// and market prices.
//
// Production costs read the game config. [Value.MoneyPerPointFromRaw]
// is verified: it matches the server's MarketStats.GetPPMoneyValueByItem
// within the price movement between calls.
package production
