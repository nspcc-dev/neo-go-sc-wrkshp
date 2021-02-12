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
Как вы видите, тут есть с чем поэкспериментировать. Давайте создадим простой смарт-контракт `1-print.go` и скомпилируем его:
 
```
package main

import (
	"github.com/nspcc-dev/neo-go/pkg/interop/runtime"
)

func Main() {
	runtime.Log("Hello, world!")
}
``` 
Используйте команду `loadgo` для компиляции:
```
NEO-GO-VM > loadgo 1-print.go
READY: loaded 22 instructions
NEO-GO-VM 0 >  
```
Теперь вы можете увидеть, сколько инструкций было сгенерировано. Также вы можете получить опкоды (opcodes) данной программы:
```
NEO-GO-VM 0 > ops
0        PUSHDATA1    48656c6c6f2c20776f726c6421 ("Hello, world!")    <<
15       SYSCALL      System.Runtime.Log (cfe74796)
20       NOP
21       RET
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

## Воркшоп. Подготовка
В этой части мы настроим окружение: запустим частную сеть, подсоединим к ней ноду neo-go и переведем немного GAS на аккаунт, с который будем использовать далее
для создания транзакций. Давайте начнем.

#### Требования
Для этого воркшопа у вам понадобятся установленные Debian 10, Docker, docker-compose и go:
- [docker](https://docs.docker.com/install/linux/docker-ce/debian/)
- [go](https://golang.org/dl/)

#### Версионирование
Как и многие другие проекты Neo, NeoGo находится на своем пути к Neo 3, поэтому в нем существуют две основные ветки - [master](https://github.com/nspcc-dev/neo-go),
в которой сейчас происходит разработка Neo 3, и [master-2.x](https://github.com/nspcc-dev/neo-go/tree/master-2.x) - стабильная реализация Neo 2. 
Данный воркшоп содержит базовый туториал для Neo 3. 
Если вы хотите продолжить с Neo 2, воспользуйтесь веткой [master-2.x branch](https://github.com/nspcc-dev/neo-go-sc-wrkshp/tree/master-2.x).

#### Шаг 1
Если у вас уже установлен neo-go или есть смарт-контракты на go, пожалуйста, обновите go modules чтобы использовать свежую версию API интеропов.
Если нет, скачайте neo-go и соберите проект (ветку master):
```
$ git clone https://github.com/nspcc-dev/neo-go.git
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
$ git clone https://github.com/CityOfZion/neo-local.git
$ cd neo-local
$ git checkout -b 4nodes 0.12
$ make start
```

#### Шаг 3
Запустите ноду neo-go, которая подключится к запущенной ранее privatenet:
```
$ ./bin/neo-go node --privnet
```

Результат:
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
Переведите немного GAS с мультисигового аккаунта на аккаунт, который мы будем использовать в дальнейшем.

1. Создадим транзакцию перевода GAS токенов:
    ```
        $ ./bin/neo-go wallet nep17 transfer -w .docker/wallets/wallet1.json --out my_tx.json -r http://localhost:20331 --from NUVPACMnKFhpuHjsRjhUvXz1XhqfGZYVtY --to NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt --token GAS --amount 29999999
    ``` 
    Где
    - `./bin/neo-go` запускает neo-go
    - `wallet nep17 transfer` - команда с аргументами в [neo-go](https://github.com/nspcc-dev/neo-go/blob/master/cli/wallet/nep17.go#L108)
    - `-w .docker/wallets/wallet1.json` - путь к [кошельку](https://github.com/nspcc-dev/neo-go/blob/master/.docker/wallets/wallet1.json) для первой ноды в созданной частной сети
    - `--out my_tx.json` - файл для записи подписанной транзакции
    - `-r http://localhost:20331` - RPC-эндпоинт ноды
    - `--from NUVPACMnKFhpuHjsRjhUvXz1XhqfGZYVtY` - мультисиговый аккаунт, являющийся владельцем GAS
    - `--to NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt` - наш аккаунт из [кошелька](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/my_wallet.json)
    - `--token GAS` - имя переводимого токена (в данном случае это GAS)
    - `--amount 29999999` - количество GAS для перевода
    
    Введите пароль `one`:
    ```
    Password >
    ```
    Результатом является транзакция, подписанная первой нодой и сохраненная в `my_tx.json`.

2. Подпишите созданную транзакцию, используя адрес второй ноды:

    ```
    $ ./bin/neo-go wallet multisig sign -w .docker/wallets/wallet2.json --in my_tx.json --out my_tx2.json --address NUVPACMnKFhpuHjsRjhUvXz1XhqfGZYVtY
    ```
    Где
    - `-w .docker/wallets/wallet2.json` - путь к [кошельку](https://github.com/nspcc-dev/neo-go/blob/master/.docker/wallets/wallet2.json) для второй ноды в частной сети
    - `--in my_tx.json` - транзакция перевода, созданная на предыдущем шаге
    - `--out my_tx2.json` - выходной файл для записи подписанной транзакции
    - `--address NUVPACMnKFhpuHjsRjhUvXz1XhqfGZYVtY` - мультисиговый аккаунт для подписи транзакции
    
    Введите пароль `two`:
    ```
    Enter password to unlock wallet and sign the transaction
    Password >
    ```
    Результатом является транзакция, подписанная первой и второй нодой частной сети.

3. Подпишите транзакцию, использую адрес третьей ноды и отправьте ее в цепочку:
    ```
    $ ./bin/neo-go wallet multisig sign -w ./.docker/wallets/wallet3.json --in my_tx2.json --out my_tx3.json --address NUVPACMnKFhpuHjsRjhUvXz1XhqfGZYVtY -r http://localhost:20331
    ```
    Введите пароль `three`:
    ```
    Enter password to unlock wallet and sign the transaction
    Password >
    ```
    Результатом является транзакция, подписанная первой, второй и третьей нодами частной сети, отправленная в цепочку.

4. Проверьте баланс:

    На данный момент на балансе аккаунта `NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt` должно находиться 29999999 GAS.
    Чтобы проверить, что трансфер прошел успешно, воспользуйтесь `getnep17transfers` RPC-вызовом:
    ```
    curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getnep17transfers", "params": ["NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt"] }' localhost:20331 | json_pp
    ```
    Результат должен выглядеть следующим образом:
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


## Воркшоп. Часть 1.
Теперь все готово для того, чтобы написать, развернуть и вызовать ваш первый смарт-контракт. Начнем!


#### Шаг 1
Создайте простой смарт-контракт "Hello World" (или используйте представленный в репозтитории воркшопа):
```
package main

import (
	"github.com/nspcc-dev/neo-go/pkg/interop/runtime"
)

func Main() {
	runtime.Log("Hello, world!")
}
```
И сохраните его как `1-print.go`.

Создайте конфигурацию для смарт-контракта:
https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/1-print.yml

#### Шаг 2
Скомпилируйте смарт-контракт "Hello World":
```
$ ./bin/neo-go contract compile -i 1-print.go -c 1-print.yml -m 1-print.manifest.json
```

Где
- `./bin/neo-go` запускает neo-go
- `contract compile` команда с аргументами из [neo-go](https://github.com/nspcc-dev/neo-go/blob/master/cli/smartcontract/smart_contract.go#L105)
- `-i 1-print.go` путь к смарт-контракту
- `-c 1-print.yml` путь к конфигурационному файлу
- `-m 1-print.manifest.json` путь к файлу манифеста, который потребуется в дальнейшем при деплое смарт-контракта

Результат: 
Скомпилированный смарт-контракт `1-pring.nef` и созданный манифест смарт-контракта `1-pring.manifest.json`

Для просмотра опкодов вы можете воспользоваться командой:
```
$ ./bin/neo-go contract inspect -i 1-print.nef
```

#### Шаг 3
Разверните смарт-контракт в запущенной ранее частной сети:
```
$ ./bin/neo-go contract deploy -i 1-print.nef -manifest 1-print.manifest.json -r http://localhost:20331 -w my_wallet.json
```

Где
- `contract deploy` - команда для развертывания
- `-i 1-print.nef` - путь к смарт-контракту
- `-manifest 1-print.manifest.json` - файл манифеста смарт-контракта
- `-r http://localhost:20331` - эндпоинт ноды
- `-w my_wallet.json` - кошелек, в котором хранится ключ для подписи транзакции (вы можете взять его из репозитория воркшопа)

Введите пароль `qwerty` для аккаунта:
```
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```

Результат:
```
Contract: fb8bf3f9cab177b1c031a654fc0c94bf3db2d9b8
725857b43b5b5a67d712946b25b96a832671b21266e16be76419927d102bf9bc
```

На данном этапе ваш контракт ‘Hello World’ развернут и может быть вызван. В следующем шаге вызовем этот контракт.

#### Шаг 4
Вызовите контракт.
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json fb8bf3f9cab177b1c031a654fc0c94bf3db2d9b8 main
```

Где
- `contract invokefunction` запускает вызов контракта с заданными параметрами
- `-r http://localhost:20331` определяет эндпоинт RPC, используемый для вызова функции
- `-w my_wallet.json` - кошелек
- `fb8bf3f9cab177b1c031a654fc0c94bf3db2d9b8` хеш контракта, полученный в результате выполнения предыдущей команды (развертывание из шага 6)
- `Main` - вызываемый метод контракта

Введите пароль `qwerty` для аккаунта:
```
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```

Результат:
```
Sent invocation transaction 31e931f151716960d03188e5ba9fdbe812d7862827d77c255921f3c25a2cdc5f
```
В консоли, где была запущена нода (шаг 5), вы увидите:
```
2020-12-17T15:29:48.790+0300	INFO	runtime log	{"script": "fa0be077e0950c380ea3cb4250ae6167a75c4886", "logs": "\"Hello, world!\""}
```
Что означает, что контракт был выполнен.

На этом все. Вам потребовалось всего 4 шага, чтобы развернуть свой контракт, и это оказалось довольно легко.
Спасибо!

## Воркшоп. Часть 2
В этой части мы выполним несколько RPC вызовов и попробуем написать, задеплоить и вызвать смарт-контракт, использующий хранилище. Начнем!

### Вызовы RPC
Давайте рассмотрим более детально, что происходит с нашим смарт-контрактом при развертывании и вызове. 
Каждая нода neo-go предоставляет API интерфейс для получения данных о блокчейне.
Данное взаимодействие осуществляется по протоколу `JSON-RPC`, использующему HTTP для общения.

Полный `NEO JSON-RPC 3.0 API` описан [здесь](https://docs.neo.org/v3/docs/en-us/reference/rpc/latest-version/api.html).

RPC-сервер ноды neo-go, запущенной на шаге 5, доступен по `localhost:20331`. Давайте выполним несколько вызовов RPC.

#### GetRawTransaction
[GetRawTransaction](https://docs.neo.org/v3/docs/en-us/reference/rpc/latest-version/api/getrawtransaction.html) возвращает информацию о транзакции по ее хешу.

Запросите информацию о нашей разворачивающей транзакции из шага 3:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getrawtransaction", "params": ["725857b43b5b5a67d712946b25b96a832671b21266e16be76419927d102bf9bc", 1] }' localhost:20331 | json_pp
```

Где:
- `"jsonrpc": "2.0"` - версия протокола
- `"id": 1` - id текущего запроса
Contract: fb8bf3f9cab177b1c031a654fc0c94bf3db2d9b8
725857b43b5b5a67d712946b25b96a832671b21266e16be76419927d102bf9bc
   - `1` это `verbose` параметр для получения детального ответа в формате json-строки 

Результат:
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
[GetApplicationLog](https://docs.neo.org/v3/docs/en-us/reference/rpc/latest-version/api/getapplicationlog.html) возвращает лог контракта по соответствующему хешу транзакции.

Запросите информацию о контракте для нашей вызывающей транзакции, полученной на шаге 4:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["79a4aae4dfc2df17ee293902fdc10a5590289114635dfe551a4cd4715f4926c6"] }' localhost:20331 | json_pp
```

Где в качестве параметра:
- `79a4aae4dfc2df17ee293902fdc10a5590289114635dfe551a4cd4715f4926c6` - хеш вызывающей транзакции из шага 4

Результат:
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

#### Другие полезные вызовы RPC
```
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getversion", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getblockcount", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getconnectioncount", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getcontractstate", "params": ["0xfb8bf3f9cab177b1c031a654fc0c94bf3db2d9b8"] }' localhost:20331
```

Список всех поддерживаемых нодой neo-go вызовов RPC вы найдете [здесь](https://github.com/nspcc-dev/neo-go/blob/master/docs/rpc.md#supported-methods).

### Смарт-контракт, использующий хранилище

Давайте изучим еще один пример смарт-контракта: [2-storage.go](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/2-storage.go).
Он достаточно простой и, так же как предыдущий, не принимает никаких аргументов.
С другой стороны, этот контракт умеет считать количество его вызовов, сохраняя целое число и увеличивая его на 1 после каждого вызова.
Подобный контракт будет интересен нам, поскольку он способен *хранить* значения, т.е. обладает *хранилищем*, которое является общим для всех вызовов данного контракта.

За наличие хранилища у нашего контракта нужно заплатить дополнительное количество GAS, которое определяется вызываемым методом (например, put) и объемом данных.

В контракте `2-storage.go` также описан специальный метод `_deploy`, который выполняется во время развертывания или обновления контракта.
Данный метод не возвращает никаких значений и принимает единственный булевый аргумент, служащий индикатором обновления контракта.
Метод `_deploy` в нашем контракте предназначен для первичной инициализации счетчика вызовов контракта во время его развертывания.

Теперь, когда мы узнали о хранилище, давайте скомпилируем, развернем и вызовем смарт-контракт.

#### Шаг #1
Скомпилируйте смарт-контракт [2-storage.go](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/2-storage.go):
```
$ ./bin/neo-go contract compile -i 2-storage.go -c 2-storage.yml -m 2-storage.manifest.json
```

Результат:

Скомпилированный смарт-контракт `2-storage.nef` и манифест `2-storage.manifest.json`

#### Шаг #2
Разверните скомпилированный смарт-контракт:
```
$ ./bin/neo-go contract deploy -i 2-storage.nef -manifest 2-storage.manifest.json -r http://localhost:20331 -w my_wallet.json
```
... введите пароль `qwerty`:
```
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```

Результат:
```
Contract: 9f9bb13864f635a540fa3c3f2ce410f267a076df
a71bd9e48b354fb2421027f79b1968d6b0f4e211fc1887325e90a90201184ee1
```   

Что означает, что наш контракт развернут, и теперь мы можем вызывать его.

Давайте проверим, что значение количества вызовов контракта было проинициализировано. Используйте для этого RPC-вызов `getapplicaionlog` с хешем развертывающей транзакции в качестве параметра:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["a71bd9e48b354fb2421027f79b1968d6b0f4e211fc1887325e90a90201184ee1"] }' localhost:20331 | json_pp
```

Результат:

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
                     ... сериализованный контракт пропущен ...
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

Обратите внимание на поле `notifications`: оно содержит два уведомления `info` с сообщениями в base64.
Чтобы декодировать сообщения, испоьзуйте команду `echo string | base64 -d`:
```
$ echo U3RvcmFnZSBrZXkgbm90IHlldCBzZXQuIFNldHRpbmcgdG8gMA== | base64 -d
```
результат: `Storage key not yet set. Setting to 0`.

```
$ echo U3RvcmFnZSBrZXkgaXMgaW5pdGlhbGlzZWQ= | base64 -d
```
результат: `Storage key is initialised`.

#### Шаг #3
Поскольку мы не вызывали наш смарт-контракт раньше, при первом вызове он должен инкрементировать лежащее в хранилище значение `0` и положить новое значение = 1 в хранилище.
Давайте проверим:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json c5470e21df594ce4d143f6f9c1b4decd5e94c049 main
```
... введите пароль `qwerty`:
```
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```
Результат:
```
Sent invocation transaction 30ada0a04423eb303c93ed7bc18b5dc7b38efcf5e856a4c169f6a025e43a8bf1
```
Для проверки значения счетчика вызовем `getapplicaionlog` вызов RPC для вызывающей транзакции:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["30ada0a04423eb303c93ed7bc18b5dc7b38efcf5e856a4c169f6a025e43a8bf1"] }' localhost:20331 | json_pp
```
Результат:
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
Обратите внимание на поле `notifications`. Оно содержит сообщения, переданные методу `runtime.Notify`.
В нашем случае в нем находятся три массива байт в base64, которые можно декодировать в 3 сообщения с помощью
`echo string | base64 -d` команды CLI, например:
```
$ echo VmFsdWUgcmVhZCBmcm9tIHN0b3JhZ2U= | base64 -d
```
даст в результате:
```
Value read from storage
```
Используя эту команду, декодируем сообщения:
  - `Value read from storage`, которое было вызвано после того как мы попытались достать значение счетчика из хранилища
  - `Storage key already set. Incrementing by 1`, которое было вызвано после того, как мы поняли, что полученное значение = 0
  - `New value written into storage`, которое было вызвано после того, как мы записали новое значение в хранилище
  
И последняя часть - поле `stack`. Данное поле содержит все возвращенные контрактом значения, поэтому здесь вы можете увидеть целое `1`,
которое является значением счетчика, определяющего количество вызовов смарт-контракта.

#### Шаг #4
Для того чтобы убедиться, что все работает как надо, давайте вызовем наш контракт еще раз и проверим, что счетчик будет увеличен: 
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json c5470e21df594ce4d143f6f9c1b4decd5e94c049 main
```
... введите пароль `qwerty`:
```
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```
Результат:
```
Sent invocation transaction 23dbfd720451dd56274910365cfb982325f4312ce213c9551737ace5757f4ed4
```
Для проверки значения счетчика, выполните `getapplicaionlog` вызов RPC для вызывающей транзакции:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["23dbfd720451dd56274910365cfb982325f4312ce213c9551737ace5757f4ed4"] }' localhost:20331 | json_pp
```
Результат:
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

Теперь поле `stack` содержит значение `2` - счетчик был увеличен, как мы и ожидали.

## Воркшоп. Часть 3
В этой части мы узнаем о стандарте токена NEP17 и попробуем написать, задеплоить и вызвать более сложный смарт-контракт. Начнем!

### NEP17
[NEP17](https://github.com/neo-project/proposals/blob/master/nep-17.mediawiki) - это стандарт токена блокчейна Neo, обеспечивающий системы обобщенным механизмом взаимодействия для токенизированных смарт-контрактов.
Пример с реализацией всех требуемых стандартом методов вы можете найти в [nep17.go](https://github.com/nspcc-dev/neo-go/blob/master/examples/token/nep17/nep17.go)
 
Давайте посмотрим на пример смарт-контракта NEP17: [token.go](https://github.com/nspcc-dev/neo-go/blob/master/examples/token/token.go)
 
Этот смарт-контракт принимает в качестве параетра строку с операцией, которая может принимать следующие значения:
- `symbol` возвращает код токена
- `decimals` возвращает количество десятичных знаков токена
- `totalSupply` возвращает общий множитель * токена
- `balanceOf` возвращает баланс токена, находящегося на указанном адресе и требует дополнительног аргумента:
  - `holder` адрес запрашиваемог аккаунта
- `transfer` перемещает токен от одного пользователя к другому и требует дополнительных аргументов:
  - `from` адрес аккаунта, с которого будет списан токен
  - `to` адрес аккаунта, на который будет зачислен токен
  - `amount` количество токена для перевода
  - `data` любая дополнительная информация, которая будет передана методу `onPayment` (если получатель является контрактом)
    
Давайте проведем несколько операций с помощью этого контракта.

#### Шаг #1
Для компиляции [token.go](https://github.com/nspcc-dev/neo-go/blob/master/examples/token/token.go)
можно использовать файл [конфигурации](https://github.com/nspcc-dev/neo-go/blob/master/examples/token/token.yml).
Поскольку данный контракт использует хранилище, необходимо установить флаг 

Скомпилируйте смарт-контракт:
```
$ ./bin/neo-go contract compile -i examples/token/token.go -c examples/token/token.yml -m examples/token/token.manifest.json
```
Разверните скомпилированный контракт:
```
$ ./bin/neo-go contract deploy -i examples/token/token.nef -manifest examples/token/token.manifest.json -r http://localhost:20331 -w my_wallet.json
```
... введите пароль `qwerty`:
```
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```

Результат:
```
Contract: 1b63c282ce2618c8882694cc4214d8faee69ac42
293c9bb4977000f4bbf1aab955a79d67c9fbbfb7286f101f6cf3e34d70636e03
```   

Что означает, что наш контракт был развернут, и теперь мы можем вызывать его.

#### Шаг #2
Давайте вызовем контракт для осуществления операций с nep17.

Для начала, запросите символ созданного токена nep17:

```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 1b63c282ce2618c8882694cc4214d8faee69ac42 symbol
```                                                                   
Где
- `1b63c282ce2618c8882694cc4214d8faee69ac42` - хеш нашего контракта, полученный на шаге #1.
- `name` - строка операции, описанная ранее и возвращающая символ токена.

... не забудьте пароль от аккаунта `qwerty`.

Результат:
```
Sent invocation transaction 48f2f4f21aa92fd77247df20c83c652db77163a02dd079960e89f031142895d4
```                                                                                         
Теперь давайте подробнее посмотрим на полученную вызывающую транзакцию с помощью `getapplicationlog` RPC-вызова:

```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["48f2f4f21aa92fd77247df20c83c652db77163a02dd079960e89f031142895d4"] }' localhost:20331 | json_pp
```               

Результат:
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

Поле `stack` полученного JSON-сообщения содержит массив байтов в base64 со значением символа токена.

Следующие команды позволят получить вам дополнительную информацию о токене:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 1b63c282ce2618c8882694cc4214d8faee69ac42 decimals
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 1b63c282ce2618c8882694cc4214d8faee69ac42 totalSupply
```

#### Шаг #3
Настало время для более интересных вещей. Для начала проверим баланс nep17 токенов на нашем счету с помощью метода `balanceOf`:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 1b63c282ce2618c8882694cc4214d8faee69ac42 balanceOf NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt
```                             
... с паролем `qwerty`. Результат:
```
Sent invocation transaction 6db68a92683ffb5f6f40700d85ebb08058952d3ce0fb23ed837228118cc14025
```
Для более детального рассмотрения транзакции используем `getapplicationlog` RPC-вызов:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["6db68a92683ffb5f6f40700d85ebb08058952d3ce0fb23ed837228118cc14025"] }' localhost:20331 | json_pp
```
Результат:
```
{
   "id" : 1,
   "result" : {
      "executions" : [
         {
            "gasconsumed" : "5423900",
            "notifications" : [],
            "stack" : [
               {
                  "type" : "Integer",
                  "value" : "0"
               }
            ],
            "vmstate" : "HALT",
            "trigger" : "Application"
         }
      ],
      "txid" : "0x6db68a92683ffb5f6f40700d85ebb08058952d3ce0fb23ed837228118cc14025"
   },
   "jsonrpc" : "2.0"
}
``` 
Как вы видите, поле `stack` содержит целое значение `0`, то есть в настоящий момент мы не обладаем токенами.
Но не стоит об этом беспокоиться, переходите к следующему шагу.

#### Шаг #4

Перед тем как мы будем способны использовать наш токен (например, попытаемся передать его кому-либо), мы должны его *выпустить*.
Другими словами, мы должны перевести все имеющееся количество токена (total supply) на чей-нибудь аккаунт.
Для этого в нашем контракте существует специальная функция - `Mint`. Однако, эта функция использует сисколл `CheckWitness`, чтобы
проверить, является ли вызывающий контракта его владельцем, и обладает ли он правами управлять начальным количеством токенов.
Для этой цели существуют *подписанты* транзакции: проверка заданного хэша осуществляется с помощью листа подписантов, прикрепленного к ней.
Чтобы пройти эту проверку, нам необходимо добавить наш аккаунт с областью CalledByEntry к подписантам транзакции перевода. Давайте выпустим токен на наш адрес:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 1b63c282ce2618c8882694cc4214d8faee69ac42 mint NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt -- f8d79b436e9d3a54a38cc877d182b71e381d7a5c:CalledByEntry
```
Где
- `--` специальный разделитель, служащий для обозначения списка подписантов транзакции
- `f8d79b436e9d3a54a38cc877d182b71e381d7a5c` сам подписант транзакции (шестнадцатиричное LE представление нашего аккаунта `NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt`)

... с паролем `qwerty`. Результат:
``` 
Sent invocation transaction 5288982aa4566f0a09d7bc2fb545d74602c097cdc580626a46b4995030bc196d
```
`getapplicationlog` RPC-вызов для этой транзакции дает нам следующий результат:
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
Обратите внимание, что поле `stack` содержит значение `1` - токен был успешно выпущен.
Давайте убедимся в этом, еще раз запросив баланс нашего аккаунта с помощью метода `balanceOf`:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 1b63c282ce2618c8882694cc4214d8faee69ac42 balanceOf NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt
```
... пароль `qwerty`. Результат:
``` 
Sent invocation transaction ce8212e93ae6171efbdbb415153184ad06e4dbdd656c15ffdfcbbcd57463d7be
```
... со следующим сообщением от `getapplicationlog` вызова RPC:
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
            "gasconsumed" : "0.0514227",
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
Теперь мы видим целое значение в поле `stack`, а именно, `1100000000000000` является значением баланса токена nep17 на нашем аккаунте.

Важно, что токен может быть выпущен лишь однажды.

#### Шаг #5

После того, как мы закончили с выпуском токена, мы можем перевести некоторое количество токена кому-нибудь.
Давайте переведем 5 токенов аккаунту с адресом `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm` с помощью функции `transfer`:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 1b63c282ce2618c8882694cc4214d8faee69ac42 transfer NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm 500000000 null -- f8d79b436e9d3a54a38cc877d182b71e381d7a5c:CalledByEntry
```
... пароль `qwerty`. Результат:
``` 
Sent invocation transaction 8fefd5861d865379b6632693a5c74f430eb4465b6abf2efef8bddfd0a17dfa70
```
Наш любимый вызов RPC `getapplicationlog` говорит нам:
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
Заметьте, что поле `stack` содержит `1`, что означает, что токен был успешно переведен с нашего аккаунта.
Теперь давайте проверим баланс аккаунта, на который был совершен перевод (`NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm`), чтобы убедиться, что количество токена на нем = 5:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 1b63c282ce2618c8882694cc4214d8faee69ac42 balanceOf NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm
```
Вызов RPC `getapplicationlog` для этой транзакции возвращает следующий результат:
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
Как и ожидалось, мы видим ровно 5 токенов в поле `stack`.
Вы можете самостоятельно убедиться, что с нашего аккаунта были списаны 5 токенов, выполнив метод `balanceOf` с аргументом `NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt`.

## Воркшоп. Часть 4
В этой части подытожим наши знания о смарт-контрактах и исследуем смарт-контракт [4-domain.go](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/4-domain.go).
Данный контракт описывает операции регистрации, переноса и удаления доменов, а также операцию получения информации о зарегистрированном домене.

Начнем!

#### Шаг #1
Давайте рассмотрим и исследуем [смарт-контракт](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/4-domain.go). В качестве первого параметра контракт принимает на вход строку - действие, одно из следующих значений:
- `register` проверяет, существует ли домен с указанным именем. В случае, если такого домена не существует, добавляет пару `[domainName, owner]` в хранилище. Данная операция требудет дополнительных аргументов:
   - `domainName` - новое имя домена.
   - `owner` - 34-значный адрес аккаунта из нашего [wallet](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/my_wallet.json), который будет использоваться для вызова контракта.
- `query` возвращает адрес аккаунта владельца запрашиваемого домена (или false, в случае, если домен с указанным именем не зарегистрирован). Требует дополнительных аргументов:
   - `domainName` - имя запрашиваемого домена.
- `transfer` переводит домен с указанным именем на другой адрес (в случае, если вы являетесь владельцем указанного домена). Требует следующих аргументов:
   - `domainName` - имя домена, который вы хотите перевести.
   - `toAddress` - адрес аккаунта, на который вы хотите перевести домен.
- `delete` удаляет домен из хранилища. Аргументы:
   - `domainName` имя домента, который вы хотите удалить.
 
 
 В следующих шагах мы скомпилируем и развернем смарт-контракт.
 После этого мы зарегистрируем новый домен, переведем его на другой аккаунт и запросим информацию о нем.

#### Шаг #2

Скомпилируйте смарт-контракт [4-domain.go](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/4-domain.go) с [конфигурацией](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/4-domain.yml)
```
$ ./bin/neo-go contract compile -i 4-domain.go -c 4-domain.yml -m 4-domain.manifest.json
```

... и разверните его:
```
$ ./bin/neo-go contract deploy -i 4-domain.avm -c 4-domain.yml -e http://localhost:20331 -w my_wallet.json -g 0.001
```
Обратите внимание, что наш контракт использует хранилище и, как и с предыдущим контрактом, необходимо, чтобы флаг `hasstorage` имел значение `true`.
Этот флаг указывается в файле [конфигурации](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/4-domain.yml).

... введите пароль `qwerty`:
```
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```

Результат:
```
Contract: 66304b7c84fd61bf6ab5b9401e8fb7fe6ee8e486
91b3d71991cc6e89923fcc3f53dfc708c6dbeefa924ecebbc3a9b7c6f8f01f94
```   
Вы догадываетесь, что это значит :)

#### Шаг #3

Вызовите контракт, чтобы зарегистрировать домен с именем `my_first_domain`: 
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 66304b7c84fd61bf6ab5b9401e8fb7fe6ee8e486 register my_first_domain NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt -- f8d79b436e9d3a54a38cc877d182b71e381d7a5c:CalledByEntry
```
... пароль: `qwerty`
```
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```
Результат:
```
Sent invocation transaction 114b0c93dfb540cd86e50470dd628c3884a391568667acb6a44bd3800ff640f1
```
Также вы можете увидеть лог-сообщение в консоли, где запускали ноду neo-go:
```
2020-12-17T17:31:43.480+0300	INFO	runtime log	{"script": "86e4e86efeb78f1e40b9b56abf61fd847c4b3066", "logs": "\"RegisterDomain: my_first_domain\""}
```
Все получилось. Теперь проверим, был ли наш домен действительно зарегистрирован, с помощью вызова RPC `getapplicationlog`:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["114b0c93dfb540cd86e50470dd628c3884a391568667acb6a44bd3800ff640f1"] }' localhost:20331 | json_pp
```
Результат:
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
Здесь мы в особенности заинтересованы в двух полях полученного json:

Первое поле - `notifications`, оно содержит одно уведомление с именем `registered`:
- `bXlfZmlyc3RfZG9tYWlu` - строка в base64, которая может быть декодирована в имя нашего домена - `my_first_domain`,
- `XHodOB63gtF3yIyjVDqdbkOb1/g=` - строка, которая декодируется в адрес нашего аккаунта `NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt`.

Второе поле - `stack`, в котором лежит `1` - значение, возвращенное смарт-контрактом.

Все эти значения дают нам понять, что наш домен был успешно зарегистрирован.    

#### Шаг #4

Вызовите контракт, чтобы запросить информацию об адресе аккаунта, зарегистрировавшего домен `my_first_domain`:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 66304b7c84fd61bf6ab5b9401e8fb7fe6ee8e486 query my_first_domain
```
... любимейший пароль `qwerty`:
```
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```
Результат:
```
Sent invocation transaction 1e941ce1641fed5593a3dfb2ddc1a10090c3b466d1b5bf636f0ff1ee458a9554
```
и лог-сообщение в консоли запущенной ноды neo-go:
```
2020-12-17T17:39:32.677+0300	INFO	runtime log	{"script": "86e4e86efeb78f1e40b9b56abf61fd847c4b3066", "logs": "\"QueryDomain: my_first_domain\""}
```
Проверим транзакцию с помощью вызова RPC `getapplicationlog`:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["1e941ce1641fed5593a3dfb2ddc1a10090c3b466d1b5bf636f0ff1ee458a9554"] }' localhost:20331 | json_pp
```
... что даст нам следующий результат:
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

с base64 представлением адреса аккаунта `NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt` на стеке и в уведомлениях, что означает, что домен `my_first_domain` был зарегистрирован владельцем с полученным адресом аккаунта.

#### Шаг #5

Вызовите контракт для передачи домена другому аккаунту (например, аккаунту с адресом `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm`):
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 66304b7c84fd61bf6ab5b9401e8fb7fe6ee8e486 transfer my_first_domain NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm -- f8d79b436e9d3a54a38cc877d182b71e381d7a5c
```
... пароль: `qwerty`
```
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```
Результат:
```
Sent invocation transaction 012dc71a593b024fd38e4d301a1ded85841f1cb0d0d94c81f42846699fda8a3a
```
и лог-сообщение:
```
2020-12-17T17:44:07.536+0300	INFO	runtime log	{"script": "86e4e86efeb78f1e40b9b56abf61fd847c4b3066", "logs": "\"TransferDomain: my_first_domain\""}
```
Отлично. И `getapplicationlog` вызов RPC...
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["012dc71a593b024fd38e4d301a1ded85841f1cb0d0d94c81f42846699fda8a3a"] }' localhost:20331 | json_pp
```
... говорит нам:
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
Поле `notifications` содержит два события:
- первое с именем `deleted` и полями с дополнительной информацией (домен `my_first_domain` был удален с аккаунта с адресом `NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt`),
- второе с именем `registered` и полями с дополнительной информацией (домен `my_first_domain` был зарегистрирован аккаунтом с адресом `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm`).
Поле `stack` содержит `1`, что значит, что домен был успешно перемещен.

#### Шаг #6

Оставшийся вызов - `delete`, вы можете попробовать выполнить его самостоятельно, создав перед этим еще один домен, например, с именем `my_second_domain`, а затем удалить его из хранилища с помощью:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 66304b7c84fd61bf6ab5b9401e8fb7fe6ee8e486 register my_second_domain NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt -- f8d79b436e9d3a54a38cc877d182b71e381d7a5c:CalledByEntry
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 66304b7c84fd61bf6ab5b9401e8fb7fe6ee8e486 delete my_second_domain -- f0a33acbf435417f3d2c18445d607d1e6f48d413
```

Спасибо!

### Полезные ссылки

* [Наш воркшоп на Medium](https://medium.com/@neospcc/%D1%81%D0%BC%D0%B0%D1%80%D1%82-%D0%BA%D0%BE%D0%BD%D1%82%D1%80%D0%B0%D0%BA%D1%82-%D0%B4%D0%BB%D1%8F-neo-769139352b65)
* [Использование NEO Blockchain Toolkit](https://medium.com/@neospcc/neogo-adds-support-for-neo-blockchain-toolkit-673ea914f661)
* [Документация NEO](https://docs.neo.org/)
* [NEO github](https://github.com/neo-project/neo/)
* [NEO-GO github](https://github.com/nspcc-dev/neo-go)
