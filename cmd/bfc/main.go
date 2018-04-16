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
var optLevel = flag.Uint("O", 0, "Optimization level")

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Printf("Usage: %s [options..] <input>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(-1)
	}

	inputFile := flag.Arg(0)

	mod := llvm.NewModule(inputFile)

	putcharFuncType := llvm.FunctionType(llvm.VoidType(), []llvm.Type{llvm.Int8Type()}, false)
	llvm.AddFunction(mod, "bf_putchar", putcharFuncType)
	getcharFuncType := llvm.FunctionType(llvm.Int8Type(), []llvm.Type{}, false)
	llvm.AddFunction(mod, "bf_getchar", getcharFuncType)
	memsetFuncType := llvm.FunctionType(llvm.PointerType(llvm.Int64Type(), 0), []llvm.Type{llvm.PointerType(llvm.Int8Type(), 0), llvm.Int8Type(), llvm.Int64Type()}, false)
	llvm.AddFunction(mod, "memset", memsetFuncType)

	if res, err := ioutil.ReadFile(inputFile); err != nil {
		fmt.Printf("Could not read input file: %s", err)
		os.Exit(-1)
	} else if err := bf.Compile(string(res), mod); err != nil {
		fmt.Printf("Could not compile: %s", err)
		os.Exit(-1)
	}

	if *outputFile == "" {
		ext := filepath.Ext(inputFile)
		*outputFile = inputFile[:len(inputFile)-len(ext)] + ".bc"
	}
	f, err := os.Create(*outputFile)
	if err != nil {
		fmt.Printf("Cant open output file: %s", err)
		os.Exit(-1)
	}
	defer f.Close()
	pm := llvm.NewPassManager()
	defer pm.Dispose()
	pmb := llvm.NewPassManagerBuilder()
	defer pmb.Dispose()
	pmb.SetOptLevel(int(*optLevel))
	pmb.Populate(pm)
	pm.Run(mod)
	if err := llvm.WriteBitcodeToFile(mod, f); err != nil {
		fmt.Printf("Failed to write bitcode to file: %s", err)
		os.Exit(-1)
	}
}
