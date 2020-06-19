package statistics

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

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
		Postcode:      fmt.Sprintf("%05d", postcode),
		From:          from,
		To:            to,
		DeliveryCount: count,
	}

	return Stat{
		UniqueRecipeCount:       uniqueRecipeCount(recipes),
		CountPerRecipe:          countPerRecipe(recipes),
		BusiestPostcode:         busiestPostcode(recipes),
		CountPerPostcodeAndTime: postcodePerTime,
		MatchByName:             matchByName(recipes, words),
	}
}

func uniqueRecipeCount(recipes []Recipe) int {
	return len(countPerRecipe(recipes))
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
	return PostcodeDeliveryCount{DeliveryCount: count, Postcode: fmt.Sprintf("%05d", postcode)}
}

func countPerPostcodeAndTime(recipes []Recipe, postcode int, limitFrom string, limitTo string) (count int, err error) {
	recipes = filter(recipes, func(recipe Recipe) bool {
		from, to, err := parseDelivery(recipe.Delivery)
		if err != nil {
			fmt.Println(err)
			return false
		}

		return recipe.Postcode == postcode &&
			fmt.Sprintf("%04s", limitFrom) <= fmt.Sprintf("%04s", from) &&
			fmt.Sprintf("%04s", to) <= fmt.Sprintf("%04s", limitTo)
	})

	return len(recipes), nil
}

func parseDelivery(str string) (from string, to string, err error) {
	re := regexp.MustCompile(`\d+(AM|PM)`)
	result := re.FindAll([]byte(str), -1)
	if len(result) != 2 {
		return "", "", fmt.Errorf("Invalid timewindow: %v", str)
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
