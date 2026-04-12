package main

import (
	"fmt"
	"unsafe"
)

func IllO11IOI1Ol(IO0II1l0IOI0 []byte) string {
	for I1l1l001IO1I := range IO0II1l0IOI0 {
		IO0II1l0IOI0[I1l1l001IO1I] ^= 0xAA
	}
	return string(IO0II1l0IOI0)
}

type l0I1OllII01O struct {
	Il1O1l0I10lI	[8]byte
	lIlO1O1OOIl1	[8]byte
}

func main() {
	IO0II1l0IOI0 := l0I1OllII01O{
		Il1O1l0I10lI:	[8]byte{'A', 'A', 'A', 'A', 'A', 'A', 'A', 'A'},
		lIlO1O1OOIl1:	[8]byte{'B', 'B', 'B', 'B', 'B', 'B', 'B', 'B'},
	}

	fmt.Println("Before write: ")
	fmt.Printf("buf1: %q\n", string(IO0II1l0IOI0.Il1O1l0I10lI[:]))
	fmt.Printf("buf2: %q\n", string(IO0II1l0IOI0.lIlO1O1OOIl1[:]))
	fmt.Println()

	var llIIOI1II0OO []byte
	fmt.Println("Enter your message:")
	fmt.Scan(&llIIOI1II0OO)
	fmt.Printf("Input (%d bytes): %q\n\n", len(llIIOI1II0OO), string(llIIOI1II0OO))

	fmt.Println("Unsafe write to buf1:")
	II1O1IlO1100(&IO0II1l0IOI0.Il1O1l0I10lI, llIIOI1II0OO)
	fmt.Printf("buf1: %q\n", string(IO0II1l0IOI0.Il1O1l0I10lI[:]))
	fmt.Printf("buf2: %q\n", string(IO0II1l0IOI0.lIlO1O1OOIl1[:]))
	fmt.Println()

	IO0II1l0IOI0 = l0I1OllII01O{
		Il1O1l0I10lI:	[8]byte{'A', 'A', 'A', 'A', 'A', 'A', 'A', 'A'},
		lIlO1O1OOIl1:	[8]byte{'B', 'B', 'B', 'B', 'B', 'B', 'B', 'B'},
	}

	fmt.Println("Safe write to buf1:")
	IlOO11I10ll1 := IO1II101l101(&IO0II1l0IOI0.Il1O1l0I10lI, llIIOI1II0OO)
	if IlOO11I10ll1 != nil {
		fmt.Printf("error: %v\n", IlOO11I10ll1)
	}
	fmt.Printf("buf1: %q\n", string(IO0II1l0IOI0.Il1O1l0I10lI[:]))
	fmt.Printf("buf2: %q\n", string(IO0II1l0IOI0.lIlO1O1OOIl1[:]))
}

func II1O1IlO1100(l0lIO01llOI0 *[8]byte, lIOOll10O0Il []byte) {
	for I1l1l001IO1I := range lIOOll10O0Il {
		Il0I0001OOIl := (*byte)(unsafe.Add(unsafe.Pointer(l0lIO01llOI0), I1l1l001IO1I))
		*Il0I0001OOIl = lIOOll10O0Il[I1l1l001IO1I]
	}
}

func IO1II101l101(l0lIO01llOI0 *[8]byte, lIOOll10O0Il []byte) error {
	var IlOO11I10ll1 error
	ll1O1I01Ol0l := len(lIOOll10O0Il)
	if ll1O1I01Ol0l > len(l0lIO01llOI0) {
		ll1O1I01Ol0l = len(l0lIO01llOI0)
		IlOO11I10ll1 = fmt.Errorf("truncated: wrote %d of %d bytes", ll1O1I01Ol0l, len(lIOOll10O0Il))
	}
	for I1l1l001IO1I := 0; I1l1l001IO1I < ll1O1I01Ol0l; I1l1l001IO1I++ {
		l0lIO01llOI0[I1l1l001IO1I] = lIOOll10O0Il[I1l1l001IO1I]
	}
	return IlOO11I10ll1
}
