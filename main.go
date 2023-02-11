package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const NumberOfPizzas = 10

var PizzasMade, PizzasFailed, total int

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

type PizzaOrder struct {
	PizzaNumber int
	Message     string
	Success     bool
}

func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++
	if pizzaNumber <= NumberOfPizzas {
		delay := rand.Intn(5) + 1
		fmt.Printf("Received order #%d!\n", pizzaNumber)

		rnd := rand.Intn(12) + 1
		msg := ""
		success := false

		if rnd < 5 {
			PizzasFailed++
		} else {
			PizzasMade++
		}
		total++

		fmt.Printf("Making pizza #%d. It will take %d seconds\n", pizzaNumber, delay)
		time.Sleep(time.Duration(delay) * time.Second)

		if rnd < 2 {
			fmt.Printf("We ran out of ingredients for pizza #%d", pizzaNumber)
		} else if rnd < 4 {
			fmt.Printf("The cook quit wile making pizza #%d", pizzaNumber)
		} else {
			success = true
			msg = fmt.Sprintf("Pizza order #%d is ready", pizzaNumber)
		}

		p := PizzaOrder{
			Message:     msg,
			PizzaNumber: pizzaNumber,
			Success:     success,
		}

		return &p
	}

	return &PizzaOrder{
		PizzaNumber: pizzaNumber,
	}
}

func pizzeria(pizzaMaker *Producer) {
	var i = 0
	// keep track which pizza we are trying to make
	// run forever until we get	quit signal
	// try to make pizzas

	for {
		color.HiYellow("in cycle %d\n", i)
		currentPizza := makePizza(i)

		if currentPizza != nil {
			i = currentPizza.PizzaNumber
			select {
			case pizzaMaker.data <- *currentPizza:
			case quitChan := <-pizzaMaker.quit:
				close(pizzaMaker.data)
				close(quitChan)

				return
			}
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	color.Cyan("Pizzeria is ready to start.")

	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	go pizzeria(pizzaJob)

	for i := range pizzaJob.data {
		if i.PizzaNumber <= NumberOfPizzas {
			if i.Success {
				color.Green(i.Message)
				color.Green("Order #%d is ready for delivery", i.PizzaNumber)
			} else {
				color.Red(i.Message)
				color.Red("The customer is really mad!")
			}
		} else {
			color.Cyan("Done making pizzas!")
			err := pizzaJob.Close()
			if err != nil {
				color.Red("*** Error closing channel")
			}
		}
	}

	color.Cyan("-------------------------\n Done of the day")

	color.Cyan("We've made %d pizzas and failed to make %d pizzas. Total attempts %d", PizzasMade, PizzasFailed, total)

	switch {
	case PizzasFailed > 9:
		color.Red("It was awful day")
	case PizzasFailed > 6:
		color.Red("It was a bad day")
	case PizzasFailed > 4:
		color.Yellow("It was okay day")
	default:
		color.Green("It was excellent day")
	}
}
