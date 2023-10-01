package telemetry

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/aclgo/grpc-mail/config"
	"github.com/aclgo/grpc-mail/pkg/logger"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/credentials/insecure"
)

type Provider struct {
	config     *config.Config
	logger     logger.Logger
	tracer     trace.Tracer
	meter      metric.Meter
	propagator propagation.TextMapPropagator
	Shutdown   func()
}

func NewProvider(config *config.Config, logger logger.Logger, attrs ...attribute.KeyValue) (*Provider, error) {

	provider := &Provider{
		config: config,
		logger: logger,
	}

	fn, err := provider.start(attrs...)
	if err != nil {
		return nil, err
	}

	provider.Shutdown = fn

	return provider, nil
}

func (p *Provider) Logger() logger.Logger {
	return p.logger
}
func (p *Provider) Tracer() trace.Tracer {
	return p.tracer
}
func (p *Provider) Meter() metric.Meter {
	return p.meter
}
func (p *Provider) Propagator() propagation.TextMapPropagator {
	return p.propagator
}

func (p *Provider) start(attrs ...attribute.KeyValue) (func(), error) {

	var (
		tr  sdktrace.SpanExporter
		mr  sdkmetric.Exporter
		err error
	)

	p.propagator = propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})

	switch exporter := strings.TrimSpace(strings.ToLower(p.config.OtelExporter)); {
	case exporter == "stdout":
		tr, err = stdouttrace.New()
		if err != nil {
			return nil, fmt.Errorf("stdouttrace: %v", err)
		}

		mr, err = stdoutmetric.New(stdoutmetric.WithEncoder(json.NewEncoder(os.Stdout)))
		if err != nil {
			return nil, fmt.Errorf("stdoutmetric: %v", err)
		}

	case exporter == "otlp":
		tr, err = otlptracegrpc.New(context.Background(), otlptracegrpc.WithTLSCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("otlptracergrpc: %v", err)
		}
		mr, err = otlpmetricgrpc.New(context.Background(), otlpmetricgrpc.WithTLSCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("otlpmetricgrpc: %v", err)
		}
	default:
		p.tracer = trace.NewNoopTracerProvider().Tracer(p.config.Tracer.Name)
		p.meter = noop.NewMeterProvider().Meter(p.config.Tracer.Name)
		return func() {}, nil
	}

	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(attrs...),
	)

	if err != nil {
		return nil, fmt.Errorf("cannot initialize tracer resource: %v", err)
	}

	tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()), sdktrace.WithResource(res), sdktrace.WithBatcher(tr))
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(sdkmetric.NewPeriodicReader(mr)))

	p.tracer = tp.Tracer(p.config.Tracer.Name)
	p.meter = mp.Meter(p.config.Tracer.Name)
	p.propagator = propagation.NewCompositeTextMapPropagator()

	return func() {
		haltctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		var w sync.WaitGroup
		w.Add(2)

		go func() {
			defer w.Done()
			if err := tp.Shutdown(haltctx); err != nil {
				p.logger.Errorf("telemetry tracer shutdown: %v", err)
			}
		}()

		go func() {
			defer w.Done()
			if err := mp.Shutdown(haltctx); err != nil {
				p.logger.Errorf("telemetry meter shutdown: %v", err)
			}
		}()

		w.Wait()
	}, nil

}
