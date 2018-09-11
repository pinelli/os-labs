package main

import (
	"fmt"
	"time"
)

var rrTime = time.Duration(80 * time.Millisecond)
var fcfsTime = time.Duration(20 * time.Millisecond)

type Process struct {
	id       int
	timeLeft time.Duration //ms
}

//
//var interTasks chan *Process
//var backgrTasks chan *Process

func (p *Process) execute(interrupt chan bool, finished chan<- bool) {
	fmt.Printf("Process %v executes, time: %v\n", p.id, p.timeLeft)

	start := time.Now()
	timer := time.NewTimer(p.timeLeft).C

	select {
	case <-timer:
		fmt.Printf("Process %v finished\n", p.id)
		finished <- true
	case <-interrupt:
		passed := time.Since(start)
		p.timeLeft = p.timeLeft - passed
		if p.timeLeft < 0 {
			p.timeLeft = 0
		}

		fmt.Printf("Process %v interrupted, left: %v\n", p.id, p.timeLeft)
		interrupt <- true
	}
}

func rrPlanner(interTasks chan *Process) {
	finished := make(chan bool)
	interrupt := make(chan bool)

	for {
		select {
		case task := <-interTasks:
			go task.execute(interrupt, finished)
			select {
			case <-time.NewTimer(rrTime).C:
				interrupt <- true
				<-interrupt //ok
				interTasks <- task
				return
			case <-finished:
			}
		default:
			return
		}

	}
}

func planner(interTasks chan *Process, backgrTasks chan *Process) {
	for {
		rrPlanner(interTasks)
		fmt.Println("Planner finished")
	}

	//for {
	//	finished := make(chan bool)
	//
	//	fmt.Println("RR works")
	//	select {
	//	case <-time.NewTimer(miliSec * rrTime * time.Millisecond).C:
	//	case <-finished:
	//	}
	//
	//	fmt.Println("FCFS works")
	//	select {
	//	case <-time.NewTimer(miliSec * fcfsTime * time.Millisecond).C:
	//	case <-finished:
	//	}
	//}
}

func main() {
	interTasks := make(chan *Process, 100)
	backgrTasks := make(chan *Process, 100)

	//p := &Process{0, time.Duration(200)}
	p1 := &Process{0, time.Duration(20 * time.Millisecond)}
	p2 := &Process{1, time.Duration(200 * time.Millisecond)}

	interTasks <- p1
	interTasks <- p2

	//close(interTasks)

	planner(interTasks, backgrTasks)

	//interrupt := make(chan bool)
	//
	//go func() {
	//	time.Sleep(300 * time.Millisecond)
	//	interrupt <- true
	//}()
	//
	//p.execute(interrupt)
	//fmt.Print("LO1L")
}
