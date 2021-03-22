<p align="center">
<img src="./pic/neo_color_dark_gopher.png" width="300px" alt="logo">
</p>

[NEO](https://neo.org/) builds smart economy and we at [NEO SPCC](https://nspcc.ru/en/) help them with that big challenge. 
In our blog you might find the latest articles [how we run NEOFS public test net](https://medium.com/@neospcc/public-neofs-testnet-launch-18f6315c5ced) 
but it’s not the only thing we’re working on.

## NEO GO
As you know network is composed of nodes. These nodes as of now have several implementations:
- https://github.com/neo-project/neo
- https://github.com/CityOfZion/neo-python
- https://github.com/nspcc-dev/neo-go

This article is about the last one since we’re developing it at NEO SPCC. 
Hope that this article will help you to get an idea of how everything is tied up and being able to start neo-go node,
 write smart contract and deploy it.
 
## What is a node?

<p align="center">
<img src="./pic/node.png" width="300px" alt="node">
</p>

The main goal of the node is to interact with each other (through P2P) and synchronize blocks in the network. 
It also allows user to compile and run smart contracts within the blockchain network. 
Node consists of Client (CLI), Network layer, Consensus, Virtual Machine, Compiler and Blockchain.
 Let’s take a closer look at each of them. 

#### Client
Client (CLI) allows users to run commands from the terminal. These commands can be divided in 4 categories:

- server operations
- smart contract operations
- vm operations
- wallet operations


For example to connect node to the running private network you can use this command:
```
 go run cli/main.go node -p
```
Here you can find more information about Private Network and how to start it. Simply speaking private network -- it’s the network that you can run locally. 
Follow the link if you are interested in more detailed description
[medium article](https://medium.com/@neospcc/neo-privatenet-auto-import-of-a-smart-contract-dbf2b9220ad2). 
Another usage example is to compile smart contract:
```
$ ./bin/neo-go vm 

    _   ____________        __________      _    ____  ___
   / | / / ____/ __ \      / ____/ __ \    | |  / /  |/  /
  /  |/ / __/ / / / /_____/ / __/ / / /____| | / / /|_/ / 
 / /|  / /___/ /_/ /_____/ /_/ / /_/ /_____/ |/ / /  / /  
/_/ |_/_____/\____/      \____/\____/      |___/_/  /_/   


NEO-GO-VM >  
```
Once we run this command we will get an interface to interact with virtual machine. 
To get a list of all supported operation you just use `help`:
```
NEO-GO-VM > help

Commands:
  astack        Show alt stack contents
  break         Place a breakpoint
  clear         clear the screen
  cont          Continue execution of the current loaded script
  estack        Show evaluation stack contents
  exit          Exit the VM prompt
  help          display help
  ip            Show current instruction
  istack        Show invocation stack contents
  loadavm       Load an avm script into the VM
  loadgo        Compile and load a Go file into the VM
  loadhex       Load a hex-encoded script string into the VM
  ops           Dump opcodes of the current loaded program
  run           Execute the current loaded script
  step          Step (n) instruction in the program
  stepinto      Stepinto instruction to take in the debugger
  stepout       Stepout instruction to take in the debugger
  stepover      Stepover instruction to take in the debugger

```
As you can see there are a lot of options to play with it. Let’s take simple smart contract `1-print.go` and compile it:
 
```
package main

import (
	"github.com/nspcc-dev/neo-go/pkg/interop/runtime"
)

func Main() {
	runtime.Log("Hello, world!")
}
``` 
Use command `loadgo` to compile it:
```
NEO-GO-VM > loadgo 1-print.go
READY: loaded 22 instructions
NEO-GO-VM 0 >  
```
And there you can see how many instructions were generated and even if you are interested in opcodes of current program you can dump them:
```
NEO-GO-VM 0 > ops
0        PUSHDATA1    48656c6c6f2c20776f726c6421 ("Hello, world!")    <<
15       SYSCALL      System.Runtime.Log (cfe74796)
20       NOP
21       RET
```
Later we will use this compiled contract in a workshop =).
You can find more information on how to use the CLI [here](https://github.com/nspcc-dev/neo-go/blob/master/docs/cli.md)

#### Network
Network layer is one of the most important parts of the node. In our case we have P2P protocol which allows nodes to communicate with each other and RPC -- which is used for getting some information from node like balance, accounts, current state, etc.
Here is the document where you can find supported [RPC calls](https://github.com/nspcc-dev/neo-go/blob/master/docs/rpc.md).

#### Consensus
Consensus is a mechanism allowing nodes to agree on a specific value (block in case of blockchain). We use our own go-implementation of dBFT algorithm.

#### Compiler
Compiler allows to build byte code, so you can write Smart Contract in your favourite Golang. All the output you saw in this example above was generated by the compiler.

#### Virtual machine
Virtual machine runs compiled byte code. NeoVM is a [stack-based virtual machine](https://docs.neo.org/docs/en-us/basic/technology/neovm.html). It has 2 stacks for performing computation.

#### Blockchain
And what is the Blockchain piece? It’s quite a big one since it contains operations with accepting/validation transactions, 
signing transactions, working with accounts, assets, storing blocks in database (or in cache).

#### Network
There are 3 types of network. 
Private net -- it’s the private one which you can run locally. Testnet and Mainnet where much of the nodes across the world now running. 
NEO has a nice monitor where you can find particular node running in the blockchain network.
[Neo Monitor](http://monitor.cityofzion.io/)

## Workshop. Preparation
In this part we will setup the environment: run private network, connect neo-go node to it and transfer some initial GAS to our basic account
in order to be able to pay for transaction deployment and invocation. Let's start.

#### Requirements
For this workshop you will need Debian 10, Docker, docker-compose, go to be installed:
- [docker](https://docs.docker.com/install/linux/docker-ce/debian/)
- [go](https://golang.org/dl/)

#### Versioning
As with many other Neo projects NeoGo is currently on its way to Neo 3, so there are two main branches there — [master](https://github.com/nspcc-dev/neo-go),
where all Neo 3 development is happening right now and [master-2.x](https://github.com/nspcc-dev/neo-go/tree/master-2.x) for stable Neo 2 implementation. 
This workshop contains basic tutorial notes for Neo 3 version. 
If you want to continue with Neo 2, please, refer to [master-2.x branch](https://github.com/nspcc-dev/neo-go-sc-wrkshp/tree/master-2.x).

#### Step 1
If you already have neo-go or go smart-contracts, please, update go modules in order to be up-to-date with the current interop API changes.
If not, download neo-go and build it (master branch):
```
$ git clone https://github.com/nspcc-dev/neo-go.git
$ cd neo-go
$ make build 
```

#### Step 2
There are 2 ways of running local private network. 
One way is using neo-local private network and other way is with neo-go private network.

#### Running with neo-go private network
```
$ make env_image
$ make env_up
```
Result: running privatenet:
```
=> Bootup environment
Creating network "neo_go_network" with the default driver
Creating volume "docker_volume_chain" with local driver
Creating neo_go_node_four  ... done
Creating neo_go_node_two   ... done
Creating neo_go_node_one   ... done
Creating neo_go_node_three ... done
```

In case you need to shutdown environment you can use:
```
$ make env_down
```

#### Running with neo local private network
```
git clone https://github.com/CityOfZion/neo-local.git
$ cd neo-local
$ git checkout -b 4nodes 0.12
$ make start
```

#### Step 3
Start neo-go node which will connect to previously started privatenet:
```
$ ./bin/neo-go node --privnet
```

Result:
```
2020-12-17T14:51:53.200+0300	INFO	no storage version found! creating genesis block
2020-12-17T14:51:53.203+0300	INFO	starting rpc-server	{"endpoint": ":20331"}
2020-12-17T14:51:53.203+0300	INFO	service is running	{"service": "Prometheus", "endpoint": ":2112"}
2020-12-17T14:51:53.203+0300	INFO	service hasn't started since it's disabled	{"service": "Pprof"}
2020-12-17T14:51:53.203+0300	INFO	node started	{"blockHeight": 0, "headerHeight": 0}

    _   ____________        __________
   / | / / ____/ __ \      / ____/ __ \
  /  |/ / __/ / / / /_____/ / __/ / / /
 / /|  / /___/ /_/ /_____/ /_/ / /_/ /
/_/ |_/_____/\____/      \____/\____/

/NEO-GO:0.91.1-pre-657-gc13d6ecc/

2020-12-17T14:51:53.204+0300	INFO	new peer connected	{"addr": "127.0.0.1:20333", "peerCount": 1}
2020-12-17T14:51:53.206+0300	INFO	started protocol	{"addr": "127.0.0.1:20333", "userAgent": "/NEO-GO:0.91.1-pre-657-gc13d6ecc/", "startHeight": 0, "id": 3172166887}
2020-12-17T14:51:54.204+0300	INFO	blockchain persist completed	{"persistedBlocks": 0, "persistedKeys": 71, "headerHeight": 0, "blockHeight": 0, "took": "765.955µs"}
2020-12-17T14:51:56.204+0300	INFO	new peer connected	{"addr": "127.0.0.1:20336", "peerCount": 2}
2020-12-17T14:51:56.204+0300	INFO	new peer connected	{"addr": "127.0.0.1:20334", "peerCount": 3}
2020-12-17T14:51:56.205+0300	INFO	new peer connected	{"addr": "127.0.0.1:20335", "peerCount": 4}
2020-12-17T14:51:56.205+0300	INFO	new peer connected	{"addr": "127.0.0.1:20333", "peerCount": 5}
2020-12-17T14:51:56.205+0300	INFO	started protocol	{"addr": "127.0.0.1:20336", "userAgent": "/NEO-GO:0.91.1-pre-657-gc13d6ecc/", "startHeight": 0, "id": 90708676}
2020-12-17T14:51:56.206+0300	WARN	peer disconnected	{"addr": "127.0.0.1:20333", "reason": "already connected", "peerCount": 4}
2020-12-17T14:51:56.206+0300	INFO	started protocol	{"addr": "127.0.0.1:20334", "userAgent": "/NEO-GO:0.91.1-pre-657-gc13d6ecc/", "startHeight": 0, "id": 410946741}
2020-12-17T14:51:56.207+0300	INFO	started protocol	{"addr": "127.0.0.1:20335", "userAgent": "/NEO-GO:0.91.1-pre-657-gc13d6ecc/", "startHeight": 0, "id": 4085957952}
2020-12-17T14:52:35.213+0300	INFO	blockchain persist completed	{"persistedBlocks": 1, "persistedKeys": 19, "headerHeight": 1, "blockHeight": 1, "took": "518.786µs"}
2020-12-17T14:52:50.217+0300	INFO	blockchain persist completed	{"persistedBlocks": 1, "persistedKeys": 19, "headerHeight": 2, "blockHeight": 2, "took": "384.966µs"}
2020-12-17T14:53:05.222+0300	INFO	blockchain persist completed	{"persistedBlocks": 1, "persistedKeys": 19, "headerHeight": 3, "blockHeight": 3, "took": "496.654µs"}
...
```

#### Step 4
Transfer some GAS from multisig account to our account.

1. Create NEP17 transfer transaction:
    ```
        $ ./bin/neo-go wallet nep17 transfer -w .docker/wallets/wallet1.json --out my_tx.json -r http://localhost:20331 --from NgEisvCqr2h8wpRxQb7bVPWUZdbVCY8Uo6 --to NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S --token GAS --amount 29999999
    ``` 
    Where
    - `./bin/neo-go` runs neo-go
    - `wallet nep17 transfer` - command with arguments in [neo-go](https://github.com/nspcc-dev/neo-go/blob/master/cli/wallet/nep17.go#L108)
    - `-w .docker/wallets/wallet1.json` - path to the [wallet](https://github.com/nspcc-dev/neo-go/blob/master/.docker/wallets/wallet1.json) for the first node in the private network
    - `--out my_tx.json` - output file for the signed transaction
    - `-r http://localhost:20331` - RPC node endpoint
    - `--from NgEisvCqr2h8wpRxQb7bVPWUZdbVCY8Uo6` - multisig account to transfer GAS from
    - `--to NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S` - our account from the [wallet](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/my_wallet.json)
    - `--token GAS` - transferred token name, which is GAS
    - `--amount 29999999` - amount of GAS to transfer
    
    Enter the password `one`:
    ```
    Password >
    ```
    The result is transaction signed by the first node `my_tx.json`.

2. Sign the created transaction using the second node address:

    ```
    $ ./bin/neo-go wallet sign -w .docker/wallets/wallet2.json --in my_tx.json --out my_tx2.json --address NgEisvCqr2h8wpRxQb7bVPWUZdbVCY8Uo6
    ```
    Where
    - `-w .docker/wallets/wallet2.json` - path to the [wallet](https://github.com/nspcc-dev/neo-go/blob/master/.docker/wallets/wallet2.json) for the second node in private network
    - `--in my_tx.json` - previously created transfer transaction
    - `--out my_tx2.json` - output file for the signed transaction
    - `--address NgEisvCqr2h8wpRxQb7bVPWUZdbVCY8Uo6` - multisig account to sign the transaction
    
    Enter the password `two`:
    ```
    Enter password to unlock wallet and sign the transaction
    Password >
    ```
    The result is transaction signed by both first and second nodes.

3. Sign the transaction using the third node address and push it to the chain:
    ```
    $ ./bin/neo-go wallet sign -w ./.docker/wallets/wallet3.json --in my_tx2.json --out my_tx3.json --address NgEisvCqr2h8wpRxQb7bVPWUZdbVCY8Uo6 -r http://localhost:20331
    ```
    Enter the password `three`:
    ```
    Enter password to unlock wallet and sign the transaction
    Password >
    ```
    The result is transaction signed by the first, second and third nodes and deployed to the chain.

4. Check the balance:

    Now you should have 29999999 GAS on the balance of `NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S` account.
    To check the transfer was successfully submitted use `getnep17transfers` RPC call:
    ```
    curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getnep17transfers", "params": ["NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S"] }' localhost:20331 | json_pp
    ```
    The result should look like the following:
```
{
   "id" : 1,
   "result" : {
      "sent" : [],
      "address" : "NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S",
      "received" : [
         {
            "timestamp" : 1616423429953,
            "txhash" : "0x1b123a0f26fdc22c94752a29edd7a669c96284c57523ddcc875f1862ce678c1d",
            "blockindex" : 3,
            "amount" : "2999999900000000",
            "assethash" : "0xd2a4cff31913016155e38e474a2c06d08be276cf",
            "transfernotifyindex" : 0,
            "transferaddress" : "NgEisvCqr2h8wpRxQb7bVPWUZdbVCY8Uo6"
         }
      ]
   },
   "jsonrpc" : "2.0"
}
```


## Workshop. Part 1
Now you have all things done to write your first smart contract, deploy and invoke it. 
Let’s go!

#### Step 1
Create basic "Hello World" smart contract (or use the one presented in this repo):
```
package main

import (
	"github.com/nspcc-dev/neo-go/pkg/interop/runtime"
)

func Main() {
	runtime.Log("Hello, world!")
}
```
And save it as `1-print.go`.

Create a configuration for it:
https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/1-print.yml

#### Step 2
Compile "Hello World" smart contract:
```
$ ./bin/neo-go contract compile -i 1-print.go -c 1-print.yml -m 1-print.manifest.json
```
Where
- `./bin/neo-go` runs neo-go
- `contract compile` command with arguments in [neo-go](https://github.com/nspcc-dev/neo-go/blob/master/cli/smartcontract/smart_contract.go#L105)
- `-i 1-print.go` path to smart contract
- `-c 1-print.yml` path to configuration file
- `-m 1-print.manifest.json` path to manifest file, which is required for smart contract deployment

Result:

Compiled smart-contract: `1-pring.nef` and smart contract manifest `1-print.manifest.json`

To dump all the opcodes, you can use:
```
$ ./bin/neo-go contract inspect -i 1-print.nef
```

#### Step 3
Deploy smart contract to the previously setup network:
```
$ ./bin/neo-go contract deploy -i 1-print.nef -manifest 1-print.manifest.json -r http://localhost:20331 -w my_wallet.json
```

Where
- `contract deploy` is a command for deployment
- `-i 1-print.nef` path to smart contract
- `-manifest 1-print.manifest.json` smart contract manifest file
- `-r http://localhost:20331` node endpoint
- `-w my_wallet.json` wallet to use to get the key for transaction signing (you can use one from the workshop repo)

Enter password `qwerty` for the account:
```
Enter account NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S password >
```

Result:
```
Contract: 2bf79b6255d27a2c13462742a545e4b4f94f2d66
ee04aa50b8084684dd49a58b0a42b586c41146626df109acfb3ecf27dee46378
```

At this point your ‘Hello World’ contract is deployed and could be invoked. Let’s do it as a final step.

#### Step 4
Invoke contract.
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 2bf79b6255d27a2c13462742a545e4b4f94f2d66 main
```

Where
- `contract invokefunction` runs invoke with provided parameters
- `-r http://localhost:20331` defines RPC endpoint used for function call
- `-w my_wallet.json` is a wallet
- `2bf79b6255d27a2c13462742a545e4b4f94f2d66` contract hash got as an output from the previous command (deployment in step 6)
- `main` - method to be called

Enter password `qwerty` for account:
```
Enter account NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S password >
```

Result:
```
Sent invocation transaction acffab15d4bab75d0b88a5cf82fe82bff6e6a5715a88d2caad9d73f480f3e635
```
In the console where you were running step #5 you will get:
```
2020-12-17T15:29:48.790+0300	INFO	runtime log	{"script": "662d4ff9b4e445a5422746132c7ad255629bf72b", "logs": "\"Hello, world!\""}
```
Which means that this contract was executed.

This is it. There are only 4 steps to make deployment and they look easy, aren’t they?
Thank you!

## Workshop. Part 2
In this part we'll look at RPC calls and try to write, deploy and invoke smart contract with storage. 
Let’s go!

### RPC calls
Let's check what's going on under the hood. 
Each neo-go node provides an API interface for obtaining blockchain data from it.
The interface is provided via `JSON-RPC`, and the underlying protocol uses HTTP for communication.

Full `NEO JSON-RPC 3.0 API` described [here](https://docs.neo.org/v3/docs/en-us/reference/rpc/latest-version/api.html).

RPC-server of started in step #3 neo-go node is available on `localhost:20331`, so let's try to perform several RPC calls.

#### GetRawTransaction
[GetRawTransaction](https://docs.neo.org/v3/docs/en-us/reference/rpc/latest-version/api/getrawtransaction.html) returns 
the corresponding transaction information, based on the specified hash value.

Request information about our deployment transaction:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getrawtransaction", "params": ["ee04aa50b8084684dd49a58b0a42b586c41146626df109acfb3ecf27dee46378", 1] }' localhost:20331 | json_pp
```

Where
- `"jsonrpc": "2.0"` is protocol version
- `"id": 1` is id of current request
- `"method": "getrawtransaction"` is requested method
- `"params": ["ee04aa50b8084684dd49a58b0a42b586c41146626df109acfb3ecf27dee46378", 1]` is an array of parameters, 
  where
   - `ee04aa50b8084684dd49a58b0a42b586c41146626df109acfb3ecf27dee46378` is deployment transaction hash
   - `1` is `verbose` parameter for detailed JSON string output
- `json_pp` just makes the JSON output prettier

Result:
```
{
   "jsonrpc" : "2.0",
   "result" : {
      "blocktime" : 1616423565029,
      "signers" : [
         {
            "account" : "0x896a32130b231858306db612477bf4d1d0ffcb79",
            "scopes" : "None"
         }
      ],
      "version" : 0,
      "vmstate" : "HALT",
      "blockhash" : "0x82c9e0b857b4ec803d5593a573e1cac9f47cc2d4ac1feaf710664ae5a11ea750",
      "attributes" : [],
      "validuntilblock" : 13,
      "sysfee" : "1001045530",
      "nonce" : 4227310155,
      "size" : 548,
      "hash" : "0xee04aa50b8084684dd49a58b0a42b586c41146626df109acfb3ecf27dee46378",
      "script" : "DPZ7Im5hbWUiOiJIZWxsb1dvcmxkIGNvbnRyYWN0IiwiYWJpIjp7Im1ldGhvZHMiOlt7Im5hbWUiOiJtYWluIiwib2Zmc2V0IjowLCJwYXJhbWV0ZXJzIjpbXSwicmV0dXJudHlwZSI6IlZvaWQiLCJzYWZlIjpmYWxzZX1dLCJldmVudHMiOltdfSwiZ3JvdXBzIjpbXSwicGVybWlzc2lvbnMiOlt7ImNvbnRyYWN0IjoiKiIsIm1ldGhvZHMiOiIqIn1dLCJzdXBwb3J0ZWRzdGFuZGFyZHMiOltdLCJ0cnVzdHMiOltdLCJleHRyYSI6bnVsbH0MZE5FRjNuZW8tZ28tMC45NC4xLXByZS00LWcyOGRhMDBmMAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAWDA1IZWxsbywgd29ybGQhQc/nR5YhQIr9WvcSwBsMBmRlcGxveQwU/aP6Q0bqUyolj8SX3a3bZDfJ/f9BYn1bUg==",
      "confirmations" : 12,
      "sender" : "NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S",
      "witnesses" : [
         {
            "verification" : "DCEDhEhWuuSSNuCc7nLsxQhI8nFlt+UfY3oP0/UkYmdH7G5BdHR2qg==",
            "invocation" : "DED8Q6jOp1locqcwsK+Je1Cy8WvJZ/29p2A3dmHBueBYrWt9plly2UkyX9pGaPOABA4pxQlhTwjzew5E/L0BifRq"
         }
      ],
      "netfee" : "1531520"
   },
   "id" : 1
}
```

#### GetApplicationLog
[GetApplicationLog](https://docs.neo.org/v3/docs/en-us/reference/rpc/latest-version/api/getapplicationlog.html) returns the contract log based on the specified transaction id.

Request application log for invocation transaction from step #4:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["acffab15d4bab75d0b88a5cf82fe82bff6e6a5715a88d2caad9d73f480f3e635"] }' localhost:20331 | json_pp
```

With a single parameter:
- `acffab15d4bab75d0b88a5cf82fe82bff6e6a5715a88d2caad9d73f480f3e635` - invocation transaction hash from step #7

Result:
```
{
   "jsonrpc" : "2.0",
   "id" : 1,
   "result" : {
      "executions" : [
         {
            "vmstate" : "HALT",
            "stack" : [
               {
                  "type" : "Any"
               }
            ],
            "notifications" : [],
            "trigger" : "Application",
            "gasconsumed" : "2028330"
         }
      ],
      "txid" : "0xacffab15d4bab75d0b88a5cf82fe82bff6e6a5715a88d2caad9d73f480f3e635"
   }
}
```

#### Other Useful RPC calls
```
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getversion", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getblockcount", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getconnectioncount", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getcontractstate", "params": ["0x2bf79b6255d27a2c13462742a545e4b4f94f2d66"] }' localhost:20331
```

List of supported by neo-go node RPC commands you can find [here](https://github.com/nspcc-dev/neo-go/blob/master/docs/rpc.md#supported-methods).

### Storage smart contract

Let's take a look at the another smart contract example: [2-storage.go](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/2-storage.go).
This contract is quite simple and, as the previous one, doesn't take any arguments.
On the other hand, it is able to count the number of its own invocations by storing an integer value and increment it after each invocation.
We are interested in this contract as far as it's able to *store* values, i.e. it has a *storage* which can be shared within all contract invocations.
We have to pay some GAS for storage usage, the amount depends on the storage operation (e.g. put) and data size.


This contract also has a special internal `_deploy` method which is executed when the contract is deployed or updated.
It should return no value and accept single bool argument which will be true on contract update.
Our `_deploy` method is aimed to initialise the storage value with `0` when the contract will be deployed.

Now, when we learned about the storage, let's try to deploy and invoke our contract!

#### Step #1
Compile smart contract [2-storage.go](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/2-storage.go):

```
$ ./bin/neo-go contract compile -i 2-storage.go -c 2-storage.yml -m 2-storage.manifest.json
```

Result:

Compiled smart-contract: `2-storage.nef` and smart contract manifest `2-storage.manifest.json`

#### Step #2
Deploy compiled smart contract:
```
$ ./bin/neo-go contract deploy -i 2-storage.nef -manifest 2-storage.manifest.json -r http://localhost:20331 -w my_wallet.json
```
... enter the password `qwerty`:
```
Enter account NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S password >
```

Result:
```
Contract: df91739a75d41b1915fa63f39420aac35a91058c
d34a30968a5ee9f166466667d01d9faa497003e74fc2bc3779e64bcb198b8142
```   

Which means that our contract was deployed and now we can invoke it.

Let's check that the storage value was initialised with `0`. Use `getapplicaionlog` RPC-call for the deployment transaction:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["d34a30968a5ee9f166466667d01d9faa497003e74fc2bc3779e64bcb198b8142"] }' localhost:20331 | json_pp
```

The JSON result is:
```
{
   "result" : {
      "txid" : "0xd34a30968a5ee9f166466667d01d9faa497003e74fc2bc3779e64bcb198b8142",
      "executions" : [
         {
            "stack" : [
               ... skipped serialized contract representation ...
                  ],
                  "type" : "Array"
               }
            ],
            "notifications" : [
               {
                  "contract" : "0xdf91739a75d41b1915fa63f39420aac35a91058c",
                  "eventname" : "info",
                  "state" : {
                     "value" : [
                        {
                           "value" : "U3RvcmFnZSBrZXkgbm90IHlldCBzZXQuIFNldHRpbmcgdG8gMA==",
                           "type" : "ByteString"
                        }
                     ],
                     "type" : "Array"
                  }
               },
               {
                  "eventname" : "info",
                  "state" : {
                     "value" : [
                        {
                           "value" : "U3RvcmFnZSBrZXkgaXMgaW5pdGlhbGlzZWQ=",
                           "type" : "ByteString"
                        }
                     ],
                     "type" : "Array"
                  },
                  "contract" : "0xdf91739a75d41b1915fa63f39420aac35a91058c"
               },
               {
                  "contract" : "0xfffdc93764dbaddd97c48f252a53ea4643faa3fd",
                  "state" : {
                     "value" : [
                        {
                           "value" : "jAWRWsOqIJTzY/oVGRvUdZpzkd8=",
                           "type" : "ByteString"
                        }
                     ],
                     "type" : "Array"
                  },
                  "eventname" : "Deploy"
               }
            ],
            "vmstate" : "HALT",
            "trigger" : "Application",
            "gasconsumed" : "1006244000"
         }
      ]
   },
   "jsonrpc" : "2.0",
   "id" : 1
}
```

Pay attention to the `notifications` JSON field: it contains two `info` notifications with base64-encoded messages.
To decode them just use `echo string | base64 -d` CLI command, e.g.:
```
$ echo U3RvcmFnZSBrZXkgbm90IHlldCBzZXQuIFNldHRpbmcgdG8gMA== | base64 -d
```
which results in `Storage key not yet set. Setting to 0` and
```
$ echo U3RvcmFnZSBrZXkgaXMgaW5pdGlhbGlzZWQ= | base64 -d
```
which is `Storage key is initialised`.

#### Step #3
Let's invoke our contract. As far as we have never invoked this contract, it should increment value from the storage (which is `0`) and put the new `1` value back into the storage.
Let's check:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json df91739a75d41b1915fa63f39420aac35a91058c main
```
... enter the password `qwerty`:
```
Enter account NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S password >
```
Result:
```
Sent invocation transaction 585ae9f91a30f3332605d23c3166f0f5aa97711ebfa742b4e9a7402cbe57e82b
```
To check the counter value, call `getapplicaionlog` RPC-call for the invocation transaction:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["585ae9f91a30f3332605d23c3166f0f5aa97711ebfa742b4e9a7402cbe57e82b"] }' localhost:20331 | json_pp
```
The JSON result is:
```
{
   "id" : 1,
   "jsonrpc" : "2.0",
   "result" : {
      "txid" : "0x585ae9f91a30f3332605d23c3166f0f5aa97711ebfa742b4e9a7402cbe57e82b",
      "executions" : [
         {
            "trigger" : "Application",
            "vmstate" : "HALT",
            "notifications" : [
               {
                  "contract" : "0xdf91739a75d41b1915fa63f39420aac35a91058c",
                  "state" : {
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "VmFsdWUgcmVhZCBmcm9tIHN0b3JhZ2U="
                        }
                     ],
                     "type" : "Array"
                  },
                  "eventname" : "info"
               },
               {
                  "contract" : "0xdf91739a75d41b1915fa63f39420aac35a91058c",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "value" : "U3RvcmFnZSBrZXkgYWxyZWFkeSBzZXQuIEluY3JlbWVudGluZyBieSAx",
                           "type" : "ByteString"
                        }
                     ]
                  },
                  "eventname" : "info"
               },
               {
                  "contract" : "0xdf91739a75d41b1915fa63f39420aac35a91058c",
                  "eventname" : "info",
                  "state" : {
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "TmV3IHZhbHVlIHdyaXR0ZW4gaW50byBzdG9yYWdl"
                        }
                     ],
                     "type" : "Array"
                  }
               }
            ],
            "gasconsumed" : "7233580",
            "stack" : [
               {
                  "value" : "1",
                  "type" : "Integer"
               }
            ]
         }
      ]
   }
}
```
Pay attention to `notifications` field. It contains messages, which where passed to `runtime.Notify` method.
This one contains base64 byte arrays which can be decoded into 3 messages. 
To decode them just use `echo string | base64 -d` CLI command, e.g.:
```
$ echo VmFsdWUgcmVhZCBmcm9tIHN0b3JhZ2U= | base64 -d
```
which results in:
```
Value read from storage
```
So, these 3 messages are:
  - `Value read from storage` which was called after we've got the counter value from storage
  - `Storage key already set. Incrementing by 1` which was called when we realised that counter value is 0
  - `New value written into storage` which was called after the counter value was put in the storage.
  
The final part is `stack` field. This field contains all returned by the contract values, so here you can see integer value `1`,
which is the counter value denoted to the number of contract invocations.

#### Step #4
To ensure that all works as expected, let's invoke the contract one more time and check, whether the counter will be incremented: 
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json df91739a75d41b1915fa63f39420aac35a91058c main
```
... enter the password `qwerty`:
```
Enter account NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S password >
```
Result:
```
Sent invocation transaction 7cb9349e5cb21c57e091db75ffdf9975ec9ea88cf63ad2a20d1c8f0b40a0ce1c
```
To check the counter value, call `getapplicaionlog` RPC-call for the invocation transaction:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["7cb9349e5cb21c57e091db75ffdf9975ec9ea88cf63ad2a20d1c8f0b40a0ce1c"] }' localhost:20331 | json_pp
```
The JSON result is:
```
{
   "jsonrpc" : "2.0",
   "id" : 1,
   "result" : {
      "executions" : [
         {
            "trigger" : "Application",
            "vmstate" : "HALT",
            "gasconsumed" : "7233580",
            "notifications" : [
               {
                  "contract" : "0xdf91739a75d41b1915fa63f39420aac35a91058c",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "VmFsdWUgcmVhZCBmcm9tIHN0b3JhZ2U="
                        }
                     ]
                  },
                  "eventname" : "info"
               },
               {
                  "state" : {
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "U3RvcmFnZSBrZXkgYWxyZWFkeSBzZXQuIEluY3JlbWVudGluZyBieSAx"
                        }
                     ],
                     "type" : "Array"
                  },
                  "contract" : "0xdf91739a75d41b1915fa63f39420aac35a91058c",
                  "eventname" : "info"
               },
               {
                  "contract" : "0xdf91739a75d41b1915fa63f39420aac35a91058c",
                  "state" : {
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "TmV3IHZhbHVlIHdyaXR0ZW4gaW50byBzdG9yYWdl"
                        }
                     ],
                     "type" : "Array"
                  },
                  "eventname" : "info"
               }
            ],
            "stack" : [
               {
                  "value" : "2",
                  "type" : "Integer"
               }
            ]
         }
      ],
      "txid" : "0x7cb9349e5cb21c57e091db75ffdf9975ec9ea88cf63ad2a20d1c8f0b40a0ce1c"
   }
}
```

The `stack` field contains now `2` integer value, so the counter was incremented as we expected.

## Workshop. Part 3
In this part we'll know about NEP5 token standard and try to write, deploy and invoke more complicated smart contract. 
Let’s go!

### NEP17
[NEP17](https://github.com/neo-project/proposals/blob/master/nep-17.mediawiki) is a token standard for the Neo blockchain that provides systems with a generalized interaction mechanism for tokenized smart contracts.
The example with implementation of all required by the standard methods you can find in [nep17.go](https://github.com/nspcc-dev/neo-go/blob/master/examples/token/nep17/nep17.go)
 
Let's take a view on the example of smart contract with NEP17: [token.go](https://github.com/nspcc-dev/neo-go/blob/master/examples/token/token.go)
 
This smart contract initialises nep17 token interface and takes operation string as a parameter, which is one of:
- `symbol` returns ticker symbol of the token
- `decimals` returns amount of decimals for the token
- `totalSupply` returns total token * multiplier
- `balanceOf` returns the token balance of a specific address and requires additional argument:
  - `account` which is requested address
- `transfer` transfers token from one user to another and requires additional arguments:
  - `from` is account which you'd like to transfer tokens from
  - `to` is account which you'd like to transfer tokens to
  - `amount` is the amount of token to transfer
  - `data` is any additional parameter which shall be passed to `onNEP17Payment` method (if the receiver is a contract)
Let's perform several operations with our contract.

#### Step #1
To compile [token.go](https://github.com/nspcc-dev/neo-go/blob/master/examples/token/token.go)
you can use [configuration](https://github.com/nspcc-dev/neo-go/blob/master/examples/token/token.yml).

Compile smart contract:
```
$ ./bin/neo-go contract compile -i examples/token/token.go -c examples/token/token.yml -m examples/token/token.manifest.json
```

Deploy smart contract:
```
$ ./bin/neo-go contract deploy -i examples/token/token.nef -manifest examples/token/token.manifest.json -r http://localhost:20331 -w my_wallet.json
```
... enter the password `qwerty`:
```
Enter account NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S password >
```

Result:
```
Contract: 3dabc1861c671bd8a4f2826eea4d37e2487f80f4
4956eb56fe2aec1c9d33599ba02c8a38c5e215866a2e5846f03458c6a077f4f3

```   

Which means that our contract was deployed and now we can invoke it.

#### Step #2
Let's invoke the contract to perform different operations.

To start with, query `Symbol` of the created nep17 token:

```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 3dabc1861c671bd8a4f2826eea4d37e2487f80f4 symbol
```                                                                   
Where
- `3dabc1861c671bd8a4f2826eea4d37e2487f80f4` is our contract hash from step #1
- `symbol` is operation string which was described earlier and returns token symbol

... and don't forget the password of your account `qwerty`.

Result:
```
Sent invocation transaction c95a56e378b005fe5e71d3e60be2d30869713b65f5722ab3b869e4071224b3b8
```                                                                                         
Now, let's take a detailed look at this invocation transaction with `getapplicationlog` RPC call:

```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["c95a56e378b005fe5e71d3e60be2d30869713b65f5722ab3b869e4071224b3b8"] }' localhost:20331 | json_pp
```               

Result:
```
{
   "id" : 1,
   "result" : {
      "executions" : [
         {
            "trigger" : "Application",
            "notifications" : [],
            "vmstate" : "HALT",
            "gasconsumed" : "4294290",
            "stack" : [
               {
                  "value" : "QU5U",
                  "type" : "ByteString"
               }
            ]
         }
      ],
      "txid" : "0xc95a56e378b005fe5e71d3e60be2d30869713b65f5722ab3b869e4071224b3b8"
   },
   "jsonrpc" : "2.0"
}
```

At least, you can see that `stack` field of JSON result is not empty: it contains base64 byte array with the symbol of our token.

Following commands able you to get some additional information about token:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 3dabc1861c671bd8a4f2826eea4d37e2487f80f4 decimals
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 3dabc1861c671bd8a4f2826eea4d37e2487f80f4 totalSupply
```

#### Step #3

Now it's time for more interesting things. First of all, let's check the balance of nep17 token on our account by using `balanceOf`:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 3dabc1861c671bd8a4f2826eea4d37e2487f80f4 balanceOf NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S
```                             
... with `qwerty` password. The result is:
```
Sent invocation transaction 9e56932c17761013f2b2b38fc0a1b8e08d3354e4c0b2889aec080b0adfed6621
```
And take a closer look at the transaction's details with `getapplicationlog` RPC-call:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["9e56932c17761013f2b2b38fc0a1b8e08d3354e4c0b2889aec080b0adfed6621"] }' localhost:20331 | json_pp
```
Result:
```
{
   "result" : {
      "txid" : "0x9e56932c17761013f2b2b38fc0a1b8e08d3354e4c0b2889aec080b0adfed6621",
      "executions" : [
         {
            "notifications" : [],
            "gasconsumed" : "5311140",
            "vmstate" : "HALT",
            "stack" : [
               {
                  "type" : "Integer",
                  "value" : "0"
               }
            ],
            "trigger" : "Application"
         }
      ]
   },
   "id" : 1,
   "jsonrpc" : "2.0"
}
``` 
As far as `stack` field contains integer value `0`, we have no token on the balance. But don't worry about that. Just follow the next step.

#### Step #4

Before we are able to start using our token (e.g. transfer it to someone else), we have to *mint* it.
In other words, we should transfer all available amount of token (total supply) to someone's account.
There's a special function for this purpose in our contract - `Mint` function. However, this function
uses `CheckWitness` runtime syscall to check whether the caller of the contract is the owner and authorized
to manage initial supply of tokens. That's the purpose of transaction's *signers*: checking given hash
against the values provided in the list of signers. To pass this check we should add our account to
transaction's signers list with CalledByEntry scope. So let's mint token to our address:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 3dabc1861c671bd8a4f2826eea4d37e2487f80f4 mint NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S -- NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S:CalledByEntry
```
Where:
- `--` is a special delimiter of transaction's cosigners list
- `NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S` is the signer itself (which is our account)

... with `qwerty` pass. The result:
``` 
Sent invocation transaction 94d1881679d7686980f7dd11d9196cc933a929c7e279e1a652c8e8487bd58fea
```
`getapplicationlog` RPC-call for this transaction tells us the following:
```
{
   "result" : {
      "executions" : [
         {
            "stack" : [
               {
                  "type" : "Boolean",
                  "value" : true
               }
            ],
            "gasconsumed" : "16522950",
            "vmstate" : "HALT",
            "notifications" : [
               {
                  "state" : {
                     "value" : [
                        {
                           "type" : "Any"
                        },
                        {
                           "value" : "ecv/0NH0e0cStm0wWBgjCxMyaok=",
                           "type" : "ByteString"
                        },
                        {
                           "type" : "Integer",
                           "value" : "1100000000000000"
                        }
                     ],
                     "type" : "Array"
                  },
                  "contract" : "0x3dabc1861c671bd8a4f2826eea4d37e2487f80f4",
                  "eventname" : "Transfer"
               }
            ],
            "trigger" : "Application"
         }
      ],
      "txid" : "0x94d1881679d7686980f7dd11d9196cc933a929c7e279e1a652c8e8487bd58fea"
   },
   "id" : 1,
   "jsonrpc" : "2.0"
}
```
Here we have `true` at the `stack` field, which means that token was successfully minted.
Let's just ensure that by querying `balanceOf` one more time:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 3dabc1861c671bd8a4f2826eea4d37e2487f80f4 balanceOf NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S
```
... with `qwerty` pass. The result:
``` 
Sent invocation transaction 3f4915fb015dc0592f9e1c6675280647344e3e099259ddd954d9c51cedf8e0ef
```
... with the following `getapplicationlog` JSON message:
```
{
   "jsonrpc" : "2.0",
   "result" : {
      "txid" : "0x3f4915fb015dc0592f9e1c6675280647344e3e099259ddd954d9c51cedf8e0ef",
      "executions" : [
         {
            "stack" : [
               {
                  "value" : "1100000000000000",
                  "type" : "Integer"
               }
            ],
            "gasconsumed" : "5557020",
            "vmstate" : "HALT",
            "notifications" : [],
            "trigger" : "Application"
         }
      ]
   },
   "id" : 1
}
```
Now we can see integer value at the `stack` field, so `1100000000000000` is the nep17 token balance of our account.

Note, that token can be minted only once.

#### Step #5

After we are done with minting, it's possible to transfer token to someone else.
Let's transfer 5 tokens from our account to `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm` with `transfer` call:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 3dabc1861c671bd8a4f2826eea4d37e2487f80f4 transfer NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm 500000000 null -- NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S:CalledByEntry
```
... with password `qwerty` and following result:
``` 
Sent invocation transaction c6eb182efb23e4238dbb7cc622c6a1e7d8a3811efbf70d262188afb0e1b40660
```
Our favourite `getapplicationlog` RPC-call tells us:
```
{
   "result" : {
      "executions" : [
         {
            "trigger" : "Application",
            "gasconsumed" : "14760950",
            "stack" : [
               {
                  "value" : true,
                  "type" : "Boolean"
               }
            ],
            "notifications" : [
               {
                  "contract" : "0x3dabc1861c671bd8a4f2826eea4d37e2487f80f4",
                  "state" : {
                     "value" : [
                        {
                           "value" : "ecv/0NH0e0cStm0wWBgjCxMyaok=",
                           "type" : "ByteString"
                        },
                        {
                           "type" : "ByteString",
                           "value" : "50l6vFaauRKm8hPVkr3Aw2CeHQs="
                        },
                        {
                           "type" : "Integer",
                           "value" : "500000000"
                        }
                     ],
                     "type" : "Array"
                  },
                  "eventname" : "Transfer"
               }
            ],
            "vmstate" : "HALT"
         }
      ],
      "txid" : "0xc6eb182efb23e4238dbb7cc622c6a1e7d8a3811efbf70d262188afb0e1b40660"
   },
   "jsonrpc" : "2.0",
   "id" : 1
}
```
Note, that `stack` field contains `true`, which means that token was successfully transferred.
Let's now check the balance of `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm` account to ensure that the amount of token on that account = 5:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 3dabc1861c671bd8a4f2826eea4d37e2487f80f4 balanceOf NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm
```
The `getapplicationlog` RPC-call for this transaction tells us the following:
```
{
   "result" : {
      "txid" : "0x9f6afc7263564f349ef01f66baddf0db6607c88d3309859a91b91b452c57345f",
      "executions" : [
         {
            "vmstate" : "HALT",
            "notifications" : [],
            "gasconsumed" : "5557020",
            "trigger" : "Application",
            "stack" : [
               {
                  "type" : "Integer",
                  "value" : "500000000"
               }
            ]
         }
      ]
   },
   "id" : 1,
   "jsonrpc" : "2.0"
}
```
Here we are! There are exactly 5 tokens at the `stack` field. You can also ensure that these 5 tokens were debited from `NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S` account by using `balanceOf` method.

## Workshop. Part 4
In this part we'll summarise our knowledge about smart contracts by investigating [4-domain](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/4-domain.go) smart contract. This contract 
contains code for domain registration, transferring, deletion and getting information about registered domains.

Let’s go!

#### Step #1
Let's take a glance at our [contract](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/4-domain.go) and inspect it. The contract takes an action string as the first parameter, which is one of the following:
- `register` checks, whether domain with the specified name already exists. If not, it also adds the pair `[domainName, owner]` to the storage. It requires additional arguments:
   - `domainName` which is the new domain name.
   - `owner` - the 34-digit account address from our [wallet](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/my_wallet.json), which will be used for contract invocation.
- `query` returns the specified domain owner address (or false, if no such domain was registered). It requires the following argument:
   - `domainName` which is requested domain name.
- `transfer` transfers domain with the specified name to the other address (of course, in case if you're the actual owner of the domain requested). It requires additional arguments:
   - `domainName` which is the name of domain you'd like to transfer.
   - `toAddress` - the account address you'd like to transfer the specified domain to.
- `delete` deletes the specified domain from the storage. The arguments:
   - `domainName` which is the name of the domain you'd like to delete.
 
 In the next steps we'll compile and deploy smart contract. 
 After that we'll try to register new domain, transfer it to another account and query information about it.

#### Step #2

Compile smart contract [4-domain.go](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/4-domain.go) with [configuration](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/4-domain.yml)
```
$ ./bin/neo-go contract compile -i 4-domain.go -c 4-domain.yml -m 4-domain.manifest.json
```

... and deploy it:
```
$ ./bin/neo-go contract deploy -i 4-domain.nef --manifest 4-domain.manifest.json -r http://localhost:20331 -w my_wallet.json
```
Just a note: our contract uses storage and, as the previous one, needs the flag `hasstorage` to be set to `true` value.
That can be done in [configuration](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/4-domain.yml) file.

... enter the password `qwerty`:
```
Enter account NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S password >
```

Result:
```
Contract: 3fc16e7ec0ba746caac13629599f9b1287808d0c
f11325384d6cd030550484c0477d58cad6fed728eac9d60d0aef4ba0b306f72f
```   
You know, what it means :)

#### Step #3

Invoke the contract to register domain with name `my_first_domain`:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 3fc16e7ec0ba746caac13629599f9b1287808d0c register my_first_domain NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S -- NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S:CalledByEntry
```
... the strongest password in the world, guess: `qwerty`
```
Enter account NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S password >
```
Result:
```
Sent invocation transaction e620cfb874102f0da2f1b1fe83820fdbb4f09156c4443ef0fe7f77591ab366a5
```
Also you can see the log message in the console, where you run neo-go node:
```
2020-12-17T17:31:43.480+0300	INFO	runtime log	{"script": "0c8d8087129b9f592936c1aa6c74bac07e6ec13f", "logs": "\"RegisterDomain: my_first_domain\""}
```
Well, that's ok. Let's check now, whether our domain was registered with `getapplicationlog` RPC-call:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["e620cfb874102f0da2f1b1fe83820fdbb4f09156c4443ef0fe7f77591ab366a5"] }' localhost:20331 | json_pp
```
The result is:
```
{
   "jsonrpc" : "2.0",
   "id" : 1,
   "result" : {
      "txid" : "0xe620cfb874102f0da2f1b1fe83820fdbb4f09156c4443ef0fe7f77591ab366a5",
      "executions" : [
         {
            "trigger" : "Application",
            "gasconsumed" : "9143210",
            "vmstate" : "HALT",
            "notifications" : [
               {
                  "eventname" : "registered",
                  "state" : {
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "ecv/0NH0e0cStm0wWBgjCxMyaok="
                        },
                        {
                           "value" : "bXlfZmlyc3RfZG9tYWlu",
                           "type" : "ByteString"
                        }
                     ],
                     "type" : "Array"
                  },
                  "contract" : "0x3fc16e7ec0ba746caac13629599f9b1287808d0c"
               }
            ],
            "stack" : [
               {
                  "value" : true,
                  "type" : "Boolean"
               }
            ]
         }
      ]
   }
}
```
Especially, we're interested in two fields of the json:

First one is `notifications` field, which contains one notification with `registered` name:
- `bXlfZmlyc3RfZG9tYWlu` byte string in base64 representation, which can be decoded to `my_first_domain` - our domain's name
- `ecv/0NH0e0cStm0wWBgjCxMyaok=` byte array, which can be decoded to the account address `NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S`.

The second field is `stack` with `true` value, which was returned by the smart contract.

All of these values let us be sure that our domain was successfully registered.  

#### Step #4

Invoke the contract to query the address information our `my_first_domain` domain:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 3fc16e7ec0ba746caac13629599f9b1287808d0c query my_first_domain
```
... the pass `qwerty`:
```
Enter account NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S password >
```
Result:
```
Sent invocation transaction 8e23be4fb2effde35d8d5b806a6d0a345ce0499d8df465122530bee9fb6c473a
```
and log-message:
```
2020-12-17T17:39:32.677+0300	INFO	runtime log	{"script": "0c8d8087129b9f592936c1aa6c74bac07e6ec13f", "logs": "\"QueryDomain: my_first_domain\""}
```
Let's check this transaction with `getapplicationlog` RPC call:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["8e23be4fb2effde35d8d5b806a6d0a345ce0499d8df465122530bee9fb6c473a"] }' localhost:20331 | json_pp
```
... which gives us the following result:
```
{
   "jsonrpc" : "2.0",
   "id" : 1,
   "result" : {
      "txid" : "0x8e23be4fb2effde35d8d5b806a6d0a345ce0499d8df465122530bee9fb6c473a",
      "executions" : [
         {
            "trigger" : "Application",
            "gasconsumed" : "4321230",
            "notifications" : [],
            "stack" : [
               {
                  "value" : "ecv/0NH0e0cStm0wWBgjCxMyaok=",
                  "type" : "ByteString"
               }
            ],
            "vmstate" : "HALT"
         }
      ]
   }
}
```

with base64 interpretation of our account address `NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S` on the stack, which means that domain `my_first_domain` was registered by the owner with received account address.

#### Step #5

Invoke the contract to transfer domain to the other account (e.g. account with `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm` address):
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 3fc16e7ec0ba746caac13629599f9b1287808d0c transfer my_first_domain NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm -- NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S
```
... the password: `qwerty`
```
Enter account NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S password >
```
Result:
```
Sent invocation transaction c30002b6ada7d56525fe6b091d52c8db4897f4773318ad46b17537ee0ee3aaf5
```
and log-message:
```
2020-12-17T17:44:07.536+0300	INFO	runtime log	{"script": "0c8d8087129b9f592936c1aa6c74bac07e6ec13f", "logs": "\"TransferDomain: my_first_domain\""}
```
Perfect. And `getapplicationlog` RPC-call...
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["c30002b6ada7d56525fe6b091d52c8db4897f4773318ad46b17537ee0ee3aaf5"] }' localhost:20331 | json_pp
```
... tells us:
```
{
   "id" : 1,
   "result" : {
      "executions" : [
         {
            "trigger" : "Application",
            "gasconsumed" : "7680110",
            "notifications" : [
               {
                  "eventname" : "deleted",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "ecv/0NH0e0cStm0wWBgjCxMyaok="
                        },
                        {
                           "value" : "bXlfZmlyc3RfZG9tYWlu",
                           "type" : "ByteString"
                        }
                     ]
                  },
                  "contract" : "0x3fc16e7ec0ba746caac13629599f9b1287808d0c"
               },
               {
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "50l6vFaauRKm8hPVkr3Aw2CeHQs="
                        },
                        {
                           "value" : "bXlfZmlyc3RfZG9tYWlu",
                           "type" : "ByteString"
                        }
                     ]
                  },
                  "contract" : "0x3fc16e7ec0ba746caac13629599f9b1287808d0c",
                  "eventname" : "registered"
               }
            ],
            "stack" : [
               {
                  "type" : "Boolean",
                  "value" : true
               }
            ],
            "vmstate" : "HALT"
         }
      ],
      "txid" : "0xc30002b6ada7d56525fe6b091d52c8db4897f4773318ad46b17537ee0ee3aaf5"
   },
   "jsonrpc" : "2.0"
}
```
The `notifications` field contains two events:
- First one with name `deleted` and additional information (domain `my_first_domain` was deleted from account `NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S`),
- Second one with name `registered`  and additional information  (domain `my_first_domain` was registered with account `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm`).
The `stack` field contains `true` value, which means that our domain was successfully transferred.

#### Step #6

The last call is `delete`, so you can try to create the other domain, e.g. `my_second_domain` and then remove it from storage with:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 3fc16e7ec0ba746caac13629599f9b1287808d0c register my_second_domain NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S -- NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S:CalledByEntry
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 3fc16e7ec0ba746caac13629599f9b1287808d0c delete my_second_domain -- NX1yL5wDx3inK2qUVLRVaqCLUxYnAbv85S
```

Thank you!

### Useful links

* [Our basic tutorial on Medium](https://medium.com/@neospcc/%D1%81%D0%BC%D0%B0%D1%80%D1%82-%D0%BA%D0%BE%D0%BD%D1%82%D1%80%D0%B0%D0%BA%D1%82-%D0%B4%D0%BB%D1%8F-neo-769139352b65)
* [Using NEO Blockchain Toolkit](https://medium.com/@neospcc/neogo-adds-support-for-neo-blockchain-toolkit-673ea914f661)
* [NEO documentation](https://docs.neo.org/)
* [NEO github](https://github.com/neo-project/neo/)
* [NEO-GO github](https://github.com/nspcc-dev/neo-go)
