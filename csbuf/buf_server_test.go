package csbuf

import (
	"testing"
	"time"
)

func TestServerBuffer(t *testing.T) {
	bp := NewBuffer()
	go func() {
		duration := 1 * time.Second
		time.Sleep(duration)
		bp.Insert(100)
		time.Sleep(duration)
		bp.Insert(200)
	}()
	if !bp.Empty() {
		t.Log("Expect empty")
		t.Fail()
	}
	x := bp.Front()
	if x != nil {
		t.Logf("Expect nil, get %d", x)
		t.Fail()
	}
	time.Sleep(1 * time.Second)
	x = bp.Front()
	if x != 100 {
		t.Logf("Expect 100, get %d", x)
		t.Fail()
	}
	x = bp.Remove()

	time.Sleep(1 * time.Second)
	x = bp.Front()
	if x != 200 {
		t.Logf("Expect 200, get %d", x)
		t.Fail()
	}
	bp.Flush()
	if !bp.Empty() {
		t.Log("Expect empty")
		t.Fail()
	}
}
