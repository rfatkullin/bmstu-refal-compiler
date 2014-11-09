#include "version.h"

#ifdef FOR_OS_DOS
char _refal_build_version [] = "Feb 01 2000";
#else
char _refal_build_version [] = __DATE__;
#endif
