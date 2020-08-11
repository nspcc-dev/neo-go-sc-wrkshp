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
        $ ./bin/neo-go wallet nep5 transfer -w .docker/wallets/wallet1.json --out my_tx.json -r http://localhost:20331 --from NUVPACMnKFhpuHjsRjhUvXz1XhqfGZYVtY --to NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt --token gas --amount 29999999
    ``` 
    Where
    - `./bin/neo-go` runs neo-go
    - `wallet nep5 transfer` - command with arguments in [neo-go](https://github.com/nspcc-dev/neo-go/blob/master/cli/wallet/nep5.go#L108)
    - `-w .docker/wallets/wallet1.json` - path to the [wallet](https://github.com/nspcc-dev/neo-go/blob/master/.docker/wallets/wallet1.json) for the first node in the private network
    - `--out my_tx.json` - output file for the signed transaction
    - `-r http://localhost:20331` - RPC node endpoint
    - `--from NUVPACMnKFhpuHjsRjhUvXz1XhqfGZYVtY` - multisig account to transfer GAS from
    - `--to NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt` - our account from the [wallet](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/my_wallet.json)
    - `--token gas` - transferred token name, which is GAS
    - `--amount 29999999` - amount of GAS to transfer
    
    Enter the password `one`:
    ```
    Password >
    ```
    The result is transaction signed by the first node `my_tx.json`.

2. Sign the created transaction using the second node address:

    ```
    $ ./bin/neo-go wallet multisig sign -w .docker/wallets/wallet2.json --in my_tx.json --out my_tx2.json --addr NUVPACMnKFhpuHjsRjhUvXz1XhqfGZYVtY
    ```
    Where
    - `-w .docker/wallets/wallet2.json` - path to the [wallet](https://github.com/nspcc-dev/neo-go/blob/master/.docker/wallets/wallet2.json) for the second node in private network
    - `--in my_tx.json` - previously created transfer transaction
    - `--out my_tx2.json` - output file for the signed transaction
    - `--addr NUVPACMnKFhpuHjsRjhUvXz1XhqfGZYVtY` - multisig account to sign the transaction
    
    Enter the password `two`:
    ```
    Enter password to unlock wallet and sign the transaction
    Password >
    ```
    The result is transaction signed by both first and second nodes.

3. Sign the transaction using the third node address and push it to the chain:
    ```
    $ ./bin/neo-go wallet multisig sign -w ./.docker/wallets/wallet3.json --in my_tx2.json --out my_tx3.json --addr NUVPACMnKFhpuHjsRjhUvXz1XhqfGZYVtY -r http://localhost:20331
    ```
    Enter the password `three`:
    ```
    Enter password to unlock wallet and sign the transaction
    Password >
    ```
    The result is transaction signed by the first, second and third nodes and deployed to the chain.

4. Check the balance:

    Now you should have 29999999 GAS on the balance of `NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt` account.
    To check the transfer was successfully submitted use `getnep5transfers` RPC call:
    ```
    curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getnep5transfers", "params": ["NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt"] }' localhost:20331 | json_pp
    ```
    The result should look like the following:
    ```
    {
       "jsonrpc" : "2.0",
       "id" : 1,
       "result" : {
          "received" : [
             {
                "amount" : "29999999",
                "transferaddress" : "NUVPACMnKFhpuHjsRjhUvXz1XhqfGZYVtY",
                "blockindex" : 13,
                "assethash" : "0x668e0c1f9d7b70a99dd9e06eadd4c784d641afbc",
                "timestamp" : 1597143962082,
                "txhash" : "0x19a4a1e58a48c9d9c65c451968a3a900e67915cf37ccd8edaa67e0866f1b5d04",
                "transfernotifyindex" : 0
             }
          ],
          "address" : "NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt",
          "sent" : []
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
Sent deployment transaction 8fc1dee25649d4292ac2cb19e7be2e65288ef5fab0add933b3f54c31477852dd for contract 28dbf93dc07a3d9b84ce6499132b874844784f9c
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
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```

Result:
```
Sent invocation transaction 5025daa84c6c41b672444e5f26fc125811900309c84295b35b2db3d70bab18b8
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

RPC-server of started in step #3 neo-go node is available on `localhost:20331`, so let's try to perform several RPC calls.

#### GetRawTransaction
[GetRawTransaction](https://docs.neo.org/v3/docs/en-us/reference/rpc/latest-version/api/getrawtransaction.html) returns 
the corresponding transaction information, based on the specified hash value.

Request information about our deployment transaction:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getrawtransaction", "params": ["8fc1dee25649d4292ac2cb19e7be2e65288ef5fab0add933b3f54c31477852dd", 1] }' localhost:20331 | json_pp
```

Where
- `"jsonrpc": "2.0"` is protocol version
- `"id": 1` is id of current request
- `"method": "getrawtransaction"` is requested method
- `"params": ["8fc1dee25649d4292ac2cb19e7be2e65288ef5fab0add933b3f54c31477852dd", 1]` is an array of parameters, where
   - `8fc1dee25649d4292ac2cb19e7be2e65288ef5fab0add933b3f54c31477852dd` is deployment transaction hash
   - `1` is `verbose` parameter for detailed JSON string output
- `json_pp` just makes the JSON output prettier

Result:
```
{
   "result" : {
      "attributes" : [],
      "blocktime" : 1597144100686,
      "size" : 511,
      "nonce" : 3765129059,
      "blockhash" : "0xef1da7ae391bc0ff88b061adced2e8b9fcb7f4fb6c27c65eef230b2a875aaa59",
      "sender" : "NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt",
      "signers" : [
         {
            "account" : "0xf8d79b436e9d3a54a38cc877d182b71e381d7a5c",
            "scopes" : "FeeOnly"
         }
      ],
      "witnesses" : [
         {
            "invocation" : "DEC4KhicN9a3V+sFm7UZfBntVmDY7QOWWpKGRLBjPXjx4m0epLmPQcpq722Y7aR1knqnLDNdcJebKqZ0FyB9tkWT",
            "verification" : "DCECqIHFxkXMSrGdFGesF8hU+x8Y9RaRmsEOWyRK4B9eK30LQZVEDXg="
         }
      ],
      "confirmations" : 23,
      "validuntilblock" : 22,
      "hash" : "0x8fc1dee25649d4292ac2cb19e7be2e65288ef5fab0add933b3f54c31477852dd",
      "netfee" : "2620420",
      "version" : 0,
      "vmstate" : "HALT",
      "sysfee" : "34013180",
      "script" : "DT4BeyJhYmkiOnsiaGFzaCI6IjB4MjhkYmY5M2RjMDdhM2Q5Yjg0Y2U2NDk5MTMyYjg3NDg0NDc4NGY5YyIsIm1ldGhvZHMiOlt7Im5hbWUiOiJtYWluIiwib2Zmc2V0IjowLCJwYXJhbWV0ZXJzIjpbXSwicmV0dXJudHlwZSI6IlZvaWQifV0sImV2ZW50cyI6W119LCJncm91cHMiOltdLCJmZWF0dXJlcyI6eyJwYXlhYmxlIjpmYWxzZSwic3RvcmFnZSI6ZmFsc2V9LCJwZXJtaXNzaW9ucyI6W3siY29udHJhY3QiOiIqIiwibWV0aG9kcyI6IioifV0sInN1cHBvcnRlZHN0YW5kYXJkcyI6W10sInRydXN0cyI6W10sInNhZmVtZXRob2RzIjpbXSwiZXh0cmEiOm51bGx9DBYMDUhlbGxvLCB3b3JsZCFBz+dHliFAQc41LIU="
   },
   "jsonrpc" : "2.0",
   "id" : 1
}
```

#### GetApplicationLog
[GetApplicationLog](https://docs.neo.org/v3/docs/en-us/reference/rpc/latest-version/api/getapplicationlog.html) returns the contract log based on the specified transaction id.

Request application log for invocation transaction from step #4:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["5025daa84c6c41b672444e5f26fc125811900309c84295b35b2db3d70bab18b8"] }' localhost:20331 | json_pp
```

With a single parameter:
- `5025daa84c6c41b672444e5f26fc125811900309c84295b35b2db3d70bab18b8` - invocation transaction hash from step #7

Result:
```
{
   "result" : {
      "txid" : "0x5025daa84c6c41b672444e5f26fc125811900309c84295b35b2db3d70bab18b8",
      "gasconsumed" : "2007600",
      "trigger" : "Application",
      "vmstate" : "HALT",
      "notifications" : [],
      "stack" : [
         {
            "type" : "Any"
         }
      ]
   },
   "id" : 1,
   "jsonrpc" : "2.0"
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
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```

Result:
```
Sent deployment transaction e3e61ca9d713bddda110a002504f9d5f94f871a15675c2e514d08efb1521ea1b for contract 79edf20dc8ee247981756787a638e9679026c16c
```   

Which means that our contract was deployed and now we can invoke it.

#### Step #3
Let's invoke our contract. As far as we have never invoked this contract, there's no value in the storage, so the contract should create a new one (which is `1`) and put it into storage.
Let's check:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 79edf20dc8ee247981756787a638e9679026c16c main
```
... enter the password `qwerty`:
```
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```
Result:
```
Sent invocation transaction 7f658d0008f7c9bb3f05d655450ead36b9e5d7a2ad0f9e5f7895a43579d2616d
```
To check the counter value, call `getapplicaionlog` RPC-call for the invocation transaction:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["7f658d0008f7c9bb3f05d655450ead36b9e5d7a2ad0f9e5f7895a43579d2616d"] }' localhost:20331 | json_pp
```
The JSON result is:
```
{
   "result" : {
      "stack" : [
         {
            "value" : "1",
            "type" : "Integer"
         }
      ],
      "notifications" : [
         {
            "contract" : "0x79edf20dc8ee247981756787a638e9679026c16c",
            "eventname" : "info",
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
            "eventname" : "info",
            "state" : {
               "type" : "Array",
               "value" : [
                  {
                     "value" : "U3RvcmFnZSBrZXkgbm90IHlldCBzZXQuIFNldHRpbmcgdG8gMQ==",
                     "type" : "ByteString"
                  }
               ]
            },
            "contract" : "0x79edf20dc8ee247981756787a638e9679026c16c"
         },
         {
            "contract" : "0x79edf20dc8ee247981756787a638e9679026c16c",
            "eventname" : "info",
            "state" : {
               "type" : "Array",
               "value" : [
                  {
                     "value" : "TmV3IHZhbHVlIHdyaXR0ZW4gaW50byBzdG9yYWdl",
                     "type" : "ByteString"
                  }
               ]
            }
         }
      ],
      "gasconsumed" : "5212830",
      "txid" : "0x7f658d0008f7c9bb3f05d655450ead36b9e5d7a2ad0f9e5f7895a43579d2616d",
      "trigger" : "Application",
      "vmstate" : "HALT"
   },
   "id" : 1,
   "jsonrpc" : "2.0"
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
  - `Storage key not yet set. Setting to 1` which was called when we realised that counter value is 0
  - `New value written into storage` which was called after the counter value was put in the storage.
  
The final part is `stack` field. This field contains all returned by the contract values, so here you can see integer value `1`,
which is the counter value denoted to the number of contract invocations.

#### Step #4
To ensure that all works as expected, let's invoke the contract one more time and check, whether the counter will be incremented: 
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 79edf20dc8ee247981756787a638e9679026c16c main
```
... enter the password `qwerty`:
```
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```
Result:
```
Sent invocation transaction bbe338b20228e7067c71da2137d14ac7cbd01158469979d882982aee7f41e316
```
To check the counter value, call `getapplicaionlog` RPC-call for the invocation transaction:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["bbe338b20228e7067c71da2137d14ac7cbd01158469979d882982aee7f41e316"] }' localhost:20331 | json_pp
```
The JSON result is:
```
{
   "result" : {
      "trigger" : "Application",
      "txid" : "0xbbe338b20228e7067c71da2137d14ac7cbd01158469979d882982aee7f41e316",
      "gasconsumed" : "5293020",
      "notifications" : [
         {
            "state" : {
               "type" : "Array",
               "value" : [
                  {
                     "value" : "VmFsdWUgcmVhZCBmcm9tIHN0b3JhZ2U=",
                     "type" : "ByteString"
                  }
               ]
            },
            "contract" : "0x79edf20dc8ee247981756787a638e9679026c16c",
            "eventname" : "info"
         },
         {
            "eventname" : "info",
            "contract" : "0x79edf20dc8ee247981756787a638e9679026c16c",
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
               "type" : "Array",
               "value" : [
                  {
                     "type" : "ByteString",
                     "value" : "TmV3IHZhbHVlIHdyaXR0ZW4gaW50byBzdG9yYWdl"
                  }
               ]
            },
            "contract" : "0x79edf20dc8ee247981756787a638e9679026c16c",
            "eventname" : "info"
         }
      ],
      "vmstate" : "HALT",
      "stack" : [
         {
            "value" : "2",
            "type" : "Integer"
         }
      ]
   },
   "id" : 1,
   "jsonrpc" : "2.0"
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
you can use [configuration](https://github.com/nspcc-dev/neo-go/blob/master/examples/token/token.yml).
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
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```

Result:
```
Sent deployment transaction d275df152996bdab279efb0942b15fd26c6a2b69965cdc77540f401f6a91aedc for contract 4f84a3b4a32125056b5b2bf5b6c7addc3616985f
```   

Which means that our contract was deployed and now we can invoke it.

#### Step #2
Let's invoke the contract to perform different operations.

To start with, query `Name` of the created nep5 token:

```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 4f84a3b4a32125056b5b2bf5b6c7addc3616985f name
```                                                                   
Where
- `4f84a3b4a32125056b5b2bf5b6c7addc3616985f` is our contract hash from step #1
- `name` is operation string which was described earlier and returns token name

... and don't forget the password of your account `qwerty`.

Result:
```
Sent invocation transaction de56bd31d785e6656a5eb20f48fd94c98d19700a9de4e721aebf68ba585490f2
```                                                                                         
Now, let's take a detailed look at this invocation transaction with `getapplicationlog` RPC call:

```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["de56bd31d785e6656a5eb20f48fd94c98d19700a9de4e721aebf68ba585490f2"] }' localhost:20331 | json_pp
```               

Result:
```
{
   "result" : {
      "stack" : [
         {
            "type" : "ByteString",
            "value" : "QXdlc29tZSBORU8gVG9rZW4="
         }
      ],
      "txid" : "0xde56bd31d785e6656a5eb20f48fd94c98d19700a9de4e721aebf68ba585490f2",
      "trigger" : "Application",
      "vmstate" : "HALT",
      "notifications" : [],
      "gasconsumed" : "4647080"
   },
   "jsonrpc" : "2.0",
   "id" : 1
}
```

At least, you can see that `stack` field of JSON result is not empty: it contains base64 byte array with the name of our token.

Following commands able you to get some additional information about token:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 4f84a3b4a32125056b5b2bf5b6c7addc3616985f symbol
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 4f84a3b4a32125056b5b2bf5b6c7addc3616985f decimals
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 4f84a3b4a32125056b5b2bf5b6c7addc3616985f totalSupply
```

#### Step #3

Now it's time for more interesting things. First of all, let's check the balance of nep5 token on our account by using `balanceOf`:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 4f84a3b4a32125056b5b2bf5b6c7addc3616985f balanceOf NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt
```                             
... with `qwerty` password. The result is:
```
Sent invocation transaction de1a69a6c0b5389458f35c44e28d1bea0a0ab4fe7b4ad0013e34b36d0995069a
```
And take a closer look at the transaction's details with `getapplicationlog` RPC-call:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["de1a69a6c0b5389458f35c44e28d1bea0a0ab4fe7b4ad0013e34b36d0995069a"] }' localhost:20331 | json_pp
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
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 4f84a3b4a32125056b5b2bf5b6c7addc3616985f mint NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt -- f8d79b436e9d3a54a38cc877d182b71e381d7a5c:CalledByEntry
```
Where:
- `--` is a special delimiter of transaction's cosigners list
- `f8d79b436e9d3a54a38cc877d182b71e381d7a5c` is the signer itself (which is hex-encoded LE representation of our `NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt` account)

... with `qwerty` pass. The result:
``` 
Sent invocation transaction 2f4eff5532808080ae4c6a1a844c32a53015ca3cad5c308c5bcf3d8261399d06
```
`getapplicationlog` RPC-call for this transaction tells us the following:
```
{
   "result" : {
      "stack" : [
         {
            "value" : "1",
            "type" : "Integer"
         }
      ],
      "vmstate" : "HALT",
      "notifications" : [
         {
            "contract" : "0x4f84a3b4a32125056b5b2bf5b6c7addc3616985f",
            "state" : {
               "value" : [
                  {
                     "value" : "",
                     "type" : "ByteString"
                  },
                  {
                     "value" : "XHodOB63gtF3yIyjVDqdbkOb1/g=",
                     "type" : "ByteString"
                  },
                  {
                     "value" : "1100000000000000",
                     "type" : "Integer"
                  }
               ],
               "type" : "Array"
            },
            "eventname" : "transfer"
         }
      ],
      "txid" : "0x2f4eff5532808080ae4c6a1a844c32a53015ca3cad5c308c5bcf3d8261399d06",
      "gasconsumed" : "8233460",
      "trigger" : "Application"
   },
   "jsonrpc" : "2.0",
   "id" : 1
}
```
Here we have `1` at the `stack` field, which means that token was successfully minted.
Let's just ensure that by querying `balanceOf` one more time:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 4f84a3b4a32125056b5b2bf5b6c7addc3616985f balanceOf NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt
```
... with `qwerty` pass. The result:
``` 
Sent invocation transaction cb3a28b2670836d1034655c19b5fbea62593f5f834977e72b1599f28c0ea79d5
```
... with the following `getapplicationlog` JSON message:
```
{
   "jsonrpc" : "2.0",
   "result" : {
      "gasconsumed" : "5504020",
      "notifications" : [],
      "vmstate" : "HALT",
      "txid" : "0xcb3a28b2670836d1034655c19b5fbea62593f5f834977e72b1599f28c0ea79d5",
      "trigger" : "Application",
      "stack" : [
         {
            "type" : "Integer",
            "value" : "1100000000000000"
         }
      ]
   },
   "id" : 1
}
```
Now we can see integer value at the `stack` field, so `1100000000000000` is the nep5 token balance of our account.

Note, that token can be minted only once.

#### Step #5

After we are done with minting, it's possible to transfer token to someone else.
Let's transfer 5 tokens from our account to `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm` with `transfer` call:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 4f84a3b4a32125056b5b2bf5b6c7addc3616985f transfer NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm 500000000 -- f8d79b436e9d3a54a38cc877d182b71e381d7a5c:CalledByEntry
```
... with password `qwerty` and following result:
``` 
Sent invocation transaction 65864bd07566318141fd781ece65a72d3c185f3f7190ac7096f7b9b7b62583ee
```
Our favourite `getapplicationlog` RPC-call tells us:
```
{
   "id" : 1,
   "result" : {
      "trigger" : "Application",
      "txid" : "0x65864bd07566318141fd781ece65a72d3c185f3f7190ac7096f7b9b7b62583ee",
      "vmstate" : "HALT",
      "gasconsumed" : "8117770",
      "notifications" : [
         {
            "state" : {
               "type" : "Array",
               "value" : [
                  {
                     "type" : "ByteString",
                     "value" : "XHodOB63gtF3yIyjVDqdbkOb1/g="
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
            },
            "eventname" : "transfer",
            "contract" : "0x4f84a3b4a32125056b5b2bf5b6c7addc3616985f"
         }
      ],
      "stack" : [
         {
            "value" : "1",
            "type" : "Integer"
         }
      ]
   },
   "jsonrpc" : "2.0"
}
```
Note, that `stack` field contains `1`, which means that token was successfully transferred.
Let's now check the balance of `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm` account to ensure that the amount of token on that account = 5:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 4f84a3b4a32125056b5b2bf5b6c7addc3616985f balanceOf NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm
```
The `getapplicationlog` RPC-call for this transaction tells us the following:
```
{
   "result" : {
      "trigger" : "Application",
      "stack" : [
         {
            "type" : "Integer",
            "value" : "500000000"
         }
      ],
      "notifications" : [],
      "vmstate" : "HALT",
      "txid" : "0x246aab29aefe81baed1c3a40243b8ee9113689a305e90aa5a89efdaba880bf0e",
      "gasconsumed" : "5504020"
   },
   "id" : 1,
   "jsonrpc" : "2.0"
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
Sent deployment transaction 2796d6d641d3e7528f7cb7c30cc8b0099e1291e2af3a14ef4a8d89f3dca14c44 for contract 459f3383e8b087a67454da141eeba34521eaccfc
```   
You know, what it means :)

#### Step #3

Invoke the contract to register domain with name `my_first_domain`:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 459f3383e8b087a67454da141eeba34521eaccfc register my_first_domain NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt -- f8d79b436e9d3a54a38cc877d182b71e381d7a5c:CalledByEntry
```
... the strongest password in the world, guess: `qwerty`
```
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```
Result:
```
Sent invocation transaction f6995b104a03fdb3dc45b236c371523276b8c17036631e3df9fdc5575d608988
```
Also you can see the log message in the console, where you run neo-go node:
```
2020-07-06T16:50:37.713+0300	INFO	runtime log	{"script": "fcccea2145a3eb1e14da5474a687b0e883339f45", "logs": "\"RegisterDomain: my_first_domain\""}
```
Well, that's ok. Let's check now, whether our domain was registered with `getapplicationlog` RPC-call:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["f6995b104a03fdb3dc45b236c371523276b8c17036631e3df9fdc5575d608988"] }' localhost:20331 | json_pp
```
The result is:
```
{
   "id" : 1,
   "jsonrpc" : "2.0",
   "result" : {
      "vmstate" : "HALT",
      "txid" : "0xf6995b104a03fdb3dc45b236c371523276b8c17036631e3df9fdc5575d608988",
      "gasconsumed" : "6223400",
      "stack" : [
         {
            "type" : "Integer",
            "value" : "1"
         }
      ],
      "trigger" : "Application",
      "notifications" : [
         {
            "contract" : "0x459f3383e8b087a67454da141eeba34521eaccfc",
            "eventname" : "registered",
            "state" : {
               "type" : "Array",
               "value" : [
                  {
                     "value" : "XHodOB63gtF3yIyjVDqdbkOb1/g=",
                     "type" : "ByteString"
                  },
                  {
                     "value" : "bXlfZmlyc3RfZG9tYWlu",
                     "type" : "ByteString"
                  }
               ]
            }
         }
      ]
   }
}
```
Especially, we're interested in two fields of the json:

First one is `notifications` field, which contains one notification with `registered` name:
- `bXlfZmlyc3RfZG9tYWlu` byte string in base64 representation, which can be decoded to `my_first_domain` - our domain's name
- `E9RIbx59YF1EGCw9f0E19Ms6o/A=` byte array, which can be decoded to the account address `NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt`.

The second field is `stack` with `1` value, which was returned by the smart contract.

All of these values let us be sure that our domain was successfully registered.  

#### Step #4

Invoke the contract to query the address information our `my_first_domain` domain:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 459f3383e8b087a67454da141eeba34521eaccfc query my_first_domain
```
... the pass `qwerty`:
```
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```
Result:
```
Sent invocation transaction 7890938840734c26a66ed642ab50ce9d5f1c075aa090a901c87fccb638d92012
```
and log-message:
```
2020-07-06T17:02:10.782+0300	INFO	runtime log	{"script": "fcccea2145a3eb1e14da5474a687b0e883339f45", "logs": "\"QueryDomain: my_first_domain\""}
```
Let's check this transaction with `getapplicationlog` RPC call:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["7890938840734c26a66ed642ab50ce9d5f1c075aa090a901c87fccb638d92012"] }' localhost:20331 | json_pp
```
... which gives us the following result:
```
{
   "jsonrpc" : "2.0",
   "id" : 1,
   "result" : {
      "trigger" : "Application",
      "vmstate" : "HALT",
      "txid" : "0x7890938840734c26a66ed642ab50ce9d5f1c075aa090a901c87fccb638d92012",
      "stack" : [
         {
            "type" : "ByteString",
            "value" : "XHodOB63gtF3yIyjVDqdbkOb1/g="
         }
      ],
      "notifications" : [],
      "gasconsumed" : "4185570"
   }
}
```

with base64 interpretation of our account address `NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt` on the stack, which means that domain `my_first_domain` was registered by the owner with received account address.

#### Step #5

Invoke the contract to transfer domain to the other account (e.g. account with `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm` address):
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 459f3383e8b087a67454da141eeba34521eaccfc transfer my_first_domain NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm -- f8d79b436e9d3a54a38cc877d182b71e381d7a5c
```
... the password: `qwerty`
```
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```
Result:
```
Sent invocation transaction 7df4b7d4f318e1e2c3007ec919bbcbff620685df71add4834e0f1a35f57ef893
```
and log-message:
```
2020-07-06T17:10:08.521+0300	INFO	runtime log	{"script": "fcccea2145a3eb1e14da5474a687b0e883339f45", "logs": "\"TransferDomain: my_first_domain\""}
```
Perfect. And `getapplicationlog` RPC-call...
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["7df4b7d4f318e1e2c3007ec919bbcbff620685df71add4834e0f1a35f57ef893"] }' localhost:20331 | json_pp
```
... tells us:
```
{
   "id" : 1,
   "jsonrpc" : "2.0",
   "result" : {
      "vmstate" : "HALT",
      "trigger" : "Application",
      "txid" : "0x7df4b7d4f318e1e2c3007ec919bbcbff620685df71add4834e0f1a35f57ef893",
      "notifications" : [
         {
            "contract" : "0x459f3383e8b087a67454da141eeba34521eaccfc",
            "state" : {
               "value" : [
                  {
                     "type" : "ByteString",
                     "value" : "XHodOB63gtF3yIyjVDqdbkOb1/g="
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
            "contract" : "0x459f3383e8b087a67454da141eeba34521eaccfc",
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
            }
         }
      ],
      "stack" : [
         {
            "value" : "1",
            "type" : "Integer"
         }
      ],
      "gasconsumed" : "5410900"
   }
}
```
The `notifications` field contains two events:
- First one with name `deleted` and additional information (domain `my_first_domain` was deleted from account `NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt`),
- Second one with name `cmVnaXN0ZXJlZA==`  and additional information  (domain `my_first_domain` was registered with account `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm`).
The `stack` field contains `1` value, which means that our domain was successfully transferred.

#### Step #6

The last call is `delete`, so you can try to create the other domain, e.g. `my_second_domain` and then remove it from storage with:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 459f3383e8b087a67454da141eeba34521eaccfc delete my_second_domain -- f0a33acbf435417f3d2c18445d607d1e6f48d413
```

Thank you!

### Useful links

* [Our basic tutorial on Medium](https://medium.com/@neospcc/%D1%81%D0%BC%D0%B0%D1%80%D1%82-%D0%BA%D0%BE%D0%BD%D1%82%D1%80%D0%B0%D0%BA%D1%82-%D0%B4%D0%BB%D1%8F-neo-769139352b65)
* [Using NEO Blockchain Toolkit](https://medium.com/@neospcc/neogo-adds-support-for-neo-blockchain-toolkit-673ea914f661)
* [NEO documentation](https://docs.neo.org/)
* [NEO github](https://github.com/neo-project/neo/)
* [NEO-GO github](https://github.com/nspcc-dev/neo-go)
