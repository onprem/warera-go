// Package character computes fight stats, prices skill builds, and grades
// equipment rolls.
//
// Skill values and costs read the game config. [StatsOf] is verified
// against the stat totals the API reports, and [SkillPointsSpent] against
// the spent skill points. Equipment grades compare rolls with the ranges in
// the game config.
package character
