package console

import (
	"context"
	"fmt"
	"strings"

	"github.com/onprem/warera-go"
)

// SelectPlayer resolves a handle or name search to one player profile. The API
// search result is lightweight; callers can use its ID for a full profile lookup.
func SelectPlayer(ctx context.Context, client *warera.Client, prompt *Prompt) (*warera.UserLite, error) {
	for {
		query, err := prompt.Ask(ctx, "Enter a player handle or part of a handle:")
		if err != nil {
			return nil, fmt.Errorf("read player handle: %w", err)
		}
		matches, err := client.Search.SearchUsers(ctx, query)
		if err != nil {
			return nil, fmt.Errorf("search players: %w", err)
		}
		if len(matches) == 0 {
			fmt.Fprintf(prompt.out, "No players matched %q. Try a more specific handle.\n", query)
			continue
		}

		// Prefer a case-insensitive exact username match to avoid asking the user
		// to choose when a search also returns similarly named players.
		exact := make([]warera.UserLite, 0)
		for _, user := range matches {
			if strings.EqualFold(user.Username, query) {
				exact = append(exact, user)
			}
		}
		if len(exact) > 0 {
			matches = exact
			if len(exact) == 1 {
				fmt.Fprintf(prompt.out, "Found exact match: %s [%s].\n", exact[0].Username, exact[0].ID)
				return &exact[0], nil
			}
		}
		// Search may return many partial matches. Keep the menu readable and invite
		// a narrower query if the desired player is outside the displayed set.
		if len(matches) > 10 {
			matches = matches[:10]
			fmt.Fprintln(prompt.out, "Showing the first 10 matches. Enter a more specific handle if yours is not listed.")
		}
		options := make([]string, len(matches))
		for i, user := range matches {
			options[i] = fmt.Sprintf("%s (level %d, ID %s)", user.Username, user.Level(), user.ID)
		}
		selected, err := chooseSingleOrPrompt(ctx, prompt, "Select the player:", options)
		if err != nil {
			return nil, err
		}
		return &matches[selected], nil
	}
}

// SelectCountry resolves a country name or code to one country record.
// Countries are loaded once because the list is small and useful for local
// name/code matching without issuing a search request for every keystroke.
func SelectCountry(ctx context.Context, client *warera.Client, prompt *Prompt) (*warera.Country, error) {
	countries, err := client.Countries.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("load countries: %w", err)
	}
	if len(countries) == 0 {
		return nil, fmt.Errorf("the API returned no countries")
	}
	for {
		query, err := prompt.Ask(ctx, "Enter a country name or code:")
		if err != nil {
			return nil, fmt.Errorf("read country query: %w", err)
		}

		exact := make([]warera.Country, 0)
		for _, country := range countries {
			if strings.EqualFold(country.Name, query) || strings.EqualFold(country.Code, query) {
				exact = append(exact, country)
			}
		}
		if len(exact) == 1 {
			fmt.Fprintf(prompt.out, "Found exact match: %s (%s) [%s].\n", exact[0].Name, exact[0].Code, exact[0].ID)
			return &exact[0], nil
		}
		matches := exact
		if len(matches) == 0 {
			// If the input was not exact, offer case-insensitive name/code
			// substring matches so partial names remain convenient to explore.
			for _, country := range countries {
				if strings.Contains(strings.ToLower(country.Name), strings.ToLower(query)) ||
					strings.Contains(strings.ToLower(country.Code), strings.ToLower(query)) {
					matches = append(matches, country)
				}
			}
		}
		if len(matches) == 0 {
			fmt.Fprintf(prompt.out, "No country matched %q. Try a country name or code.\n", query)
			continue
		}
		if len(matches) > 20 {
			matches = matches[:20]
			fmt.Fprintln(prompt.out, "Showing the first 20 matches. Enter a more specific country name or code if yours is not listed.")
		}
		options := make([]string, len(matches))
		for i, country := range matches {
			options[i] = fmt.Sprintf("%s (%s), population %d, ID %s", country.Name, country.Code, country.CurrentPopulation, country.ID)
		}
		selected, err := chooseSingleOrPrompt(ctx, prompt, "Select the country:", options)
		if err != nil {
			return nil, err
		}
		return &matches[selected], nil
	}
}

// SelectMilitaryUnit resolves a unit name or partial name through the global
// search endpoint and hydrates the returned IDs into full military-unit records.
func SelectMilitaryUnit(ctx context.Context, client *warera.Client, prompt *Prompt) (*warera.MilitaryUnit, error) {
	for {
		query, err := prompt.Ask(ctx, "Enter a military-unit name or part of a name:")
		if err != nil {
			return nil, fmt.Errorf("read military-unit query: %w", err)
		}
		result, err := client.Search.SearchAnything(ctx, query)
		if err != nil {
			return nil, fmt.Errorf("search military units: %w", err)
		}
		if len(result.MUIDs) == 0 {
			fmt.Fprintf(prompt.out, "No military units matched %q. Try a more specific name.\n", query)
			continue
		}
		units, err := client.MilitaryUnits.GetMany(ctx, result.MUIDs)
		if err != nil {
			return nil, fmt.Errorf("load matching military units: %w", err)
		}
		if len(units) > 20 {
			units = units[:20]
			fmt.Fprintln(prompt.out, "Showing the first 20 matches. Enter a more specific name if yours is not listed.")
		}
		options := make([]string, len(units))
		for i, unit := range units {
			options[i] = fmt.Sprintf("%s | country %s | region %s | %d listed members | ID %s",
				unit.Name, unit.CountryID, unit.Region, len(unit.Members), unit.ID)
		}
		selected, err := chooseSingleOrPrompt(ctx, prompt, "Select the military unit:", options)
		if err != nil {
			return nil, err
		}
		return &units[selected], nil
	}
}

func chooseSingleOrPrompt(ctx context.Context, prompt *Prompt, title string, options []string) (int, error) {
	return prompt.Choose(ctx, title, options)
}
