// file:test.ref

l_term* Go(vec_header* vecData) 
{
 struct v_term* data = vecData.data;
 uint32_t length = vecData.size;
 uint32_t ok = 0;
 ok = 0;
 ok = 1;

 if (ok == 1)
%!s(MISSING){%!s(MISSING)}
 ok = 0;
  if (follow_0 >= length) /*Откат*/;
  int start_1 = follow_1 = follow_0;
  if (data[start_1]->tag == V_TERM_SYMBOL_TAG) follow_1++;
  	else /*Откат*/;
//--------------------------------------

 if (ok == 1)
%!s(MISSING){%!s(MISSING)}
} // Go

l_term* Func(vec_header* vecData) 
{
 struct v_term* data = vecData.data;
 uint32_t length = vecData.size;
 uint32_t ok = 0;
 ok = 0;
 int start_0 = follow_0 = 0;
 if (data[start_0]->tag == V_TERM_SYMBOL_TAG) follow_0++;
 	else /*Откат*/;
//--------------------------------------
  if (follow_0 >= length) /*Откат*/;
  int start_1 = follow_1 = follow_0;
  if (data[start_1]->tag == V_TERM_SYMBOL_TAG) follow_1++;
  	else /*Откат*/;
//--------------------------------------
      if (follow_4 >= length) /*Откат*/;
      int start_5 = follow_5 = follow_4;
      if (data[start_5]->tag == V_TERM_SYMBOL_TAG) follow_5++;
      	else /*Откат*/;
//--------------------------------------

 if (ok == 1)
%!s(MISSING){%!s(MISSING)}
 ok = 0;
 int start_0 = follow_0 = 0;
 if (data[start_0]->tag == V_TERM_SYMBOL_TAG) follow_0++;
 	else /*Откат*/;
//--------------------------------------
  if (follow_0 >= length) /*Откат*/;
  int start_1 = follow_1 = follow_0;
  if (data[start_1]->tag == V_TERM_SYMBOL_TAG) follow_1++;
  	else /*Откат*/;
//--------------------------------------

 if (ok == 1)
%!s(MISSING){%!s(MISSING)}
 ok = 0;
 ok = 1;

 if (ok == 1)
%!s(MISSING){%!s(MISSING)}
} // Func

l_term* Func2(vec_header* vecData) 
{
 struct v_term* data = vecData.data;
 uint32_t length = vecData.size;
 uint32_t ok = 0;
 ok = 0;
 int start_0 = follow_0 = 0;
 if (data[start_0]->tag == V_TERM_SYMBOL_TAG) follow_0++;
 	else /*Откат*/;
//--------------------------------------
  if (follow_0 >= length) /*Откат*/;
  int start_1 = follow_1 = follow_0;
  if (data[start_1]->tag == V_TERM_SYMBOL_TAG) follow_1++;
  	else /*Откат*/;
//--------------------------------------

 if (ok == 1)
%!s(MISSING){%!s(MISSING)}
 ok = 0;
 int start_0 = follow_0 = 0;
 if (data[start_0]->tag == V_TERM_SYMBOL_TAG) follow_0++;
 	else /*Откат*/;
//--------------------------------------

 if (ok == 1)
%!s(MISSING){%!s(MISSING)}
} // Func2

