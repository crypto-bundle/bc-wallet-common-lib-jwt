package nats

import (
	"context"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

// jsQueueConsumerWorkerPool is a minimal Worker implementation that simply wraps a
type jsQueueConsumerWorkerPool struct {
	msgChannel chan *nats.Msg

	jsInfo         *nats.StreamInfo
	jsConfig       *nats.StreamConfig
	jsConsumerConn nats.JetStreamContext

	subjectName    string
	streamName     string
	queueGroupName string
	durable        bool

	handler consumerHandler
	workers []*jsConsumerWorkerWrapper

	logger *zap.Logger
}

func (wp *jsQueueConsumerWorkerPool) Init(ctx context.Context) error {
	_, err := wp.jsConsumerConn.ChanQueueSubscribe(wp.subjectName, wp.queueGroupName, wp.msgChannel)
	if err != nil {
		return err
	}

	return nil
}

func (wp *jsQueueConsumerWorkerPool) Run(ctx context.Context) error {
	wp.run()

	return nil
}

func (wp *jsQueueConsumerWorkerPool) run() {
	for _, w := range wp.workers {
		w.msgChannel = wp.msgChannel

		go w.Start()
	}
}

func (wp *jsQueueConsumerWorkerPool) Shutdown(ctx context.Context) error {
	for _, w := range wp.workers {
		w.Stop()
	}

	return nil
}

func NewJsQueueConsumerWorkersPool(logger *zap.Logger,
	msgChannel chan *nats.Msg,
	streamName string,
	subjectName string,
	queueGroupName string,
	workersCount uint16,
	handler consumerHandler,
	jsNatsConn nats.JetStreamContext,
) *jsQueueConsumerWorkerPool {
	l := logger.Named("queue_consumer_pool.service")

	workersPool := &jsQueueConsumerWorkerPool{
		handler:        handler,
		logger:         l,
		msgChannel:     msgChannel,
		subjectName:    subjectName,
		streamName:     streamName,
		queueGroupName: queueGroupName,
		jsConsumerConn: jsNatsConn,
	}

	for i := uint16(0); i < workersCount; i++ {
		ww := &jsConsumerWorkerWrapper{
			msgChannel:       workersPool.msgChannel,
			stopWorkerChanel: make(chan bool),
			handler:          workersPool.handler,
			logger:           l.With(zap.Uint16(WorkerUnitNumberTag, i)),
		}

		workersPool.workers = append(workersPool.workers, ww)
	}

	return workersPool
}
