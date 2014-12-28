#!/bin/bash
#Run from build directory

TestsDir="../Tests"
TmpRefSourceFile="source.ref"
TmpCSourceFile="source.c"
RealOutputFile="output.txt"

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

for sourceFile in `ls ${TestsDir}/*.ref`
do	
	cp ${sourceFile} ${TmpRefSourceFile}
	
	#Компилируем рефал программу
	refalc ${TmpRefSourceFile} 1>/dev/null 
	AssertSuccess "Can't compile refal source ${sourceFile}"		
	
	#Собираем весь проект - линкуем сгенерированный файл с библиотекой исполнения.
	cp ${TmpCSourceFile} ../Project/main.c
	make 1>/dev/null
	AssertSuccess "Can't build project!"
	
	#Запускаем испольняемый файл.
	./Project > ${RealOutputFile}
	AssertSuccess "Bad execuatable file!"
	
	#Проверям ожидаемое с полученным
	cmp -s ${RealOutputFile} ${sourceFile%.*}.out 
	AssertSuccess "Check by command: vim -d ${RealOutputFile} ${sourceFile%.*}.out"		
	
	echo -e "${green}[OK]: ${RealOutputFile} ${sourceFile%.*}.out ${NC}"
done