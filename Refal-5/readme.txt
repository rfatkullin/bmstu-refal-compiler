Here are sources of Refal-5 VERSION-PZ.
Version 28 October 2004.

There are a number of makefiles: 
 makefile.win, makefile.lin .

************************** TO COMPILE. **********************************
Windows-NT, Windows 2000, Windows XP:
    -1- Windows-NT> copy makefile.win makefile
    -2- Windows-NT> nmake

Linux/Intel:
    -1- Linux> cp makefile.lin to makefile
    -2- Linux> make

FreeBSD/Intel:
    -1- Linux> cp makefile.lin to makefile
    -2- Linux> gmake

*************************************************************************

---Content---

========== *.c =======
arithm.c
bif.c
bif_lex.c
dcom.c
dimpl.c
dtrace.c
freeze.c
func1.c
func2.c
lex.c
load.c
macros.c
parser.c
pass2.c
rc.c
rc5.c
rcaux.c
rcleft.c
rcopt.c
refaux.c
refgo.c
refio.c
ri.c
rti.c
see1.c
see2.c
see2a.c
sem.c
sysfun.c
test.c
test1.c
trace.c
version.c
vyvod.c
xml2ref.c
xxx.c
xxx1.c

========= *.h ==============
arithm.h
bif_lex.h
bif_lex.bak
cdecl.h
cfunc.h
ddecl.h
decl.h
dmacro.h
fileio.h
freeze.h
ifunc.h
junk.h
ldecl.h
macros.h
memory.h
rasl.h
tfunc.h
version.h
xmlparse.h

========== *.dll =============
              --- Dynamic XML-libraries for Windows NT.
xmlparse.dll
xmltok.dll

========== *.lib =============
              --- XML-libraries for Windows NT.
xmlparse.lib
xmltok.lib

========== linux/libxml* ===========
              --- XML-libraries for LINUX.
libxml.so.1.8.7
libxmlparse.a
libxmlparse.so.1.1
libxmltok.a
libxmltok.so.1.1

========== linux/libxml* ===========
              --- Batch modules for LINUX.
LinuxXMLLib       -- To link libraries.
LinuxXMLLib.del   -- To unlink libraries.
LinuxXMLLib.root  -- To open libraries.

========== *.lnk =============
refc.lnk
refgo.lnk
reftr.lnk

========== *.ref =============
reflib.ref    -- Refal5 library.
e.ref         -- Evaluator of Refal expressions.
mbprep.ref    -- Multibracket preprocessor.
test.ref      -- A test of the Refal-5 system.

========== *.txt =============
readme.txt    -- This file.
bf.txt        -- Guide to Writing a Built in Function for Refal-5.
rdhelp.txt    -- Refal-5 tracer commands.
news.txt      -- News of the version.
readme        -- Help how to install and to use Refal-5.
install.txt   -- Help how to install Refal-5.
copyright.txt -- The REFAL-5 VERSION-PZ Copyright.

***************** NOTES **************************************
 - This version is not supported under MS-DOS.
**************************************************************
