#!/usr/bin/env bash

echo
echo "------------------------------------------------------------------------------------------------------------"
echo "Running process_round.sh"
echo "------------------------------------------------------------------------------------------------------------"

if [ -z "$1" ]
then
    echo "Usage ./process_round <round>"
    exit
fi

export TB_DEFAULT_FMT=json
export ROUND=$1
export DIR=Round_$1

#------------------------------------
# Start with a clean slate
#------------------------------------
mkdir -p $DIR

#------------------------------------
# Start with a clean slate
#------------------------------------
./scripts/clean.sh $ROUND
#exit

#------------------------------------
# Extract the transaction ids and network from the round's data
#------------------------------------
giveth data eligible -r $ROUND | jq '.data[] | "\(.txHash) \(.network)"' | sed 's/xDAI/gnosis/' | sed 's/\"//g' | tr ' ' ',' >tmp
#exit

#------------------------------------
# Separate the data by chain
#------------------------------------
GNOSIS_ELIGIBLE=$DIR"/results/gnosis/txhashes_"$ROUND".csv"
MAINNET_ELIGIBLE=$DIR"/results/mainnet/txhashes_"$ROUND".csv"
echo $GNOSIS_ELIGIBLE
echo $MAINNET_ELIGIBLE

cat tmp | grep gnosis  | cut -f1 -d, >$GNOSIS_ELIGIBLE
cat tmp | grep mainnet | cut -f1 -d, >$MAINNET_ELIGIBLE
rm -f tmp
#exit

#------------------------------------
# Extract, from the node, the actual transactional data
#------------------------------------
GNOSIS_TXS=$DIR"/results/gnosis/txs_"$ROUND".csv"
MAINNET_TXS=$DIR"/results/mainnet/txs_"$ROUND".csv"
echo $GNOSIS_TXS
echo $MAINNET_TXS

chifra transactions --no_header --file $GNOSIS_ELIGIBLE  --fmt csv --chain gnosis  | sed 's/\"//g' | tee $GNOSIS_TXS
chifra transactions --no_header --file $MAINNET_ELIGIBLE --fmt csv --chain mainnet | sed 's/\"//g' | tee $MAINNET_TXS
#exit

#------------------------------------
# Get some stats, so we can see what's going on
#------------------------------------
GNOSIS_DONOR_COUNTS=$DIR"/results/gnosis/donor_counts_"$ROUND".csv"
MAINNET_DONOR_COUNTS=$DIR"/results/mainnet/donor_counts_"$ROUND".csv"
echo $GNOSIS_DONOR_COUNTS
echo $MAINNET_DONOR_COUNTS

cat $GNOSIS_TXS  | cut -f5 -d, | sort | uniq -c | sort -n -r | tee $GNOSIS_DONOR_COUNTS
cat $MAINNET_TXS | cut -f5 -d, | sort | uniq -c | sort -n -r | tee $MAINNET_DONOR_COUNTS
#exit

#------------------------------------
# Extract just the from addresses in all those transactions
#------------------------------------
GNOSIS_DONOR=$DIR"/results/gnosis/donors_"$ROUND".csv"
MAINNET_DONOR=$DIR"/results/mainnet/donors_"$ROUND".csv"
echo $GNOSIS_DONOR
echo $MAINNET_DONOR

cat $GNOSIS_TXS  | cut -f5 -d, | sort -u | tee $GNOSIS_DONOR
cat $MAINNET_TXS | cut -f5 -d, | sort -u | tee $MAINNET_DONOR
#exit

#------------------------------------
# Create shell scripts to extract the neighbors lists
#------------------------------------
GNOSIS_SCRIPT="scripts/get_neighbors_gnosis_"$ROUND".sh"
MAINNET_SCRIPT="scripts/get_neighbors_mainnet_"$ROUND".sh"
echo $GNOSIS_SCRIPT
echo $MAINNET_SCRIPT

cat $GNOSIS_DONOR  | sed 's/^/.\/scripts\/get_neighbors.sh gnosis '$ROUND' /'  | tee $GNOSIS_SCRIPT
cat $MAINNET_DONOR | sed 's/^/.\/scripts\/get_neighbors.sh mainnet '$ROUND' /'  | tee $MAINNET_SCRIPT
#exit

#------------------------------------
# Create the neighbors files
#------------------------------------
FIRST=13858106 source $GNOSIS_SCRIPT
FIRST=13858106 source $MAINNET_SCRIPT

head $DIR/results/gnosis/*counts* | tee $DIR/stats/gnosis.txt
find $DIR/results/gnosis -name "*.csv" -exec wc {} ';' | grep -v counts | tee -a $DIR/stats/gnosis.txt

head $DIR/results/mainnet/*counts* | tee $DIR/stats/mainnet.txt
find $DIR/results/mainnet -name "*.csv" -exec wc {} ';' | grep -v counts | tee -a $DIR/stats/mainnet.txt

wc $DIR/neighbors/*/*/*
