package bit

// Run
//	 /usr/local/go/pkg/tool/darwin_amd64/6g -S div_vs_bitshift.go
// to verify that these produce identical code.

func div(u uint8) (uint8, uint8) {
	return u / 64, u % 64
}

func bitsx(u uint8) (uint8, uint8) {
	return u >> 6, u & 0x3f
}
