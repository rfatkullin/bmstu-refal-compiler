#include <stdlib.h>
#include "allocator.h"

void initAllocator(uint32_t newMaxSize)
{
	heap.maxSize = newMaxSize / 2;
	heap.head = malloc(heap.maxSize);
	heap.swapHead = malloc(heap.maxSize);
	heap.offset = 0;
}

struct v_term* allocate(struct environment* env, struct l_term* expr)
{

}
