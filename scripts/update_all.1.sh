#!/usr/bin/env bash

#giveth data --update calc-givback --round $1
#sleep 1
giveth data --update eligible --round $1
sleep 1
giveth data --update not-eligible --round $1
sleep 1
giveth data --update purple-verified --round $1
sleep 1
giveth data --update purple-list --round $1
sleep 1

giveth summarize --round $1

cat data/purpleList/purpleList.json | jq -S >y
mv y data/purpleList/purpleList.json

rm -f x y
