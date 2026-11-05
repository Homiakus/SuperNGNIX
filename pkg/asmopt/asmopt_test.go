package asmopt

import "testing"

func TestAVX2ChecksumBytes(t *testing.T) {
	data := []byte("GET /api/v1/resource HTTP/1.1\r\nHost: gateway.internal\r\n\r\n")
	res := AVX2ChecksumBytes(data)
	if res == 0 {
		t.Fatalf("expected non-zero checksum, got 0")
	}
}

func TestFastStringMatchSearch(t *testing.T) {
	text := []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nEdgeGateway Response")
	pattern := []byte("EdgeGateway")
	idx := FastStringMatchSearch(text, pattern)
	if idx < 0 {
		t.Fatalf("expected to find pattern, got %d", idx)
	}

	notExist := FastStringMatchSearch(text, []byte("NonExistentHeader"))
	if notExist != -1 {
		t.Fatalf("expected -1 for non-existent pattern, got %d", notExist)
	}
}

func TestMemoryZeroAligned(t *testing.T) {
	buf := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	MemoryZeroAligned(buf)
	for i, b := range buf {
		if b != 0 {
			t.Fatalf("byte at index %d is not zero: %d", i, b)
		}
	}
}
