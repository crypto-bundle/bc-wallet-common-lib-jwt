package nats

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

// jsProducerWorkerPool is a minimal Worker implementation that simply wraps a
type jsProducerWorkerPool struct {
	logger *zap.Logger

	msgChannel chan *nats.Msg
	handler    ProducerWorkerTask
	streamName string
	subject    []string

	storage          nats.StorageType
	jsInfo           *nats.StreamInfo
	jsConfig         *nats.StreamConfig
	natsProducerConn nats.JetStreamContext

	workers      []*jsProducerWorkerWrapper
	workersCount uint32
	rr           uint32 // round-robin index
}

func (wp *jsProducerWorkerPool) Init(ctx context.Context) error {
	streamInfo, err := wp.getStreamInfo(ctx)
	if err != nil {
		return err
	}

	if streamInfo == nil {
		return ErrReturnedNilStreamInfo
	}

	for i, _ := range wp.workers {
		wp.workers[i].SetStreamInfo(streamInfo)
	}

	return nil
}

func (wp *jsProducerWorkerPool) getStreamInfo(ctx context.Context) (*nats.StreamInfo, error) {
	streamInfo, err := wp.natsProducerConn.StreamInfo(wp.jsConfig.Name)
	if err != nil {
		if errors.Is(err, nats.ErrStreamNotFound) {
			wp.logger.Error("stream not found", zap.Error(err))
		}

		return nil, err
	}

	return streamInfo, nil
}

func (wp *jsProducerWorkerPool) Run(ctx context.Context) error {
	wp.run()

	return nil
}

func (wp *jsProducerWorkerPool) run() {
	for i, _ := range wp.workers {
		go wp.workers[i].Run()
	}
}

func (wp *jsProducerWorkerPool) Shutdown(ctx context.Context) error {
	for _, w := range wp.workers {
		w.Stop()
	}

	return nil
}

func (wp *jsProducerWorkerPool) Produce(ctx context.Context, msg *nats.Msg) {
	wp.msgChannel <- msg
}

func (wp *jsProducerWorkerPool) ProduceSync(ctx context.Context, msg *nats.Msg) error {
	n := atomic.AddUint32(&wp.rr, 1)
	return wp.workers[n%wp.workersCount].PublishMsg(msg)
}

func NewJsProducerWorkersPool(
	logger *zap.Logger,
	workersCount uint16,
	streamName string,
	subjects []string,
	storage nats.StorageType,
	natsProducerConn nats.JetStreamContext,
) (*jsProducerWorkerPool, error) {
	l := logger.Named("producer.service").
		With(zap.String(QueueStreamNameTag, streamName))

	streamChannel := make(chan *nats.Msg, workersCount)

	jsConfig := &nats.StreamConfig{
		Name:     streamName,
		Subjects: subjects,
		Storage:  storage,
	}

	workersPool := &jsProducerWorkerPool{
		logger:           l,
		msgChannel:       streamChannel,
		jsConfig:         jsConfig,
		streamName:       streamName,
		subject:          subjects,
		storage:          storage,
		natsProducerConn: natsProducerConn,

		workersCount: uint32(workersCount),
		rr:           1, // round-robin index
	}

	for i := uint16(0); i < workersCount; i++ {
		ww := newJsProducerWorker(logger, i, streamChannel, streamName, subjects,
			natsProducerConn, make(chan bool))

		workersPool.workers = append(workersPool.workers, ww)
	}

	return workersPool, nil
}
