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

	stat "github.com/ricardoalmeida/golang-cli/cmd/statistics"
	"github.com/spf13/cobra"
)

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
		recipes := []stat.Recipe{
			{"10224", "Creamy Dill Chicken", "Wednesday 1AM - 7PM"},
			{"10208", "Speedy Steak Fajitas", "Thursday 7AM - 5PM"},
		}
		postcode := cmd.Flag("postcode").Value.String()
		from := cmd.Flag("from").Value.String()
		to := cmd.Flag("to").Value.String()

		stat := stat.Stats(recipes, postcode, from, to)
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
