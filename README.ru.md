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
   "id" : 1,
   "result" : {
      "sent" : [],
      "address" : "NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB",
      "received" : [
         {
            "timestamp" : 1616423429953,
            "txhash" : "0x1b123a0f26fdc22c94752a29edd7a669c96284c57523ddcc875f1862ce678c1d",
            "blockindex" : 3,
            "amount" : "2999999900000000",
            "assethash" : "0xd2a4cff31913016155e38e474a2c06d08be276cf",
            "transfernotifyindex" : 0,
            "transferaddress" : "NVTiAjNgagDkTr5HTzDmQP9kPwPHN5BgVq"
         }
      ]
   },
   "jsonrpc" : "2.0"
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

Результат:
```
Contract: ecdd946811bcfe48feefb91c927234a6f18e341c
b164f03a5dfc61273f7ebaf8943ff49a3ee1971babd25c7a817d46a2374f624a
```

На данном этапе ваш контракт ‘Hello World’ развернут и может быть вызван. В следующем шаге вызовем этот контракт.

#### Шаг 4
Вызовите контракт.
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json ecdd946811bcfe48feefb91c927234a6f18e341c main
```

Где
- `contract invokefunction` запускает вызов контракта с заданными параметрами
- `-r http://localhost:20331` определяет эндпоинт RPC, используемый для вызова функции
- `-w my_wallet.json` - кошелек
- `ecdd946811bcfe48feefb91c927234a6f18e341c` хеш контракта, полученный в результате выполнения предыдущей команды (развертывание из шага 6)
- `Main` - вызываемый метод контракта

Введите пароль `qwerty` для аккаунта:
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```

Результат:
```
Sent invocation transaction 31e931f151716960d03188e5ba9fdbe812d7862827d77c255921f3c25a2cdc5f
```
В консоли, где была запущена нода (шаг 5), вы увидите:
```
2020-12-17T15:29:48.790+0300	INFO	runtime log	{"tx": "bfb0398f22ae15628a1353c3b84afba6ff994e48cd376b840825314abf9bc291", "script": "ecdd946811bcfe48feefb91c927234a6f18e341c", "msg": "Hello, world!"}
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
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getrawtransaction", "params": ["b164f03a5dfc61273f7ebaf8943ff49a3ee1971babd25c7a817d46a2374f624a", 1] }' localhost:20331 | json_pp
```

Где:
- `"jsonrpc": "2.0"` - версия протокола
- `"id": 1` - id текущего запроса
Contract: ecdd946811bcfe48feefb91c927234a6f18e341c
b164f03a5dfc61273f7ebaf8943ff49a3ee1971babd25c7a817d46a2374f624a
   - `1` это `verbose` параметр для получения детального ответа в формате json-строки 

Результат:
```
{
   "jsonrpc" : "2.0",
   "result" : {
      "confirmations" : 11,
      "netfee" : "1546520",
      "script" : "DQQBeyJuYW1lIjoiSGVsbG9Xb3JsZCBjb250cmFjdCIsImFiaSI6eyJtZXRob2RzIjpbeyJuYW1lIjoibWFpbiIsIm9mZnNldCI6MCwicGFyYW1ldGVycyI6W10sInJldHVybnR5cGUiOiJWb2lkIiwic2FmZSI6ZmFsc2V9XSwiZXZlbnRzIjpbXX0sImZlYXR1cmVzIjp7fSwiZ3JvdXBzIjpbXSwicGVybWlzc2lvbnMiOlt7ImNvbnRyYWN0IjoiKiIsIm1ldGhvZHMiOiIqIn1dLCJzdXBwb3J0ZWRzdGFuZGFyZHMiOltdLCJ0cnVzdHMiOltdLCJleHRyYSI6bnVsbH0MZE5FRjNuZW8tZ28tMC45NS4xLXByZQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAWDA1IZWxsbywgd29ybGQhQc/nR5YhQLbNunASwB8MBmRlcGxveQwU/aP6Q0bqUyolj8SX3a3bZDfJ/f9BYn1bUg==",
      "sysfee" : "1001060650",
      "blocktime" : 1621345578680,
      "witnesses" : [
         {
            "invocation" : "DECFNf91SCJf0Xj5MqwQY9lEFMKYiwkm/wHxvx8B/1gT0TkfH2eL8sSHM4b99QklcRAUgNfniPYacMtMaOCfrTlv",
            "verification" : "DCEDhEhWuuSSNuCc7nLsxQhI8nFlt+UfY3oP0/UkYmdH7G5BVuezJw=="
         }
      ],
      "attributes" : [],
      "vmstate" : "HALT",
      "hash" : "0xb164f03a5dfc61273f7ebaf8943ff49a3ee1971babd25c7a817d46a2374f624a",
      "nonce" : 1906296755,
      "validuntilblock" : 13,
      "blockhash" : "0x75e8dd246c40806b49502471d2d6244fcf8ae881216989119a3e32f0ddeb6959",
      "sender" : "NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB",
      "signers" : [
         {
            "scopes" : "None",
            "account" : "0x410b5658f92f9937ed7bdd4ba04c665d3bdbd8ae"
         }
      ],
      "size" : 563,
      "version" : 0
   },
   "id" : 1
}
```

#### GetApplicationLog
[GetApplicationLog](https://docs.neo.org/v3/docs/en-us/reference/rpc/latest-version/api/getapplicationlog.html) возвращает лог контракта по соответствующему хешу транзакции.

Запросите информацию о контракте для нашей вызывающей транзакции, полученной на шаге 4:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["bfb0398f22ae15628a1353c3b84afba6ff994e48cd376b840825314abf9bc291"] }' localhost:20331 | json_pp
```

Где в качестве параметра:
- `bfb0398f22ae15628a1353c3b84afba6ff994e48cd376b840825314abf9bc291` - хеш вызывающей транзакции из шага 4

Результат:
```
{
   "jsonrpc" : "2.0",
   "result" : {
      "executions" : [
         {
            "trigger" : "Application",
            "notifications" : [],
            "gasconsumed" : "2028330",
            "stack" : [
               {
                  "type" : "Any"
               }
            ],
            "vmstate" : "HALT"
         }
      ],
      "txid" : "0xbfb0398f22ae15628a1353c3b84afba6ff994e48cd376b840825314abf9bc291"
   },
   "id" : 1
}
```

#### Другие полезные вызовы RPC
```
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getversion", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getblockcount", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getconnectioncount", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getcontractstate", "params": ["0xecdd946811bcfe48feefb91c927234a6f18e341c"] }' localhost:20331
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
Contract: 1b2fb1dc5b32934abae1ad1706b0e43513c44e66
1ea2a1ed6c2f651436e9fcf41023119f934cd0a4f38cac16ac9042c124345f0c
```   

Что означает, что наш контракт развернут, и теперь мы можем вызывать его.

Давайте проверим, что значение количества вызовов контракта было проинициализировано. Используйте для этого RPC-вызов `getapplicaionlog` с хешем развертывающей транзакции в качестве параметра:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["1ea2a1ed6c2f651436e9fcf41023119f934cd0a4f38cac16ac9042c124345f0c"] }' localhost:20331 | json_pp
```

Результат:

```
{
   "result" : {
      "txid" : "0x1ea2a1ed6c2f651436e9fcf41023119f934cd0a4f38cac16ac9042c124345f0c",
      "executions" : [
         {
            "gasconsumed" : "1006244000",
            "notifications" : [
               {
                  "eventname" : "info",
                  "state" : {
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "U3RvcmFnZSBrZXkgbm90IHlldCBzZXQuIFNldHRpbmcgdG8gMA=="
                        }
                     ],
                     "type" : "Array"
                  },
                  "contract" : "0x1b2fb1dc5b32934abae1ad1706b0e43513c44e66"
               },
               {
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "U3RvcmFnZSBrZXkgaXMgaW5pdGlhbGlzZWQ="
                        }
                     ]
                  },
                  "eventname" : "info",
                  "contract" : "0x1b2fb1dc5b32934abae1ad1706b0e43513c44e66"
               },
               {
                  "contract" : "0xfffdc93764dbaddd97c48f252a53ea4643faa3fd",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "Zk7EEzXksAYXreG6SpMyW9yxLxs="
                        }
                     ]
                  },
                  "eventname" : "Deploy"
               }
            ],
            "vmstate" : "HALT",
            "stack" : [
            ... skipped serialized contract representation ...
            ],
            "trigger" : "Application"
         }
      ]
   },
   "id" : 1,
   "jsonrpc" : "2.0"
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
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 1b2fb1dc5b32934abae1ad1706b0e43513c44e66 main
```
... введите пароль `qwerty`:
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```
Результат:
```
Sent invocation transaction ecc5aecf2334b06e0b5d76494bfdf3ddbaca675858f254aa4e4633756b86a40d
```
Для проверки значения счетчика вызовем `getapplicaionlog` вызов RPC для вызывающей транзакции:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["ecc5aecf2334b06e0b5d76494bfdf3ddbaca675858f254aa4e4633756b86a40d"] }' localhost:20331 | json_pp
```
Результат:
```
{
   "result" : {
      "txid" : "0xecc5aecf2334b06e0b5d76494bfdf3ddbaca675858f254aa4e4633756b86a40d",
      "executions" : [
         {
            "notifications" : [
               {
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "VmFsdWUgcmVhZCBmcm9tIHN0b3JhZ2U="
                        }
                     ]
                  },
                  "contract" : "0x1b2fb1dc5b32934abae1ad1706b0e43513c44e66",
                  "eventname" : "info"
               },
               {
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "U3RvcmFnZSBrZXkgYWxyZWFkeSBzZXQuIEluY3JlbWVudGluZyBieSAx"
                        }
                     ]
                  },
                  "contract" : "0x1b2fb1dc5b32934abae1ad1706b0e43513c44e66",
                  "eventname" : "info"
               },
               {
                  "contract" : "0x1b2fb1dc5b32934abae1ad1706b0e43513c44e66",
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
            "gasconsumed" : "7233580",
            "vmstate" : "HALT",
            "trigger" : "Application"
         }
      ]
   },
   "id" : 1,
   "jsonrpc" : "2.0"
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
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 1b2fb1dc5b32934abae1ad1706b0e43513c44e66 main
```
... введите пароль `qwerty`:
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```
Результат:
```
Sent invocation transaction 593a7c887d9d216adb46b322be05b3a0ba4d6be8450f478e3d71c6304189328c
```
Для проверки значения счетчика, выполните `getapplicaionlog` вызов RPC для вызывающей транзакции:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["593a7c887d9d216adb46b322be05b3a0ba4d6be8450f478e3d71c6304189328c"] }' localhost:20331 | json_pp
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
            "stack" : [
               {
                  "value" : "2",
                  "type" : "Integer"
               }
            ],
            "notifications" : [
               {
                  "contract" : "0x1b2fb1dc5b32934abae1ad1706b0e43513c44e66",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "VmFsdWUgcmVhZCBmcm9tIHN0b3JhZ2U="
                        }
                     ]
                  },
                  "eventname" : "info"
               },
               {
                  "eventname" : "info",
                  "state" : {
                     "value" : [
                        {
                           "value" : "U3RvcmFnZSBrZXkgYWxyZWFkeSBzZXQuIEluY3JlbWVudGluZyBieSAx",
                           "type" : "ByteString"
                        }
                     ],
                     "type" : "Array"
                  },
                  "contract" : "0x1b2fb1dc5b32934abae1ad1706b0e43513c44e66"
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
                  "contract" : "0x1b2fb1dc5b32934abae1ad1706b0e43513c44e66",
                  "eventname" : "info"
               }
            ],
            "gasconsumed" : "7233580",
            "trigger" : "Application"
         }
      ],
      "txid" : "0x593a7c887d9d216adb46b322be05b3a0ba4d6be8450f478e3d71c6304189328c"
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
Contract: 13175d9c27074057cf4d8c50183ce9d4dceaf95c
e0d85465a2fcfbb5280f068dc22979259d59d0f4d9a871a6de96fc0b92eaa3a5
```   

Что означает, что наш контракт был развернут, и теперь мы можем вызывать его.

#### Шаг #2
Давайте вызовем контракт для осуществления операций с nep17.

Для начала, запросите символ созданного токена nep17:

```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 13175d9c27074057cf4d8c50183ce9d4dceaf95c symbol
```                                                                   
Где
- `13175d9c27074057cf4d8c50183ce9d4dceaf95c` - хеш нашего контракта, полученный на шаге #1.
- `name` - строка операции, описанная ранее и возвращающая символ токена.

... не забудьте пароль от аккаунта `qwerty`.

Результат:
```
Sent invocation transaction bd85005d02c383c400595fcefb237a9f0a0919d16f7e0dbe72336c592f1951b6
```                                                                                         
Теперь давайте подробнее посмотрим на полученную вызывающую транзакцию с помощью `getapplicationlog` RPC-вызова:

```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["bd85005d02c383c400595fcefb237a9f0a0919d16f7e0dbe72336c592f1951b6"] }' localhost:20331 | json_pp
```               

Результат:
```
{
   "result" : {
      "executions" : [
         {
            "notifications" : [],
            "vmstate" : "HALT",
            "trigger" : "Application",
            "stack" : [
               {
                  "value" : "QU5U",
                  "type" : "ByteString"
               }
            ],
            "gasconsumed" : "4294290"
         }
      ],
      "txid" : "0xbd85005d02c383c400595fcefb237a9f0a0919d16f7e0dbe72336c592f1951b6"
   },
   "jsonrpc" : "2.0",
   "id" : 1
}
```

Поле `stack` полученного JSON-сообщения содержит массив байтов в base64 со значением символа токена.

Следующие команды позволят получить вам дополнительную информацию о токене:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 13175d9c27074057cf4d8c50183ce9d4dceaf95c decimals
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 13175d9c27074057cf4d8c50183ce9d4dceaf95c totalSupply
```

#### Шаг #3
Настало время для более интересных вещей. Для начала проверим баланс nep17 токенов на нашем счету с помощью метода `balanceOf`:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 13175d9c27074057cf4d8c50183ce9d4dceaf95c balanceOf NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB
```                             
... с паролем `qwerty`. Результат:
```
Sent invocation transaction 84e819e08a05a0709aa5cfc02e613b7aaa4650f324cd37eda4e2941605354498
```
Для более детального рассмотрения транзакции используем `getapplicationlog` RPC-вызов:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["84e819e08a05a0709aa5cfc02e613b7aaa4650f324cd37eda4e2941605354498"] }' localhost:20331 | json_pp
```
Результат:
```
{
   "result" : {
      "executions" : [
         {
            "stack" : [
               {
                  "value" : "0",
                  "type" : "Integer"
               }
            ],
            "notifications" : [],
            "trigger" : "Application",
            "vmstate" : "HALT",
            "gasconsumed" : "5311140"
         }
      ],
      "txid" : "0x84e819e08a05a0709aa5cfc02e613b7aaa4650f324cd37eda4e2941605354498"
   },
   "jsonrpc" : "2.0",
   "id" : 1
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
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 13175d9c27074057cf4d8c50183ce9d4dceaf95c mint NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB -- NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB:CalledByEntry
```
Где
- `--` специальный разделитель, служащий для обозначения списка подписантов транзакции
- `NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB` сам подписант транзакции (наш аккаунт)

... с паролем `qwerty`. Результат:
``` 
Sent invocation transaction 17ceda88c215876c45e2e824cd62639c113f2be235278649024124d8f0213da3
```
`getapplicationlog` RPC-вызов для этой транзакции дает нам следующий результат:
```
{
   "result" : {
      "txid" : "0x17ceda88c215876c45e2e824cd62639c113f2be235278649024124d8f0213da3",
      "executions" : [
         {
            "trigger" : "Application",
            "gasconsumed" : "16522950",
            "vmstate" : "HALT",
            "notifications" : [
               {
                  "contract" : "0x13175d9c27074057cf4d8c50183ce9d4dceaf95c",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "type" : "Any"
                        },
                        {
                           "value" : "rtjbO11mTKBL3XvtN5kv+VhWC0E=",
                           "type" : "ByteString"
                        },
                        {
                           "value" : "1100000000000000",
                           "type" : "Integer"
                        }
                     ]
                  },
                  "eventname" : "Transfer"
               }
            ],
            "stack" : [
               {
                  "value" : true,
                  "type" : "Boolean"
               }
            ]
         }
      ]
   },
   "id" : 1,
   "jsonrpc" : "2.0"
}
```
Обратите внимание, что поле `stack` содержит значение `true` - токен был успешно выпущен.
Давайте убедимся в этом, еще раз запросив баланс нашего аккаунта с помощью метода `balanceOf`:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 13175d9c27074057cf4d8c50183ce9d4dceaf95c balanceOf NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB
```
... пароль `qwerty`. Результат:
``` 
Sent invocation transaction 95fed6e224071b09e885994aa869bee9b3dd815e89cc11b871b4432a35d04b0a
```
... со следующим сообщением от `getapplicationlog` вызова RPC:
```
{
   "result" : {
      "executions" : [
         {
            "stack" : [
               {
                  "value" : "1100000000000000",
                  "type" : "Integer"
               }
            ],
            "trigger" : "Application",
            "vmstate" : "HALT",
            "gasconsumed" : "5557020",
            "notifications" : []
         }
      ],
      "txid" : "0x95fed6e224071b09e885994aa869bee9b3dd815e89cc11b871b4432a35d04b0a"
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
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 13175d9c27074057cf4d8c50183ce9d4dceaf95c transfer NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm 500000000 null -- NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB:CalledByEntry
```
... пароль `qwerty`. Результат:
``` 
Sent invocation transaction 5b8c23f4816b67e67e869c4305942fb81041671e70bed6a3888e7493883db8bd
```
Наш любимый вызов RPC `getapplicationlog` говорит нам:
```
{
   "result" : {
      "executions" : [
         {
            "vmstate" : "HALT",
            "notifications" : [
               {
                  "contract" : "0x13175d9c27074057cf4d8c50183ce9d4dceaf95c",
                  "state" : {
                     "type" : "Array",
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
                     ]
                  },
                  "eventname" : "Transfer"
               }
            ],
            "trigger" : "Application",
            "gasconsumed" : "14760950",
            "stack" : [
               {
                  "type" : "Boolean",
                  "value" : true
               }
            ]
         }
      ],
      "txid" : "0x5b8c23f4816b67e67e869c4305942fb81041671e70bed6a3888e7493883db8bd"
   },
   "jsonrpc" : "2.0",
   "id" : 1
}
```
Заметьте, что поле `stack` содержит `true`, что означает, что токен был успешно переведен с нашего аккаунта.
Теперь давайте проверим баланс аккаунта, на который был совершен перевод (`NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm`), чтобы убедиться, что количество токена на нем = 5:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 13175d9c27074057cf4d8c50183ce9d4dceaf95c balanceOf NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm
```
Вызов RPC `getapplicationlog` для этой транзакции возвращает следующий результат:
```
{
   "id" : 1,
   "jsonrpc" : "2.0",
   "result" : {
      "executions" : [
         {
            "trigger" : "Application",
            "stack" : [
               {
                  "value" : "500000000",
                  "type" : "Integer"
               }
            ],
            "gasconsumed" : "5557020",
            "notifications" : [],
            "vmstate" : "HALT"
         }
      ],
      "txid" : "0x022ade3087cb32fed650b1ab6a799ac9861539308a27f54ba6f38ffde26bd424"
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
Contract: 5e6d360e472798ee10c676022761bc5a2c3828f5
a84ea858045b498b73319d8acdf0f95c9bdba4227b0999006176c0136306b80f
```   
Вы догадываетесь, что это значит :)

#### Шаг #3

Вызовите контракт, чтобы зарегистрировать домен с именем `my_first_domain`: 
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 5e6d360e472798ee10c676022761bc5a2c3828f5 register my_first_domain NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB -- NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB:CalledByEntry
```
... пароль: `qwerty`
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```
Результат:
```
Sent invocation transaction a7f868e07898592de9846549b4395190dfe7efbbdd521178c2b72d96c25dd831
```
Также вы можете увидеть лог-сообщение в консоли, где запускали ноду neo-go:
```
2020-12-17T17:31:43.480+0300	INFO	runtime log	{"tx": "a7f868e07898592de9846549b4395190dfe7efbbdd521178c2b72d96c25dd831", "script": "5e6d360e472798ee10c676022761bc5a2c3828f5", "msg": "RegisterDomain: my_first_domain"}
```
Все получилось. Теперь проверим, был ли наш домен действительно зарегистрирован, с помощью вызова RPC `getapplicationlog`:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["a7f868e07898592de9846549b4395190dfe7efbbdd521178c2b72d96c25dd831"] }' localhost:20331 | json_pp
```
Результат:
```
{
   "jsonrpc" : "2.0",
   "result" : {
      "txid" : "0xa7f868e07898592de9846549b4395190dfe7efbbdd521178c2b72d96c25dd831",
      "executions" : [
         {
            "stack" : [
               {
                  "type" : "Boolean",
                  "value" : true
               }
            ],
            "gasconsumed" : "9143210",
            "trigger" : "Application",
            "notifications" : [
               {
                  "contract" : "0x5e6d360e472798ee10c676022761bc5a2c3828f5",
                  "eventname" : "registered",
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
                  }
               }
            ],
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
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 5e6d360e472798ee10c676022761bc5a2c3828f5 query my_first_domain
```
... любимейший пароль `qwerty`:
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```
Результат:
```
Sent invocation transaction b8b1e8e473437badfad13355c65da6bcd1e868bd30fbedfcd8f9fd5daef52bfb
```
и лог-сообщение в консоли запущенной ноды neo-go:
```
2020-12-17T17:39:32.677+0300	INFO	runtime log	{"tx": "7693ddedee55e1ccf2914a049d3cf5c3d1b29d5fb8ecbd3e9bd5672a242170a6", "script": "5e6d360e472798ee10c676022761bc5a2c3828f5", "msg": "QueryDomain: my_first_domain"}
```
Проверим транзакцию с помощью вызова RPC `getapplicationlog`:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["b8b1e8e473437badfad13355c65da6bcd1e868bd30fbedfcd8f9fd5daef52bfb"] }' localhost:20331 | json_pp
```
... что даст нам следующий результат:
```
{
   "jsonrpc" : "2.0",
   "id" : 1,
   "result" : {
      "txid" : "0xb8b1e8e473437badfad13355c65da6bcd1e868bd30fbedfcd8f9fd5daef52bfb",
      "executions" : [
         {
            "vmstate" : "HALT",
            "stack" : [
               {
                  "type" : "ByteString",
                  "value" : "rtjbO11mTKBL3XvtN5kv+VhWC0E="
               }
            ],
            "gasconsumed" : "4321230",
            "trigger" : "Application",
            "notifications" : []
         }
      ]
   }
}
```

с base64 представлением адреса аккаунта `NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB` на стеке и в уведомлениях, что означает, что домен `my_first_domain` был зарегистрирован владельцем с полученным адресом аккаунта.

#### Шаг #5

Вызовите контракт для передачи домена другому аккаунту (например, аккаунту с адресом `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm`):
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 5e6d360e472798ee10c676022761bc5a2c3828f5 transfer my_first_domain NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm -- NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB
```
... пароль: `qwerty`
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```
Результат:
```
Sent invocation transaction c1a97f8263cfea75ea27b31b53a4a2088d4a0f48d03c9967f92aa96fbf74a603
```
и лог-сообщение:
```
2020-12-17T17:44:07.536+0300	INFO	runtime log		{"tx": "5d06117399a2ffc35d1f4a16fc16de7ab5136406b16dbfb8cda033d31cb59bcc", "script": "5e6d360e472798ee10c676022761bc5a2c3828f5", "msg": "TransferDomain: my_first_domain"}
```
Отлично. И `getapplicationlog` вызов RPC...
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["c1a97f8263cfea75ea27b31b53a4a2088d4a0f48d03c9967f92aa96fbf74a603"] }' localhost:20331 | json_pp
```
... говорит нам:
```
{
   "id" : 1,
   "jsonrpc" : "2.0",
   "result" : {
      "executions" : [
         {
            "stack" : [
               {
                  "value" : true,
                  "type" : "Boolean"
               }
            ],
            "gasconsumed" : "7680110",
            "vmstate" : "HALT",
            "trigger" : "Application",
            "notifications" : [
               {
                  "eventname" : "deleted",
                  "state" : {
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "rtjbO11mTKBL3XvtN5kv+VhWC0E="
                        },
                        {
                           "type" : "ByteString",
                           "value" : "bXlfZmlyc3RfZG9tYWlu"
                        }
                     ],
                     "type" : "Array"
                  },
                  "contract" : "0x5e6d360e472798ee10c676022761bc5a2c3828f5"
               },
               {
                  "state" : {
                     "value" : [
                        {
                           "value" : "50l6vFaauRKm8hPVkr3Aw2CeHQs=",
                           "type" : "ByteString"
                        },
                        {
                           "type" : "ByteString",
                           "value" : "bXlfZmlyc3RfZG9tYWlu"
                        }
                     ],
                     "type" : "Array"
                  },
                  "contract" : "0x5e6d360e472798ee10c676022761bc5a2c3828f5",
                  "eventname" : "registered"
               }
            ]
         }
      ],
      "txid" : "0xc1a97f8263cfea75ea27b31b53a4a2088d4a0f48d03c9967f92aa96fbf74a603"
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
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 5e6d360e472798ee10c676022761bc5a2c3828f5 register my_second_domain NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB -- NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB:CalledByEntry
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 5e6d360e472798ee10c676022761bc5a2c3828f5 delete my_second_domain -- NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB
```

Спасибо!

### Полезные ссылки

* [Наш воркшоп на Medium](https://medium.com/@neospcc/%D1%81%D0%BC%D0%B0%D1%80%D1%82-%D0%BA%D0%BE%D0%BD%D1%82%D1%80%D0%B0%D0%BA%D1%82-%D0%B4%D0%BB%D1%8F-neo-769139352b65)
* [Использование NEO Blockchain Toolkit](https://medium.com/@neospcc/neogo-adds-support-for-neo-blockchain-toolkit-673ea914f661)
* [Документация NEO](https://docs.neo.org/)
* [NEO github](https://github.com/neo-project/neo/)
* [NEO-GO github](https://github.com/nspcc-dev/neo-go)
