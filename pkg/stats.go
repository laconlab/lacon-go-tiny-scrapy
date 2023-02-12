package lacon

import (
	"fmt"
	"strings"
	"sync"
)

type StatusCodeStats struct {
	Stats map[int]int // code:count
}

func (s *StatusCodeStats) Record(code int) {
	s.Stats[code]++
}

func (s *StatusCodeStats) String() string {
	total := 0
	for _, v := range s.Stats {
		total += v
	}

	out := make([]string, 0)
	for code, count := range s.Stats {
		out = append(out, fmt.Sprintf("%d:%.2f", code, float64(count)/float64(total)))
	}
	return strings.Join(out, "\t")
}

type DownloadStats struct {
	ResponseCodesBySite map[string]*StatusCodeStats
	lock                *sync.Mutex
}

func NewStats() *DownloadStats {
	return &DownloadStats{
		ResponseCodesBySite: make(map[string]*StatusCodeStats),
		lock:                &sync.Mutex{},
	}
}

func (s *DownloadStats) String() string {
	s.lock.Lock()
	defer s.lock.Unlock()
	out := make([]string, 0)
	for site, stats := range s.ResponseCodesBySite {
		out = append(out, fmt.Sprintf("%s\t%s", site, stats.String()))
	}
	return strings.Join(out, "\n")
}

func (s *DownloadStats) Record(name string, code int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	stats, ok := s.ResponseCodesBySite[name]
	if !ok {
		stats = &StatusCodeStats{make(map[int]int)}
	}
	stats.Record(code)
	s.ResponseCodesBySite[name] = stats
}
