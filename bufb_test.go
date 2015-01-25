package bufb

import (
	"encoding/json"
	//"github.com/goby/godist/bufb"
	"math/rand"
	"testing"
)

func i2b(i int) []byte {
	ret, _ := json.Marshal(i)
	return ret
}

func b2i(b []byte) int {
	var i int
	json.Unmarshal(b, &i)
	return i
}

func TestBuffer(t *testing.T) {

}

func runtest(t *testing.T, bp *Buffer) {
	inserted := 0
	removed := 0
	emptyCount := 0
	for removed < 10 {
		if bp.Empty() {
			emptyCount++
		}
		insert := !(inserted == 10)
		if inserted > removed && rand.Int31n(2) == 0 {
			insert = false
		}
		if insert {
			bp.Insert(i2b(inserted))
			inserted++
		} else {
			b, err := bp.Remove()
			if err != nil {
				t.Logf("Attempt to remove from empty buffer.\n")
				t.Fail()
			}
			v := b2i(b)
			if v != removed {
				t.Logf("Removed %d, Expected %d\n", v, removed)
				t.Fail()
			}
			removed++
		}
	}
}
