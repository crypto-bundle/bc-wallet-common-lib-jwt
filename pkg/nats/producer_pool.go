package nats

import (
	"context"
	"sync/atomic"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

// producerWorkerPool is a minimal Worker implementation that simply wraps a
type producerWorkerPool struct {
	logger *zap.Logger

	msgChannel chan *nats.Msg

	subjectName string
	groupName   string

	natsProducerConn *nats.Conn
	workers          []*producerWorkerWrapper

	workersCount uint32
	rr           uint32 // round-robin index
}

func (wp *producerWorkerPool) Init(ctx context.Context) error {

	return nil
}

func (wp *producerWorkerPool) Run(ctx context.Context) error {
	wp.run()

	return nil
}

func (wp *producerWorkerPool) run() {
	for i, _ := range wp.workers {
		go wp.workers[i].Run()
	}
}

func (wp *producerWorkerPool) Shutdown(ctx context.Context) error {
	for _, w := range wp.workers {
		w.Stop()
	}

	return nil
}

func (wp *producerWorkerPool) Produce(ctx context.Context, msg *nats.Msg) {
	wp.msgChannel <- msg
}

func (wp *producerWorkerPool) ProduceSync(ctx context.Context, msg *nats.Msg) error {
	n := atomic.AddUint32(&wp.rr, 1)
	return wp.workers[n%wp.workersCount].PublishMsg(msg)
}

func NewProducerWorkersPool(
	logger *zap.Logger,
	workersCount uint16,
	subjectName string,
	groupName string,
	natsProducerConn *nats.Conn,
) (*producerWorkerPool, error) {
	l := logger.Named("producer.service").
		With(zap.String(QueueSubjectNameTag, subjectName))

	msgChannel := make(chan *nats.Msg, workersCount)

	workersPool := &producerWorkerPool{
		logger: l,

		subjectName: subjectName,
		groupName:   groupName,

		msgChannel:       msgChannel,
		natsProducerConn: natsProducerConn,
		workers:          make([]*producerWorkerWrapper, workersCount),
		workersCount:     uint32(workersCount),
		rr:               1, // round-robin index
	}

	for i := uint16(0); i < workersCount; i++ {
		ww := newProducerWorker(logger, i, msgChannel, subjectName,
			natsProducerConn, make(chan bool))

		workersPool.workers[i] = ww
	}

	return workersPool, nil
}
