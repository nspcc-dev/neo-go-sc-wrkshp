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
More information how to use cli you can find [here](https://github.com/nspcc-dev/neo-go/blob/master-2.x/docs/cli.md)

#### Network
Network layer is one of the most important parts of the node. In our case we have P2P protocol which allows nodes to communicate with each other and RPC -- which is used for getting some information from node like balance, accounts, current state, etc.
Here is the document where you can find supported [RPC calls](https://github.com/nspcc-dev/neo-go/blob/master-2.x/docs/rpc.md).

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

## Workshop. Part 1
Now it’s time to run your private network. Connect neo-go node to it, write smart contract, deploy and invoke it. 
Let’s go!
#### Requirements
For this workshop you will need Debian 10, Docker, docker-compose, go to be installed:
- [docker](https://docs.docker.com/install/linux/docker-ce/debian/)
- [go](https://golang.org/dl/)

#### Step 1
Download neo-go and build it (master-2.x branch):
```
git clone https://github.com/nspcc-dev/neo-go.git
$ cd neo-go
$ git checkout -b master-2.x origin/master-2.x
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
Create basic "Hello World" smart contract(or use the one presented in this repo):
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

#### Step 4 
Run "Hello World" smart contract:
```
$ ./bin/neo-go contract compile -i '/1-print.go'
```
Where
- `./bin/neo-go` runs neo-go
- `contract compile` command with arguments in [neo-go](https://github.com/nspcc-dev/neo-go/blob/master-2.x/cli/smartcontract/smart_contract.go#L43)
- `-i '/1-print.go'` path to smart contract

Result:

Compiled smart-contract: `1-pring.avm`

To dump all the opcodes, you can use:
```
$ ./bin/neo-go contract inspect -i '/1-print.go'
```

#### Step 5
Start neo-go node which will connect to previously started privatenet:
```
$ ./bin/neo-go node --privnet
```

Result:
```
INFO[0000] no storage version found! creating genesis block 
INFO[0000] Pprof service hasn't started since it's disabled 
INFO[0000] Prometheus service hasn't started since it's disabled 
INFO[0000] bad MinPeers configured, using the default value  MinPeers actual=5 MinPeers configured=0
INFO[0000] bad AttemptConnPeers configured, using the default value  AttemptConnPeers actual=20 AttemptConnPeers configured=0

    _   ____________        __________
   / | / / ____/ __ \      / ____/ __ \
  /  |/ / __/ / / / /_____/ / __/ / / /
 / /|  / /___/ /_/ /_____/ /_/ / /_/ /
/_/ |_/_____/\____/      \____/\____/

/NEO-GO:0.70.2-pre-11-g735b937/

INFO[0000] RPC server is not enabled                    
INFO[0000] node started                                  blockHeight=0 headerHeight=0
INFO[0000] new peer connected                            addr="127.0.0.1:20336"
...
```

#### Step 6
Deploy smart contract:
```
./bin/neo-go contract deploy -i 1-print.avm -c 1-print.yml -e \
http://localhost:20331 -w my_wallet.json -g 0.001
```

Where
- `contract deploy` is a command for deployment
- `-i '/1-print.avm'` path to smart contract
- `-c 1-print.yml` configuration input file
- `-e http://localhost:20331` node endpoint
- `-w my_wallet.json` wallet to use to get the key for transaction signing (you can use one from the workshop repo)
- `-g 0.001` amount of gas to pay for contract deployment

Enter password `qwerty` for account:
```
Enter account AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y password >
```

Result:
```
Sent deployment transaction ea93196802fe3517d2d028e4d4f244aa734a9b1988456740d96f2c3336140fda for contract 6d1eeca891ee93de2b7a77eb91c26f3b3c04d6cf
```

At this point your ‘Hello World’ contract is deployed and could be invoked. Let’s do it as a final step.

#### Step 7
Invoke contract.
```
$ ./bin/neo-go contract invoke -e http://localhost:20331 -w my_wallet.json -g 0.00001 6d1eeca891ee93de2b7a77eb91c26f3b3c04d6cf
```

Where
- `contract invoke` runs invoke with provided parameters
- `-e http://localhost:20331` defines RPC endpoint used for function call
- `-w my_wallet.json` is a wallet
- `-g 0.00001` defines amount of GAS to be used for invoke operation
- `6d1eeca891ee93de2b7a77eb91c26f3b3c04d6cf` contract hash got as an output from the previous command (deployment in step 6)

Enter password `qwerty` for account:
```
Enter account AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y password >
```

Result:
In the console where you were running step #5 you will get:
```
INFO[0227] script cfd6043c3b6fc291eb777a2bde93ee91a8ec1e6d logs: "Hello, world!"
```
Which means that this contract was executed.

This is it. There are only 5 steps to make deployment and they look easy, aren’t they?
Thank you!

## Workshop. Part 2
In this part we'll look at RPC calls and try to write, deploy and invoke smart contract with storage. 
Let’s go!

### RPC calls
Let's check what's going on under the hood. 
Each neo-go node provides an API interface for obtaining blockchain data from it.
The interface is provided via `JSON-RPC`, and the underlying protocol uses HTTP for communication.

Full `NEO JSON-RPC 2.0 API` described [here](https://docs.neo.org/docs/en-us/reference/rpc/latest-version/api.html).

RPC-server of started in step #5 neo-go node is available on `localhost:20331`, so let's try to perform several RPC calls.

#### GetRawTransaction
[GetRawTransaction](https://docs.neo.org/docs/en-us/reference/rpc/latest-version/api/getrawtransaction.html) returns 
the corresponding transaction information, based on the specified hash value.

Request information about our deployment transaction:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getrawtransaction", "params": ["ea93196802fe3517d2d028e4d4f244aa734a9b1988456740d96f2c3336140fda", 1] }' localhost:20331 | json_pp
```

Where
- `"jsonrpc": "2.0"` is protocol version
- `"id": 1` is id of current request
- `"method": "getrawtransaction"` is requested method
- `"params": ["ea93196802fe3517d2d028e4d4f244aa734a9b1988456740d96f2c3336140fda", 1]` is array of parameters, where
   - `ea93196802fe3517d2d028e4d4f244aa734a9b1988456740d96f2c3336140fda` is deployment transaction hash
   - `1` is `verbose` parameter for detailed JSON string output
- `json_pp` just makes the JSON output prettier

Result:
```
{
   "result" : {
      "blocktime" : 1584027315,
      "net_fee" : "100",
      "type" : "InvocationTransaction",
      "version" : 1,
      "vin" : [
         {
            "vout" : 0,
            "txid" : "0x9aa010ea9618b34dd2dc42d35015280701c35a292cfb702ad32e86d40e9239cb"
         }
      ],
      "txid" : "0xea93196802fe3517d2d028e4d4f244aa734a9b1988456740d96f2c3336140fda",
      "sys_fee" : "100",
      "size" : 341,
      "confirmations" : 70,
      "blockhash" : "0xa0ed247b4161426532fd20c7c21da633a6994493216c71f3a9940d2256d3cee4",
      "vout" : [
         {
            "n" : 0,
            "value" : "19952",
            "asset" : "0x602c79718b16e442de58778e148d0b1084e3b2dffd5de6b7b16cee7969282de7",
            "address" : "AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y"
         }
      ],
      "scripts" : [
         {
            "invocation" : "4001a7e30b4e3b4ec2f50bb7af86bf2b8ef03ae55fbf302e4f5db36df12fe1e417a01d92949376131c0be5c472cbc47c61669fc8c7abc691a071107c1271b2a460",
            "verification" : "21031a6c6fbbdf02ca351745fa86b9ba5a9452d785ac4f7fc2b7548ca2a46c4fcf4aac"
         }
      ],
      "attributes" : []
   },
   "jsonrpc" : "2.0",
   "id" : 1
}
```

#### GetApplicationLog
[GetApplicationLog](https://docs.neo.org/docs/en-us/reference/rpc/latest-version/api/getapplicationlog.html) returns the contract log based on the specified transaction id.

Request application log for invocation transaction from step #7:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["42879c901de7d180f2ff609aded2da2b9c63d5ed7d1b22f18e246c16fc5d0c57"] }' localhost:20331 | json_pp
```

With a single parameter:
- `42879c901de7d180f2ff609aded2da2b9c63d5ed7d1b22f18e246c16fc5d0c57` - invocation transaction hash from step #7

Result:
```
{
   "id" : 1,
   "result" : {
      "txid" : "0x42879c901de7d180f2ff609aded2da2b9c63d5ed7d1b22f18e246c16fc5d0c57",
      "executions" : [
         {
            "gas_consumed" : "0.017",
            "contract" : "0x9dbb827c329d765240569058a5bdd8176aab4cb6",
            "stack" : [
               {
                  "type" : "Array",
                  "value" : []
               },
               {
                  "type" : "ByteArray",
                  "value" : ""
               }
            ],
            "trigger" : "Application",
            "notifications" : [],
            "vmstate" : "HALT"
         }
      ]
   },
   "jsonrpc" : "2.0"
}
```

#### Other Useful RPC calls
```
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getversion", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getblockcount", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getconnectioncount", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getaccountstate", "params": ["AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y"] }' localhost:20331
```

List of supported by neo-go node RPC commands you can find [here](https://github.com/nspcc-dev/neo-go/blob/master-2.x/docs/rpc.md#supported-methods).

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
./bin/neo-go contract compile -i 2-storage.go
```

Result:

Compiled smart-contract: `2-storage.avm`

#### Step #2
Deploy compiled smart contract with [configuration](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/2-storage.yml):
```
./bin/neo-go contract deploy -i 2-storage.avm -c 2-storage.yaml -e http://localhost:20331 -w my_wallet.json -g 0.001
```
... enter the password `qwerty`:
```
Enter account AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y password >
```

Result:
```
Sent deployment transaction d79e5e39d87c2911b623d3efe98842cfde41eddc34d85cee394783b7320813e8 for contract 85cf2075f3e297d489ff3c4c1745ca80d44e2a68
```   

Which means that our contract was deployed and now we can invoke it.

#### Step #3
Let's invoke our contract. As far as we have never invoked this contract, there's no value in the storage, so the contract should create a new one (which is `1`) and put it into storage.
Let's check:
```
./bin/neo-go contract invoke -e http://localhost:20331 -w my_wallet.json -g 0.00001 85cf2075f3e297d489ff3c4c1745ca80d44e2a68
```
... enter the password `qwerty`:
```
Enter account AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y password >
```
Result:
```
Sent invocation transaction 6f27f523da9c71297f4a81a254274b0c2a78f893b81c500429be6230254be0bf
```
To check the counter value, call `getapplicaionlog` RPC-call for the invocation transaction:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["6f27f523da9c71297f4a81a254274b0c2a78f893b81c500429be6230254be0bf"] }' localhost:20331 | json_pp
```
The JSON result is:
```
{
   "id" : 1,
   "result" : {
      "executions" : [
         {
            "vmstate" : "HALT",
            "gas_consumed" : "1.159",
            "trigger" : "Application",
            "contract" : "0x343b284abf1e6a441a1361c5de76bdcb15b8e332",
            "notifications" : [
               {
                  "contract" : "0x85cf2075f3e297d489ff3c4c1745ca80d44e2a68",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "type" : "ByteArray",
                           "value" : "56616c756520726561642066726f6d2073746f72616765"
                        }
                     ]
                  }
               },
               {
                  "contract" : "0x85cf2075f3e297d489ff3c4c1745ca80d44e2a68",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "type" : "ByteArray",
                           "value" : "53746f72616765206b6579206e6f7420796574207365742e2053657474696e6720746f2031"
                        }
                     ]
                  }
               },
               {
                  "state" : {
                     "value" : [
                        {
                           "value" : "4e65772076616c7565207772697474656e20696e746f2073746f72616765",
                           "type" : "ByteArray"
                        }
                     ],
                     "type" : "Array"
                  },
                  "contract" : "0x85cf2075f3e297d489ff3c4c1745ca80d44e2a68"
               }
            ],
            "stack" : [
               {
                  "value" : "1",
                  "type" : "Integer"
               }
            ]
         }
      ],
      "txid" : "0x6f27f523da9c71297f4a81a254274b0c2a78f893b81c500429be6230254be0bf"
   },
   "jsonrpc" : "2.0"
}
```
Pay attention to `notification` field. It contains messages, which where passed to `runtime.Notify` method.
This one contains hexadecimal byte arrays which can be decoded into 3 messages:
  - `Value read from storage` which was called after we've got the counter value from storage
  - `Storage key not yet set. Setting to 1` which was called when we realised that counter value is 0
  - `New value written into storage` which was called after the counter value was put in the storage.
  
The final part is `stack` field. This field contains all returned by the contract values, so here you can see integer value `1`,
which is the counter value denoted to the number of contract invocations.

#### Step #4
To ensure that all works as expected, let's invoke the contract one more time and check, whether the counter will be incremented: 
```
./bin/neo-go contract invoke -e http://localhost:20331 -w my_wallet.json -g 0.00001 85cf2075f3e297d489ff3c4c1745ca80d44e2a68
```
... enter the password `qwerty`:
```
Enter account AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y password >
```
Result:
```
Sent invocation transaction dcf53cdb69c816f0f9ab27ba509fb656b6eddd2dad5a658e9f534e9dba38462b
```
To check the counter value, call `getapplicaionlog` RPC-call for the invocation transaction:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["dcf53cdb69c816f0f9ab27ba509fb656b6eddd2dad5a658e9f534e9dba38462b"] }' localhost:20331 | json_pp
```
The JSON result is:
```
{
   "jsonrpc" : "2.0",
   "id" : 1,
   "result" : {
      "executions" : [
         {
            "vmstate" : "HALT",
            "contract" : "0x343b284abf1e6a441a1361c5de76bdcb15b8e332",
            "trigger" : "Application",
            "notifications" : [
               {
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "value" : "56616c756520726561642066726f6d2073746f72616765",
                           "type" : "ByteArray"
                        }
                     ]
                  },
                  "contract" : "0x85cf2075f3e297d489ff3c4c1745ca80d44e2a68"
               },
               {
                  "contract" : "0x85cf2075f3e297d489ff3c4c1745ca80d44e2a68",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "type" : "ByteArray",
                           "value" : "53746f72616765206b657920616c7265616479207365742e20496e6372656d656e74696e672062792031"
                        }
                     ]
                  }
               },
               {
                  "contract" : "0x85cf2075f3e297d489ff3c4c1745ca80d44e2a68",
                  "state" : {
                     "value" : [
                        {
                           "value" : "4e65772076616c7565207772697474656e20696e746f2073746f72616765",
                           "type" : "ByteArray"
                        }
                     ],
                     "type" : "Array"
                  }
               }
            ],
            "stack" : [
               {
                  "type" : "Integer",
                  "value" : "2"
               }
            ],
            "gas_consumed" : "1.161"
         }
      ],
      "txid" : "0xdcf53cdb69c816f0f9ab27ba509fb656b6eddd2dad5a658e9f534e9dba38462b"
   }
}
```

The `stack` field contains now `2` integer value, so the counter was incremented as we expected.

## Workshop. Part 3
In this part we'll know about NEP5 token standard and try to write, deploy and invoke more complicated smart contract. 
Let’s go!

### NEP5
[NEP5](https://docs.neo.org/docs/en-us/sc/write/nep5.html) is a token standard for the Neo blockchain that provides systems with a generalized interaction mechanism for tokenized smart contracts.
The example with implementation of all required by the standard methods you can find in [nep5.go](https://github.com/nspcc-dev/neo-go/blob/master-2.x/examples/token/nep5/nep5.go)
 
Let's take a view on the example of smart contract with NEP5: [token.go](https://github.com/nspcc-dev/neo-go/blob/master-2.x/examples/token/token.go)
 
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
Compile smart contract [token.go](https://github.com/nspcc-dev/neo-go/blob/master-2.x/examples/token/token.go):
```
./bin/neo-go contract compile -i examples/token/token.go
```

Note: as a deploy configuration you can use [configuration](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/1-print.yml) from Part 1 of the workshop with following changes:
as far as our contract uses storage, the flag `hasstorage` should be set to `true`:
```
hasstorage: true
```

Deploy smart contract with modified configuration:
```
./bin/neo-go contract deploy -i examples/token/token.avm -c nep5.yml -e http://localhost:20331 -w my_wallet.json -g 0.001
```
... enter the password `qwerty`:
```
Enter account AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y password >
```

Result:
```
Sent deployment transaction 6e46e2f2c8e799bf0fce679fa3c8acbc046f595252ae58b39abea839da01067d for contract f84d6a337fbc3d3a201d41da99e86b479e7a2554
```   

Which means that our contract was deployed and now we can invoke it.

#### Step #2
Let's invoke the contract to perform different operations.

To start with, query `name` of the created nep5 token:

```
./bin/neo-go contract invokefunction -e http://localhost:20331 -w my_wallet.json -g 0.00001 f84d6a337fbc3d3a201d41da99e86b479e7a2554 name
```                                                                   
Where
- `f84d6a337fbc3d3a201d41da99e86b479e7a2554` is our contract hash from step #1
- `name` is operation string which was described earlier and returns token name

... and don't forget the password of your account `qwerty`.

Result:
```
Sent invocation transaction cfdf4c2883d71a4375ab94fbb302a386828e7934541aa222fc38b8bc67e6a2b4
```                                                                                         
Now, let's take a detailed look at this invocation transaction with `getapplicationlog` RPC call:

```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["cfdf4c2883d71a4375ab94fbb302a386828e7934541aa222fc38b8bc67e6a2b4"] }' localhost:20331 | json_pp
```               

Result:
```
"result" : {
      "txid" : "0xcfdf4c2883d71a4375ab94fbb302a386828e7934541aa222fc38b8bc67e6a2b4",
      "executions" : [
         {
            "contract" : "0xc1c04480d56e04689c0ac4db973516f9a44f277d",
            "gas_consumed" : "0.059",
            "notifications" : [],
            "stack" : [
               {
                  "type" : "ByteArray",
                  "value" : "417765736f6d65204e454f20546f6b656e"
               }
            ],
            "trigger" : "Application",
            "vmstate" : "HALT"
         }
      ]
   },
   "id" : 1,
   "jsonrpc" : "2.0"
}
```

At least, you can see that `stack` field of JSON result is not empty: it contains byte array with the name of our token.

Following commands able you to get some additional information about token:
```
./bin/neo-go contract invokefunction -e http://localhost:20331 -w my_wallet.json -g 0.00001 f84d6a337fbc3d3a201d41da99e86b479e7a2554 symbol
./bin/neo-go contract invokefunction -e http://localhost:20331 -w my_wallet.json -g 0.00001 f84d6a337fbc3d3a201d41da99e86b479e7a2554 decimals
./bin/neo-go contract invokefunction -e http://localhost:20331 -w my_wallet.json -g 0.00001 f84d6a337fbc3d3a201d41da99e86b479e7a2554 totalSupply
```

#### Step #3

Now it's time for more interesting things. First of all, let's check the balance of nep5 token on our account by using `balanceOf`:
```
./bin/neo-go contract invokefunction -e http://localhost:20331 -w my_wallet.json -g 0.00001 f84d6a337fbc3d3a201d41da99e86b479e7a2554 balanceOf AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y
```                             
... with `qwerty` password. The result is:
```
Sent invocation transaction 30de65dec667d68ca1b385c590c9c5fccf82e2d4831540e6bb5875afa57c5cbe
```
And take a closer look at the transaction's details with `getapplicationlog` RPC-call:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["30de65dec667d68ca1b385c590c9c5fccf82e2d4831540e6bb5875afa57c5cbe"] }' localhost:20331 | json_pp
```
Result:
```
{
   "jsonrpc" : "2.0",
   "id" : 1,
   "result" : {
      "executions" : [
         {
            "trigger" : "Application",
            "vmstate" : "HALT",
            "stack" : [
               {
                  "value" : "",
                  "type" : "ByteArray"
               }
            ],
            "notifications" : [],
            "gas_consumed" : "0.209",
            "contract" : "0x762ca50a574b7140961283e9d45fc67d1482b0ba"
         }
      ],
      "txid" : "0x30de65dec667d68ca1b385c590c9c5fccf82e2d4831540e6bb5875afa57c5cbe"
   }
}
``` 
As far as `stack` field contains an empty byte array, we have no token on the balance. But don't worry about that. Just follow the next step.

#### Step #4

Before we are able to start using our token (e.g. transfer it to someone else), we have to *mint* it.
In other words, we should transfer all available amount of token (total supply) to someone's account.
There's a special function for this purpose in our contract - `mint` function, so let's mint token to our address:
```
./bin/neo-go contract invokefunction -e http://localhost:20331 -w my_wallet.json -g 0.00001 f84d6a337fbc3d3a201d41da99e86b479e7a2554 mint AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y
```
... with `qwerty` pass. The result:
``` 
Sent invocation transaction a571adebdbdabfd087f34867da94524649003c6b851ed0cc5da7a30ff843bc1e
```
`getapplicationlog` RPC-call for this transaction tells us the following:
```
{
   "result" : {
      "executions" : [
         {
            "stack" : [
               {
                  "value" : "1",
                  "type" : "Integer"
               }
            ],
            "trigger" : "Application",
            "vmstate" : "HALT",
            "contract" : "0x4176c0e2f8b5b23910dac91a77cd97784e618c73",
            "gas_consumed" : "2.489",
            "notifications" : [
               {
                  "contract" : "0xf6ac2777b1cbd227bed1fa5735bd06befdee6d34",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "value" : "7472616e73666572",
                           "type" : "ByteArray"
                        },
                        {
                           "type" : "ByteArray",
                           "value" : ""
                        },
                        {
                           "value" : "23ba2703c53263e8d6e522dc32203339dcd8eee9",
                           "type" : "ByteArray"
                        },
                        {
                           "type" : "ByteArray",
                           "value" : "00c040b571e803"
                        }
                     ]
                  }
               }
            ]
         }
      ],
      "txid" : "0xa571adebdbdabfd087f34867da94524649003c6b851ed0cc5da7a30ff843bc1e"
   },
   "jsonrpc" : "2.0",
   "id" : 1
}
```
Here we have `1` at the `stack` field, which means that token was successfully minted.
Let's just ensure that by querying `balanceOf` one more time:
```
./bin/neo-go contract invokefunction -e http://localhost:20331 -w my_wallet.json -g 0.00001 f84d6a337fbc3d3a201d41da99e86b479e7a2554 balanceOf AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y
```
... with `qwerty` pass. The result:
``` 
Sent invocation transaction c56f469cd9d47c6a4195a742752621c4898447aa4cfd7f550046bdd10d297c12
```
... with the following `getapplicationlog` JSON message:
```
{
   "result" : {
      "executions" : [
         {
            "notifications" : [],
            "stack" : [
               {
                  "value" : "00c040b571e803",
                  "type" : "ByteArray"
               }
            ],
            "contract" : "0x762ca50a574b7140961283e9d45fc67d1482b0ba",
            "vmstate" : "HALT",
            "trigger" : "Application",
            "gas_consumed" : "0.209"
         }
      ],
      "txid" : "0xc56f469cd9d47c6a4195a742752621c4898447aa4cfd7f550046bdd10d297c12"
   },
   "id" : 1,
   "jsonrpc" : "2.0"
}
```
Now we can see non-empty byte array at the `stack` field, so `00c040b571e803` is a hexadecimal representation of the nep5 token balance of our account.

Note, that token can be minted only ones.

#### Step #5

After we are done with minting, it's possible to transfer token to someone else.
Let's transfer 5 tokens from our account to `AKkkumHbBipZ46UMZJoFynJMXzSRnBvKcs` with `transfer` call:
```
./bin/neo-go contract invokefunction -e http://localhost:20331 -w my_wallet.json -g 0.00001 f84d6a337fbc3d3a201d41da99e86b479e7a2554 transfer AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y AKkkumHbBipZ46UMZJoFynJMXzSRnBvKcs 5
```
... with password `qwerty` and following result:
``` 
Sent invocation transaction 11a272f4a0d7912f7219979bab7d094df3b404b89e903337ee72a90249cc448d
```
Our favourite `getapplicationlog` RPC-call tells us:
```
{
   "id" : 1,
   "result" : {
      "executions" : [
         {
            "vmstate" : "HALT",
            "contract" : "0x102277f9ab76c0bc0452e890652c7e272ce9c94a",
            "gas_consumed" : "2.485",
            "notifications" : [
               {
                  "state" : {
                     "value" : [
                        {
                           "value" : "7472616e73666572",
                           "type" : "ByteArray"
                        },
                        {
                           "type" : "ByteArray",
                           "value" : "23ba2703c53263e8d6e522dc32203339dcd8eee9"
                        },
                        {
                           "type" : "ByteArray",
                           "value" : "2baa76ad534b886cb87c6b3720a34943d9000fa9"
                        },
                        {
                           "value" : "5",
                           "type" : "Integer"
                        }
                     ],
                     "type" : "Array"
                  },
                  "contract" : "0xf6ac2777b1cbd227bed1fa5735bd06befdee6d34"
               }
            ],
            "stack" : [
               {
                  "value" : "1",
                  "type" : "Integer"
               }
            ],
            "trigger" : "Application"
         }
      ],
      "txid" : "0x11a272f4a0d7912f7219979bab7d094df3b404b89e903337ee72a90249cc448d"
   },
   "jsonrpc" : "2.0"
}
```
Note, that `stack` field contains `1`, which means that token was successfully transferred.
Let's now check the balance of `AKkkumHbBipZ46UMZJoFynJMXzSRnBvKcs` account to ensure that the amount of token on that account = 5:
```
./bin/neo-go contract invokefunction -e http://localhost:20331 -w my_wallet.json -g 0.00001 f84d6a337fbc3d3a201d41da99e86b479e7a2554 balanceOf AKkkumHbBipZ46UMZJoFynJMXzSRnBvKcs
```
The `getapplicationlog` RPC-call for this transaction tells us the following:
```
{
   "result" : {
      "txid" : "0x172c5074646ce043095e612d31e5c5c3d00c7c8a4e8c01873cc732692d5152f5",
      "executions" : [
         {
            "stack" : [
               {
                  "value" : "05",
                  "type" : "ByteArray"
               }
            ],
            "gas_consumed" : "0.209",
            "notifications" : [],
            "vmstate" : "HALT",
            "trigger" : "Application",
            "contract" : "0xbaaafb6be440d1de3d298ba556ad23aa0209ef2f"
         }
      ]
   },
   "id" : 1,
   "jsonrpc" : "2.0"
}
```
Here we are! There are exactly 5 tokens at the `stack` field. You can also ensure that these 5 tokens were debited from `AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y` account by using `balanceOf` method.

## Workshop. Part 4
In this part we'll summarise our knowledge about smart contracts by investigating [4-domain](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/4-domain.go) smart contract. This contract 
contains code for domain registration, transferring, deletion and getting information about registered domains.

Let’s go!

#### Step #1
Let's take a glance at our [contract](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/4-domain.go) and inspect it. The contract takes an action string as the first parameter, which is one of the following:
- `register` checks, whether domain with the specified name already exists. If not, it also adds the pair `[domain_name, owner]` to the storage. It requires additional arguments:
   - `domain_name` which is the new domain name.
   - `owner` - the 34-digit account address from our [wallet](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/my_wallet.json), which will be used for contract invocation.
- `query` returns the specified domain owner address (or false, if no such domain was registered). It requires the following argument:
   - `domain_name` which is requested domain name.
- `transfer` transfers domain with the specified name to the other address (of course, in case if you're the actual owner of the domain requested). It requires additional arguments:
   - `domain_name` which is the name of domain you'd like to transfer.
   - `to_address` - the account address you'd like to transfer the specified domain to.
- `delete` deletes specified domain from the storage. The arguments:
   - `domain_name` which is the name of the domain you'd like to delete.
 
 In the next steps we'll compile and deploy smart contract. 
 After that we'll try to register new domain, transfer it to another account and query information about it.

#### Step #2

Compile smart contract [4-domain.go](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/4-domain.go)
```
./bin/neo-go contract compile -i 4-domain.go
```

... and deploy it with [configuration](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/4-domain.yml):
```
./bin/neo-go contract deploy -i 4-domain.avm -c 4-domain.yml -e http://localhost:20331 -w my_wallet.json -g 0.001
```
Just a note: our contract uses storage and, as the previous one, needs the flag `hasstorage` to be set to `true` value.
That can be done in [configuration](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/4-domain.yml) file.

... enter the password `qwerty`:
```
Enter account AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y password >
```

Result:
```
Sent deployment transaction e2647e15a855dd174984fde05a3a193cf8c45deb5c741d3045cb0b35e9e462ce for contract 96f996eb681f65d380c596949d33d7f29897ad27
```   
You know, what it means :)

#### Step #3

Invoke the contract to register domain with name `my_first_domain`:
```
./bin/neo-go contract invokefunction -e http://localhost:20331 -w my_wallet.json -g 0.00001 96f996eb681f65d380c596949d33d7f29897ad27 register my_first_domain AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y
```
... the strongest password in the world, guess: `qwerty`
```
Enter account AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y password >
```
Result:
```
Sent invocation transaction 2d91a11f47217b7d7346d712456db79bdcaf982eea76ce105b41de859f7fc2e1
```
Also you can see the log message in the console, where you run neo-go node:
```
2020-04-06T11:57:26.889+0300	INFO	runtime log	{"script": "27ad9798f2d7339d9496c580d3651f68eb96f996", "logs": "\"RegisterDomain: my_first_domain\""}
```
Well, that's ok. Let's check now, whether our domain was registered with `getapplicationlog` RPC-call:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["2d91a11f47217b7d7346d712456db79bdcaf982eea76ce105b41de859f7fc2e1"] }' localhost:20331 | json_pp
```
The result is:
```
{
   "jsonrpc" : "2.0",
   "result" : {
      "txid" : "0x2d91a11f47217b7d7346d712456db79bdcaf982eea76ce105b41de859f7fc2e1",
      "executions" : [
         {
            "vmstate" : "HALT",
            "notifications" : [
               {
                  "contract" : "0x96f996eb681f65d380c596949d33d7f29897ad27",
                  "state" : {
                     "value" : [
                        {
                           "type" : "ByteArray",
                           "value" : "72656769737465726564"
                        },
                        {
                           "value" : "23ba2703c53263e8d6e522dc32203339dcd8eee9",
                           "type" : "ByteArray"
                        },
                        {
                           "value" : "6d795f66697273745f646f6d61696e",
                           "type" : "ByteArray"
                        }
                     ],
                     "type" : "Array"
                  }
               }
            ],
            "gas_consumed" : "1.397",
            "stack" : [
               {
                  "value" : "1",
                  "type" : "Integer"
               }
            ],
            "contract" : "0x965ec3b2648fab249ef9b732cb6c26d781b34d62",
            "trigger" : "Application"
         }
      ]
   },
   "id" : 1
}
```
Especially, we're interested in two fields of the json:

First one is `notifications` field, which contains 3 values:
- `72656769737465726564` byte array with hexadecimal value. A bit of decoding magic with the following script:
```
package main

import (
	"encoding/hex"
	"fmt"
)

func main() {
	bytes, _ := hex.DecodeString("72656769737465726564")
	fmt.Println(string(bytes))
}
```
... let us see the actual notification message: `registered`.

- `6d795f66697273745f646f6d61696e` byte array, which can be decoded to `my_first_domain` - our domain's name
- `23ba2703c53263e8d6e522dc32203339dcd8eee9` byte array, which can be decoded to the account address `AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y` by following script:
```
package main

import (
	"encoding/hex"
	"fmt"

	"github.com/nspcc-dev/neo-go/pkg/encoding/address"
	"github.com/nspcc-dev/neo-go/pkg/util"
)

func main() {
	addressBytes, _ := hex.DecodeString("23ba2703c53263e8d6e522dc32203339dcd8eee9")
	addressUint160, _ := util.Uint160DecodeBytesBE(addressBytes)
	fmt.Println(address.Uint160ToString(addressUint160))
}
```
The second field is `stack` with `1` value, which was returned by the smart contract.

All of these values let us be sure that our domain was successfully registered.  

#### Step #4

Invoke the contract to query the address information our `my_first_domain` domain:
```
./bin/neo-go contract invokefunction -e http://localhost:20331 -w my_wallet.json -g 0.00001 96f996eb681f65d380c596949d33d7f29897ad27 query my_first_domain
```
... the pass `qwerty`:
```
Enter account AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y password >
```
Result:
```
Sent invocation transaction 35c72ab78aa261242371426d2c33433312c291ef886be4f78e0ccbbe34d17349
```
and log-message:
```
2020-04-06T12:01:27.212+0300	INFO	runtime log	{"script": "27ad9798f2d7339d9496c580d3651f68eb96f996", "logs": "\"QueryDomain: my_first_domain\""}
```
Let's check this transaction with `getapplicationlog` RPC call:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["35c72ab78aa261242371426d2c33433312c291ef886be4f78e0ccbbe34d17349"] }' localhost:20331 | json_pp
```
... which gives us the following result:
```
{
   "id" : 1,
   "result" : {
      "executions" : [
         {
            "vmstate" : "HALT",
            "stack" : [
               {
                  "value" : "23ba2703c53263e8d6e522dc32203339dcd8eee9",
                  "type" : "ByteArray"
               }
            ],
            "gas_consumed" : "0.17",
            "trigger" : "Application",
            "notifications" : [],
            "contract" : "0x383b542b266a0953a30129e509fca8070dcbb668"
         }
      ],
      "txid" : "0x35c72ab78aa261242371426d2c33433312c291ef886be4f78e0ccbbe34d17349"
   },
   "jsonrpc" : "2.0"
}

```

with hexadecimal interpretation of our account address `AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y` on the stack, which means that domain `my_first_domain` was registered by the owner with received account address.

#### Step #5

Invoke the contract to transfer domain to the other account (e.g. account with `AKkkumHbBipZ46UMZJoFynJMXzSRnBvKcs` address):
```
./bin/neo-go contract invokefunction -e http://localhost:20331 -w my_wallet.json -g 0.00001 96f996eb681f65d380c596949d33d7f29897ad27 transfer my_first_domain AKkkumHbBipZ46UMZJoFynJMXzSRnBvKcs
```
... the password: `qwerty`
```
Enter account AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y password >
```
Result:
```
Sent invocation transaction da3653ab7eb23eab511b05ea64c2cdc3872b0c576d86b924a7738986cfd63462
```
and log-message:
```
2020-04-06T12:04:42.483+0300	INFO	runtime log	{"script": "27ad9798f2d7339d9496c580d3651f68eb96f996", "logs": "\"TransferDomain: my_first_domain\""}
```
Perfect. And `getapplicationlog` RPC-call...
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["da3653ab7eb23eab511b05ea64c2cdc3872b0c576d86b924a7738986cfd63462"] }' localhost:20331 | json_pp
```
... tells us:
```
{
   "result" : {
      "executions" : [
         {
            "trigger" : "Application",
            "contract" : "0x5e1e5f65562dfe38d214a8c4be36418c40614cb8",
            "stack" : [
               {
                  "type" : "Integer",
                  "value" : "1"
               }
            ],
            "vmstate" : "HALT",
            "notifications" : [
               {
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "value" : "64656c65746564",
                           "type" : "ByteArray"
                        },
                        {
                           "type" : "ByteArray",
                           "value" : "23ba2703c53263e8d6e522dc32203339dcd8eee9"
                        },
                        {
                           "type" : "ByteArray",
                           "value" : "6d795f66697273745f646f6d61696e"
                        }
                     ]
                  },
                  "contract" : "0x96f996eb681f65d380c596949d33d7f29897ad27"
               },
               {
                  "state" : {
                     "value" : [
                        {
                           "type" : "ByteArray",
                           "value" : "72656769737465726564"
                        },
                        {
                           "type" : "ByteArray",
                           "value" : "2baa76ad534b886cb87c6b3720a34943d9000fa9"
                        },
                        {
                           "type" : "ByteArray",
                           "value" : "6d795f66697273745f646f6d61696e"
                        }
                     ],
                     "type" : "Array"
                  },
                  "contract" : "0x96f996eb681f65d380c596949d33d7f29897ad27"
               }
            ],
            "gas_consumed" : "1.408"
         }
      ],
      "txid" : "0xda3653ab7eb23eab511b05ea64c2cdc3872b0c576d86b924a7738986cfd63462"
   },
   "id" : 1,
   "jsonrpc" : "2.0"
}
```
Quite a detailed one. The `notifications` field contains two arrays:
- First one with `64656c65746564` byte array, which is `deleted` with additional information (domain `my_first_domain` was deleted from account `AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y`),
- Second one with `72656769737465726564` byte array, which is `registered` with additional information  (domain `my_first_domain` was registered with account `AKkkumHbBipZ46UMZJoFynJMXzSRnBvKcs`).
The `stack` field contains `1` value, which means that our domain was successfully transferred.

#### Step #6

The last call is `delete`, so you can try to create the other domain, e.g. `my_second_domain` and then remove it from storage with:
```
./bin/neo-go contract invokefunction -e http://localhost:20331 -w my_wallet.json -g 0.00001 96f996eb681f65d380c596949d33d7f29897ad27 delete my_second_domain
```

Thank you!

### Useful links

* [Our basic tutorial on Medium](https://medium.com/@neospcc/%D1%81%D0%BC%D0%B0%D1%80%D1%82-%D0%BA%D0%BE%D0%BD%D1%82%D1%80%D0%B0%D0%BA%D1%82-%D0%B4%D0%BB%D1%8F-neo-769139352b65)
* [NEO documentation](https://docs.neo.org/)
* [NEO github](https://github.com/neo-project/neo/)
* [NEO-GO github](https://github.com/nspcc-dev/neo-go)
