# Bococoin
Boco Coin is a cryptocurrency and worldwide payment system. ﻿It is the first decentralized digital currency, as the system works without a central bank or single administrator.

# Build 
- Install [Git](https://git-scm.com/download)
- Install [GoLang](https://golang.org/dl/) 1.14.1 or higher 
- Clone Bococoin core repository 

```
md /home/user/bococoin
cd /home/user/bococoin
git clone https://github.com/Bococoin/core
```
- Build sources
```
cd ./core
make all
```
# Run node
```
$ cd GOPATH/bin
$ chmod +x ./bocod
$ chmod +x ./bococli
$ ./bocod init <node_moniker> --chain-id=boco-02
$ wget -O genesis.json -P /home/user/.bococoin/config https://github.com/Bococoin/core/releases/download/1.0/genesis.json
$ ./bocod start
```

# Run Light Client rest-server
```
$ cd GOPATH/bin
$ ./bococli config output json
$ ./bococli config indent true
$ ./bococli config trust-node true
$ ./bococli config chain-id boco-02
$ ./bococli config node "https://rpc.bococoin.com:443" #remove this line for connect to local node
$ ./bococli config keyring-backend "file"
$ ./bococli rest-server
```