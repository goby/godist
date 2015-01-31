package chanbuf

import (
	"testing"
	"time"
)

func TestChanBuffer(t *testing.T) {
	bp := NewChanBuffer()
	go func() {
		duration := 1 * time.Second
		time.Sleep(duration)
		bp.Insert(100)
		bp.Remove()
		time.Sleep(duration)
		bp.Insert(200)
	}()
	if !bp.Empty() {
		t.Log("Expect empty")
		t.Fail()
	}
	x := bp.Front()
	if x != 100 {
		t.Log("fail")
		t.Fail()
	}
	//time.Sleep(100 * time.Microsecond)
	x = bp.Remove()
	if x != 200 {
		t.Logf("Expect 200, get %d", x)
		t.Fail()
	}
}
