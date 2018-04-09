package main

import (
	"C"
	"os"
)

//export bf_putchar
func bf_putchar(c C.uchar) {
	if _, err := os.Stdout.Write([]byte{byte(c)}); err != nil {
		panic(err)
	}
}

//export bf_getchar
func bf_getchar() C.uchar {
	buff := make([]byte, 1)
	if _, err := os.Stdout.Read(buff); err != nil {
		panic(err)
	}

	return C.uchar(buff[0])
}
