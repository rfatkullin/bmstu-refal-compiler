#include <stdio.h>
#include <stdlib.h>

#include "memory_manager.h"

int main()
{
	initAllocator(512, 2);

	markInSegmentTree(0, 2);
	markInSegmentTree(2, 3);
	markInSegmentTree(4, 4);

	sumInSegmentTree(0, 4);

	printf("sum: %u\n", sumInSegmentTree(0, 4));

	return 0;
}
