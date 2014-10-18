#include "environment.h"

int addEnvVar(struct environment* env, char* name, struct v_term* value)
{
	return 1;
}

struct v_term* getEnvVar(struct environment* env, char* name)
{
	return 0;
}

void removeEnvVar(struct environment* env, char* name)
{
}

void clearEnv(struct environment* env)
{
}
