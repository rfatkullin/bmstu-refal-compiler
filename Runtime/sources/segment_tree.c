#include <string.h>
#include <math.h>

#include "segment_tree.h"

static void mark(uint32_t curr, uint32_t currBegin, uint32_t currEnd, uint32_t needBegin, uint32_t needEnd);
static uint32_t sum(uint32_t curr, uint32_t currBegin, uint32_t currEnd, uint32_t needBegin, uint32_t needEnd);

void buildSegmentTree(uint32_t n)
{
	int treeSize = 4 * n / N;	
    memset(memoryManager.segmentTree->tree, 0, sizeof(uint32_t) * treeSize);
    memset(memoryManager.segmentTree->elements, 0, sizeof(uint32_t) * memoryManager.segmentTree->n);
}

void markInSegmentTree(uint32_t begin, uint32_t end)
{
    sum(1, 0, memoryManager.segmentTree->n, begin, end);
} 

uint32_t sumInSegmentTree(uint32_t begin, uint32_t end)
{
    sum(1, 0, memoryManager.segmentTree->n, begin, end);
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


static void push(uint32_t curr)
{
    if (memoryManager.segmentTree->tree[curr] == 0)
		return;
	
    memoryManager.segmentTree->tree[curr * 2] = memoryManager.segmentTree->tree[curr];
    memoryManager.segmentTree->tree[curr * 2 + 1] = memoryManager.segmentTree->tree[curr];
    memoryManager.segmentTree->tree[curr] = 0;
}

static void mark(uint32_t curr, uint32_t currBegin, uint32_t currEnd, uint32_t needBegin, uint32_t needEnd)
{
	if (needBegin > needEnd)
		return;
	
	if (currBegin == needBegin && currEnd == needEnd)
	{
        memoryManager.segmentTree->tree[curr] = 1;
		return;
	}
	
	push(curr);
	
	uint32_t currMiddle = (currBegin + currEnd) / 2;
	
	mark(curr * 2, currBegin, currMiddle, needBegin, min(currMiddle, needEnd));
	mark(curr * 2 + 1, currMiddle + 1, currEnd, max(currMiddle + 1, needBegin), needEnd);
}

static uint32_t sum(uint32_t curr, uint32_t currBegin, uint32_t currEnd, uint32_t needBegin, uint32_t needEnd)
{
	if (needBegin > needEnd)
		return 0;
	
    if (memoryManager.segmentTree->tree[curr] == 1)
		return needEnd - needBegin + 1;
	
	uint32_t currMiddle = (currBegin + currEnd) / 2;
	
	return sum(curr * 2, currBegin, currMiddle, needBegin, min(currMiddle, needEnd)) +	
		sum(curr * 2 + 1, currMiddle + 1, currEnd, max(currMiddle + 1, needBegin), needEnd);
}
