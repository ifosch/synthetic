package slack

import (
	"github.com/ifosch/synthetic/pkg/synthetic"
)

// IMessageProcessor is an interface to represent message processor
// functions.
type IMessageProcessor interface {
	Name() string
	Run(synthetic.Message)
}

// MessageProcessor is an implementation of IMessageProcessor.
type MessageProcessor struct {
	name          string
	processorFunc func(synthetic.Message)
}

// NewMessageProcessor is the constructor for MessageProcessor.
func NewMessageProcessor(name string, processorFunc func(synthetic.Message)) *MessageProcessor {
	return &MessageProcessor{
		name:          name,
		processorFunc: processorFunc,
	}
}

// Name returns the name of the MessageProcessor.
func (mp *MessageProcessor) Name() string {
	return mp.name
}

// Run executes the MessageProcessor function.
func (mp *MessageProcessor) Run(msg synthetic.Message) {
	mp.processorFunc(msg)
}
