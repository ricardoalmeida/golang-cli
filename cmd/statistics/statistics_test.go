package statistics

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestStats(t *testing.T) {
	const jsonStream = `
	[
		{"postcode": "10208", "recipe": "Speedy Steak Fajitas", "delivery": "Thursday 7AM - 5PM"},
		{"postcode": "10208", "recipe": "Speedy Steak Fajitas", "delivery": "Thursday 7AM - 5PM"},
		{"postcode": "10224", "recipe": "Creamy Dill Chicken", "delivery": "Wednesday 1AM - 7PM"},
		{"postcode": "10220", "recipe": "Spinach Artichoke Pasta Bake", "delivery": "Monday 5AM - 4PM"},
		{"postcode": "10120", "recipe": "Meatloaf à La Mom", "delivery": "Saturday 11AM - 3PM"}
	]
`
	dec := json.NewDecoder(strings.NewReader(jsonStream))

	want := Stat{
		UniqueRecipeCount: 4,
		CountPerRecipe: []RecipeCount{
			{Recipe: "Creamy Dill Chicken", Count: 1},
			{Recipe: "Meatloaf à La Mom", Count: 1},
			{Recipe: "Speedy Steak Fajitas", Count: 2},
			{Recipe: "Spinach Artichoke Pasta Bake", Count: 1},
		},
		BusiestPostcode: PostcodeDeliveryCount{
			Postcode:      "10208",
			DeliveryCount: 2,
		},
		CountPerPostcodeAndTime: PostcodePerTime{
			Postcode:      "10120",
			From:          "11AM",
			To:            "3PM",
			DeliveryCount: 1,
		},
		MatchByName: []string{"Meatloaf à La Mom", "Speedy Steak Fajitas"},
	}

	got := Stats(dec, 10120, "11AM", "3PM", []string{"Meatloaf", "Fajitas"})
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("stats() = %v; want %v", got, want)
	}
}

func TestStats_DeliveryHasFormatError(t *testing.T) {
	const jsonStream = `
	[
		{"postcode": "10208", "recipe": "Speedy Steak Fajitas", "delivery": "Thursday 7 AM - 5PM"},
		{"postcode": "10208", "recipe": "Speedy Steak Fajitas", "delivery": "Thursday 7AM - 5PM"}
	]
`
	dec := json.NewDecoder(strings.NewReader(jsonStream))

	testStdout, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() err = %v; want %v", err, nil)
	}
	osStdout := os.Stdout // keep backup of the real stdout
	os.Stdout = writer
	defer func() {
		// Undo what we changed when this test is done.
		os.Stdout = osStdout
	}()

	want := Stat{
		UniqueRecipeCount: 1,
		CountPerRecipe: []RecipeCount{
			{Recipe: "Speedy Steak Fajitas", Count: 1},
		},
		BusiestPostcode: PostcodeDeliveryCount{
			Postcode:      "10208",
			DeliveryCount: 1,
		},
		CountPerPostcodeAndTime: PostcodePerTime{
			Postcode:      "10120",
			From:          "11AM",
			To:            "3PM",
			DeliveryCount: 0,
		},
		MatchByName: []string{"Speedy Steak Fajitas"},
	}

	got := Stats(dec, 10120, "11AM", "3PM", []string{"Meatloaf", "Fajitas"})
	writer.Close()

	var buf bytes.Buffer
	io.Copy(&buf, testStdout)
	errorMessage := buf.String()
	if errorMessage != "Error getting time window from delivery: Thursday 7 AM - 5PM\n" {
		t.Fatalf("stats() invalid hour (7 AM) has unexpected space and should be ignored")
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("stats() = %v; want %v", got, want)
	}

}

func TestCountPerRecipe(t *testing.T) {
	input := map[string]int{
		"Meatloaf à La Mom":            3,
		"Speedy Steak Fajitas":         1,
		"Creamy Dill Chicken":          2,
		"Spinach Artichoke Pasta Bake": 1,
	}
	got := countPerRecipe(input)

	want := []RecipeCount{
		{Recipe: "Creamy Dill Chicken", Count: 2},
		{Recipe: "Meatloaf à La Mom", Count: 3},
		{Recipe: "Speedy Steak Fajitas", Count: 1},
		{Recipe: "Spinach Artichoke Pasta Bake", Count: 1},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("countPerRecipe() = %v; want %v", got, want)
	}
}
