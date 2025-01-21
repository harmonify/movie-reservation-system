package logger

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type lokiPusherImpl struct {
	quit      chan struct{}
	entry     chan logEntry
	waitGroup sync.WaitGroup

	config *LokiZapConfig
	client *http.Client
}

func newLokiPusher(cfg *LokiZapConfig) *lokiPusherImpl {
	c := &http.Client{}
	cfg.Url = strings.TrimSuffix(cfg.Url, "/")
	cfg.Url = fmt.Sprintf("%s/loki/api/v1/push", cfg.Url)

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	lp := &lokiPusherImpl{
		config: cfg,
		client: c,
		quit:   make(chan struct{}),
		entry:  make(chan logEntry, cfg.BatchMaxSize*2),
	}

	lp.waitGroup.Add(1)
	go lp.run()
	return lp
}

// stop stops the loki pusher
func (lp *lokiPusherImpl) stop() {
	close(lp.quit)
	lp.waitGroup.Wait()
}

// Run
func (lp *lokiPusherImpl) run() {
	ticker := time.NewTicker(lp.config.BatchMaxWait)
	defer ticker.Stop()

	localBatch := make([]streamValue, 0, lp.config.BatchMaxSize)

	// Cleanup
	defer lp.cleanup(localBatch)

	for {
		select {
		case <-lp.quit:
			return
		case entry := <-lp.entry:
			localBatch = append(localBatch, newLog(entry))
			if len(localBatch) >= lp.config.BatchMaxSize {
				err := lp.send(localBatch)
				if err != nil {
					slog.Error("Failed to send batch", slog.Any("error", err))
				}
				localBatch = make([]streamValue, 0, lp.config.BatchMaxSize) // Reset batch
				ticker.Reset(lp.config.BatchMaxWait)
			}
		case <-ticker.C:
			if len(localBatch) > 0 {
				err := lp.send(localBatch)
				if err != nil {
					slog.Error("Failed to send batch", slog.Any("error", err))
				}
				localBatch = make([]streamValue, 0, lp.config.BatchMaxSize) // Reset batch
				ticker.Reset(lp.config.BatchMaxWait)
			}
		}
	}
}

func (lp *lokiPusherImpl) cleanup(remainingBatch []streamValue) {
	err := lp.send(remainingBatch)
	if err != nil {
		slog.Error("Failed to send batch", slog.Any("error", err))
	}
	lp.waitGroup.Done()
}

func (lp *lokiPusherImpl) send(batch []streamValue) error {
	buf := bytes.NewBuffer(nil)
	gz := gzip.NewWriter(buf)

	err := json.NewEncoder(gz).Encode(lokiPushRequest{
		Streams: []stream{{
			Stream: lp.config.Labels,
			Values: batch,
		}},
	})
	if err != nil {
		return fmt.Errorf("failed to encode logs: %w", err)
	}

	if err := gz.Close(); err != nil {
		return fmt.Errorf("failed to close gzip writer: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, lp.config.Url, buf)
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
		return fmt.Errorf("unexpected response code from Loki: %s", resp.Status)
	}

	return nil
}
