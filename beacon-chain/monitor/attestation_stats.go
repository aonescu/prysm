package monitor

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

type AttestationStats struct {
	TotalVerified  int
	TotalFailed    int
	FailureReasons map[string]int
	mu             sync.Mutex
}

func NewAttestationStats() *AttestationStats {
	return &AttestationStats{
		FailureReasons: make(map[string]int),
	}
}

func (s *AttestationStats) LogSuccess() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.TotalVerified++
	log.WithField("verified", s.TotalVerified).Info("Verified attestation")
}

func (s *AttestationStats) LogFailure(reason string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.TotalFailed++
	s.FailureReasons[reason]++
	log.WithFields(logrus.Fields{
		"reason":       reason,
		"total_failed": s.TotalFailed,
	}).Warn("Failed attestation recorded")
}

func (s *AttestationStats) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.TotalVerified = 0
	s.TotalFailed = 0
	s.FailureReasons = make(map[string]int)
	log.Info("Attestation stats reset for new epoch")
}

func (s *AttestationStats) Summary() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	summary := fmt.Sprintf("Epoch Summary: Total Verified: %d, Total Failed: %d, Failure Reasons: %v",
		s.TotalVerified, s.TotalFailed, s.FailureReasons)
	log.Info(summary)
	return summary
}
