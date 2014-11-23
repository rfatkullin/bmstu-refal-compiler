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
	/// Указатель на выделенный участок памяти
	uint8_t* mainHeap;

	/// Указатель на начало свободного места в куче.
	uint8_t*  currHeapPointer;

	/// Размер выделенного участка
	uint32_t totalSize;

	/// Указатель на начало свободного места в куче для
	/// литеральных v_term
	struct v_term* literalTermsHeap;

	struct v_term* activeTermsHeap;
	struct v_term* inactiveTermsHeap;

	struct segment_tree* segmentTree;

	uint8_t* constDataHeap;
	uint8_t* activeDataHeap;
	uint8_t* inactiveDataHeap;

	uint8_t* lTermsHeap;

	uint32_t vtermsOffset;
	uint32_t dataOffset;
	uint32_t ltermsOffset;

	uint32_t vtermsCount;

	//Количество элементов в листе дерева отрезков
	uint32_t segmentLen;

	//На какое число v_term'ов хватит памяти
	uint32_t maxVTermCount;
};

struct memory_manager memMngr;

/// Выделяет память размера size
/// и сохраняет указатель на выделенный участок
/// в переменной mainHeap.
void initAllocator(uint32_t size);

/// Распределеяет память для типов данных
/// т.е. инциализирует поля activeTermsHeap, inactiveTermsHeap и т.д.
void initHeaps(uint32_t newSegmentLen);

/// Собирает l_term и возвращает результат.
struct v_term* allocate(struct l_term* expr);

/// Собирает мусор.
void collectGarbage(struct l_term* expr);

/// Выделяет память под строку и возвращает результат.
struct l_term* allocateVector(int strLen, char* str);

#endif
