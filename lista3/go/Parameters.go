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

	CeoSpeed           = 20
	WorkerSpeed        = 200
	MachineSpeed       = 200
	ClientSpeed        = 300
	ServiceWorkerSpeed = 300

	Workers        = 10
	Clients        = 10
	Machines       = 3
	MaxArgument    = 1000
	ServiceWorkers = 2

	WorkerImpatientDelay  = 200
	ServiceImpatientDelay = 200

	BrokeProbability = 0.5
)
