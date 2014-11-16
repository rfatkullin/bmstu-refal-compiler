#include <stdio.h>
#include <stdlib.h>
#include <assert.h>
#include <time.h>
#include <string.h>

#include "memory_manager.h"

static float byte2KByte(uint32_t bytes);
static void swapBuffers();
static void markVTerms(struct l_term* expr);

static uint32_t allocateMemoryForSegmentTree(uint32_t size);
static uint32_t allocateMemoryForVTerms(uint32_t size);
static uint32_t allocateMemoryForData(uint32_t size);
static uint32_t allocateMemoryForLTerms(uint32_t size);
static struct l_term* allocateLTerm(uint32_t offset, uint32_t len);

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

struct v_term* allocate(struct l_term* expr)
{

}

struct l_term* allocateVector(int strLen, char* str)
{
	uint32_t i = 0;

	if (memoryManager.vtermsOffset + strLen >= memoryManager.maxVTermCount)
	{
		//TO FIX: нельзя передавать нулевой указатель.
		collectGarbage(NULL);

		if (memoryManager.vtermsOffset + strLen >= memoryManager.maxVTermCount)
		{
			printf("[Memory manager]Fatal error: Can't allocate memory!\n");
			exit(1);
		}
	}

	for (i = 0; i < strLen; ++i, ++memoryManager.vtermsOffset)
	{
		struct v_term* term = memoryManager.activeTermsHeap + memoryManager.vtermsOffset;
		term->tag =V_CHAR_TAG;
		term->ch = str[i];
	}

	return allocateLTerm(memoryManager.vtermsOffset - strLen, strLen);
}

void initAllocator(uint32_t size, uint32_t N)
{
	uint32_t dataHeapSize = DATA_HEAP_SIZE_FACTOR * size;
	uint32_t vtermsHeapSize = V_TERMS_HEAP_SIZE_FACTOR * size;
	uint32_t ltermsHeapSize = size - dataHeapSize - vtermsHeapSize;
	uint32_t usedMemory = 0;
	memoryManager.SegmentLen = N;

	printf("\nAllocation size:                      %.2f Kb\n", byte2KByte(size));
	printf("\nAllocation ratios and sizes:         Ratio\t   Size\n");
	printf("\t For data:                    %.2f\t %.2f Kb\n", DATA_HEAP_SIZE_FACTOR, byte2KByte(dataHeapSize));
	printf("\t For vterms:                  %.2f\t %.2f Kb\n", V_TERMS_HEAP_SIZE_FACTOR, byte2KByte(vtermsHeapSize));
	printf("\t For lterms:                  %.2f\t %.2f Kb\n", 1.0f - (DATA_HEAP_SIZE_FACTOR + V_TERMS_HEAP_SIZE_FACTOR), byte2KByte(ltermsHeapSize));

	usedMemory += allocateMemoryForData(dataHeapSize);
	usedMemory += allocateMemoryForVTerms(vtermsHeapSize);
	usedMemory += allocateMemoryForLTerms(ltermsHeapSize);

	assert(usedMemory < size);

	memoryManager.totalSize = usedMemory;
	memoryManager.vtermsCount = 0;

	memoryManager.vtermsOffset = 0;
	memoryManager.dataOffset = 0;
	memoryManager.ltermsOffset = 0;

	printf("Total used memory:                    %.2f Kb\n", byte2KByte(usedMemory));
}

static struct l_term* allocateLTerm(uint32_t offset, uint32_t len)
{
	struct l_term* term = (struct l_term*)malloc(sizeof(struct l_term));

	term->tag = L_TERM_RANGE_TAG;
	term->range = (struct fragment*)malloc(sizeof(struct fragment));
	term->range->offset = offset;
	term->range->length = len;

	return term;
}

static float byte2KByte(uint32_t bytes)
{
	return bytes / 1024.0f;
}

//Значение n выводится из формулы:
//size = 2 * n * sizeof(struct v_term) + (4 * n / N + n) * sizeof(uint32_t)
static uint32_t allocateMemoryForSegmentTree(uint32_t size)
{
	uint32_t chunck = memoryManager.SegmentLen;
	uint32_t memSizeWithoutHeader = size - sizeof(struct segment_tree);
	uint32_t n = (chunck * memSizeWithoutHeader) / (2 * chunck * sizeof(struct v_term) + (4 + chunck) * sizeof(uint32_t));
	uint32_t memSizeForTree = 4 * n / chunck * sizeof(uint32_t);
	uint32_t memSizeForElements = n * sizeof(uint32_t);

	memoryManager.segmentTree = (struct segment_tree*)malloc(sizeof(struct segment_tree));
	memoryManager.segmentTree->tree = (int32_t*)malloc(memSizeForTree);
	memoryManager.segmentTree->elements = (int32_t*)malloc(memSizeForElements);
	memoryManager.maxVTermCount = n;

	return memSizeForTree + memSizeForElements + sizeof(struct segment_tree);
}

static uint32_t allocateMemoryForVTerms(uint32_t size)
{
	uint32_t segmentTreeHeapSize = allocateMemoryForSegmentTree(size);
	uint32_t activeTermsHeapSize = memoryManager.maxVTermCount * sizeof(struct v_term);
	uint32_t inactiveTermsHeapSize = memoryManager.maxVTermCount * sizeof(struct v_term);

	uint32_t usedMemory = activeTermsHeapSize + inactiveTermsHeapSize + segmentTreeHeapSize;

	assert(usedMemory <= size);

	memoryManager.activeTermsHeap = (struct v_term*)malloc(activeTermsHeapSize);
	memoryManager.inactiveTermsHeap = (struct v_term*)malloc(inactiveTermsHeapSize);


	printf("\nMemory allocation for vterms:\n");
	printf("\tMemory enough terms count:    %d \n", memoryManager.maxVTermCount);
	printf("\tActive vterms heap size:      %.2f Kb\n", byte2KByte(activeTermsHeapSize));
	printf("\tInactive vterms heap size:    %.2f Kb\n", byte2KByte(inactiveTermsHeapSize));
	printf("\tSegment tree size:            %.2f Kb\n", byte2KByte(segmentTreeHeapSize));
	printf("\tUsed memory:                  %.2f Kb\n", byte2KByte(usedMemory));
	printf("\tLost memory:                  %.2f Kb\n", byte2KByte(size - usedMemory));

	return usedMemory;
}

static uint32_t allocateMemoryForData(uint32_t size)
{
	uint32_t singleDataHeapSize = size / 2;
	uint32_t usedMemory = 2 * singleDataHeapSize;

	memoryManager.activeDataHeap = (uint8_t*)malloc(singleDataHeapSize);
	memoryManager.inactiveDataHeap = (uint8_t*)malloc(singleDataHeapSize);

	printf("\nMemory allocation for data:\n");
	printf("\tMemory for single data heap:  %.2f Kb\n", byte2KByte(singleDataHeapSize));
	printf("\tUsed memory:                  %.2f Kb\n", byte2KByte(usedMemory));
	printf("\tLost memory:                  %.2f Kb\n", byte2KByte(size - usedMemory));

	return usedMemory;
}

static uint32_t allocateMemoryForLTerms(uint32_t size)
{
	memoryManager.lTermsHeap = (uint8_t*)malloc(size);

	printf("\nMemory allocated for lterms:          %.2f Kb\n", byte2KByte(size));

	return size;
}

//Возвращает сколько байтов было использовано
static uint32_t copyVTerm(struct v_term* term)
{
	uint8_t* data = memoryManager.inactiveDataHeap;
	uint32_t memSize = 0;

	switch (term->tag)
	{
		case V_CHAR_TAG:
			data[0] = term->str[0];
			memSize = 1;
			break;

		case V_IDENT_TAG:
			memSize = strlen(term->str) + 1;
			memcpy(term->str, data, memSize);
			break;

		case V_NUMBER_TAG:
			((int*)data)[0] = term->num;
			memSize = sizeof(int);
			break;

		case V_CLOSURE_TAG:
			//Пока ничего не делаем
			break;
	}

	return memSize;
}

static void swapBuffers()
{
	struct v_term* vterms  = memoryManager.inactiveTermsHeap;
	uint32_t dataOffset = 0;
	int i = 0;

	for (i; i < memoryManager.vtermsCount; ++i)
	{
		if (sumInSegmentTree(i, i + 1) > 0)
		{
			//Копируем vterm
			memcpy((void*)(memoryManager.activeTermsHeap + i),
				(void*)vterms,
				sizeof(struct v_term)
			);
			++vterms;

			//Копируем данные vterm'а
			dataOffset += copyVTerm(memoryManager.activeTermsHeap + i);
		}
	}

	memoryManager.activeTermsHeap = memoryManager.inactiveTermsHeap;
	memoryManager.activeDataHeap = memoryManager.inactiveDataHeap;
	memoryManager.vtermsCount = vterms - memoryManager.inactiveTermsHeap;
	memoryManager.vtermsOffset = vterms - memoryManager.inactiveTermsHeap;
	memoryManager.dataOffset = dataOffset;
}

static void markVTerms(struct l_term* expr)
{
	struct l_term* currTerm = expr;

	while (currTerm)
	{
		switch (currTerm->tag)
		{
			case L_TERM_RANGE_TAG:
				markInSegmentTree(currTerm->range->offset, currTerm->range->offset + currTerm->range->length - 1);
				break;

			case L_TERM_CHAIN_TAG:
				markVTerms(currTerm->chain);
				break;
		}

		currTerm = currTerm->next;
	}
}
