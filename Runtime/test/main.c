#include <stdio.h>
#include <stdlib.h>

#include "memory_manager.h"

int main(int argc, char** argv)
{
	initAllocator(1024 * 1024 * 1024);

	*(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_INT_NUM_TAG, .intNum = 100};

	initHeaps(2);

	if (argc > 1)
		freopen(argv[1], "r", stdin);
	else
		freopen("input.txt", "r", stdin);

	uint32_t i;
	uint32_t n;
	uint32_t begin;
	uint32_t end;
	uint32_t sum;
	uint32_t correctSum;

	buildSegmentTree(memMngr.maxVTermCount - 1);

	while (scanf("%u", &n) == 1)
	{
		clearSegmentTree();

		for (i = 0; i < n; ++i)
		{
			scanf("%u%u", &begin, &end);
			markInSegmentTree(begin, end);
		}

		scanf("%u%u%u", &begin, &end, &correctSum);

		sum = sumInSegmentTree(begin, end);

		if (sum == correctSum)
			printf("Ok!\n");
		else
			printf("Fail: expected sum: %u real sum: %u\n", correctSum, sum);
	}

	return 0;
}
