// file:FAB.ref

l_term* Go(vec_header* vecData) 
{
 struct v_term* data = vecData.data;
 uint32_t length = vecData.size;
 uint32_t ok = 0;
 ok = 0;
 ok = 1;

} // Go

l_term* Pal(vec_header* vecData) 
{
 struct v_term* data = vecData.data;
 uint32_t length = vecData.size;
 uint32_t ok = 0;
 ok = 0;
 int start_0 = follow_0 = 0;
 if (data[start_0]->tag == V_TERM_SYMBOL_TAG) follow_0++;
 	else /*Откат*/;
//--------------------------------------
   if (follow_1 >= length) /*Откат*/;
   int start_2 = follow_2 = follow_1;
   if (data[start_2]->tag == V_TERM_SYMBOL_TAG) follow_2++;
   	else /*Откат*/;
//--------------------------------------

 ok = 0;
 int start_0 = follow_0 = 0;
 if (data[start_0]->tag == V_TERM_SYMBOL_TAG) follow_0++;
 	else /*Откат*/;
//--------------------------------------

 ok = 0;
 ok = 1;

 ok = 0;

} // Pal

