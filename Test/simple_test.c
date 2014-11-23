// file:simple_test.ref

l_term* Go(vec_header* vecData) 
{
 struct v_term* data = vecData.data;
 uint32_t length = vecData.size;
 uint32_t ok = 0;
 ok = 0;
 ok = 1;

 if (ok == 1)
%!s(MISSING){%!s(MISSING)}
} // Go

int main()
{
 memoryManager.constTermOffset = 0;
 int i = 0;

 memoryManager.constTermsHeap[memoryManager.constTermOffset++] = {.tag = V_INT_NUM_TAG, .intNum = 1};
 memoryManager.constTermsHeap[memoryManager.constTermOffset++] = {.tag = V_INT_NUM_TAG, .intNum = 23};
 memoryManager.constTermsHeap[memoryManager.constTermOffset++] = {.tag = V_INT_NUM_TAG, .intNum = 123};
 memoryManager.constTermsHeap[memoryManager.constTermOffset++] = {.tag = V_IDENT_TAG, .str = qweqw};
 memoryManager.constTermsHeap[memoryManager.constTermOffset++] = {.tag = V_CHAR_TAG, .ch = A};
 memoryManager.constTermsHeap[memoryManager.constTermOffset++] = {.tag = V_CHAR_TAG, .ch = A};
 memoryManager.constTermsHeap[memoryManager.constTermOffset++] = {.tag = V_CHAR_TAG, .ch = A};
 memoryManager.constTermsHeap[memoryManager.constTermOffset++] = {.tag = V_CHAR_TAG, .ch = A};
 memoryManager.constTermsHeap[memoryManager.constTermOffset++] = {.tag = V_CHAR_TAG, .ch = A};
 memoryManager.constTermsHeap[memoryManager.constTermOffset++] = {.tag = V_FLOAT_NUM_TAG, .floatNum = 12.000000};
 memoryManager.constTermsHeap[memoryManager.constTermOffset++] = {.tag = V_IDENT_TAG, .str = d};
 memoryManager.constTermsHeap[memoryManager.constTermOffset++] = {.tag = V_IDENT_TAG, .str = hell};
 memoryManager.constTermsHeap[memoryManager.constTermOffset++] = {.tag = V_IDENT_TAG, .str = asdasd};
 memoryManager.constTermsHeap[memoryManager.constTermOffset++] = {.tag = V_IDENT_TAG, .str = asdasd};
 memoryManager.constTermsHeap[memoryManager.constTermOffset++] = {.tag = V_IDENT_TAG, .str = a};
 memoryManager.constTermsHeap[memoryManager.constTermOffset++] = {.tag = V_INT_NUM_TAG, .intNum = 111};

 return 0;
}
