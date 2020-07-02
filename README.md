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
As you can see there are a lot of options to play with it. Let’s take simple smart contract(1-print.go) and compile it:
 
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
NEO-GO-VM > loadgo test.go
READY: loaded 38 instructions
NEO-GO-VM 0 >  
```
And there you can see how many instructions were generated and even if you are interested in opcodes of current program you can dump them:
```
NEO-GO-VM 0 > ops
INDEX    OPCODE          PARAMETER                                       
0        PUSH1                                                           <<
1        NEWARRAY                                                        
2        TOALTSTACK                                                      
3        PUSHBYTES13     48656c6c6f2c20776f726c6421 ("Hello, world!")    
17       SYSCALL         "Neo.Runtime.Log"                               
34       NOP                                                             
35       FROMALTSTACK                                                    
36       DROP                                                            
37       RET     
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
Download neo-go and build it (master branch):
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
2020-06-30T16:26:47.549+0300	INFO	restoring blockchain	{"version": "0.1.0"}

    _   ____________        __________
   / | / / ____/ __ \      / ____/ __ \
  /  |/ / __/ / / / /_____/ / __/ / / /
 / /|  / /___/ /_/ /_____/ /_/ / /_/ /
/_/ |_/_____/\____/      \____/\____/

/NEO-GO:0.90.0-pre-610-g59d1013a/

2020-06-30T16:26:47.549+0300	INFO	service is running	{"service": "Prometheus", "endpoint": ":2112"}
2020-06-30T16:26:47.549+0300	INFO	service hasn't started since it's disabled	{"service": "Pprof"}
2020-06-30T16:26:47.550+0300	INFO	starting rpc-server	{"endpoint": ":20331"}
2020-06-30T16:26:47.550+0300	INFO	node started	{"blockHeight": 0, "headerHeight": 0}
2020-06-30T16:26:47.551+0300	INFO	new peer connected	{"addr": "127.0.0.1:20333", "peerCount": 1}
2020-06-30T16:26:47.551+0300	INFO	new peer connected	{"addr": "127.0.0.1:20334", "peerCount": 2}
2020-06-30T16:26:47.552+0300	INFO	new peer connected	{"addr": "127.0.0.1:20335", "peerCount": 3}
2020-06-30T16:26:47.552+0300	INFO	new peer connected	{"addr": "127.0.0.1:20336", "peerCount": 4}
2020-06-30T16:26:47.553+0300	INFO	started protocol	{"addr": "127.0.0.1:20334", "userAgent": "/NEO-GO:0.90.0-pre-610-g59d1013a/", "startHeight": 3, "id": 2280194870}
2020-06-30T16:26:47.554+0300	INFO	started protocol	{"addr": "127.0.0.1:20333", "userAgent": "/NEO-GO:0.90.0-pre-610-g59d1013a/", "startHeight": 3, "id": 3666776256}
2020-06-30T16:26:47.555+0300	INFO	started protocol	{"addr": "127.0.0.1:20336", "userAgent": "/NEO-GO:0.90.0-pre-610-g59d1013a/", "startHeight": 3, "id": 1699156200}
2020-06-30T16:26:47.555+0300	INFO	started protocol	{"addr": "127.0.0.1:20335", "userAgent": "/NEO-GO:0.90.0-pre-610-g59d1013a/", "startHeight": 3, "id": 874998449}
2020-06-30T16:26:48.550+0300	INFO	blockchain persist completed	{"persistedBlocks": 0, "persistedKeys": 4, "headerHeight": 3, "blockHeight": 0, "took": "437.278µs"}
...
```

#### Step 4
Transfer some GAS from multisig account to our account.

1. Create NEP5 transfer transaction:
    ```
        $ ./bin/neo-go wallet nep5 transfer -w .docker/wallets/wallet1.json --out my_tx.json -r http://localhost:20331 --from Nbb1qkwcwNSBs9pAnrVVrnFbWnbWBk91U2 --to NMipL5VsNoLUBUJKPKLhxaEbPQVCZnyJyB --token gas --amount 29999999
    ``` 
    Where
    - `./bin/neo-go` runs neo-go
    - `wallet nep5 transfer` - command with arguments in [neo-go](https://github.com/nspcc-dev/neo-go/blob/master/cli/wallet/nep5.go#L108)
    - `-w .docker/wallets/wallet1.json` - path to the [wallet](https://github.com/nspcc-dev/neo-go/blob/master/.docker/wallets/wallet1.json) for the first node in the private network
    - `--out my_tx.json` - output file for the signed transaction
    - `-r http://localhost:20331` - RPC node endpoint
    - `--from Nbb1qkwcwNSBs9pAnrVVrnFbWnbWBk91U2` - multisig account to transfer GAS from
    - `--to NMipL5VsNoLUBUJKPKLhxaEbPQVCZnyJyB` - our account from the [wallet](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/my_wallet.json)
    - `--token gas` - transferred token name, which is GAS
    - `--amount 29999999` - amount of GAS to transfer
    
    Enter the password `one`:
    ```
    Password >
    ```
    The result is transaction signed by the first node `my_tx.json`.

2. Sign the created transaction using the second node address:

    ```
    $ ./bin/neo-go wallet multisig sign -w .docker/wallets/wallet2.json --in my_tx.json --out my_tx2.json --addr Nbb1qkwcwNSBs9pAnrVVrnFbWnbWBk91U2
    ```
    Where
    - `-w .docker/wallets/wallet2.json` - path to the [wallet](https://github.com/nspcc-dev/neo-go/blob/master/.docker/wallets/wallet2.json) for the second node in private network
    - `--in my_tx.json` - previously created transfer transaction
    - `--out my_tx2.json` - output file for the signed transaction
    - `--addr Nbb1qkwcwNSBs9pAnrVVrnFbWnbWBk91U2` - multisig account to sign the transaction
    
    Enter the password `two`:
    ```
    Enter password to unlock wallet and sign the transaction
    Password >
    ```
    The result is transaction signed by both first and second nodes.

3. Sign the transaction using the third node address and push it to the chain:
    ```
    $ ./bin/neo-go wallet multisig sign -w ./.docker/wallets/wallet3.json --in my_tx2.json --out my_tx3.json --addr Nbb1qkwcwNSBs9pAnrVVrnFbWnbWBk91U2 -r http://localhost:20331
    ```
    Enter the password `three`:
    ```
    Enter password to unlock wallet and sign the transaction
    Password >
    ```
    The result is transaction signed by the first, second and third nodes and deployed to the chain.

4. Check the balance:

    Now you should have 29999999 GAS on the balance of `NMipL5VsNoLUBUJKPKLhxaEbPQVCZnyJyB` account.
    To check the transfer was successfully submitted use `getnep5transfers` RPC call:
    ```
    curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getnep5transfers", "params": ["NMipL5VsNoLUBUJKPKLhxaEbPQVCZnyJyB"] }' localhost:20331 | json_pp
    ```
    The result should look like the following:
    ```
    {
       "id" : 1,
       "jsonrpc" : "2.0",
       "result" : {
          "sent" : [],
          "address" : "NMipL5VsNoLUBUJKPKLhxaEbPQVCZnyJyB",
          "received" : [
             {
                "amount" : "29999999",
                "timestamp" : 1594055861999,
                "tx_hash" : "0xc6d51a1434065f7e0fa4a56993be4020d083303e9d79a49176cecf800ca136f7",
                "asset_hash" : "0x8c23f196d8a1bfd103a9dcb1f9ccf0c611377d3b",
                "transfer_address" : "Nbb1qkwcwNSBs9pAnrVVrnFbWnbWBk91U2",
                "block_index" : 7,
                "transfer_notify_index" : 0
             }
          ]
       }
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

Create configuration for it:
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
Enter account NMipL5VsNoLUBUJKPKLhxaEbPQVCZnyJyB password >
```

Result:
```
Sent deployment transaction 88a81d4acde6b302352e22ba5b0addcbdd4e5c284185c2b926930d83c5dc4128 for contract 28dbf93dc07a3d9b84ce6499132b874844784f9c
```

At this point your ‘Hello World’ contract is deployed and could be invoked. Let’s do it as a final step.

#### Step 4
Invoke contract.
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 28dbf93dc07a3d9b84ce6499132b874844784f9c main
```

Where
- `contract invokefunction` runs invoke with provided parameters
- `-r http://localhost:20331` defines RPC endpoint used for function call
- `-w my_wallet.json` is a wallet
- `28dbf93dc07a3d9b84ce6499132b874844784f9c` contract hash got as an output from the previous command (deployment in step 6)
- `main` - method to be called

Enter password `qwerty` for account:
```
Enter account NMipL5VsNoLUBUJKPKLhxaEbPQVCZnyJyB password >
```

Result:
```
Sent invocation transaction fa092821e3d5cc0a19ccc374f749f78e283771d0d734a348137d7f0a465bc308
```
In the console where you were running step #5 you will get:
```
2020-07-02T17:01:46.175+0300	INFO	runtime log	{"script": "9c4f784448872b139964ce849b3d7ac03df9db28", "logs": "\"Hello, world!\""}
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

RPC-server of started in step #5 neo-go node is available on `localhost:20331`, so let's try to perform several RPC calls.

#### GetRawTransaction
[GetRawTransaction](https://docs.neo.org/v3/docs/en-us/reference/rpc/latest-version/api/getrawtransaction.html) returns 
the corresponding transaction information, based on the specified hash value.

Request information about our deployment transaction:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getrawtransaction", "params": ["88a81d4acde6b302352e22ba5b0addcbdd4e5c284185c2b926930d83c5dc4128", 1] }' localhost:20331 | json_pp
```

Where
- `"jsonrpc": "2.0"` is protocol version
- `"id": 1` is id of current request
- `"method": "getrawtransaction"` is requested method
- `"params": ["88a81d4acde6b302352e22ba5b0addcbdd4e5c284185c2b926930d83c5dc4128", 1]` is an array of parameters, where
   - `88a81d4acde6b302352e22ba5b0addcbdd4e5c284185c2b926930d83c5dc4128` is deployment transaction hash
   - `1` is `verbose` parameter for detailed JSON string output
- `json_pp` just makes the JSON output prettier

Result:
```
{
   "jsonrpc" : "2.0",
   "id" : 1,
   "result" : {
      "script" : "DSkBeyJhYmkiOnsiaGFzaCI6IjB4MjhkYmY5M2RjMDdhM2Q5Yjg0Y2U2NDk5MTMyYjg3NDg0NDc4NGY5YyIsImVudHJ5UG9pbnQiOnsibmFtZSI6Ik1haW4iLCJwYXJhbWV0ZXJzIjpbXSwicmV0dXJuVHlwZSI6IlZvaWQifSwibWV0aG9kcyI6W10sImV2ZW50cyI6W119LCJncm91cHMiOltdLCJmZWF0dXJlcyI6eyJwYXlhYmxlIjpmYWxzZSwic3RvcmFnZSI6ZmFsc2V9LCJwZXJtaXNzaW9ucyI6W3siY29udHJhY3QiOiIqIiwibWV0aG9kcyI6IioifV0sInRydXN0cyI6W10sInNhZmVNZXRob2RzIjpbXSwiZXh0cmEiOm51bGx9DBYMDUhlbGxvLCB3b3JsZCFBz+dHliFAQc41LIU=",
      "blocktime" : 1594038738541,
      "confirmations" : 25,
      "sys_fee" : "31913180",
      "net_fee" : "589210",
      "size" : 489,
      "cosigners" : [],
      "hash" : "0x88a81d4acde6b302352e22ba5b0addcbdd4e5c284185c2b926930d83c5dc4128",
      "attributes" : [],
      "nonce" : 3862925904,
      "scripts" : [
         {
            "verification" : "DCEDtxJbrd1Yndtn9EiQ9Gy8B1aIOIYTheXMnRdnv3023lwLQQqQatQ=",
            "invocation" : "DEBSemvRhCt8IWB2hQVN3YpR8NzAPI10j1G/kGx3w66aC26sdyzxyptiskqpu2zrhFF8yLKuxbeYiAqhfy7KS0B2"
         }
      ],
      "version" : 0,
      "sender" : "NMipL5VsNoLUBUJKPKLhxaEbPQVCZnyJyB",
      "valid_until_block" : 52,
      "blockhash" : "0x3cf26d2eff8b27f54a0207bc245ebbd2dfcc62851949a9b612088bffdcf94c82"
   }
}
```

#### GetApplicationLog
[GetApplicationLog](https://docs.neo.org/v3/docs/en-us/reference/rpc/latest-version/api/getapplicationlog.html) returns the contract log based on the specified transaction id.

Request application log for invocation transaction from step #7:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["fa092821e3d5cc0a19ccc374f749f78e283771d0d734a348137d7f0a465bc308"] }' localhost:20331 | json_pp
```

With a single parameter:
- `fa092821e3d5cc0a19ccc374f749f78e283771d0d734a348137d7f0a465bc308` - invocation transaction hash from step #7

Result:
```
{
   "id" : 1,
   "jsonrpc" : "2.0",
   "result" : {
      "stack" : [
         {
            "value" : [],
            "type" : "Array"
         },
         {
            "type" : "ByteString",
            "value" : "bWFpbg=="
         }
      ],
      "trigger" : "Application",
      "gas_consumed" : "2007600",
      "vmstate" : "HALT",
      "txid" : "0xfa092821e3d5cc0a19ccc374f749f78e283771d0d734a348137d7f0a465bc308",
      "notifications" : []
   }
}
```

#### Other Useful RPC calls
```
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getversion", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getblockcount", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getconnectioncount", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getcontractstate", "params": ["0x9c33bbf2f5afbbc8fe271dd37508acd93573cffc"] }' localhost:20331
```

List of supported by neo-go node RPC commands you can find [here](https://github.com/nspcc-dev/neo-go/blob/master/docs/rpc.md#supported-methods).

### Storage smart contract

Let's take a look at the another smart contract example: [2-storage.go](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/2-storage.go).
This contract is quite simple and, as the previous one, doesn't take any arguments.
On the other hand, it is able to count the number of its own invocations by storing an integer value and increment it after each invocation.
We are interested in this contract as far as it's able to *store* values, i.e. it has a *storage* which can be shared within all contract invocations.

Unfortunately, we have to pay some GAS for storage usage, so it should be noted in the contract configuration [2-storage.yml](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/2-storage.yml) by the following parameter:
```
hasstorage: true
```
If you don't set the flag `hasstorage` to `true` value, you won't be able to use the storage in your contract.

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
Enter account NMipL5VsNoLUBUJKPKLhxaEbPQVCZnyJyB password >
```

Result:
```
Sent deployment transaction 86b19bab4bcca590de3dcdb245c4ef088943c4f08865d79e653ca81abda5d7f0 for contract 5776525e8c2c4ebccca99e71ee852e3d1bd073f9
```   

Which means that our contract was deployed and now we can invoke it.

#### Step #3
Let's invoke our contract. As far as we have never invoked this contract, there's no value in the storage, so the contract should create a new one (which is `1`) and put it into storage.
Let's check:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 5776525e8c2c4ebccca99e71ee852e3d1bd073f9 main
```
... enter the password `qwerty`:
```
Enter account NMipL5VsNoLUBUJKPKLhxaEbPQVCZnyJyB password >
```
Result:
```
Sent invocation transaction 74c1cabe482fa46c64f081d6b0af49b7d897758e98d10adea59e3b8ef78e6a38
```
To check the counter value, call `getapplicaionlog` RPC-call for the invocation transaction:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["74c1cabe482fa46c64f081d6b0af49b7d897758e98d10adea59e3b8ef78e6a38"] }' localhost:20331 | json_pp
```
The JSON result is:
```
{
   "jsonrpc" : "2.0",
   "id" : 1,
   "result" : {
      "vmstate" : "HALT",
      "gas_consumed" : "5131710",
      "trigger" : "Application",
      "txid" : "0x74c1cabe482fa46c64f081d6b0af49b7d897758e98d10adea59e3b8ef78e6a38",
      "stack" : [
         {
            "type" : "Array",
            "value" : []
         },
         {
            "type" : "ByteString",
            "value" : "bWFpbg=="
         },
         {
            "value" : "1",
            "type" : "Integer"
         }
      ],
      "notifications" : [
         {
            "contract" : "0x5776525e8c2c4ebccca99e71ee852e3d1bd073f9",
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
            "contract" : "0x5776525e8c2c4ebccca99e71ee852e3d1bd073f9",
            "state" : {
               "type" : "Array",
               "value" : [
                  {
                     "value" : "U3RvcmFnZSBrZXkgbm90IHlldCBzZXQuIFNldHRpbmcgdG8gMQ==",
                     "type" : "ByteString"
                  }
               ]
            }
         },
         {
            "state" : {
               "value" : [
                  {
                     "type" : "ByteString",
                     "value" : "TmV3IHZhbHVlIHdyaXR0ZW4gaW50byBzdG9yYWdl"
                  }
               ],
               "type" : "Array"
            },
            "contract" : "0x5776525e8c2c4ebccca99e71ee852e3d1bd073f9"
         }
      ]
   }
}
```
Pay attention to `notification` field. It contains messages, which where passed to `runtime.Notify` method.
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
  - `Storage key not yet set. Setting to 1` which was called when we realised that counter value is 0
  - `New value written into storage` which was called after the counter value was put in the storage.
  
The final part is `stack` field. This field contains all returned by the contract values, so here you can see integer value `1`,
which is the counter value denoted to the number of contract invocations.

#### Step #4
To ensure that all works as expected, let's invoke the contract one more time and check, whether the counter will be incremented: 
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 5776525e8c2c4ebccca99e71ee852e3d1bd073f9 main
```
... enter the password `qwerty`:
```
Enter account NMipL5VsNoLUBUJKPKLhxaEbPQVCZnyJyB password >
```
Result:
```
Sent invocation transaction 55ea69629975d3f97b9463e0ea1dddf7724c633795baea7c34372bb88698e487
```
To check the counter value, call `getapplicaionlog` RPC-call for the invocation transaction:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["55ea69629975d3f97b9463e0ea1dddf7724c633795baea7c34372bb88698e487"] }' localhost:20331 | json_pp
```
The JSON result is:
```
{
   "result" : {
      "gas_consumed" : "5211900",
      "trigger" : "Application",
      "notifications" : [
         {
            "contract" : "0x5776525e8c2c4ebccca99e71ee852e3d1bd073f9",
            "state" : {
               "value" : [
                  {
                     "type" : "ByteString",
                     "value" : "VmFsdWUgcmVhZCBmcm9tIHN0b3JhZ2U="
                  }
               ],
               "type" : "Array"
            }
         },
         {
            "contract" : "0x5776525e8c2c4ebccca99e71ee852e3d1bd073f9",
            "state" : {
               "value" : [
                  {
                     "type" : "ByteString",
                     "value" : "U3RvcmFnZSBrZXkgYWxyZWFkeSBzZXQuIEluY3JlbWVudGluZyBieSAx"
                  }
               ],
               "type" : "Array"
            }
         },
         {
            "state" : {
               "value" : [
                  {
                     "value" : "TmV3IHZhbHVlIHdyaXR0ZW4gaW50byBzdG9yYWdl",
                     "type" : "ByteString"
                  }
               ],
               "type" : "Array"
            },
            "contract" : "0x5776525e8c2c4ebccca99e71ee852e3d1bd073f9"
         }
      ],
      "vmstate" : "HALT",
      "txid" : "0x55ea69629975d3f97b9463e0ea1dddf7724c633795baea7c34372bb88698e487",
      "stack" : [
         {
            "type" : "Array",
            "value" : []
         },
         {
            "value" : "bWFpbg==",
            "type" : "ByteString"
         },
         {
            "type" : "Integer",
            "value" : "2"
         }
      ]
   },
   "jsonrpc" : "2.0",
   "id" : 1
}
```

The `stack` field contains now `2` integer value, so the counter was incremented as we expected.

## Workshop. Part 3
In this part we'll know about NEP5 token standard and try to write, deploy and invoke more complicated smart contract. 
Let’s go!

### NEP5
[NEP5](https://docs.neo.org/docs/en-us/sc/write/nep5.html) is a token standard for the Neo blockchain that provides systems with a generalized interaction mechanism for tokenized smart contracts.
The example with implementation of all required by the standard methods you can find in [nep5.go](https://github.com/nspcc-dev/neo-go/blob/master/examples/token/nep5/nep5.go)
 
Let's take a view on the example of smart contract with NEP5: [token.go](https://github.com/nspcc-dev/neo-go/blob/master/examples/token/token.go)
 
This smart contract initialises nep5 token interface and takes operation string as a parameter, which is one of:
- `name` returns name of created nep5 token 
- `symbol` returns ticker symbol of the token
- `decimals` returns amount of decimals for the token
- `totalSupply` returns total token * multiplier
- `balanceOf` returns the token balance of a specific address and requires additional argument:
  - `holder` which is requested address
- `transfer` transfers token from one user to another and requires additional arguments:
  - `from` is account which you'd like to transfer tokens from
  - `to` is account which you'd like to transfer tokens to
  - `amount` is the amount of token to transfer
- `mint` supplies initial amount of token to account and requires additional arguments:
  - `to` is account address which you'd like to transfer initial token to
Let's perform several operations with our contract.

#### Step #1
To compile [token.go](https://github.com/nspcc-dev/neo-go/blob/master/examples/token/token.go)
you can use [configuration](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/3-token.yml).
As far as our contract uses storage, the flag `hasstorage` should be set to `true`:
```
hasstorage: true
```
Compile smart contract:
```
$ ./bin/neo-go contract compile -i examples/token/token.go -c 3-token.yml -m examples/token/token.manifest.json
```

Deploy smart contract:
```
$ ./bin/neo-go contract deploy -i examples/token/token.nef -manifest examples/token/token.manifest.json -r http://localhost:20331 -w my_wallet.json
```
... enter the password `qwerty`:
```
Enter account NMipL5VsNoLUBUJKPKLhxaEbPQVCZnyJyB password >
```

Result:
```
Sent deployment transaction 110c2b4a5c4abe79e7b19b7bdf737041fe73d45f3679911ac44932e9cbe978e8 for contract 081480331bcd3183cd6c4d0c035717bba1b4daae
```   

Which means that our contract was deployed and now we can invoke it.

#### Step #2
Let's invoke the contract to perform different operations.

To start with, query `Name` of the created nep5 token:

```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 081480331bcd3183cd6c4d0c035717bba1b4daae name
```                                                                   
Where
- `081480331bcd3183cd6c4d0c035717bba1b4daae` is our contract hash from step #1
- `name` is operation string which was described earlier and returns token name

... and don't forget the password of your account `qwerty`.

Result:
```
Sent invocation transaction e3a9b58c98b47b005ee68ae613942d5019617cc71c6ef4377649e800ffe9dc10
```                                                                                         
Now, let's take a detailed look at this invocation transaction with `getapplicationlog` RPC call:

```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["e3a9b58c98b47b005ee68ae613942d5019617cc71c6ef4377649e800ffe9dc10"] }' localhost:20331 | json_pp
```               

Result:
```
{
   "jsonrpc" : "2.0",
   "result" : {
      "txid" : "0xe3a9b58c98b47b005ee68ae613942d5019617cc71c6ef4377649e800ffe9dc10",
      "notifications" : [],
      "trigger" : "Application",
      "gas_consumed" : "3041410",
      "vmstate" : "HALT",
      "stack" : [
         {
            "type" : "ByteString",
            "value" : "QXdlc29tZSBORU8gVG9rZW4="
         }
      ]
   },
   "id" : 1
}
```

At least, you can see that `stack` field of JSON result is not empty: it contains base64 byte array with the name of our token.

Following commands able you to get some additional information about token:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 081480331bcd3183cd6c4d0c035717bba1b4daae symbol
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 081480331bcd3183cd6c4d0c035717bba1b4daae decimals
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 081480331bcd3183cd6c4d0c035717bba1b4daae totalSupply
```

#### Step #3

Now it's time for more interesting things. First of all, let's check the balance of nep5 token on our account by using `balanceOf`:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 081480331bcd3183cd6c4d0c035717bba1b4daae balanceOf NMipL5VsNoLUBUJKPKLhxaEbPQVCZnyJyB
```                             
... with `qwerty` password. The result is:
```
Sent invocation transaction 4e6efb13da6e396030e85936c17911b0bc5fe684992deae9c8889b4c3c84d774
```
And take a closer look at the transaction's details with `getapplicationlog` RPC-call:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["4e6efb13da6e396030e85936c17911b0bc5fe684992deae9c8889b4c3c84d774"] }' localhost:20331 | json_pp
```
Result:
```
{
   "result" : {
      "gas_consumed" : "4147990",
      "notifications" : [],
      "vmstate" : "HALT",
      "stack" : [
         {
            "value" : "0",
            "type" : "Integer"
         }
      ],
      "txid" : "0x4e6efb13da6e396030e85936c17911b0bc5fe684992deae9c8889b4c3c84d774",
      "trigger" : "Application"
   },
   "jsonrpc" : "2.0",
   "id" : 1
}
``` 
As far as `stack` field contains integer value `0`, we have no token on the balance. But don't worry about that. Just follow the next step.

#### Step #4

Before we are able to start using our token (e.g. transfer it to someone else), we have to *mint* it.
In other words, we should transfer all available amount of token (total supply) to someone's account.
There's a special function for this purpose in our contract - `Mint` function. However, this function
uses `CheckWitness` runtime syscall to check whether the caller of the contract is the owner and authorized
to manage initial supply of tokens. That's the purpose of transaction's *cosigners*: checking given hash
against the values provided in the list of cosigners. To pass this check we should add our account to
transaction's cosigners list. So let's mint token to our address:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 081480331bcd3183cd6c4d0c035717bba1b4daae mint NMipL5VsNoLUBUJKPKLhxaEbPQVCZnyJyB -- f0a33acbf435417f3d2c18445d607d1e6f48d413
```
Where:
- `--` is a special delimiter of transaction's cosigners list
- `f0a33acbf435417f3d2c18445d607d1e6f48d413` is the cosigner itself (which is hex-encoded LE representation of our `NMipL5VsNoLUBUJKPKLhxaEbPQVCZnyJyB` account)

... with `qwerty` pass. The result:
``` 
Sent invocation transaction 828ee09b1e875cdf59f4cbaa0c68941acd4fbbf17c30854dbf711b25cb0bce7c
```
`getapplicationlog` RPC-call for this transaction tells us the following:
```
{
   "jsonrpc" : "2.0",
   "result" : {
      "trigger" : "Application",
      "txid" : "0x828ee09b1e875cdf59f4cbaa0c68941acd4fbbf17c30854dbf711b25cb0bce7c",
      "vmstate" : "HALT",
      "notifications" : [
         {
            "contract" : "0x081480331bcd3183cd6c4d0c035717bba1b4daae",
            "state" : {
               "type" : "Array",
               "value" : [
                  {
                     "value" : "dHJhbnNmZXI=",
                     "type" : "ByteString"
                  },
                  {
                     "value" : "",
                     "type" : "ByteString"
                  },
                  {
                     "value" : "E9RIbx59YF1EGCw9f0E19Ms6o/A=",
                     "type" : "ByteString"
                  },
                  {
                     "value" : "1100000000000000",
                     "type" : "Integer"
                  }
               ]
            }
         }
      ],
      "stack" : [
         {
            "value" : "1",
            "type" : "Integer"
         }
      ],
      "gas_consumed" : "6982010"
   },
   "id" : 1
}
```
Here we have `1` at the `stack` field, which means that token was successfully minted.
Let's just ensure that by querying `balanceOf` one more time:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 081480331bcd3183cd6c4d0c035717bba1b4daae balanceOf NMipL5VsNoLUBUJKPKLhxaEbPQVCZnyJyB
```
... with `qwerty` pass. The result:
``` 
Sent invocation transaction 6a7ceac2a294d6136a69e22ebf8dc40ccc3129d3be8ece965b1d9b852d9f7d0c
```
... with the following `getapplicationlog` JSON message:
```
{
   "result" : {
      "stack" : [
         {
            "type" : "Integer",
            "value" : "1100000000000000"
         }
      ],
      "txid" : "0x6a7ceac2a294d6136a69e22ebf8dc40ccc3129d3be8ece965b1d9b852d9f7d0c",
      "trigger" : "Application",
      "gas_consumed" : "4228110",
      "vmstate" : "HALT",
      "notifications" : []
   },
   "id" : 1,
   "jsonrpc" : "2.0"
}
```
Now we can see integer value at the `stack` field, so `1100000000000000` is the nep5 token balance of our account.

Note, that token can be minted only once.

#### Step #5

After we are done with minting, it's possible to transfer token to someone else.
Let's transfer 5 tokens from our account to `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm` with `transfer` call:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 081480331bcd3183cd6c4d0c035717bba1b4daae transfer NMipL5VsNoLUBUJKPKLhxaEbPQVCZnyJyB NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm 500000000 -- f0a33acbf435417f3d2c18445d607d1e6f48d413
```
... with password `qwerty` and following result:
``` 
Sent invocation transaction 6b87a6ce22af49734f4d13c114d7449171f2298f7959e88c9b4aac071b0886a8
```
Our favourite `getapplicationlog` RPC-call tells us:
```
{
   "jsonrpc" : "2.0",
   "id" : 1,
   "result" : {
      "stack" : [
         {
            "value" : "1",
            "type" : "Integer"
         }
      ],
      "vmstate" : "HALT",
      "trigger" : "Application",
      "txid" : "0x6b87a6ce22af49734f4d13c114d7449171f2298f7959e88c9b4aac071b0886a8",
      "gas_consumed" : "7565310",
      "notifications" : [
         {
            "contract" : "0x081480331bcd3183cd6c4d0c035717bba1b4daae",
            "state" : {
               "value" : [
                  {
                     "value" : "dHJhbnNmZXI=",
                     "type" : "ByteString"
                  },
                  {
                     "type" : "ByteString",
                     "value" : "E9RIbx59YF1EGCw9f0E19Ms6o/A="
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
            }
         }
      ]
   }
}
```
Note, that `stack` field contains `1`, which means that token was successfully transferred.
Let's now check the balance of `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm` account to ensure that the amount of token on that account = 5:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 081480331bcd3183cd6c4d0c035717bba1b4daae balanceOf NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm
```
The `getapplicationlog` RPC-call for this transaction tells us the following:
```
{
   "id" : 1,
   "result" : {
      "gas_consumed" : "4228110",
      "txid" : "0x41e24302e393012b2a23859367adf50c931d66ded48926d3e2b9f8b2acc00f86",
      "stack" : [
         {
            "value" : "500000000",
            "type" : "Integer"
         }
      ],
      "trigger" : "Application",
      "vmstate" : "HALT",
      "notifications" : []
   },
   "jsonrpc" : "2.0"
}
```
Here we are! There are exactly 5 tokens at the `stack` field. You can also ensure that these 5 tokens were debited from `NMipL5VsNoLUBUJKPKLhxaEbPQVCZnyJyB` account by using `balanceOf` method.

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
- `delete` deletes specified domain from the storage. The arguments:
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
Enter account NMipL5VsNoLUBUJKPKLhxaEbPQVCZnyJyB password >
```

Result:
```
Sent deployment transaction 242a04bac1b0803874c2cc14b6aeb98e7f7caa77878edb55d78fcc0d7448e510 for contract 9ca0c9181131c52fb39470f57e64d72a83115b6f
```   
You know, what it means :)

#### Step #3

Invoke the contract to register domain with name `my_first_domain`:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 9ca0c9181131c52fb39470f57e64d72a83115b6f register my_first_domain NMipL5VsNoLUBUJKPKLhxaEbPQVCZnyJyB -- f0a33acbf435417f3d2c18445d607d1e6f48d413
```
... the strongest password in the world, guess: `qwerty`
```
Enter account NMipL5VsNoLUBUJKPKLhxaEbPQVCZnyJyB password >
```
Result:
```
Sent invocation transaction 08f14a3d439d0f024e281a80800717478dd7afa2b587f84b68e3e5f8b345afbf
```
Also you can see the log message in the console, where you run neo-go node:
```
2020-07-06T16:50:37.713+0300	INFO	runtime log	{"script": "6f5b11832ad7647ef57094b32fc5311118c9a09c", "logs": "\"RegisterDomain: my_first_domain\""}
```
Well, that's ok. Let's check now, whether our domain was registered with `getapplicationlog` RPC-call:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["08f14a3d439d0f024e281a80800717478dd7afa2b587f84b68e3e5f8b345afbf"] }' localhost:20331 | json_pp
```
The result is:
```
{
   "jsonrpc" : "2.0",
   "result" : {
      "txid" : "0x08f14a3d439d0f024e281a80800717478dd7afa2b587f84b68e3e5f8b345afbf",
      "vmstate" : "HALT",
      "stack" : [
         {
            "value" : "1",
            "type" : "Integer"
         }
      ],
      "notifications" : [
         {
            "state" : {
               "value" : [
                  {
                     "value" : "cmVnaXN0ZXJlZA==",
                     "type" : "ByteString"
                  },
                  {
                     "type" : "ByteString",
                     "value" : "E9RIbx59YF1EGCw9f0E19Ms6o/A="
                  },
                  {
                     "type" : "ByteString",
                     "value" : "bXlfZmlyc3RfZG9tYWlu"
                  }
               ],
               "type" : "Array"
            },
            "contract" : "0x9ca0c9181131c52fb39470f57e64d72a83115b6f"
         }
      ],
      "gas_consumed" : "7203080",
      "trigger" : "Application"
   },
   "id" : 1
}
```
Especially, we're interested in two fields of the json:

First one is `notifications` field, which contains 3 values:
- `cmVnaXN0ZXJlZA==` byte string in base64 representation which can be decoded into `registered`,
- `bXlfZmlyc3RfZG9tYWlu` byte string, which can be decoded to `my_first_domain` - our domain's name
- `E9RIbx59YF1EGCw9f0E19Ms6o/A=` byte array, which can be decoded to the account address `NMipL5VsNoLUBUJKPKLhxaEbPQVCZnyJyB`.

The second field is `stack` with `1` value, which was returned by the smart contract.

All of these values let us be sure that our domain was successfully registered.  

#### Step #4

Invoke the contract to query the address information our `my_first_domain` domain:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 9ca0c9181131c52fb39470f57e64d72a83115b6f query my_first_domain
```
... the pass `qwerty`:
```
Enter account NMipL5VsNoLUBUJKPKLhxaEbPQVCZnyJyB password >
```
Result:
```
Sent invocation transaction 061e429962e0fdbfc11a7e89b6fe72f64a8815bd39ae0cd62f1a93ccfaef32d4
```
and log-message:
```
2020-07-06T17:02:10.782+0300	INFO	runtime log	{"script": "6f5b11832ad7647ef57094b32fc5311118c9a09c", "logs": "\"QueryDomain: my_first_domain\""}
```
Let's check this transaction with `getapplicationlog` RPC call:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["061e429962e0fdbfc11a7e89b6fe72f64a8815bd39ae0cd62f1a93ccfaef32d4"] }' localhost:20331 | json_pp
```
... which gives us the following result:
```
{
   "result" : {
      "vmstate" : "HALT",
      "stack" : [
         {
            "type" : "ByteString",
            "value" : "E9RIbx59YF1EGCw9f0E19Ms6o/A="
         }
      ],
      "gas_consumed" : "4893900",
      "notifications" : [],
      "txid" : "0x061e429962e0fdbfc11a7e89b6fe72f64a8815bd39ae0cd62f1a93ccfaef32d4",
      "trigger" : "Application"
   },
   "id" : 1,
   "jsonrpc" : "2.0"
}
```

with base64 interpretation of our account address `NMipL5VsNoLUBUJKPKLhxaEbPQVCZnyJyB` on the stack, which means that domain `my_first_domain` was registered by the owner with received account address.

#### Step #5

Invoke the contract to transfer domain to the other account (e.g. account with `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm` address):
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 9ca0c9181131c52fb39470f57e64d72a83115b6f transfer my_first_domain NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm -- f0a33acbf435417f3d2c18445d607d1e6f48d413
```
... the password: `qwerty`
```
Enter account NMipL5VsNoLUBUJKPKLhxaEbPQVCZnyJyB password >
```
Result:
```
Sent invocation transaction 5c1da38f000abec53dc135dc404e768a1613749c3377b73d960782fac464d941
```
and log-message:
```
2020-07-06T17:10:08.521+0300	INFO	runtime log	{"script": "6f5b11832ad7647ef57094b32fc5311118c9a09c", "logs": "\"TransferDomain: my_first_domain\""}
```
Perfect. And `getapplicationlog` RPC-call...
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["5c1da38f000abec53dc135dc404e768a1613749c3377b73d960782fac464d941"] }' localhost:20331 | json_pp
```
... tells us:
```
{
   "result" : {
      "gas_consumed" : "6390910",
      "stack" : [
         {
            "value" : "1",
            "type" : "Integer"
         }
      ],
      "trigger" : "Application",
      "vmstate" : "HALT",
      "txid" : "0x5c1da38f000abec53dc135dc404e768a1613749c3377b73d960782fac464d941",
      "notifications" : [
         {
            "state" : {
               "type" : "Array",
               "value" : [
                  {
                     "value" : "ZGVsZXRlZA==",
                     "type" : "ByteString"
                  },
                  {
                     "type" : "ByteString",
                     "value" : "E9RIbx59YF1EGCw9f0E19Ms6o/A="
                  },
                  {
                     "type" : "ByteString",
                     "value" : "bXlfZmlyc3RfZG9tYWlu"
                  }
               ]
            },
            "contract" : "0x9ca0c9181131c52fb39470f57e64d72a83115b6f"
         },
         {
            "state" : {
               "value" : [
                  {
                     "value" : "cmVnaXN0ZXJlZA==",
                     "type" : "ByteString"
                  },
                  {
                     "type" : "ByteString",
                     "value" : "50l6vFaauRKm8hPVkr3Aw2CeHQs="
                  },
                  {
                     "type" : "ByteString",
                     "value" : "bXlfZmlyc3RfZG9tYWlu"
                  }
               ],
               "type" : "Array"
            },
            "contract" : "0x9ca0c9181131c52fb39470f57e64d72a83115b6f"
         }
      ]
   },
   "jsonrpc" : "2.0",
   "id" : 1
}
```
Quite a detailed one. The `notifications` field contains two arrays:
- First one with `ZGVsZXRlZA==` byte string, which is `deleted` with additional information (domain `my_first_domain` was deleted from account `NMipL5VsNoLUBUJKPKLhxaEbPQVCZnyJyB`),
- Second one with `cmVnaXN0ZXJlZA==` byte string, which is `registered` with additional information  (domain `my_first_domain` was registered with account `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm`).
The `stack` field contains `1` value, which means that our domain was successfully transferred.

#### Step #6

The last call is `delete`, so you can try to create the other domain, e.g. `my_second_domain` and then remove it from storage with:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 9ca0c9181131c52fb39470f57e64d72a83115b6f delete my_second_domain -- f0a33acbf435417f3d2c18445d607d1e6f48d413
```

Thank you!

### Useful links

* [Our basic tutorial on Medium](https://medium.com/@neospcc/%D1%81%D0%BC%D0%B0%D1%80%D1%82-%D0%BA%D0%BE%D0%BD%D1%82%D1%80%D0%B0%D0%BA%D1%82-%D0%B4%D0%BB%D1%8F-neo-769139352b65)
* [Using NEO Blockchain Toolkit](https://medium.com/@neospcc/neogo-adds-support-for-neo-blockchain-toolkit-673ea914f661)
* [NEO documentation](https://docs.neo.org/)
* [NEO github](https://github.com/neo-project/neo/)
* [NEO-GO github](https://github.com/nspcc-dev/neo-go)
