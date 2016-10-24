package fork

import (
	"math/rand"
	"testing"
	"time"
)

var v = time.Millisecond * 10

func testFork(t *testing.T, r, s, min int) {
	f := NewFork(r)
	o := time.Now()
	for i := 0; i != s; i++ {
		f.Push(func() {
			time.Sleep(v + time.Duration(rand.Int63n(int64(v/5))))
		})
	}
	f.Join()
	su := time.Now().Sub(o)
	in := v * time.Duration(min)
	ax := in + v*4
	t.Log("(", r, ",", s, ")", in, "<", su, "<", ax)
	if su < in || su > ax {
		t.Error("fork error")
		//t.Fatalf("fork error")
	}
}

func testLoop(t *testing.T, j, r, s int) {
	for i := 0; i != j; i++ {
		u := rand.Uint32()%uint32(r) + 1
		m := rand.Uint32()%(u*uint32(s)) + 1
		min := int(m / u)
		if m%u != 0 {
			min++
		}
		testFork(t, int(u), int(m), min)
	}
}

func TestFork(t *testing.T) {
	for i := 0; i != 10; i++ {
		testFork(t, 9, 10, 2)
	}
	testLoop(t, 10, 1000, 3)
	testLoop(t, 10, 100, 3)
	testLoop(t, 10, 10, 3)
}

func TestCCC(t *testing.T) {
	f := NewFork(1)
	f.Push(func() {
		time.Sleep(1 * time.Second)
	})
	f.Join()
}
