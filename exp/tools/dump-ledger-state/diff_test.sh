ENTRIES=(accounts accountdata offers trustlines claimablebalances pools)

echo "Sorting dump-ledger-state output files..."
for i in "${ENTRIES[@]}"
do
  if test -f "${i}_sorted.csv"; then
    echo "Skipping, ${i}_sorted.csv exists (remove if out of date to sort again)"
    continue
  fi
  wc -l ${i}.csv
  sort -S 500M -o ${i}_sorted.csv ${i}.csv
done

echo "Sorting diamcircle-core output files..."
for i in "${ENTRIES[@]}"
do
  if test -f "${i}_core_sorted.csv"; then
    echo "Skipping, ${i}_core_sorted.csv exists (remove if out of date to sort again)"
    continue
  fi
  wc -l ${i}_core.csv
  sort -S 500M -o ${i}_core_sorted.csv ${i}_core.csv
done

echo "Checking diffs..."
for type in "${ENTRIES[@]}"
do
  diff -q ${type}_core_sorted.csv ${type}_sorted.csv
  if [ "$?" -ne "0" ]
  then
    echo "ERROR: $type does NOT match";
    exit -1
  else
    echo "$type OK";
  fi
done
