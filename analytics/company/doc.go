// Package company computes what employees, managers, and companies produce
// and cost, and what a company or a workforce makes per hour.
//
// Employee production and wages, manager pay and bonus, and the happiness
// effects are verified against the numbers the API reports for recorded
// rosters, so they can answer what-if questions the API cannot, such as
// what a different manager would add. [HourlyOf] values a company's output
// at a [market.Basis]: reference prices, or what the order books pay now.
package company
