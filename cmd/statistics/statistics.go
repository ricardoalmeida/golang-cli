package statistics

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Recipe struct {
	Postcode string `json:"postcode"`
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

func Stats(recipes []Recipe, postcode string, from string, to string, words []string) Stat {
	count, _ := countPerPostcodeAndTime(recipes, PostcodePerTime{
		Postcode: postcode,
		From:     from,
		To:       to,
	}) // TODO error

	postcodePerTime := PostcodePerTime{
		Postcode:      postcode,
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
	res := PostcodeDeliveryCount{}
	groupedByPostcode := map[string]int{}
	for _, recipe := range recipes {
		groupedByPostcode[recipe.Postcode]++
		if groupedByPostcode[recipe.Postcode] > res.DeliveryCount {
			res.DeliveryCount = groupedByPostcode[recipe.Postcode]
			res.Postcode = recipe.Postcode
		}
	}
	return res
}

func countPerPostcodeAndTime(recipes []Recipe, postcodePerTime PostcodePerTime) (count int, err error) {
	limitFrom, err := parseTime(postcodePerTime.From)
	if err != nil {
		return
	}
	limitTo, err := parseTime(postcodePerTime.To)
	if err != nil {
		fmt.Println(err)
		return
	}

	recipes = filter(recipes, func(recipe Recipe) bool {
		timeFrom, timeTo, err := parseDelivery(recipe.Delivery)
		if err != nil {
			fmt.Println(err)
			return false
		}

		return recipe.Postcode == postcodePerTime.Postcode &&
			(timeFrom.After(limitFrom) || timeFrom.Equal(limitFrom)) &&
			(timeTo.Before(limitTo) || timeTo.Equal(limitTo))
	})

	return len(recipes), nil
}

func parseDelivery(str string) (from time.Time, to time.Time, err error) {
	re := regexp.MustCompile(`\d+(AM|PM)`)
	result := re.FindAll([]byte(str), -1)
	if len(result) != 2 {
		return time.Time{}, time.Time{}, fmt.Errorf("Invalid timewindow: %v", str)
	}
	timeFrom, _ := parseTime(string(result[0]))
	timeTo, _ := parseTime(string(result[1]))
	return timeFrom, timeTo, nil
}

func parseTime(str string) (result time.Time, err error) {
	form := "3PM"
	result, err = time.Parse(form, str)
	if err != nil {
		return
	}
	return
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
