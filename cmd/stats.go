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
	"os"
	"strconv"
	"strings"

	stat "github.com/ricardoalmeida/golang-cli/cmd/statistics"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Process input (json file) and returns Recipe Stats json (stdout)",
	Long: `Process input (json file) and returns Recipe Stats json (stdout).
If no flags are provides we are assuming the defaults described in the test requirement.`,
	Args: func(cmd *cobra.Command, args []string) error {
		from := cmd.Flag("from").Value.String()
		to := cmd.Flag("to").Value.String()

		err := stat.ValidateTimeWindow(from, to)
		if err != nil {
			return fmt.Errorf("%s", err)
		}
		_, err = strconv.Atoi(cmd.Flag("postcode").Value.String())
		if err != nil {
			return fmt.Errorf("Invalid Postcode: %s", err)
		}

		if _, err := os.Stat(cmd.Flag("file").Value.String()); err != nil {
			return fmt.Errorf("File does not exist")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		from := cmd.Flag("from").Value.String()
		to := cmd.Flag("to").Value.String()
		filePath := cmd.Flag("file").Value.String()
		words := strings.Split(cmd.Flag("word").Value.String(), ",")
		postcode, err := strconv.Atoi(cmd.Flag("postcode").Value.String())

		input, _ := os.Open(filePath)
		dec := json.NewDecoder(input)
		stat := stat.Stats(dec, postcode, from, to, words)

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

	statsCmd.PersistentFlags().String("postcode", "10120", "Postcode number used for counting per postcode and time")
	statsCmd.PersistentFlags().String("from", "10AM", "Format {h}AM or {h}PM used for counting per postcode and time")
	statsCmd.PersistentFlags().String("to", "3PM", "Format {h}AM or {h}PM used for counting per postcode and time")
	statsCmd.PersistentFlags().StringP("file", "f", "./test_calculation_small.json", "Default file for processing")
	statsCmd.PersistentFlags().StringP("word", "w", "Potato, Veggie, Mushroom", "List of words to match recipe names (alphabetically ordered)")
}
