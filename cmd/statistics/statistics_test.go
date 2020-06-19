package statistics

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
		{"10120", "Meatloaf à La Mom", "Saturday 11AM - 3PM"},
	}

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

	got := Stats(recipes, "10120", "11AM", "3PM", []string{"Meatloaf", "Fajitas"})
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

func TestBusiestPostcode(t *testing.T) {
	got := busiestPostcode([]Recipe{
		{"10224", "Creamy Dill Chicken", "Wednesday 1AM - 7PM"},
		{"10208", "Speedy Steak Fajitas", "Thursday 7AM - 5PM"},
		{"10161", "Meatloaf à La Mom", "Saturday 10AM - 6PM"},
		{"10224", "Creamy Dill Chicken", "Wednesday 1AM - 7PM"},
		{"10220", "Spinach Artichoke Pasta Bake", "Monday 5AM - 4PM"},
		{"10161", "Meatloaf à La Mom", "Saturday 10AM - 6PM"},
		{"10161", "Meatloaf à La Mom", "Saturday 10AM - 6PM"},
	})

	want := PostcodeDeliveryCount{Postcode: "10161", DeliveryCount: 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("busiestPostcode() = %v; want %v", got, want)
	}
}

func TestCountPerPostcodeAndTime(t *testing.T) {
	recipes := []Recipe{
		{"10224", "Creamy Dill Chicken", "Wednesday 10AM - 2PM"},
		{"10224", "Creamy Dill Chicken", "Saturday 9AM - 2PM"},
		{"10224", "Spinach Artichoke Pasta Bake", "Wednesday 11AM - 2PM"},
		{"10208", "Speedy Steak Fajitas", "Thursday 10AM - 1PM"},
	}
	got, _ := countPerPostcodeAndTime(recipes, PostcodePerTime{
		Postcode: "10224",
		From:     "10AM",
		To:       "3PM",
	})

	if got != 2 {
		t.Fatalf("countPerPostcodeAndTime() = %v; want %v", got, 2)
	}
}

func TestCountPerPostcodeAndTime_DeliveryFormatError(t *testing.T) {
	recipes := []Recipe{
		{"10224", "Creamy Dill Chicken", "Wednesday 10 AM - 2PM"},
		{"10224", "Creamy Dill Chicken", "Wednesday 10AM - 2PM"},
	}

	want := "Invalid timewindow: Wednesday 10 AM - 2PM"
	got, err := countPerPostcodeAndTime(recipes, PostcodePerTime{
		Postcode: "10224",
		From:     "10AM",
		To:       "3PM",
	})
	if err != nil {
		t.Fatalf("countPerPostcodeAndTime() err %v; want nil", err)
	}
	if got != 1 {
		t.Fatalf("countPerPostcodeAndTime() error = %v; want %v", got, want)
	}
}

func TestMatchByName(t *testing.T) {
	recipes := []Recipe{
		{"10224", "Spinach Artichoke Pasta Bake", "Wednesday 11AM - 2PM"},
		{"10224", "Creamy Dill Chicken", "Wednesday 10AM - 2PM"},
		{"10224", "Creamy Dill Chicken", "Saturday 9AM - 2PM"},
		{"10208", "Speedy Steak Fajitas", "Thursday 10AM - 1PM"},
	}

	want := []string{"Creamy Dill Chicken", "Spinach Artichoke Pasta Bake"}
	got := matchByName(recipes, []string{"chicken", "pasta"})
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("matchByName() = %v; want %v", got, want)
	}
}
