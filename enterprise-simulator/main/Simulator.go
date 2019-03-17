package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

//TRUE when Talkative Mode, FALSE when Calm Mode
var (
	mode          = true
	taskListMutex = &sync.Mutex{}
	storageMutex  = &sync.Mutex{}
)

//operations
const (
	Addition       = 0
	Subtraction    = 1
	Multiplication = 2
	Division       = 3
	MaxOperations  = 4
)

const (
	MaxTasks        = 255
	StorageCapacity = 255
	CeoSpeed        = 500
	WorkerSpeed     = 1000
	ClientSpeed     = 1000
	Workers         = 2
	Clients         = 2
	MaxArgument     = 1000
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

//noinspection GoBoolExpressions
func CEO(newTasks chan<- task) {
	for {
		time.Sleep(time.Millisecond * time.Duration(CeoSpeed))
		arg1 := rand.Intn(MaxArgument)
		arg2 := rand.Intn(MaxArgument) + 1
		op := rand.Intn(MaxOperations)
		newTask := task{arg1, arg2, op}
		newTasks <- newTask
		var opName string
		if mode {
			switch newTask.op {
			case Addition:
				opName = "+"
			case Subtraction:
				opName = "-"
			case Multiplication:
				opName = "*"
			case Division:
				opName = "/"
			}
			fmt.Println("CEO made new task:", arg1, opName, arg2)
		}
	}
}

func taskLogger(newTasks <-chan task, deletedTasks <-chan task, loggedTasks chan<- task, taskLog *[]task) {
	for {
		taskListMutex.Lock()
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
		taskListMutex.Unlock()
	}
}

func productLogger(newProducts <-chan product, boughtProducts <-chan product, loggedProducts chan<- product, productLog *[]product) {
	for {
		storageMutex.Lock()
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
		storageMutex.Unlock()
	}
}

//noinspection GoBoolExpressions
func worker(id int, loggedTasks <-chan task, deletedTasks chan<- task, newProducts chan<- product) {
	for {
		time.Sleep(time.Millisecond * time.Duration(WorkerSpeed))
		var result int
		var opName string
		task := <-loggedTasks
		deletedTasks <- task
		switch task.op {
		case Addition:
			result = task.arg1 + task.arg2
			opName = "+"
		case Subtraction:
			result = task.arg1 - task.arg2
			opName = "-"
		case Multiplication:
			result = task.arg1 * task.arg2
			opName = "*"
		case Division:
			result = task.arg1 / task.arg2
			opName = "/"
		}
		product := product{rand.Uint64(), result}
		newProducts <- product
		fmt.Println("Worker", id, "created product", product.id, "from task", task.arg1, opName, task.arg2, "=", result)
	}
}

//noinspection GoBoolExpressions
func client(id int, loggedProducts <-chan product, boughtProducts chan<- product) {
	for {
		time.Sleep(time.Millisecond * time.Duration(ClientSpeed))
		result := <-loggedProducts
		boughtProducts <- result
		if mode {
			fmt.Println("Client", id, "took product no", result.id, "from loggedProducts")
		}

	}
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

	go taskLogger(newTasks, deletedTasks, loggedTasks, &taskList)
	go productLogger(newProducts, boughtProducts, loggedProducts, &productList)
	go CEO(newTasks)
	for i := 1; i <= Workers; i++ {
		go worker(i, loggedTasks, deletedTasks, newProducts)
	}
	for i := 1; i <= Clients; i++ {
		go client(i, loggedProducts, boughtProducts)
	}

	done := make(chan bool, 1)
	<-done
}
