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
	"io/ioutil"
	"strconv"
	"strings"

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
		from := cmd.Flag("from").Value.String()
		to := cmd.Flag("to").Value.String()
		filePath := cmd.Flag("file").Value.String()
		words := strings.Split(cmd.Flag("word").Value.String(), ",")
		postcode, err := strconv.Atoi(cmd.Flag("postcode").Value.String())
		if err != nil {
			panic(err)
		}
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			panic(err)
		}

		recipes := []stat.Recipe{}
		err = json.Unmarshal(data, &recipes)
		if err != nil {
			panic(err)
		}

		stat := stat.Stats(recipes, postcode, from, to, words)
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
	statsCmd.PersistentFlags().StringP("file", "f", "./test_calculation_small.json", "A help for to")
	statsCmd.PersistentFlags().StringP("word", "w", "Potato, Veggie, Mushroom", "A help for word")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
