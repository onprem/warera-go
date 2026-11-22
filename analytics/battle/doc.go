// Package battle computes battle loot pools and round timing: how many items
// a pool holds, which ranks get one and what rarity to expect, and how many
// ticks a round takes.
//
// Pool sizes, the ranks that get items, and round ticks are verified
// against recorded battles. The rarity odds of a pool item are the Warera
// Simulator's model, which matches the frequencies of recorded items.
package battle
