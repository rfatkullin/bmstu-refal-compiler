#include <stdio.h>
#include <stdlib.h>
#include <assert.h>
#include <time.h>
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
	
	memoryManager.totalSize = usedMemory;
	memoryManager.vtermsCount = 0;
	
	memoryManager.vtermsOffset = 0;
	memoryManager.dataOffset = 0;
	memoryManager.ltermsOffset = 0;	
	
	printf("Total used memory: %d\n", usedMemory);
}

void collectGarbage(struct l_term* expr)
{
	clock_t start, end;	
	printf("Start garbage collection.\n");
	start = clock();
	
	buildSegmentTree(memoryManager.vtermsCount);
	
	markVTerms(expr);
	
	swapBuffers();
	
	end = clock();
	
	printf("End garbage collection. Time elapsed: %f\n", ((float)(end - start)) / CLOCKS_PER_SEC);
}

static void swapBuffers()
{	
	struct v_term* vterms  = memoryManager.inactiveTermsHeap;	
	uint32_t dataOffset = 0;
	uint32_t newCount = 0;
	int i = 0;	
	
	for (i; i < memoryManager.vtermsCount; ++i)
	{
		if (sumInSegmentTree(i, i + 1) > 0)
		{
			//Копируем vterm
			memcpy((uint8_t*)(memoryManager.activeTermsHeap + i), 
				(uint8_t*)vterms,
				sizeof(struct v_term)
			);			
			++vterms;
			
			//Копируем данные vterm'а			
			dataOffset += copyVTermData(memoryManager.activeTermsHeap[i]);
		}
	}
	
	memoryManager.activeTermsHeap = memoryManager.inactiveTermsHeap;
	memoryManager.activeDataHeap = memoryManager.inactiveDataHeap;
	memoryManager.vtermsCount = vterms - memoryManager.inactiveTermsHeap;
	memoryManager.vtermsOffset = vterms - memoryManager.inactiveTermsHeap;
	memoryManager.dataOffset = offset;
}

static uint32_t copyVTermData(struct v_term* term)
{	
	switch (term.tag)
	{
		case V_TERM_SYMBOL_TAG:
			//Подправляем указатель на данные vterm'а
			term->symbol = memoryManager.inactiveDataHeap;
			return copySymbol(term.symbol)
			break;
			
		case V_TERM_BRACKETS_TAG:
			//В этом случае ничего делать не нужно
			return 0;
			break;
	}
	
	retrun 0;
}

//Возвращает сколько байтов было использовано
static uint32_t copySymbol(struct v_symbol* symbol)
{
	uint8_t* data = memoryManager.inactiveDataHeap;
	
	switch (symbol->tag)
	{
		case V_SYMBOL_CHAR_TAG:
			data[0] = symbol->str[0];
			++data;
			break;
			
		case V_SYMBOL_STR_TAG:
			uint32_t memSize = strlen(symbol->str) + 1;
			memcpy(symbol->str, data, memSize);
			data += memSize;
			break;
			
		case V_SYMBOL_NUMBER_TAG:
			((int*)data)[0] = symbol->number;
			++((int*)data);
			break;
			
		case V_SYMBOL_CLOSURE_TAG
			//Пока ничего не делаем
			break;
	}
	
	return data - memoryManager.inactiveDataHeap;
}

static void markVTerms(struct l_term* expr)
{	
	struct l_term* currTerm = expr;
	
	while (currTerm)
	{
		switch (currTerm->tag)
		{
			case L_TERM_RANGE_TAG:
				markInSegmentTree(currTerm->rang->offset, currTerm->rang->offset + currTerm->rang->length - 1);
				break;
				
			case L_TERM_CHAIN_TAG:
				markVTerms(currTerm->chain);
				break;
		}
		
		currTerm = currTerm->next;
	}
}

struct v_term* allocate(struct l_term* expr)
{

}
