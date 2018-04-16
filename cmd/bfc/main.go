package main

import (
	"../../bf"
	"flag"
	"fmt"
	"io/ioutil"
	"llvm.org/llvm/bindings/go/llvm"
	"os"
	"path/filepath"
)

var outputFile = flag.String("o", "", "Output file")

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Printf("Usage: %s [options..] <input>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(-1)
	}

	mod := llvm.NewModule("bf")

	putcharFuncType := llvm.FunctionType(llvm.VoidType(), []llvm.Type{llvm.Int8Type()}, false)
	llvm.AddFunction(mod, "bf_putchar", putcharFuncType)
	getcharFuncType := llvm.FunctionType(llvm.Int8Type(), []llvm.Type{}, false)
	llvm.AddFunction(mod, "bf_getchar", getcharFuncType)
	memsetFuncType := llvm.FunctionType(llvm.PointerType(llvm.Int64Type(), 0), []llvm.Type{llvm.PointerType(llvm.Int8Type(), 0), llvm.Int8Type(), llvm.Int64Type()}, false)
	llvm.AddFunction(mod, "memset", memsetFuncType)

	if res, err := ioutil.ReadFile(flag.Arg(0)); err != nil {
		fmt.Printf("Could not read input file: %s", err)
		os.Exit(-1)
	} else if err := bf.Compile(string(res), mod); err != nil {
		fmt.Printf("Could not compile: %s", err)
		os.Exit(-1)
	}

	f, err := os.Create(*outputFile)
	if err != nil {
		fmt.Printf("Cant open output file: %s", err)
		os.Exit(-1)
	}
	defer f.Close()
	if err := llvm.WriteBitcodeToFile(mod, f); err != nil {
		fmt.Printf("Failed to write bitcode to file: %s", err)
		os.Exit(-1)
	}
}
