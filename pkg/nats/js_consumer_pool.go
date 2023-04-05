package nats

import (
	"context"
	"errors"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

var (
	ErrReturnedNilConsumerInfo = errors.New("returned nil consumer info")
	ErrReturnedNilStreamInfo   = errors.New("returned nil stream info")
)

// jsConsumerWorkerPool is a minimal Worker implementation that simply wraps a
type jsConsumerWorkerPool struct {
	msgChannel chan *nats.Msg

	jsInfo         *nats.StreamInfo
	jsConfig       *nats.StreamConfig
	jsConsumerConn nats.JetStreamContext

	subjectName string
	streamName  string
	durable     bool

	handler consumerHandler
	workers []*jsConsumerWorkerWrapper

	logger *zap.Logger
}

func (wp *jsConsumerWorkerPool) Init(ctx context.Context) error {
	streamInfo, err := wp.getOrCreateStream(ctx)
	if err != nil {
		return err
	}

	if streamInfo == nil {
		return ErrReturnedNilStreamInfo
	}

	//consumerInfo, err := wp.getOrCreateSubscriber(ctx)
	//if err != nil {
	//	return err
	//}
	//
	//if consumerInfo == nil {
	//	return ErrReturnedNilConsumerInfo
	//}

	_, err = wp.jsConsumerConn.ChanSubscribe(wp.subjectName, wp.msgChannel)
	if err != nil {
		return err
	}

	return nil
}

func (wp *jsConsumerWorkerPool) getOrCreateStream(ctx context.Context) (*nats.StreamInfo, error) {
	streamInfo, err := wp.jsConsumerConn.StreamInfo(wp.streamName)
	if err != nil {
		if errors.Is(err, nats.ErrStreamNotFound) {
			stream, addStreamErr := wp.jsConsumerConn.AddStream(wp.jsConfig)
			if addStreamErr != nil {
				return nil, addStreamErr
			}

			streamInfo = stream
		}

		return nil, err
	}

	return streamInfo, nil
}

func (wp *jsConsumerWorkerPool) getOrCreateSubscriber(ctx context.Context) (*nats.ConsumerInfo, error) {
	consumerInfo, err := wp.jsConsumerConn.ConsumerInfo(wp.streamName, wp.subjectName)
	if err != nil {
		if errors.Is(err, nats.ErrConsumerNotActive) {
			consumerConfig := &nats.ConsumerConfig{
				Durable: wp.subjectName,
			}

			consumer, addErr := wp.jsConsumerConn.AddConsumer(wp.streamName, consumerConfig)
			if err != nil {
				return nil, addErr
			}

			consumerInfo = consumer
		}

		return nil, err
	}

	return consumerInfo, nil
}

func (wp *jsConsumerWorkerPool) Run(ctx context.Context) error {
	wp.run()

	return nil
}

func (wp *jsConsumerWorkerPool) run() {
	for _, w := range wp.workers {
		w.msgChannel = wp.msgChannel

		go w.Start()
	}
}

func (wp *jsConsumerWorkerPool) Shutdown(ctx context.Context) error {
	for _, w := range wp.workers {
		w.Stop()
	}

	return nil
}

func NewJsConsumerWorkersPool(logger *zap.Logger,
	msgChannel chan *nats.Msg,
	streamName string,
	subjectName string,
	workersCount uint16,
	handler consumerHandler,
	jsNatsConn nats.JetStreamContext,
) *jsConsumerWorkerPool {
	l := logger.Named("consumer_pool.service")

	workersPool := &jsConsumerWorkerPool{
		handler:        handler,
		logger:         l,
		msgChannel:     msgChannel,
		subjectName:    subjectName,
		streamName:     streamName,
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
