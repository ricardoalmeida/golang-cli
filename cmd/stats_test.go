package cmd

import (
	"reflect"
	"testing"
)

func TestStats(t *testing.T) {
	recipes := []Recipe{
		{"10208", "Speedy Steak Fajitas", "Thursday 7AM - 5PM"},
		{"10208", "Speedy Steak Fajitas", "Thursday 7AM - 5PM"},
		{"10224", "Creamy Dill Chicken", "Wednesday 1AM - 7PM"},
		{"10220", "Spinach Artichoke Pasta Bake", "Monday 5AM - 4PM"},
		{"10161", "Meatloaf à La Mom", "Saturday 10AM - 6PM"},
	}

	want := Stat{
		UniqueRecipeCount: 4,
		CountPerRecipe: []RecipeCount{
			{Recipe: "Creamy Dill Chicken", Count: 1},
			{Recipe: "Meatloaf à La Mom", Count: 1},
			{Recipe: "Speedy Steak Fajitas", Count: 2},
			{Recipe: "Spinach Artichoke Pasta Bake", Count: 1},
		},
		BusiestPostcode: DeliveryCount{
			Postcode:      "10120",
			DeliveryCount: 1000,
		},
		CountPerPostcodeAndTime: PostcodePerTime{
			Postcode:      "10120",
			From:          "11AM",
			To:            "3PM",
			DeliveryCount: 500,
		},
	}

	got := stats(recipes)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("stats() = %v; want %v", got, want)
	}
}

func TestUniqueRecipeCount(t *testing.T) {
	got := uniqueRecipeCount([]Recipe{
		{"10224", "Creamy Dill Chicken", "Wednesday 1AM - 7PM"},
		{"10208", "Speedy Steak Fajitas", "Thursday 7AM - 5PM"},
		{"10220", "Spinach Artichoke Pasta Bake", "Monday 5AM - 4PM"},
		{"10161", "Meatloaf à La Mom", "Saturday 10AM - 6PM"},
		{"10224", "Creamy Dill Chicken", "Wednesday 1AM - 7PM"},
		{"10161", "Meatloaf à La Mom", "Saturday 10AM - 6PM"},
	})

	if got != 4 {
		t.Fatalf("uniqueRecipeCount() = %v; want %v", got, 4)
	}
}

func TestCountPerRecipe(t *testing.T) {
	got := countPerRecipe([]Recipe{
		{"10224", "Creamy Dill Chicken", "Wednesday 1AM - 7PM"},
		{"10208", "Speedy Steak Fajitas", "Thursday 7AM - 5PM"},
		{"10161", "Meatloaf à La Mom", "Saturday 10AM - 6PM"},
		{"10224", "Creamy Dill Chicken", "Wednesday 1AM - 7PM"},
		{"10220", "Spinach Artichoke Pasta Bake", "Monday 5AM - 4PM"},
		{"10161", "Meatloaf à La Mom", "Saturday 10AM - 6PM"},
		{"10161", "Meatloaf à La Mom", "Saturday 10AM - 6PM"},
	})

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
