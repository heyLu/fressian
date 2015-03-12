#include <errno.h>
#include <stdio.h>
#include <string.h>
#include <zlib.h>

int main(int argc, char **argv) {
	if (argc != 2) {
		fprintf(stderr, "Usage: %s file\n", argv[0]);
		return 1;
	}

	FILE *fp = NULL;
	if (strcmp(argv[1], "-") == 0) {
		fp = stdin;
	} else {
		fp = fopen(argv[1], "r");
		if (fp == NULL) {
			fprintf(stderr, "Error: %s\n", strerror(errno));
			return 1;
		}
	}

	// initialize adler
	uLong adler = adler32(0L, Z_NULL, 0);

	Bytef buf[100];
	size_t n;
	int total = 0;
	while ((n = fread(buf, 1, 100, fp)) == 100) {
		adler = adler32(adler, buf, n);
		total += n;
	}
	adler = adler32(adler, buf, n);
	total += n;

	printf("%d bytes read\n", total);
	printf("adler32: %lu\n", adler);

	fclose(fp);

	return 0;
}
