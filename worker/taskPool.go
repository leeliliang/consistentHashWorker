package worker

import (
	"go.uber.org/zap"
	"sync"
)

type WorkerPool struct {
	poolSize   int
	taskName   []string
	tasks      map[string]chan Task // 每个 hash 值对应一个任务通道
	wg         sync.WaitGroup
	consistent *ConsistentHash
	logger     *zap.Logger
}

func NewWorkerPool(poolSize int, taskName []string, logger *zap.Logger) *WorkerPool {
	tasks := make(map[string]chan Task, poolSize)
	for _, value := range taskName {
		tasks[value] = make(chan Task, 1000)
	}

	consistentIns := New(poolSize, MurmurHash)
	consistentIns.Add(taskName)
	logger.Info("taskName", zap.Any("taskName", taskName))
	return &WorkerPool{
		poolSize:   poolSize,
		tasks:      tasks,
		taskName:   taskName,
		consistent: consistentIns,
		logger:     logger,
	}
}

func (wp *WorkerPool) Start() {
	for index, value := range wp.taskName {
		wp.wg.Add(1)
		wp.logger.Info("start", zap.Int("index", index), zap.String("value", value))
		go wp.worker(value)
	}
}

func (wp *WorkerPool) worker(workerID string) {
	wp.logger.Debug("task", zap.String("workerID", workerID))
	defer wp.wg.Done()
	for task := range wp.tasks[workerID] {
		wp.logger.Debug("task", zap.String("workerID", workerID))
		err := task.HandlerMessage()
		if err != nil {
			wp.logger.Error("Failed to process record message", zap.Error(err))
			continue
		}

	}
}

func (wp *WorkerPool) AddTask(task Task, hashKey string) {
	key := wp.consistent.Get(hashKey)
	wp.logger.Info("msg", zap.String("key", key), zap.String("hashKey", hashKey))
	wp.tasks[key] <- task
}

func (wp *WorkerPool) Wait() {
	for _, ch := range wp.tasks {
		close(ch)
	}
	wp.wg.Wait()
}
