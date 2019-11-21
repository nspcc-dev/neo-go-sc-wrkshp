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
<p align="center">
<img src="https://gfycat.com/legitimatepertinentindianrhinoceros" width="300px" alt="network_monitor">
</p>

## Workshop
Now it’s time to run your private network. Connect neo-go node to it, write smart contract and deploy it. Let’s go!

#### Step 1
Run local privnet:
```
git clone https://github.com/CityOfZion/neo-local.git
$ cd neo-local
$ git checkout -b 4nodes 0.12
$ make start
```

#### Step 2
Download neo-go
```
https://github.com/nspcc-dev/neo-go.git
$ make build
```

#### Step 3
Create basic Hello World smart contract(or use the one presented in this repo):
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

#### Step 4 
Run Hello World smart contract:
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
Copy smart contract to docker environment:

`$ sudo docker cp smart-contracts/1-print.avm neo-python:/neo-python`

#### Step 6
Deploy smart contract:
```
./bin/neo-go contract deploy -i 1-print.avm -c neo-go.yml -e 
http://localhost:20331 -w KxDgvEKzgSBPPfuVfw67oPQBSjidEiqTHURKSDL1R7yGaGYAeYnr -g 100
```

Where
- `contract deploy` is a command for deployment
- `-i '/1-print.avm'` path to smart contract
- `-c neo-go.yml` configuration input file
- `-e http://localhost:20331` node endpoint
- `-w KxDgvEKzgSBPPfuVfw67oPQBSjidEiqTHURKSDL1R7yGaGYAeYnr` key to sign deployed transaction
- `-g 100` amount of gas to pay for contract deployment

Result:
```
Sent deployment transaction 98d33630d98fa6e171c2659bf9028497574aca9ccf3f398624173b7d445fc0d6 for contract 50befd26fdf6e4d957c11e078b24ebce6291456f
```
#### Step 7
Invoke contract from neo-python environment.
```
sc invoke 0x4cf87bd149748abc3cf5fe46040bb8b4c2c0b2c7
```

Where `0x4cf87bd149748abc3cf5fe46040bb8b4c2c0b2c7` is a smart contract hashcode which you can get from the log in the previous step.

Result:
```
Used 500.0 Gas 

-------------------------------------------------------------------------------------------------------------------------------------
Test deploy invoke successful
Total operations executed: 11
Results:
[<neo.Core.State.ContractState.ContractState object at 0x7f4164c67358>]
Deploy Invoke TX GAS cost: 490.0
Deploy Invoke TX Fee: 0.0
-------------------------------------------------------------------------------------------------------------------------------------

Priority Fee (1.5) + Deploy Invoke TX Fee (0.0) = 1.5

Enter your password to continue and deploy this contract
```

Enter `coz` password to continue.
After you will see printed `Hello world!`

This is it. There are only 5 steps to make deployment and they look easy, aren’t they?
Thank you!

