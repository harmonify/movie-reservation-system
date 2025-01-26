package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/failsafe-go/failsafe-go"
	"github.com/failsafe-go/failsafe-go/retrypolicy"
	"github.com/failsafe-go/failsafe-go/timeout"
	"go.uber.org/zap"
)

type DLQError struct {
	Error   error
	RouteID string
}

const DLQOriginalHeaderKey = "dlq_original"
const DLQErrorHeaderKey = "dlq_errors"

type DLQOriginalHeader struct {
	Headers        []*DLQOriginalHeaderHeader `json:"headers,omitempty"`
	Timestamp      time.Time                  `json:"timestamp,omitempty"`
	BlockTimestamp time.Time                  `json:"block_timestamp,omitempty"`
	Key            []byte                     `json:"key,omitempty"`
	Topic          string                     `json:"topic,omitempty"`
	Partition      int32                      `json:"partition,omitempty"`
	Offset         int64                      `json:"offset,omitempty"`
}

type DLQOriginalHeaderHeader struct {
	Key   []byte `json:"key,omitempty"`
	Value []byte `json:"value,omitempty"`
}

type DLQErrorHeader struct {
	RouteID         string `json:"route_id,omitempty"`
	ErrorType       string `json:"error_type,omitempty"`
	ErrorMessage    string `json:"error_message,omitempty"`
	ErrorStackTrace string `json:"error_stack_trace,omitempty"`
}

type DLQHeaders struct {
	Original DLQOriginalHeader `json:"original,omitempty"`
	Errors   []DLQErrorHeader  `json:"errors,omitempty"`
}

// KafkaDLQProducer wraps a KafkaProducer and DLQ handling logic.
type KafkaDLQProducer struct {
	producer                 *KafkaProducer
	moveMessageToDLQExecutor failsafe.Executor[any]
}

func NewKafkaDLQProducer(p KafkaProducerParam) (*KafkaDLQProducer, error) {
	kp, err := NewKafkaProducer(p)
	if err != nil {
		return nil, err
	}
	return &KafkaDLQProducer{
		producer: kp,
		moveMessageToDLQExecutor: failsafe.NewExecutor(
			retrypolicy.Builder[any]().
				WithBackoff(100*time.Millisecond, time.Second).
				WithJitterFactor(0.2).
				WithMaxRetries(3).
				Build(),
			timeout.With[any](5*time.Second),
		),
	}, nil
}

func (dlq *KafkaDLQProducer) MoveMessageToDLQ(ctx context.Context, message *sarama.ConsumerMessage, dlqErrors []DLQError) {
	ctx, span := dlq.producer.tracer.StartSpanWithCaller(ctx)
	defer span.End()

	headers, err := constructDlqHeaders(message, dlqErrors)
	if err != nil {
		dlq.producer.logger.WithCtx(ctx).Error("Failed to construct DQL message headers", zap.Error(err))
	}

	// Send the message to the DLQ topic with original value and updated headers
	err = dlq.moveMessageToDLQExecutor.Run(func() error {
		return dlq.producer.SendMessage(ctx, &sarama.ProducerMessage{
			Topic:   message.Topic + ".dlq",
			Key:     sarama.ByteEncoder(message.Key),
			Value:   sarama.ByteEncoder(message.Value), // Original message body
			Headers: headers,
		})
	})
	if err != nil {
		dlq.producer.logger.WithCtx(ctx).Error("Failed to send message to DLQ", zap.Error(err))
	}
}

func constructDlqHeaders(message *sarama.ConsumerMessage, dlqErrors []DLQError) ([]sarama.RecordHeader, error) {
	// Serialize the original headers into JSON
	originalHeadersBytes, err := json.Marshal(DLQOriginalHeader{
		Headers:        convertHeaders(message.Headers),
		Timestamp:      message.Timestamp,
		BlockTimestamp: message.BlockTimestamp,
		Key:            message.Key,
		Topic:          message.Topic,
		Partition:      message.Partition,
		Offset:         message.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to serialize original headers: %w", err)
	}

	// Prepare error details as a JSON string for headers
	var errorDetails []DLQErrorHeader
	for _, err := range dlqErrors {
		errorDetails = append(errorDetails, DLQErrorHeader{
			RouteID:      err.RouteID,
			ErrorType:    reflect.ValueOf(err.Error).Type().Name(),
			ErrorMessage: err.Error.Error(),
			// ErrorStackTrace: stackTrace(3, 5),
		})
	}
	errorDetailsBytes, err := json.Marshal(errorDetails)
	if err != nil {
		return nil, err
	}

	return []sarama.RecordHeader{
		{Key: []byte(DLQOriginalHeaderKey), Value: originalHeadersBytes},
		{Key: []byte(DLQErrorHeaderKey), Value: errorDetailsBytes},
	}, nil
}

func convertHeaders(headers []*sarama.RecordHeader) []*DLQOriginalHeaderHeader {
	var convertedHeaders []*DLQOriginalHeaderHeader
	for _, header := range headers {
		convertedHeaders = append(convertedHeaders, &DLQOriginalHeaderHeader{
			Key:   header.Key,
			Value: header.Value,
		})
	}
	return convertedHeaders
}

func stackTrace(skip int, maxLines int) string {
	var stacks []string
	for i := 0; i < maxLines; i++ {
		pc, path, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		stacks = append(stacks, fmt.Sprintf("%s:%d %s()", path, line, fn.Name()))
	}
	return strings.Join(stacks, "\n")
}

func (dlq *KafkaDLQProducer) ExtractDLQHeaders(message *sarama.ConsumerMessage) (*DLQHeaders, error) {
	extracted := &DLQHeaders{}

	for _, header := range message.Headers {
		switch string(header.Key) {
		case DLQOriginalHeaderKey:
			// Deserialize original headers
			var original DLQOriginalHeader
			if err := json.Unmarshal(header.Value, &original); err != nil {
				return nil, fmt.Errorf("failed to deserialize original headers: %w", err)
			}
			extracted.Original = original

		case DLQErrorHeaderKey:
			// Deserialize error details
			var errors []DLQErrorHeader
			if err := json.Unmarshal(header.Value, &errors); err != nil {
				return nil, fmt.Errorf("failed to deserialize error details: %w", err)
			}
			extracted.Errors = errors
		}
	}

	return extracted, nil
}
