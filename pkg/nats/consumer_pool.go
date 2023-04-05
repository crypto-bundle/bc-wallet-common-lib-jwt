package nats

import (
	"context"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

// consumerWorkerPool is a minimal Worker implementation that simply wraps a
type consumerWorkerPool struct {
	natsConn   *nats.Conn
	msgChannel chan *nats.Msg

	subjectName string
	groupName   string

	handler consumerHandler
	workers []*consumerWorkerWrapper

	logger *zap.Logger
}

func (wp *consumerWorkerPool) Init(ctx context.Context) error {
	_, err := wp.natsConn.ChanQueueSubscribe(wp.subjectName, wp.groupName, wp.msgChannel)
	if err != nil {
		return err
	}

	return nil
}

func (wp *consumerWorkerPool) Run(ctx context.Context) error {
	wp.run()

	return nil
}

func (wp *consumerWorkerPool) run() {
	for _, w := range wp.workers {
		w.msgChannel = wp.msgChannel

		go w.Start()
	}
}

func (wp *consumerWorkerPool) Shutdown(ctx context.Context) error {
	for _, w := range wp.workers {
		w.Stop()
	}

	return nil
}

func NewConsumerWorkersPool(logger *zap.Logger,
	msgChannel chan *nats.Msg,
	subjectName string,
	queueGroupName string,
	workersCount uint16,
	handler consumerHandler,
	natsConn *nats.Conn,
) *consumerWorkerPool {
	l := logger.Named("consumer_pool.service")

	workersPool := &consumerWorkerPool{
		handler:     handler,
		logger:      l,
		msgChannel:  msgChannel,
		subjectName: subjectName,
		groupName:   queueGroupName,
		natsConn:    natsConn,
	}

	for i := uint16(0); i < workersCount; i++ {
		ww := &consumerWorkerWrapper{
			msgChannel:       workersPool.msgChannel,
			stopWorkerChanel: make(chan bool),
			handler:          workersPool.handler,
			logger:           l.With(zap.Uint16(WorkerUnitNumberTag, i)),
		}

		workersPool.workers = append(workersPool.workers, ww)
	}

	return workersPool
}
