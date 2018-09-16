package main

import (
	"fmt"
	"os"
	"sync"
)

type Changer struct {
	*sync.WaitGroup
	bank    *Bank
	queries chan ChangerQuery
}

type ChangerQuery struct {
	terminalId int
	change     int
	resp       chan []int
}

func (changer *Changer) run() {
	for msg := range changer.queries {
		//	fmt.Fprintf("%v\n", msg)
		coins := changer.bank.getCoins(msg.change)
		fmt.Fprint(os.Stderr, "Changer: \n     give coins to the terminal #", msg.terminalId, ":")
		fmt.Fprintln(os.Stderr, "	", coins)
		fmt.Fprintln(os.Stderr, "	 bank:")

		changer.bank.print()
		msg.resp <- coins
	}
	//changer.Done()
}

type TerminalQuery struct {
	id    int
	price int
}

type Terminal struct {
	*sync.WaitGroup
	id      int
	changer Changer
	queries chan TerminalQuery
}

func printCoins(coins []int) {
	if coins == nil {
		fmt.Fprintln(os.Stderr, "No coins")
	} else {
		fmt.Fprintln(os.Stderr, coins)
	}
}
func (terminal *Terminal) run() {
	for q := range terminal.queries {
		change := 100 - q.price
		fmt.Fprintln(os.Stderr, "Terminal #", terminal.id, ": \n "+
			"	query #", q.id, "| price -", q.price, "| change - ", change)

		resp := make(chan []int)
		terminal.changer.queries <- ChangerQuery{terminal.id, change, resp}

		coins := <-resp
		fmt.Fprint(os.Stderr, "Terminal #", terminal.id, ":\n	 received coins: ")
		printCoins(coins)
	}
	//close(terminal.changer.queries)
	terminal.Done()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2) //terminals

	//Bank
	coins := []int{1, 2, 5, 10, 25, 50, 100, 200}
	amounts := []int{2, 4, 6, 8, 1, 1, 1, 1}
	bank := &Bank{coins, amounts}
	fmt.Fprintln(os.Stderr, "Init bank:")
	bank.print()

	//Changer
	queries := make(chan ChangerQuery, 100)
	changer := Changer{&wg, bank, queries}
	go changer.run()

	//Terminal1
	term1Queries := make(chan TerminalQuery, 100)
	terminal1 := Terminal{&wg, 1, changer, term1Queries}

	term1Queries <- TerminalQuery{0, 45}
	term1Queries <- TerminalQuery{1, 35}
	term1Queries <- TerminalQuery{2, 55}
	term1Queries <- TerminalQuery{3, 27}
	term1Queries <- TerminalQuery{4, 13}

	close(term1Queries)

	go terminal1.run()

	//Terminal2
	term2Queries := make(chan TerminalQuery, 100)
	terminal2 := Terminal{&wg, 2, changer, term2Queries}
	go terminal2.run()

	term2Queries <- TerminalQuery{10, 18}
	term2Queries <- TerminalQuery{11, 98}
	term2Queries <- TerminalQuery{12, 13}

	close(term2Queries)

	//Finish
	wg.Wait()
	close(queries)
}
