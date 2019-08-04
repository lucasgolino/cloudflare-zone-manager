#!/bin/bash

for file in */ ; do
    target=${file%/}
    go build -buildmode=plugin -o $target.so $file/$target.go
done