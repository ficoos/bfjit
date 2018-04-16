/* Sample runtime for bfc objects */

#include <stdint.h>
#include <unistd.h>

void bf_putchar(uint8_t c) {
	write(1, &c, 1);
}

uint8_t bf_getchar() {
	uint8_t c = 0;
	read(0, &c, 1);
	return c;
}
