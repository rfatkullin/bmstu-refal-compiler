// file:simple_test.ref

#include <memory_manager.h>
#include <v_machine.h>

struct func_result_t Go(int entryPoint, struct env_t* env, struct field_view_t* fieldOfView) 
{
 struct fresult_t result;
 if (entryPoint == 0)
 {
  env.locals = (struct l_term*)malloc(1 * sizeof(struct l_term));
  fieldOfView.backups = (struct l_term_chain_t*)malloc(1 * sizeof(struct l_term_chain_t));
 }
 switch (entryPoint)
 {
  case 0: 

   struct l_term* currTerm = 0;
   struct l_term* tmpTerm = 0;
   struct l_term_chain_t* chain0 = (struct l_term_chain_t*)malloc(sizeof(struct l_term_chain_t));
   tmpTerm = currTerm;
   currTerm = (struct l_term*)malloc(sizeof(struct l_term));
   currTerm->tag = L_TERM_FRAGMENT_TAG;
   currTerm->fragment = (struct fragment*)malloc(sizeof(struct fragment));
   if (tmpTerm != 0) {
    tmpTerm->next = currTerm;
    currTerm->prev = tmpTerm;
   }
   chainTerm0->chain->begin = currTerm;
   currTerm->fragment->offset = 0;
   currTerm->fragment->length = 1;
   struct l_term* chainTerm1 = (struct l_term*)malloc(struct l_term);
   chainTerm1->tag = L_TERM_CHAIN_TAG;
   chainTerm1->chain = (struct l_term_chain_t*)malloc(struct l_term_chain_t);
   chainTerm1->prev = currTerm;
   currTerm->next = chainTerm1;
   struct l_term* chainTerm2 = (struct l_term*)malloc(struct l_term);
   chainTerm2->tag = L_TERM_CHAIN_TAG;
   chainTerm2->chain = (struct l_term_chain_t*)malloc(struct l_term_chain_t);
   chainTerm1->chain->begin = chainTerm2;
   chainTerm2->prev = 0;
   tmpTerm = currTerm;
   currTerm = (struct l_term*)malloc(sizeof(struct l_term));
   currTerm->tag = L_TERM_FRAGMENT_TAG;
   currTerm->fragment = (struct fragment*)malloc(sizeof(struct fragment));
   if (tmpTerm != 0) {
    tmpTerm->next = currTerm;
    currTerm->prev = tmpTerm;
   }
   chainTerm2->chain->begin = currTerm;
   currTerm->fragment->offset = 0;
   currTerm->fragment->length = 2;
   chainTerm2->chain->end = currTerm;
   currTerm = chainTerm2;
   tmpTerm = currTerm;
   currTerm = (struct l_term*)malloc(sizeof(struct l_term));
   currTerm->tag = L_TERM_FRAGMENT_TAG;
   currTerm->fragment = (struct fragment*)malloc(sizeof(struct fragment));
   if (tmpTerm != 0) {
    tmpTerm->next = currTerm;
    currTerm->prev = tmpTerm;
   }
   currTerm->fragment->offset = 0;
   currTerm->fragment->length = 1;
   chainTerm1->chain->end = currTerm;
   currTerm = chainTerm1;
   chainTerm0->end = currTerm;
  break;
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

 initHeaps(2);
} // __initLiteralData()

int main()
{
 __initLiteralData();
 mainLoop(Go);
 return 0;
}
