package loki_logger

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	logger_shared "github.com/harmonify/movie-reservation-system/pkg/logger/shared"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLoki interface {
	Hook(e zapcore.Entry) error
	Sink(u *url.URL) (zap.Sink, error)
	Stop()
	WithCreateLogger(cfg zap.Config, opts ...zap.Option) (*zap.Logger, error)
}

type lokiPusher struct {
	config    *logger_shared.LokiConfig
	ctx       context.Context
	client    *http.Client
	quit      chan struct{}
	entries   chan logEntry
	waitGroup sync.WaitGroup
}

type lokiPushRequest struct {
	Streams []stream `json:"streams"`
}

type stream struct {
	Stream map[string]string `json:"stream"`
	Values [][2]string       `json:"values"`
}

type logEntry struct {
	Level     string  `json:"level"`
	Timestamp float64 `json:"ts"`
	Message   string  `json:"msg"`
	Caller    string  `json:"caller"`
	TraceID   string  `json:"traceId"`
	raw       string
}

var sinkRegistered bool

func NewZapLoki(ctx context.Context, cfg logger_shared.LokiConfig) ZapLoki {
	c := &http.Client{}
	cfg.Url = strings.TrimSuffix(cfg.Url, "/")
	cfg.Url = fmt.Sprintf("%s/loki/api/v1/push", cfg.Url)

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	pusher := &lokiPusher{
		config:  &cfg,
		ctx:     ctx,
		client:  c,
		quit:    make(chan struct{}),
		entries: make(chan logEntry),
	}

	pusher.waitGroup.Add(1)
	go pusher.run()
	return pusher
}

// Hook is a function that can be used as a zap hook to write log lines to loki
func (lp *lokiPusher) Hook(e zapcore.Entry) error {
	lp.entries <- logEntry{
		Level:     e.Level.String(),
		Timestamp: float64(e.Time.UnixMilli()),
		Message:   e.Message,
		Caller:    e.Caller.TrimmedPath(),
	}
	return nil
}

// Sink returns a new loki zap sink
func (lp *lokiPusher) Sink(_ *url.URL) (zap.Sink, error) {
	return newSink(lp), nil
}

// Stop stops the loki pusher
func (lp *lokiPusher) Stop() {
	close(lp.quit)
	lp.waitGroup.Wait()
}

// WithCreateLogger creates a new zap logger with a loki sink from a zap config
func (lp *lokiPusher) WithCreateLogger(cfg zap.Config, opts ...zap.Option) (*zap.Logger, error) {
	if !sinkRegistered {
		err := zap.RegisterSink(lokiSinkKey, lp.Sink)
		if err != nil {
			log.Fatal(err)
		}
		sinkRegistered = true
	}

	fullSinkKey := fmt.Sprintf("%s://", lokiSinkKey)

	if cfg.OutputPaths == nil {
		cfg.OutputPaths = []string{fullSinkKey}
	} else {
		cfg.OutputPaths = append(cfg.OutputPaths, fullSinkKey)
	}

	return cfg.Build(opts...)
}

func (lp *lokiPusher) run() {
	var batch []logEntry
	ticker := time.NewTimer(lp.config.BatchMaxWait)
	defer func() {
		if len(batch) > 0 {
			lp.send(batch)
		}

		lp.waitGroup.Done()
	}()

	for {
		select {
		case <-lp.ctx.Done():
			return
		case <-lp.quit:
			return
		case entry := <-lp.entries:
			batch = append(batch, entry)
			if len(batch) >= lp.config.BatchMaxSize {
				lp.send(batch)
				batch = make([]logEntry, 0)
				ticker.Reset(lp.config.BatchMaxWait)
			}
		case <-ticker.C:
			if len(batch) > 0 {
				lp.send(batch)
				batch = make([]logEntry, 0)
			}
			ticker.Reset(lp.config.BatchMaxWait)
		}
	}
}

func (lp *lokiPusher) send(batch []logEntry) error {
	data := lokiPushRequest{}
	var logs [][2]string
	for _, entry := range batch {
		ts := time.Unix(int64(entry.Timestamp), 0)
		v := [2]string{strconv.FormatInt(ts.UnixNano(), 10), entry.raw}
		logs = append(logs, v)

		// Add custom label to loki
		lp.config.Labels["level"] = entry.Level
		lp.config.Labels["caller"] = entry.Caller
	}

	data.Streams = append(data.Streams, stream{
		Stream: lp.config.Labels,
		Values: logs,
	})

	msg, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(msg); err != nil {
		return fmt.Errorf("failed to gzip json: %w", err)
	}
	if err := gz.Close(); err != nil {
		return fmt.Errorf("failed to close gzip writer: %w", err)
	}

	req, err := http.NewRequest("POST", lp.config.Url, &buf)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	if lp.config.Username != "" && lp.config.Password != "" {
		req.SetBasicAuth(lp.config.Username, lp.config.Password)
	}

	resp, err := lp.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("received unexpected response code from Loki: %s", resp.Status)
	}

	return nil
}
