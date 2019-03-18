package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type task struct {
	arg1 int
	arg2 int
	op   int
}

type product struct {
	id    uint64
	value int
}

func opToString(id int) string {
	switch id {
	case Addition:
		return "+"
	case Subtraction:
		return "-"
	case Multiplication:
		return "*"
	case Division:
		return "/"
	}
	return ""
}

func calculate(arg1 int, arg2 int, op int) int {
	switch op {
	case Addition:
		return arg1 + arg2
	case Subtraction:
		return arg1 - arg2
	case Multiplication:
		return arg1 * arg2
	case Division:
		return arg1 / arg2
	}
	return 0
}

//noinspection GoBoolExpressions
func CEO(newTasks chan<- task) {
	for {
		//Sleeping for given time
		time.Sleep(time.Millisecond * time.Duration(CeoSpeed))
		//Creating task
		newTask := task{rand.Intn(MaxArgument), rand.Intn(MaxArgument-1) + 1, rand.Intn(MaxOperations)}
		//Sending to task logger
		newTasks <- newTask
		if Mode {
			fmt.Println("CEO made new task:", newTask.arg1, opToString(newTask.op), newTask.arg2)
		}
	}
}

func taskLogger(newTasks <-chan task, deletedTasks <-chan task, loggedTasks chan<- task, taskLog *[]task) {
	for {
		select {
		case newTask := <-newTasks:
			*taskLog = append(*taskLog, newTask)
			loggedTasks <- newTask
		case deleteTask := <-deletedTasks:
			for idx, val := range *taskLog {
				if val == deleteTask {
					*taskLog = append((*taskLog)[:idx], (*taskLog)[idx+1:]...)
				}
			}
		}
	}
}

func productLogger(newProducts <-chan product, boughtProducts <-chan product, loggedProducts chan<- product, productLog *[]product) {
	for {
		select {
		case newProduct := <-newProducts:
			*productLog = append(*productLog, newProduct)
			loggedProducts <- newProduct
		case boughtProduct := <-boughtProducts:
			for idx, val := range *productLog {
				if val == boughtProduct {
					*productLog = append((*productLog)[:idx], (*productLog)[idx+1:]...)
				}
			}
		}
	}
}

//noinspection GoBoolExpressions
func worker(id int, loggedTasks <-chan task, deletedTasks chan<- task, newProducts chan<- product) {
	for {
		time.Sleep(time.Millisecond * time.Duration(WorkerSpeed))
		task := <-loggedTasks
		deletedTasks <- task
		product := product{rand.Uint64(), calculate(task.arg1, task.arg2, task.op)}
		newProducts <- product
		if Mode {
			fmt.Println("Worker", id, "created product", product.id, "from task", task.arg1, opToString(task.op), task.arg2, "=", product.value)
		}
	}
}

//noinspection GoBoolExpressions
func client(id int, loggedProducts <-chan product, boughtProducts chan<- product) {
	for {
		time.Sleep(time.Millisecond * time.Duration(ClientSpeed))
		result := <-loggedProducts
		boughtProducts <- result
		if Mode {
			fmt.Println("Client", id, "took product no", result.id, "from storage")
		}

	}
}

func viewHelp() {
	fmt.Println("Possible commands:")
	fmt.Println("help - view possible commands")
	fmt.Println("talk - switch to talkative mode")
	fmt.Println("tasklist - view tasklist")
	fmt.Println("storage - view storage")
	fmt.Println()
}

func main() {
	newTasks := make(chan task, MaxTasks)
	loggedTasks := make(chan task, MaxTasks)
	deletedTasks := make(chan task, MaxTasks)

	newProducts := make(chan product, StorageCapacity)
	loggedProducts := make(chan product, StorageCapacity)
	boughtProducts := make(chan product, StorageCapacity)

	taskList := make([]task, 0)
	productList := make([]product, 0)

	scanner := bufio.NewScanner(os.Stdin)

	viewHelp()
	fmt.Println("Starting in", Delay, "ms...")
	time.Sleep(time.Millisecond * Delay)
	fmt.Println("Started.")
	fmt.Println()

	go taskLogger(newTasks, deletedTasks, loggedTasks, &taskList)
	go productLogger(newProducts, boughtProducts, loggedProducts, &productList)
	go CEO(newTasks)
	for i := 1; i <= Workers; i++ {
		go worker(i, loggedTasks, deletedTasks, newProducts)
	}
	for i := 1; i <= Clients; i++ {
		go client(i, loggedProducts, boughtProducts)
	}

	for scanner.Scan() {
		switch scanner.Text() {
		case "help":
			viewHelp()
		case "talk":
			Mode = true
		case "tasklist":
			fmt.Println("Active tasks:", taskList)
		case "storage":
			fmt.Println("Storage:", productList)
		}
	}
}
