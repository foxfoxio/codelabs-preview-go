/*
 * Copyright (c) adixity 2020. https://github.com/adixity
 */

package stopwatch

import (
	"github.com/foxfoxio/codelabs-preview-go/internal/logger"
	"time"
)

type StopWatch struct {
	start time.Time
	log   *logger.Logger
}

type Stopper interface {
	Stop() time.Duration
}

func (s *StopWatch) Stop() (diff time.Duration) {
	defer func() {
		if s.log != nil {
			s.log.Printf("elapsed time: %v", diff)
		}
	}()

	return time.Now().Sub(s.start)
}

func StartWithDefaultLogger() Stopper {
	return StartWithLogger(logger.New("StopWatch"))
}

func Start() Stopper {
	return StartWithLogger(nil)
}

func StartWithLogger(log *logger.Logger) Stopper {
	if log != nil {
		log.Println("start")
	}
	return &StopWatch{start: time.Now(), log: log}
}
