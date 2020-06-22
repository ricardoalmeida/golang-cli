package statistics

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var hours = [...]string{
	"12AM", "1AM", "2AM", "3AM", "4AM", "5AM", "6AM", "7AM", "8AM", "9AM", "10AM", "11AM",
	"12PM", "1PM", "2PM", "3PM", "4PM", "5PM", "6PM", "7PM", "8PM", "9PM", "10PM", "11PM",
}

type Recipe struct {
	Postcode int    `json:"postcode,string"`
	Recipe   string `json:"recipe"`
	Delivery string `json:"delivery"`
}

type RecipeCount struct {
	Recipe string `json:"recipe"`
	Count  int    `json:"count"`
}

type PostcodeDeliveryCount struct {
	Postcode      string `json:"postcode"`
	DeliveryCount int    `json:"delivery_count"`
}

type PostcodePerTime struct {
	Postcode      string `json:"postcode"`
	From          string `json:"from"`
	To            string `json:"to"`
	DeliveryCount int    `json:"delivery_count"`
}

type Stat struct {
	UniqueRecipeCount       int                   `json:"unique_recipe_count"`
	CountPerRecipe          []RecipeCount         `json:"count_per_recipe"`
	BusiestPostcode         PostcodeDeliveryCount `json:"busiest_postcode"`
	CountPerPostcodeAndTime PostcodePerTime       `json:"count_per_postcode_and_time"`
	MatchByName             []string              `json:"match_by_name"`
}

func Stats(recipes []Recipe, postcode int, from string, to string, words []string) Stat {
	count, _ := countPerPostcodeAndTime(recipes, postcode, from, to) // TODO error

	postcodePerTime := PostcodePerTime{
		Postcode:      formatPostcode(postcode),
		From:          from,
		To:            to,
		DeliveryCount: count,
	}

	countPerRecipe := countPerRecipe(recipes)

	return Stat{
		UniqueRecipeCount:       len(countPerRecipe),
		CountPerRecipe:          countPerRecipe,
		BusiestPostcode:         busiestPostcode(recipes),
		CountPerPostcodeAndTime: postcodePerTime,
		MatchByName:             matchByName(recipes, words),
	}
}

func countPerRecipe(recipes []Recipe) []RecipeCount {
	groupedByName := map[string]int{}
	for _, recipe := range recipes {
		groupedByName[recipe.Recipe]++
	}
	names := make([]string, 0, len(groupedByName))
	for k := range groupedByName {
		names = append(names, k)
	}

	sort.Strings(names)
	res := []RecipeCount{}
	for _, k := range names {
		res = append(res, RecipeCount{Recipe: k, Count: groupedByName[k]})
	}
	return res
}

func busiestPostcode(recipes []Recipe) PostcodeDeliveryCount {
	var count int
	var postcode int
	groupedByPostcode := map[int]int{}
	for _, recipe := range recipes {
		groupedByPostcode[recipe.Postcode]++
		if groupedByPostcode[recipe.Postcode] > count {
			count = groupedByPostcode[recipe.Postcode]
			postcode = recipe.Postcode
		}
	}
	return PostcodeDeliveryCount{DeliveryCount: count, Postcode: formatPostcode(postcode)}
}

func countPerPostcodeAndTime(recipes []Recipe, postcode int, limitFrom string, limitTo string) (count int, err error) {
	recipes = filter(recipes, func(recipe Recipe) bool {
		from, to, err := parseDelivery(recipe.Delivery)
		if err != nil {
			fmt.Println(err)
			return false
		}

		i, _ := indexHour(limitFrom)
		j, _ := indexHour(from)
		k, _ := indexHour(to)
		l, _ := indexHour(limitTo)

		return recipe.Postcode == postcode && i <= j && k <= l
	})

	return len(recipes), nil
}

func parseDelivery(str string) (from string, to string, err error) {
	result := matchHoursInDelivery(str)
	if len(result) != 2 {
		return "", "", fmt.Errorf("Error getting time window from delivery: %v", str)
	}
	return string(result[0]), string(result[1]), nil
}

func filter(recipes []Recipe, f func(Recipe) bool) []Recipe {
	result := make([]Recipe, 0)
	for _, recipe := range recipes {
		if f(recipe) {
			result = append(result, recipe)
		}
	}
	return result
}

func matchByName(recipes []Recipe, words []string) []string {
	namesMap := map[string]bool{}
	names := []string{}
	for _, recipe := range recipes {
		for _, w := range words {
			if namesMap[recipe.Recipe] == false {
				if strings.Contains(strings.ToLower(recipe.Recipe), strings.ToLower(w)) {
					names = append(names, recipe.Recipe)
					namesMap[recipe.Recipe] = true
				}
			}
		}
	}
	sort.Strings(names)
	return names
}

func formatPostcode(postcode int) string {
	return fmt.Sprintf("%05d", postcode)
}

func indexHour(hour string) (int, error) {
	for i, s := range hours {
		if s == hour {
			return i, nil
		}
	}
	return 0, fmt.Errorf("Invalid hour value: %v", hour)
}

func matchHoursInDelivery(hour string) [][]byte {
	re := regexp.MustCompile(`\d+(AM|PM)`)
	return re.FindAll([]byte(hour), -1)
}

func ValidateTimeWindow(from string, to string) (err error) {
	var i int
	var j int
	if i, err = indexHour(from); err != nil {
		return fmt.Errorf("Invalid 'from' value: %v", from)
	}
	if j, err = indexHour(to); err != nil {
		return fmt.Errorf("Invalid 'to' value: %v", to)
	}
	if i >= j {
		return fmt.Errorf("Invalid time window: 'from' %v 'to' %v. A 'from' should be before 'to'", from, to)
	}
	return nil
}
