#!/usr/bin/env bash

if [ -z "$1" ]
then
    echo "Usage ./clean.sh <round>"
    exit
fi

DIR=Round_$1

rm -fR $DIR/results
rm -fR $DIR/neighbors
rm -fR $DIR/names
rm -fR $DIR/stats

mkdir -p $DIR/results/gnosis
mkdir -p $DIR/results/mainnet

mkdir -p $DIR/neighbors/gnosis/counts
mkdir -p $DIR/neighbors/gnosis/senders
mkdir -p $DIR/neighbors/gnosis/raw
mkdir -p $DIR/neighbors/gnosis/raw_senders
mkdir -p $DIR/neighbors/mainnet/counts
mkdir -p $DIR/neighbors/mainnet/senders
mkdir -p $DIR/neighbors/mainnet/raw
mkdir -p $DIR/neighbors/mainnet/raw_senders

mkdir -p $DIR/names/gnosis
mkdir -p $DIR/names/mainnet

mkdir -p $DIR/stats/gnosis
mkdir -p $DIR/stats/mainnet
