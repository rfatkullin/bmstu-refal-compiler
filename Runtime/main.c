#include <stdio.h>
#include <stdlib.h>

#include "allocator.h"

extern void initAllocator(uint32_t newMaxSize);

int main()
{
	initAllocator(200);
	
	return 0;
}