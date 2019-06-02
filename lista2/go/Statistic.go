package main

import "fmt"

type increaseTasksNumber struct {
	id 			int
	resp 		chan bool
}

func StatisticManager(increaseRequest chan *increaseTasksNumber, requestTaskList <-chan bool) {
	numberOfTasks := make(map[int]int)
	for {
		select {
		case set := <-increaseRequest:
			numberOfTasks[set.id] = numberOfTasks[set.id] + 1
			set.resp <- true
		case <-requestTaskList:
			fmt.Println("WORKER STATISTIC")
			for k, v := range numberOfTasks {
				fmt.Println("Worker", k, "made", v, "tasks")
			}
		}
	}
}
