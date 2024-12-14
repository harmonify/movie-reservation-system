package tracer

import (
	"context"

	"github.com/agiledragon/gomonkey/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type span struct {
}

func (s *span) End(...trace.SpanEndOption) {
}

func (s *span) AddEvent(string, ...trace.EventOption) {
}

func (s *span) IsRecording() bool {
	return false
}

func (s *span) RecordError(error, ...trace.EventOption) {
}

func (s *span) SpanContext() trace.SpanContext {
	return trace.SpanContext{}
}

func (s *span) SetStatus(codes.Code, string) {
}

func (s *span) SetName(string) {
}

func (s *span) SetAttributes(...attribute.KeyValue) {
}

func (s *span) TracerProvider() trace.TracerProvider {
	return nil
}

func GetSpanMock() trace.Span {
	var (
		SpanMock *span
		patches  = gomonkey.NewPatches()
	)

	patches.ApplyMethodReturn(SpanMock, "End")

	return SpanMock
}

func GetTracerMock() *Tracer {
	var (
		TracerMock *Tracer
		patches    = gomonkey.NewPatches()
	)

	patches.ApplyMethodReturn(TracerMock, "Start", context.TODO(), GetSpanMock())

	return TracerMock
}
