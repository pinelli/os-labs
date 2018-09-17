package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

const terminalsNum = 2

var logger *log.Logger

//Changer is a process that gives a change
type Changer struct {
	bank    *Bank
	queries chan ChangerQuery
}

type ChangerQuery struct {
	terminalId int
	queryId    int
	change     int
	resp       chan []int
}

func newChanger() Changer {
	//Bank
	coins := []int{1, 2, 5, 10, 25, 50, 100, 200}
	amounts := []int{2, 4, 6, 8, 1, 1, 1, 1}
	bank := &Bank{coins, amounts}
	logger.Println("Init bank:")
	fmt.Println(bank.String())

	//Changer
	queries := make(chan ChangerQuery, 100)
	return Changer{bank, queries}
}

func (changer *Changer) run() {
	for msg := range changer.queries {
		coins := changer.bank.getCoins(msg.change)
		str := fmt.Sprint("Changer: \n     give COINS to the terminal #", msg.terminalId, " (query #", msg.queryId, ")", ":")
		str += fmt.Sprintln("	", coins)
		str += fmt.Sprintln("	 bank:\n", changer.bank.String())
		logger.Print(str)
		msg.resp <- coins
	}
}

type Terminal struct {
	*sync.WaitGroup
	id      int
	changer Changer
	queries chan Query
}

type Query struct {
	id    int
	price int
}

func newTerminal(id int, changer Changer, group *sync.WaitGroup) Terminal {
	queries := make(chan Query, 100)
	return Terminal{group, id, changer, queries}
}

func (terminal *Terminal) run() {
	for q := range terminal.queries {
		change := 100 - q.price
		logger.Println("Terminal #", terminal.id, ": \n "+
			"	received QUERY #", q.id, "| price -", q.price, "| change - ", change)

		resp := make(chan []int)
		terminal.changer.queries <- ChangerQuery{terminal.id, q.id, change, resp}

		coins := <-resp
		logger.Println("Terminal #", terminal.id, ":\n	 received COINS for query #", q.id, " :", coinsToStr(coins))
	}
	terminal.Done()
}

func coinsToStr(coins []int) string {
	if coins == nil {
		return "no coins"
	} else {
		return fmt.Sprintln(coins)
	}
}

func main() {
	logger = log.New(os.Stdout, "", 0)

	var terminalGroup sync.WaitGroup
	terminalGroup.Add(terminalsNum)

	changer := newChanger()
	go changer.run()

	terminal1 := newTerminal(1, changer, &terminalGroup)

	terminal1.queries <- Query{0, 45}
	terminal1.queries <- Query{1, 35}
	terminal1.queries <- Query{2, 55}
	terminal1.queries <- Query{3, 27}
	terminal1.queries <- Query{4, 13}
	close(terminal1.queries)

	go terminal1.run()

	//Terminal2

	terminal2 := newTerminal(2, changer, &terminalGroup)

	terminal2.queries <- Query{10, 18}
	terminal2.queries <- Query{11, 98}
	terminal2.queries <- Query{12, 13}
	close(terminal2.queries)

	go terminal2.run()

	//Finish
	terminalGroup.Wait()
	close(changer.queries)
}
