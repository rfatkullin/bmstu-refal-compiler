#ifndef __ENVIRONMENT_H__
#define __ENVIRONMENT_H__

#include "v_term.h"

struct environment
{
	struct environment* parentEnv;
	int hashTable[];
};

int addEnvVar(struct environment* env, char* name, struct v_term* value);
struct v_term* getEnvVar(struct environment* env, char* name);
void removeEnvVar(struct environment* env, char* name);
void clearEnv(struct environment* env);

#endif