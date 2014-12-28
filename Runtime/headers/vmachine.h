#ifndef __V_MACHINE_H__
#define __V_MACHINE_H__

#include "lterm.h"
#include "func_call.h"

void mainLoop(struct func_result_t (*firstFuncPtr)(int entryPoint, struct env_t* env, struct field_view_t* fieldOfView));
//struct l_term* createLTermFuncCall(const char* funcName, struct l_term* prev, struct l_term* (*func)(void* args), struct l_term* args, void* stackArgs);

#endif
