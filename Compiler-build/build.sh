#!/bin/bash

red='\e[0;31m'
green='\e[0;32m'
NC='\e[0m' # No Color

refalSource=${1}
sourceBaseName=${refalSource%.*}

function Fail
{
	echo -e "${red}[FAIL]: ${refalSource}${NC}"
	exit 1;
}

go install -compiler gccgo ../Compiler/src/refalc/refalc.go

if [ "$?" != 0 ] ; then
	Fail "Can't build compiler"
fi
	
refalc --ptree ${sourceBaseName}.ref

if [ "$?" != 0 ] ; then
	Fail "Can't compile ${refalSource}"
fi

#cat ${sourceBaseName}.ptree
#cat ${sourceBaseName}.c
