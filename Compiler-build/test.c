// file:test.ref

void __initLiteralData()
{
 initAllocator(1024 * 1024 * 1024);
 *(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_IDENT_TAG, .str = Goodbye};
 *(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_IDENT_TAG, .str = VAr};
 *(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_IDENT_TAG, .str = aasd};
 *(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_IDENT_TAG, .str = asdasd};
 *(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_INT_NUM_TAG, .intNum = 123};
 *(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_INT_NUM_TAG, .intNum = 231};

 initHeaps(2);
} // __initLiteralData()

int main()
{
 return 0;
}
