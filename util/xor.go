// Originally from the crypto/ciper/xor.go package and modified to use pointers

package util

import (
	"github.com/goarchit/archit/log"
	"runtime"
	"unsafe"
)

const wordSize = int(unsafe.Sizeof(uintptr(0)))
const supportsUnaligned = runtime.GOARCH == "386" || runtime.GOARCH == "amd64" || runtime.GOARCH == "ppc64" || runtime.GOARCH == "ppc64le" || runtime.GOARCH == "s390x"

// fastXORBytes xors in bulk. It only works on architectures that
// support unaligned read/writes.
func fastXORBytes(dst, b []byte ) int {
	n := len(b)

	w := n / wordSize
	if w > 0 {
		dw := *(*[]uintptr)(unsafe.Pointer(&dst))
		bw := *(*[]uintptr)(unsafe.Pointer(&b))
		for i := 0; i < w; i++ {
			dw[i] ^= bw[i]
		}
	}

	for i := (n - n%wordSize); i < n; i++ {
		dst[i] ^= b[i]
	}

	return n
}

func safeXORBytes(dst, b []byte) int {
	n := len(b)

	for i := 0; i < n; i++ {
		dst[i] ^= b[i]
	}
	return n
}

// xorBytes xors the bytes in a and b. The destination is assumed to have enough
// space. Returns the number of bytes xor'd.
var warnonce bool
func XorBytes(dst, b []byte) int {
	if supportsUnaligned {
		return fastXORBytes(dst, b)
	} else {
		if !warnonce {
			log.Warning("Architecture does not support unaligned bytes - program will run slower")
			warnonce = true
		}
		return safeXORBytes(dst, b)
	}
}
