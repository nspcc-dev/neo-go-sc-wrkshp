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
        $ ./bin/neo-go wallet nep17 transfer -w .docker/wallets/wallet1.json --out my_tx.json -r http://localhost:20331 --from NVTiAjNgagDkTr5HTzDmQP9kPwPHN5BgVq --to NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB --token GAS --amount 29999999
    ``` 
    Где
    - `./bin/neo-go` запускает neo-go
    - `wallet nep17 transfer` - команда с аргументами в [neo-go](https://github.com/nspcc-dev/neo-go/blob/master/cli/wallet/nep17.go#L108)
    - `-w .docker/wallets/wallet1.json` - путь к [кошельку](https://github.com/nspcc-dev/neo-go/blob/master/.docker/wallets/wallet1.json) для первой ноды в созданной частной сети
    - `--out my_tx.json` - файл для записи подписанной транзакции
    - `-r http://localhost:20331` - RPC-эндпоинт ноды
    - `--from NVTiAjNgagDkTr5HTzDmQP9kPwPHN5BgVq` - мультисиговый аккаунт, являющийся владельцем GAS
    - `--to NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB` - наш аккаунт из [кошелька](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/my_wallet.json)
    - `--token GAS` - имя переводимого токена (в данном случае это GAS)
    - `--amount 29999999` - количество GAS для перевода
    
    Введите пароль `one`:
    ```
    Password >
    ```
    Результатом является транзакция, подписанная первой нодой и сохраненная в `my_tx.json`.

2. Подпишите созданную транзакцию, используя адрес второй ноды:

    ```
    $ ./bin/neo-go wallet sign -w .docker/wallets/wallet2.json --in my_tx.json --out my_tx2.json --address NVTiAjNgagDkTr5HTzDmQP9kPwPHN5BgVq
    ```
    Где
    - `-w .docker/wallets/wallet2.json` - путь к [кошельку](https://github.com/nspcc-dev/neo-go/blob/master/.docker/wallets/wallet2.json) для второй ноды в частной сети
    - `--in my_tx.json` - транзакция перевода, созданная на предыдущем шаге
    - `--out my_tx2.json` - выходной файл для записи подписанной транзакции
    - `--address NVTiAjNgagDkTr5HTzDmQP9kPwPHN5BgVq` - мультисиговый аккаунт для подписи транзакции
    
    Введите пароль `two`:
    ```
    Enter password to unlock wallet and sign the transaction
    Password >
    ```
    Результатом является транзакция, подписанная первой и второй нодой частной сети.

3. Подпишите транзакцию, использую адрес третьей ноды и отправьте ее в цепочку:
    ```
    $ ./bin/neo-go wallet sign -w ./.docker/wallets/wallet3.json --in my_tx2.json --out my_tx3.json --address NVTiAjNgagDkTr5HTzDmQP9kPwPHN5BgVq -r http://localhost:20331
    ```
    Введите пароль `three`:
    ```
    Enter password to unlock wallet and sign the transaction
    Password >
    ```
    Результатом является транзакция, подписанная первой, второй и третьей нодами частной сети, отправленная в цепочку.

4. Проверьте баланс:

    На данный момент на балансе аккаунта `NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB` должно находиться 29999999 GAS.
    Чтобы проверить, что трансфер прошел успешно, воспользуйтесь `getnep17transfers` RPC-вызовом:
    ```
    curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getnep17transfers", "params": ["NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB"] }' localhost:20331 | json_pp
    ```
    Результат должен выглядеть следующим образом:
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
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```
CLI предложит подтвердить отправку транзакции с указанными комиссиямию Здесь и
далее введите `y` для подтверждения:
```
Network fee: 1515520
System fee: 1001045530
Total fee: 1002561050
Relay transaction (y|N)>
```

Результат:
```
Sent invocation transaction 28b26283ea2689dc5abf30bf6f0605b3819089f7fbf07bc26e41d62e1a9f5841
Contract: a48467c9bf559524575cf0d3b25cd97e67b01bc5
```

На данном этапе ваш контракт ‘Hello World’ развернут и может быть вызван. В следующем шаге вызовем этот контракт.

#### Шаг 4
Вызовите контракт.
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json a48467c9bf559524575cf0d3b25cd97e67b01bc5 main
```

Где
- `contract invokefunction` запускает вызов контракта с заданными параметрами
- `-r http://localhost:20331` определяет эндпоинт RPC, используемый для вызова функции
- `-w my_wallet.json` - кошелек
- `a48467c9bf559524575cf0d3b25cd97e67b01bc5` хеш контракта, полученный в результате выполнения предыдущей команды (развертывание из шага 6)
- `Main` - вызываемый метод контракта

Введите пароль `qwerty` для аккаунта:
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```
Результат:
```
Sent invocation transaction bd23c836f7bdd62a0d9c5ecb3f5bdbf2d38ec9a5e2e3935ca543d8c18ed5479d
```
В консоли, где была запущена нода (шаг 5), вы увидите:
```
2021-11-29T17:02:44.395+0300	INFO	runtime log	{"tx": "bd23c836f7bdd62a0d9c5ecb3f5bdbf2d38ec9a5e2e3935ca543d8c18ed5479d", "script": "a48467c9bf559524575cf0d3b25cd97e67b01bc5", "msg": "Hello, world!"}
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
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getrawtransaction", "params": ["28b26283ea2689dc5abf30bf6f0605b3819089f7fbf07bc26e41d62e1a9f5841", 1] }' localhost:20331 | json_pp
```

Где:
- `"jsonrpc": "2.0"` - версия протокола
- `"id": 1` - id текущего запроса
Contract: a48467c9bf559524575cf0d3b25cd97e67b01bc5
28b26283ea2689dc5abf30bf6f0605b3819089f7fbf07bc26e41d62e1a9f5841
   - `1` это `verbose` параметр для получения детального ответа в формате json-строки 

Результат:
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
[GetApplicationLog](https://docs.neo.org/v3/docs/en-us/reference/rpc/latest-version/api/getapplicationlog.html) возвращает лог контракта по соответствующему хешу транзакции.

Запросите информацию о контракте для нашей вызывающей транзакции, полученной на шаге 4:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["bd23c836f7bdd62a0d9c5ecb3f5bdbf2d38ec9a5e2e3935ca543d8c18ed5479d"] }' localhost:20331 | json_pp
```

Где в качестве параметра:
- `bd23c836f7bdd62a0d9c5ecb3f5bdbf2d38ec9a5e2e3935ca543d8c18ed5479d` - хеш вызывающей транзакции из шага 4

Результат:
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

#### Другие полезные вызовы RPC
```
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getversion", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getblockcount", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getconnectioncount", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getcontractstate", "params": ["0xa48467c9bf559524575cf0d3b25cd97e67b01bc5"] }' localhost:20331
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
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```

Результат:
```
Sent invocation transaction a0722f99edb590b789cee2589e74a09f93f36efeb06b8f5da7abde85c789a2d3
Contract: aa9c0d6006eccb53ee76688722898617606a88aa
```   

Что означает, что наш контракт развернут, и теперь мы можем вызывать его.

Давайте проверим, что значение количества вызовов контракта было проинициализировано. Используйте для этого RPC-вызов `getapplicaionlog` с хешем развертывающей транзакции в качестве параметра:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["a0722f99edb590b789cee2589e74a09f93f36efeb06b8f5da7abde85c789a2d3"] }' localhost:20331 | json_pp
```

Результат:

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
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json aa9c0d6006eccb53ee76688722898617606a88aa main
```
... введите пароль `qwerty`:
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```
Результат:
```
Sent invocation transaction a58140ee3ebee1f4fb844311b73ac86454d458122eec9c4cea19725a106a260f
```
Для проверки значения счетчика вызовем `getapplicaionlog` вызов RPC для вызывающей транзакции:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["a58140ee3ebee1f4fb844311b73ac86454d458122eec9c4cea19725a106a260f"] }' localhost:20331 | json_pp
```
Результат:
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
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json aa9c0d6006eccb53ee76688722898617606a88aa main
```
... введите пароль `qwerty`:
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```
Результат:
```
Sent invocation transaction 157ca5e5b8cf8f84c9660502a3270b346011612bded1514a6847f877c433a9bb
```
Для проверки значения счетчика, выполните `getapplicaionlog` вызов RPC для вызывающей транзакции:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["157ca5e5b8cf8f84c9660502a3270b346011612bded1514a6847f877c433a9bb"] }' localhost:20331 | json_pp
```
Результат:
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
  - `account` адрес запрашиваемого аккаунта
- `transfer` переводит токен от одного пользователя к другому и требует дополнительных аргументов:
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
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```

Результат:
```
Sent invocation transaction 7af616aacc798760274a449700f14e4e25d5c3b262d200303dc701f8ea41707c
Contract: 27502a01e2fb013e1e4c428abb7b360df9f3f0cb
```   

Что означает, что наш контракт был развернут, и теперь мы можем вызывать его.

#### Шаг #2
Давайте вызовем контракт для осуществления операций с nep17.

Для начала, запросите символ созданного токена nep17:

```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 27502a01e2fb013e1e4c428abb7b360df9f3f0cb symbol
```                                                                   
Где
- `27502a01e2fb013e1e4c428abb7b360df9f3f0cb` - хеш нашего контракта, полученный на шаге #1.
- `name` - строка операции, описанная ранее и возвращающая символ токена.

... не забудьте пароль от аккаунта `qwerty`.

Результат:
```
Sent invocation transaction 535bccc585698c531cc58677b116ea7c567604194bf3202c6be7ac4d420b85af
```                                                                                         
Теперь давайте подробнее посмотрим на полученную вызывающую транзакцию с помощью `getapplicationlog` RPC-вызова:

```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["535bccc585698c531cc58677b116ea7c567604194bf3202c6be7ac4d420b85af"] }' localhost:20331 | json_pp
```               

Результат:
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

Поле `stack` полученного JSON-сообщения содержит массив байтов в base64 со значением символа токена.

Следующие команды позволят получить вам дополнительную информацию о токене:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 27502a01e2fb013e1e4c428abb7b360df9f3f0cb decimals
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 27502a01e2fb013e1e4c428abb7b360df9f3f0cb totalSupply
```

#### Шаг #3
Настало время для более интересных вещей. Для начала проверим баланс nep17 токенов на нашем счету с помощью метода `balanceOf`:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 27502a01e2fb013e1e4c428abb7b360df9f3f0cb balanceOf NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB
```                             
... с паролем `qwerty`. Результат:
```
Sent invocation transaction e15871ad735a216a5e55f86dbb31ed4b4e928f4531f2788e547cc881e8532a8a
```
Для более детального рассмотрения транзакции используем `getapplicationlog` RPC-вызов:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["e15871ad735a216a5e55f86dbb31ed4b4e928f4531f2788e547cc881e8532a8a"] }' localhost:20331 | json_pp
```
Результат:
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
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 27502a01e2fb013e1e4c428abb7b360df9f3f0cb mint NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB -- NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB:CalledByEntry
```
Где
- `--` специальный разделитель, служащий для обозначения списка подписантов транзакции
- `NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB` сам подписант транзакции (наш аккаунт)

... с паролем `qwerty`. Результат:
``` 
Sent invocation transaction 296cb753f9afeef7ace3690eed05c08336329200f86ff82b63a67726bac5ec4c
```
`getapplicationlog` RPC-вызов для этой транзакции дает нам следующий результат:
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
Обратите внимание, что поле `stack` содержит значение `true` - токен был успешно выпущен.
Давайте убедимся в этом, еще раз запросив баланс нашего аккаунта с помощью метода `balanceOf`:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 27502a01e2fb013e1e4c428abb7b360df9f3f0cb balanceOf NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB
```
... пароль `qwerty`. Результат:
``` 
Sent invocation transaction 82ed056d7f9d27c5366561eb897d08a382747cd54b98e5c05fa82c30818f363b
```
... со следующим сообщением от `getapplicationlog` вызова RPC:
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
Теперь мы видим целое значение в поле `stack`, а именно, `1100000000000000` является значением баланса токена nep17 на нашем аккаунте.

Важно, что токен может быть выпущен лишь однажды.

#### Шаг #5

После того, как мы закончили с выпуском токена, мы можем перевести некоторое количество токена кому-нибудь.
Давайте переведем 5 токенов аккаунту с адресом `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm` с помощью функции `transfer`:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 27502a01e2fb013e1e4c428abb7b360df9f3f0cb transfer NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm 500000000 null -- NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB:CalledByEntry
```
... пароль `qwerty`. Результат:
``` 
Sent invocation transaction a8dac2052450664697f02e852b82485225f7b3a1d1017eda2b4362fbc0cc962d
```
Наш любимый вызов RPC `getapplicationlog` говорит нам:
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
Заметьте, что поле `stack` содержит `true`, что означает, что токен был успешно переведен с нашего аккаунта.
Теперь давайте проверим баланс аккаунта, на который был совершен перевод (`NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm`), чтобы убедиться, что количество токена на нем = 5:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 27502a01e2fb013e1e4c428abb7b360df9f3f0cb balanceOf NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm
```
Вызов RPC `getapplicationlog` для этой транзакции возвращает следующий результат:
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
Как и ожидалось, мы видим ровно 5 токенов в поле `stack`.
Вы можете самостоятельно убедиться, что с нашего аккаунта были списаны 5 токенов, выполнив метод `balanceOf` с аргументом `NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB`.

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
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```

Результат:
```
Sent invocation transaction 1306887c24fef841cbcc3dee3dbea734a0084c5f698ca62244bfda8f0dec4aba
Contract: a4ded8036fd90cf75daeefa7828498b80eee3e97
```   
Вы догадываетесь, что это значит :)

#### Шаг #3

Вызовите контракт, чтобы зарегистрировать домен с именем `my_first_domain`: 
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json a4ded8036fd90cf75daeefa7828498b80eee3e97 register my_first_domain NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB -- NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB:CalledByEntry
```
... пароль: `qwerty`
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```
Результат:
```
Sent invocation transaction aca960d485f63fd0feca3fa4e5217f38350a7cfec0528a16cbff1aa67279ab34
```
Также вы можете увидеть лог-сообщение в консоли, где запускали ноду neo-go:
```
2021-11-29T17:23:45.356+0300	INFO	runtime log	{"tx": "aca960d485f63fd0feca3fa4e5217f38350a7cfec0528a16cbff1aa67279ab34", "script": "a4ded8036fd90cf75daeefa7828498b80eee3e97", "msg": "RegisterDomain: my_first_domain"}
```
Все получилось. Теперь проверим, был ли наш домен действительно зарегистрирован, с помощью вызова RPC `getapplicationlog`:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["aca960d485f63fd0feca3fa4e5217f38350a7cfec0528a16cbff1aa67279ab34"] }' localhost:20331 | json_pp
```
Результат:
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
Здесь мы в особенности заинтересованы в двух полях полученного json:

Первое поле - `notifications`, оно содержит одно уведомление с именем `registered`:
- `bXlfZmlyc3RfZG9tYWlu` - строка в base64, которая может быть декодирована в имя нашего домена - `my_first_domain`,
- `ecv/0NH0e0cStm0wWBgjCxMyaok=` - строка, которая декодируется в адрес нашего аккаунта `NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB`.

Второе поле - `stack`, в котором лежит `true` - значение, возвращенное смарт-контрактом.

Все эти значения дают нам понять, что наш домен был успешно зарегистрирован.    

#### Шаг #4

Вызовите контракт, чтобы запросить информацию об адресе аккаунта, зарегистрировавшего домен `my_first_domain`:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json a4ded8036fd90cf75daeefa7828498b80eee3e97 query my_first_domain
```
... любимейший пароль `qwerty`:
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```
Результат:
```
Sent invocation transaction ddec59cd213a43e70f489e6e68ef76924f8a26538cd12b169d72ce78555c9d3a
```
и лог-сообщение в консоли запущенной ноды neo-go:
```
2021-11-29T17:26:30.476+0300	INFO	runtime log	{"tx": "ddec59cd213a43e70f489e6e68ef76924f8a26538cd12b169d72ce78555c9d3a", "script": "a4ded8036fd90cf75daeefa7828498b80eee3e97", "msg": "QueryDomain: my_first_domain"}
```
Проверим транзакцию с помощью вызова RPC `getapplicationlog`:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["ddec59cd213a43e70f489e6e68ef76924f8a26538cd12b169d72ce78555c9d3a"] }' localhost:20331 | json_pp
```
... что даст нам следующий результат:
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

с base64 представлением адреса аккаунта `NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB` на стеке и в уведомлениях, что означает, что домен `my_first_domain` был зарегистрирован владельцем с полученным адресом аккаунта.

#### Шаг #5

Вызовите контракт для передачи домена другому аккаунту (например, аккаунту с адресом `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm`):
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json a4ded8036fd90cf75daeefa7828498b80eee3e97 transfer my_first_domain NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm -- NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB
```
... пароль: `qwerty`
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```
Результат:
```
Sent invocation transaction 108d62cefd64e3adea47025bf23e3749b604aa05422e515aecaaaaa3d0b6c9a3
```
и лог-сообщение:
```
2021-11-29T17:28:00.535+0300	INFO	runtime log	{"tx": "108d62cefd64e3adea47025bf23e3749b604aa05422e515aecaaaaa3d0b6c9a3", "script": "a4ded8036fd90cf75daeefa7828498b80eee3e97", "msg": "TransferDomain: my_first_domain"}
```
Отлично. И `getapplicationlog` вызов RPC...
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["108d62cefd64e3adea47025bf23e3749b604aa05422e515aecaaaaa3d0b6c9a3"] }' localhost:20331 | json_pp
```
... говорит нам:
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
Поле `notifications` содержит два события:
- первое с именем `deleted` и полями с дополнительной информацией (домен `my_first_domain` был удален с аккаунта с адресом `NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB`),
- второе с именем `registered` и полями с дополнительной информацией (домен `my_first_domain` был зарегистрирован аккаунтом с адресом `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm`).
Поле `stack` содержит `true`, что значит, что домен был успешно перемещен.

#### Шаг #6

Оставшийся вызов - `delete`, вы можете попробовать выполнить его самостоятельно, создав перед этим еще один домен, например, с именем `my_second_domain`, а затем удалить его из хранилища с помощью:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json a4ded8036fd90cf75daeefa7828498b80eee3e97 register my_second_domain NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB -- NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB:CalledByEntry
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json a4ded8036fd90cf75daeefa7828498b80eee3e97 delete my_second_domain -- NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB
```

Спасибо!

### Полезные ссылки

* [Наш воркшоп на Medium](https://medium.com/@neospcc/%D1%81%D0%BC%D0%B0%D1%80%D1%82-%D0%BA%D0%BE%D0%BD%D1%82%D1%80%D0%B0%D0%BA%D1%82-%D0%B4%D0%BB%D1%8F-neo-769139352b65)
* [Использование NEO Blockchain Toolkit](https://medium.com/@neospcc/neogo-adds-support-for-neo-blockchain-toolkit-673ea914f661)
* [Документация NEO](https://docs.neo.org/)
* [NEO github](https://github.com/neo-project/neo/)
* [NEO-GO github](https://github.com/nspcc-dev/neo-go)
