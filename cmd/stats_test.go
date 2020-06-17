package cmd

import (
	"reflect"
	"testing"
)

func TestStats(t *testing.T) {
	recipes := []Recipe{
		{"10224", "Creamy Dill Chicken", "Wednesday 1AM - 7PM"},
		{"10208", "Speedy Steak Fajitas", "Thursday 7AM - 5PM"},
		{"10220", "Spinach Artichoke Pasta Bake", "Monday 5AM - 4PM"},
		{"10161", "Meatloaf à La Mom", "Saturday 10AM - 6PM"},
	}

	stat := Stat{
		UniqueRecipeCount: 15,
		CountPerRecipe: []RecipeCount{
			{Recipe: "Speedy Steak Fajitas", Count: 1},
			{Recipe: "Tex-Mex Tilapia", Count: 3},
			{Recipe: "Mediterranean Baked Veggies", Count: 1},
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

	res := stats(recipes)
	if !reflect.DeepEqual(res, stat) {
		t.Fatalf("stats() = %v; want %v", res, stat)
	}
}
