package amqpreceiver

import (
	"context"
	"errors"

	"github.com/syron/nodiniteexporter/nodiniteexporter/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

const (
	defaultTracesQueue  = "otlp_spans"
	defaultMetricsQueue = "otlp_metrics"
	defaultLogsQueue    = "otlp_logs"
	defaultEncoding     = "otlp_proto"
)

var (
	typeStr = component.MustNewType("amqp")
)

var errInvalidConfig = errors.New("config was not a nodinite exporter config")

type amqpReceiverFactory struct{}

// FactoryOption applies changes to kafkaExporterFactory.
type FactoryOption func(factory *amqpReceiverFactory)

// NewFactory creates Kafka receiver factory.
func NewFactory(options ...FactoryOption) receiver.Factory {
	f := &amqpReceiverFactory{}
	for _, o := range options {
		o(f)
	}
	return receiver.NewFactory(
		metadata.Type,
		createDefaultConfig,
		receiver.WithTraces(f.createTracesReceiver, metadata.TracesStability),
		receiver.WithMetrics(f.createMetricsReceiver, metadata.MetricsStability),
		receiver.WithLogs(f.createLogsReceiver, metadata.LogsStability),
	)
}

func createDefaultConfig() component.Config {
	c := &Config{
		Endpoint: "",
	}
	return c
}

func (f *amqpReceiverFactory) createTracesReceiver(
	_ context.Context,
	set receiver.Settings,
	cfg component.Config,
	nextConsumer consumer.Traces,
) (receiver.Traces, error) {
	oCfg := *(cfg.(*Config))
	if oCfg.Address == "" {
		oCfg.Address = defaultTracesQueue
	}

	r, err := newTracesReceiver(oCfg, set, nextConsumer)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (f *amqpReceiverFactory) createMetricsReceiver(
	_ context.Context,
	set receiver.Settings,
	cfg component.Config,
	nextConsumer consumer.Metrics,
) (receiver.Metrics, error) {
	oCfg := *(cfg.(*Config))
	if oCfg.Address == "" {
		oCfg.Address = defaultMetricsQueue
	}

	r, err := newMetricsReceiver(oCfg, set, nextConsumer)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (f *amqpReceiverFactory) createLogsReceiver(
	_ context.Context,
	set receiver.Settings,
	cfg component.Config,
	nextConsumer consumer.Logs,
) (receiver.Logs, error) {
	oCfg := *(cfg.(*Config))
	if oCfg.Address == "" {
		oCfg.Address = defaultLogsQueue
	}

	r, err := newLogsReceiver(oCfg, set, nextConsumer)
	if err != nil {
		return nil, err
	}
	return r, nil
}
