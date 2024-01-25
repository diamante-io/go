# Get state from diamcircle-core DB, colums match CSV printer
# FETCH_COUNT is there for circleci to use cursor-based method of getting rows (less RAM usage):
# https://dba.stackexchange.com/a/101510

echo "Fetching accounts from diamcircle-core DB..."
psql -d core -t -A -F"," --variable="FETCH_COUNT=10000" -c "select accountid, balance, seqnum, numsubentries, inflationdest, homedomain, thresholds, flags, COALESCE(extension, 'AAAAAA=='), signers, ledgerext from accounts" > accounts_core.csv
rm accounts_core_sorted.csv || true # Remove if exist in case original files are rebuilt

echo "Fetching accountdata from diamcircle-core DB..."
psql -d core -t -A -F"," --variable="FETCH_COUNT=10000" -c "select accountid, dataname, datavalue, COALESCE(extension, 'AAAAAA=='), ledgerext from accountdata" > accountdata_core.csv
rm accountdata_core_sorted.csv || true # Remove if exist in case original files are rebuilt

echo "Fetching offers from diamcircle-core DB..."
psql -d core -t -A -F"," --variable="FETCH_COUNT=10000" -c "select sellerid, offerid, sellingasset, buyingasset, amount, pricen, priced, flags, COALESCE(extension, 'AAAAAA=='), ledgerext from offers" > offers_core.csv
rm offers_core_sorted.csv || true # Remove if exist in case original files are rebuilt

echo "Fetching trustlines from diamcircle-core DB..."
psql -d core -t -A -F"," --variable="FETCH_COUNT=10000" -c "select ledgerentry from trustlines" > trustlines_core.csv
rm trustlines_core_sorted.csv || true # Remove if exist in case original files are rebuilt

echo "Fetching claimable balances from diamcircle-core DB..."
psql -d core -t -A -F"," --variable="FETCH_COUNT=10000" -c "select balanceid, ledgerentry from claimablebalance" > claimablebalances_core.csv
rm claimablebalances_core_sorted.csv || true # Remove if exist in case original files are rebuilt

echo "Fetching liquidity pools from diamcircle-core DB..."
psql -d core -t -A -F"," --variable="FETCH_COUNT=10000" -c "select ledgerentry from liquiditypool" > pools_core.csv
rm pools_core_sorted.csv || true # Remove if exist in case original files are rebuilt