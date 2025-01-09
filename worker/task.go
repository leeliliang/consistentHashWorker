package worker

type Task interface {
	HandlerMessage() error
}
