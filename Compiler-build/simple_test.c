// file:simple_test.ref

#include <memory_manager.h>
#include <v_machine.h>

struct func_result_t Go(int entryPoint, struct env_t* env, struct field_view_t* fieldOfView) 
{
 if (entryPoint == 0)
 {
  env.locals = (struct l_term*)malloc(1 * sizeof(struct l_term));
  fieldOfView.backups = (struct l_term_chain_t*)malloc(1 * sizeof(struct l_term_chain_t));
 }
 switch (entryPoint)
 {
  case 0: 
 ok = 0;
 ok = 1;

 if (ok == 1)
%!s(MISSING){ %s}
 } // case block end
 if (res != CALL_RESULT)
 {
  free(env.locals);
  free(fieldOfView.backups);
 }
} // Go

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
 mainLoop(Go);
 return 0;
}
