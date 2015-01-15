package bit

func CountOnes64(n uint64) int {
	// Bit population count, see
	// http://graphics.stanford.edu/~seander/bithacks.html#CountBitsSetParallel
	// This is faster than table lookup:
	//	BenchmarkPopcountBitflip500000000       4.8 ns/op
	//	BenchmarkPopcountTable100000000        13.0 ns/op
	n -= (n >> 1) & 0x5555555555555555
	n = (n>>2)&0x3333333333333333 + n&0x3333333333333333
	n += n >> 4
	n &= 0x0f0f0f0f0f0f0f0f
	n *= 0x0101010101010101
	return int(n >> 56)
}
