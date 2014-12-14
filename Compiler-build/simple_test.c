// file:simple_test.ref

#include <memory_manager.h>
#include <v_machine.h>

struct func_result_t Func(int entryPoint, struct env_t* env, struct field_view_t* fieldOfView) 
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

   result.status = OK_RESULT;
   result.mainChain = 0;
   result.callChain = 0;
  break;
 } // case block end
 if (res != CALL_RESULT)
 {
  free(env.locals);
  free(fieldOfView.backups);
 }
} // Func

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

   struct l_term* helper = (struct l_term*)malloc(2 * sizeof(struct l_term));
   int i;
   for (i = 0; i < 2; ++i)
   {
    helper[i]->tag = L_TERM_CHAIN_TAG;
    helper[i]->chain = (struct l_term_chain_t*)malloc(sizeof(struct l_term_chain_t));
   }
   struct l_term* currTerm = 0;
   /*Start expr 1 with 1 terms*/;
   currTerm = (struct l_term*)malloc(sizeof(struct l_term));
   currTerm->tag = L_TERM_FRAGMENT_TAG;
   currTerm->fragment = (struct fragment*)malloc(sizeof(struct fragment));
   currTerm->fragment->offset = 0;
   currTerm->fragment->length = 1;
   currTerm->prev = 0;
   helper[1]->chain->begin = currTerm;
   helper[1]->chain->end = currTerm;
   helper[1]->chain->end->next = 0;
   helper[0]->chain->begin = helper[1];
   helper[0]->chain->end = helper[1];
   /*End expr 1*/;
   helper[0]->chain->end->next = 0;
   /*End expr 0*/;
   helper[0]->chain->begin->prev = 0;
   helper[0]->chain->end->next = 0;
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

 initHeaps(2);
} // __initLiteralData()

int main()
{
 __initLiteralData();
 mainLoop(Go);
 return 0;
}
