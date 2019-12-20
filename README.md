<p align="center">
<img src="./pic/neo_color_dark_gopher.png" width="300px" alt="logo">
</p>

[NEO](https://neo.org/) builds smart economy and we at [NEO SPCC](https://nspcc.ru/en/) help them with that big challenge. 
In our blog you might find the latest articles [how we run NEOFS public test net](https://medium.com/@neospcc/public-neofs-testnet-launch-18f6315c5ced) 
but it’s not the only thing we’re working on.

## NEO GO
As you know network is composed of nodes. These nodes as of now have several implementations:
- https://github.com/CityOfZion/neo-sharp
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
	"github.com/CityOfZion/neo-go/pkg/interop/runtime"
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
More information how to use cli you can find [here](https://github.com/nspcc-dev/neo-go/blob/master/docs/cli.md)

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

## Workshop
Now it’s time to run your private network. Connect neo-go node to it, write smart contract, deploy and invoke it. 
Let’s go!

#### Step 1
Download neo-go and build it
```
git clone https://github.com/nspcc-dev/neo-go.git
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
Create basic "Hello World" smart contract(or use the one presented in this repo):
```
package main

import (
	"github.com/CityOfZion/neo-go/pkg/interop/runtime"
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
- `contract compile` command with arguments in [neo-go](https://github.com/nspcc-dev/neo-go/blob/master/cli/smartcontract/smart_contract.go#L43)
- `-i '/1-print.go'` path to smart contract

Result:

Compiled smart-contract: `1-pring.avm`

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
./bin/neo-go contract deploy -i 1-print.avm -c 1-print.yml -e 
http://localhost:20331 -w KxDgvEKzgSBPPfuVfw67oPQBSjidEiqTHURKSDL1R7yGaGYAeYnr -g 100
```

Where
- `contract deploy` is a command for deployment
- `-i '/1-print.avm'` path to smart contract
- `-c 1-print.yml` configuration input file
- `-e http://localhost:20331` node endpoint
- `-w KxDgvEKzgSBPPfuVfw67oPQBSjidEiqTHURKSDL1R7yGaGYAeYnr` key to sign deployed transaction
- `-g 100` amount of gas to pay for contract deployment

Result:
```
Sent deployment transaction 26d0a206e724e402ee1d4bcd794e82e43ca436888c50dbc5a2216e1ba08ecd0d for contract 6d1eeca891ee93de2b7a77eb91c26f3b3c04d6cf
```

At this point your ‘Hello World’ contract is deployed and could be invoked. Let’s do it as a final step.

#### Step 7
Invoke contract.
```
$ ./bin/neo-go contract invokefunction -e http://localhost:20331 -w KxDgvEKzgSBPPfuVfw67oPQBSjidEiqTHURKSDL1R7yGaGYAeYnr -g 0.00001 6d1eeca891ee93de2b7a77eb91c26f3b3c04d6cf
```

Where
- `contract invokefunction` runs invoke with provided parameters
- `-e http://localhost:20331` defines RPC endpoint used for function call
- `-w KxDgvEKzgSBPPfuVfw67oPQBSjidEiqTHURKSDL1R7yGaGYAeYnr` is a wallet
- `-g 0.00001` defines amount of GAS to be used for invoke operation
- `6d1eeca891ee93de2b7a77eb91c26f3b3c04d6cf` contract hash got as an output from the previous command (deployment in step 6)

Result:
In the console where you were running step #5 you will get:
```
INFO[0227] script cfd6043c3b6fc291eb777a2bde93ee91a8ec1e6d logs: "Hello, world!"
```
Which means that this contract was executed.

This is it. There are only 5 steps to make deployment and they look easy, aren’t they?
Thank you!

