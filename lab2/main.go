package main

import (
	"fmt"
	"log"
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
	queryId    int
	change     int
	resp       chan []int
}

func (changer *Changer) run() {
	for msg := range changer.queries {
		//	fmt.Fprintf("%v\n", msg)

		coins := changer.bank.getCoins(msg.change)
		str := fmt.Sprint("Changer: \n     give COINS to the terminal #", msg.terminalId, " (query #", msg.queryId, ")", ":")
		str += fmt.Sprintln("	", coins)
		str += fmt.Sprintln("	 bank:\n", changer.bank.String())
		logger.Print(str)
		//fmt.Fprint(os.Stderr, "Changer: \n     give coins to the terminal #", msg.terminalId, ":")
		//fmt.Fprintln(os.Stderr, "	", coins)
		//fmt.Fprintln(os.Stderr, "	 bank:")

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

func printCoins(coins []int) string {
	if coins == nil {
		//fmt.Fprintln(os.Stderr, "No coins")
		return "no coins"
	} else {
		//fmt.Fprintln(os.Stderr, coins)
		return fmt.Sprintln(coins)
	}
}
func (terminal *Terminal) run() {
	for q := range terminal.queries {
		change := 100 - q.price
		logger.Println("Terminal #", terminal.id, ": \n "+
			"	received QUERY #", q.id, "| price -", q.price, "| change - ", change)

		resp := make(chan []int)
		terminal.changer.queries <- ChangerQuery{terminal.id, q.id, change, resp}

		coins := <-resp
		logger.Println("Terminal #", terminal.id, ":\n	 received COINS for query #", q.id, " :", printCoins(coins))
	}
	//close(terminal.changer.queries)
	terminal.Done()
}

var logger *log.Logger

func main() {
	logger = log.New(os.Stdout, "", 0)

	var wg sync.WaitGroup
	wg.Add(2) //terminals

	//Bank
	coins := []int{1, 2, 5, 10, 25, 50, 100, 200}
	amounts := []int{2, 4, 6, 8, 1, 1, 1, 1}
	bank := &Bank{coins, amounts}
	//fmt.Fprintln(os.Stderr, "Init bank:")
	logger.Println("Init bank:")
	fmt.Println(bank.String())

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
