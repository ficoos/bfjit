package bf

import (
	"errors"
	"fmt"
	"llvm.org/llvm/bindings/go/llvm"
	"unicode/utf8"
)

const MAX_LOOP_DEPTH = 100

type BFCompilerContext struct {
	mod     llvm.Module
	builder llvm.Builder
	prog    llvm.Value
}

var mainFunctionType llvm.Type

func init() {
	mainFunctionType = llvm.FunctionType(llvm.VoidType(), []llvm.Type{}, false)
}

func Compile(code string, mod llvm.Module) error {
	itemType := llvm.Int8Type()
	one := llvm.ConstInt(itemType, 1, false)
	minusOne := llvm.ConstNeg(one)
	zero := llvm.ConstInt(itemType, 0, false)
	putchar := mod.NamedFunction("bf_putchar")
	getchar := mod.NamedFunction("bf_getchar")
	memset := mod.NamedFunction("memset")
	memSize := llvm.ConstInt(llvm.Int64Type(), 4096, false)

	prog := llvm.AddFunction(mod, "main", mainFunctionType)
	block := llvm.AddBasicBlock(prog, "entry")
	builder := llvm.NewBuilder()
	defer builder.Dispose()
	builder.SetInsertPoint(block, block.FirstInstruction())
	tape := builder.CreateArrayMalloc(itemType, memSize, "tape")
	builder.CreateCall(memset, []llvm.Value{tape, zero, memSize}, "")
	ptrType := llvm.PointerType(itemType, tape.Type().PointerAddressSpace())
	ptr := builder.CreateAlloca(ptrType, "ptr")
	builder.CreateStore(tape, ptr)

	loopCounter := 0

	loopDepth := 0
	loopBlocks := make([]llvm.BasicBlock, MAX_LOOP_DEPTH)

	for i := 0; i < len(code); {
		c, l := utf8.DecodeRune([]byte(code[i:]))
		i += l
		switch c {
		case '>':
			p := builder.CreateLoad(ptr, "")
			np := builder.CreateGEP(p, []llvm.Value{one}, "")
			builder.CreateStore(np, ptr)
		case '<':
			p := builder.CreateLoad(ptr, "")
			np := builder.CreateGEP(p, []llvm.Value{minusOne}, "")
			builder.CreateStore(np, ptr)
		case '+':
			p := builder.CreateLoad(ptr, "")
			v := builder.CreateLoad(p, "tmp")
			v1 := builder.CreateAdd(v, one, "tmp")
			builder.CreateStore(v1, p)
		case '-':
			p := builder.CreateLoad(ptr, "")
			v := builder.CreateLoad(p, "tmp")
			v1 := builder.CreateSub(v, one, "tmp")
			builder.CreateStore(v1, p)
		case '.':
			p := builder.CreateLoad(ptr, "")
			v := builder.CreateLoad(p, "tmp")
			builder.CreateCall(putchar, []llvm.Value{v}, "")
		case ',':
			v := builder.CreateCall(getchar, []llvm.Value{}, "tmp")
			p := builder.CreateLoad(ptr, "")
			builder.CreateStore(v, p)
		case '[':
			innerBlock := llvm.AddBasicBlock(prog, fmt.Sprintf("loop_%d_enter", loopCounter))
			outerBlock := llvm.AddBasicBlock(prog, fmt.Sprintf("loop_%d_exit", loopCounter))
			innerBlock.MoveAfter(builder.GetInsertBlock())
			outerBlock.MoveAfter(innerBlock)
			if loopDepth == MAX_LOOP_DEPTH {
				return errors.New("To many nesting []")
			}
			loopBlocks[loopDepth] = innerBlock
			loopDepth++
			v := builder.CreateLoad(builder.CreateLoad(ptr, ""), "tmp")
			builder.CreateCondBr(builder.CreateICmp(llvm.IntEQ, v, zero, ""), outerBlock, innerBlock)
			builder.SetInsertPoint(innerBlock, innerBlock.FirstInstruction())
		case ']':
			loopDepth--
			if loopDepth < 0 {
				return errors.New("] without a matching [")
			}
			nextBlock := llvm.NextBasicBlock(builder.GetInsertBlock())
			v := builder.CreateLoad(builder.CreateLoad(ptr, ""), "tmp")
			builder.CreateCondBr(builder.CreateICmp(llvm.IntNE, v, zero, ""), loopBlocks[loopDepth], nextBlock)
			builder.SetInsertPoint(nextBlock, nextBlock.FirstInstruction())
		}
	}

	if loopDepth != 0 {
		return errors.New("[ without matching ]")
	}

	builder.CreateRetVoid()

	if err := llvm.VerifyFunction(prog, llvm.PrintMessageAction); err != nil {
		return err
	}

	return nil
}
