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
	CeoSpeed        = 10
	WorkerSpeed     = 10
	ClientSpeed     = 10
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
func CEO(taskList chan<- task) {
	for {
		time.Sleep(time.Millisecond * time.Duration(CeoSpeed))
		arg1 := rand.Intn(MaxArgument)
		arg2 := rand.Intn(MaxArgument) + 1
		op := rand.Intn(MaxOperations)
		newTask := task{arg1, arg2, op}
		taskList <- newTask
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

//noinspection GoBoolExpressions
func worker(id int, taskList <-chan task, storage chan<- product) {
	for {
		time.Sleep(time.Millisecond * time.Duration(WorkerSpeed))
		var result int
		var opName string
		task := <-taskList
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
		storageMutex.Lock()
		storage <- product
		storageMutex.Unlock()
		fmt.Println("Worker", id, "created product", product.id, "from task", task.arg1, opName, task.arg2, "=", result)
	}
}

func client(id int, storage <-chan product) {
	for {
		time.Sleep(time.Millisecond * time.Duration(ClientSpeed))
		result := <-storage
		if mode {
			fmt.Println("Client", id, "took product no", result.id, "from storage")
		}
	}
}

func main() {
	taskList := make(chan task, MaxTasks)
	storage := make(chan product, StorageCapacity)
	//done := make(chan bool, 1)

	go CEO(taskList)
	for i := 1; i <= Workers; i++ {
		go worker(i, taskList, storage)
	}
	for i := 1; i <= Clients; i++ {
		go client(i, storage)
	}

	time.Sleep(time.Second * 100)
}
