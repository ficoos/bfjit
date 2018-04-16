package main

import (
	"../../bf"
	"flag"
	"fmt"
	"io/ioutil"
	"llvm.org/llvm/bindings/go/llvm"
	"os"
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Printf("Usage: %s <input>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(-1)
	}

	llvm.LinkInMCJIT()
	if err := llvm.InitializeNativeTarget(); err != nil {
		panic(err)
	}

	if err := llvm.InitializeNativeAsmPrinter(); err != nil {
		panic(err)
	}

	mod := llvm.NewModule("bf")
	opts := llvm.MCJITCompilerOptions{}
	opts.SetMCJITCodeModel(llvm.CodeModelJITDefault)
	opts.SetMCJITOptimizationLevel(3)

	pm := llvm.NewPassManager()
	defer pm.Dispose()
	pmb := llvm.NewPassManagerBuilder()
	defer pmb.Dispose()
	pmb.SetOptLevel(3)
	pmb.Populate(pm)

	engine, err := llvm.NewMCJITCompiler(mod, opts)
	if err != nil {
		panic(err)
	}

	putcharFuncType := llvm.FunctionType(llvm.VoidType(), []llvm.Type{llvm.Int8Type()}, false)
	putchar := llvm.AddFunction(mod, "bf_putchar", putcharFuncType)
	putchar.SetLinkage(llvm.ExternalLinkage)
	getcharFuncType := llvm.FunctionType(llvm.Int8Type(), []llvm.Type{}, false)
	getchar := llvm.AddFunction(mod, "bf_getchar", getcharFuncType)
	getchar.SetLinkage(llvm.ExternalLinkage)
	memsetFuncType := llvm.FunctionType(llvm.PointerType(llvm.Int64Type(), 0), []llvm.Type{llvm.PointerType(llvm.Int8Type(), 0), llvm.Int8Type(), llvm.Int64Type()}, false)
	llvm.AddFunction(mod, "memset", memsetFuncType)

	if res, err := ioutil.ReadFile(flag.Arg(0)); err != nil {
		fmt.Printf("Could not read input file: %s", err)
		os.Exit(-1)
	} else if err := bf.Compile(string(res), mod); err != nil {
		fmt.Printf("Could not compile: %s", err)
		os.Exit(-1)
	}

	pm.Run(mod)

	engine.RunStaticConstructors()
	engine.RunFunction(mod.NamedFunction("main"), []llvm.GenericValue{})
	engine.RunStaticDestructors()
}
