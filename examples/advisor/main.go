// advisor prints a war report for one player, built on the analytics
// packages: their fight stats next to the best war build for their level,
// expected and simulated damage over a day on a pill, what fighting costs
// and earns at market prices, and the gear swaps and skill points that pay
// off. It is read-only. Run:
//
//	go run ./examples/advisor [flags] <username>
//
// By default the player takes a pill before fighting and fights under a
// country order and an MU order, for their own country, with their MU's
// headquarters. Flags change these assumptions; -h lists them. Set NO_COLOR
// to print without colors.
package main

import (
	"cmp"
	"context"
	"flag"
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
	fightHours  = 8
	simulations = 500
	width       = 78
)

// options are the report's assumptions, set by flags.
type options struct {
	pill                  bool
	food                  string
	goal                  string
	bounty                float64
	countryOrder, muOrder bool
	patriotic, enemy      bool
	hq, base, bunker      int
}

var fightSkills = []string{"attack", "precision", "criticalChance", "criticalDamages", "armor", "dodge", "health", "hunger", "lootChance"}

var skillNames = map[string]string{
	"attack": "Attack", "precision": "Precision", "criticalChance": "Crit chance", "criticalDamages": "Crit damage",
	"armor": "Armor", "dodge": "Dodge", "health": "Health", "hunger": "Hunger", "lootChance": "Loot chance",
}

func main() {
	var o options
	flag.BoolVar(&o.pill, "pill", true, "take a pill (cocain) before fighting: more attack for 8 hours, less for 15.5 after")
	flag.StringVar(&o.food, "food", "cookedFish", "food item code to eat while fighting")
	flag.StringVar(&o.goal, "goal", "damage", `what the war build and skill picks aim for: "damage" (most 8-hour damage) or "cost" (least gold per 1k damage)`)
	flag.Float64Var(&o.bounty, "bounty", 0, "gold paid per 1k damage, to judge gear swaps; 0 judges them against your current cost per 1k")
	flag.BoolVar(&o.countryOrder, "country-order", true, "fight under a country order")
	flag.BoolVar(&o.muOrder, "mu-order", true, "fight under an MU order (needs an MU)")
	flag.BoolVar(&o.patriotic, "patriotic", true, "fight for your own country")
	flag.BoolVar(&o.enemy, "enemy", false, "fight your country's sworn enemy")
	flag.IntVar(&o.hq, "hq", -1, "MU headquarters level; -1 reads it from the player's MU")
	flag.IntVar(&o.base, "base", 0, "military base level in the battle's region")
	flag.IntVar(&o.bunker, "bunker", 0, "bunker level in the battle's region")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: go run ./examples/advisor [flags] <username>")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 1 || strings.TrimSpace(flag.Arg(0)) == "" || (o.goal != "damage" && o.goal != "cost") {
		flag.Usage()
		os.Exit(2)
	}
	if err := run(flag.Arg(0), o); err != nil {
		log.Fatal(err)
	}
}

// report holds everything the sections print.
type report struct {
	cfg       *warera.GameConfig
	user      *warera.User
	equipment *warera.CurrentEquipment
	prices    warera.ItemPrices
	gear      map[string]float64 // average market price by equipment code
	o         options
	a         combat.Assumptions

	// yours and best are fight loadouts: battle bonuses applied and, with
	// -pill, the pill active.
	yours, best character.Loadout
	kit         combat.Prices
	metric      optimize.Metric
}

func run(username string, o options) error {
	client := warera.NewClient(warera.WithAPIKey(os.Getenv("WARERA_API_KEY")))
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	r := &report{o: o, a: combat.DefaultAssumptions()}
	var err error
	if r.user, err = findUser(ctx, client, username); err != nil {
		return err
	}
	if r.equipment, err = client.Inventory.GetCurrentEquipment(ctx, r.user.ID); err != nil {
		return fmt.Errorf("load equipment: %w", err)
	}
	if r.cfg, err = client.GameConfig.Get(ctx); err != nil {
		return fmt.Errorf("load game config: %w", err)
	}
	if r.prices, err = client.Trading.GetPrices(ctx); err != nil {
		return fmt.Errorf("load prices: %w", err)
	}
	if _, ok := r.cfg.Items[o.food]; !ok {
		return fmt.Errorf("unknown food %q", o.food)
	}
	if r.gear, err = equipmentPrices(ctx, client, r.cfg, r.equipment); err != nil {
		return err
	}

	// Battle bonuses: the MU's headquarters level unless -hq says otherwise.
	muName := ""
	if r.user.MUID == "" {
		r.o.muOrder = false
		r.o.hq = 0
	} else if mu, err := client.MilitaryUnits.GetByID(ctx, r.user.MUID); err == nil {
		muName = mu.Name
		if r.o.hq < 0 {
			r.o.hq = mu.ActiveUpgradeLevels.Headquarters
		}
	}
	r.o.hq = max(0, r.o.hq)
	bonus, err := combat.BattlePercent(r.cfg, combat.Bonuses{
		CountryOrder: r.o.countryOrder, MUOrder: r.o.muOrder, Patriotic: r.o.patriotic, Enemy: r.o.enemy,
		HeadquartersLevel: r.o.hq, BaseLevel: r.o.base, BunkerLevel: r.o.bunker,
	})
	if err != nil {
		return err
	}

	r.yours = character.LoadoutOf(r.user.Skills, r.equipment)
	r.yours.BattlePercent, r.yours.DebuffPercent, r.yours.BuffPercent = bonus, 0, 0
	if o.pill {
		r.yours.BuffPercent = r.cfg.Items["cocain"].FlatStats["percentAttack"]
	}
	r.kit = r.kitPrices()
	r.metric = r.goalMetric()

	// The best war build: every point of the player's level spent on fight
	// skills from scratch, keeping their prestige levels and gear.
	fresh := r.yours
	fresh.Levels = map[string]int{}
	if _, r.best, err = optimize.Roadmap(r.cfg, fresh, 0, r.user.Leveling.Level, r.metric); err != nil {
		return fmt.Errorf("plan war build: %w", err)
	}

	r.header(ctx, client, muName, bonus)
	steps := []func() error{r.statsSection, r.skillsSection, r.damageSection, r.fightSection, r.gearSection, r.nextPointsSection}
	for _, step := range steps {
		if err := step(); err != nil {
			return err
		}
	}
	fmt.Println()
	fmt.Println(dim("Damage per swing, wear, and case drops follow the community Warera Simulator's model."))
	return nil
}

// kitPrices values what the player's current kit uses up: gear at average
// market prices (or what crafting it takes), ammo, food, and a pill.
func (r *report) kitPrices() combat.Prices {
	var p combat.Prices
	eq := r.equipment
	if eq.Weapon != nil {
		p.Weapon = r.itemValue(eq.Weapon.Code)
	}
	for _, piece := range []warera.EquipmentItem{eq.Helmet, eq.Chest, eq.Pants, eq.Boots, eq.Gloves} {
		if piece.Code != "" {
			p.Armor += r.itemValue(piece.Code)
		}
	}
	if eq.Ammo != "" {
		p.Ammo = r.prices[eq.Ammo]
	}
	p.Food = r.prices[r.o.food]
	if r.o.pill {
		p.Boosts = r.prices["cocain"]
	}
	return p
}

func (r *report) itemValue(code string) float64 {
	if p, ok := r.gear[code]; ok {
		return p
	}
	v, _ := combat.CraftCost(r.cfg, r.cfg.Items[code].Rarity, r.prices)
	return v
}

// economics is what an 8-hour fight with some stats deals and nets.
type economics struct {
	fight  combat.Fight
	cost   combat.Cost
	income combat.Income
	// net is income minus cost, and per1k what 1k damage costs net of
	// cases, before any bounty.
	net, per1k float64
}

func (r *report) economicsOf(s character.Stats, kit combat.Prices) (economics, error) {
	f, err := combat.Window(r.cfg, s, r.o.food, fightHours, r.a)
	if err != nil {
		return economics{}, err
	}
	e := economics{fight: f, cost: combat.CostOf(f, kit, r.a)}
	e.income = combat.IncomeOf(f, s.LootChance, combat.Rewards{Case1: r.prices["case1"], Case2: r.prices["case2"], BountyPer1k: r.o.bounty})
	e.net = e.income.Total - e.cost.Total
	e.per1k = f.PerThousand(e.cost.Total - e.income.Cases)
	return e, nil
}

// goalMetric scores stats by -goal: 8-hour damage, or the negative net
// cost of 1k damage with the current kit.
func (r *report) goalMetric() optimize.Metric {
	if r.o.goal == "damage" {
		return optimize.Damage(r.cfg, r.o.food, fightHours, r.a)
	}
	return func(s character.Stats) (float64, error) {
		e, err := r.economicsOf(s, r.kit)
		return -e.per1k, err
	}
}

func (r *report) stats(l character.Loadout) character.Stats {
	s, err := character.StatsOf(r.cfg, l)
	if err != nil {
		panic(err) // the loadouts come from the game config
	}
	return s
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

func (r *report) header(ctx context.Context, client *warera.Client, muName string, bonus float64) {
	u, o := r.user, r.o
	fmt.Println()
	title := fmt.Sprintf(" WAR REPORT · %s ", u.Username)
	fmt.Println(bold(cyan("╭" + title + strings.Repeat("─", width-2-len([]rune(title))) + "╮")))
	where := ""
	if country, err := client.Countries.GetByID(ctx, u.CountryID); err == nil {
		where = " · " + country.Name
	}
	if muName != "" {
		where += " · MU " + muName
	}
	fmt.Printf("  Level %d%s · military rank %d (%s attack)\n", u.Leveling.Level, where, u.MilitaryRank,
		green(fmt.Sprintf("+%g%%", u.Skills.Attack.MilitaryRankPercent)))
	if b := u.Buffs; b != nil {
		now := time.Now()
		switch {
		case len(b.BuffCodes) > 0 && b.BuffEndAt.After(now):
			fmt.Printf("  %s Pilled up right now, for %s more.\n", yellow("●"), duration(b.BuffEndAt.Sub(now)))
		case len(b.DebuffCodes) > 0 && b.DebuffEndAt.After(now):
			fmt.Printf("  %s In a pill's after-effect right now, for %s more.\n", red("●"), duration(b.DebuffEndAt.Sub(now)))
		}
	}

	fmt.Println()
	var parts []string
	add := func(on bool, name string, v float64) {
		if on {
			parts = append(parts, fmt.Sprintf("%s %+g", name, v))
		}
	}
	c := r.cfg.Battle
	add(o.countryOrder, "country order", c.CountryOrderBonusPercent)
	add(o.muOrder, "MU order", c.MuOrderBonusPercent)
	add(o.patriotic, "patriotic", c.PatrioticBonusPercent)
	add(o.enemy, "sworn enemy", c.EnemyDamagesBonusPercent)
	add(o.hq > 0, fmt.Sprintf("HQ level %d", o.hq), r.upgradeBonus(warera.UpgradeTypeHeadquarters, o.hq))
	add(o.base > 0, fmt.Sprintf("base level %d", o.base), r.upgradeBonus(warera.UpgradeTypeBase, o.base))
	add(o.bunker > 0, fmt.Sprintf("bunker level %d", o.bunker), r.upgradeBonus(warera.UpgradeTypeBunker, o.bunker))
	if len(parts) == 0 {
		parts = []string{"none"}
	}
	fmt.Printf("  Battle bonus  %s  %s\n", bold(green(fmt.Sprintf("%+g%%", bonus))), dim("("+strings.Join(parts, ", ")+")"))
	if o.pill {
		pill := r.cfg.Items["cocain"].FlatStats
		fmt.Printf("  Pill          %s  %s\n", bold(green(fmt.Sprintf("+%g%%", pill["percentAttack"]))),
			dim(fmt.Sprintf("for %gh, then −%d%% for %gh; %s gold each", pill["buffDurationHours"], combat.PillAfterEffectPercent, pill["debuffDurationHours"], money(r.prices["cocain"]))))
	} else {
		fmt.Printf("  Pill          %s\n", dim("none (-pill=false)"))
	}
	goal := "the most 8-hour damage"
	if o.goal == "cost" {
		goal = "the least gold per 1k damage"
	}
	fmt.Printf("  Goal          %s  %s\n", bold(goal), dim("(-goal)"))
	fmt.Println(bold(cyan("╰" + strings.Repeat("─", width-2) + "╯")))
}

// upgradeBonus returns the bonus one building level gives on its own.
func (r *report) upgradeBonus(t warera.UpgradeType, level int) float64 {
	var b combat.Bonuses
	switch t {
	case warera.UpgradeTypeHeadquarters:
		b.HeadquartersLevel = level
	case warera.UpgradeTypeBase:
		b.BaseLevel = level
	case warera.UpgradeTypeBunker:
		b.BunkerLevel = level
	}
	percent, _ := combat.BattlePercent(r.cfg, b)
	return percent
}

func (r *report) statsSection() error {
	mine, best := r.stats(r.yours), r.stats(r.best)
	section("FIGHT STATS", "yours vs. the best war build for your level, while fighting")
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
	for _, row := range rows {
		fmt.Printf("  %-14s %14s %14s %12s\n", row.name, row.format(row.mine), row.format(row.best), delta(row.best-row.mine, row.format))
	}
	return nil
}

func (r *report) skillsSection() error {
	section("SKILL LEVELS", "the war build spends every point on fight skills")
	fmt.Printf("  %-14s %8s %10s\n", "", "yours", "war build")
	var fightPoints float64
	for _, s := range fightSkills {
		l, err := character.SkillLevel(r.cfg, s, r.yours.Levels[s])
		if err != nil {
			return err
		}
		fightPoints += l.TotalCost
		change := ""
		if d := r.best.Levels[s] - r.yours.Levels[s]; d > 0 {
			change = green(fmt.Sprintf("▲ %d", d))
		} else if d < 0 {
			change = red(fmt.Sprintf("▼ %d", -d))
		}
		fmt.Printf("  %-14s %8d %10d   %s\n", skillNames[s], r.yours.Levels[s], r.best.Levels[s], change)
	}
	lv := r.user.Leveling
	spent := float64(lv.SpentSkillPoints)
	fmt.Printf("\n  You spend %s of %d points on fight skills and %s on economy skills.\n",
		bold(num(fightPoints)), lv.TotalSkillPoints, bold(num(spent-fightPoints)))
	if lv.FreeReset > 0 {
		fmt.Println("  Switching is free: you have a free skill reset.")
	} else {
		fmt.Printf("  Switching costs a skill reset: %s gold for %s points.\n", bold(money(spent*r.cfg.User.ResetSkillsCostPerPoint)), num(spent))
	}
	return nil
}

func (r *report) damageSection() error {
	mine, best := r.stats(r.yours), r.stats(r.best)
	section("DAMAGE", "expected damage from full health and hunger, eating "+r.o.food)
	type row struct {
		name       string
		mine, best float64
		note       string
	}
	var rows []row
	for _, w := range []struct {
		name  string
		hours float64
	}{{"Burst", 0}, {"8 hours", fightHours}} {
		m, err := combat.Window(r.cfg, mine, r.o.food, w.hours, r.a)
		if err != nil {
			return err
		}
		b, err := combat.Window(r.cfg, best, r.o.food, w.hours, r.a)
		if err != nil {
			return err
		}
		rows = append(rows, row{w.name, m.Damage, b.Damage, fmt.Sprintf("%.0f swings", m.Swings)})
	}
	var day combat.Day
	if r.o.pill {
		var bestDay combat.Day
		var err error
		if day, err = combat.PillDay(r.cfg, unboosted(r.yours), r.o.food, r.a); err != nil {
			return err
		}
		if bestDay, err = combat.PillDay(r.cfg, unboosted(r.best), r.o.food, r.a); err != nil {
			return err
		}
		rows = append(rows, row{"Pill day", day.Damage, bestDay.Damage, "8h pilled + 16h after"})
	} else {
		m, err := combat.Window(r.cfg, mine, r.o.food, 24, r.a)
		if err != nil {
			return err
		}
		b, err := combat.Window(r.cfg, best, r.o.food, 24, r.a)
		if err != nil {
			return err
		}
		rows = append(rows, row{"24 hours", m.Damage, b.Damage, fmt.Sprintf("%.0f swings", m.Swings)})
	}
	var top float64
	for _, row := range rows {
		top = max(top, row.mine, row.best)
	}
	for _, row := range rows {
		fmt.Printf("  %-9s yours %s %8s  %s\n", row.name, cyan(bar(row.mine, top, 24)), num(row.mine), dim(row.note))
		fmt.Printf("  %-9s build %s %8s  %s\n", "", magenta(bar(row.best, top, 24)), num(row.best),
			delta(row.best/row.mine*100-100, func(v float64) string { return fmt.Sprintf("%.0f%%", v) }))
	}
	if r.o.pill {
		fmt.Printf("\n  Your pill day: %s while pilled, %s in the after-effect, %s after it.\n",
			bold(num(day.Boosted.Damage)), num(day.AfterEffect.Damage), num(day.Plain.Damage))
	}

	// Simulate the 8-hour fight swing by swing for the spread around the
	// expected damage.
	damages := make([]float64, simulations)
	for i := range damages {
		f, err := combat.Simulate(r.cfg, mine, r.o.food, fightHours, r.a, rand.New(rand.NewPCG(uint64(i), 0)))
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

func unboosted(l character.Loadout) character.Loadout {
	l.BuffPercent, l.DebuffPercent = 0, 0
	return l
}

func (r *report) fightSection() error {
	title := "ONE 8-HOUR FIGHT"
	if r.o.pill {
		title = "ONE 8-HOUR PILLED FIGHT"
	}
	section(title, "what your current kit uses up and earns at market prices")
	mine := r.stats(r.yours)
	e, err := r.economicsOf(mine, r.kit)
	if err != nil {
		return err
	}
	f, p := e.fight, r.kit
	line := func(name string, v float64, note string) {
		fmt.Printf("  %-22s %12s  %s\n", name, signedMoney(v), dim(note))
	}
	line("Weapon wear", -f.Swings/r.a.ItemUses*p.Weapon, fmt.Sprintf("%.1f weapons at %s", f.Swings/r.a.ItemUses, money(p.Weapon)))
	line("Armor wear", -f.Undodged/r.a.ItemUses*p.Armor, fmt.Sprintf("%.1f armor sets at %s (%.0f%% of swings dodged)", f.Undodged/r.a.ItemUses, money(p.Armor), mine.DodgeChance*100))
	line("Ammo", -e.cost.Ammo, fmt.Sprintf("%.0f × %s", f.Swings, orNone(r.equipment.Ammo)))
	line("Food", -e.cost.Food, fmt.Sprintf("%.0f × %s", f.Food, r.o.food))
	if r.o.pill {
		line("Pill", -p.Boosts, "1 × cocain")
	}
	line("Cases dropped", e.income.Cases, fmt.Sprintf("%.1f case1, %.2f case2", e.income.Case1, e.income.Case2))
	if r.o.bounty > 0 {
		line("Bounty", e.income.Bounty, fmt.Sprintf("%s per 1k damage", money(r.o.bounty)))
	}
	fmt.Println("  " + strings.Repeat("─", 36))
	fmt.Printf("  %-22s %12s  %s\n", bold("Net"), signedMoney(e.net),
		dim(fmt.Sprintf("for %s damage: %s gold per 1k before bounties", num(f.Damage), money(e.per1k))))

	if r.o.pill {
		// The after-effect phase: the same wear for far less damage.
		day, err := combat.PillDay(r.cfg, unboosted(r.yours), r.o.food, r.a)
		if err != nil {
			return err
		}
		after := day.AfterEffect
		kit := p
		kit.Boosts = 0
		cost := combat.CostOf(after, kit, r.a)
		income := combat.IncomeOf(after, mine.LootChance, combat.Rewards{Case1: r.prices["case1"], Case2: r.prices["case2"], BountyPer1k: r.o.bounty})
		per1k := after.PerThousand(cost.Total - income.Cases)
		verdict := green("worth it")
		if per1k > 1.5*e.per1k {
			verdict = red("usually not worth it")
		}
		fmt.Printf("\n  Fighting in the after-effect too: %s more damage at %s gold per 1k, %s.\n",
			num(after.Damage), money(per1k), verdict)
	}
	return nil
}

func (r *report) gearSection() error {
	section("GEAR", "swapping one rarity up or down, with wear priced at average market prices")
	changes, err := optimize.GearChanges(r.cfg, r.yours, r.equipment, r.gear, optimize.Damage(r.cfg, r.o.food, fightHours, r.a))
	if err != nil {
		return err
	}
	if len(changes) == 0 {
		fmt.Println("  No priced swaps for your current gear.")
		return nil
	}
	now, err := r.economicsOf(r.stats(r.yours), r.kit)
	if err != nil {
		return err
	}
	// The bar a swap must clear: the bounty when given, or what 1k damage
	// costs you now.
	threshold, thresholdName := now.per1k, "your current cost"
	if r.o.bounty > 0 {
		threshold, thresholdName = r.o.bounty, "the bounty"
	}
	type swap struct {
		optimize.GearChange
		damage, net, marginal float64
	}
	var ups, downs []swap
	for _, c := range changes {
		kit := r.kit
		if c.Slot == "weapon" {
			kit.Weapon += c.Cost
		} else {
			kit.Armor += c.Cost
		}
		e, err := r.economicsOf(r.stats(c.Loadout), kit)
		if err != nil {
			return err
		}
		s := swap{GearChange: c, damage: e.fight.Damage - now.fight.Damage, net: e.net - now.net}
		// Gold per 1k damage the swap adds (upgrade) or saves (downgrade),
		// before bounties.
		spent := (e.cost.Total - e.income.Cases) - (now.cost.Total - now.income.Cases)
		s.marginal = spent / s.damage * 1000
		if c.Cost > 0 {
			ups = append(ups, s)
		} else {
			downs = append(downs, s)
		}
	}
	slices.SortFunc(ups, func(a, b swap) int { return cmp.Compare(a.marginal, b.marginal) })
	slices.SortFunc(downs, func(a, b swap) int { return cmp.Compare(b.marginal, a.marginal) })

	fmt.Printf("  1k damage costs you %s gold now; swaps are judged against %s, %s.\n", bold(money(now.per1k)), thresholdName, money(threshold))
	if r.o.bounty == 0 {
		fmt.Println(dim("  Pass -bounty with what you are paid per 1k damage to judge swaps by profit instead."))
	}
	fmt.Println()
	fmt.Printf("  %-26s %10s %12s %16s\n", "", "damage", "net gold", "gold per 1k")
	show := func(s swap, verdict string) {
		name := fmt.Sprintf("%s → %s", s.From, s.To)
		fmt.Printf("  %-26s %s %s %16s  %s\n", name, delta(s.damage, num), padSigned(s.net), money(s.marginal), verdict)
	}
	if len(ups) > 0 {
		fmt.Println("  " + bold("Upgrades") + dim(" (per 8-hour fight; gold per 1k is what the extra damage costs)"))
		for _, s := range ups {
			if s.marginal <= threshold {
				show(s, green("✔ pays off"))
			} else {
				show(s, red("✘ costs more than it brings"))
			}
		}
	}
	if len(downs) > 0 {
		fmt.Println("\n  " + bold("Downgrades") + dim(" (gold per 1k is what each 1k damage given up saves)"))
		for _, s := range downs {
			if s.marginal > threshold {
				show(s, green("✔ saves more than it loses"))
			} else {
				show(s, dim("✘ keep yours"))
			}
		}
	}
	return nil
}

func (r *report) nextPointsSection() error {
	section("NEXT SKILL POINTS", "for your current build, by your goal")
	gains, err := optimize.SkillGains(r.cfg, r.yours, r.user.Leveling.Level, r.metric)
	if err != nil {
		return err
	}
	shown := 0
	for _, g := range gains {
		if g.Gain <= 0 || shown == 3 {
			continue
		}
		shown++
		// With -goal cost the metric is minus gold per 1k, so a gain is that
		// much cheaper damage.
		effect := "+" + num(g.Gain) + " 8-hour damage"
		if r.o.goal == "cost" {
			effect = fmt.Sprintf("−%.3f gold per 1k damage", g.Gain)
		}
		fmt.Printf("  %d. %-12s → level %-3d %4s points   %s\n", shown, skillNames[g.Skill], g.Level, num(g.Points), green(effect))
	}
	if shown == 0 {
		fmt.Println("  No fight skill level you can buy improves your goal.")
	}
	points := "skill points"
	if r.user.Leveling.AvailableSkillPoints == 1 {
		points = "skill point"
	}
	fmt.Printf("\n  You have %s unspent.\n", bold(fmt.Sprintf("%d %s", r.user.Leveling.AvailableSkillPoints, points)))
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

func padSigned(v float64) string {
	s := fmt.Sprintf("%12s", fmt.Sprintf("%+.2f", v))
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
