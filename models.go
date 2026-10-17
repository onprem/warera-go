//go:generate go run ./cmd/generate/main.go

package warera

import (
	"encoding/json"
	"sort"
	"time"
)

// HasNext reports whether there is a subsequent page.
func (p *CursorPage[T]) HasNext() bool {
	return p != nil && p.NextCursor != "" && p.NextCursor != "undefined"
}

// Level returns the player's level from Leveling.
func (u *User) Level() int {
	return u.Leveling.Level
}

// Level returns the player's level from Leveling.
func (u *UserLite) Level() int {
	return u.Leveling.Level
}

// NextDayTime parses NextDayAt as time.Time.
func (d *GameDates) NextDayTime() (time.Time, error) {
	return time.Parse(time.RFC3339, d.NextDayAt)
}

// NextRegenTime parses NextRegenAt as time.Time.
func (d *GameDates) NextRegenTime() (time.Time, error) {
	return time.Parse(time.RFC3339, d.NextRegenAt)
}

// AttackerCountry returns the attacker's country ID.
func (b *Battle) AttackerCountry() string {
	return b.Attacker.CountryID
}

// DefenderCountry returns the defender's country ID.
func (b *Battle) DefenderCountry() string {
	return b.Defender.CountryID
}

// DefenderRegion returns the defender's region ID.
func (b *Battle) DefenderRegion() string {
	return b.Defender.Region
}

// List returns a sorted slice of ItemPrice from ItemPrices.
func (p ItemPrices) List() []ItemPrice {
	list := make([]ItemPrice, 0, len(p))
	for k, v := range p {
		list = append(list, ItemPrice{ItemCode: k, Price: v})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].ItemCode < list[j].ItemCode })
	return list
}

// BestBuy returns the highest buy order, or nil if no buy orders exist.
func (t *TopOrders) BestBuy() *TopOrder {
	if t == nil || len(t.BuyOrders) == 0 {
		return nil
	}
	return &t.BuyOrders[0]
}

// BestSell returns the lowest sell order, or nil if no sell orders exist.
func (t *TopOrders) BestSell() *TopOrder {
	if t == nil || len(t.SellOrders) == 0 {
		return nil
	}
	return &t.SellOrders[0]
}

// UnmarshalJSON handles both populated BattleRound objects and unpopulated round ID strings.
func (r *BattleRound) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	if data[0] == '"' {
		var id string
		if err := json.Unmarshal(data, &id); err != nil {
			return err
		}
		r.ID = id
		return nil
	}
	type Alias BattleRound
	var a Alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	*r = BattleRound(a)
	return nil
}
