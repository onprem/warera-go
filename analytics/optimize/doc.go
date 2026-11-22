// Package optimize suggests how to spend skill points and gold: which skill
// level gains the most per point, a skill plan for the levels ahead, and
// which gear swap gains the most per gold.
//
// It scores builds with a [Metric] over the fight stats of
// [github.com/onprem/warera-go/analytics/character], such as [Damage],
// the expected damage of a fight. Its plans are only as good as the metric:
// the fight model behind [Damage] is the community Warera Simulator's.
package optimize
