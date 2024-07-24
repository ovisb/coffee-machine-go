package coffeeMachine

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"strings"
	"testing"
)

func TestNewCoffeeMachine(t *testing.T) {
	t.Parallel()
	got := NewCoffeeMachine(coffeeMachineStartingWater, coffeeMachineStartingMilk, coffeeMachineStartingBeans, coffeeMachineStartingCups, coffeeMachineStartingCash)
	want := &CoffeeMachine{coffeeMachineStartingWater, coffeeMachineStartingMilk, coffeeMachineStartingBeans, coffeeMachineStartingCups, coffeeMachineStartingCash}

	if !cmp.Equal(want, got, cmp.AllowUnexported(CoffeeMachine{})) {
		t.Error(cmp.Diff(want, got, cmp.AllowUnexported(CoffeeMachine{})))
	}
}

func TestGetInput(t *testing.T) {
	t.Run("test integer parsed and returned correctly", func(t *testing.T) {
		number := "25\n"
		reader := strings.NewReader(number)
		want := 25

		got, err := GetInt(reader)

		assertError(t, err, nil)

		if got != want {
			t.Errorf("GetInput() returns %v, want %v", got, want)
		}
	})

	t.Run("test input incorrect / cannot be parsed as integer", func(t *testing.T) {
		number := "asd\n"
		reader := strings.NewReader(number)

		_, err := GetInt(reader)

		if err == nil {
			t.Fatalf("expected error, got success")
		}

		assertError(t, err, expectedIntegerError)
	})

	t.Run("test string parsed and returned correctly", func(t *testing.T) {
		number := "25\n"
		reader := strings.NewReader(number)
		want := "25"

		got, err := GetString(reader)

		assertError(t, err, nil)

		if got != want {
			t.Errorf("GetInput() returns %v, want %v", got, want)
		}
	})
}

func TestBuySingleCoffeeType(t *testing.T) {
	tests := []struct {
		name   string
		coffee coffeeType
		want   *CoffeeMachine
	}{
		{
			"Test success buying espresso",
			coffeeType{250, 0, 16, 4},
			NewCoffeeMachine(150, 540, 104, 8, 554),
		},
		{
			"Test success buying latte",
			coffeeType{350, 75, 20, 7},
			NewCoffeeMachine(50, 465, 100, 8, 557),
		},
		{
			"Test success buying cappuccino",
			coffeeType{200, 100, 12, 6},
			NewCoffeeMachine(200, 440, 108, 8, 556),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coffeeMachine := NewCoffeeMachine(400, 540, 120, 9, 550)
			coffeeMachine.Buy(tt.coffee)

			if !cmp.Equal(coffeeMachine, tt.want, cmp.AllowUnexported(CoffeeMachine{})) {
				t.Error(cmp.Diff(coffeeMachine, tt.want, cmp.AllowUnexported(CoffeeMachine{})))
			}
		})
	}
}

func TestBuyMultipleCoffeeTypes(t *testing.T) {
	coffeeMachine := NewCoffeeMachine(600, 540, 120, 9, 550)
	tests := []struct {
		name   string
		coffee coffeeType
		want   *CoffeeMachine
	}{
		{
			"Test buy 1st coffee a espresso",
			coffeeType{250, 0, 16, 4},
			NewCoffeeMachine(350, 540, 104, 8, 554),
		},
		{
			"Test buy 2nd coffee a latte",
			coffeeType{350, 75, 20, 7},
			NewCoffeeMachine(0, 465, 84, 7, 561),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coffeeMachine.Buy(tt.coffee)
			if !cmp.Equal(coffeeMachine, tt.want, cmp.AllowUnexported(CoffeeMachine{})) {
				t.Error(cmp.Diff(coffeeMachine, tt.want, cmp.AllowUnexported(CoffeeMachine{})))
			}
		})
	}
}

func TestTake(t *testing.T) {
	coffeeMachine := NewCoffeeMachine(600, 540, 120, 9, 550)
	coffeeMachine.Take()
	want := 0
	if coffeeMachine.cash != want {
		t.Errorf("Take returned %d, wanted %d", coffeeMachine.cash, want)
	}
}

func TestFill(t *testing.T) {
	coffeeMachine := NewCoffeeMachine(400, 0, 120, 9, 550)
	resources := coffeeMachineResources{100, 200, 50, 10}
	coffeeMachine.Fill(resources)
	want := NewCoffeeMachine(500, 200, 170, 19, 550)

	if !cmp.Equal(coffeeMachine, want, cmp.AllowUnexported(CoffeeMachine{})) {
		t.Error(cmp.Diff(coffeeMachine, want, cmp.AllowUnexported(CoffeeMachine{})))
	}
}

func TestBuyShouldFailOnLowMachineResources(t *testing.T) {
	tests := []struct {
		name          string
		coffeeMachine CoffeeMachine
		coffee        coffeeType
	}{
		{
			"fail with low on water",
			CoffeeMachine{0, 540, 120, 9, 550},
			coffeeType{250, 0, 16, 4},
		},
		{
			"fail with low on milk",
			CoffeeMachine{500, 0, 120, 9, 550},
			coffeeType{250, 75, 16, 4},
		},
		{
			"fail with low on beans",
			CoffeeMachine{500, 100, 0, 9, 550},
			coffeeType{250, 75, 16, 4},
		},
		{
			"fail with low on beans",
			CoffeeMachine{500, 100, 0, 9, 550},
			coffeeType{250, 75, 16, 4},
		},
		{
			"fail with low on cups",
			CoffeeMachine{500, 100, 100, 0, 550},
			coffeeType{250, 75, 16, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.coffeeMachine.Buy(tt.coffee)
			if !cmp.Equal(tt.coffeeMachine, tt.coffeeMachine, cmp.AllowUnexported(CoffeeMachine{})) {
				t.Error(cmp.Diff(tt.coffeeMachine, tt.coffeeMachine, cmp.AllowUnexported(CoffeeMachine{})))
			}
		})
	}
}

func assertError(t *testing.T, got, want error) {
	t.Helper()
	if !errors.Is(got, want) {
		t.Errorf("Got error %v, but want %v", got, want)
	}
}
