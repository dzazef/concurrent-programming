package main

import "fmt"

type taskWrite struct {
	t    *task
	resp chan bool
}

type taskRead struct {
	value chan *task
}

func taskWriteGuard(condition bool, c <-chan *taskWrite) <-chan *taskWrite {
	if condition {
		return c
	}
	return nil
}
func taskReadGuard(condition bool, c <-chan *taskRead) <-chan *taskRead {
	if condition {
		return c
	}
	return nil
}

func TaskManager(newTaskWrite chan *taskWrite, newTaskRead chan *taskRead, requestTaskList <-chan bool) {
	taskList := make([]*task, 0)
	for {
		select {
		case read := <-taskReadGuard(len(taskList) > 0, newTaskRead):
			task := (taskList)[0]
			taskList = (taskList)[1:]
			read.value <- task
		case write := <-taskWriteGuard(len(taskList) < MaxTasks, newTaskWrite):
			taskList = append(taskList, write.t)
			write.resp <- true
		case <-requestTaskList:
			fmt.Println("INFO", taskList)
		}
	}
}
