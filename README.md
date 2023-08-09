<p align="center">
<img src="./pic/neo_color_dark_gopher.png" width="300px" alt="logo">
</p>

[Neo](https://neo.org/) builds smart economy and we at [NeoSPCC](https://nspcc.ru/en/) help them with that big challenge. 
In our blog you might find the latest articles [how we run NeoFS public test net](https://medium.com/@neospcc/public-neofs-testnet-launch-18f6315c5ced) 
but it’s not the only thing we’re working on.

## NeoGo
As you know network is composed of nodes. These nodes as of now have several implementations:
- https://github.com/neo-project/neo
- https://github.com/nspcc-dev/neo-go

This article is about the last one since we’re developing it at NeoSPCC. 
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
NAME:
   VM CLI - Official VM CLI for Neo-Go

USAGE:
    [global options] command [command options] [arguments...]

VERSION:
   0.98.2

COMMANDS:
   exit        Exit the VM prompt
   ip          Show current instruction
   break       Place a breakpoint
   estack      Show evaluation stack contents
   istack      Show invocation stack contents
   sslot       Show static slot contents
   lslot       Show local slot contents
   aslot       Show arguments slot contents
   loadnef     Load a NEF-consistent script into the VM
   loadbase64  Load a base64-encoded script string into the VM
   loadhex     Load a hex-encoded script string into the VM
   loadgo      Compile and load a Go file with the manifest into the VM
   reset       Unload compiled script from the VM
   parse       Parse provided argument and convert it into other possible formats
   run         Execute the current loaded script
   cont        Continue execution of the current loaded script
   step        Step (n) instruction in the program
   stepinto    Stepinto instruction to take in the debugger
   stepout     Stepout instruction to take in the debugger
   stepover    Stepover instruction to take in the debugger
   ops         Dump opcodes of the current loaded program
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
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
READY: loaded 21 instructions
NEO-GO-VM 0 >  
```
And there you can see how many instructions were generated and even if you are interested in opcodes of current program you can dump them:
```
NEO-GO-VM 0 > ops
INDEX    OPCODE       PARAMETER
0        PUSHDATA1    48656c6c6f2c20776f726c6421 ("Hello, world!")    <<
15       SYSCALL      System.Runtime.Log (cfe74796)
20       RET
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
Neo has a nice monitor where you can find particular node running in the blockchain network.
[Neo Monitor](http://monitor.cityofzion.io/)

## Workshop. Content
The workshop contains a guide, examples and tips on Neo smart contracts development
and dApps development for Neo ecosystem using development kit offered by NeoGo.
The workshop is divided into multiple parts:
1. [**Preparation.**](#workshop-preparation) Learn how to run a Neo private network,
  transfer some funds from multi-signature account to a simple signature account
  using NeoGo CLI and check the balance via JSON RPC call.
2. [**Part 1.**](#workshop-part-1) Compile, inspect, deploy and invoke simple
  `Hello, world!` smart contract written in Go.
3. [**Part 2.**](#workshop-part-2) Investigate Neo JSON RPC protocol and NeoGo CLI
   utilities to retrieve information from Neo RPC nodes. Get acquainted with the
   concept of a smart contract storage. Compile, deploy and invoke contract that
   demonstrates how to use its storage.
4. [**Part 3.**](#workshop-part-3) Learn more about NEP-17 token standard and take
   a look at the NEP-17 compatible contract.
5. [**Part 4.**](#workshop-part-4) Summarize knowledge about smart contracts. Create,
   compile and invoke more complicated contract.
6. [**Part 5.**](#workshop-part-5) Learn how to develop simple decentralized application
   for the Neo ecosystem using tools provided by NeoGo.
 
## Workshop. Preparation
In this part we will setup the environment: run private network, connect neo-go node to it and transfer some initial GAS to our basic account
in order to be able to pay for transaction deployment and invocation. Let's start.

#### Requirements
For this workshop you will need Debian 10, Docker, docker-compose, go to be installed:
- [docker](https://docs.docker.com/install/linux/docker-ce/debian/)
- [go](https://golang.org/dl/)

#### Step 1
If you already have neo-go or go smart-contracts, please, update go modules in order to be up-to-date with the current interop API changes.
If not, download neo-go and build it (master branch):
```
$ git clone https://github.com/nspcc-dev/neo-go.git
$ cd neo-go
$ make build 
```

#### Step 2
To run NeoGo-based 4-node private network use this commands:

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

1. Create NEP-17 transfer transaction:
    ```
        $ ./bin/neo-go wallet nep17 transfer -w .docker/wallets/wallet1.json --out my_tx.json -r http://localhost:20331 --from NVTiAjNgagDkTr5HTzDmQP9kPwPHN5BgVq --to NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB --token GAS --amount 29999999
    ``` 
    Where
    - `./bin/neo-go` runs neo-go
    - `wallet nep17 transfer` - command with arguments in [neo-go](https://github.com/nspcc-dev/neo-go/blob/master/cli/wallet/nep17.go#L108)
    - `-w .docker/wallets/wallet1.json` - path to the [wallet](https://github.com/nspcc-dev/neo-go/blob/master/.docker/wallets/wallet1.json) for the first node in the private network
    - `--out my_tx.json` - output file for the signed transaction
    - `-r http://localhost:20331` - RPC node endpoint
    - `--from NVTiAjNgagDkTr5HTzDmQP9kPwPHN5BgVq` - multisig account to transfer GAS from
    - `--to NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB` - our account from the [wallet](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/my_wallet.json)
    - `--token GAS` - transferred token name, which is GAS
    - `--amount 29999999` - amount of GAS to transfer
    
    Enter the password `one`:
    ```
    Enter password >
    ```
    The result is transaction signed by the first node `my_tx.json`.

2. Sign the created transaction using the second node address:

    ```
    $ ./bin/neo-go wallet sign -w .docker/wallets/wallet2.json --in my_tx.json --out my_tx2.json --address NVTiAjNgagDkTr5HTzDmQP9kPwPHN5BgVq
    ```
    Where
    - `-w .docker/wallets/wallet2.json` - path to the [wallet](https://github.com/nspcc-dev/neo-go/blob/master/.docker/wallets/wallet2.json) for the second node in private network
    - `--in my_tx.json` - previously created transfer transaction
    - `--out my_tx2.json` - output file for the signed transaction
    - `--address NVTiAjNgagDkTr5HTzDmQP9kPwPHN5BgVq` - multisig account to sign the transaction
    
    Enter the password `two`:
    ```
    Enter password >
    ```
    The result is transaction signed by both first and second nodes.

3. Sign the transaction using the third node address and push it to the chain:
    ```
    $ ./bin/neo-go wallet sign -w ./.docker/wallets/wallet3.json --in my_tx2.json --out my_tx3.json --address NVTiAjNgagDkTr5HTzDmQP9kPwPHN5BgVq -r http://localhost:20331
    ```
    Enter the password `three`:
    ```
    Enter password >
    ```
    The result is transaction signed by the first, second and third nodes and deployed to the chain.

4. Check the balance:

    Now you should have 29999999 GAS on the balance of `NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB` account.
    To check the transfer was successfully submitted use `getnep17transfers` RPC call:
    ```
    curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getnep17transfers", "params": ["NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB"] }' localhost:20331 | json_pp
    ```
    The result should look like the following:
```
{
   "result" : {
      "received" : [
         {
            "txhash" : "0x7f1a2c41f0c03107f7a44ac510fa95fe11dde4c4994d30d61439f73f27e70f0d",
            "transfernotifyindex" : 0,
            "transferaddress" : "NVTiAjNgagDkTr5HTzDmQP9kPwPHN5BgVq",
            "blockindex" : 27,
            "timestamp" : 1657014304108,
            "amount" : "2999999900000000",
            "assethash" : "0xd2a4cff31913016155e38e474a2c06d08be276cf"
         }
      ],
      "address" : "NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB",
      "sent" : []
   },
   "jsonrpc" : "2.0",
   "id" : 1
}
```


## Workshop. Part 1
Now you have all things done to write your first smart contract, deploy and invoke it. 
Let’s go!

#### Step 1
Use the "Hello World" smart contract contained in the repository in its own
[1-print
directory](https://github.com/nspcc-dev/neo-go-sc-wrkshp/tree/master/1-print). The
code is rather simple:
```
package main

import (
	"github.com/nspcc-dev/neo-go/pkg/interop/runtime"
)

func Main() {
	runtime.Log("Hello, world!")
}
```

Contract configuration is available in the same directory, [1-print.yml](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/1-print/1-print.yml).

#### Step 2
Compile "Hello World" smart contract:
```
$ ./bin/neo-go contract compile -i 1-print/1-print.go -c 1-print/1-print.yml -m 1-print/1-print.manifest.json
```
Where
- `./bin/neo-go` runs neo-go
- `contract compile` command with arguments in [neo-go](https://github.com/nspcc-dev/neo-go/blob/master/cli/smartcontract/smart_contract.go#L105)
- `-i 1-print/1-print.go` path to smart contract
- `-c 1-print/1-print.yml` path to configuration file
- `-m 1-print/1-print.manifest.json` path to manifest file, which is required for smart contract deployment

Result:

Compiled smart-contract: `1-pring.nef` and smart contract manifest `1-print.manifest.json`

To dump all the opcodes, you can use:
```
$ ./bin/neo-go contract inspect -i 1-print/1-print.nef
```

#### Step 3
Deploy smart contract to the previously setup network:
```
$ ./bin/neo-go contract deploy -i 1-print/1-print.nef -manifest 1-print/1-print.manifest.json -r http://localhost:20331 -w my_wallet.json
```

Where
- `contract deploy` is a command for deployment
- `-i 1-print/1-print.nef` path to smart contract
- `-manifest 1-print/1-print.manifest.json` smart contract manifest file
- `-r http://localhost:20331` node endpoint
- `-w my_wallet.json` wallet to use to get the key for transaction signing (you can use one from the workshop repo)

Enter password `qwerty` for the account:
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```
CLI will print transaction fees and ask for confirmation. Here and later enter `y`
to relay the transaction:
```
Network fee: 0.0151452
System fee: 10.0104553
Total fee: 10.0256005
Relay transaction (y|N)> y
```
Result:
```
Sent invocation transaction b0436603d27d14e3aa27280e1bc2cdb17d4def8cb8cda2204c3b6a203203e6bf
Contract: bfad19135422aaddf2fc86f86ec5d4b1371e8e93
```

At this point your ‘Hello World’ contract is deployed and could be invoked. Let’s do it as a final step.

#### Step 4
Invoke contract.
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json bfad19135422aaddf2fc86f86ec5d4b1371e8e93 main
```

Where
- `contract invokefunction` runs invoke with provided parameters
- `-r http://localhost:20331` defines RPC endpoint used for function call
- `-w my_wallet.json` is a wallet
- `bfad19135422aaddf2fc86f86ec5d4b1371e8e93` contract hash got as an output from the previous command (deployment in step 6)
- `main` - method to be called

Enter password `qwerty` for account:
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```

Result:
```
Network fee: 0.0117652
System fee: 0.0196731
Total fee: 0.0314383
Relay transaction (y|N)> y
Sent invocation transaction 60fbf79b06714e34a3f55e782bf509eedcc602661e49848f5611aaaf2f3442de
```
In the console where you were running step #5 you will get:
```
2022-07-05T12:52:49.413+0300	INFO	runtime log	{"tx": "60fbf79b06714e34a3f55e782bf509eedcc602661e49848f5611aaaf2f3442de", "script": "bfad19135422aaddf2fc86f86ec5d4b1371e8e93", "msg": "Hello, world!"}
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
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getrawtransaction", "params": ["b0436603d27d14e3aa27280e1bc2cdb17d4def8cb8cda2204c3b6a203203e6bf", 1] }' localhost:20331 | json_pp
```

Where
- `"jsonrpc": "2.0"` is protocol version
- `"id": 1` is id of current request
- `"method": "getrawtransaction"` is requested method
- `"params": ["b0436603d27d14e3aa27280e1bc2cdb17d4def8cb8cda2204c3b6a203203e6bf", 1]` is an array of parameters, 
  where
   - `b0436603d27d14e3aa27280e1bc2cdb17d4def8cb8cda2204c3b6a203203e6bf` is deployment transaction hash
   - `1` is `verbose` parameter for detailed JSON string output
- `json_pp` just makes the JSON output prettier

Result:
```
{
   "id" : 1,
   "result" : {
      "nonce" : 2596996162,
      "sysfee" : "1001045530",
      "size" : 531,
      "attributes" : [],
      "blocktime" : 1657014649346,
      "script" : "DOZ7Im5hbWUiOiJIZWxsb1dvcmxkIGNvbnRyYWN0IiwiYWJpIjp7Im1ldGhvZHMiOlt7Im5hbWUiOiJtYWluIiwib2Zmc2V0IjowLCJwYXJhbWV0ZXJzIjpbXSwicmV0dXJudHlwZSI6IlZvaWQiLCJzYWZlIjpmYWxzZX1dLCJldmVudHMiOltdfSwiZmVhdHVyZXMiOnt9LCJncm91cHMiOltdLCJwZXJtaXNzaW9ucyI6W10sInN1cHBvcnRlZHN0YW5kYXJkcyI6W10sInRydXN0cyI6W10sImV4dHJhIjpudWxsfQxjTkVGM25lby1nby0wLjk5LjEtcHJlLTEwMy1nM2ZiYzEzMzEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABUMDUhlbGxvLCB3b3JsZCFBz+dHlkCTKBNVEsAfDAZkZXBsb3kMFP2j+kNG6lMqJY/El92t22Q3yf3/QWJ9W1I=",
      "validuntilblock" : 51,
      "signers" : [
         {
            "account" : "0x410b5658f92f9937ed7bdd4ba04c665d3bdbd8ae",
            "scopes" : "CalledByEntry"
         }
      ],
      "netfee" : "1514520",
      "sender" : "NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB",
      "version" : 0,
      "confirmations" : 15,
      "witnesses" : [
         {
            "invocation" : "DEDAFsqnEFjXighqESMUGAAZxR2vaDBpYfbMgH55C1Q8TFNl5AfQA7Cder+MCDjPLNu7S1KHqwp97XlZK2OpZGnf",
            "verification" : "DCEDhEhWuuSSNuCc7nLsxQhI8nFlt+UfY3oP0/UkYmdH7G5BVuezJw=="
         }
      ],
      "blockhash" : "0x1fbfb494e669c03a4666fb8b2da9ed2f8205b07aabe55125311bb1f569e83d92",
      "hash" : "0xb0436603d27d14e3aa27280e1bc2cdb17d4def8cb8cda2204c3b6a203203e6bf",
      "vmstate" : "HALT"
   },
   "jsonrpc" : "2.0"
}
```

#### GetApplicationLog
[GetApplicationLog](https://docs.neo.org/v3/docs/en-us/reference/rpc/latest-version/api/getapplicationlog.html) returns the contract log based on the specified transaction id.

Request application log for invocation transaction from step #4:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["60fbf79b06714e34a3f55e782bf509eedcc602661e49848f5611aaaf2f3442de"] }' localhost:20331 | json_pp
```

With a single parameter:
- `60fbf79b06714e34a3f55e782bf509eedcc602661e49848f5611aaaf2f3442de` - invocation transaction hash from step #7

Result:
```
{
   "jsonrpc" : "2.0",
   "id" : 1,
   "result" : {
      "txid" : "0x60fbf79b06714e34a3f55e782bf509eedcc602661e49848f5611aaaf2f3442de",
      "executions" : [
         {
            "gasconsumed" : "1967310",
            "trigger" : "Application",
            "stack" : [
               {
                  "type" : "Any"
               }
            ],
            "vmstate" : "HALT",
            "notifications" : []
         }
      ]
   }
}
```

#### Other Useful RPC calls
```
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getversion", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getblockcount", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getconnectioncount", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getcontractstate", "params": ["0xbfad19135422aaddf2fc86f86ec5d4b1371e8e93"] }' localhost:20331
```

List of supported by neo-go node RPC commands you can find [here](https://github.com/nspcc-dev/neo-go/blob/master/docs/rpc.md#supported-methods).

#### Utilities

neo-go CLI provides `query tx` utility to check the transaction state. It uses
`getrawtransaction` and `getapplicationlog` RPC calls under the hood and prints details
of transaction invocation. Use `query tx` command to ensure transaction was accepted
to the chain:
```
./bin/neo-go query tx 60fbf79b06714e34a3f55e782bf509eedcc602661e49848f5611aaaf2f3442de -r http://localhost:20331 -v
```
where
- `60fbf79b06714e34a3f55e782bf509eedcc602661e49848f5611aaaf2f3442de` - invocation transaction hash from step #7
- `-r http://localhost:20331` - RPC node endpoint
- `-v` - verbose flag (enables transaction's signers, fees and script dumps)

The result is:
```
Hash:			60fbf79b06714e34a3f55e782bf509eedcc602661e49848f5611aaaf2f3442de
OnChain:		true
BlockHash:		da4b4959936208b4658136e338bec1608772cdd16124675efc4862d1b576cf74
Success:		true
Signer:			NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB (None)
SystemFee:		0.0196731 GAS
NetworkFee:		0.0117652 GAS
Script:			wh8MBG1haW4MFJOOHjex1MVu+Ib88t2qIlQTGa2/QWJ9W1I=
INDEX    OPCODE       PARAMETER
0        NEWARRAY0        <<
1        PUSH15       
2        PUSHDATA1    6d61696e ("main")
8        PUSHDATA1    938e1e37b1d4c56ef886fcf2ddaa22541319adbf
30       SYSCALL      System.Contract.Call (627d5b52)
```

From the `OnChain` field we can see that transaction was successfully accepted to the chain.
`Success` field tells us whether transaction script was successfully executed, i.e. changes
in tx has been persisted on chain and VM has `HALT` state after script execution. 


### Storage smart contract

Let's take a look at the another smart contract example: [2-storage.go](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/2-storage/2-storage.go).
This contract is quite simple and, as the previous one, doesn't take any arguments.
On the other hand, it is able to count the number of its own invocations by storing an integer value and increment it after each invocation.
We are interested in this contract as far as it's able to *store* values, i.e. it has a *storage* which can be shared within all contract invocations.
We have to pay some GAS for storage usage, the amount depends on the storage operation (e.g. put) and data size.


This contract also has a special internal `_deploy` method which is executed when the contract is deployed or updated.
It should return no value and accept single bool argument which will be true on contract update.
Our `_deploy` method is aimed to initialise the storage value with `0` when the contract will be deployed.

Now, when we learned about the storage, let's try to deploy and invoke our contract!

#### Step #1
Compile smart contract [2-storage.go](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/2-storage/2-storage.go):

```
$ ./bin/neo-go contract compile -i 2-storage/2-storage.go -c 2-storage/2-storage.yml -m 2-storage/2-storage.manifest.json
```

Result:

Compiled smart-contract: `2-storage.nef` and smart contract manifest `2-storage.manifest.json`

#### Step #2
Deploy compiled smart contract:
```
$ ./bin/neo-go contract deploy -i 2-storage/2-storage.nef -manifest 2-storage/2-storage.manifest.json -r http://localhost:20331 -w my_wallet.json
```
... enter the password `qwerty`:
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```

Result:
```
Network fee: 0.0210952
System fee: 10.0624424
Total fee: 10.0835376
Relay transaction (y|N)> y
Sent invocation transaction 58b98332ad8456da1ab6c9d162c1c557ed2ef85a67ae6bbec17e33cdad599c25
Contract: ccd533440d0317e9f366c50648d0013540e82741
```   

Which means that our contract was deployed and now we can invoke it.

Let's check that the storage value was initialised with `0`. Use `getapplicaionlog` RPC-call for the deployment transaction:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["58b98332ad8456da1ab6c9d162c1c557ed2ef85a67ae6bbec17e33cdad599c25"] }' localhost:20331 | json_pp
```

The JSON result is:
```
{
   "jsonrpc" : "2.0",
   "result" : {
      "txid" : "0x58b98332ad8456da1ab6c9d162c1c557ed2ef85a67ae6bbec17e33cdad599c25",
      "executions" : [
         {
            "vmstate" : "HALT",
            "notifications" : [
               {
                  "contract" : "0xccd533440d0317e9f366c50648d0013540e82741",
                  "state" : {
                     "value" : [
                        {
                           "type" : "Buffer",
                           "value" : "U3RvcmFnZSBrZXkgbm90IHlldCBzZXQuIFNldHRpbmcgdG8gMA=="
                        }
                     ],
                     "type" : "Array"
                  },
                  "eventname" : "info"
               },
               {
                  "contract" : "0xccd533440d0317e9f366c50648d0013540e82741",
                  "eventname" : "info",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "value" : "U3RvcmFnZSBrZXkgaXMgaW5pdGlhbGlzZWQ=",
                           "type" : "Buffer"
                        }
                     ]
                  }
               },
               {
                  "contract" : "0xfffdc93764dbaddd97c48f252a53ea4643faa3fd",
                  "state" : {
                     "value" : [
                        {
                           "value" : "QSfoQDUB0EgGxWbz6RcDDUQz1cw=",
                           "type" : "ByteString"
                        }
                     ],
                     "type" : "Array"
                  },
                  "eventname" : "Deploy"
               }
            ],
            "stack" : [
               ...skipped serialized contract representation...
            ],
            "gasconsumed" : "1006244000",
            "trigger" : "Application"
         }
      ]
   },
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
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json ccd533440d0317e9f366c50648d0013540e82741 main
```
... enter the password `qwerty`:
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```
Result:
```
Network fee: 0.0117652
System fee: 0.0717313
Total fee: 0.0834965
Relay transaction (y|N)> y
Sent invocation transaction bf92dbe258d9113f2a0684b6c782566b5b7bcda7524b349e085beb99147d8dc1
```
To check the counter value, call `getapplicaionlog` RPC-call for the invocation transaction:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["bf92dbe258d9113f2a0684b6c782566b5b7bcda7524b349e085beb99147d8dc1"] }' localhost:20331 | json_pp
```
The JSON result is:
```
{
   "id" : 1,
   "jsonrpc" : "2.0",
   "result" : {
      "txid" : "0xbf92dbe258d9113f2a0684b6c782566b5b7bcda7524b349e085beb99147d8dc1",
      "executions" : [
         {
            "trigger" : "Application",
            "gasconsumed" : "7173130",
            "notifications" : [
               {
                  "eventname" : "info",
                  "contract" : "0xccd533440d0317e9f366c50648d0013540e82741",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "value" : "VmFsdWUgcmVhZCBmcm9tIHN0b3JhZ2U=",
                           "type" : "ByteString"
                        }
                     ]
                  }
               },
               {
                  "state" : {
                     "value" : [
                        {
                           "value" : "U3RvcmFnZSBrZXkgYWxyZWFkeSBzZXQuIEluY3JlbWVudGluZyBieSAx",
                           "type" : "ByteString"
                        }
                     ],
                     "type" : "Array"
                  },
                  "eventname" : "info",
                  "contract" : "0xccd533440d0317e9f366c50648d0013540e82741"
               },
               {
                  "eventname" : "info",
                  "contract" : "0xccd533440d0317e9f366c50648d0013540e82741",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "TmV3IHZhbHVlIHdyaXR0ZW4gaW50byBzdG9yYWdl"
                        }
                     ]
                  }
               }
            ],
            "vmstate" : "HALT",
            "stack" : [
               {
                  "type" : "Integer",
                  "value" : "1"
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
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json ccd533440d0317e9f366c50648d0013540e82741 main
```
... enter the password `qwerty`:
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```
Result:
```
Network fee: 0.0117652
System fee: 0.0717313
Total fee: 0.0834965
Relay transaction (y|N)> y
Sent invocation transaction 4edf5fcceef8fddd6c145cafaad52fee6ebe83f166d28c1ad862349182f5550d
```
To check the counter value, call `getapplicaionlog` RPC-call for the invocation transaction:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["4edf5fcceef8fddd6c145cafaad52fee6ebe83f166d28c1ad862349182f5550d"] }' localhost:20331 | json_pp
```
The JSON result is:
```
{
   "id" : 1,
   "result" : {
      "txid" : "0x4edf5fcceef8fddd6c145cafaad52fee6ebe83f166d28c1ad862349182f5550d",
      "executions" : [
         {
            "vmstate" : "HALT",
            "notifications" : [
               {
                  "state" : {
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "VmFsdWUgcmVhZCBmcm9tIHN0b3JhZ2U="
                        }
                     ],
                     "type" : "Array"
                  },
                  "contract" : "0xccd533440d0317e9f366c50648d0013540e82741",
                  "eventname" : "info"
               },
               {
                  "eventname" : "info",
                  "contract" : "0xccd533440d0317e9f366c50648d0013540e82741",
                  "state" : {
                     "value" : [
                        {
                           "value" : "U3RvcmFnZSBrZXkgYWxyZWFkeSBzZXQuIEluY3JlbWVudGluZyBieSAx",
                           "type" : "ByteString"
                        }
                     ],
                     "type" : "Array"
                  }
               },
               {
                  "eventname" : "info",
                  "contract" : "0xccd533440d0317e9f366c50648d0013540e82741",
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
            "stack" : [
               {
                  "value" : "2",
                  "type" : "Integer"
               }
            ],
            "trigger" : "Application",
            "gasconsumed" : "7173130"
         }
      ]
   },
   "jsonrpc" : "2.0"
}
```

The `stack` field contains now `2` integer value, so the counter was incremented as we expected.

## Workshop. Part 3
In this part we'll know about NEP-17 token standard and try to write, deploy and invoke more complicated smart contract. 
Let’s go!

### NEP-17
[NEP-17](https://github.com/neo-project/proposals/blob/master/nep-17.mediawiki) is a token standard for the Neo blockchain that provides systems with a generalized interaction mechanism for tokenized smart contracts.
The example with implementation of all required by the standard methods you can find in [nep17.go](https://github.com/nspcc-dev/neo-go/blob/master/examples/token/nep17/nep17.go)
 
Let's take a view on the example of smart contract with NEP-17: [token.go](https://github.com/nspcc-dev/neo-go/blob/master/examples/token/token.go)
 
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
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```

Result:
```
Network fee: 0.0308452
System fee: 10.0107577
Total fee: 10.0416029
Relay transaction (y|N)> y
Sent invocation transaction ab6f934a5e2137d008613977b41b0a791e5497c2e97a2a84aed0bb684af2c5c3
Contract: c36534b6b81621178980438c18796f23a463441a
```   

Which means that our contract was deployed and now we can invoke it.

#### Step #2
Let's invoke the contract to perform different operations.

To start with, query `Symbol` of the created NEP-17 token:

```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json c36534b6b81621178980438c18796f23a463441a symbol
```                                                                   
Where
- `c36534b6b81621178980438c18796f23a463441a` is our contract hash from step #1
- `symbol` is operation string which was described earlier and returns token symbol

... and don't forget the password of your account `qwerty`.

Result:
```
Network fee: 0.0117852
System fee: 0.0141954
Total fee: 0.0259806
Relay transaction (y|N)> y
Sent invocation transaction 1e9d018b4ea8fc3442229ae437bc8451e876216dd04a9c071672753437c9ada5
```                                                                                         
Now, let's take a detailed look at this invocation transaction with `getapplicationlog` RPC call:

```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["1e9d018b4ea8fc3442229ae437bc8451e876216dd04a9c071672753437c9ada5"] }' localhost:20331 | json_pp
```               

Result:
```
{
   "id" : 1,
   "result" : {
      "txid" : "0x1e9d018b4ea8fc3442229ae437bc8451e876216dd04a9c071672753437c9ada5",
      "executions" : [
         {
            "vmstate" : "HALT",
            "stack" : [
               {
                  "value" : "QU5U",
                  "type" : "ByteString"
               }
            ],
            "notifications" : [],
            "trigger" : "Application",
            "gasconsumed" : "1419540"
         }
      ]
   },
   "jsonrpc" : "2.0"
}
```

At least, you can see that `stack` field of JSON result is not empty: it contains base64 byte array with the symbol of our token.

Following commands able you to get some additional information about token:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json c36534b6b81621178980438c18796f23a463441a decimals
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json c36534b6b81621178980438c18796f23a463441a totalSupply
```

#### Step #3

Now it's time for more interesting things. First of all, let's check the balance of NEP-17 token on our account by using `balanceOf`:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json c36534b6b81621178980438c18796f23a463441a balanceOf NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB
```                             
... with `qwerty` password. The result is:
```
Network fee: 0.0120452
System fee: 0.0249927
Total fee: 0.0370379
Relay transaction (y|N)> y
Sent invocation transaction 2f9d830cfc52747c6d0658aeb25c75e334afd1e4badd1dc946c706a22a7d1e10
```
And take a closer look at the transaction's details with `getapplicationlog` RPC-call:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["2f9d830cfc52747c6d0658aeb25c75e334afd1e4badd1dc946c706a22a7d1e10"] }' localhost:20331 | json_pp
```
Result:
```
{
   "id" : 1,
   "jsonrpc" : "2.0",
   "result" : {
      "txid" : "0x2f9d830cfc52747c6d0658aeb25c75e334afd1e4badd1dc946c706a22a7d1e10",
      "executions" : [
         {
            "vmstate" : "HALT",
            "stack" : [
               {
                  "type" : "Integer",
                  "value" : "0"
               }
            ],
            "notifications" : [],
            "gasconsumed" : "2499270",
            "trigger" : "Application"
         }
      ]
   }
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
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json c36534b6b81621178980438c18796f23a463441a mint NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB -- NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB:CalledByEntry
```
Where:
- `--` is a special delimiter of transaction's cosigners list
- `NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB` is the signer itself (which is our account)

... with `qwerty` pass. The result:
``` 
Network fee: 0.0119952
System fee: 0.1371123
Total fee: 0.1491075
Relay transaction (y|N)> y
Sent invocation transaction 9a54e07e54550e57ab9d7d1a1a001516ae2514bae23f6691632f2a3bc1c2d8b7
```
`getapplicationlog` RPC-call for this transaction tells us the following:
```
{
   "jsonrpc" : "2.0",
   "result" : {
      "txid" : "0x9a54e07e54550e57ab9d7d1a1a001516ae2514bae23f6691632f2a3bc1c2d8b7",
      "executions" : [
         {
            "stack" : [
               {
                  "value" : true,
                  "type" : "Boolean"
               }
            ],
            "notifications" : [
               {
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "type" : "Any"
                        },
                        {
                           "value" : "rtjbO11mTKBL3XvtN5kv+VhWC0E=",
                           "type" : "ByteString"
                        },
                        {
                           "type" : "Integer",
                           "value" : "1100000000000000"
                        }
                     ]
                  },
                  "contract" : "0xc36534b6b81621178980438c18796f23a463441a",
                  "eventname" : "Transfer"
               }
            ],
            "trigger" : "Application",
            "vmstate" : "HALT",
            "gasconsumed" : "13711230"
         }
      ]
   },
   "id" : 1
}
```
Here we have `true` at the `stack` field, which means that token was successfully minted.
Let's just ensure that by querying `balanceOf` one more time:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json c36534b6b81621178980438c18796f23a463441a balanceOf NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB
```
... with `qwerty` pass. The result:
``` 
Network fee: 0.0120452
System fee: 0.0274533
Total fee: 0.0394985
Relay transaction (y|N)> y
Sent invocation transaction 870607e5bbffdaef9adb38cf4ca08125481554bf674d6a63e79c3779c924017c
```
... with the following `getapplicationlog` JSON message:
```
{
   "id" : 1,
   "result" : {
      "executions" : [
         {
            "trigger" : "Application",
            "gasconsumed" : "2745330",
            "stack" : [
               {
                  "type" : "Integer",
                  "value" : "1100000000000000"
               }
            ],
            "vmstate" : "HALT",
            "notifications" : []
         }
      ],
      "txid" : "0x870607e5bbffdaef9adb38cf4ca08125481554bf674d6a63e79c3779c924017c"
   },
   "jsonrpc" : "2.0"
}
```
Now we can see integer value at the `stack` field, so `1100000000000000` is the NEP-17 token balance of our account.

Note, that token can be minted only once.

#### Step #5

After we are done with minting, it's possible to transfer token to someone else.
Let's transfer 5 tokens from our account to `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm` with `transfer` call:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json c36534b6b81621178980438c18796f23a463441a transfer NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm 500000000 null -- NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB:CalledByEntry
```
... with password `qwerty` and following result:
``` 
Network fee: 0.0123652
System fee: 0.1188695
Total fee: 0.1312347
Relay transaction (y|N)> y
Sent invocation transaction a796fc3d5b75f6c5289aca9d2f77d6e50c5d9fdbd860068fb5771b99ff747e96
```
Our favourite `getapplicationlog` RPC-call tells us:
```
{
   "id" : 1,
   "jsonrpc" : "2.0",
   "result" : {
      "txid" : "0xa796fc3d5b75f6c5289aca9d2f77d6e50c5d9fdbd860068fb5771b99ff747e96",
      "executions" : [
         {
            "stack" : [
               {
                  "type" : "Boolean",
                  "value" : true
               }
            ],
            "gasconsumed" : "11886950",
            "vmstate" : "HALT",
            "notifications" : [
               {
                  "eventname" : "Transfer",
                  "contract" : "0xc36534b6b81621178980438c18796f23a463441a",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "rtjbO11mTKBL3XvtN5kv+VhWC0E="
                        },
                        {
                           "type" : "ByteString",
                           "value" : "50l6vFaauRKm8hPVkr3Aw2CeHQs="
                        },
                        {
                           "type" : "Integer",
                           "value" : "500000000"
                        }
                     ]
                  }
               }
            ],
            "trigger" : "Application"
         }
      ]
   }
}
```
Note, that `stack` field contains `true`, which means that token was successfully transferred.
Let's now check the balance of `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm` account to ensure that the amount of token on that account = 5:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json c36534b6b81621178980438c18796f23a463441a balanceOf NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm
```
The `getapplicationlog` RPC-call for this transaction tells us the following:
```
{
   "id" : 1,
   "jsonrpc" : "2.0",
   "result" : {
      "txid" : "0x40f3f3c12d3eeba7e282bbaf76af944310b82504dcb3c09db3ea6c6d8418bb6b",
      "executions" : [
         {
            "gasconsumed" : "2745330",
            "trigger" : "Application",
            "vmstate" : "HALT",
            "stack" : [
               {
                  "type" : "Integer",
                  "value" : "500000000"
               }
            ],
            "notifications" : []
         }
      ]
   }
}
```
Here we are! There are exactly 5 tokens at the `stack` field. You can also ensure that these 5 tokens were debited from `NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB` account by using `balanceOf` method.

## Workshop. Part 4
In this part we'll summarise our knowledge about smart contracts by investigating [4-domain](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/4-domain/4-domain.go) smart contract. This contract
contains code for domain registration, transferring, deletion and getting information about registered domains.

Let’s go!

#### Step #1
Let's take a glance at our [contract](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/4-domain/4-domain.go) and inspect it. The contract takes an action string as the first parameter, which is one of the following:
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

Compile smart contract [4-domain.go](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/4-domain/4-domain.go) with [configuration](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/4-domain/4-domain.yml)
```
$ ./bin/neo-go contract compile -i 4-domain/4-domain.go -c 4-domain/4-domain.yml -m 4-domain/4-domain.manifest.json
```

... and deploy it:
```
$ ./bin/neo-go contract deploy -i 4-domain/4-domain.nef --manifest 4-domain/4-domain.manifest.json -r http://localhost:20331 -w my_wallet.json
```
... enter the password `qwerty`:
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```

Result:
```
Network fee: 0.0303652
System fee: 10.0107577
Total fee: 10.0411229
Relay transaction (y|N)> y
Sent invocation transaction 69548dfecf70c190e2bc872aa210d53ff7faa7956074154c90a27d4c94420562
Contract: 9042814f07d65d2b835fa1f07d21c22c6e1cbdf7
```   
You know, what it means :)

#### Step #3

Invoke the contract to register domain with name `my_first_domain`:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 9042814f07d65d2b835fa1f07d21c22c6e1cbdf7 register my_first_domain NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB -- NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB:CalledByEntry
```
... the strongest password in the world, guess: `qwerty`
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```
Result:
```
Network fee: 0.0122052
System fee: 0.0894353
Total fee: 0.1016405
Relay transaction (y|N)> y
Sent invocation transaction d9f05af09ee497ceb523be3274483143f83f7f24038f37799a962c8dee640357
```
Also you can see the log message in the console, where you run neo-go node:
```
2022-07-05T13:42:36.592+0300	INFO	runtime log	{"tx": "d9f05af09ee497ceb523be3274483143f83f7f24038f37799a962c8dee640357", "script": "9042814f07d65d2b835fa1f07d21c22c6e1cbdf7", "msg": "RegisterDomain: my_first_domain"}
```
Well, that's ok. Let's check now, whether our domain was registered with `getapplicationlog` RPC-call:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["d9f05af09ee497ceb523be3274483143f83f7f24038f37799a962c8dee640357"] }' localhost:20331 | json_pp
```
The result is:
```
{
   "jsonrpc" : "2.0",
   "result" : {
      "executions" : [
         {
            "stack" : [
               {
                  "type" : "Boolean",
                  "value" : true
               }
            ],
            "vmstate" : "HALT",
            "notifications" : [
               {
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "value" : "rtjbO11mTKBL3XvtN5kv+VhWC0E=",
                           "type" : "ByteString"
                        },
                        {
                           "value" : "bXlfZmlyc3RfZG9tYWlu",
                           "type" : "ByteString"
                        }
                     ]
                  },
                  "contract" : "0x9042814f07d65d2b835fa1f07d21c22c6e1cbdf7",
                  "eventname" : "registered"
               }
            ],
            "trigger" : "Application",
            "gasconsumed" : "8943530"
         }
      ],
      "txid" : "0xd9f05af09ee497ceb523be3274483143f83f7f24038f37799a962c8dee640357"
   },
   "id" : 1
}
```
Especially, we're interested in two fields of the json:

First one is `notifications` field, which contains one notification with `registered` name:
- `bXlfZmlyc3RfZG9tYWlu` byte string in base64 representation, which can be decoded to `my_first_domain` - our domain's name
- `ecv/0NH0e0cStm0wWBgjCxMyaok=` byte array, which can be decoded to the account address `NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB`.

The second field is `stack` with `true` value, which was returned by the smart contract.

All of these values let us be sure that our domain was successfully registered.  

#### Step #4

Invoke the contract to query the address information our `my_first_domain` domain:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 9042814f07d65d2b835fa1f07d21c22c6e1cbdf7 query my_first_domain
```
... the pass `qwerty`:
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```
Result:
```
Network fee: 0.0119552
System fee: 0.0412161
Total fee: 0.0531713
Relay transaction (y|N)> y
Sent invocation transaction 4d1b15e3c891b4e5efe6f093e54c1090476dc0ed0069b228f58578c9360ee1f2
```
and log-message:
```
2022-07-05T13:44:21.681+0300	INFO	runtime log	{"tx": "4d1b15e3c891b4e5efe6f093e54c1090476dc0ed0069b228f58578c9360ee1f2", "script": "9042814f07d65d2b835fa1f07d21c22c6e1cbdf7", "msg": "QueryDomain: my_first_domain"}
```
Let's check this transaction with `getapplicationlog` RPC call:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["4d1b15e3c891b4e5efe6f093e54c1090476dc0ed0069b228f58578c9360ee1f2"] }' localhost:20331 | json_pp
```
... which gives us the following result:
```
{
   "result" : {
      "executions" : [
         {
            "stack" : [
               {
                  "value" : "rtjbO11mTKBL3XvtN5kv+VhWC0E=",
                  "type" : "ByteString"
               }
            ],
            "trigger" : "Application",
            "gasconsumed" : "4121610",
            "notifications" : [],
            "vmstate" : "HALT"
         }
      ],
      "txid" : "0x4d1b15e3c891b4e5efe6f093e54c1090476dc0ed0069b228f58578c9360ee1f2"
   },
   "id" : 1,
   "jsonrpc" : "2.0"
}
```

with base64 interpretation of our account address `NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB` on the stack, which means that domain `my_first_domain` was registered by the owner with received account address.

#### Step #5

Invoke the contract to transfer domain to the other account (e.g. account with `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm` address):
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 9042814f07d65d2b835fa1f07d21c22c6e1cbdf7 transfer my_first_domain NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm -- NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB
```
... the password: `qwerty`
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```
Result:
```
Network fee: 0.0122052
System fee: 0.0748064
Total fee: 0.0870116
Relay transaction (y|N)> y
Sent invocation transaction 937ec7539a31246ff88eb0bdda74cf7d9613d4a8ad1b7b33f0a785e458d76a14
```
and log-message:
```
2022-07-05T13:46:06.746+0300	INFO	runtime log	{"tx": "937ec7539a31246ff88eb0bdda74cf7d9613d4a8ad1b7b33f0a785e458d76a14", "script": "9042814f07d65d2b835fa1f07d21c22c6e1cbdf7", "msg": "TransferDomain: my_first_domain"}
```
Perfect. And `getapplicationlog` RPC-call...
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["937ec7539a31246ff88eb0bdda74cf7d9613d4a8ad1b7b33f0a785e458d76a14"] }' localhost:20331 | json_pp
```
... tells us:
```
{
   "jsonrpc" : "2.0",
   "result" : {
      "executions" : [
         {
            "vmstate" : "HALT",
            "gasconsumed" : "7480640",
            "trigger" : "Application",
            "stack" : [
               {
                  "value" : true,
                  "type" : "Boolean"
               }
            ],
            "notifications" : [
               {
                  "contract" : "0x9042814f07d65d2b835fa1f07d21c22c6e1cbdf7",
                  "eventname" : "deleted",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "rtjbO11mTKBL3XvtN5kv+VhWC0E="
                        },
                        {
                           "type" : "ByteString",
                           "value" : "bXlfZmlyc3RfZG9tYWlu"
                        }
                     ]
                  }
               },
               {
                  "eventname" : "registered",
                  "state" : {
                     "value" : [
                        {
                           "value" : "50l6vFaauRKm8hPVkr3Aw2CeHQs=",
                           "type" : "ByteString"
                        },
                        {
                           "value" : "bXlfZmlyc3RfZG9tYWlu",
                           "type" : "ByteString"
                        }
                     ],
                     "type" : "Array"
                  },
                  "contract" : "0x9042814f07d65d2b835fa1f07d21c22c6e1cbdf7"
               }
            ]
         }
      ],
      "txid" : "0x937ec7539a31246ff88eb0bdda74cf7d9613d4a8ad1b7b33f0a785e458d76a14"
   },
   "id" : 1
}
```
The `notifications` field contains two events:
- First one with name `deleted` and additional information (domain `my_first_domain` was deleted from account `NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB`),
- Second one with name `registered`  and additional information  (domain `my_first_domain` was registered with account `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm`).
The `stack` field contains `true` value, which means that our domain was successfully transferred.

#### Step #6

The last call is `delete`, so you can try to create the other domain, e.g. `my_second_domain` and then remove it from storage with:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 9042814f07d65d2b835fa1f07d21c22c6e1cbdf7 register my_second_domain NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB -- NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB:CalledByEntry
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 9042814f07d65d2b835fa1f07d21c22c6e1cbdf7 delete my_second_domain -- NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB
```

## Workshop. Part 5

In this part we'll take a look at the example of a decentralized application for the
Neo ecosystem written in Go. The example demonstrates how to use NeoGo RPC Client
and a set of helpers to work with wallets, deployments, invocations, data retrieval,
blockchain events, auto-generated smart contract RPC bindings and such.

#### Step #1

Ensure that private network created at the [Preparation part](#workshop-preparation)
step is up and running. We'll also need a `NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB`
account from the [wallet](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/e11949bf5f1dc1ce4e3b6551b6ae22032945c75d/my_wallet.json)
with some GAS on it. Take a look at the dApp example: [dApp.go](./dApp/dApp.go).
The example includes all basic API usages you need to be aware of to start your
dApp development for Neo:
* [JSON-RPC client](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient) creation, initialization and example of simple Neo JSON RPC APIs usage.
* Neo wallet and account managing with [NeoGo `wallet` package](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/wallet).
* [Extensions](https://github.com/nspcc-dev/neo-go/blob/5fc61be5f6c5349d8de8b61967380feee6b51c55/docs/rpc.md#extensions) offered by NeoGo JSON RPC server and supported by the NeoGo RPC client.
* NeoGo JSON RPC server [web-socket extension](https://github.com/nspcc-dev/neo-go/blob/5fc61be5f6c5349d8de8b61967380feee6b51c55/docs/rpc.md#websocket-server). [Web-socket JSON RPC Client](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient#WSClient) creation, initialization and usage example. 
* [`unwrap` helper package](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/unwrap) for test invocation results unwrapping.
* NeoGo JSON RPC [Notification subsystem](https://github.com/nspcc-dev/neo-go/blob/5fc61be5f6c5349d8de8b61967380feee6b51c55/docs/notifications.md) usage example. Subscribe for a various set of blockchain node events and receive notifications over web-socket channel.
* NeoGo `invoker` package usage example. Perform test invocation of smart contract, script or verification method with the powerful [Invoker API](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/invoker).
* How to work with witnesses and [witness scopes](https://neospcc.medium.com/thou-shalt-check-their-witnesses-485d2bf8375d).
* [Historic invocations functionality](https://github.com/nspcc-dev/neo-go/blob/5fc61be5f6c5349d8de8b61967380feee6b51c55/docs/rpc.md#invokecontractverifyhistoric-invokefunctionhistoric-and-invokescripthistoric-calls) provided by NeoGo RPC server and client extensions.
* NeoGo `actor` package usage example. Build, tune, sign, send and await transactions with our flexible [Actor API](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/actor).
* Work with a set of NEP-specific and native-specific packages:
  * [`nep17`](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/nep17) and [`nep11`](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/nep11) packages for NEP-17 or NEP-11 compatible Neo token management.
  * Native contracts specific actors: [`gas`](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/gas), [`neo`](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/neo), [`management`](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/management), [`policy`](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/policy), [`oracle`](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/oracle), [`rolemgmt`](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/rolemgmt), [`notary`](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/notary) Actor packages.
* Deploy and invoke an [example contract](https://github.com/nspcc-dev/neo-go/tree/5fc61be5f6c5349d8de8b61967380feee6b51c55/examples/storage) with the help of native [`ContractManagement` specific Actor package](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/management).
* Contract storage iterator API demonstration:
  * Script-based iterator unwrapping via [`CallAndExpandIterator` Invoker API](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/invoker#Invoker.CallAndExpandIterator).
  * Session-based iterator traversal via [`TraverseIterator` Invoker API](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/invoker#Invoker.TraverseIterator).
* Usage of autogenerated [RPC smart contract binding](https://github.com/nspcc-dev/neo-go/blob/5fc61be5f6c5349d8de8b61967380feee6b51c55/docs/compiler.md#generating-rpc-contract-bindings).

Carefully read the dApp example and comments to the blocks of code. Follow the
provided links and read the corresponding documentation. Before running the
example, clone the NeoGo repo to the same folder where the workshop repo is
placed and compile the example Storage contract:
```
$ git clone https://github.com/nspcc-dev/neo-go
$ cd neo-go
$ make build
$ ./bin/neo-go contract compile -i ./examples/storage/storage.go -c ./examples/storage/storage.yml -o examples/storage/storage.nef -m ./examples/storage/storage.manifest.json
```  

The last preparation step is to check that variable `transferTxH` contains the
actual hash of GAS transfer transaction from the Preparation part.

#### Step #2

Change working directory to the `dApp` package and run the `dApp.go` example:
```
$ cd ./dApp
$ go run dApp.go 
``` 

Inspect the resulting output in the console, investigate the dApp example code
snippets and build your own dApp for Neo!

Thank you!

### Useful links

* [Our basic tutorial on Medium](https://medium.com/@neospcc/%D1%81%D0%BC%D0%B0%D1%80%D1%82-%D0%BA%D0%BE%D0%BD%D1%82%D1%80%D0%B0%D0%BA%D1%82-%D0%B4%D0%BB%D1%8F-neo-769139352b65)
* [Go smart contracts development workshop on YouTube (as a part of the Neo Polaris Launchpad)](https://www.youtube.com/watch?v=o38fXiLG7EM)
* [Go smart contracts development workshop on YouTube (as a part of the Neo Asia-Pacific Tour)](https://www.youtube.com/watch?v=q_TMlpx1-0M)
* [Go dApp development workshop on YouTube](https://www.youtube.com/live/8zVBIrVQa58)
* [Using Neo Blockchain Toolkit](https://medium.com/@neospcc/neogo-adds-support-for-neo-blockchain-toolkit-673ea914f661)
* [Neo documentation](https://docs.neo.org/)
* [Neo github](https://github.com/neo-project/neo/)
* [NeoGo github](https://github.com/nspcc-dev/neo-go)
