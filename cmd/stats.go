/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"time"

	"github.com/spf13/cobra"
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
}

func stats(recipes []Recipe, postcodePerTime PostcodePerTime) Stat {
	count, _ := countPerPostcodeAndTime(recipes, postcodePerTime) // TODO error
	postcodePerTime.DeliveryCount = count

	return Stat{
		UniqueRecipeCount:       uniqueRecipeCount(recipes),
		CountPerRecipe:          countPerRecipe(recipes),
		BusiestPostcode:         busiestPostcode(recipes),
		CountPerPostcodeAndTime: postcodePerTime,
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

// statsCmd represents the stats command
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		recipes := []Recipe{
			{"10224", "Creamy Dill Chicken", "Wednesday 1AM - 7PM"},
			{"10208", "Speedy Steak Fajitas", "Thursday 7AM - 5PM"},
		}
		postcode := cmd.Flag("postcode").Value.String()
		from := cmd.Flag("from").Value.String()
		to := cmd.Flag("to").Value.String()
		postcodePerTime := PostcodePerTime{Postcode: postcode, From: from, To: to}
		stat := stats(recipes, postcodePerTime)
		b, err := json.Marshal(&stat)
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}
		fmt.Println(string(b))
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	statsCmd.PersistentFlags().String("postcode", "10120", "A help for postcode")
	statsCmd.PersistentFlags().String("from", "10AM", "A help for from")
	statsCmd.PersistentFlags().String("to", "3PM", "A help for to")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
