#include <stdio.h>
#include <stdlib.h>

#include "memory_manager.h"

//extern void initAllocator(uint32_t newMaxSize);

int main()
{
	initAllocator(1024 * 1024);
	
	return 0;
}