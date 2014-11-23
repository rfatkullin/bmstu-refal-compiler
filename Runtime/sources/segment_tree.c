#include <string.h>
#include <math.h>

#include "segment_tree.h"

static void mark(uint32_t curr, uint32_t currBegin, uint32_t currEnd, uint32_t needBegin, uint32_t needEnd);
static uint32_t sum(uint32_t curr, uint32_t currBegin, uint32_t currEnd, uint32_t needBegin, uint32_t needEnd);

void buildSegmentTree(uint32_t n)
{
	memMngr.segmentTree->currN = n;
	uint32_t treeSize = 4 * memMngr.segmentTree->currN / memMngr.segmentLen;
	uint32_t i = 0;

	//Кажется, можно просто обнулять. Но пока работает)))
	for (i = 0; i < treeSize; ++i)
		memMngr.segmentTree->tree[i] = -1;

	memset(memMngr.segmentTree->elements, 0, sizeof(uint32_t) * memMngr.segmentTree->currN);
}

void clearSegmentTree()
{
	uint32_t treeSize = 4 * memMngr.segmentTree->currN / memMngr.segmentLen;
	uint32_t i = 0;

	for (i = 0; i < treeSize; ++i)
		memMngr.segmentTree->tree[i] = -1;

	memset(memMngr.segmentTree->elements, 0, sizeof(uint32_t) * memMngr.segmentTree->currN);
}

void markInSegmentTree(uint32_t begin, uint32_t end)
{
	mark(1, 0, memMngr.segmentTree->currN - 1, begin, end);
}

uint32_t sumInSegmentTree(uint32_t begin, uint32_t end)
{
	sum(1, 0, memMngr.segmentTree->currN - 1, begin, end);
}

static uint32_t max(uint32_t a, uint32_t b)
{
	if (a > b)
		return a;

	return b;
}

static uint32_t min(uint32_t a, uint32_t b)
{
	if (a > b)
		return b;

	return a;
}

static void markElements(uint32_t needBegin, uint32_t needEnd, int32_t mark)
{
	int i = 0;
	uint32_t segmentNum = needBegin / memMngr.segmentLen;

	for (i = segmentNum * memMngr.segmentLen; i < (segmentNum + 1) * memMngr.segmentLen; ++i)
	{
		if (i >= needBegin && i <= needEnd)
			memMngr.segmentTree->elements[i] = 1;
		else
			memMngr.segmentTree->elements[i] = mark;
	}
}

static uint32_t sumElements(uint32_t needBegin, uint32_t needEnd)
{
	uint32_t i = 0;
	uint32_t sum = 0;

	for (i = needBegin; i <= needEnd; ++i)
		sum += memMngr.segmentTree->elements[i];

	return sum;
}

static void push(uint32_t curr)
{
	if (memMngr.segmentTree->tree[curr] == -1)
				return;

	memMngr.segmentTree->tree[curr * 2] = memMngr.segmentTree->tree[curr];
	memMngr.segmentTree->tree[curr * 2 + 1] = memMngr.segmentTree->tree[curr];
	memMngr.segmentTree->tree[curr] = -1;
}

static uint32_t left(uint32_t val)
{
	return val * memMngr.segmentLen;
}

static uint32_t right(uint32_t val)
{
	return (val + 1) * memMngr.segmentLen - 1;
}

static void mark(uint32_t curr, uint32_t currBegin, uint32_t currEnd, uint32_t needBegin, uint32_t needEnd)
{
	if (needBegin > needEnd)
		return;

	//Наткнулись на узел, который совпадает с текущим отрезком
	if (needBegin == left(currBegin) && needEnd == right(currEnd))
	{
		memMngr.segmentTree->tree[curr] = 1;
		return;
	}

	//Дошли до листа
	if (currBegin == currEnd)
	{
		uint32_t label = 0;

		if (memMngr.segmentTree->tree[currBegin] > 0)
			label = 1;

		markElements(needBegin, needEnd, label);
		memMngr.segmentTree->tree[curr] = -1;
		return;
	}

	push(curr);

	uint32_t currMiddle = (currBegin + currEnd) / 2;

	mark(curr * 2, currBegin, currMiddle, needBegin, min(right(currMiddle), needEnd));
	mark(curr * 2 + 1, currMiddle + 1, currEnd, max(left(currMiddle + 1), needBegin), needEnd);
}

static uint32_t sum(uint32_t curr, uint32_t currBegin, uint32_t currEnd, uint32_t needBegin, uint32_t needEnd)
{
	if (needBegin > needEnd)
		return 0;

	if (memMngr.segmentTree->tree[curr] != -1)
		return memMngr.segmentTree->tree[curr] * (needEnd - needBegin + 1);

	//Дошли до листа
	if (currBegin == currEnd)
		return sumElements(needBegin, needEnd);

	uint32_t currMiddle = (currBegin + currEnd) / 2;

	return sum(curr * 2, currBegin, currMiddle, needBegin, min(right(currMiddle), needEnd)) +
		sum(curr * 2 + 1, currMiddle + 1, currEnd, max(left(currMiddle + 1), needBegin), needEnd);
}
