package queue

import "context"

type Queue interface {
	Push(ctx context.Context, queue string, msg Message) error
	Pop(ctx context.Context, queue string) (<-chan Message, error)
}

type Message struct {
	ID          string
	ContentType string
	Body        []byte
	AckFunc     func() error
	NackFunc    func() error
}

func (m *Message) Ack() error {
	return m.AckFunc()
}

func (m *Message) Nack() error {
	return m.NackFunc()
}
