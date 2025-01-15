package mocks

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/embedded"
)

type mockSpan struct {
	embedded.Span
}

func (s *mockSpan) span() {
}

func (s *mockSpan) AddLink(link trace.Link) {
}

func (s *mockSpan) End(...trace.SpanEndOption) {
}

func (s *mockSpan) AddEvent(string, ...trace.EventOption) {
}

func (s *mockSpan) IsRecording() bool {
	return false
}

func (s *mockSpan) RecordError(error, ...trace.EventOption) {
}

func (s *mockSpan) SpanContext() trace.SpanContext {
	return trace.SpanContext{}
}

func (s *mockSpan) SetStatus(codes.Code, string) {
}

func (s *mockSpan) SetName(string) {
}

func (s *mockSpan) SetAttributes(...attribute.KeyValue) {
}

func (s *mockSpan) TracerProvider() trace.TracerProvider {
	return nil
}

func NewMockSpan() trace.Span {
	var (
		SpanMock *mockSpan
	)

	return SpanMock
}
