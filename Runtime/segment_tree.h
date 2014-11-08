#ifndef __SEGMENT_TREE_H__
#define __SEGMENT_TREE_H__

#include "memory_manager.h"

void buildSegmentTree(uint32_t n);
void markInSegmentTree(uint32_t begin, uint32_t end); 
uint32_t sumInSegmentTree(uint32_t begin, uint32_t end);

#endif