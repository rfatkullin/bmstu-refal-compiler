#!/bin/bash

rm a.out
gcc ../main.c ../memory_manager.c ../segment_tree.c
./a.out