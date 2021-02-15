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
        $ ./bin/neo-go wallet nep17 transfer -w .docker/wallets/wallet1.json --out my_tx.json -r http://localhost:20331 --from NUVPACMnKFhpuHjsRjhUvXz1XhqfGZYVtY --to NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt --token GAS --amount 29999999
    ``` 
    Where
    - `./bin/neo-go` runs neo-go
    - `wallet nep17 transfer` - command with arguments in [neo-go](https://github.com/nspcc-dev/neo-go/blob/master/cli/wallet/nep17.go#L108)
    - `-w .docker/wallets/wallet1.json` - path to the [wallet](https://github.com/nspcc-dev/neo-go/blob/master/.docker/wallets/wallet1.json) for the first node in the private network
    - `--out my_tx.json` - output file for the signed transaction
    - `-r http://localhost:20331` - RPC node endpoint
    - `--from NUVPACMnKFhpuHjsRjhUvXz1XhqfGZYVtY` - multisig account to transfer GAS from
    - `--to NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt` - our account from the [wallet](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/my_wallet.json)
    - `--token GAS` - transferred token name, which is GAS
    - `--amount 29999999` - amount of GAS to transfer
    
    Enter the password `one`:
    ```
    Password >
    ```
    The result is transaction signed by the first node `my_tx.json`.

2. Sign the created transaction using the second node address:

    ```
    $ ./bin/neo-go wallet multisig sign -w .docker/wallets/wallet2.json --in my_tx.json --out my_tx2.json --address NUVPACMnKFhpuHjsRjhUvXz1XhqfGZYVtY
    ```
    Where
    - `-w .docker/wallets/wallet2.json` - path to the [wallet](https://github.com/nspcc-dev/neo-go/blob/master/.docker/wallets/wallet2.json) for the second node in private network
    - `--in my_tx.json` - previously created transfer transaction
    - `--out my_tx2.json` - output file for the signed transaction
    - `--address NUVPACMnKFhpuHjsRjhUvXz1XhqfGZYVtY` - multisig account to sign the transaction
    
    Enter the password `two`:
    ```
    Enter password to unlock wallet and sign the transaction
    Password >
    ```
    The result is transaction signed by both first and second nodes.

3. Sign the transaction using the third node address and push it to the chain:
    ```
    $ ./bin/neo-go wallet multisig sign -w ./.docker/wallets/wallet3.json --in my_tx2.json --out my_tx3.json --address NUVPACMnKFhpuHjsRjhUvXz1XhqfGZYVtY -r http://localhost:20331
    ```
    Enter the password `three`:
    ```
    Enter password to unlock wallet and sign the transaction
    Password >
    ```
    The result is transaction signed by the first, second and third nodes and deployed to the chain.

4. Check the balance:

    Now you should have 29999999 GAS on the balance of `NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt` account.
    To check the transfer was successfully submitted use `getnep17transfers` RPC call:
    ```
    curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getnep17transfers", "params": ["NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt"] }' localhost:20331 | json_pp
    ```
    The result should look like the following:
```
    {
       "jsonrpc" : "2.0",
       "result" : {
          "address" : "NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt",
          "sent" : [],
          "received" : [
             {
                "assethash" : "0xa6a6c15dcdc9b997dac448b6926522d22efeedfb",
                "transfernotifyindex" : 0,
                "amount" : "2999999900000000",
                "txhash" : "0xd44399f508e1ea6ebfd01cf9252a52e0e8699a1b799f60c03f3daae8f74af741",
                "blockindex" : 41,
                "transferaddress" : "NUVPACMnKFhpuHjsRjhUvXz1XhqfGZYVtY",
                "timestamp" : 1608206555652
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
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```

Result:
```
Contract: fb8bf3f9cab177b1c031a654fc0c94bf3db2d9b8
725857b43b5b5a67d712946b25b96a832671b21266e16be76419927d102bf9bc
```

At this point your ‘Hello World’ contract is deployed and could be invoked. Let’s do it as a final step.

#### Step 4
Invoke contract.
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json fb8bf3f9cab177b1c031a654fc0c94bf3db2d9b8 main
```

Where
- `contract invokefunction` runs invoke with provided parameters
- `-r http://localhost:20331` defines RPC endpoint used for function call
- `-w my_wallet.json` is a wallet
- `fb8bf3f9cab177b1c031a654fc0c94bf3db2d9b8` contract hash got as an output from the previous command (deployment in step 6)
- `main` - method to be called

Enter password `qwerty` for account:
```
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```

Result:
```
Sent invocation transaction 79a4aae4dfc2df17ee293902fdc10a5590289114635dfe551a4cd4715f4926c6
```
In the console where you were running step #5 you will get:
```
2020-12-17T15:29:48.790+0300	INFO	runtime log	{"script": "b8d9b23dbf940cfc54a631c0b177b1caf9f38bfb", "logs": "\"Hello, world!\""}
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
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getrawtransaction", "params": ["725857b43b5b5a67d712946b25b96a832671b21266e16be76419927d102bf9bc", 1] }' localhost:20331 | json_pp
```

Where
- `"jsonrpc": "2.0"` is protocol version
- `"id": 1` is id of current request
- `"method": "getrawtransaction"` is requested method
- `"params": ["725857b43b5b5a67d712946b25b96a832671b21266e16be76419927d102bf9bc", 1]` is an array of parameters, 
  where
   - `725857b43b5b5a67d712946b25b96a832671b21266e16be76419927d102bf9bc` is deployment transaction hash
   - `1` is `verbose` parameter for detailed JSON string output
- `json_pp` just makes the JSON output prettier

Result:
```
{
   "id" : 1,
   "jsonrpc" : "2.0",
   "result" : {
      "attributes" : [],
      "blockhash" : "0x7670186e0ea4ab665b41b11365c6c671ec424be695e163084cb59dde3de691cd",
      "blocktime" : 1613129835686,
      "confirmations" : 2,
      "hash" : "0x725857b43b5b5a67d712946b25b96a832671b21266e16be76419927d102bf9bc",
      "netfee" : "1532550",
      "nonce" : 4113748543,
      "script" : "DPZ7Im5hbWUiOiJIZWxsb1dvcmxkIGNvbnRyYWN0IiwiYWJpIjp7Im1ldGhvZHMiOlt7Im5hbWUiOiJtYWluIiwib2Zmc2V0IjowLCJwYXJhbWV0ZXJzIjpbXSwicmV0dXJudHlwZSI6IlZvaWQiLCJzYWZlIjpmYWxzZX1dLCJldmVudHMiOltdfSwiZ3JvdXBzIjpbXSwicGVybWlzc2lvbnMiOlt7ImNvbnRyYWN0IjoiKiIsIm1ldGhvZHMiOiIqIn1dLCJzdXBwb3J0ZWRzdGFuZGFyZHMiOltdLCJ0cnVzdHMiOltdLCJleHRyYSI6bnVsbH0MZE5FRjNuZW8tZ28tMC45My4wLXByZS0yMTItZ2M2NmFjNzgwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAWDA1IZWxsbywgd29ybGQhQc/nR5YhQM/mETQSwBsMBmRlcGxveQwUQw6fb7MTqNOit2E7Z4MJ0dfXAaVBYn1bUg==",
      "sender" : "NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt",
      "signers" : [
         {
            "account" : "0xf8d79b436e9d3a54a38cc877d182b71e381d7a5c",
            "scopes" : "None"
         }
      ],
      "size" : 549,
      "sysfee" : "1000999450",
      "validuntilblock" : 21,
      "version" : 0,
      "vmstate" : "HALT",
      "witnesses" : [
         {
            "invocation" : "DEDk+J0oFHKvw6zRrgrZBafchMg6g8NSii9ICM7oa+uOmImCe5lIqz9x4Jg6UP3tM97JUab/DvMNQoYE0HtAaRKP",
            "verification" : "DCECqIHFxkXMSrGdFGesF8hU+x8Y9RaRmsEOWyRK4B9eK30LQZVEDXg="
         }
      ]
   }
}
```

#### GetApplicationLog
[GetApplicationLog](https://docs.neo.org/v3/docs/en-us/reference/rpc/latest-version/api/getapplicationlog.html) returns the contract log based on the specified transaction id.

Request application log for invocation transaction from step #4:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["79a4aae4dfc2df17ee293902fdc10a5590289114635dfe551a4cd4715f4926c6"] }' localhost:20331 | json_pp
```

With a single parameter:
- `79a4aae4dfc2df17ee293902fdc10a5590289114635dfe551a4cd4715f4926c6` - invocation transaction hash from step #7

Result:
```
{
   "id" : 1,
   "jsonrpc" : "2.0",
   "result" : {
      "executions" : [
         {
            "gasconsumed" : "1982250",
            "notifications" : [],
            "stack" : [
               {
                  "type" : "Any"
               }
            ],
            "trigger" : "Application",
            "vmstate" : "HALT"
         }
      ],
      "txid" : "0x79a4aae4dfc2df17ee293902fdc10a5590289114635dfe551a4cd4715f4926c6"
   }
}
```

#### Other Useful RPC calls
```
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getversion", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getblockcount", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getconnectioncount", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getcontractstate", "params": ["0xfb8bf3f9cab177b1c031a654fc0c94bf3db2d9b8"] }' localhost:20331
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
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```

Result:
```
Contract: 9f9bb13864f635a540fa3c3f2ce410f267a076df
a71bd9e48b354fb2421027f79b1968d6b0f4e211fc1887325e90a90201184ee1
```   

Which means that our contract was deployed and now we can invoke it.

Let's check that the storage value was initialised with `0`. Use `getapplicaionlog` RPC-call for the deployment transaction:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["a71bd9e48b354fb2421027f79b1968d6b0f4e211fc1887325e90a90201184ee1"] }' localhost:20331 | json_pp
```

The JSON result is:
```
{
   "id" : 1,
   "jsonrpc" : "2.0",
   "result" : {
      "executions" : [
         {
            "gasconsumed" : "1006781960",
            "notifications" : [
               {
                  "contract" : "0xc5470e21df594ce4d143f6f9c1b4decd5e94c049",
                  "eventname" : "info",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "U3RvcmFnZSBrZXkgbm90IHlldCBzZXQuIFNldHRpbmcgdG8gMA=="
                        }
                     ]
                  }
               },
               {
                  "contract" : "0xc5470e21df594ce4d143f6f9c1b4decd5e94c049",
                  "eventname" : "info",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "U3RvcmFnZSBrZXkgaXMgaW5pdGlhbGlzZWQ="
                        }
                     ]
                  }
               },
               {
                  "contract" : "0xa501d7d7d10983673b61b7a2d3a813b36f9f0e43",
                  "eventname" : "Deploy",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "ScCUXs3etMH59kPR5ExZ3yEOR8U="
                        }
                     ]
                  }
               }
            ],
            "stack" : [
               {
                  "type" : "Array",
                  "value" : [
                     ... skipped serialized contract representation ...
                  ]
               }
            ],
            "trigger" : "Application",
            "vmstate" : "HALT"
         }
      ],
      "txid" : "0xa71bd9e48b354fb2421027f79b1968d6b0f4e211fc1887325e90a90201184ee1"
   }
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
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json c5470e21df594ce4d143f6f9c1b4decd5e94c049 main
```
... enter the password `qwerty`:
```
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```
Result:
```
Sent invocation transaction 30ada0a04423eb303c93ed7bc18b5dc7b38efcf5e856a4c169f6a025e43a8bf1
```
To check the counter value, call `getapplicaionlog` RPC-call for the invocation transaction:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["30ada0a04423eb303c93ed7bc18b5dc7b38efcf5e856a4c169f6a025e43a8bf1"] }' localhost:20331 | json_pp
```
The JSON result is:
```
{
   "id" : 1,
   "jsonrpc" : "2.0",
   "result" : {
      "executions" : [
         {
            "gasconsumed" : "5244260",
            "notifications" : [
               {
                  "contract" : "0xc5470e21df594ce4d143f6f9c1b4decd5e94c049",
                  "eventname" : "info",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "VmFsdWUgcmVhZCBmcm9tIHN0b3JhZ2U="
                        }
                     ]
                  }
               },
               {
                  "contract" : "0xc5470e21df594ce4d143f6f9c1b4decd5e94c049",
                  "eventname" : "info",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "U3RvcmFnZSBrZXkgYWxyZWFkeSBzZXQuIEluY3JlbWVudGluZyBieSAx"
                        }
                     ]
                  }
               },
               {
                  "contract" : "0xc5470e21df594ce4d143f6f9c1b4decd5e94c049",
                  "eventname" : "info",
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
            "stack" : [
               {
                  "type" : "Integer",
                  "value" : "1"
               }
            ],
            "trigger" : "Application",
            "vmstate" : "HALT"
         }
      ],
      "txid" : "0x30ada0a04423eb303c93ed7bc18b5dc7b38efcf5e856a4c169f6a025e43a8bf1"
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
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json c5470e21df594ce4d143f6f9c1b4decd5e94c049 main
```
... enter the password `qwerty`:
```
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```
Result:
```
Sent invocation transaction 23dbfd720451dd56274910365cfb982325f4312ce213c9551737ace5757f4ed4
```
To check the counter value, call `getapplicaionlog` RPC-call for the invocation transaction:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["23dbfd720451dd56274910365cfb982325f4312ce213c9551737ace5757f4ed4"] }' localhost:20331 | json_pp
```
The JSON result is:
```
{
   "id" : 1,
   "jsonrpc" : "2.0",
   "result" : {
      "executions" : [
         {
            "gasconsumed" : "5144260",
            "notifications" : [
               {
                  "contract" : "0xc5470e21df594ce4d143f6f9c1b4decd5e94c049",
                  "eventname" : "info",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "VmFsdWUgcmVhZCBmcm9tIHN0b3JhZ2U="
                        }
                     ]
                  }
               },
               {
                  "contract" : "0xc5470e21df594ce4d143f6f9c1b4decd5e94c049",
                  "eventname" : "info",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "U3RvcmFnZSBrZXkgYWxyZWFkeSBzZXQuIEluY3JlbWVudGluZyBieSAx"
                        }
                     ]
                  }
               },
               {
                  "contract" : "0xc5470e21df594ce4d143f6f9c1b4decd5e94c049",
                  "eventname" : "info",
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
            "stack" : [
               {
                  "type" : "Integer",
                  "value" : "2"
               }
            ],
            "trigger" : "Application",
            "vmstate" : "HALT"
         }
      ],
      "txid" : "0x23dbfd720451dd56274910365cfb982325f4312ce213c9551737ace5757f4ed4"
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
  - `holder` which is requested address
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
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```

Result:
```
Contract: 1b63c282ce2618c8882694cc4214d8faee69ac42
293c9bb4977000f4bbf1aab955a79d67c9fbbfb7286f101f6cf3e34d70636e03
```   

Which means that our contract was deployed and now we can invoke it.

#### Step #2
Let's invoke the contract to perform different operations.

To start with, query `Symbol` of the created nep17 token:

```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 1b63c282ce2618c8882694cc4214d8faee69ac42 symbol
```                                                                   
Where
- `1b63c282ce2618c8882694cc4214d8faee69ac42` is our contract hash from step #1
- `symbol` is operation string which was described earlier and returns token symbol

... and don't forget the password of your account `qwerty`.

Result:
```
Sent invocation transaction 48f2f4f21aa92fd77247df20c83c652db77163a02dd079960e89f031142895d4
```                                                                                         
Now, let's take a detailed look at this invocation transaction with `getapplicationlog` RPC call:

```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["48f2f4f21aa92fd77247df20c83c652db77163a02dd079960e89f031142895d4"] }' localhost:20331 | json_pp
```               

Result:
```
{
   "jsonrpc" : "2.0",
   "id" : 1,
   "result" : {
      "txid" : "0x48f2f4f21aa92fd77247df20c83c652db77163a02dd079960e89f031142895d4",
      "executions" : [
         {
            "stack" : [
               {
                  "type" : "ByteString",
                  "value" : "QU5U"
               }
            ],
            "notifications" : [],
            "trigger" : "Application",
            "gasconsumed" : "40638600",
            "vmstate" : "HALT"
         }
      ]
   }
}
```

At least, you can see that `stack` field of JSON result is not empty: it contains base64 byte array with the symbol of our token.

Following commands able you to get some additional information about token:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 1b63c282ce2618c8882694cc4214d8faee69ac42 decimals
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 1b63c282ce2618c8882694cc4214d8faee69ac42 totalSupply
```

#### Step #3

Now it's time for more interesting things. First of all, let's check the balance of nep17 token on our account by using `balanceOf`:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 1b63c282ce2618c8882694cc4214d8faee69ac42 balanceOf NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt
```                             
... with `qwerty` password. The result is:
```
Sent invocation transaction 6db68a92683ffb5f6f40700d85ebb08058952d3ce0fb23ed837228118cc14025
```
And take a closer look at the transaction's details with `getapplicationlog` RPC-call:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["6db68a92683ffb5f6f40700d85ebb08058952d3ce0fb23ed837228118cc14025"] }' localhost:20331 | json_pp
```
Result:
```
{
   "id" : 1,
   "result" : {
      "vmstate" : "HALT",
      "notifications" : [],
      "trigger" : "Application",
      "gasconsumed" : "5423900",
      "txid" : "0xde1a69a6c0b5389458f35c44e28d1bea0a0ab4fe7b4ad0013e34b36d0995069a",
      "stack" : [
         {
            "type" : "Integer",
            "value" : "0"
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
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 1b63c282ce2618c8882694cc4214d8faee69ac42 mint NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt -- f8d79b436e9d3a54a38cc877d182b71e381d7a5c:CalledByEntry
```
Where:
- `--` is a special delimiter of transaction's cosigners list
- `f8d79b436e9d3a54a38cc877d182b71e381d7a5c` is the signer itself (which is hex-encoded LE representation of our `NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt` account)

... with `qwerty` pass. The result:
``` 
Sent invocation transaction 5288982aa4566f0a09d7bc2fb545d74602c097cdc580626a46b4995030bc196d
```
`getapplicationlog` RPC-call for this transaction tells us the following:
```
{
   "id" : 1,
   "jsonrpc" : "2.0",
   "result" : {
      "executions" : [
         {
            "gasconsumed" : "12006390",
            "notifications" : [
               {
                  "contract" : "0x1b63c282ce2618c8882694cc4214d8faee69ac42",
                  "eventname" : "Transfer",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "type" : "Any"
                        },
                        {
                           "type" : "ByteString",
                           "value" : "XHodOB63gtF3yIyjVDqdbkOb1/g="
                        },
                        {
                           "type" : "Integer",
                           "value" : "1100000000000000"
                        }
                     ]
                  }
               }
            ],
            "stack" : [
               {
                  "type" : "Integer",
                  "value" : "1"
               }
            ],
            "trigger" : "Application",
            "vmstate" : "HALT"
         }
      ],
      "txid" : "0x5288982aa4566f0a09d7bc2fb545d74602c097cdc580626a46b4995030bc196d"
   }
}
```
Here we have `1` at the `stack` field, which means that token was successfully minted.
Let's just ensure that by querying `balanceOf` one more time:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 1b63c282ce2618c8882694cc4214d8faee69ac42 balanceOf NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt
```
... with `qwerty` pass. The result:
``` 
Sent invocation transaction ce8212e93ae6171efbdbb415153184ad06e4dbdd656c15ffdfcbbcd57463d7be
```
... with the following `getapplicationlog` JSON message:
```
{
   "jsonrpc" : "2.0",
   "id" : 1,
   "result" : {
      "txid" : "0xce8212e93ae6171efbdbb415153184ad06e4dbdd656c15ffdfcbbcd57463d7be",
      "executions" : [
         {
            "vmstate" : "HALT",
            "notifications" : [],
            "gasconsumed" : "51422700",
            "trigger" : "Application",
            "stack" : [
               {
                  "type" : "Integer",
                  "value" : "1100000000000000"
               }
            ]
         }
      ]
   }
}
```
Now we can see integer value at the `stack` field, so `1100000000000000` is the nep17 token balance of our account.

Note, that token can be minted only once.

#### Step #5

After we are done with minting, it's possible to transfer token to someone else.
Let's transfer 5 tokens from our account to `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm` with `transfer` call:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 1b63c282ce2618c8882694cc4214d8faee69ac42 transfer NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm 500000000 null -- f8d79b436e9d3a54a38cc877d182b71e381d7a5c:CalledByEntry
```
... with password `qwerty` and following result:
``` 
Sent invocation transaction 8fefd5861d865379b6632693a5c74f430eb4465b6abf2efef8bddfd0a17dfa70
```
Our favourite `getapplicationlog` RPC-call tells us:
```
{
   "jsonrpc" : "2.0",
   "result" : {
      "executions" : [
         {
            "notifications" : [
               {
                  "state" : {
                     "value" : [
                        {
                           "value" : "XHodOB63gtF3yIyjVDqdbkOb1/g=",
                           "type" : "ByteString"
                        },
                        {
                           "value" : "50l6vFaauRKm8hPVkr3Aw2CeHQs=",
                           "type" : "ByteString"
                        },
                        {
                           "type" : "Integer",
                           "value" : "500000000"
                        }
                     ],
                     "type" : "Array"
                  },
                  "eventname" : "Transfer",
                  "contract" : "0x1b63c282ce2618c8882694cc4214d8faee69ac42"
               }
            ],
            "gasconsumed" : "98136200",
            "trigger" : "Application",
            "vmstate" : "HALT",
            "stack" : [
               {
                  "value" : "1",
                  "type" : "Integer"
               }
            ]
         }
      ],
      "txid" : "0x8fefd5861d865379b6632693a5c74f430eb4465b6abf2efef8bddfd0a17dfa70"
   },
   "id" : 1
}
```
Note, that `stack` field contains `1`, which means that token was successfully transferred.
Let's now check the balance of `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm` account to ensure that the amount of token on that account = 5:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 1b63c282ce2618c8882694cc4214d8faee69ac42 balanceOf NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm
```
The `getapplicationlog` RPC-call for this transaction tells us the following:
```
{
   "result" : {
      "txid" : "0x0cb6db18b98a188d15010f8d9002948168897eea119e1b28f3254a4c554feab4",
      "executions" : [
         {
            "gasconsumed" : "5142270",
            "trigger" : "Application",
            "vmstate" : "HALT",
            "notifications" : [],
            "stack" : [
               {
                  "value" : "500000000",
                  "type" : "Integer"
               }
            ]
         }
      ]
   },
   "jsonrpc" : "2.0",
   "id" : 1
}
```
Here we are! There are exactly 5 tokens at the `stack` field. You can also ensure that these 5 tokens were debited from `NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt` account by using `balanceOf` method.

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
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```

Result:
```
Contract: 66304b7c84fd61bf6ab5b9401e8fb7fe6ee8e486
91b3d71991cc6e89923fcc3f53dfc708c6dbeefa924ecebbc3a9b7c6f8f01f94
```   
You know, what it means :)

#### Step #3

Invoke the contract to register domain with name `my_first_domain`:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 66304b7c84fd61bf6ab5b9401e8fb7fe6ee8e486 register my_first_domain NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt -- f8d79b436e9d3a54a38cc877d182b71e381d7a5c:CalledByEntry
```
... the strongest password in the world, guess: `qwerty`
```
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```
Result:
```
Sent invocation transaction 114b0c93dfb540cd86e50470dd628c3884a391568667acb6a44bd3800ff640f1
```
Also you can see the log message in the console, where you run neo-go node:
```
2020-12-17T17:31:43.480+0300	INFO	runtime log	{"script": "86e4e86efeb78f1e40b9b56abf61fd847c4b3066", "logs": "\"RegisterDomain: my_first_domain\""}
```
Well, that's ok. Let's check now, whether our domain was registered with `getapplicationlog` RPC-call:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["114b0c93dfb540cd86e50470dd628c3884a391568667acb6a44bd3800ff640f1"] }' localhost:20331 | json_pp
```
The result is:
```
{
   "id" : 1,
   "jsonrpc" : "2.0",
   "result" : {
      "executions" : [
         {
            "gasconsumed" : "7637810",
            "notifications" : [
               {
                  "contract" : "0x66304b7c84fd61bf6ab5b9401e8fb7fe6ee8e486",
                  "eventname" : "registered",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "XHodOB63gtF3yIyjVDqdbkOb1/g="
                        },
                        {
                           "type" : "ByteString",
                           "value" : "bXlfZmlyc3RfZG9tYWlu"
                        }
                     ]
                  }
               }
            ],
            "stack" : [
               {
                  "type" : "Integer",
                  "value" : "1"
               }
            ],
            "trigger" : "Application",
            "vmstate" : "HALT"
         }
      ],
      "txid" : "0x114b0c93dfb540cd86e50470dd628c3884a391568667acb6a44bd3800ff640f1"
   }
}
```
Especially, we're interested in two fields of the json:

First one is `notifications` field, which contains one notification with `registered` name:
- `bXlfZmlyc3RfZG9tYWlu` byte string in base64 representation, which can be decoded to `my_first_domain` - our domain's name
- `XHodOB63gtF3yIyjVDqdbkOb1/g=` byte array, which can be decoded to the account address `NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt`.

The second field is `stack` with `1` value, which was returned by the smart contract.

All of these values let us be sure that our domain was successfully registered.  

#### Step #4

Invoke the contract to query the address information our `my_first_domain` domain:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 66304b7c84fd61bf6ab5b9401e8fb7fe6ee8e486 query my_first_domain
```
... the pass `qwerty`:
```
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```
Result:
```
Sent invocation transaction 1e941ce1641fed5593a3dfb2ddc1a10090c3b466d1b5bf636f0ff1ee458a9554
```
and log-message:
```
2020-12-17T17:39:32.677+0300	INFO	runtime log	{"script": "86e4e86efeb78f1e40b9b56abf61fd847c4b3066", "logs": "\"QueryDomain: my_first_domain\""}
```
Let's check this transaction with `getapplicationlog` RPC call:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["1e941ce1641fed5593a3dfb2ddc1a10090c3b466d1b5bf636f0ff1ee458a9554"] }' localhost:20331 | json_pp
```
... which gives us the following result:
```
{
   "id" : 1,
   "jsonrpc" : "2.0",
   "result" : {
      "executions" : [
         {
            "gasconsumed" : "4090830",
            "notifications" : [],
            "stack" : [
               {
                  "type" : "ByteString",
                  "value" : "XHodOB63gtF3yIyjVDqdbkOb1/g="
               }
            ],
            "trigger" : "Application",
            "vmstate" : "HALT"
         }
      ],
      "txid" : "0x1e941ce1641fed5593a3dfb2ddc1a10090c3b466d1b5bf636f0ff1ee458a9554"
   }
}
```

with base64 interpretation of our account address `NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt` on the stack, which means that domain `my_first_domain` was registered by the owner with received account address.

#### Step #5

Invoke the contract to transfer domain to the other account (e.g. account with `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm` address):
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 66304b7c84fd61bf6ab5b9401e8fb7fe6ee8e486 transfer my_first_domain NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm -- f8d79b436e9d3a54a38cc877d182b71e381d7a5c
```
... the password: `qwerty`
```
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```
Result:
```
Sent invocation transaction 012dc71a593b024fd38e4d301a1ded85841f1cb0d0d94c81f42846699fda8a3a
```
and log-message:
```
2020-12-17T17:44:07.536+0300	INFO	runtime log	{"script": "86e4e86efeb78f1e40b9b56abf61fd847c4b3066", "logs": "\"TransferDomain: my_first_domain\""}
```
Perfect. And `getapplicationlog` RPC-call...
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["012dc71a593b024fd38e4d301a1ded85841f1cb0d0d94c81f42846699fda8a3a"] }' localhost:20331 | json_pp
```
... tells us:
```
{
   "id" : 1,
   "jsonrpc" : "2.0",
   "result" : {
      "executions" : [
         {
            "gasconsumed" : "5759750",
            "notifications" : [
               {
                  "contract" : "0x66304b7c84fd61bf6ab5b9401e8fb7fe6ee8e486",
                  "eventname" : "deleted",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "XHodOB63gtF3yIyjVDqdbkOb1/g="
                        },
                        {
                           "type" : "ByteString",
                           "value" : "bXlfZmlyc3RfZG9tYWlu"
                        }
                     ]
                  }
               },
               {
                  "contract" : "0x66304b7c84fd61bf6ab5b9401e8fb7fe6ee8e486",
                  "eventname" : "registered",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "50l6vFaauRKm8hPVkr3Aw2CeHQs="
                        },
                        {
                           "type" : "ByteString",
                           "value" : "bXlfZmlyc3RfZG9tYWlu"
                        }
                     ]
                  }
               }
            ],
            "stack" : [
               {
                  "type" : "Integer",
                  "value" : "1"
               }
            ],
            "trigger" : "Application",
            "vmstate" : "HALT"
         }
      ],
      "txid" : "0x012dc71a593b024fd38e4d301a1ded85841f1cb0d0d94c81f42846699fda8a3a"
   }
}
```
The `notifications` field contains two events:
- First one with name `deleted` and additional information (domain `my_first_domain` was deleted from account `NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt`),
- Second one with name `registered`  and additional information  (domain `my_first_domain` was registered with account `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm`).
The `stack` field contains `1` value, which means that our domain was successfully transferred.

#### Step #6

The last call is `delete`, so you can try to create the other domain, e.g. `my_second_domain` and then remove it from storage with:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 66304b7c84fd61bf6ab5b9401e8fb7fe6ee8e486 register my_second_domain NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt -- f8d79b436e9d3a54a38cc877d182b71e381d7a5c:CalledByEntry
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 66304b7c84fd61bf6ab5b9401e8fb7fe6ee8e486 delete my_second_domain -- f0a33acbf435417f3d2c18445d607d1e6f48d413
```

Thank you!

### Useful links

* [Our basic tutorial on Medium](https://medium.com/@neospcc/%D1%81%D0%BC%D0%B0%D1%80%D1%82-%D0%BA%D0%BE%D0%BD%D1%82%D1%80%D0%B0%D0%BA%D1%82-%D0%B4%D0%BB%D1%8F-neo-769139352b65)
* [Using NEO Blockchain Toolkit](https://medium.com/@neospcc/neogo-adds-support-for-neo-blockchain-toolkit-673ea914f661)
* [NEO documentation](https://docs.neo.org/)
* [NEO github](https://github.com/neo-project/neo/)
* [NEO-GO github](https://github.com/nspcc-dev/neo-go)
