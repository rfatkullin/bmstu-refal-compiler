#include "segment_tree.h"

void buildSegmentTree(uint32_t n)
{
	int treeSize = 4 * n / N;	
	memset(memoryManager.segmentTree, 0, sizeof(uint32_t) * treeSize);	
}

void markInSegmentTree(uint32_t begin, uint32_t end)
{
		
} 

uint32_t sumInSegmentTreee(uint32_t begin, uint32_t end)
{
	
}