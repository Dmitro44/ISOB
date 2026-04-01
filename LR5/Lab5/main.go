package main

import (
	"fmt"
	"unsafe"
)

type Buffers struct {
	buf1 [8]byte
	buf2 [8]byte
}

func main() {
	b := Buffers{
		buf1: [8]byte{'A', 'A', 'A', 'A', 'A', 'A', 'A', 'A'},
		buf2: [8]byte{'B', 'B', 'B', 'B', 'B', 'B', 'B', 'B'},
	}

	fmt.Println("Before write:")
	fmt.Printf("buf1: %q\n", string(b.buf1[:]))
	fmt.Printf("buf2: %q\n", string(b.buf2[:]))
	fmt.Println()

	var input []byte
	fmt.Println("Enter your message:")
	fmt.Scan(&input)
	fmt.Printf("Input (%d bytes): %q\n\n", len(input), string(input))

	fmt.Println("Unsafe write to buf1:")
	unsafeWrite(&b.buf1, input)
	fmt.Printf("buf1: %q\n", string(b.buf1[:]))
	fmt.Printf("buf2: %q\n", string(b.buf2[:]))
	fmt.Println()

	b = Buffers{
		buf1: [8]byte{'A', 'A', 'A', 'A', 'A', 'A', 'A', 'A'},
		buf2: [8]byte{'B', 'B', 'B', 'B', 'B', 'B', 'B', 'B'},
	}

	fmt.Println("Safe write to buf1:")
	err := safeWrite(&b.buf1, input)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("buf1: %q\n", string(b.buf1[:]))
	fmt.Printf("buf2: %q\n", string(b.buf2[:]))
}

func unsafeWrite(buf *[8]byte, data []byte) {
	for i := range data {
		ptr := (*byte)(unsafe.Add(unsafe.Pointer(buf), i))
		*ptr = data[i]
	}
}

func safeWrite(buf *[8]byte, data []byte) error {
	if len(data) > len(buf) {
		return fmt.Errorf("overflow: data (%d bytes) > buffer (%d bytes)", len(data), len(buf))
	}
	copy(buf[:], data)
	return nil
}
