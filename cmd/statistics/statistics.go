package statistics

import (
	"encoding/json"
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

func Stats(dec *json.Decoder, postcode int, limitFrom string, limitTo string, words []string) Stat {
	countPerPostcodeAndTime := 0

	groupedByName := map[string]int{}

	var busiestCount int
	var busiestPostcode int
	groupedByPostcode := map[int]int{}

	namesMap := map[string]bool{}
	namesMatch := []string{}

	_, err := dec.Token()
	if err != nil {
		fmt.Println(err)
	}
	for dec.More() {
		var recipe Recipe
		err := dec.Decode(&recipe)
		if err != nil {
			fmt.Println(err)
		}

		from, to, err := parseDelivery(recipe.Delivery)
		if err != nil {
			fmt.Println(err)
			continue
		}

		i, _ := indexHour(limitFrom)
		j, _ := indexHour(from)
		k, _ := indexHour(to)
		l, _ := indexHour(limitTo)

		if recipe.Postcode == postcode && i <= j && k <= l {
			countPerPostcodeAndTime++
		}

		groupedByName[recipe.Recipe]++

		groupedByPostcode[recipe.Postcode]++
		if groupedByPostcode[recipe.Postcode] > busiestCount {
			busiestCount = groupedByPostcode[recipe.Postcode]
			busiestPostcode = recipe.Postcode
		}

		for _, w := range words {
			if namesMap[recipe.Recipe] == false {
				if strings.Contains(strings.ToLower(recipe.Recipe), strings.ToLower(w)) {
					namesMatch = append(namesMatch, recipe.Recipe)
					namesMap[recipe.Recipe] = true
				}
			}
		}
	}

	_, err = dec.Token()
	if err != nil {
		fmt.Println(err)
	}

	postcodePerTime := PostcodePerTime{
		Postcode:      formatPostcode(postcode),
		From:          limitFrom,
		To:            limitTo,
		DeliveryCount: countPerPostcodeAndTime,
	}

	countPerRecipe := countPerRecipe(groupedByName)

	postcodeDeliveryCount := PostcodeDeliveryCount{
		DeliveryCount: busiestCount,
		Postcode:      formatPostcode(busiestPostcode),
	}

	sort.Strings(namesMatch)

	return Stat{
		UniqueRecipeCount:       len(countPerRecipe),
		CountPerRecipe:          countPerRecipe,
		BusiestPostcode:         postcodeDeliveryCount,
		CountPerPostcodeAndTime: postcodePerTime,
		MatchByName:             namesMatch,
	}
}

func countPerRecipe(groupedByName map[string]int) []RecipeCount {
	names := make([]string, 0, len(groupedByName))
	for k := range groupedByName {
		names = append(names, k)
	}

	sort.Strings(names)
	countPerRecipe := []RecipeCount{}
	for _, k := range names {
		countPerRecipe = append(countPerRecipe, RecipeCount{Recipe: k, Count: groupedByName[k]})
	}
	return countPerRecipe
}

func parseDelivery(str string) (from string, to string, err error) {
	result := matchHoursInDelivery(str)
	if len(result) != 2 {
		return "", "", fmt.Errorf("Error getting time window from delivery: %v", str)
	}
	return string(result[0]), string(result[1]), nil
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
