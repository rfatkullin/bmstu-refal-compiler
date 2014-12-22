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

	buildSegmentTree(memMngr.vtermsCount);

	markVTerms(expr);

	swapBuffers();

	end = clock();

	printf("End garbage collection. Time elapsed: %f\n", ((float)(end - start)) / CLOCKS_PER_SEC);
}

//TO FIX: сделать проверку переполнения памяти.
uint32_t allocateSymbol(char ch)
{
	struct v_term* term = memMngr.activeTermsHeap + memMngr.vtermsOffset;
	term->tag =V_CHAR_TAG;
	term->ch = ch;

	return memMngr.vtermsOffset++;
}

//TO FIX: сделать проверку переполнения памяти.
void allocateVTerms(struct fragment* frag)
{
	uint32_t i = 0;
	for (i = 0; i < frag->length; ++i)
	{
		memMngr.activeTermsHeap[memMngr.vtermsOffset].tag = memMngr.activeTermsHeap[frag->offset + i].tag;

		switch (memMngr.activeTermsHeap[frag->offset + i].tag)
		{
			case V_CHAR_TAG:
				memMngr.activeTermsHeap[memMngr.vtermsOffset].ch = memMngr.activeTermsHeap[frag->offset + i].ch;
				break;

			case V_IDENT_TAG :
				memMngr.activeTermsHeap[memMngr.vtermsOffset].str = memMngr.activeTermsHeap[frag->offset + i].str;
				break;

			case V_INT_NUM_TAG:
				memMngr.activeTermsHeap[memMngr.vtermsOffset].intNum = memMngr.activeTermsHeap[frag->offset + i].intNum;
				break;

			case V_FLOAT_NUM_TAG:
				memMngr.activeTermsHeap[memMngr.vtermsOffset].floatNum = memMngr.activeTermsHeap[frag->offset + i].floatNum;
				break;

			case V_CLOSURE_TAG:
				//TO DO:
				break;

			case V_BRACKET_TAG:
				memMngr.activeTermsHeap[memMngr.vtermsOffset].inBracketLength = memMngr.activeTermsHeap[frag->offset + i].inBracketLength;
				break;
		}
		memMngr.vtermsOffset++;
	}
}

//TO FIX: сделать проверку переполнения памяти.
uint32_t allocateBracketVTerm(uint32_t length)
{
	memMngr.activeTermsHeap[memMngr.vtermsOffset].tag = V_BRACKET_TAG;
	memMngr.activeTermsHeap[memMngr.vtermsOffset].inBracketLength = length;

	return memMngr.vtermsOffset++;
}

void changeBracketLength(uint32_t offset, uint32_t newLength)
{
	memMngr.activeTermsHeap[offset].inBracketLength = newLength;
}


struct l_term* allocateVector(int strLen, char* str)
{
	uint32_t i = 0;

	if (memMngr.vtermsOffset + strLen >= memMngr.maxVTermCount)
	{
		//TO FIX: нельзя передавать нулевой указатель.
		collectGarbage(NULL);

		if (memMngr.vtermsOffset + strLen >= memMngr.maxVTermCount)
		{
			printf("[Memory manager]Fatal error: Can't allocate memory!\n");
			exit(1);
		}
	}

	for (i = 0; i < strLen; ++i, ++memMngr.vtermsOffset)
	{
		struct v_term* term = memMngr.activeTermsHeap + memMngr.vtermsOffset;
		term->tag =V_CHAR_TAG;
		term->ch = str[i];
	}

	return allocateLTerm(memMngr.vtermsOffset - strLen, strLen);
}

void initAllocator(uint32_t size)
{
	memMngr.mainHeap = (uint8_t*)malloc(size);
	memMngr.currHeapPointer = memMngr.mainHeap;
	memMngr.totalSize = size;
	memMngr.literalTermsHeap = (struct v_term*)memMngr.mainHeap;
}

void initHeaps(uint32_t newSegmentLen)
{
	uint32_t size = memMngr.totalSize - (memMngr.literalTermsHeap - (struct v_term*)memMngr.mainHeap);
	uint32_t dataHeapSize = DATA_HEAP_SIZE_FACTOR * size;
	uint32_t vtermsHeapSize = V_TERMS_HEAP_SIZE_FACTOR * size;
	uint32_t ltermsHeapSize = size - dataHeapSize - vtermsHeapSize;
	uint32_t usedMemory = 0;
	memMngr.segmentLen = newSegmentLen;

	printf("\nAllocation size:                      %.2f Kb\n", byte2KByte(size));
	printf("\nAllocation ratios and sizes:         Ratio\t   Size\n");
	printf("\t For data:                    %.2f\t %.2f Kb\n", DATA_HEAP_SIZE_FACTOR, byte2KByte(dataHeapSize));
	printf("\t For vterms:                  %.2f\t %.2f Kb\n", V_TERMS_HEAP_SIZE_FACTOR, byte2KByte(vtermsHeapSize));
	printf("\t For lterms:                  %.2f\t %.2f Kb\n", 1.0f - (DATA_HEAP_SIZE_FACTOR + V_TERMS_HEAP_SIZE_FACTOR), byte2KByte(ltermsHeapSize));

	usedMemory += allocateMemoryForData(dataHeapSize);
	usedMemory += allocateMemoryForVTerms(vtermsHeapSize);
	usedMemory += allocateMemoryForLTerms(ltermsHeapSize);

	assert(usedMemory < size);

	memMngr.vtermsCount = 0;
	memMngr.vtermsOffset = 0;
	memMngr.dataOffset = 0;
	memMngr.ltermsOffset = 0;

	printf("Total used memory:                    %.2f Kb\n", byte2KByte(usedMemory));

	memMngr.literalVTermsNumber = memMngr.literalTermsHeap - (struct v_term*)memMngr.mainHeap;
	memMngr.literalTermsHeap = (struct v_term*)memMngr.mainHeap;
}

void debugLiteralsPrint()
{
	printf("vterms debug print:\n\t");
	int i;
	for (i = 0; i < memMngr.literalVTermsNumber; ++i)
	{
		printSymbol(memMngr.literalTermsHeap + i);
	}

	printf("\n");
}

static struct l_term* allocateLTerm(uint32_t offset, uint32_t len)
{
	struct l_term* term = (struct l_term*)malloc(sizeof(struct l_term));

	term->tag = L_TERM_FRAGMENT_TAG;
	term->fragment = (struct fragment*)malloc(sizeof(struct fragment));
	term->fragment->offset = offset;
	term->fragment->length = len;

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
	uint32_t chunck = memMngr.segmentLen;
	uint32_t memSizeWithoutHeader = size - sizeof(struct segment_tree);
	uint32_t n = (chunck * memSizeWithoutHeader) / (2 * chunck * sizeof(struct v_term) + (4 + chunck) * sizeof(uint32_t));
	uint32_t memSizeForTree = 4 * n / chunck * sizeof(uint32_t);
	uint32_t memSizeForElements = n * sizeof(uint32_t);

	memMngr.segmentTree = (struct segment_tree*)(memMngr.currHeapPointer);
	memMngr.segmentTree->tree = (int32_t*)(memMngr.currHeapPointer + sizeof(struct segment_tree));
	memMngr.segmentTree->elements = (int32_t*)(memMngr.currHeapPointer + memSizeForTree + sizeof(struct segment_tree));
	memMngr.maxVTermCount = n;

	memMngr.currHeapPointer += memSizeForTree + memSizeForElements + sizeof(struct segment_tree);

	return memSizeForTree + memSizeForElements + sizeof(struct segment_tree);
}

static uint32_t allocateMemoryForVTerms(uint32_t size)
{
	uint32_t segmentTreeHeapSize = allocateMemoryForSegmentTree(size);
	uint32_t termsHeapSize = memMngr.maxVTermCount * sizeof(struct v_term);

	uint32_t usedMemory = 2 * termsHeapSize + segmentTreeHeapSize;

	assert(usedMemory <= size);

	memMngr.activeTermsHeap = (struct v_term*)(memMngr.currHeapPointer);
	memMngr.inactiveTermsHeap = (struct v_term*)(memMngr.currHeapPointer + termsHeapSize);

	memMngr.currHeapPointer += 2 * termsHeapSize;

	printf("\nMemory allocation for vterms:\n");
	printf("\tMemory enough terms count:    %d \n", memMngr.maxVTermCount);
	printf("\tActive vterms heap size:      %.2f Kb\n", byte2KByte(termsHeapSize));
	printf("\tInactive vterms heap size:    %.2f Kb\n", byte2KByte(termsHeapSize));
	printf("\tSegment tree size:            %.2f Kb\n", byte2KByte(segmentTreeHeapSize));
	printf("\tUsed memory:                  %.2f Kb\n", byte2KByte(usedMemory));
	printf("\tLost memory:                  %.2f Kb\n", byte2KByte(size - usedMemory));

	return usedMemory;
}

static uint32_t allocateMemoryForData(uint32_t size)
{
	uint32_t singleDataHeapSize = size / 2;
	uint32_t usedMemory = 2 * singleDataHeapSize;

	memMngr.activeDataHeap = memMngr.currHeapPointer;
	memMngr.inactiveDataHeap = memMngr.currHeapPointer + singleDataHeapSize;
	memMngr.currHeapPointer +=usedMemory;

	printf("\nMemory allocation for data:\n");
	printf("\tMemory for single data heap:  %.2f Kb\n", byte2KByte(singleDataHeapSize));
	printf("\tUsed memory:                  %.2f Kb\n", byte2KByte(usedMemory));
	printf("\tLost memory:                  %.2f Kb\n", byte2KByte(size - usedMemory));

	return usedMemory;
}

static uint32_t allocateMemoryForLTerms(uint32_t size)
{
	memMngr.lTermsHeap = memMngr.currHeapPointer;
	memMngr.currHeapPointer += size;

	printf("\nMemory allocated for lterms:          %.2f Kb\n", byte2KByte(size));

	return size;
}

//Возвращает сколько байтов было использовано
static uint32_t copyVTerm(struct v_term* term)
{
	uint8_t* data = memMngr.inactiveDataHeap;
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

		case V_INT_NUM_TAG:
			((int*)data)[0] = term->intNum;
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
	struct v_term* vterms  = memMngr.inactiveTermsHeap;
	uint32_t dataOffset = 0;
	int i = 0;

	for (i; i < memMngr.vtermsCount; ++i)
	{
		if (sumInSegmentTree(i, i + 1) > 0)
		{
			//Копируем vterm
			memcpy((void*)(memMngr.activeTermsHeap + i),
				(void*)vterms,
				sizeof(struct v_term)
			);
			++vterms;

			//Копируем данные vterm'а
			dataOffset += copyVTerm(memMngr.activeTermsHeap + i);
		}
	}

	memMngr.activeTermsHeap = memMngr.inactiveTermsHeap;
	memMngr.activeDataHeap = memMngr.inactiveDataHeap;
	memMngr.vtermsCount = vterms - memMngr.inactiveTermsHeap;
	memMngr.vtermsOffset = vterms - memMngr.inactiveTermsHeap;
	memMngr.dataOffset = dataOffset;
}

static void markVTerms(struct l_term* expr)
{
	struct l_term* currTerm = expr;

	while (currTerm)
	{
		switch (currTerm->tag)
		{
			case L_TERM_FRAGMENT_TAG:
				markInSegmentTree(currTerm->fragment->offset, currTerm->fragment->offset + currTerm->fragment->length - 1);
				break;

			case L_TERM_CHAIN_TAG:
				markVTerms(currTerm->chain);
				break;
		}

		currTerm = currTerm->next;
	}
}
