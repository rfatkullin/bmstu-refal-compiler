#include "xmlparse.h"

#include <stdlib.h>
#include "decl.h"
#include "macros.h"
#include "fileio.h"
#include "ifunc.h"

#define CLEAN			0
#define NEED_TO_CLOSE		1
#define NEXT_LINE		2

static FILE * fp;

#ifndef FOR_OS_WINDOWSNT
static inline int 
is_empty (const char * ccp, int i_len ) {

  for (; i_len !=0; ccp++, i_len--) {
    if (!isspace(*ccp)) return 0;
  }
  return 1;
}
#endif

void startElement(void *userData, const char *name, const char **atts)
{
  int *tag = userData;

  if (*tag == NEED_TO_CLOSE) {
    fprintf(fp, "</>)");
  }
  /* fprintf (fp, "((#%s )", name); */
  fprintf (fp, "((%s )", name);

  /*
  if (* atts != NULL) {
    char ** cpp = atts;

    for (cpp = atts; * cpp != NULL; cpp += 2) {
      printf ("Name: \'%s\', Value: \'%s\'\n", * cpp, * (cpp + 1) );
    }
  }
  */

  *tag = CLEAN;
}

void endElement(void *userData, const char *name)
{
  int *tag = userData;
  if (*tag == NEED_TO_CLOSE) {
    fprintf(fp, "</>)");
  }

  *tag = CLEAN;
  fprintf(fp, ")");
}

void charData(void *userData, const char * cp_s, int i_len)
{
  int *tag = userData;

#ifndef FOR_OS_WINDOWSNT
  if (is_empty (cp_s, i_len)) return;
#else
  { 
    char * cp_tmp;
    int i_len_tmp;
    for (cp_tmp = cp_s, i_len_tmp = i_len; i_len_tmp !=0; cp_tmp++, i_len_tmp--) {
      if (!isspace(*cp_tmp)) goto l_char_data_is_not_empty;
    }
    return;
  }
l_char_data_is_not_empty:
#endif

  if (i_len == 1 && *cp_s == '\n') { 
  }else{
  	char ca_w [10];
    int i = 0;

    if (*tag != NEED_TO_CLOSE) {
      ca_w [i++] = '('; 
      if (* tag != NEXT_LINE) {
		ca_w [i++] = '('; ca_w [i++] = '0'; ca_w [i++] = ' '; ca_w [i++] = ')'; 
		ca_w [i++] = '<'; ca_w [i++] = '/'; ca_w [i++] = '>';
      }
    }
    fwrite (ca_w, sizeof (char), i, fp);
    fwrite (cp_s, sizeof (char), i_len, fp);

    *tag = NEED_TO_CLOSE;
  }
}

int ri_xml2ref (FILE * fp_in, FILE * fp_out) {

#ifdef FOR_OS_WINDOWSNT 
  char buf[BUFSIZ];

  XML_Parser parser = XML_ParserCreate(NULL);
  int done;
  int tag = 0;

  fp = fp_out;

  XML_SetUserData(parser, &tag);
  XML_SetElementHandler(parser, startElement, endElement);
  XML_SetCharacterDataHandler(parser, charData);
  do {
    size_t length = fread(buf, 1, sizeof(buf), fp_in);
    done = length < sizeof(buf);
    if (!XML_Parse(parser, buf, length, done)) {
      fprintf(stderr,
	      "%s at line %d\n",
	      XML_ErrorString(XML_GetErrorCode(parser)),
	      XML_GetCurrentLineNumber(parser));
      return 1;
    }
  } while (!done);
  XML_ParserFree(parser);
#endif 

  return 0;
}
