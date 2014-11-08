#ifndef __SEGMENT_TREE_H__
#define __SEGMENT_TREE_H__

#include "memory_manager.h"

struct segment_tree
{
	uint32_t n;
	int32_t* tree;
	uint32_t* elements;
};

void buildSegmentTree(uint32_t n);
void clearSegmentTree();
void markInSegmentTree(uint32_t begin, uint32_t end);
uint32_t sumInSegmentTree(uint32_t begin, uint32_t end);

#endif
