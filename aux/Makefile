# Sample Makefile for buliding bf binaries
BFC=../bfc
LLC=llc

BFC_FLAGS=-O 3
LLC_FLAGS=-O3
CFLAGS=-O3

TARGETS=bench hello fib squares

.PHONY: all
all: ${TARGETS}

%.bc: %.b
	$(BFC) ${BFC_FLAGS} -o $@ $<

%.s: %.bc
	$(LLC) ${LLC_FLAGS} -o $@ $<

rt.o: rt.c
	$(CC) ${CFLAGS} -c -o $@ $<

%.o: %.s
	$(AS) -o $@ $<

${TARGETS}: %: %.o rt.o
	$(CC) ${LDFLAGS} $< rt.o -o $@

.PHONY: clean
clean:
	$(RM) *.o *.bc *.s ${TARGETS}
