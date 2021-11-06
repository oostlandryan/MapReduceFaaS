#!/bin/bash

go mod download
go run mapreduce_demo.go -word=you -m=25 -r=30 -createindex=true