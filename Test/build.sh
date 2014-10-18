#!/bin/bash

go install ../Compiler/src/refalc/refalc.go
refalc /home/rustam/Diploma/Test/FAB.ref
cat FAB.c
