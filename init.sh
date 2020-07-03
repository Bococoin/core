./bocod init newNode1 --chain-id=boco-test-01

./bococli config output json
./bococli config indent true
./bococli config trust-node true
./bococli config chain-id boco-test-01

./bococli keys add premine

./bocod add-genesis-account $(./bococli keys show premine -a) 20000000000000ubcc

./bocod gentx --name premine --amount 5000000000000ubcc  #DefaultMinValidatorSelfDelegation

echo "Collecting genesis txs..."
./bocod collect-gentxs

echo "Validating genesis file..."
./bocod validate-genesis