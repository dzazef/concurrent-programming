package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type task struct {
	arg1  int
	arg2  int
	op    int
	value int
}

type CEO struct {
	newTaskWrite chan<- *taskWrite
}

type Machine struct {
	ID       int
	taskToDo chan *taskWrite
	backdoor chan bool
	working  bool
	inRepair bool
}

type Worker struct {
	ID            int
	patient       bool
	newTaskRead   chan<- *taskRead
	machines      [][]*Machine
	doneTaskWrite chan *taskWrite
	increase      chan *increaseTasksNumber
	service       Service
}

type ServiceWorker struct {
	ID      int
	request chan *ServiceRepairRequest
}

type Service struct {
	machineToRepair chan *Machine
	workers         []ServiceWorker
}

type ServiceRepairRequest struct {
	machineToRepair *Machine
	done            chan bool
}

func (c *CEO) run() {
	for {
		time.Sleep(time.Millisecond * time.Duration(CeoSpeed))
		newTask := &task{rand.Intn(MaxArgument), rand.Intn(MaxArgument-1) + 1, rand.Intn(MaxOperations), 0}
		write := &taskWrite{newTask, make(chan bool)}
		c.newTaskWrite <- write
		<-write.resp
		info("CEO: Made new task:", newTask.arg1, opToString(newTask.op), newTask.arg2)
	}
}

func (m *Machine) run() {
	m.working = true
	for {
		select {
		case newTaskCalculateRequest := <-m.taskToDo:
			info("MACHINE", m.ID, ": Got new task ", newTaskCalculateRequest.t.arg1, opToString(newTaskCalculateRequest.t.op), newTaskCalculateRequest.t.arg2)

			if m.working && (rand.Float64() < BrokeProbability) { //Breaking
				m.working = false
				info("MACHINE", m.ID, ": Broken")
			}

			if m.working { //If broken set value to -1
				newTaskCalculateRequest.t.value = calculate(*newTaskCalculateRequest.t)
			} else {
				newTaskCalculateRequest.t.value = -1
			}
			time.Sleep(time.Millisecond * time.Duration(MachineSpeed))
			newTaskCalculateRequest.resp <- true
			info("MACHINE", m.ID, ": Done task ", newTaskCalculateRequest.t.arg1, opToString(newTaskCalculateRequest.t.op), newTaskCalculateRequest.t.arg2)
		case <-m.backdoor:
			info("MACHINE", m.ID, ": Repaired!")
			m.working = true
			m.inRepair = false
		}
	}
}

func (w *Worker) run() {
	for {
		time.Sleep(time.Millisecond * time.Duration(WorkerSpeed))
		//Send read request to TaskManager
		read := &taskRead{make(chan *task)}
		w.newTaskRead <- read
		result := <-read.value
		info("WORKER(", w.ID, getWorkerType(w.patient), "): Taken task ", result.arg1, opToString(result.op), result.arg2, " from taskList")
		//Create calculate request
		calculateRequest := &taskWrite{result, make(chan bool)}
		if w.patient {
			for {
				machine := w.machines[result.op][rand.Intn(Machines)]
				info("WORKER(", w.ID, "patient): Put task ", result.arg1, opToString(result.op), result.arg2, " in Machine", machine.ID)
				machine.taskToDo <- calculateRequest
				<-calculateRequest.resp

				if result.value != -1 { //If result is correct, break
					break
				} else { //If it's incorrect, put machine in Service, and do task again
					info("WORKER(", w.ID, "patient): Machine is broken, I wil put it in Service and try do the task again")
					w.service.machineToRepair <- machine
				}
			}
		} else {
			for {
				var machine *Machine
				condition := true
				for condition {
					machine = w.machines[result.op][rand.Intn(Machines)]
					select {
					case machine.taskToDo <- calculateRequest:
						info("WORKER(", w.ID, "impatient(!)): Put task ", result.arg1, opToString(result.op), result.arg2, " in Machine", machine.ID)
						<-calculateRequest.resp
						condition = false
					case <-time.After(time.Millisecond * time.Duration(WorkerImpatientDelay)):
						info("WORKER(", w.ID, "impatient(!)): Looking for another Machine with task ", result.arg1, opToString(result.op), result.arg2)
					}
				}
				if result.value != -1 {
					break
				} else {
					info("WORKER(", w.ID, "patient): Machine is broken, I wil put it in Service and try do the task again")
					w.service.machineToRepair <- machine
				}
			}
		}
		//Send to Storage
		write := &taskWrite{calculateRequest.t, make(chan bool)}
		w.doneTaskWrite <- write
		<-write.resp
		info("WORKER(", w.ID, getWorkerType(w.patient), "): Put task ", result.arg1, opToString(result.op), result.arg2, " in storage")
		//Update number
		set := &increaseTasksNumber{w.ID, make(chan bool)}
		w.increase <- set
		<-set.resp
	}
}

func (s *Service) run() {
	for {
		machine := <-s.machineToRepair //Get machine to repair
		if machine.inRepair {          //If it's already being repaired, ignore
			info("SERVICE: Ignored request to repair machine", machine.ID)
			continue
		} else {
			machine.inRepair = true
			sWorkerID := rand.Intn(ServiceWorkers)
			info("SERVICE: Accepted request to repair machine", machine.ID, "service worker", sWorkerID+1)
			feedback := make(chan bool)
			s.workers[sWorkerID].request <- &ServiceRepairRequest{
				machine,
				feedback,
			}
			<-feedback
			info("SERVICE: Repaired (feedback from ", sWorkerID, ")")
			//condition := true															//Uncomment to make Service impatient
			//for condition {
			//	select {
			//	case s.workers[rand.Intn(ServiceWorkers)].machineToRepair <- machine:
			//		condition = false
			//	case <-time.After(time.Millisecond * time.Duration(ServiceImpatientDelay)):
			//	}
			//}
		}
	}
}

func (sw *ServiceWorker) run() {
	for {
		request := <-sw.request
		info("SERVICE WORKER", sw.ID, ": Repairing machine", request.machineToRepair.ID)
		time.Sleep(time.Millisecond * time.Duration(ServiceWorkerSpeed))
		request.machineToRepair.backdoor <- true
		request.done <- true
	}
}

func Client(storageRead chan *taskRead) {
	for {
		time.Sleep(time.Millisecond * time.Duration(ClientSpeed))
		read := &taskRead{make(chan *task)}
		storageRead <- read
		result := <-read.value
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
	fmt.Println("workers - view Worker types")
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

	CEO := CEO{taskListWrite}
	go CEO.run()

	machineAddArray := make([]*Machine, 0)
	machineMultiplyArray := make([]*Machine, 0)
	machineSet := make([][]*Machine, 0)

	for i := 1; i <= Machines; i++ {
		newMachineAdd := &Machine{
			2 * (i - 1),
			make(chan *taskWrite),
			make(chan bool),
			true,
			false,
		}
		newMachineMultiply := &Machine{
			2*(i-1) + 1,
			make(chan *taskWrite),
			make(chan bool),
			true,
			false,
		}
		machineAddArray = append(machineAddArray, newMachineAdd)
		machineMultiplyArray = append(machineMultiplyArray, newMachineMultiply)
		go newMachineAdd.run()
		go newMachineMultiply.run()
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

	serviceWorkers := make([]ServiceWorker, 0)
	for i := 1; i <= ServiceWorkers; i++ {
		serviceWorker := ServiceWorker{i, make(chan *ServiceRepairRequest)}
		serviceWorkers = append(serviceWorkers, serviceWorker)
		go serviceWorker.run()
	}

	service := Service{
		make(chan *Machine),
		serviceWorkers,
	}
	go service.run()

	workerType := make(map[int]bool)
	for i := 1; i <= Workers; i++ {
		workerType[i] = randBool()
		worker := Worker{
			i,
			workerType[i],
			taskListRead,
			machineSet,
			storageWrite,
			increaseTaskNumberChan,
			service,
		}
		go worker.run()
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
			requestPrintStorage <- true
		case "statistic":
			requestPrintStatistic <- true
		case "workers":
			showWorkerType(workerType)
		}
	}
}
