package telemetry

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/aclgo/grpc-mail/config"
	"github.com/aclgo/grpc-mail/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Provider struct {
	config         *config.Config
	logger         logger.Logger
	TracerProvider trace.TracerProvider
	MeterProvider  metric.MeterProvider
	propagator     propagation.TextMapPropagator
	Shutdown       func()
}

func NewProvider(config *config.Config, logger logger.Logger, attrs ...attribute.KeyValue) *Provider {

	provider := &Provider{
		config: config,
		logger: logger,
	}

	fn := provider.start(attrs...)

	provider.Shutdown = fn

	return provider
}

func (p *Provider) start(attrs ...attribute.KeyValue) func() {

	var (
		tr  sdktrace.SpanExporter
		mr  sdkmetric.Exporter
		err error
	)

	p.propagator = propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})

	if p.config.OtelExporter == "stdout" {
		tr, err = stdouttrace.New()
		if err != nil {
			return nil
		}

		mr, err = stdoutmetric.New(stdoutmetric.WithEncoder(json.NewEncoder(os.Stdout)))
		if err != nil {
			return nil
		}

	} else {
		tr, err = zipkin.New(
			// ctxTracer,
			p.config.TracerExporterURL,
			// grpc.WithTransportCredentials(insecure.NewCredentials()),
			// grpc.WithBlock(),
		)

		if err != nil {

			p.logger.Errorf("start.DialContext: %v", err)
			return nil
		}

		ctx := context.Background()

		expMeter, err := grpc.DialContext(
			ctx,
			p.config.MeterExporterURL,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)

		if err != nil {
			p.logger.Errorf("start.DialContext: %v", err)
			return nil
		}

		mr, err = otlpmetricgrpc.New(context.Background(), otlpmetricgrpc.WithGRPCConn(expMeter))
		if err != nil {
			p.logger.Errorf("otlpmetricgrpc: %v", err)
			return nil
		}

	}

	ctx := context.Background()

	res, err := resource.New(
		ctx,
		resource.WithAttributes(attrs...),
	)

	if err != nil {
		p.logger.Errorf("cannot initialize tracer resource: %v", err)
		return nil
	}

	tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()), sdktrace.WithResource(res), sdktrace.WithBatcher(tr))
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(sdkmetric.NewPeriodicReader(mr)))

	p.TracerProvider = tp
	p.MeterProvider = mp
	p.propagator = propagation.NewCompositeTextMapPropagator()

	otel.SetMeterProvider(mp)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(p.propagator)

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
	}

}
