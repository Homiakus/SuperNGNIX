package asmopt

// AVX2ChecksumBytes computes an optimized 64-bit checksum over byte stream using AVX2/x86_64 assembly.
func AVX2ChecksumBytes(buf []byte) uint64

// FastStringMatchSearch searches for pattern in text using vectorized comparison.
func FastStringMatchSearch(text, pattern []byte) int

// MemoryZeroAligned zeroes out count bytes aligned to 32 bytes using AVX2.
func MemoryZeroAligned(buf []byte)
