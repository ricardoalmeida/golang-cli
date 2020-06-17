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
	"sort"

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

type DeliveryCount struct {
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
	UniqueRecipeCount       int             `json:"unique_recipe_count"`
	CountPerRecipe          []RecipeCount   `json:"count_per_recipe"`
	BusiestPostcode         DeliveryCount   `json:"busiest_postcode"`
	CountPerPostcodeAndTime PostcodePerTime `json:"count_per_postcode_and_time"`
}

func stats(recipes []Recipe) Stat {
	return Stat{
		UniqueRecipeCount:       uniqueRecipeCount(recipes),
		CountPerRecipe:          countPerRecipe(recipes),
		BusiestPostcode:         busiestPostcode(recipes),
		CountPerPostcodeAndTime: countPerPostcodeAndTime(recipes),
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

func busiestPostcode(recipes []Recipe) DeliveryCount {
	return DeliveryCount{
		Postcode:      "10120",
		DeliveryCount: 1000,
	}
}

func countPerPostcodeAndTime(recipes []Recipe) PostcodePerTime {
	return PostcodePerTime{
		Postcode:      "10120",
		From:          "11AM",
		To:            "3PM",
		DeliveryCount: 500,
	}
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
		stat := stats(recipes)
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
	// statsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
