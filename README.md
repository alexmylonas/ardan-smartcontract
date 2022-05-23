# Smart Contract

This is a work in progress to understand how to write, debug, and maintain Ethereum smart contracts. This repository is expected to have several apps that help to discover and understand the Ethereum and smart contract development. Please look at the makefile for more details and help.

The following will guide you in deploying the basic smart contract to a local Ethereum service and execute its API.

## Installing Tooling

You must install Ethereum and Solidity into your development environment. If you are using MacOS or Linux, the makefile has brew commands to help you with this. If you can't use brew, there are links to documents in the makefile.

makefile
```
# https://geth.ethereum.org/docs/install-and-build/installing-geth
install-geth:
	brew update
	brew tap ethereum/ethereum
	brew install ethereum

# https://docs.soliditylang.org/en/v0.8.11/installing-solidity.html
install-solc:
	brew update
	brew tap ethereum/ethereum
	brew install solidity
```

## Running Ethereum

Ethereum is configured to run in a developer mode that will allow for the deployment and execution of smart contracts. Use the following `make` command to start the Ethereum service. Make sure you start Ethereum in a dedicated terminal session.

```
$ make geth-up
```

To stop the Ethereum service, use the down command.

```
$ make geth-down
```

## Building Smart Contract

The smart contract can be found in `app/basic/contracts/src/store.sol` and is a smart contract written in the Ethereum programming language called Solidity. It will be helpful to install the Solidity extension by Juan Blanco. To compile and build the smart contract, run the following make command.

```
$ make basic-build
```

That will produce an `abi`, `bin`, and `go` file inside the `app/basic/contracts/store` directory. Any time code for the smart contract is changed, it must be built again for the Go applications.

## Deploying Smart Contract

With the Ethereum service up and running, run the following make command to deploy the smart contract.

```
$ make basic-deploy
```

If this is successful, the following output should result.

```
go run app/basic/cmd/deploy/main.go
fromAddress: 0x6327A38415C53FFb36c11db55Ea74cc9cB4976Fd
transaction: &{0x6327A38415C53FFb36c11db55Ea74cc9cB4976Fd 2 0x101065e90 0 897051399 <nil> <nil> 300000 context.Background false}
tx sent          : 0x515e9a3f97b9764267d6d13eb3a10f483b506abc9e1d5cbfe6bf8913590fcf36
tx gas price     : 0.000000000897051399
tx cost          : 0.000269115419700000
tx gas allowed   : 300000
tx gas used      : 299056
tx status        : 1
Contract Address : 0x87A061ED19dcA76EC5B01643b054f5eae2730a85
balance before   : 115792089237316195423570985008687907853269984665640564039457.583671928961867053
balance after    : 115792089237316195423570985008687907853269984665640564039457.583469638337138349
```

Each time you deploy the contract, you will get a new contract ID. Every code change needs to be mined into a new block and therefore the API moves. Be aware what version of the API you are using.

## Executing Smart Contract API

To validate everything is working, run the follow make command.

```
$ make basic-write
```

When you run this command, attempt to update the contracts map with the key `name` and value `brianna`. You should see the following output if everything is working correctly.

```
go run app/basic/cmd/write/main.go
address: 0x6327A38415C53FFb36c11db55Ea74cc9cB4976Fd
contractID: 0x1da2185Ff21DE4291E8690407A4de5459E09de16
version: 1.1
tx sent        : 0x9e54198a243f99f6f433284b3ac0fe78ff8ad213e6295c9fa8a81a31fe403496
tx gas price   : 0.000000000000000001
tx gas allowed : 250000
tx value       : 0.000000000000000000
tx max cost    : 0.000000000000250000  // gas * gasPrice + value
tx gas used    : 43695
tx act cost    : 0.000000000000043695  // gasUsed * gasPrice + value
tx status      : 1
================ LOGS =================
================ LOGS =================
balance before : 35779.000000000000000000
balance after  : 35781.000000000000000000
balance diff   : -2.000000000000000000
```

To see if the value was written to the map for that key, run the follow make command.

```
$ make basic-read
```

You should see the following output if everything is working correctly.

```
go run app/basic/cmd/read/main.go
address: 0x6327A38415C53FFb36c11db55Ea74cc9cB4976Fd
contractID: 0x1da2185Ff21DE4291E8690407A4de5459E09de16
version: 1.1
value: brianna
```

## What's Next

Check out the `simplecoin` app which is a more complext smart contract. This is a work in progress.
The steps are similar: `scoin-build`, `scoin-deploy`, `scoin-transfer` and `scoin-trancheck`.
