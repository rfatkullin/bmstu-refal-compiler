#!/bin/bash

red='\e[0;31m'
green='\e[0;32m'
NC='\e[0m' # No Color

function Fail
{
	echo -e "${red}[FAIL]: $1${NC}"
	exit 1;
}

go install ../Compiler/src/refalc/refalc.go

if [ "$?" != 0 ] ; then
	Fail "Can't build compiler"
fi
	
refalc --ptree $1.ref

if [ "$?" != 0 ] ; then
	Fail "Can't compile $1"
fi

cat $1.ptree
