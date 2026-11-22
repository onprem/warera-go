// Package optimize suggests how to spend skill points: which skill level
// gains the most per point, and a plan for the levels ahead.
//
// It scores builds with a [Metric] over the fight stats of
// [github.com/onprem/warera-go/analytics/character], such as [Damage],
// the expected damage of a fight. Its plans are only as good as the metric:
// the fight model behind [Damage] is the community Warera Simulator's.
package optimize
