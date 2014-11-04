#include <stdio.h>
#include <stdlib.h>
#include <assert.h>
#include "allocator.h"

static uint32_t allocateMemoryForVTerms(uint32_t size)
{
	uint32_t n = (N * size) / (2 * N * sizeof(struct v_term) + 4 * sizeof(uint32_t));
	uint32_t activeTermsHeapSize = n * sizeof(struct v_term);
	uint32_t inactiveTermsHeapSize = n * sizeof(struct v_term);
	uint32_t segmentTreeHeapSize = n * sizeof(uint32_t);
	uint32_t usedMemory = activeTermsHeapSize + inactiveTermsHeapSize + segmentTreeHeapSize;
	
	assert(usedMemory <= size);
	
	memoryManager.activeTermsHeap = (uint8_t*)malloc(activeTermsHeapSize);
	memoryManager.inactiveTermsHeap = (uint8_t*)malloc(inactiveTermsHeapSize);	
	memoryManager.segmentTree = (uint32_t*)malloc(segmentTreeHeapSize);
	
	printf("Memory allocation for vterms:\n");
	printf("\tMemory enough for %d terms\n", n);
	printf("\tActive vterms heap size: %d\n", activeTermsHeapSize);
	printf("\tInactive vterms heap size: %d\n", inactiveTermsHeapSize);
	printf("\tSegment tree size: %d\n", segmentTreeHeapSize);
	printf("\tUsed memory: %d\n", usedMemory);
	printf("\tLost memory: %d\n", size - usedMemory);		
	
	return usedMemory;
}

static uint32_t allocateMemoryForData(uint32_t size)
{
	uint32_t singleDataHeapSize = size / 2;
	uint32_t usedMemory = 2 * singleDataHeapSize;
	
	memoryManager.activeDataHeap = (uint8_t*)malloc(singleDataHeapSize);
	memoryManager.inactiveTermsHeap = (uint8_t*)malloc(singleDataHeapSize);
	
	printf("Memory allocation for data:\n");
	printf("\tMemory for single data heap: %d\n", singleDataHeapSize);
	printf("\tUsed memory: %d\n", usedMemory);
	printf("\tLost memory: %d\n", size - usedMemory);			
	
	return usedMemory;
}

static uint32_t allocateMemoryForLTerms(uint32_t size)
{
	memoryManager.lTermsHeap = (uint8_t*)malloc(size);
	
	printf("Memory allocated for lterms: %d\n", size);
	
	return size;
}

void initAllocator(uint32_t size)
{	
	uint32_t dataHeapSize = DATA_HEAP_SIZE_FACTOR * size;
	uint32_t vtermsHeapSize = V_TERMS_HEAP_SIZE_FACTOR * size;
	uint32_t ltermsHeapSize = size - dataHeapSize - vtermsHeapSize;
	uint32_t usedMemory = 0;
	
	usedMemory += allocateMemoryForData(dataHeapSize);
	usedMemory += allocateMemoryForVTerms(vtermsHeapSize);	
	usedMemory += allocateMemoryForLTerms(ltermsHeapSize);
	
	assert(usedMemory < size);
	
	memoryManager.size = usedMemory;
	
	memoryManager.vtermsOffset = 0;
	memoryManager.dataOffset = 0;
	memoryManager.ltermsOffset = 0;	
	
	printf("Total used memory: %d\n", usedMemory);
}

struct v_term* allocate(struct l_term* expr)
{

}
