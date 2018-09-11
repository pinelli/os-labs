package main

import (
	"fmt"
	"time"
)

var rrTime = time.Duration(80 * time.Millisecond)
var rrOneTaskTime = time.Duration(20 * time.Millisecond)
var fcfsTime = time.Duration(20 * time.Millisecond)

type Process struct {
	id       int
	timeLeft time.Duration //ms
}

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

func rrPlanner(tasks chan *Process) {
	for {
		finished := make(chan bool)
		interrupt := make(chan bool)
		select {
		case task := <-tasks:
			go task.execute(interrupt, finished)
			select {
			case <-time.NewTimer(rrTime).C:
				select {
				case interrupt <- true:
					<-interrupt //ok
					tasks <- task
				case <-finished:
				}
				return
			case <-finished:
			}

		default:
			return
		}

	}
}

func rrPlannerTask(tasks chan *Process) {
	for {
		finished := make(chan bool)
		interrupt := make(chan bool)
		select {
		case task := <-tasks:
			go task.execute(interrupt, finished)
			select {
			case <-time.NewTimer(rrTime).C:
				select {
				case interrupt <- true:
					<-interrupt //ok
					tasks <- task
				case <-finished:
				}
				return
			case <-finished:
			}

		default:
			return
		}

	}
}

var lastTask *Process = nil

func fcfsPlanner(tasks chan *Process) {
	for {
		finished := make(chan bool)
		interrupt := make(chan bool)
		var task *Process
		if lastTask != nil {
			task = lastTask
		} else {
			select {
			case task = <-tasks:
			default:
				return
			}
		}

		go task.execute(interrupt, finished)

		select {
		case <-time.NewTimer(fcfsTime).C:
			select {
			case interrupt <- true:
				<-interrupt //ok
				lastTask = task
			case <-finished:
			}
			return
		case <-finished:
			lastTask = nil
		}
	}
}

func planner(interTasks chan *Process, backgrTasks chan *Process) {
	for {
		rrPlannerTask(interTasks)
		fcfsPlanner(backgrTasks)
	}
}

func main() {
	interTasks := make(chan *Process, 100)
	p1 := &Process{0, time.Duration(20 * time.Millisecond)}
	p2 := &Process{1, time.Duration(200 * time.Millisecond)}
	interTasks <- p1
	interTasks <- p2

	backgrTasks := make(chan *Process, 100)
	p3 := &Process{3, time.Duration(20 * time.Millisecond)}
	p4 := &Process{4, time.Duration(200 * time.Millisecond)}
	backgrTasks <- p3
	backgrTasks <- p4

	planner(interTasks, backgrTasks)
}
