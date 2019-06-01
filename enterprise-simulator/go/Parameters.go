package main

//TRUE when Talkative Mode, FALSE when Silent Mode
var (
	Mode = true
)

const (
	Addition       = 0
	Multiplication = 1
	MaxOperations  = 2
)

const (
	Delay    = 2000
	MaxTasks = 255

	CeoSpeed           = 200
	WorkerSpeed        = 2000
	MachineSpeed       = 2000
	ClientSpeed        = 3000
	ServiceWorkerSpeed = 3000

	Workers        = 10
	Clients        = 10
	Machines       = 3
	MaxArgument    = 1000
	ServiceWorkers = 2

	WorkerImpatientDelay  = 200
	ServiceImpatientDelay = 200

	BrokeProbability = 1
)
