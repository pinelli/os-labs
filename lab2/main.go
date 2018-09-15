package main

import (
	"fmt"
	"sync"
)

type Changer struct {
	*sync.WaitGroup
	bank    *Bank
	queries chan ChangerQuery
}

type ChangerQuery struct {
	change int
	resp   chan []int
}

func (changer *Changer) run() {
	for msg := range changer.queries {
		fmt.Printf("%v\n", msg)
		coins := changer.bank.getCoins(msg.change)
		msg.resp <- coins
		changer.bank.print()
	}
	changer.Done()
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
		fmt.Println("Terminal[term N]: no change")
	} else {
		fmt.Println("Giving a change to the query [N] :", coins)
	}
}
func (terminal *Terminal) run() {
	for q := range terminal.queries {
		change := 100 - q.price
		fmt.Println("Terminal (", terminal.id, "): \n "+
			"	query #", q.id, "| price -", q.price, "| change - ", change)

		resp := make(chan []int)
		terminal.changer.queries <- ChangerQuery{q.price, resp}

		coins := <-resp
		printCoins(coins)
	}
	close(terminal.changer.queries)
	terminal.Done()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2) //terminals

	coins := []int{1, 2, 5, 10, 25, 50, 100, 200}
	amounts := []int{2, 4, 6, 8, 1, 1, 1, 1}
	bank := &Bank{coins, amounts}

	queries := make(chan ChangerQuery, 100)
	changer := Changer{&wg, bank, queries}
	go changer.run()

	termQueries := make(chan TerminalQuery, 100)
	terminal := Terminal{&wg, 0, changer, termQueries}
	go terminal.run()

	termQueries <- TerminalQuery{0, 45}
	close(termQueries)
	//
	//term1Response := make(chan []int)
	//queries <- ChangerQuery{50, term1Response}
	//
	//term2Response := make(chan []int)
	//queries <- ChangerQuery{70, term2Response}
	//
	//term3Response := make(chan []int)
	//queries <- ChangerQuery{0, term3Response}
	//
	//r1 := <-term1Response
	//r2 := <-term2Response
	//r3 := <-term3Response
	//fmt.Println("resp", r1)
	//fmt.Println("resp", r2)
	//fmt.Println("resp", r3)
	//close(queries)
	wg.Wait()
}
