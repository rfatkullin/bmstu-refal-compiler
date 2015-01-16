#include <stdio.h>
#include <stdlib.h>
#include <assert.h>
#include <time.h>
#include <string.h>

#include "memory_manager.h"

static float byte2KByte(uint32_t bytes);
static void swapBuffers();
static void markVTerms(struct lterm_t* expr);

static void allocateMemoryForSegmentTree(uint32_t termsNumber, uint8_t** pointer);
static void allocateMemoryForVTerms(uint32_t size, uint8_t** pointer);
static void allocateMemoryForData(uint32_t size, uint8_t** pointer);
static void allocateMemoryForLTerms(uint32_t size, uint8_t** pointer);
static struct lterm_t* allocateLTerm(uint32_t offset, uint32_t len);

void collectGarbage(struct lterm_t* expr)
{
	clock_t start, end;
	printf("Start garbage collection.\n");
	start = clock();

	buildSegmentTree(memMngr.vtermsOffset);

	markVTerms(expr);

	swapBuffers();

	end = clock();

	printf("End garbage collection. Time elapsed: %f\n", ((float)(end - start)) / CLOCKS_PER_SEC);
}

//TO FIX: сделать проверку переполнения памяти.
uint32_t allocateSymbol(char ch)
{
	struct v_term* term = memMngr.vterms + memMngr.vtermsOffset;
	term->tag =V_CHAR_TAG;
	term->ch = ch;

	return memMngr.vtermsOffset++;
}

//TO FIX: сделать проверку переполнения памяти.
void allocateVTerms(struct fragment_t* frag)
{
	uint32_t i = 0;
	for (i = 0; i < frag->length; ++i)
	{
		memMngr.vterms[memMngr.vtermsOffset].tag = memMngr.vterms[frag->offset + i].tag;

		switch (memMngr.vterms[frag->offset + i].tag)
		{
			case V_CHAR_TAG:
				memMngr.vterms[memMngr.vtermsOffset].ch = memMngr.vterms[frag->offset + i].ch;
				break;

			case V_IDENT_TAG :
				memMngr.vterms[memMngr.vtermsOffset].str = memMngr.vterms[frag->offset + i].str;
				break;

			case V_INT_NUM_TAG:
				memMngr.vterms[memMngr.vtermsOffset].intNum = memMngr.vterms[frag->offset + i].intNum;
				break;

			case V_FLOAT_NUM_TAG:
				memMngr.vterms[memMngr.vtermsOffset].floatNum = memMngr.vterms[frag->offset + i].floatNum;
				break;

			case V_CLOSURE_TAG:
				//TO DO:
				break;

			case V_BRACKET_TAG:
				memMngr.vterms[memMngr.vtermsOffset].inBracketLength = memMngr.vterms[frag->offset + i].inBracketLength;
				break;
		}
		memMngr.vtermsOffset++;
	}
}

//TO FIX: сделать проверку переполнения памяти.
uint32_t allocateBracketVTerm(uint32_t length)
{
	memMngr.vterms[memMngr.vtermsOffset].tag = V_BRACKET_TAG;
	memMngr.vterms[memMngr.vtermsOffset].inBracketLength = length;

	return memMngr.vtermsOffset++;
}

void changeBracketLength(uint32_t offset, uint32_t newLength)
{
	memMngr.vterms[offset].inBracketLength = newLength;
}

struct lterm_t* allocateVector(int strLen, char* str)
{
	uint32_t i = 0;

	if (memMngr.vtermsOffset + strLen >= memMngr.maxTermsNumber)
	{
		//TO FIX: нельзя передавать нулевой указатель.
		collectGarbage(NULL);

		if (memMngr.vtermsOffset + strLen >= memMngr.maxTermsNumber)
		{
			printf("[Memory manager]Fatal error: Can't allocate memory!\n");
			exit(1);
		}
	}

	for (i = 0; i < strLen; ++i, ++memMngr.vtermsOffset)
	{
		struct v_term* term = memMngr.vterms + memMngr.vtermsOffset;
		term->tag =V_CHAR_TAG;
		term->ch = str[i];
	}

	return allocateLTerm(memMngr.vtermsOffset - strLen, strLen);
}

void initAllocator(uint32_t size)
{
	memMngr.vterms = (struct v_term*)malloc(size);
	memMngr.totalSize = size;
}

void initHeaps(uint32_t segmentLen, uint32_t literalsNumber)
{
	//TO FIX: Если инициализируем данные, то их тоже нужно тут учитывать.
	uint32_t size = memMngr.totalSize - literalsNumber * sizeof(struct v_term);
	uint32_t dataHeapSize = DATA_HEAP_SIZE_FACTOR * size;
	uint32_t vtermsHeapSize = V_TERMS_HEAP_SIZE_FACTOR * size;
	uint32_t ltermsHeapSize = size - dataHeapSize - vtermsHeapSize;


	printf("\nAllocation size:                      %.2f Kb\n", byte2KByte(size));
	printf("\nAllocation ratios and sizes:         Ratio\t   Size\n");
	printf("\t For data:                    %.2f\t %.2f Kb\n", DATA_HEAP_SIZE_FACTOR, byte2KByte(dataHeapSize));
	printf("\t For vterms:                  %.2f\t %.2f Kb\n", V_TERMS_HEAP_SIZE_FACTOR, byte2KByte(vtermsHeapSize));
	printf("\t For lterms:                  %.2f\t %.2f Kb\n", 1.0f - (DATA_HEAP_SIZE_FACTOR + V_TERMS_HEAP_SIZE_FACTOR), byte2KByte(ltermsHeapSize));

	memMngr.segmentLen = segmentLen;
	memMngr.literalsNumber = literalsNumber;

	uint8_t* pointer = (uint8_t*)memMngr.vterms;

	allocateMemoryForVTerms(vtermsHeapSize, &pointer);
	allocateMemoryForData(dataHeapSize, &pointer);
	allocateMemoryForLTerms(ltermsHeapSize, &pointer);

	memMngr.vtermsOffset = memMngr.activeOffset;
	memMngr.dataOffset = 0;
	memMngr.ltermsOffset = 0;
}

void debugLiteralsPrint()
{
	printf("vterms debug print:\n\t");
	int i;
	for (i = 0; i < memMngr.literalsNumber; ++i)
	{
		printSymbol(memMngr.vterms + i);
	}

	printf("\n");
}

static struct lterm_t* allocateLTerm(uint32_t offset, uint32_t len)
{
	struct lterm_t* term = (struct lterm_t*)malloc(sizeof(struct lterm_t));

	term->tag = L_TERM_FRAGMENT_TAG;
	term->fragment = (struct fragment_t*)malloc(sizeof(struct fragment_t));
	term->fragment->offset = offset;
	term->fragment->length = len;

	return term;
}

static float byte2KByte(uint32_t bytes)
{
	return bytes / 1024.0f;
}

// Значение n выводится из формулы:
// size = 2 * n * sizeof(struct v_term) + (4 * n / N + n) * sizeof(uint32_t)
// где N - длина отрезка в листьях дерева.
static uint32_t getTermsMaxNumber(uint32_t size)
{
	uint32_t memSizeWithoutHeader = size - sizeof(struct segment_tree);
	uint32_t N = memMngr.segmentLen;
	return (N * memSizeWithoutHeader) / (2 * N * sizeof(struct v_term) + (4 + N) * sizeof(uint32_t));
}

static void allocateMemoryForSegmentTree(uint32_t termsNumber, uint8_t** pointer)
{
	uint32_t chunck = memMngr.segmentLen;
	uint32_t memSizeForTree = 4 * termsNumber / chunck * sizeof(uint32_t);
	uint32_t memSizeForElements = termsNumber * sizeof(uint32_t);

	memMngr.segmentTree = (struct segment_tree*)(*pointer);
	memMngr.segmentTree->tree = (int32_t*)(*pointer + sizeof(struct segment_tree));
	memMngr.segmentTree->elements = (int32_t*)(*pointer + memSizeForTree + sizeof(struct segment_tree));

	*pointer += memSizeForTree + memSizeForElements + sizeof(struct segment_tree);
}

static void allocateMemoryForVTerms(uint32_t size, uint8_t** pointer)
{
	memMngr.maxTermsNumber = getTermsMaxNumber(size);
	memMngr.activeOffset = memMngr.literalsNumber;
	memMngr.inactiveOffset = memMngr.activeOffset + memMngr.maxTermsNumber;

	*pointer += memMngr.literalsNumber * sizeof(struct v_term);
	*pointer += 2 * memMngr.maxTermsNumber * sizeof(struct v_term);

	allocateMemoryForSegmentTree(memMngr.maxTermsNumber, pointer);

	printf("\tMax terms: %d\n", memMngr.maxTermsNumber);
}

static void allocateMemoryForData(uint32_t size, uint8_t** pointer)
{
	uint32_t singleDataHeapSize = size / 2;

	memMngr.activeDataHeap = *pointer;
	memMngr.inactiveDataHeap = *pointer + singleDataHeapSize;

	*pointer += size;
}

static void allocateMemoryForLTerms(uint32_t size, uint8_t** pointer)
{
	memMngr.lTermsHeap = *pointer;

	*pointer += size;
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
	uint32_t dataOffset = 0;
	int i = 0;

//	for (i = memMngr.activeOffset; i < memMngr.vtermsOffset; ++i)
//	{
//		if (sumInSegmentTree(i, i + 1) > 0)
//		{
//			//Копируем vterm
//			memcpy((void*)(memMngr.activeTermsHeap + i),
//				(void*)vterms,
//				sizeof(struct v_term)
//			);
//			++vterms;

//			//Копируем данные vterm'а
//			dataOffset += copyVTerm(memMngr.activeTermsHeap + i);
//		}
//	}

//	memMngr.vtermsOffset = memMngr.activeOffset;

//	memMngr.activeTermsHeap = memMngr.inactiveTermsHeap;
//	memMngr.activeDataHeap = memMngr.inactiveDataHeap;
//	memMngr.vtermsCount = vterms - memMngr.inactiveTermsHeap;
//	memMngr.vtermsOffset = vterms - memMngr.inactiveTermsHeap;
//	memMngr.dataOffset = dataOffset;
}

static void markVTerms(struct lterm_t* expr)
{
	struct lterm_t* currTerm = expr;

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
