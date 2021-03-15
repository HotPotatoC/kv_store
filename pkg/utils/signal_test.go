package utils_test

import (
	"syscall"
	"testing"
	"time"

	"github.com/HotPotatoC/kvstore/pkg/utils"
)

func TestWaitForSignals(t *testing.T) {
	go func() {
		time.Sleep(250 * time.Millisecond)
		syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	}()

	signal := <-utils.WaitForSignals(syscall.SIGTERM)
	if signal != syscall.SIGTERM {
		t.Errorf("Failed TestWaitForSignals -> Expected: %s | Got: %s", "SIGTERM", signal.String())
	}
}