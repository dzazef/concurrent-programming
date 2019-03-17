# concurrent-programming
Tasks for Concurrent Programming course

## Task 1
The goal of the first task is to create a concurrent program, working as a simulator of small enterprise.
There are three levels of the enterprise:
* CEO: create tasks for workers, task contains two arguments and an operator
* Worker: solves tasks and puts them in storage
* Client: buys products from storage

The simulator is working in two modes:
* Talkative mode: every information is being printed on stdout on the run
* Silent mode: simulator is waiting for users commands 
