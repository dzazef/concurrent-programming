package main

import (
	"fmt"
	"math/rand"
)

func randBool() bool {
	return rand.Float32() < 0.5
}

func getWorkerType(t bool) string {
	if t {
		return "patient"
	}
	return "impatient"
}

func opToString(id int) string {
	switch id {
	case Addition:
		return "+"
	case Multiplication:
		return "*"
	}
	return ""
}

func calculate(t task) int {
	switch t.op {
	case Addition:
		return t.arg1 + t.arg2
	case Multiplication:
		return t.arg1 * t.arg2
	}
	return 0
}

func info(message... interface{}) {
	if Mode { fmt.Println(message...) }
}