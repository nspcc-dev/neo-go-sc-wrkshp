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

#### Step 3
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

#### Step 4 
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
./bin/neo-go contract deploy -i 1-print.avm -c 1-print.yml -e http://localhost:20331 -w my_wallet.json -g 100
```

Где
- `contract deploy` - команда для развертывания
- `-i '/1-print.avm'` - путть к смарт-контракту
- `-c 1-print.yml` - файл конфигурации
- `-e http://localhost:20331` - эндпоинт ноды
- `-w my_wallet.json` - бумажник, в котором хранится ключ для подписи транзакции (вы можете взять его из репозитория воркшопа)
- `-g 100` - количество газа для оплаты развертывания контракта

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
В этой части мы выполним несколько RPC вызовов и попробуем написать, задеплоить и вызвать более сложный смарт-контракт. Начнем!

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

### Полезные ссылки

* [Наш воркшоп на Medium](https://medium.com/@neospcc/%D1%81%D0%BC%D0%B0%D1%80%D1%82-%D0%BA%D0%BE%D0%BD%D1%82%D1%80%D0%B0%D0%BA%D1%82-%D0%B4%D0%BB%D1%8F-neo-769139352b65)
* [Документация NEO](https://docs.neo.org/)
* [NEO github](https://github.com/neo-project/neo/)
* [NEO-GO github](https://github.com/nspcc-dev/neo-go)
