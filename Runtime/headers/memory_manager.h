#ifndef __ALLOCATOR_H__
#define __ALLOCATOR_H__

#include <stdint.h>

#include "v_term.h"
#include "l_term.h"
#include "segment_tree.h"

//Сколько процентов памяти выдялется тому или иному типу данных
//Память для дерева отрезков берется из памяти для v_term'ов.
#define V_TERMS_HEAP_SIZE_FACTOR 		0.3f
#define DATA_HEAP_SIZE_FACTOR 			0.4f
#define L_TERMS_HEAP_SIZE_FACTOR		0.3f

struct memory_manager
{
	struct v_term* constTermsHeap;
	struct v_term* activeTermsHeap;
	struct v_term* inactiveTermsHeap;

	struct segment_tree* segmentTree;

	uint8_t* constDataHeap;
	uint8_t* activeDataHeap;
	uint8_t* inactiveDataHeap;

	uint8_t* lTermsHeap;

	uint32_t constTermOffset;
	uint32_t constDataOffset;
	
	uint32_t vtermsOffset;
	uint32_t dataOffset;
	uint32_t ltermsOffset;

	uint32_t vtermsCount;
	uint32_t totalSize;

	//Количество элементов в листе дерева отрезков
	uint32_t SegmentLen;

	//На какое число v_term'ов хватит памяти
	uint32_t maxVTermCount;
};

struct memory_manager memoryManager;

void initAllocator(uint32_t size, uint32_t N);
void markTerms(struct l_term* term);
struct v_term* allocate(struct l_term* expr);
void collectGarbage(struct l_term* expr);
struct l_term* allocateVector(int strLen, char* str);

#endif
