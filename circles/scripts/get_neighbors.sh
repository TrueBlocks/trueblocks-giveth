#!/usr/bin/env bash

echo
echo "------------------------------------------------------------------------------------------------------------"
echo "Running steps.sh"
echo "------------------------------------------------------------------------------------------------------------"

if [ -z "$2" ]
then
    echo "Usage ./get_neighbors <chain> <address>"
    exit
fi
if [ -z "$FIRST" ]
then
    echo "Usage: \$FIRST is required."
    exit
fi

export CHAIN=$1
export ROUND=$2
export ADDR=$3
export DIR=Round_$2

export DEST1=$DIR"/neighbors/"$CHAIN"/raw/"$ADDR".txt"
export DEST2=$DIR"/neighbors/"$CHAIN"/raw_senders/"$ADDR".txt"
export DEST3=$DIR"/neighbors/"$CHAIN"/counts/"$ADDR".txt"
export DEST4=$DIR"/neighbors/"$CHAIN"/senders/"$ADDR".txt"

echo $DEST1
echo $DEST2
echo $DEST3
echo $DEST4
#exit

chifra export --cache --neighbors --first_block $FIRST $ADDR | tee $DEST1
cat $DEST1 | grep from | grep -v $ADDR | tee $DEST2
cat $DEST2 | cut -f3 | sort | uniq -c | sort -n | tee $DEST3
cat $DEST2 | cut -f3 | sort -u | tee $DEST4
