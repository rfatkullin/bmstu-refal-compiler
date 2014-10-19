#!/bin/bash

go install ../Compiler/src/refalc/refalc.go
refalc --ptree /home/rustam/Diploma/Test/FAB.ref
cat FAB.c
