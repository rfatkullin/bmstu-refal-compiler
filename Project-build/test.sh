#!/bin/bash
#Run from build directory

TestsDir="../Tests"
TmpRefSourceFile="source.ref"
TmpCSourceFile="source.c"

red='\e[0;31m'
green='\e[0;32m'
NC='\e[0m' # No Color

function AssertSuccess
{
	if [ "$?" != 0 ] ; then
		echo -e "${red}[FAIL] $1${NC}"
		exit 1;
	fi
}

#Собираем библиотеку рантайма.
cd ../Runtime-build 
make 1>/dev/null 
AssertSuccess "Runtime-build error" 
cd - 1>/dev/null

#Компилируем компилятор! В итоге получаем исполняемый файл refalc, который кладется в папку, прописанную в переменной PATH.
go install -compiler gccgo ../Compiler/src/refalc/refalc.go 
AssertSuccess "Can't build compiler"

cp ${1} ${TmpRefSourceFile}

#Компилируем рефал программу
refalc ${TmpRefSourceFile} 1>/dev/null 
AssertSuccess "Can't compile refal source ${1}"

#Собираем весь проект - линкуем сгенерированный файл с библиотекой исполнения.
cp ${TmpCSourceFile} ../Project/main.c
make 1>/dev/null
AssertSuccess "Can't build project!"

#Запускаем испольняемый файл.
./Project


