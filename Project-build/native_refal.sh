#!/bin/bash

TmpRefSourceFile="source.ref"
TmpRslSourceFile="source.rsl"

cp ${1} ${TmpRefSourceFile}
refc ${TmpRefSourceFile}
refgo ${TmpRslSourceFile}

