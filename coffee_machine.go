package coffeeMachine

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	coffeeMachineStartingWater = 400
	coffeeMachineStartingMilk  = 540
	coffeeMachineStartingBeans = 120
	coffeeMachineStartingCups  = 9
	coffeeMachineStartingCash  = 550

	wrongInputError      = coffeeMachineError("Wrong input. Please enter one of the available options")
	expectedIntegerError = coffeeMachineError("expected integer")
	notEnoughWater       = coffeeMachineError("Sorry, not enough water!")
	notEnoughMilk        = coffeeMachineError("Sorry, not enough milk!")
	notEnoughBeans       = coffeeMachineError("Sorry, not enough beans!")
	notEnoughCoups       = coffeeMachineError("Sorry, not enough cups!")
)

const (
	Buy       machineMenuOption = "buy"
	Fill      machineMenuOption = "fill"
	Take      machineMenuOption = "take"
	Remaining machineMenuOption = "remaining"
	Exit      machineMenuOption = "exit"

	Espresso   coffeeTypeOption = "1"
	Late       coffeeTypeOption = "2"
	Cappuccino coffeeTypeOption = "3"
)

type machineMenuOption string

type coffeeTypeOption string

type coffeeMachineError string

func (e coffeeMachineError) Error() string {
	return string(e)
}

type coffeeType struct {
	water int
	milk  int
	beans int
	cost  int
}

type coffeeMachineResources struct {
	water int
	milk  int
	beans int
	cups  int
}

type CoffeeMachine struct {
	water int
	milk  int
	beans int
	cups  int
	cash  int
}

func (c *CoffeeMachine) String() string {
	return fmt.Sprintf("The coffee machine has:\n"+
		"%d ml of water\n"+
		"%d ml of milk\n"+
		"%d g of coffee beans\n"+
		"%d disposable cups\n"+
		"$%d of money",
		c.water, c.milk, c.beans, c.cups, c.cash)
}

func NewCoffeeMachine(water, milk, beans, cups, cash int) *CoffeeMachine {
	return &CoffeeMachine{
		water: water,
		milk:  milk,
		beans: beans,
		cups:  cups,
		cash:  cash,
	}
}

func GetString(reader io.Reader) (string, error) {
	var input string
	scanner := bufio.NewScanner(reader)
	if scanner.Scan() {
		_, err := fmt.Sscanf(scanner.Text(), "%s", &input)
		if err != nil {
			return "", scanner.Err()
		}
	}
	return input, nil
}

func GetInt(reader io.Reader) (int, error) {
	var input int
	scanner := bufio.NewScanner(reader)
	if scanner.Scan() {
		_, err := fmt.Sscanf(scanner.Text(), "%d", &input)
		if err != nil {
			return 0, expectedIntegerError
		}
	}
	return input, nil
}

func (c *CoffeeMachine) Buy(coffee coffeeType) {

	if err := c.isActionValid(coffee); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("I have enough resources, making you a coffee!")
	c.water -= coffee.water
	c.milk -= coffee.milk
	c.beans -= coffee.beans
	c.cash += coffee.cost
	c.cups -= 1
}

func (c *CoffeeMachine) isActionValid(coffee coffeeType) error {
	switch {
	case c.water < coffee.water:
		return notEnoughWater
	case c.milk < coffee.milk:
		return notEnoughMilk
	case c.beans < coffee.beans:
		return notEnoughBeans
	case c.cups < 1:
		return notEnoughCoups
	default:
		return nil
	}
}

func (c *CoffeeMachine) Fill(resources coffeeMachineResources) {

	c.water += resources.water
	c.milk += resources.milk
	c.beans += resources.beans
	c.cups += resources.cups
}

func (c *CoffeeMachine) Take() {
	fmt.Printf("I gave you $%d\n", c.cash)
	c.cash = 0
}

func (c *CoffeeMachine) machineMenu() {
	for {
		fmt.Println("Write action (buy, fill, take, remaining, exit): ")
		input, _ := GetString(os.Stdin)
		option := machineMenuOption(input)
		switch option {
		case Buy:
			c.coffeeMenu()
		case Fill:
			fmt.Println("Write how many ml of water you want to add: ")
			water, _ := GetInt(os.Stdin)
			fmt.Println("Write how many ml of milk you want to add: ")
			milk, _ := GetInt(os.Stdin)
			fmt.Println("Write how many grams of coffee beans you want to add: ")
			beans, _ := GetInt(os.Stdin)
			fmt.Println("Write how many disposable cups you want to add: ")
			cups, _ := GetInt(os.Stdin)
			resources := coffeeMachineResources{water, milk, beans, cups}
			c.Fill(resources)
		case Take:
			c.Take()
		case Remaining:
			fmt.Println(c)
		case Exit:
			return
		default:
			fmt.Println(wrongInputError)
		}
		fmt.Println()
	}

}

func (c *CoffeeMachine) coffeeMenu() {
	fmt.Println("What do you want to buy? 1 - espresso, 2 - latte, 3 - cappuccino, back - to main menu: ")
	input, _ := GetString(os.Stdin)
	choice := coffeeTypeOption(input)

	switch choice {
	case Espresso:
		c.Buy(coffeeType{250, 0, 16, 4})
	case Late:
		c.Buy(coffeeType{350, 75, 20, 7})
	case Cappuccino:
		c.Buy(coffeeType{200, 100, 12, 6})
	default:
		return
	}
}

func StartCoffeeMachine() {
	coffeeMachine := NewCoffeeMachine(coffeeMachineStartingWater, coffeeMachineStartingMilk, coffeeMachineStartingBeans, coffeeMachineStartingCups, coffeeMachineStartingCash)

	coffeeMachine.machineMenu()
}
