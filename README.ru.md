<p align="center">
<img src="./pic/neo_color_dark_gopher.png" width="300px" alt="logo">
</p>

[NEO](https://neo.org/) разрабатывает системы умной экономики, и мы в [NEO SPCC](https://nspcc.ru/) помогаем им с этой нелегкой задачей. 
В нашем блоге вы можете найти статью [how we run NEOFS public test net](https://medium.com/@neospcc/public-neofs-testnet-launch-18f6315c5ced), 
но это не единственная вещь, над которой мы работаем.

## NEO GO
Как вы знаете, сеть состоит из нод. В текущий момент ноды имеют несколько реализаций:
- https://github.com/neo-project/neo.git
- https://github.com/CityOfZion/neo-python
- https://github.com/nspcc-dev/neo-go

Данная статья посвящена последней реализации, поскольку мы в NEO SPCC занимаемся ее разработкой.
Мы надеемся, что данная статья позволит вам понять, как устроена нода neo-go, и поможет научиться писать и разворачивать смрт-контракты.

## Что такое нода?

<p align="center">
<img src="./pic/node.png" width="300px" alt="node">
</p>

Главная цель нод - взаимодействовие друг с другом по протоколу P2P и синхронизация блоков в сети.
Кроме того, ноды позволяют пользователям компилировать и запускать смарт-контракты в сети блокчейн.
Нода состоит из Клиента (CLI), Сетевого слоя, Консенсуса, Виртуальной машины, Компилятора и Блокчейна.
Рассмотрим каждую компоненту более подробно.

#### Клиент
Клиент (CLI) позволяет пользователям запускать команды в терминале. Команды делятся на 4 категории:

- серверные операции
- операции со смарт контрактами
- операции виртуальной машины
- wallet-операции

Например, чтобы подключить ноду к запущенной частной сети (Private Network), вы можете использовать следующую команду:
```
 go run cli/main.go node -p
```
[Здесь](https://medium.com/@neospcc/neo-privatenet-auto-import-of-a-smart-contract-dbf2b9220ad2) вы можете найти больше информации о Private Network и ее запуске. Проще говоря, private network - это сеть, которую вы можете запустить локально. 

Другой пример использования CLI - компиляция смарт-контракта:
```
$ ./bin/neo-go vm 

    _   ____________        __________      _    ____  ___
   / | / / ____/ __ \      / ____/ __ \    | |  / /  |/  /
  /  |/ / __/ / / / /_____/ / __/ / / /____| | / / /|_/ / 
 / /|  / /___/ /_/ /_____/ /_/ / /_/ /_____/ |/ / /  / /  
/_/ |_/_____/\____/      \____/\____/      |___/_/  /_/   


NEO-GO-VM >  
```
После запуска данной команды мы можем взаимодействовать с виртуальной машиной. 
Для получения списка поддерживаемых операций используйте `help`:
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
Как вы видите, тут есть с чем поэкспериментировать. Давайте создадим простой смарт-контракт `1-print.go`и скомпилируем его:
 
```
package main

import (
	"github.com/CityOfZion/neo-go/pkg/interop/runtime"
)

func Main() {
	runtime.Log("Hello, world!")
}
``` 
Используйте команду `loadgo` для компиляции:
```
NEO-GO-VM > loadgo test.go
READY: loaded 38 instructions
NEO-GO-VM 0 >  
```
Теперь вы можете увидеть, сколько инструкций было сгенерировано. Также вы можете получить опкоды (opcodes) данной программы:
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
Этот скомпилированный контракт пригодится нам позже =).
Больше информации об использовании CLI [здесь](https://github.com/nspcc-dev/neo-go/blob/master/docs/cli.md).

#### Сетевой слой
Network-слой - один из самых важных частей ноды. В нашем случае поддерживаются два протокола: протокол P2P позволяет нодам взаимодействовать друг с другом, а протокол RPC используется для получения информации о балансе, аккаунтах, текущем состоянии чейна и т.д.
Здесь вы найдете поддерживаемые [вызовы RPC](https://github.com/nspcc-dev/neo-go/blob/master/docs/rpc.md).

#### Консенсус
Консенсус - это механизм, позволяющий нодам приходить к общему значению (блокам в случае блокчейна). Мы используем нашу собственную реализацию алгоритма dBFT.

#### Компилятор
Компилятор позволяет генерировать байт-код, благодаря чему вы можете писать смарт-контракты на Golang. Все примеры в данном воркшопе были сгенерированы компилятором.

#### Виртуальная машина
Виртуальная машина запускает скомпилированный байт-код. Виртуальная машина Neo является [стековой](https://docs.neo.org/docs/en-us/basic/technology/neovm.html). Для вычислений в ней содержится два стека.

#### Блокчейн
Блокчейн - достаточно большая часть NEO GO, содержащая в себе операции по принятию и валидации транзакций, их подписи,
работе с аккаунтами, ассетами, хранению блоков в базе данных (или в кэше).

#### Сеть
Существует 3 типа сетей.
Частная сеть (Private net) - это сеть, которую вы можете запустить локально. Тестовая сеть (Testnet) и Основная сеть (Mainnet) - сети, в которых запущены большинство нод NEO по всему миру.
Каждую ноду, запущенную в сети блокчейн, вы можете найти в [Neo Monitor](http://monitor.cityofzion.io/)

## Воркшоп. Часть 1.
В этой части мы запустим вашу private network, подключим ноду neo-go к ней, напишем, развернем и вызовем смарт-контракт. Начнем!

#### Требования
Для этого воркшопа у вам понадобятся установленные Debian 10, Docker, docker-compose и go:
- [docker](https://docs.docker.com/install/linux/docker-ce/debian/)
- [go](https://golang.org/dl/)

#### Шаг 1
Скачайте neo-go и соберите проект:
```
git clone https://github.com/nspcc-dev/neo-go.git
$ cd neo-go
$ make build 
```

#### Шаг 2
Существует два способа запуска локальной private network. 
Первый - запуск neo-local private network, второй - запуск neo-go private network.

#### Запуск neo-go private network
```
$ make env_image
$ make env_up
```

В результате должна запуститься локальная privatenet:
```
=> Bootup environment
Creating network "neo_go_network" with the default driver
Creating volume "docker_volume_chain" with local driver
Creating neo_go_node_four  ... done
Creating neo_go_node_two   ... done
Creating neo_go_node_one   ... done
Creating neo_go_node_three ... done
```

Для остановки используйте:
```
$ make env_down
```

#### Запуск neo local private network
```
git clone https://github.com/CityOfZion/neo-local.git
$ cd neo-local
$ git checkout -b 4nodes 0.12
$ make start
```

#### Шаг 3
Создайте простой смарт-контракт "Hello World" (или используйте представленный в репозтитории воркшопа):
```
package main

import (
	"github.com/CityOfZion/neo-go/pkg/interop/runtime"
)

func Main() {
	runtime.Log("Hello, world!")
}
```
И сохраните его как `1-print.go`.

Создайте конфигурацию для смарт-контракта:
https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/1-print.yml

#### Шаг 4 
Запустите смарт-контракт "Hello World":
```
$ ./bin/neo-go contract compile -i '/1-print.go'
```

Где
- `./bin/neo-go` запускает neo-go
- `contract compile` команда с аргументами из [neo-go](https://github.com/nspcc-dev/neo-go/blob/master/cli/smartcontract/smart_contract.go#L43)
- `-i '/1-print.go'` путь к смарт-контракту

Результат: 
Скомпилированный смарт-контракт `1-pring.avm`

Для просмотра опкодов вы можете воспользоваться командой:
```
$ ./bin/neo-go contract inspect -i '/1-print.go'
```

#### Шаг 5
Запустите ноду neo-go, которая подключится к запущенной ранее privatenet:
```
$ ./bin/neo-go node --privnet
```

Результат:
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

#### Шаг 6
Разверните смарт-контракт:
```
./bin/neo-go contract deploy -i 1-print.avm -c 1-print.yml -e http://localhost:20331 -w my_wallet.json -g 0.001
```

Где
- `contract deploy` - команда для развертывания
- `-i '/1-print.avm'` - путть к смарт-контракту
- `-c 1-print.yml` - файл конфигурации
- `-e http://localhost:20331` - эндпоинт ноды
- `-w my_wallet.json` - бумажник, в котором хранится ключ для подписи транзакции (вы можете взять его из репозитория воркшопа)
- `-g 0.001` - количество газа для оплаты развертывания контракта

Введите пароль `qwerty` для аккаунта:
```
Enter account AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y password >
```

Результат:
```
Sent deployment transaction ea93196802fe3517d2d028e4d4f244aa734a9b1988456740d96f2c3336140fda for contract 6d1eeca891ee93de2b7a77eb91c26f3b3c04d6cf
```

На данном этапе ваш контракт ‘Hello World’ развернут и может быть вызван. В следующем шаге вызовем этот контракт.

#### Шаг 7
Вызовите контракт.
```
$ ./bin/neo-go contract invokefunction -e http://localhost:20331 -w my_wallet.json -g 0.00001 6d1eeca891ee93de2b7a77eb91c26f3b3c04d6cf
```

Где
- `contract invokefunction` запускает вызов контракта с заданными параметрами
- `-e http://localhost:20331` определяет эндпоинт RPC, используемый для вызова функции
- `-w my_wallet.json` - бумажник
- `-g 0.00001` определяет количество газа, которое будет использовано для вызова смарт-контракта
- `6d1eeca891ee93de2b7a77eb91c26f3b3c04d6cf` хеш контракта, полученный в результате выполнения предыдущей команды (развертывание из шага 6)

Введите пароль `qwerty` для аккаунта:
```
Enter account AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y password >
```

Результат:

В консоли, где была запущена нода (шаг 5), вы увидите:
```
INFO[0227] script cfd6043c3b6fc291eb777a2bde93ee91a8ec1e6d logs: "Hello, world!"
```
Что означает, что контракт был выполнен.

На этом все. Вам потребовалось всего 5 шагов, чтобы развернуть свой контракт, и это оказалось довольно легко.
Спасибо!

## Воркшоп. Часть 2
В этой части мы выполним несколько RPC вызовов и попробуем написать, задеплоить и вызвать смарт-контракт, использующий хранилище. Начнем!

### Вызовы RPC
Давайте рассмотрим более детально, что происходит с нашим смарт-контрактом при развертывании и вызове. 
Каждая нода neo-go предоставляет API интерфейс для получения данных о блокчейне.
Данное взаимодействие осуществляется по протоколу `JSON-RPC`, использующему HTTP для общения.

Полный `NEO JSON-RPC 2.0 API` описан [здесь](https://docs.neo.org/docs/en-us/reference/rpc/latest-version/api.html).

RPC-сервер ноды neo-go, запущенной на шаге 5, доступен по `localhost:20331`. Давайте выполним несколько вызовов RPC.

#### GetRawTransaction
[GetRawTransaction](https://docs.neo.org/docs/en-us/reference/rpc/latest-version/api/getrawtransaction.html) возвращает информацию о транзакции по ее хешу.

Запросите информацию о нашей разворачивающей транзакции из шага 6:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getrawtransaction", "params": ["ea93196802fe3517d2d028e4d4f244aa734a9b1988456740d96f2c3336140fda", 1] }' localhost:20331 | json_pp
```

Где:
- `"jsonrpc": "2.0"` - версия протокола
- `"id": 1` - id текущего запроса
- `"method": "getrawtransaction"` - запрашиваемый метод
- `"params": ["ea93196802fe3517d2d028e4d4f244aa734a9b1988456740d96f2c3336140fda", 1]` массив параметров запроса, где
   - `ea93196802fe3517d2d028e4d4f244aa734a9b1988456740d96f2c3336140fda` - хеш разворачивающеей транзакции, полученный после выполнения шага 6
   - `1` это `verbose` параметр для получения детального ответа в формате json-строки 

Результат:
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
[GetApplicationLog](https://docs.neo.org/docs/en-us/reference/rpc/latest-version/api/getapplicationlog.html) возвращает лог контракта по соответствующему хешу транзакции.

Запросите информацию о контракте для нашей вызывающей транзакции, полученной на шаге 7:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["42879c901de7d180f2ff609aded2da2b9c63d5ed7d1b22f18e246c16fc5d0c57"] }' localhost:20331 | json_pp
```

Где в качестве параметра:
- `42879c901de7d180f2ff609aded2da2b9c63d5ed7d1b22f18e246c16fc5d0c57` - хеш вызывающей транзакции из шага 7

Результат:
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

#### Другие полезные вызовы RPC
```
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getversion", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getblockcount", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getconnectioncount", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getaccountstate", "params": ["AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y"] }' localhost:20331
```

Список всех поддерживаемых нодой neo-go вызовов RPC вы найдете [здесь](https://github.com/nspcc-dev/neo-go/blob/master/docs/rpc.md#supported-methods).

### Смарт-контракт, использующий хранилище

Давайте изучим еще один пример смарт-контракта: [2-storage.go](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/2-storage.go).
Он достаточно простой и, так же как предыдущий, не принимает никаких аргументов.
С другой стороны, этот контракт умеет считать количество его вызовов, сохраняя целое число и увеличивая его на 1 после каждого вызова.
Подобный контракт будет интересен нам, поскольку он способен *хранить* значения, т.е. обладает *хранилищем*, которое является общим для всех вызовов данного контракта.

К сожалению, за все хорошее нужно платить, в том числе и за наличие хранилища у нашего контракта.
Чтобы обозначить, что контракт имеет хранилище, в его конфигурации [2-storage.yml](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/2-storage.yml) мы обязаны установить значение следующего флага:
```
hasstorage: true
```
В противном случае мы не сможем воспользоваться хранилищем в контракте.

Теперь, когда мы узнали о хранилище, давайте скомпилируем, развернем и вызовем смарт-контракт.

#### Шаг #1
Скомпилируйте смарт-контракт [2-storage.go](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/2-storage.go):
```
./bin/neo-go contract compile -i 2-storage.go
```

Результат:

Скомпилированный смарт-контракт: `2-storage.avm`

#### Шаг #2
Разверните скомпилированный смарт-контракт с [конфигурацией](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/2-storage.yml):
```
./bin/neo-go contract deploy -i 2-storage.avm -c 2-storage.yaml -e http://localhost:20331 -w my_wallet.json -g 0.001
```
... введите пароль `qwerty`:
```
Enter account AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y password >
```

Результат:
```
Sent deployment transaction d79e5e39d87c2911b623d3efe98842cfde41eddc34d85cee394783b7320813e8 for contract 85cf2075f3e297d489ff3c4c1745ca80d44e2a68
```   

Что означает, что наш контракт развернуть и теперь мы можем вызывать его.

#### Шаг #3
Поскольку мы не вызывали наш смарт-контракт раньше, в его хранилище нет никаких значений, поэтому при первом вызове он должен создать новое значение (равное `1`) и положить его в хранилище.
Давайте проверим:
```
./bin/neo-go contract invokefunction -e http://localhost:20331 -w my_wallet.json -g 0.00001 85cf2075f3e297d489ff3c4c1745ca80d44e2a68
```
... введите пароль `qwerty`:
```
Enter account AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y password >
```
Результат:
```
Sent invocation transaction 6f27f523da9c71297f4a81a254274b0c2a78f893b81c500429be6230254be0bf
```
Для проверки значения счетчика вызовем `getapplicaionlog` вызов RPC для вызывающей транзакции:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["6f27f523da9c71297f4a81a254274b0c2a78f893b81c500429be6230254be0bf"] }' localhost:20331 | json_pp
```
Результат:
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
Обратите внимание на поле `notification`. Оно содержит сообщения, переданные методу `runtime.Notify`.
В нашем случае у нем находятся три шестнадцатеричных массива байт, которые можно декодировать в следующие сообщения:
  - `Value read from storage`, которое было вызвано после того как мы попытались достать значение счетчика из хранилища
  - `Storage key not yet set. Setting to 1`, которое было вызвано после того, как мы поняли, что полученное значение = 0
  - `New value written into storage`, которое было вызвано после того, как мы записали новое значение в хранилище
  
И последняя часть - поле `stack`. Данное поле содержит все возвращенные контрактом значения, поэтому здесь вы можете увидеть целое `1`,
которое является значением счетчика, определяющего количество вызовов смарт-контракта.

#### Шаг #4
Для того чтобы убедиться, что все работает как надо, давайте вызовем наш контракт еще раз и проверим, что счетчик будет увеличен: 
```
./bin/neo-go contract invokefunction -e http://localhost:20331 -w my_wallet.json -g 0.00001 85cf2075f3e297d489ff3c4c1745ca80d44e2a68
```
... введите пароль `qwerty`:
```
Enter account AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y password >
```
Результат:
```
Sent invocation transaction dcf53cdb69c816f0f9ab27ba509fb656b6eddd2dad5a658e9f534e9dba38462b
```
Для проверки значения счетчика, выполните `getapplicaionlog` вызов RPC для вызывающей транзакции:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["dcf53cdb69c816f0f9ab27ba509fb656b6eddd2dad5a658e9f534e9dba38462b"] }' localhost:20331 | json_pp
```
Результат:
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

Теперь поле `stack` содержит значение `2` - счетчик был увеличен, как мы и ожидали.

## Воркшоп. Часть 3
В этой части мы узнаем о стандарте токена NEP5 и попробуем написать, задеплоить и вызвать более сложный смарт-контракт. Начнем!

### NEP5
[NEP5](https://docs.neo.org/docs/en-us/sc/write/nep5.html) - это стандарт токена блокчейна Neo, обеспечивающий системы обобщенным механизмом взаимодействия для токенизированных смарт-контрактов.
Пример с реализацией всех требуемых стандартом методов вы можете найти в [nep5.go](https://github.com/nspcc-dev/neo-go/blob/master/examples/token/nep5/nep5.go)
 
Давайте посмотрим на пример смарт-контракта NEP5: [token.go](https://github.com/nspcc-dev/neo-go/blob/master/examples/token/token.go)
 
Этот смарт-контракт принимает в качестве параетра строку с операцией, которая может принимать следующие значения:
- `name` возвращает имя созданного токена nep5 
- `symbol` возвращает код токена
- `decimals` возвращает количество десятичных знаков токена
- `totalSupply` возвращает общий множитель * токена
- `balanceOf` возвращает баланс токена, находящегося на указанном адресе и требует дополнительног аргумента:
  - `holder` адрес запрашиваемог аккаунта
- `transfer` перемещает токен от одного пользователя к другому и требует дополнительных аргументов:
  - `from` адрес аккаунта, с которого будет списан токен
  - `to` адрес аккаунта, на который будет зачислен токен
  - `amount` количество токена для перевода
- `mint` выпускает начальное количество токенов на аккаунт и требует дополнительных аргументов:
  - `to` адрес аккаунта, на который вы бы хотели выпустить токены  

Давайте проведем несколько операций с помощью этого контракта.

#### Шаг #1
Скомпилируйте смарт-контракт [token.go](https://github.com/nspcc-dev/neo-go/blob/master/examples/token/token.go):
```
./bin/neo-go contract compile -i examples/token/token.go
```
Для развертывания можно использовать отредактированный файл [конфигурации](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/1-print.yml) из части 1 со следующими изменениями:
поскольку данный контракт использует storage, необходимо установить флаг 
```
hasstorage: true
```
Разверните скомпилированный контракт с отредактированной конфигурацией:
```
./bin/neo-go contract deploy -i examples/token/token.avm -c nep5.yml -e http://localhost:20331 -w my_wallet.json -g 0.001
```
... введите пароль `qwerty`:
```
Enter account AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y password >
```

Результат:
```
Sent deployment transaction 6e46e2f2c8e799bf0fce679fa3c8acbc046f595252ae58b39abea839da01067d for contract f84d6a337fbc3d3a201d41da99e86b479e7a2554
```   

Что означает, что наш контракт был развернут, и теперь мы можем вызывать его.

#### Шаг #2
Давайте вызовем контракт для осуществления операций с nep5.

Для начала, запросите имя созданного токена nep5:

```
./bin/neo-go contract invokefunction -e http://localhost:20331 -w my_wallet.json -g 0.00001 f84d6a337fbc3d3a201d41da99e86b479e7a2554 name
```                                                                   
Где
- `f84d6a337fbc3d3a201d41da99e86b479e7a2554` - хеш нашего контракта, полученный на шаге #1.
- `name` - строка операции, описанная ранее и возвращающая имя токена.

... не забудьте пароль от аккаунта `qwerty`.

Результат:
```
Sent invocation transaction cfdf4c2883d71a4375ab94fbb302a386828e7934541aa222fc38b8bc67e6a2b4
```                                                                                         
Теперь давайте подробнее посмотрим на полученную вызывающую транзакцию с помощью `getapplicationlog` RPC-вызова:

```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["cfdf4c2883d71a4375ab94fbb302a386828e7934541aa222fc38b8bc67e6a2b4"] }' localhost:20331 | json_pp
```               

Результат:
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

Поле `stack` полученного JSON-сообщения содержит массив байтов со значением имени токена.

Следующие команды позволят получить вам дополнительную информацию о токене:
```
./bin/neo-go contract invokefunction -e http://localhost:20331 -w my_wallet.json -g 0.00001 f84d6a337fbc3d3a201d41da99e86b479e7a2554 symbol
./bin/neo-go contract invokefunction -e http://localhost:20331 -w my_wallet.json -g 0.00001 f84d6a337fbc3d3a201d41da99e86b479e7a2554 decimals
./bin/neo-go contract invokefunction -e http://localhost:20331 -w my_wallet.json -g 0.00001 f84d6a337fbc3d3a201d41da99e86b479e7a2554 totalSupply
```

#### Шаг #3
Настало время для более интересных вещей. Для начала проверим баланс nep5 токенов на нашем счету с помощью метода `balanceOf`:
```
./bin/neo-go contract invokefunction -e http://localhost:20331 -w my_wallet.json -g 0.00001 f84d6a337fbc3d3a201d41da99e86b479e7a2554 balanceOf AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y
```                             
... с паролем `qwerty`. Результат:
```
Sent invocation transaction 30de65dec667d68ca1b385c590c9c5fccf82e2d4831540e6bb5875afa57c5cbe
```
Для более детального рассмотрения транзакции используем `getapplicationlog` RPC-вызов:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["30de65dec667d68ca1b385c590c9c5fccf82e2d4831540e6bb5875afa57c5cbe"] }' localhost:20331 | json_pp
```
Результат:
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
Как вы видите, поле `stack` содержит пустой массив, то есть в настоящий момент мы не обладаем токенами.
Но не стоит об этом беспокоиться, переходите к следующему шагу.

#### Шаг #4

Перед тем как мы будем способны использовать наш токен (например, попытаемся передать его кому-либо), мы должны его *выпустить*.
Другими словами, мы должны перевести все имеющееся количество токена (total supply) на чей-нибудь аккаунт.
Для этого в нашем контракте существует специальная функция - `mint`. Давайте выпустим токен на наш адрес:
```
./bin/neo-go contract invokefunction -e http://localhost:20331 -w my_wallet.json -g 0.00001 f84d6a337fbc3d3a201d41da99e86b479e7a2554 mint AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y
```
... с паролем `qwerty`. Результат:
``` 
Sent invocation transaction a571adebdbdabfd087f34867da94524649003c6b851ed0cc5da7a30ff843bc1e
```
`getapplicationlog` RPC-вызов для этой транзакции даст нам следующее:
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
Обратите внимание, что поле `stack` содержит значение `1` - токен был успешно выпущен.
Давайте убедимся в этом, еще раз запросив баланс нашего аккаунта с помощью метода `balanceOf`:
```
./bin/neo-go contract invokefunction -e http://localhost:20331 -w my_wallet.json -g 0.00001 f84d6a337fbc3d3a201d41da99e86b479e7a2554 balanceOf AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y
```
... пароль `qwerty`. Результат:
``` 
Sent invocation transaction c56f469cd9d47c6a4195a742752621c4898447aa4cfd7f550046bdd10d297c12
```
... со следующим сообщением от `getapplicationlog` вызова RPC:
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
Теперь мы видим непустой массив байт в поле `stack`, а именно, `00c040b571e803` является шестнадцатеричным представлением баланса токена nep5 на нашем аккаунте.

Важно, что токен может быть выпущен лишь однажды.

#### Шаг #5

После того, как мы закончили с выпуском токена, мы можем перевести некоторое количество токена кому-нибудь.
Давайте переведем 5 токенов аккаунту с адресом `AKkkumHbBipZ46UMZJoFynJMXzSRnBvKcs` с помощью функции `transfer`:
```
./bin/neo-go contract invokefunction -e http://localhost:20331 -w my_wallet.json -g 0.00001 f84d6a337fbc3d3a201d41da99e86b479e7a2554 transfer AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y AKkkumHbBipZ46UMZJoFynJMXzSRnBvKcs 5
```
... пароль `qwerty`. Результат:
``` 
Sent invocation transaction 11a272f4a0d7912f7219979bab7d094df3b404b89e903337ee72a90249cc448d
```
Наш любимый вызов RPC `getapplicationlog` говорит нам:
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
Заметьте, что поле `stack` содержит `1`, что означает, что токен был успешно переведен с нашего аккаунта.
Теперь давайте проверим баланс аккаунта, на который был совершен перевод (`AKkkumHbBipZ46UMZJoFynJMXzSRnBvKcs`), чтобы убедиться, что количество токена на нем = 5:
```
./bin/neo-go contract invokefunction -e http://localhost:20331 -w my_wallet.json -g 0.00001 f84d6a337fbc3d3a201d41da99e86b479e7a2554 balanceOf AKkkumHbBipZ46UMZJoFynJMXzSRnBvKcs
```
Вызов RPC `getapplicationlog` для этой транзакции возвращает следующий результат:
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
Как и ожидалось, мы видим ровно 5 токенов в поле `stack`.
Вы можете самостоятельно убедиться, что с нашего аккаунта были списаны 5 токенов, выполнив метод `balanceOf` с аргументом `AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y`.

## Воркшоп. Часть 4
В этой части подытожим наши знания о смарт-контрактах и исследуем смарт-контракт [4-domain.go](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/4-domain.go).
Данный контракт описывает операции регистрации, переноса и удаления доменов, а также операцию получения информации о зарегистрированном домене.

Начнем!

#### Шаг #1
Давайте рассмотрим и исследуем [смарт-контракт](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/4-domain.go). В качестве первого параметра контракт принимает на вход строку - действие, одно из следующих значений:
- `register` проверяет, существует ли домен с указанным именем. В случае, если такого домена не существует, добавляет пару `[domain_name, owner]` в хранилище. Данная операция требудет дополнительных аргументов:
   - `domain_name` - новое имя домена.
   - `owner` - 34-значный адрес аккаунта из нашего [wallet](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/my_wallet.json), который будет использоваться для вызова контракта.
- `query` возвращает адрес аккаунта владельца запрашиваемого домена (или false, в случае, если домен с указанным именем не зарегистрирован). Требует дополнительных аргументов:
   - `domain_name` - имя запрашиваемого домена.
- `transfer` переводит домен с указанным именем на другой адрес (в случае, если вы являетесь владельцем указанного домена). Требует следующих аргументов:
   - `domain_name` - имя домена, который вы хотите перевести.
   - `to_address` - адрес аккаунта, на который вы хотите перевести домен.
- `delete` удаляет домен из хранилища. Аргументы:
   - `domain_name` имя домента, который вы хотите удалить.
 
 
 В следующих шагах мы скомпилируем и развернем смарт-контракт.
 После этого мы зарегистрируем новый домен, переведем его на другой аккаунт и запросим информацию о нем.

#### Шаг #2

Скомпилируйте смарт-контракт [4-domain.go](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/4-domain.go)
```
./bin/neo-go contract compile -i 4-domain.go
```

... и разверните его с [конфигурацией](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/4-domain.yml):
```
./bin/neo-go contract deploy -i 4-domain.avm -c 4-domain.yml -e http://localhost:20331 -w my_wallet.json -g 0.001
```
Обратите внимание, что наш контракт использует хранилище и, как и с предыдущим контрактом, необходимо, чтобы флаг `hasstorage` имел значение `true`.
Этот флаг указывается в файле [конфигурации](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/4-domain.yml).

... введите пароль `qwerty`:
```
Enter account AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y password >
```

Результат:
```
Sent deployment transaction e2647e15a855dd174984fde05a3a193cf8c45deb5c741d3045cb0b35e9e462ce for contract 96f996eb681f65d380c596949d33d7f29897ad27
```   
Вы догадываетесь, что это значит :)

#### Шаг #3

Вызовите контракт, чтобы зарегистрировать домен с именем `my_first_domain`: 
```
./bin/neo-go contract invokefunction -e http://localhost:20331 -w my_wallet.json -g 0.00001 96f996eb681f65d380c596949d33d7f29897ad27 register my_first_domain AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y
```
... пароль: `qwerty`
```
Enter account AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y password >
```
Результат:
```
Sent invocation transaction 2d91a11f47217b7d7346d712456db79bdcaf982eea76ce105b41de859f7fc2e1
```
Также вы можете увидеть лог-сообщение в консоли, где запускали ноду neo-go:
```
2020-04-06T11:57:26.889+0300	INFO	runtime log	{"script": "27ad9798f2d7339d9496c580d3651f68eb96f996", "logs": "\"RegisterDomain: my_first_domain\""}
```
Все получилось. Теперь проверим, был ли наш домен действительно зарегистрирован, с помощью вызова RPC `getapplicationlog`:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["2d91a11f47217b7d7346d712456db79bdcaf982eea76ce105b41de859f7fc2e1"] }' localhost:20331 | json_pp
```
Результат:
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
Здесь мы в особенности заинтересованы в двух полях полученного json:

Первое поле - `notifications`, оно содержит 3 значния:
- `72656769737465726564` - массив байт. С помощью следующего скрипта можно декодировать это значение:
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
... что позволит нам увидеть сообщение, содержащееся в уведомлении: `registered`.

- `6d795f66697273745f646f6d61696e` - массив байт, который может быть декодирован в имя нашего домена - `my_first_domain`,
- `23ba2703c53263e8d6e522dc32203339dcd8eee9` - массив байт, который декодируется в адрес нашего аккаунта `AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y` с помощью следующего скрипта:
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

Второе поле - `stack`, в котором лежит `1` - значение, возвращенное смарт-контрактом.

Все эти значения дают нам понять, что наш домен был успещно зарегистрирован.    

#### Шаг #4

Вызовите контракт, чтобы запросить информацию об адресе аккаунта, зарегистрировавшего домен `my_first_domain`:
```
./bin/neo-go contract invokefunction -e http://localhost:20331 -w my_wallet.json -g 0.00001 96f996eb681f65d380c596949d33d7f29897ad27 query my_first_domain
```
... любимейший пароль `qwerty`:
```
Enter account AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y password >
```
Результат:
```
Sent invocation transaction 35c72ab78aa261242371426d2c33433312c291ef886be4f78e0ccbbe34d17349
```
и лог-сообщение в консоли запущенной ноды neo-go:
```
2020-04-06T12:01:27.212+0300	INFO	runtime log	{"script": "27ad9798f2d7339d9496c580d3651f68eb96f996", "logs": "\"QueryDomain: my_first_domain\""}
```
Проверим транзакцию с помощью вызова RPC `getapplicationlog`:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["35c72ab78aa261242371426d2c33433312c291ef886be4f78e0ccbbe34d17349"] }' localhost:20331 | json_pp
```
... что даст нам следующий результат:
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

с шестнадцатеричным представлением адреса аккаунта `AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y` на стеке и в уведомлениях, что означает, что домен `my_first_domain` был зарегистрирован владельцем с полученным адресом аккаунта.

#### Шаг #5

Вызовите контракт для передачи домена другому аккаунту (например, аккаунту с адресом `AKkkumHbBipZ46UMZJoFynJMXzSRnBvKcs`):
```
./bin/neo-go contract invokefunction -e http://localhost:20331 -w my_wallet.json -g 0.00001 96f996eb681f65d380c596949d33d7f29897ad27 transfer my_first_domain AKkkumHbBipZ46UMZJoFynJMXzSRnBvKcs
```
... пароль: `qwerty`
```
Enter account AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y password >
```
Результат:
```
Sent invocation transaction da3653ab7eb23eab511b05ea64c2cdc3872b0c576d86b924a7738986cfd63462
```
и лог-сообщение:
```
2020-04-06T12:04:42.483+0300	INFO	runtime log	{"script": "27ad9798f2d7339d9496c580d3651f68eb96f996", "logs": "\"TransferDomain: my_first_domain\""}
```
Отлично. И `getapplicationlog` вызов RPC...
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["da3653ab7eb23eab511b05ea64c2cdc3872b0c576d86b924a7738986cfd63462"] }' localhost:20331 | json_pp
```
... говорит нам:
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
Достаточно детальный json. Поле `notifications` содержит два значения:
- массив с полем `64656c65746564`, декодируемым в `deleted`, и полями с дополнительной информацией (домен `my_first_domain` был удален с аккаунта с адресом `AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y`),
- массив с полем `72656769737465726564`, декодируемым в `registered`, и полями с дополнительной информацией (домен `my_first_domain` бфл зарегистрирован аккаунтом с адресом `AKkkumHbBipZ46UMZJoFynJMXzSRnBvKcs`).
Поле `stack` содержит `1`, что значит, что домен был успешно перемещен.

#### Шаг #6

Оставшийся вызов - `delete`, вы можете попробовать выполнить его самостоятельно, создав перед этим еще один домен, например, с именем `my_second_domain`, а затем удалить его из хранилища с помощью:
```
./bin/neo-go contract invokefunction -e http://localhost:20331 -w my_wallet.json -g 0.00001 96f996eb681f65d380c596949d33d7f29897ad27 delete my_second_domain
```

Спасибо!

### Полезные ссылки

* [Наш воркшоп на Medium](https://medium.com/@neospcc/%D1%81%D0%BC%D0%B0%D1%80%D1%82-%D0%BA%D0%BE%D0%BD%D1%82%D1%80%D0%B0%D0%BA%D1%82-%D0%B4%D0%BB%D1%8F-neo-769139352b65)
* [Документация NEO](https://docs.neo.org/)
* [NEO github](https://github.com/neo-project/neo/)
* [NEO-GO github](https://github.com/nspcc-dev/neo-go)
