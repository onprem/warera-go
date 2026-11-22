// advisor prints a war report for one player, built on the analytics
// packages: their fight stats next to the best war build for their level,
// expected and simulated damage, what an 8-hour fight costs and earns at
// market prices, and the gear swaps and skill points that gain the most.
// It is read-only. Run:
//
//	go run ./examples/advisor <username>
//
// Set NO_COLOR to print without colors.
package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/onprem/warera-go"
	"github.com/onprem/warera-go/analytics/character"
	"github.com/onprem/warera-go/analytics/combat"
	"github.com/onprem/warera-go/analytics/optimize"
)

const (
	// food is what the fights below eat; cooked fish heals the most.
	food        = "cookedFish"
	fightHours  = 8
	simulations = 500
	width       = 74
)

var fightSkills = []string{"attack", "precision", "criticalChance", "criticalDamages", "armor", "dodge", "health", "hunger", "lootChance"}

var skillNames = map[string]string{
	"attack": "Attack", "precision": "Precision", "criticalChance": "Crit chance", "criticalDamages": "Crit damage",
	"armor": "Armor", "dodge": "Dodge", "health": "Health", "hunger": "Hunger", "lootChance": "Loot chance",
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run ./examples/advisor <username>")
		os.Exit(2)
	}
	if err := run(os.Args[1]); err != nil {
		log.Fatal(err)
	}
}

func run(username string) error {
	client := warera.NewClient(warera.WithAPIKey(os.Getenv("WARERA_API_KEY")))
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	user, err := findUser(ctx, client, username)
	if err != nil {
		return err
	}
	equipment, err := client.Inventory.GetCurrentEquipment(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("load equipment: %w", err)
	}
	cfg, err := client.GameConfig.Get(ctx)
	if err != nil {
		return fmt.Errorf("load game config: %w", err)
	}
	prices, err := client.Trading.GetPrices(ctx)
	if err != nil {
		return fmt.Errorf("load prices: %w", err)
	}
	gearPrices, err := equipmentPrices(ctx, client, cfg, equipment)
	if err != nil {
		return err
	}

	// Compare builds without boosts, so a pill or its after-effect does not
	// skew the numbers; the header says when one is active.
	a := combat.DefaultAssumptions()
	metric := optimize.Damage(cfg, food, fightHours, a)
	yours := character.LoadoutOf(user.Skills, equipment)
	yours.BuffPercent, yours.DebuffPercent = 0, 0
	level := user.Leveling.Level

	// The best war build: every point of the player's level spent on fight
	// skills from scratch, keeping their prestige levels and gear.
	fresh := yours
	fresh.Levels = map[string]int{}
	_, best, err := optimize.Roadmap(cfg, fresh, 0, level, metric)
	if err != nil {
		return fmt.Errorf("plan war build: %w", err)
	}
	myStats, err := character.StatsOf(cfg, yours)
	if err != nil {
		return err
	}
	bestStats, err := character.StatsOf(cfg, best)
	if err != nil {
		return err
	}

	header(ctx, client, user)
	statsSection(myStats, bestStats)
	if err := skillsSection(cfg, user, yours, best); err != nil {
		return err
	}
	if err := damageSection(cfg, myStats, bestStats, a); err != nil {
		return err
	}
	if err := fightSection(cfg, myStats, equipment, prices, gearPrices, a); err != nil {
		return err
	}
	if err := gearSection(cfg, yours, equipment, gearPrices, metric); err != nil {
		return err
	}
	if err := nextPointsSection(cfg, user, yours, metric); err != nil {
		return err
	}
	fmt.Println()
	fmt.Println(dim(fmt.Sprintf("Fights eat %s and spend all health and food; damage per swing follows the community Warera Simulator's model.", food)))
	return nil
}

// findUser resolves a username to a full profile, preferring an exact
// match over partial ones.
func findUser(ctx context.Context, client *warera.Client, username string) (*warera.User, error) {
	matches, err := client.Search.SearchUsers(ctx, username)
	if err != nil && len(matches) == 0 {
		return nil, fmt.Errorf("search players: %w", err)
	}
	i := slices.IndexFunc(matches, func(u warera.UserLite) bool { return strings.EqualFold(u.Username, username) })
	if i < 0 {
		if len(matches) == 0 {
			return nil, fmt.Errorf("no player is called %q", username)
		}
		names := make([]string, 0, 5)
		for _, m := range matches[:min(5, len(matches))] {
			names = append(names, m.Username)
		}
		return nil, fmt.Errorf("no player is called %q; did you mean %s?", username, strings.Join(names, ", "))
	}
	return client.Users.GetByID(ctx, matches[i].ID)
}

// equipmentPrices returns the average market price of each equipped item's
// code and the codes one rarity above and below it.
func equipmentPrices(ctx context.Context, client *warera.Client, cfg *warera.GameConfig, equipment *warera.CurrentEquipment) (map[string]float64, error) {
	rarities := []string{"common", "uncommon", "rare", "epic", "legendary", "mythic"}
	grades, err := character.GradeEquipment(cfg, equipment)
	if err != nil {
		return nil, err
	}
	prices := map[string]float64{}
	for _, g := range grades {
		usage := cfg.Items[g.ItemCode].Usage
		r := slices.Index(rarities, g.Rarity)
		for code, item := range cfg.Items {
			if item.Usage != usage || len(item.DynamicStats) == 0 || abs(slices.Index(rarities, item.Rarity)-r) > 1 {
				continue
			}
			p, err := client.MarketStats.GetEquipmentAveragePrice(ctx, code)
			if err != nil {
				return nil, fmt.Errorf("load price of %s: %w", code, err)
			}
			if p > 0 {
				prices[code] = p
			}
		}
	}
	return prices, nil
}

func header(ctx context.Context, client *warera.Client, user *warera.User) {
	fmt.Println()
	title := fmt.Sprintf(" WAR REPORT · %s ", user.Username)
	fmt.Println(bold(cyan("╭" + title + strings.Repeat("─", width-2-len([]rune(title))) + "╮")))
	where := ""
	if country, err := client.Countries.GetByID(ctx, user.CountryID); err == nil {
		where = " · " + country.Name
	}
	fmt.Printf("  Level %d%s · military rank %d (%s attack)\n", user.Leveling.Level, where, user.MilitaryRank,
		green(fmt.Sprintf("+%g%%", user.Skills.Attack.MilitaryRankPercent)))
	if b := user.Buffs; b != nil {
		now := time.Now()
		switch {
		case len(b.BuffCodes) > 0 && b.BuffEndAt.After(now):
			fmt.Printf("  %s pill active, +%g%% attack for %s more\n", yellow("●"), user.Skills.Attack.BuffsPercent, duration(b.BuffEndAt.Sub(now)))
		case len(b.DebuffCodes) > 0 && b.DebuffEndAt.After(now):
			fmt.Printf("  %s pill after-effect, −%g%% attack for %s more\n", red("●"), user.Skills.Attack.DebuffsPercent, duration(b.DebuffEndAt.Sub(now)))
		}
	}
	fmt.Println(dim("  Numbers below leave boosts out, so builds compare fairly."))
	fmt.Println(bold(cyan("╰" + strings.Repeat("─", width-2) + "╯")))
}

func statsSection(mine, best character.Stats) {
	section("FIGHT STATS", "yours vs. the build with the most 8-hour damage at your level")
	fmt.Printf("  %-14s %14s %14s %12s\n", "", "yours", "war build", "")
	rows := []struct {
		name       string
		mine, best float64
		format     func(float64) string
	}{
		{"Attack", mine.Attack, best.Attack, num},
		{"Precision", mine.Precision, best.Precision, pct},
		{"Crit chance", mine.CriticalChance, best.CriticalChance, pct},
		{"Crit damage", mine.CriticalDamages, best.CriticalDamages, pct},
		{"Armor", mine.ArmorReduction * 100, best.ArmorReduction * 100, pct},
		{"Dodge", mine.DodgeChance * 100, best.DodgeChance * 100, pct},
		{"Health", mine.Health, best.Health, num},
		{"Hunger", mine.Hunger, best.Hunger, num},
		{"Loot chance", mine.LootChance, best.LootChance, pct},
	}
	for _, r := range rows {
		fmt.Printf("  %-14s %14s %14s %12s\n", r.name, r.format(r.mine), r.format(r.best), delta(r.best-r.mine, r.format))
	}
}

func skillsSection(cfg *warera.GameConfig, user *warera.User, mine, best character.Loadout) error {
	section("SKILL LEVELS", "the war build spends every point on fight skills")
	fmt.Printf("  %-14s %8s %10s\n", "", "yours", "war build")
	var fightPoints float64
	for _, s := range fightSkills {
		l, err := character.SkillLevel(cfg, s, mine.Levels[s])
		if err != nil {
			return err
		}
		fightPoints += l.TotalCost
		change := ""
		if d := best.Levels[s] - mine.Levels[s]; d > 0 {
			change = green(fmt.Sprintf("▲ %d", d))
		} else if d < 0 {
			change = red(fmt.Sprintf("▼ %d", -d))
		}
		fmt.Printf("  %-14s %8d %10d   %s\n", skillNames[s], mine.Levels[s], best.Levels[s], change)
	}
	spent := float64(user.Leveling.SpentSkillPoints)
	fmt.Printf("\n  You spend %s of %d points on fight skills and %s on economy skills.\n",
		bold(num(fightPoints)), user.Leveling.TotalSkillPoints, bold(num(spent-fightPoints)))
	reset := spent * cfg.User.ResetSkillsCostPerPoint
	if user.Leveling.FreeReset > 0 {
		fmt.Printf("  Switching is free: you have a free skill reset.\n")
	} else {
		fmt.Printf("  Switching costs a skill reset: %s gold for %s points.\n", bold(money(reset)), num(spent))
	}
	return nil
}

func damageSection(cfg *warera.GameConfig, mine, best character.Stats, a combat.Assumptions) error {
	section("DAMAGE", "expected damage from full health and hunger, eating "+food)
	windows := []struct {
		name  string
		hours float64
	}{{"Burst", 0}, {"8 hours", 8}, {"24 hours", 24}}
	var top float64
	type row struct{ mine, best combat.Fight }
	rows := make([]row, len(windows))
	for i, w := range windows {
		m, err := combat.Window(cfg, mine, food, w.hours, a)
		if err != nil {
			return err
		}
		b, err := combat.Window(cfg, best, food, w.hours, a)
		if err != nil {
			return err
		}
		rows[i] = row{m, b}
		top = max(top, m.Damage, b.Damage)
	}
	for i, w := range windows {
		r := rows[i]
		fmt.Printf("  %-9s yours %s %8s  %s\n", w.name, cyan(bar(r.mine.Damage, top, 24)), num(r.mine.Damage), dim(fmt.Sprintf("%.0f swings", r.mine.Swings)))
		fmt.Printf("  %-9s build %s %8s  %s\n", "", magenta(bar(r.best.Damage, top, 24)), num(r.best.Damage), delta(r.best.Damage/r.mine.Damage*100-100, func(v float64) string { return fmt.Sprintf("%.0f%%", v) }))
	}

	// Simulate the 8-hour fight swing by swing for the spread around the
	// expected damage.
	damages := make([]float64, simulations)
	for i := range damages {
		f, err := combat.Simulate(cfg, mine, food, fightHours, a, rand.New(rand.NewPCG(uint64(i), 0)))
		if err != nil {
			return err
		}
		damages[i] = f.Damage
	}
	slices.Sort(damages)
	q := func(p float64) float64 { return damages[int(p*float64(len(damages)-1))] }
	fmt.Printf("\n  %d simulated 8-hour fights: unlucky %s · typical %s · lucky %s\n", simulations,
		red(num(q(0.1))), bold(num(q(0.5))), green(num(q(0.9))))
	fmt.Println(dim(fmt.Sprintf("  (10th, 50th, and 90th percentile; worst %s, best %s)", num(damages[0]), num(damages[len(damages)-1]))))
	return nil
}

func fightSection(cfg *warera.GameConfig, s character.Stats, eq *warera.CurrentEquipment, prices warera.ItemPrices, gear map[string]float64, a combat.Assumptions) error {
	section("ONE 8-HOUR FIGHT", "what your current kit uses up and earns at market prices")
	value := func(code string) float64 {
		if p, ok := gear[code]; ok {
			return p
		}
		v, _ := combat.CraftCost(cfg, cfg.Items[code].Rarity, prices)
		return v
	}
	var p combat.Prices
	if eq.Weapon != nil {
		p.Weapon = value(eq.Weapon.Code)
	}
	for _, piece := range []warera.EquipmentItem{eq.Helmet, eq.Chest, eq.Pants, eq.Boots, eq.Gloves} {
		if piece.Code != "" {
			p.Armor += value(piece.Code)
		}
	}
	if eq.Ammo != "" {
		p.Ammo = prices[eq.Ammo]
	}
	p.Food = prices[food]
	f, err := combat.Window(cfg, s, food, fightHours, a)
	if err != nil {
		return err
	}
	cost := combat.CostOf(f, p, a)
	income := combat.IncomeOf(f, s.LootChance, combat.Rewards{Case1: prices["case1"], Case2: prices["case2"]})
	line := func(name string, v float64, note string) {
		fmt.Printf("  %-22s %12s  %s\n", name, signedMoney(v), dim(note))
	}
	line("Gear wear", -cost.Wear, fmt.Sprintf("%.1f weapons (%s each), %.1f armor sets (%s each)",
		f.Swings/a.ItemUses, money(p.Weapon), f.Undodged/a.ItemUses, money(p.Armor)))
	line("Ammo", -cost.Ammo, fmt.Sprintf("%.0f × %s", f.Swings, orNone(eq.Ammo)))
	line("Food", -cost.Food, fmt.Sprintf("%.0f × %s", f.Food, food))
	line("Cases dropped", income.Cases, fmt.Sprintf("%.1f case1, %.2f case2", income.Case1, income.Case2))
	fmt.Println("  " + strings.Repeat("─", 36))
	net := income.Total - cost.Total
	fmt.Printf("  %-22s %12s  %s\n", bold("Net"), signedMoney(net), dim(fmt.Sprintf("%s gold per 1k damage before bounties", money(f.PerThousand(-net)))))
	return nil
}

func gearSection(cfg *warera.GameConfig, l character.Loadout, eq *warera.CurrentEquipment, prices map[string]float64, metric optimize.Metric) error {
	section("GEAR", "swapping one rarity up or down, at average market prices")
	changes, err := optimize.GearChanges(cfg, l, eq, prices, metric)
	if err != nil {
		return err
	}
	if len(changes) == 0 {
		fmt.Println("  No priced swaps for your current gear.")
		return nil
	}
	var ups, downs []optimize.GearChange
	for _, c := range changes {
		if c.Cost > 0 {
			ups = append(ups, c)
		} else {
			downs = append(downs, c)
		}
	}
	show := func(title string, cs []optimize.GearChange, n int, per func(optimize.GearChange) string) {
		if len(cs) == 0 {
			return
		}
		fmt.Println("  " + bold(title))
		for i, c := range cs[:min(n, len(cs))] {
			mark := "  "
			if i == 0 {
				mark = yellow("★ ")
			}
			fmt.Printf("  %s%-7s %-8s → %-8s %10s %13s  %s\n", mark, c.Slot, c.From, c.To,
				delta(c.Gain, num), signedMoney(-c.Cost), dim(per(c)))
		}
	}
	show("Best upgrades (8-hour damage, gold to buy)", ups, 3, func(c optimize.GearChange) string {
		return fmt.Sprintf("%s damage per gold", num(c.Gain/c.Cost))
	})
	if len(ups) > 0 && len(downs) > 0 {
		fmt.Println()
	}
	show("Cheapest to sell down", downs, 2, func(c optimize.GearChange) string {
		return fmt.Sprintf("%s damage lost per gold", num(c.Gain/c.Cost))
	})
	return nil
}

func nextPointsSection(cfg *warera.GameConfig, user *warera.User, l character.Loadout, metric optimize.Metric) error {
	section("NEXT SKILL POINTS", "for your current build")
	gains, err := optimize.SkillGains(cfg, l, user.Leveling.Level, metric)
	if err != nil {
		return err
	}
	shown := 0
	for _, g := range gains {
		if g.Gain <= 0 || shown == 3 {
			continue
		}
		shown++
		fmt.Printf("  %d. %-12s → level %-3d %4s points  %10s 8-hour damage\n", shown, skillNames[g.Skill], g.Level, num(g.Points), delta(g.Gain, num))
	}
	if shown == 0 {
		fmt.Println("  Every fight skill you can raise is at its highest level.")
	}
	points := "skill points"
	if user.Leveling.AvailableSkillPoints == 1 {
		points = "skill point"
	}
	fmt.Printf("\n  You have %s unspent.\n", bold(fmt.Sprintf("%d %s", user.Leveling.AvailableSkillPoints, points)))
	return nil
}

// --- formatting ---

var color = os.Getenv("NO_COLOR") == "" && isTerminal()

func isTerminal() bool {
	fi, err := os.Stdout.Stat()
	return err == nil && fi.Mode()&os.ModeCharDevice != 0
}

func paint(code, s string) string {
	if !color {
		return s
	}
	return "\x1b[" + code + "m" + s + "\x1b[0m"
}

func bold(s string) string    { return paint("1", s) }
func dim(s string) string     { return paint("2", s) }
func red(s string) string     { return paint("31", s) }
func green(s string) string   { return paint("32", s) }
func yellow(s string) string  { return paint("33", s) }
func magenta(s string) string { return paint("35", s) }
func cyan(s string) string    { return paint("36", s) }

func section(title, note string) {
	fmt.Println()
	fmt.Println(bold(cyan("━━ "+title+" ")) + dim(note))
}

// bar draws value as a share of top, in eighths of a cell.
func bar(value, top float64, cells int) string {
	eighths := 0
	if top > 0 {
		eighths = int(math.Round(value / top * float64(cells*8)))
	}
	partial := []string{"", "▏", "▎", "▍", "▌", "▋", "▊", "▉"}
	s := strings.Repeat("█", eighths/8) + partial[eighths%8]
	return s + strings.Repeat(" ", cells-len([]rune(s)))
}

// num formats a number compactly: 1,234 · 56.7k · 1.23M.
func num(v float64) string {
	a := math.Abs(v)
	switch {
	case a >= 1e6:
		return fmt.Sprintf("%.2fM", v/1e6)
	case a >= 1e4:
		return fmt.Sprintf("%.1fk", v/1e3)
	case a >= 100 || v == math.Trunc(v):
		return fmt.Sprintf("%.0f", v)
	default:
		return fmt.Sprintf("%.2f", v)
	}
}

func pct(v float64) string   { return fmt.Sprintf("%.0f%%", v) }
func money(v float64) string { return fmt.Sprintf("%.2f", v) }

func signedMoney(v float64) string {
	s := fmt.Sprintf("%+.2f", v)
	if v < 0 {
		return red(s)
	}
	return green(s)
}

// delta formats a change with a sign, green when positive and red when
// negative, padded before coloring so columns stay aligned.
func delta(v float64, format func(float64) string) string {
	if math.Abs(v) < 1e-9 {
		return dim(fmt.Sprintf("%10s", "="))
	}
	s := format(v)
	if v > 0 {
		s = "+" + s
	}
	s = fmt.Sprintf("%10s", s)
	if v > 0 {
		return green(s)
	}
	return red(s)
}

// duration formats a duration as hours and minutes: 2h 05m · 58m.
func duration(d time.Duration) string {
	m := int(d.Round(time.Minute).Minutes())
	if m < 60 {
		return fmt.Sprintf("%dm", m)
	}
	return fmt.Sprintf("%dh %02dm", m/60, m%60)
}

func orNone(s string) string {
	if s == "" {
		return "no ammo"
	}
	return s
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
