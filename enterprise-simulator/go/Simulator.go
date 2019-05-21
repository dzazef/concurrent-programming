package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type task struct {
	arg1	int
	arg2 	int
	op   	int
	value 	int
}

func CEO(newTaskWrite chan<- *taskWrite) {
	for {
		time.Sleep(time.Millisecond * time.Duration(CeoSpeed))
		newTask := task{rand.Intn(MaxArgument), rand.Intn(MaxArgument-1) + 1, rand.Intn(MaxOperations), 0}
		write := &taskWrite{newTask, make(chan bool)}
		newTaskWrite <- write
		<- write.resp
		info("CEO: Made new task:", newTask.arg1, opToString(newTask.op), newTask.arg2)
	}
}

func Machine(taskToDo chan *taskWrite) {
	for {
		newTaskCalculateRequest := <- taskToDo
		info("MACHINE: Got new task ", newTaskCalculateRequest.t.arg1, opToString(newTaskCalculateRequest.t.op), newTaskCalculateRequest.t.arg2)
		newTaskCalculateRequest.t.value = calculate(newTaskCalculateRequest.t)
		time.Sleep(time.Millisecond * time.Duration(MachineSpeed))
		newTaskCalculateRequest.resp <- true
		info("MACHINE: Done task ", newTaskCalculateRequest.t.arg1, opToString(newTaskCalculateRequest.t.op), newTaskCalculateRequest.t.arg2)
	}
}

func Worker(ID int, patient bool, newTaskRead chan<- *taskRead, machines [][]chan *taskWrite, doneTaskWrite chan *taskWrite, increase chan *increaseTasksNumber) {
	for {
		time.Sleep(time.Millisecond * time.Duration(WorkerSpeed))
		//Send read request to TaskManager
		read := &taskRead{make(chan task)}
		newTaskRead <- read
		result := <- read.value
		info("WORKER("+getWorkerType(patient)+"): Taken task ", result.arg1, opToString(result.op), result.arg2, " from taskList")
		//Create calculate request
		calculateRequest := &taskWrite{result, make(chan bool)}
		if patient {
			info("WORKER(patient): Put task ", result.arg1, opToString(result.op), result.arg2, " in machine")
			machines[result.op][rand.Intn(Machines)] <- calculateRequest
			<- calculateRequest.resp
		} else {
			condition := true
			for condition {
				select {
				case machines[result.op][rand.Intn(Machines)] <- calculateRequest:
					info("WORKER(impatient(!)): Put task ", result.arg1, opToString(result.op), result.arg2, " in machine")
					<- calculateRequest.resp
					condition = false
				default:
					info("WORKER(impatient(!)): Looking for another machine with task ", result.arg1, opToString(result.op), result.arg2)
				}
				time.Sleep(time.Millisecond * time.Duration(ImpatientTime))
			}
		}
		//Send to Storage
		write := &taskWrite{calculateRequest.t, make(chan bool)}
		doneTaskWrite <- write
		<-write.resp
		info("WORKER("+getWorkerType(patient)+"): Put task ", result.arg1, opToString(result.op), result.arg2, " in storage")
		//Update number
		set := &increaseTasksNumber{ID, make(chan bool)}
		increase <- set
		<- set.resp
	}
}

func Client(storageRead chan *taskRead) {
	for {
		time.Sleep(time.Millisecond * time.Duration(ClientSpeed))
		read := &taskRead{make(chan task)}
		storageRead <- read
		result := <- read.value
		info("CLIENT: Taken product ", result.arg1, opToString(result.op), result.arg2, " = ", result.value, " from storage.")
	}
}

func viewHelp() {
	fmt.Println("Possible commands:")
	fmt.Println("help - view possible commands")
	fmt.Println("talk - switch to talkative mode")
	fmt.Println("tasklist - view tasklist")
	fmt.Println("storage - view storage")
	fmt.Println("statistic - view statistics for workers")
	fmt.Println("workers - view worker types")
	fmt.Println()
}

func showWorkerType(workerType map[int]bool) {
	for k, v := range workerType {
		fmt.Println("Worker", k, getWorkerType(v))
	}
}

func main() {
	taskListWrite := make(chan *taskWrite)
	taskListRead := make(chan *taskRead)
	requestPrintTaskList := make(chan bool)
	go TaskManager(taskListWrite, taskListRead, requestPrintTaskList)
	go CEO(taskListWrite)

	machineAddArray := make([]chan *taskWrite, 0)
	machineMultiplyArray := make([]chan *taskWrite, 0)
	machineSet := make([][]chan *taskWrite, 0)
	for i := 1; i<= Machines; i++ {
		newAddChannel := make(chan *taskWrite)
		newMultiplyChannel := make(chan *taskWrite)
		machineAddArray = append(machineAddArray, newAddChannel)
		machineMultiplyArray = append(machineMultiplyArray, newMultiplyChannel)
		go Machine(newAddChannel)
		go Machine(newMultiplyChannel)
	}
	machineSet = append(machineSet, machineAddArray)
	machineSet = append(machineSet, machineMultiplyArray)

	storageWrite := make(chan *taskWrite)
	storageRead := make(chan *taskRead)
	requestPrintStorage := make(chan bool)
	go TaskManager(storageWrite, storageRead, requestPrintStorage)

	increaseTaskNumberChan := make(chan *increaseTasksNumber)
	requestPrintStatistic := make(chan bool)
	go StatisticManager(increaseTaskNumberChan, requestPrintStatistic)

	workerType := make(map[int]bool)
	for i := 1; i <= Workers; i++ {
		workerType[i] = randBool()
		go Worker(i, workerType[i], taskListRead, machineSet, storageWrite, increaseTaskNumberChan)
	}

	for i := 1; i <= Clients; i++ {
		go Client(storageRead)
	}

	scanner := bufio.NewScanner(os.Stdin)
	viewHelp()

	for scanner.Scan() {
		switch scanner.Text() {
		case "help":
			viewHelp()
		case "talk":
			Mode = true
		case "tasklist":
			requestPrintTaskList <- true
		case "storage":
			requestPrintStorage <-true
		case "statistic":
			requestPrintStatistic <-true
		case "workers":
			showWorkerType(workerType)
		}
	}
}