#ifndef __BUILTINS_H__
#define __BUILTINS_H__

#include "l_term.h"
#include "memory_manager.h"
#include "func_call_t.h"

struct func_result_t prout(int entryPoint, struct env_t* env, struct field_view_t* filedOfView);
struct l_term* card(struct l_term* expr);

#endif
