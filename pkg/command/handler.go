package command

import (
	"fmt"
	"log"
	"sync"

	"github.com/ifosch/synthetic/pkg/synthetic"
)

// ExecutorFunc is the signature of any Command Executor
type ExecutorFunc func(*Command)

// Handler routes the individual Command instances to execution
type Handler struct {
	inventory map[string]ExecutorFunc
}

// NewHandler returns a default Handler
func NewHandler() *Handler {
	return &Handler{
		inventory: make(map[string]ExecutorFunc),
	}
}

// Register adds an Executor with a name to the existing Handler
func (c *Handler) Register(name string, executor ExecutorFunc) error {
	if _, ok := c.inventory[name]; ok {
		return fmt.Errorf("command already registered under `%s` name", name)
	}
	c.inventory[name] = executor
	return nil
}

// Dispatch routes a Command through all registered Executors
func (c *Handler) Dispatch(command *Command) {
	var wg sync.WaitGroup
	for name, executor := range c.inventory {
		wg.Add(1)
		log.Printf("Invoking processor %v", name)
		go func(executor ExecutorFunc) {
			executor(command)
			wg.Done()
		}(executor)
	}
	wg.Wait()
}

// ParseMessage creates a Command from a synthetic.Message
func (c *Handler) ParseMessage(message synthetic.Message) (*Command, error) {
	return NewCommand(message), nil
}

// EventLoop runs an infinite loop that reads messages from a channel
// of synthetic.Message and calls Dispatch on each synthetic.Message
func (c *Handler) EventLoop(messageChannel chan (synthetic.Message)) {
	for message := range messageChannel {
		command, err := c.ParseMessage(message)
		if err != nil {
			log.Printf(
				"error parsing message: %v for message: %#v",
				err.Error(),
				message,
			)
		}
		c.Dispatch(command)
	}
}
