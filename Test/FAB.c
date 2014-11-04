// file:/home/rustam/Diploma/Test/FAB.ref

l_term* Go(vec_header* vecData) 
{
 struct v_term* data = vecData.data;
 uint32_t length = vecData.size;

} // Go

l_term* Pal(vec_header* vecData) 
{
 struct v_term* data = vecData.data;
 uint32_t length = vecData.size;
 int start_0 = follow_0 = 0;
 if (data[start_0]->tag == V_TERM_SYMBOL_TAG) follow_0++;
 	else /*Откат*/;
//--------------------------------------
   if (follow_1 >= length) /*Откат*/;
   int start_2 = follow_2 = follow_1;
   if (data[start_2]->tag == V_TERM_SYMBOL_TAG) follow_2++;
   	else /*Откат*/;
//--------------------------------------

 int start_0 = follow_0 = 0;
 if (data[start_0]->tag == V_TERM_SYMBOL_TAG) follow_0++;
 	else /*Откат*/;
//--------------------------------------



} // Pal

