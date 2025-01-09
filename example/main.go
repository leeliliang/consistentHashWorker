package main

import (
	"fmt"
	"github.com/leeliliang/consistentHashWorker/worker"
	"go.uber.org/zap"
)

// ExampleTask is a concrete task that implements the Task interface.
type ExampleTask struct {
	Name string
}

// Execute is the method that will be called to execute the task.
func (t *ExampleTask) HandlerMessage() error {
	fmt.Println("Executing task:", t.Name)
	return nil
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	taskNames := []string{"task1", "task2"}
	wp := worker.NewWorkerPool(5, taskNames, logger)

	wp.Start()

	task1 := &ExampleTask{Name: "Task 1"}
	task2 := &ExampleTask{Name: "Task 2"}

	// Add tasks to the WorkerPool
	wp.AddTask(task1, "hashKey1")
	wp.AddTask(task2, "hashKey2")
	// Add tasks and other operations...
	wp.Wait()
}
