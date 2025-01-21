package services

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
	"time"
)

const (
	ArticleStatsServiceBufferSize    = 100
	ArticleStatsServiceWriteInterval = 2 // In Seconds
)

type ArticleStats struct {
	Views int `json:"views"`
}

type ArticleStatsService struct {
	log       *log.Logger
	statsFile string
	stats     map[string]*ArticleStats
	lock      sync.Mutex
	events    chan string
}

func NewArticleStatsService(
	ctx context.Context,
	log *log.Logger,
	statsFile string,
) (*ArticleStatsService, error) {
	service := &ArticleStatsService{
		log:       log,
		stats:     make(map[string]*ArticleStats),
		statsFile: statsFile,
		lock:      sync.Mutex{},
		events:    make(chan string, ArticleStatsServiceBufferSize),
	}
	if err := service.loadCurrentStats(); err != nil {
		return nil, err
	}
	go service.run(ctx)
	return service, nil
}

func (ss *ArticleStatsService) run(ctx context.Context) {
	ticker := time.NewTicker(ArticleStatsServiceWriteInterval * time.Second)
	defer ticker.Stop()

	for {
		select {
		case article := <-ss.events:
			ss.lock.Lock()
			if _, exists := ss.stats[article]; !exists {
				ss.stats[article] = &ArticleStats{}
			}
			ss.stats[article].Views += 1
			ss.lock.Unlock()

		case <-ticker.C:
			if err := ss.writeToDisk(); err != nil {
				ss.log.Printf("ERROR: Could not write stats to disk: %s\n", err)
				continue
			}

		case <-ctx.Done():
			ss.log.Printf("Shutting down ArticleStatsService")
			if err := ss.writeToDisk(); err != nil {
				ss.log.Printf("ERROR: Could not write stats to disk: %s\n", err)
			}
			return
		}
	}
}

func (ss *ArticleStatsService) IncrementView(article string) {
	select {
	case ss.events <- article:
	default:
		ss.log.Printf("Failed to queue view increment for article '%s' because buffer is full.\n", article)
	}
}

func (ss *ArticleStatsService) loadCurrentStats() error {
	ss.lock.Lock()
	defer ss.lock.Unlock()

	if _, err := os.Stat(ss.statsFile); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(ss.statsFile)
		if err != nil {
			return err
		}
		file.Close()
		return nil
	}

	data, err := os.ReadFile(ss.statsFile)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	stats := make(map[string]*ArticleStats)
	if err = json.Unmarshal(data, &stats); err != nil {
		return err
	}

	ss.stats = stats
	ss.log.Printf("Loaded Statistics for %d Articles\n", len(stats))

	return nil

}

func (ss *ArticleStatsService) writeToDisk() error {
	ss.lock.Lock()
	defer ss.lock.Unlock()

	file, err := os.Create(ss.statsFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(ss.stats); err != nil {
		return err
	}

	return nil
}
