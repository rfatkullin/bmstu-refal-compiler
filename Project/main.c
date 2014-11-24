// file:../Compiler-build/simple_test.ref

#include <memory_manager.h>
void __initLiteralData()
{
 initAllocator(1024 * 1024 * 1024);
 *(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_INT_NUM_TAG, .intNum = 1};
 *(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_INT_NUM_TAG, .intNum = 23};
 *(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_INT_NUM_TAG, .intNum = 123};
 *(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_IDENT_TAG, .str = "qweqw"};
 *(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_CHAR_TAG, .ch = 'A'};
 *(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_CHAR_TAG, .ch = 'A'};
 *(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_CHAR_TAG, .ch = 'A'};
 *(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_CHAR_TAG, .ch = 'A'};
 *(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_CHAR_TAG, .ch = 'A'};
 *(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_FLOAT_NUM_TAG, .floatNum = 12.000000};
 *(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_IDENT_TAG, .str = "d"};
 *(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_IDENT_TAG, .str = "hell"};
 *(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_IDENT_TAG, .str = "asdasd"};
 *(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_IDENT_TAG, .str = "asdasd"};
 *(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_IDENT_TAG, .str = "a"};
 *(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_INT_NUM_TAG, .intNum = 111};

 initHeaps(2);
} // __initLiteralData()

int main()
{
 __initLiteralData();
 return 0;
}
