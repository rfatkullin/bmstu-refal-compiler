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

../Compiler-build/build.sh ../Compiler-build/${refalSource}

if [ "$?" != 0 ] ; then
	Fail "Compiler-build error"
fi

cp ../Compiler-build/${sourceBaseName}.c ../Project/main.c

make

if [ "$?" != 0 ] ; then
	Fail "Can't build project!"
fi

./Project

