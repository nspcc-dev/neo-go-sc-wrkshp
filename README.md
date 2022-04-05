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
Neo has a nice monitor where you can find particular node running in the blockchain network.
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
    Password >
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
    Enter password to unlock wallet and sign the transaction
    Password >
    ```
    The result is transaction signed by both first and second nodes.

3. Sign the transaction using the third node address and push it to the chain:
    ```
    $ ./bin/neo-go wallet sign -w ./.docker/wallets/wallet3.json --in my_tx2.json --out my_tx3.json --address NVTiAjNgagDkTr5HTzDmQP9kPwPHN5BgVq -r http://localhost:20331
    ```
    Enter the password `three`:
    ```
    Enter password to unlock wallet and sign the transaction
    Password >
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
   "jsonrpc" : "2.0",
   "result" : {
      "address" : "NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB",
      "sent" : [],
      "received" : [
         {
            "amount" : "2999999900000000",
            "txhash" : "0xb0d0cb55fe68fef89b071d4dfdbd19974250b10a8a257f50dd568f76c4886d30",
            "assethash" : "0xd2a4cff31913016155e38e474a2c06d08be276cf",
            "transfernotifyindex" : 0,
            "transferaddress" : "NVTiAjNgagDkTr5HTzDmQP9kPwPHN5BgVq",
            "blockindex" : 49,
            "timestamp" : 1638194279180
         }
      ]
   },
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
Network fee: 1515520
System fee: 1001045530
Total fee: 1002561050
Relay transaction (y|N)>
```
Result:
```
Sent invocation transaction 28b26283ea2689dc5abf30bf6f0605b3819089f7fbf07bc26e41d62e1a9f5841
Contract: a48467c9bf559524575cf0d3b25cd97e67b01bc5
```

At this point your ‘Hello World’ contract is deployed and could be invoked. Let’s do it as a final step.

#### Step 4
Invoke contract.
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json a48467c9bf559524575cf0d3b25cd97e67b01bc5 main
```

Where
- `contract invokefunction` runs invoke with provided parameters
- `-r http://localhost:20331` defines RPC endpoint used for function call
- `-w my_wallet.json` is a wallet
- `a48467c9bf559524575cf0d3b25cd97e67b01bc5` contract hash got as an output from the previous command (deployment in step 6)
- `main` - method to be called

Enter password `qwerty` for account:
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```

Result:
```
Sent invocation transaction bd23c836f7bdd62a0d9c5ecb3f5bdbf2d38ec9a5e2e3935ca543d8c18ed5479d
```
In the console where you were running step #5 you will get:
```
2021-11-29T17:02:44.395+0300	INFO	runtime log	{"tx": "bd23c836f7bdd62a0d9c5ecb3f5bdbf2d38ec9a5e2e3935ca543d8c18ed5479d", "script": "a48467c9bf559524575cf0d3b25cd97e67b01bc5", "msg": "Hello, world!"}
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
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getrawtransaction", "params": ["28b26283ea2689dc5abf30bf6f0605b3819089f7fbf07bc26e41d62e1a9f5841", 1] }' localhost:20331 | json_pp
```

Where
- `"jsonrpc": "2.0"` is protocol version
- `"id": 1` is id of current request
- `"method": "getrawtransaction"` is requested method
- `"params": ["28b26283ea2689dc5abf30bf6f0605b3819089f7fbf07bc26e41d62e1a9f5841", 1]` is an array of parameters, 
  where
   - `28b26283ea2689dc5abf30bf6f0605b3819089f7fbf07bc26e41d62e1a9f5841` is deployment transaction hash
   - `1` is `verbose` parameter for detailed JSON string output
- `json_pp` just makes the JSON output prettier

Result:
```
{
   "result" : {
      "vmstate" : "HALT",
      "signers" : [
         {
            "account" : "0x410b5658f92f9937ed7bdd4ba04c665d3bdbd8ae",
            "scopes" : "CalledByEntry"
         }
      ],
      "nonce" : 2714712230,
      "sender" : "NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB",
      "confirmations" : 31,
      "validuntilblock" : 64,
      "netfee" : "1515520",
      "size" : 532,
      "version" : 0,
      "hash" : "0x28b26283ea2689dc5abf30bf6f0605b3819089f7fbf07bc26e41d62e1a9f5841",
      "sysfee" : "1001045530",
      "script" : "DOZ7Im5hbWUiOiJIZWxsb1dvcmxkIGNvbnRyYWN0IiwiYWJpIjp7Im1ldGhvZHMiOlt7Im5hbWUiOiJtYWluIiwib2Zmc2V0IjowLCJwYXJhbWV0ZXJzIjpbXSwicmV0dXJudHlwZSI6IlZvaWQiLCJzYWZlIjpmYWxzZX1dLCJldmVudHMiOltdfSwiZmVhdHVyZXMiOnt9LCJncm91cHMiOltdLCJwZXJtaXNzaW9ucyI6W10sInN1cHBvcnRlZHN0YW5kYXJkcyI6W10sInRydXN0cyI6W10sImV4dHJhIjpudWxsfQxkTkVGM25lby1nby0wLjk3LjMAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABYMDUhlbGxvLCB3b3JsZCFBz+dHliFAUFIoRhLAHwwGZGVwbG95DBT9o/pDRupTKiWPxJfdrdtkN8n9/0FifVtS",
      "attributes" : [],
      "witnesses" : [
         {
            "invocation" : "DED+3Mj7PjuEIyO7zMSQBevEnaJi/Z+XKBKLEyooMdQsmZBPYcj/L+nktVhXu63Vw8ynTna1RROV9wKOXuGfJCz/",
            "verification" : "DCEDhEhWuuSSNuCc7nLsxQhI8nFlt+UfY3oP0/UkYmdH7G5BVuezJw=="
         }
      ],
      "blocktime" : 1638194489321,
      "blockhash" : "0x5bfb4b1adf2ecab1e6c6bd49bf807350a21f481fed2073d30959b56286e6ab0d"
   },
   "id" : 1,
   "jsonrpc" : "2.0"
}
```

#### GetApplicationLog
[GetApplicationLog](https://docs.neo.org/v3/docs/en-us/reference/rpc/latest-version/api/getapplicationlog.html) returns the contract log based on the specified transaction id.

Request application log for invocation transaction from step #4:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["bd23c836f7bdd62a0d9c5ecb3f5bdbf2d38ec9a5e2e3935ca543d8c18ed5479d"] }' localhost:20331 | json_pp
```

With a single parameter:
- `bd23c836f7bdd62a0d9c5ecb3f5bdbf2d38ec9a5e2e3935ca543d8c18ed5479d` - invocation transaction hash from step #7

Result:
```
{
   "jsonrpc" : "2.0",
   "id" : 1,
   "result" : {
      "executions" : [
         {
            "vmstate" : "HALT",
            "notifications" : [],
            "trigger" : "Application",
            "gasconsumed" : "2028330",
            "stack" : [
               {
                  "type" : "Any"
               }
            ]
         }
      ],
      "txid" : "0xbd23c836f7bdd62a0d9c5ecb3f5bdbf2d38ec9a5e2e3935ca543d8c18ed5479d"
   }
}
```

#### Other Useful RPC calls
```
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getversion", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getblockcount", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getconnectioncount", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getcontractstate", "params": ["0xa48467c9bf559524575cf0d3b25cd97e67b01bc5"] }' localhost:20331
```

List of supported by neo-go node RPC commands you can find [here](https://github.com/nspcc-dev/neo-go/blob/master/docs/rpc.md#supported-methods).

#### Utilities

neo-go CLI provides `query tx` utility to check the transaction state. It uses
`getrawtransaction` and `getapplicationlog` RPC calls under the hood and prints details
of transaction invocation. Use `query tx` command to ensure transaction was accepted
to the chain:
```
./bin/neo-go query tx bd23c836f7bdd62a0d9c5ecb3f5bdbf2d38ec9a5e2e3935ca543d8c18ed5479d -r http://localhost:20331 -v
```
where
- `bd23c836f7bdd62a0d9c5ecb3f5bdbf2d38ec9a5e2e3935ca543d8c18ed5479d` - invocation transaction hash from step #7
- `-r http://localhost:20331` - RPC node endpoint
- `-v` - verbose flag (enables transaction's signers, fees and script dumps)

The result is:
```
Hash:			bd23c836f7bdd62a0d9c5ecb3f5bdbf2d38ec9a5e2e3935ca543d8c18ed5479d
OnChain:		true
BlockHash:		c72e82e1dc4274a1c6d587370bfe56f359968ce18c27a7988883d34ebf415496
Success:		true
Signer:			NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB (None)
SystemFee:		0.0202833 GAS
NetworkFee:		0.0117752 GAS
Script:			EMAfDARtYWluDBTFG7BnftlcstPwXFcklVW/yWeEpEFifVtS
INDEX    OPCODE       PARAMETER                                   
0        PUSH0                                                    <<
1        PACK                                                     
2        PUSH15                                                   
3        PUSHDATA1    6d61696e ("main")                           
9        PUSHDATA1    c51bb0677ed95cb2d3f05c57249555bfc96784a4    
31       SYSCALL      System.Contract.Call (627d5b52)             

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
Sent invocation transaction a0722f99edb590b789cee2589e74a09f93f36efeb06b8f5da7abde85c789a2d3
Contract: aa9c0d6006eccb53ee76688722898617606a88aa
```   

Which means that our contract was deployed and now we can invoke it.

Let's check that the storage value was initialised with `0`. Use `getapplicaionlog` RPC-call for the deployment transaction:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["a0722f99edb590b789cee2589e74a09f93f36efeb06b8f5da7abde85c789a2d3"] }' localhost:20331 | json_pp
```

The JSON result is:
```
{
   "jsonrpc" : "2.0",
   "result" : {
      "txid" : "0xa0722f99edb590b789cee2589e74a09f93f36efeb06b8f5da7abde85c789a2d3",
      "executions" : [
         {
            "vmstate" : "HALT",
            "notifications" : [
               {
                  "contract" : "0xaa9c0d6006eccb53ee76688722898617606a88aa",
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
                  "contract" : "0xaa9c0d6006eccb53ee76688722898617606a88aa",
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
                           "value" : "qohqYBeGiSKHaHbuU8vsBmANnKo=",
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
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json aa9c0d6006eccb53ee76688722898617606a88aa main
```
... enter the password `qwerty`:
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```
Result:
```
Sent invocation transaction a58140ee3ebee1f4fb844311b73ac86454d458122eec9c4cea19725a106a260f
```
To check the counter value, call `getapplicaionlog` RPC-call for the invocation transaction:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["a58140ee3ebee1f4fb844311b73ac86454d458122eec9c4cea19725a106a260f"] }' localhost:20331 | json_pp
```
The JSON result is:
```
{
   "jsonrpc" : "2.0",
   "id" : 1,
   "result" : {
      "txid" : "0xa58140ee3ebee1f4fb844311b73ac86454d458122eec9c4cea19725a106a260f",
      "executions" : [
         {
            "notifications" : [
               {
                  "eventname" : "info",
                  "state" : {
                     "value" : [
                        {
                           "type" : "Buffer",
                           "value" : "VmFsdWUgcmVhZCBmcm9tIHN0b3JhZ2U="
                        }
                     ],
                     "type" : "Array"
                  },
                  "contract" : "0xaa9c0d6006eccb53ee76688722898617606a88aa"
               },
               {
                  "contract" : "0xaa9c0d6006eccb53ee76688722898617606a88aa",
                  "state" : {
                     "value" : [
                        {
                           "type" : "Buffer",
                           "value" : "U3RvcmFnZSBrZXkgYWxyZWFkeSBzZXQuIEluY3JlbWVudGluZyBieSAx"
                        }
                     ],
                     "type" : "Array"
                  },
                  "eventname" : "info"
               },
               {
                  "contract" : "0xaa9c0d6006eccb53ee76688722898617606a88aa",
                  "state" : {
                     "value" : [
                        {
                           "type" : "Buffer",
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
                  "value" : "1",
                  "type" : "Integer"
               }
            ],
            "trigger" : "Application",
            "gasconsumed" : "7233580",
            "vmstate" : "HALT"
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
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json aa9c0d6006eccb53ee76688722898617606a88aa main
```
... enter the password `qwerty`:
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```
Result:
```
Sent invocation transaction 157ca5e5b8cf8f84c9660502a3270b346011612bded1514a6847f877c433a9bb
```
To check the counter value, call `getapplicaionlog` RPC-call for the invocation transaction:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["157ca5e5b8cf8f84c9660502a3270b346011612bded1514a6847f877c433a9bb"] }' localhost:20331 | json_pp
```
The JSON result is:
```
{
   "jsonrpc" : "2.0",
   "result" : {
      "txid" : "0x157ca5e5b8cf8f84c9660502a3270b346011612bded1514a6847f877c433a9bb",
      "executions" : [
         {
            "gasconsumed" : "7233580",
            "notifications" : [
               {
                  "state" : {
                     "value" : [
                        {
                           "type" : "Buffer",
                           "value" : "VmFsdWUgcmVhZCBmcm9tIHN0b3JhZ2U="
                        }
                     ],
                     "type" : "Array"
                  },
                  "eventname" : "info",
                  "contract" : "0xaa9c0d6006eccb53ee76688722898617606a88aa"
               },
               {
                  "contract" : "0xaa9c0d6006eccb53ee76688722898617606a88aa",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "value" : "U3RvcmFnZSBrZXkgYWxyZWFkeSBzZXQuIEluY3JlbWVudGluZyBieSAx",
                           "type" : "Buffer"
                        }
                     ]
                  },
                  "eventname" : "info"
               },
               {
                  "contract" : "0xaa9c0d6006eccb53ee76688722898617606a88aa",
                  "state" : {
                     "value" : [
                        {
                           "value" : "TmV3IHZhbHVlIHdyaXR0ZW4gaW50byBzdG9yYWdl",
                           "type" : "Buffer"
                        }
                     ],
                     "type" : "Array"
                  },
                  "eventname" : "info"
               }
            ],
            "trigger" : "Application",
            "stack" : [
               {
                  "value" : "2",
                  "type" : "Integer"
               }
            ],
            "vmstate" : "HALT"
         }
      ]
   },
   "id" : 1
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
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```

Result:
```
Sent invocation transaction 7af616aacc798760274a449700f14e4e25d5c3b262d200303dc701f8ea41707c
Contract: 27502a01e2fb013e1e4c428abb7b360df9f3f0cb
```   

Which means that our contract was deployed and now we can invoke it.

#### Step #2
Let's invoke the contract to perform different operations.

To start with, query `Symbol` of the created nep17 token:

```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 27502a01e2fb013e1e4c428abb7b360df9f3f0cb symbol
```                                                                   
Where
- `27502a01e2fb013e1e4c428abb7b360df9f3f0cb` is our contract hash from step #1
- `symbol` is operation string which was described earlier and returns token symbol

... and don't forget the password of your account `qwerty`.

Result:
```
Sent invocation transaction 535bccc585698c531cc58677b116ea7c567604194bf3202c6be7ac4d420b85af
```                                                                                         
Now, let's take a detailed look at this invocation transaction with `getapplicationlog` RPC call:

```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["535bccc585698c531cc58677b116ea7c567604194bf3202c6be7ac4d420b85af"] }' localhost:20331 | json_pp
```               

Result:
```
{
   "result" : {
      "txid" : "0x535bccc585698c531cc58677b116ea7c567604194bf3202c6be7ac4d420b85af",
      "executions" : [
         {
            "notifications" : [],
            "trigger" : "Application",
            "gasconsumed" : "4292370",
            "vmstate" : "HALT",
            "stack" : [
               {
                  "type" : "ByteString",
                  "value" : "QU5U"
               }
            ]
         }
      ]
   },
   "jsonrpc" : "2.0",
   "id" : 1
}
```

At least, you can see that `stack` field of JSON result is not empty: it contains base64 byte array with the symbol of our token.

Following commands able you to get some additional information about token:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 27502a01e2fb013e1e4c428abb7b360df9f3f0cb decimals
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 27502a01e2fb013e1e4c428abb7b360df9f3f0cb totalSupply
```

#### Step #3

Now it's time for more interesting things. First of all, let's check the balance of nep17 token on our account by using `balanceOf`:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 27502a01e2fb013e1e4c428abb7b360df9f3f0cb balanceOf NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB
```                             
... with `qwerty` password. The result is:
```
Sent invocation transaction e15871ad735a216a5e55f86dbb31ed4b4e928f4531f2788e547cc881e8532a8a
```
And take a closer look at the transaction's details with `getapplicationlog` RPC-call:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["e15871ad735a216a5e55f86dbb31ed4b4e928f4531f2788e547cc881e8532a8a"] }' localhost:20331 | json_pp
```
Result:
```
{
   "id" : 1,
   "result" : {
      "txid" : "0xe15871ad735a216a5e55f86dbb31ed4b4e928f4531f2788e547cc881e8532a8a",
      "executions" : [
         {
            "stack" : [
               {
                  "value" : "0",
                  "type" : "Integer"
               }
            ],
            "trigger" : "Application",
            "gasconsumed" : "5311140",
            "vmstate" : "HALT",
            "notifications" : []
         }
      ]
   },
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
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 27502a01e2fb013e1e4c428abb7b360df9f3f0cb mint NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB -- NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB:CalledByEntry
```
Where:
- `--` is a special delimiter of transaction's cosigners list
- `NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB` is the signer itself (which is our account)

... with `qwerty` pass. The result:
``` 
Sent invocation transaction 296cb753f9afeef7ace3690eed05c08336329200f86ff82b63a67726bac5ec4c
```
`getapplicationlog` RPC-call for this transaction tells us the following:
```
{
   "jsonrpc" : "2.0",
   "result" : {
      "executions" : [
         {
            "gasconsumed" : "16522710",
            "vmstate" : "HALT",
            "notifications" : [
               {
                  "eventname" : "Transfer",
                  "state" : {
                     "value" : [
                        {
                           "type" : "Any"
                        },
                        {
                           "type" : "ByteString",
                           "value" : "rtjbO11mTKBL3XvtN5kv+VhWC0E="
                        },
                        {
                           "value" : "1100000000000000",
                           "type" : "Integer"
                        }
                     ],
                     "type" : "Array"
                  },
                  "contract" : "0x27502a01e2fb013e1e4c428abb7b360df9f3f0cb"
               }
            ],
            "trigger" : "Application",
            "stack" : [
               {
                  "type" : "Boolean",
                  "value" : true
               }
            ]
         }
      ],
      "txid" : "0x296cb753f9afeef7ace3690eed05c08336329200f86ff82b63a67726bac5ec4c"
   },
   "id" : 1
}
```
Here we have `true` at the `stack` field, which means that token was successfully minted.
Let's just ensure that by querying `balanceOf` one more time:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 27502a01e2fb013e1e4c428abb7b360df9f3f0cb balanceOf NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB
```
... with `qwerty` pass. The result:
``` 
Sent invocation transaction 82ed056d7f9d27c5366561eb897d08a382747cd54b98e5c05fa82c30818f363b
```
... with the following `getapplicationlog` JSON message:
```
{
   "result" : {
      "executions" : [
         {
            "vmstate" : "HALT",
            "notifications" : [],
            "trigger" : "Application",
            "gasconsumed" : "5557020",
            "stack" : [
               {
                  "value" : "1100000000000000",
                  "type" : "Integer"
               }
            ]
         }
      ],
      "txid" : "0x82ed056d7f9d27c5366561eb897d08a382747cd54b98e5c05fa82c30818f363b"
   },
   "jsonrpc" : "2.0",
   "id" : 1
}
```
Now we can see integer value at the `stack` field, so `1100000000000000` is the nep17 token balance of our account.

Note, that token can be minted only once.

#### Step #5

After we are done with minting, it's possible to transfer token to someone else.
Let's transfer 5 tokens from our account to `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm` with `transfer` call:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 27502a01e2fb013e1e4c428abb7b360df9f3f0cb transfer NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm 500000000 null -- NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB:CalledByEntry
```
... with password `qwerty` and following result:
``` 
Sent invocation transaction a8dac2052450664697f02e852b82485225f7b3a1d1017eda2b4362fbc0cc962d
```
Our favourite `getapplicationlog` RPC-call tells us:
```
{
   "result" : {
      "executions" : [
         {
            "notifications" : [
               {
                  "contract" : "0x27502a01e2fb013e1e4c428abb7b360df9f3f0cb",
                  "state" : {
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
                     ],
                     "type" : "Array"
                  },
                  "eventname" : "Transfer"
               }
            ],
            "stack" : [
               {
                  "value" : true,
                  "type" : "Boolean"
               }
            ],
            "vmstate" : "HALT",
            "gasconsumed" : "14760830",
            "trigger" : "Application"
         }
      ],
      "txid" : "0xa8dac2052450664697f02e852b82485225f7b3a1d1017eda2b4362fbc0cc962d"
   },
   "jsonrpc" : "2.0",
   "id" : 1
}
```
Note, that `stack` field contains `true`, which means that token was successfully transferred.
Let's now check the balance of `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm` account to ensure that the amount of token on that account = 5:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 27502a01e2fb013e1e4c428abb7b360df9f3f0cb balanceOf NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm
```
The `getapplicationlog` RPC-call for this transaction tells us the following:
```
{
   "id" : 1,
   "jsonrpc" : "2.0",
   "result" : {
      "executions" : [
         {
            "notifications" : [],
            "gasconsumed" : "5557020",
            "trigger" : "Application",
            "vmstate" : "HALT",
            "stack" : [
               {
                  "type" : "Integer",
                  "value" : "500000000"
               }
            ]
         }
      ],
      "txid" : "0x950c454ed7f2e79124a89a8f7cdcc16205fa544c5935012033c328cc497e834d"
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
Sent invocation transaction 1306887c24fef841cbcc3dee3dbea734a0084c5f698ca62244bfda8f0dec4aba
Contract: a4ded8036fd90cf75daeefa7828498b80eee3e97
```   
You know, what it means :)

#### Step #3

Invoke the contract to register domain with name `my_first_domain`:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json a4ded8036fd90cf75daeefa7828498b80eee3e97 register my_first_domain NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB -- NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB:CalledByEntry
```
... the strongest password in the world, guess: `qwerty`
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```
Result:
```
Sent invocation transaction aca960d485f63fd0feca3fa4e5217f38350a7cfec0528a16cbff1aa67279ab34
```
Also you can see the log message in the console, where you run neo-go node:
```
2021-11-29T17:23:45.356+0300	INFO	runtime log	{"tx": "aca960d485f63fd0feca3fa4e5217f38350a7cfec0528a16cbff1aa67279ab34", "script": "a4ded8036fd90cf75daeefa7828498b80eee3e97", "msg": "RegisterDomain: my_first_domain"}
```
Well, that's ok. Let's check now, whether our domain was registered with `getapplicationlog` RPC-call:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["aca960d485f63fd0feca3fa4e5217f38350a7cfec0528a16cbff1aa67279ab34"] }' localhost:20331 | json_pp
```
The result is:
```
{
   "jsonrpc" : "2.0",
   "result" : {
      "txid" : "0xaca960d485f63fd0feca3fa4e5217f38350a7cfec0528a16cbff1aa67279ab34",
      "executions" : [
         {
            "notifications" : [
               {
                  "contract" : "0xa4ded8036fd90cf75daeefa7828498b80eee3e97",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "value" : "rtjbO11mTKBL3XvtN5kv+VhWC0E=",
                           "type" : "ByteString"
                        },
                        {
                           "type" : "ByteString",
                           "value" : "bXlfZmlyc3RfZG9tYWlu"
                        }
                     ]
                  },
                  "eventname" : "registered"
               }
            ],
            "gasconsumed" : "9143210",
            "stack" : [
               {
                  "type" : "Boolean",
                  "value" : true
               }
            ],
            "trigger" : "Application",
            "vmstate" : "HALT"
         }
      ]
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
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json a4ded8036fd90cf75daeefa7828498b80eee3e97 query my_first_domain
```
... the pass `qwerty`:
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```
Result:
```
Sent invocation transaction ddec59cd213a43e70f489e6e68ef76924f8a26538cd12b169d72ce78555c9d3a
```
and log-message:
```
2021-11-29T17:26:30.476+0300	INFO	runtime log	{"tx": "ddec59cd213a43e70f489e6e68ef76924f8a26538cd12b169d72ce78555c9d3a", "script": "a4ded8036fd90cf75daeefa7828498b80eee3e97", "msg": "QueryDomain: my_first_domain"}
```
Let's check this transaction with `getapplicationlog` RPC call:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["ddec59cd213a43e70f489e6e68ef76924f8a26538cd12b169d72ce78555c9d3a"] }' localhost:20331 | json_pp
```
... which gives us the following result:
```
{
   "jsonrpc" : "2.0",
   "id" : 1,
   "result" : {
      "executions" : [
         {
            "gasconsumed" : "4321230",
            "vmstate" : "HALT",
            "notifications" : [],
            "trigger" : "Application",
            "stack" : [
               {
                  "type" : "ByteString",
                  "value" : "rtjbO11mTKBL3XvtN5kv+VhWC0E="
               }
            ]
         }
      ],
      "txid" : "0xddec59cd213a43e70f489e6e68ef76924f8a26538cd12b169d72ce78555c9d3a"
   }
}
```

with base64 interpretation of our account address `NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB` on the stack, which means that domain `my_first_domain` was registered by the owner with received account address.

#### Step #5

Invoke the contract to transfer domain to the other account (e.g. account with `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm` address):
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json a4ded8036fd90cf75daeefa7828498b80eee3e97 transfer my_first_domain NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm -- NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB
```
... the password: `qwerty`
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```
Result:
```
Sent invocation transaction 108d62cefd64e3adea47025bf23e3749b604aa05422e515aecaaaaa3d0b6c9a3
```
and log-message:
```
2021-11-29T17:28:00.535+0300	INFO	runtime log	{"tx": "108d62cefd64e3adea47025bf23e3749b604aa05422e515aecaaaaa3d0b6c9a3", "script": "a4ded8036fd90cf75daeefa7828498b80eee3e97", "msg": "TransferDomain: my_first_domain"}
```
Perfect. And `getapplicationlog` RPC-call...
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["108d62cefd64e3adea47025bf23e3749b604aa05422e515aecaaaaa3d0b6c9a3"] }' localhost:20331 | json_pp
```
... tells us:
```
{
   "jsonrpc" : "2.0",
   "id" : 1,
   "result" : {
      "txid" : "0x108d62cefd64e3adea47025bf23e3749b604aa05422e515aecaaaaa3d0b6c9a3",
      "executions" : [
         {
            "notifications" : [
               {
                  "contract" : "0xa4ded8036fd90cf75daeefa7828498b80eee3e97",
                  "state" : {
                     "value" : [
                        {
                           "type" : "Buffer",
                           "value" : "rtjbO11mTKBL3XvtN5kv+VhWC0E="
                        },
                        {
                           "value" : "bXlfZmlyc3RfZG9tYWlu",
                           "type" : "ByteString"
                        }
                     ],
                     "type" : "Array"
                  },
                  "eventname" : "deleted"
               },
               {
                  "contract" : "0xa4ded8036fd90cf75daeefa7828498b80eee3e97",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "value" : "50l6vFaauRKm8hPVkr3Aw2CeHQs=",
                           "type" : "ByteString"
                        },
                        {
                           "value" : "bXlfZmlyc3RfZG9tYWlu",
                           "type" : "ByteString"
                        }
                     ]
                  },
                  "eventname" : "registered"
               }
            ],
            "vmstate" : "HALT",
            "stack" : [
               {
                  "value" : true,
                  "type" : "Boolean"
               }
            ],
            "trigger" : "Application",
            "gasconsumed" : "7679990"
         }
      ]
   }
}
```
The `notifications` field contains two events:
- First one with name `deleted` and additional information (domain `my_first_domain` was deleted from account `NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB`),
- Second one with name `registered`  and additional information  (domain `my_first_domain` was registered with account `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm`).
The `stack` field contains `true` value, which means that our domain was successfully transferred.

#### Step #6

The last call is `delete`, so you can try to create the other domain, e.g. `my_second_domain` and then remove it from storage with:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json a4ded8036fd90cf75daeefa7828498b80eee3e97 register my_second_domain NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB -- NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB:CalledByEntry
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json a4ded8036fd90cf75daeefa7828498b80eee3e97 delete my_second_domain -- NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB
```

Thank you!

### Useful links

* [Our basic tutorial on Medium](https://medium.com/@neospcc/%D1%81%D0%BC%D0%B0%D1%80%D1%82-%D0%BA%D0%BE%D0%BD%D1%82%D1%80%D0%B0%D0%BA%D1%82-%D0%B4%D0%BB%D1%8F-neo-769139352b65)
* [Using Neo Blockchain Toolkit](https://medium.com/@neospcc/neogo-adds-support-for-neo-blockchain-toolkit-673ea914f661)
* [Neo documentation](https://docs.neo.org/)
* [Neo github](https://github.com/neo-project/neo/)
* [NeoGo github](https://github.com/nspcc-dev/neo-go)
