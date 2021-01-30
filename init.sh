./bocod init validator1 --chain-id=boco-02

./bococli config chain-id boco-02
./bococli config output json
./bococli config indent true
./bococli config trust-node true
./bococli config keyring-backend file


./bococli keys add premine

./bocod add-genesis-account $(./bococli keys show premine -a) 25000000000000ubcc

./bocod gentx --name premine --amount 5000000000000ubcc  #DefaultMinValidatorSelfDelegation

echo "Collecting genesis txs..."
./bocod collect-gentxs

echo "Validating genesis file..."
./bocod validate-genesis