package main

import (
	"fmt"
	"time"
)

const (
	SUSPENDED = 0
	FINISHED  = 1
)

var rrTime = time.Duration(80 * time.Millisecond)
var rrOneTaskTime = time.Duration(10 * time.Millisecond)
var fcfsTime = time.Duration(20 * time.Millisecond)

type Process struct {
	id       int
	timeLeft time.Duration //ms
}

func (p *Process) execute(interrupt chan bool, tellPlanner chan<- int) {
	fmt.Printf("Process %v executes, time: %v \n", p.id, p.timeLeft)
	start := time.Now()
	select {
	case <-time.NewTimer(p.timeLeft).C:
		p.timeLeft = 0
	case <-interrupt:
		p.timeLeft = p.timeLeft - time.Since(start)
	}
	if p.timeLeft <= 0 {
		fmt.Printf("Process %v finished\n", p.id)
		fmt.Println("-----------------")
		tellPlanner <- FINISHED
	} else {
		fmt.Printf("Process %v suspended, left %v\n", p.id, p.timeLeft)
		fmt.Println("-----------------")
		tellPlanner <- SUSPENDED
	}
}

func rrPlanner(tasks chan *Process) {
	plannerTimer := time.NewTimer(rrTime).C
	for {
		var task *Process
		select {
		case task = <-tasks:
		default: // no more tasks
			return
		}

		processIn := make(chan int)
		interrupt := make(chan bool, 1)
		go task.execute(interrupt, processIn)

		select {
		case <-plannerTimer:
			interrupt <- true //non blocking
			msg := <-processIn
			if msg == SUSPENDED {
				tasks <- task
			}
			return
		case <-time.NewTimer(rrOneTaskTime).C:
			interrupt <- true //non blocking
			msg := <-processIn
			if msg == SUSPENDED {
				tasks <- task
			}
		case <-processIn: //got FINISHED
		}
	}
}

var currentTask *Process = nil

func setCurrentTask(tasks chan *Process) {
	if currentTask == nil {
		select {
		case task := <-tasks:
			currentTask = task
		default:
		}
	}
}

func fcfsPlanner(tasks chan *Process) {
	timer := time.NewTimer(fcfsTime).C
	for {
		if setCurrentTask(tasks); currentTask == nil {
			return
		}
		processIn := make(chan int)
		interrupt := make(chan bool, 1)
		go currentTask.execute(interrupt, processIn)

		select {
		case <-timer:
			interrupt <- true //non blocking
			msg := <-processIn
			if msg == FINISHED {
				currentTask = nil
			}
			return

		case <-processIn: //got FINISHED
			currentTask = nil
		}
	}
}

func planner(interTasks chan *Process, backgrTasks chan *Process) {
	for {
		rrPlanner(interTasks)
		fcfsPlanner(backgrTasks)
	}
}

func main() {
	interTasks := make(chan *Process, 100)
	p1 := &Process{0, time.Duration(60 * time.Millisecond)}
	p2 := &Process{1, time.Duration(200 * time.Millisecond)}
	interTasks <- p1
	interTasks <- p2

	backgrTasks := make(chan *Process, 100)
	p3 := &Process{2, time.Duration(60 * time.Millisecond)}
	p4 := &Process{3, time.Duration(200 * time.Millisecond)}
	backgrTasks <- p3
	backgrTasks <- p4

	planner(interTasks, backgrTasks)
}
