# For Linux
OBJS_REFGO = refgo.o refaux.o bif.o load.o ri.o refio.o func1.o func2.o macros.o sysfun.o freeze.o arithm.o bif_lex.o xxx.o xml2ref.o
OBJS_REFTR = trace.o load.o rcaux.o rti.o bif.o arithm.o refaux.o func1.o func2.o macros.o lex.o dtrace.o \
             dcom.o sysfun.o dimpl.o freeze.o rcleft.o refio.o bif_lex.o xxx.o xml2ref.o
OBJS_REFC = rc5.o parser.o lex.o sem.o pass2.o rc.o rcleft.o rcopt.o rcaux.o vyvod.o bif_lex.o

CC = gcc
# C_FLAGS = -O2 -c -W -Wreturn-type -Wunused -Wshadow -funsigned-char 
# C_FLAGS = -O2 -c -Wall -funsigned-char 
# C_FLAGS = -g3 -c -Wall -funsigned-char 
C_FLAGS = -c -DFOR_OS_LINUX -Wall -funsigned-char 
# LFLAGS = -L./linux/ -lxmlparse -lxmltok
LFLAGS = -L.
# LFLAGS = -L
rm = rm -f 

.c.o:
	$(CC) $(C_FLAGS) -c $<

.c.exe:
	$(CC) $(C_FLAGS) $<

# all: refgo reftr refc see1 see2
all: refgo reftr refc

#refgo: refgo.o refaux.o bif.o load.o ri.o refio.o func1.o func2.o macros.o \
#    sysfun.o freeze.o arithm.o
#	$(CC) $(LFLAGS) version.c `cat refgo.lnk` -o refgo
#	$(CC_PATH)\strip refgo
#	$(CC_PATH)\aout2exe refgo
refgo: $(OBJS_REFGO)
	$(CC) $(LFLAGS) version.c $(OBJS_REFGO) -o refgo


#reftr: trace.o load.o rcaux.o rti.o bif.o arithm.o refaux.o func1.o func2.o \
#   macros.o lex.o dtrace.o dcom.o sysfun.o dimpl.o freeze.o rcleft.o refio.o
#	$(CC) $(LFLAGS) version.c `cat reftr.lnk` -o reftr
#	$(CC_PATH)\strip reftr
#	$(CC_PATH)\aout2exe reftr
reftr: $(OBJS_REFTR)
	$(CC) $(LFLAGS) version.c $(OBJS_REFTR) -o reftr

#refc: rc5.o parser.o lex.o sem.o pass2.o \
#	rc.o rcleft.o rcopt.o rcaux.o vyvod.o
#	$(CC) $(LFLAGS) version.c `cat rc5.lnk` -o refc
#	$(CC_PATH)\strip refc
#	$(CC_PATH)\aout2exe refc
refc: $(OBJS_REFC)
	$(CC) $(LFLAGS) version.c $(OBJS_REFC) -o refc

see1: see1.c
	$(CC) -o see1 see1.c
#	$(CC_PATH)\strip see1
#	$(CC_PATH)\aout2exe see1

see2: see2.c
	$(CC) -o see2 see2.c
#	$(CC_PATH)\strip see2
#	$(CC_PATH)\aout2exe see2

clean: 
	$(rm) *.o
	$(rm) a.out reftr refgo refc see1 see2

