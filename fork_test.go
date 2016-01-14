package fork

import (
	"math/rand"
	"testing"
	"time"
)

var v = time.Millisecond * 15

func testFork(t *testing.T, r, s, min int) {
	f := NewFork(r)
	o := time.Now()
	for i := 0; i != s; i++ {
		f.Puah(func() {
			time.Sleep(v * 1)
		})
	}
	f.Join()
	su := time.Now().Sub(o)
	in := v * time.Duration(min)
	ax := in + v*2
	t.Log("(", r, ",", s, ")", in, "<", su, "<", ax)
	if su < in || su > ax {
		t.Fatalf("fork error")
	}
}

func testLoop(t *testing.T, j, r, s int) {
	for i := 0; i != j; i++ {
		u := rand.Uint32()%uint32(r) + 1
		m := rand.Uint32()%(u*uint32(s)) + 1
		testFork(t, int(u), int(m), int(m/u))
	}
}

func TestFork(t *testing.T) {
	testLoop(t, 10, 10000, 3)
	testLoop(t, 10, 100, 3)
	testLoop(t, 2, 10, 3)
}
