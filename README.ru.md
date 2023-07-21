<p align="center">
<img src="./pic/neo_color_dark_gopher.png" width="300px" alt="logo">
</p>

[Neo](https://neo.org/) разрабатывает системы умной экономики, и мы в [NeoSPCC](https://nspcc.ru/) помогаем им с этой нелегкой задачей. 
В нашем блоге вы можете найти статью [how we run NeoFS public test net](https://medium.com/@neospcc/public-neofs-testnet-launch-18f6315c5ced), 
но это не единственная вещь, над которой мы работаем.

## NeoGo
Как вы знаете, сеть состоит из нод. В текущий момент ноды имеют несколько реализаций:
- https://github.com/neo-project/neo.git
- https://github.com/nspcc-dev/neo-go

Данная статья посвящена последней реализации, поскольку мы в NeoSPCC занимаемся ее разработкой.
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
NAME:
   VM CLI - Official VM CLI for Neo-Go

USAGE:
    [global options] command [command options] [arguments...]

VERSION:
   0.98.2

COMMANDS:
   exit        Exit the VM prompt
   ip          Show current instruction
   break       Place a breakpoint
   estack      Show evaluation stack contents
   istack      Show invocation stack contents
   sslot       Show static slot contents
   lslot       Show local slot contents
   aslot       Show arguments slot contents
   loadnef     Load a NEF-consistent script into the VM
   loadbase64  Load a base64-encoded script string into the VM
   loadhex     Load a hex-encoded script string into the VM
   loadgo      Compile and load a Go file with the manifest into the VM
   reset       Unload compiled script from the VM
   parse       Parse provided argument and convert it into other possible formats
   run         Execute the current loaded script
   cont        Continue execution of the current loaded script
   step        Step (n) instruction in the program
   stepinto    Stepinto instruction to take in the debugger
   stepout     Stepout instruction to take in the debugger
   stepover    Stepover instruction to take in the debugger
   ops         Dump opcodes of the current loaded program
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
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
READY: loaded 21 instructions
NEO-GO-VM 0 >  
```
Теперь вы можете увидеть, сколько инструкций было сгенерировано. Также вы можете получить опкоды (opcodes) данной программы:
```
NEO-GO-VM 0 > ops
INDEX    OPCODE       PARAMETER
0        PUSHDATA1    48656c6c6f2c20776f726c6421 ("Hello, world!")    <<
15       SYSCALL      System.Runtime.Log (cfe74796)
20       RET
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
Блокчейн - достаточно большая часть NeoGo, содержащая в себе операции по принятию и валидации транзакций, их подписи,
работе с аккаунтами, ассетами, хранению блоков в базе данных (или в кэше).

#### Сеть
Существует 3 типа сетей.
Частная сеть (Private net) - это сеть, которую вы можете запустить локально. Тестовая сеть (Testnet) и Основная сеть (Mainnet) - сети, в которых запущены большинство нод Neo по всему миру.
Каждую ноду, запущенную в сети блокчейн, вы можете найти в [Neo Monitor](http://monitor.cityofzion.io/)

## Воркшоп. Содержание
Воркшоп содержит обучающее руководство, примеры и подсказки, помогающие начать
разработку смарт контрактов и децентрализованных приложений для экосистемы Neo
с помощью инструментов, предлагаемых проектом NeoGo. Воркшоп содержит несколько
частей:

1. [**Подготовка.**](#workshop-preparation) В этой части рассказывается, как запустить локальную сеть Neo, перевести
  средства на аккаунт с помощью NeoGo CLI и проверить баланс с помощью вызова
  JSON RPC.
2. [**Часть 1.**](#workshop-part-1) Содержит инструкции для компиляции, исследования, развертывания и вызовов простого
  пример контракта `Hello, world!`, написанного на Go.
3. [**Часть 2.**](#workshop-part-2) Помогает познакомиться с протоколом Neo JSON-RPC и утилитами NeoGo
  CLI для получения информации от RPC-нод Neo. В эту часть включено описание концепции хранилища смарт
  контрактов. Кроме того, в часть включена инструкция по компиляции, развертыванию и вызовам смарт контракта, который
  демонстрирует вариант использования своего хранилища.
4. [**Часть 3.**](#workshop-part-3) Содержит описание стандарта токенов NEP-17 и пример
  контракта, поддерживающего стандарт NEP-17.
5. [**Часть 4.**](#workshop-part-4) Резюмирует знания, полученные о смарт контрактах. Содержит инструкции по
  компиляции, развертыванию и вызовам более сложного смарт контракта.
6. [**Часть 5.**](#workshop-part-5) Описывает, как разработать простое децентрализованное
  приложение для экосистемы Neo, используя инструменты NeoGo.

## Воркшоп. Подготовка
В этой части мы настроим окружение: запустим частную сеть, подсоединим к ней ноду neo-go и переведем немного GAS на аккаунт, с который будем использовать далее
для создания транзакций. Давайте начнем.

#### Требования
Для этого воркшопа у вам понадобятся установленные Debian 10, Docker, docker-compose и go:
- [docker](https://docs.docker.com/install/linux/docker-ce/debian/)
- [go](https://golang.org/dl/)

#### Шаг 1
Если у вас уже установлен neo-go или есть смарт-контракты на go, пожалуйста, обновите go modules чтобы использовать свежую версию API интеропов.
Если нет, скачайте neo-go и соберите проект (ветку master):
```
$ git clone https://github.com/nspcc-dev/neo-go.git
$ cd neo-go
$ make build 
```

#### Шаг 2
Запустим локальную сеть из четырёх узлов в реализации NeoGo.

```
$ make env_image
$ make env_up
```

В результате должна запуститься приватная сеть:
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
    Enter password >
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
    Enter password >
    ```
    Результатом является транзакция, подписанная первой и второй нодой частной сети.

3. Подпишите транзакцию, использую адрес третьей ноды и отправьте ее в цепочку:
    ```
    $ ./bin/neo-go wallet sign -w ./.docker/wallets/wallet3.json --in my_tx2.json --out my_tx3.json --address NVTiAjNgagDkTr5HTzDmQP9kPwPHN5BgVq -r http://localhost:20331
    ```
    Введите пароль `three`:
    ```
    Enter password >
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
   "result" : {
      "received" : [
         {
            "txhash" : "0x7f1a2c41f0c03107f7a44ac510fa95fe11dde4c4994d30d61439f73f27e70f0d",
            "transfernotifyindex" : 0,
            "transferaddress" : "NVTiAjNgagDkTr5HTzDmQP9kPwPHN5BgVq",
            "blockindex" : 27,
            "timestamp" : 1657014304108,
            "amount" : "2999999900000000",
            "assethash" : "0xd2a4cff31913016155e38e474a2c06d08be276cf"
         }
      ],
      "address" : "NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB",
      "sent" : []
   },
   "jsonrpc" : "2.0",
   "id" : 1
}
```


## Воркшоп. Часть 1.
Теперь все готово для того, чтобы написать, развернуть и вызовать ваш первый смарт-контракт. Начнем!


#### Шаг 1
Используйте представленный в репозитории воркшопа контракт
[1-print.go](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/1-print/1-print.go).
Его код прост, это "Hello World" из смарт-контракта:
```
package main

import (
	"github.com/nspcc-dev/neo-go/pkg/interop/runtime"
)

func Main() {
	runtime.Log("Hello, world!")
}
```

Конфигурация контракта доступна в том же каталоге, [1-print.yml](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/1-print/1-print.yml).

#### Шаг 2
Скомпилируйте смарт-контракт "Hello World":
```
$ ./bin/neo-go contract compile -i 1-print/1-print.go -c 1-print/1-print.yml -m 1-print/1-print.manifest.json
```

Где
- `./bin/neo-go` запускает neo-go
- `contract compile` команда с аргументами из [neo-go](https://github.com/nspcc-dev/neo-go/blob/master/cli/smartcontract/smart_contract.go#L105)
- `-i 1-print/1-print.go` путь к смарт-контракту
- `-c 1-print/1-print.yml` путь к конфигурационному файлу
- `-m 1-print/1-print.manifest.json` путь к файлу манифеста, который потребуется в дальнейшем при деплое смарт-контракта

Результат: 
Скомпилированный смарт-контракт `1-print.nef` и созданный манифест смарт-контракта `1-print.manifest.json`.

Для просмотра опкодов вы можете воспользоваться командой:
```
$ ./bin/neo-go contract inspect -i 1-print/1-print.nef
```

#### Шаг 3
Разверните смарт-контракт в запущенной ранее частной сети:
```
$ ./bin/neo-go contract deploy -i 1-print/1-print.nef -manifest 1-print/1-print.manifest.json -r http://localhost:20331 -w my_wallet.json
```

Где
- `contract deploy` - команда для развертывания
- `-i 1-print/1-print.nef` - путь к смарт-контракту
- `-manifest 1-print/1-print.manifest.json` - файл манифеста смарт-контракта
- `-r http://localhost:20331` - эндпоинт ноды
- `-w my_wallet.json` - кошелек, в котором хранится ключ для подписи транзакции (вы можете взять его из репозитория воркшопа)

Введите пароль `qwerty` для аккаунта:
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```
CLI предложит подтвердить отправку транзакции с указанными комиссиями. Здесь и
далее введите `y` для подтверждения:
```
Network fee: 0.0151452
System fee: 10.0104553
Total fee: 10.0256005
Relay transaction (y|N)> y
```

Результат:
```
Sent invocation transaction b0436603d27d14e3aa27280e1bc2cdb17d4def8cb8cda2204c3b6a203203e6bf
Contract: bfad19135422aaddf2fc86f86ec5d4b1371e8e93
```

На данном этапе ваш контракт ‘Hello World’ развернут и может быть вызван. В следующем шаге вызовем этот контракт.

#### Шаг 4
Вызовите контракт.
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json bfad19135422aaddf2fc86f86ec5d4b1371e8e93 main
```

Где
- `contract invokefunction` запускает вызов контракта с заданными параметрами
- `-r http://localhost:20331` определяет эндпоинт RPC, используемый для вызова функции
- `-w my_wallet.json` - кошелек
- `bfad19135422aaddf2fc86f86ec5d4b1371e8e93` хеш контракта, полученный в результате выполнения предыдущей команды (развертывание из шага 6)
- `Main` - вызываемый метод контракта

Введите пароль `qwerty` для аккаунта:
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```
Результат:
```
Network fee: 0.0117652
System fee: 0.0196731
Total fee: 0.0314383
Relay transaction (y|N)> y
Sent invocation transaction 60fbf79b06714e34a3f55e782bf509eedcc602661e49848f5611aaaf2f3442de
```
В консоли, где была запущена нода (шаг 5), вы увидите:
```
2022-07-05T12:52:49.413+0300	INFO	runtime log	{"tx": "60fbf79b06714e34a3f55e782bf509eedcc602661e49848f5611aaaf2f3442de", "script": "bfad19135422aaddf2fc86f86ec5d4b1371e8e93", "msg": "Hello, world!"}
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

Полный `Neo JSON-RPC 3.0 API` описан [здесь](https://docs.neo.org/v3/docs/en-us/reference/rpc/latest-version/api.html).

RPC-сервер ноды neo-go, запущенной на шаге 5, доступен по `localhost:20331`. Давайте выполним несколько вызовов RPC.

#### GetRawTransaction
[GetRawTransaction](https://docs.neo.org/v3/docs/en-us/reference/rpc/latest-version/api/getrawtransaction.html) возвращает информацию о транзакции по ее хешу.

Запросите информацию о нашей разворачивающей транзакции из шага 3:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getrawtransaction", "params": ["b0436603d27d14e3aa27280e1bc2cdb17d4def8cb8cda2204c3b6a203203e6bf", 1] }' localhost:20331 | json_pp
```

Где:
- `"jsonrpc": "2.0"` - версия протокола
- `"id": 1` - id текущего запроса
Contract: bfad19135422aaddf2fc86f86ec5d4b1371e8e93
b0436603d27d14e3aa27280e1bc2cdb17d4def8cb8cda2204c3b6a203203e6bf
   - `1` это `verbose` параметр для получения детального ответа в формате json-строки 

Результат:
```
{
   "id" : 1,
   "result" : {
      "nonce" : 2596996162,
      "sysfee" : "1001045530",
      "size" : 531,
      "attributes" : [],
      "blocktime" : 1657014649346,
      "script" : "DOZ7Im5hbWUiOiJIZWxsb1dvcmxkIGNvbnRyYWN0IiwiYWJpIjp7Im1ldGhvZHMiOlt7Im5hbWUiOiJtYWluIiwib2Zmc2V0IjowLCJwYXJhbWV0ZXJzIjpbXSwicmV0dXJudHlwZSI6IlZvaWQiLCJzYWZlIjpmYWxzZX1dLCJldmVudHMiOltdfSwiZmVhdHVyZXMiOnt9LCJncm91cHMiOltdLCJwZXJtaXNzaW9ucyI6W10sInN1cHBvcnRlZHN0YW5kYXJkcyI6W10sInRydXN0cyI6W10sImV4dHJhIjpudWxsfQxjTkVGM25lby1nby0wLjk5LjEtcHJlLTEwMy1nM2ZiYzEzMzEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABUMDUhlbGxvLCB3b3JsZCFBz+dHlkCTKBNVEsAfDAZkZXBsb3kMFP2j+kNG6lMqJY/El92t22Q3yf3/QWJ9W1I=",
      "validuntilblock" : 51,
      "signers" : [
         {
            "account" : "0x410b5658f92f9937ed7bdd4ba04c665d3bdbd8ae",
            "scopes" : "CalledByEntry"
         }
      ],
      "netfee" : "1514520",
      "sender" : "NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB",
      "version" : 0,
      "confirmations" : 15,
      "witnesses" : [
         {
            "invocation" : "DEDAFsqnEFjXighqESMUGAAZxR2vaDBpYfbMgH55C1Q8TFNl5AfQA7Cder+MCDjPLNu7S1KHqwp97XlZK2OpZGnf",
            "verification" : "DCEDhEhWuuSSNuCc7nLsxQhI8nFlt+UfY3oP0/UkYmdH7G5BVuezJw=="
         }
      ],
      "blockhash" : "0x1fbfb494e669c03a4666fb8b2da9ed2f8205b07aabe55125311bb1f569e83d92",
      "hash" : "0xb0436603d27d14e3aa27280e1bc2cdb17d4def8cb8cda2204c3b6a203203e6bf",
      "vmstate" : "HALT"
   },
   "jsonrpc" : "2.0"
}
```

#### GetApplicationLog
[GetApplicationLog](https://docs.neo.org/v3/docs/en-us/reference/rpc/latest-version/api/getapplicationlog.html) возвращает лог контракта по соответствующему хешу транзакции.

Запросите информацию о контракте для нашей вызывающей транзакции, полученной на шаге 4:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["60fbf79b06714e34a3f55e782bf509eedcc602661e49848f5611aaaf2f3442de"] }' localhost:20331 | json_pp
```

Где в качестве параметра:
- `60fbf79b06714e34a3f55e782bf509eedcc602661e49848f5611aaaf2f3442de` - хеш вызывающей транзакции из шага 4

Результат:
```
{
   "jsonrpc" : "2.0",
   "id" : 1,
   "result" : {
      "txid" : "0x60fbf79b06714e34a3f55e782bf509eedcc602661e49848f5611aaaf2f3442de",
      "executions" : [
         {
            "gasconsumed" : "1967310",
            "trigger" : "Application",
            "stack" : [
               {
                  "type" : "Any"
               }
            ],
            "vmstate" : "HALT",
            "notifications" : []
         }
      ]
   }
}
```

#### Другие полезные вызовы RPC
```
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getversion", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getblockcount", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getconnectioncount", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getcontractstate", "params": ["0xbfad19135422aaddf2fc86f86ec5d4b1371e8e93"] }' localhost:20331
```

Список всех поддерживаемых нодой neo-go вызовов RPC вы найдете [здесь](https://github.com/nspcc-dev/neo-go/blob/master/docs/rpc.md#supported-methods).

#### Вспомогательные инструменты

neo-go CLI предоставляет утилиту `query tx` для проверки состояния транзакции. Данная утилита
ипользует вызовы RPC `getrawtransaction` и `getapplicationlog` для определения состояния транзакции
и выводит детали выполнения транзакции на экран. Используйте команду `query tx` для того, чтобы 
убедиться, что транзакция была успешно принята в блок:
```
./bin/neo-go query tx 60fbf79b06714e34a3f55e782bf509eedcc602661e49848f5611aaaf2f3442de -r http://localhost:20331 -v
```
где
- `60fbf79b06714e34a3f55e782bf509eedcc602661e49848f5611aaaf2f3442de` - хеш вызывающей транзакции из шага #7
- `-r http://localhost:20331` - адрес и порт RPC узла
- `-v` - флаг для вывода более детальной информации о транзакции (подписантов, комиссий и скрипта)

Результат выполнения данной команды:
```
Hash:			60fbf79b06714e34a3f55e782bf509eedcc602661e49848f5611aaaf2f3442de
OnChain:		true
BlockHash:		da4b4959936208b4658136e338bec1608772cdd16124675efc4862d1b576cf74
Success:		true
Signer:			NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB (None)
SystemFee:		0.0196731 GAS
NetworkFee:		0.0117652 GAS
Script:			wh8MBG1haW4MFJOOHjex1MVu+Ib88t2qIlQTGa2/QWJ9W1I=
INDEX    OPCODE       PARAMETER
0        NEWARRAY0        <<
1        PUSH15       
2        PUSHDATA1    6d61696e ("main")
8        PUSHDATA1    938e1e37b1d4c56ef886fcf2ddaa22541319adbf
30       SYSCALL      System.Contract.Call (627d5b52)
```

Поле `OnChain` говорит о том, была ли транзакция принята в блок. Поле `Success` является индикатором
того, был ли скрипт транзакции выполнен без ошибок, т.е. были ли сохранены изменения, вносимые транзакцией в чейн
и осталась ли VM в состоянии `HALT` после выполнения скрипта транзакции.


### Смарт-контракт, использующий хранилище

Давайте изучим еще один пример смарт-контракта: [2-storage.go](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/2-storage/2-storage.go).
Он достаточно простой и, так же как предыдущий, не принимает никаких аргументов.
С другой стороны, этот контракт умеет считать количество его вызовов, сохраняя целое число и увеличивая его на 1 после каждого вызова.
Подобный контракт будет интересен нам, поскольку он способен *хранить* значения, т.е. обладает *хранилищем*, которое является общим для всех вызовов данного контракта.

За наличие хранилища у нашего контракта нужно заплатить дополнительное количество GAS, которое определяется вызываемым методом (например, put) и объемом данных.

В контракте `2-storage.go` также описан специальный метод `_deploy`, который выполняется во время развертывания или обновления контракта.
Данный метод не возвращает никаких значений и принимает единственный булевый аргумент, служащий индикатором обновления контракта.
Метод `_deploy` в нашем контракте предназначен для первичной инициализации счетчика вызовов контракта во время его развертывания.

Теперь, когда мы узнали о хранилище, давайте скомпилируем, развернем и вызовем смарт-контракт.

#### Шаг #1
Скомпилируйте смарт-контракт [2-storage.go](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/2-storage/2-storage.go):
```
$ ./bin/neo-go contract compile -i 2-storage/2-storage.go -c 2-storage/2-storage.yml -m 2-storage/2-storage.manifest.json
```

Результат:

Скомпилированный смарт-контракт `2-storage.nef` и манифест `2-storage.manifest.json` в каталоге `2-storage`.

#### Шаг #2
Разверните скомпилированный смарт-контракт:
```
$ ./bin/neo-go contract deploy -i 2-storage/2-storage.nef -manifest 2-storage/2-storage.manifest.json -r http://localhost:20331 -w my_wallet.json
```
... введите пароль `qwerty`:
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```

Результат:
```
Network fee: 0.0210952
System fee: 10.0624424
Total fee: 10.0835376
Relay transaction (y|N)> y
Sent invocation transaction 58b98332ad8456da1ab6c9d162c1c557ed2ef85a67ae6bbec17e33cdad599c25
Contract: ccd533440d0317e9f366c50648d0013540e82741
```   

Что означает, что наш контракт развернут, и теперь мы можем вызывать его.

Давайте проверим, что значение количества вызовов контракта было проинициализировано. Используйте для этого RPC-вызов `getapplicaionlog` с хешем развертывающей транзакции в качестве параметра:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["58b98332ad8456da1ab6c9d162c1c557ed2ef85a67ae6bbec17e33cdad599c25"] }' localhost:20331 | json_pp
```

Результат:

```
{
   "jsonrpc" : "2.0",
   "result" : {
      "txid" : "0x58b98332ad8456da1ab6c9d162c1c557ed2ef85a67ae6bbec17e33cdad599c25",
      "executions" : [
         {
            "vmstate" : "HALT",
            "notifications" : [
               {
                  "contract" : "0xccd533440d0317e9f366c50648d0013540e82741",
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
                  "contract" : "0xccd533440d0317e9f366c50648d0013540e82741",
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
                           "value" : "QSfoQDUB0EgGxWbz6RcDDUQz1cw=",
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
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json ccd533440d0317e9f366c50648d0013540e82741 main
```
... введите пароль `qwerty`:
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```
Результат:
```
Network fee: 0.0117652
System fee: 0.0717313
Total fee: 0.0834965
Relay transaction (y|N)> y
Sent invocation transaction bf92dbe258d9113f2a0684b6c782566b5b7bcda7524b349e085beb99147d8dc1
```
Для проверки значения счетчика вызовем `getapplicaionlog` вызов RPC для вызывающей транзакции:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["bf92dbe258d9113f2a0684b6c782566b5b7bcda7524b349e085beb99147d8dc1"] }' localhost:20331 | json_pp
```
Результат:
```
{
   "id" : 1,
   "jsonrpc" : "2.0",
   "result" : {
      "txid" : "0xbf92dbe258d9113f2a0684b6c782566b5b7bcda7524b349e085beb99147d8dc1",
      "executions" : [
         {
            "trigger" : "Application",
            "gasconsumed" : "7173130",
            "notifications" : [
               {
                  "eventname" : "info",
                  "contract" : "0xccd533440d0317e9f366c50648d0013540e82741",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "value" : "VmFsdWUgcmVhZCBmcm9tIHN0b3JhZ2U=",
                           "type" : "ByteString"
                        }
                     ]
                  }
               },
               {
                  "state" : {
                     "value" : [
                        {
                           "value" : "U3RvcmFnZSBrZXkgYWxyZWFkeSBzZXQuIEluY3JlbWVudGluZyBieSAx",
                           "type" : "ByteString"
                        }
                     ],
                     "type" : "Array"
                  },
                  "eventname" : "info",
                  "contract" : "0xccd533440d0317e9f366c50648d0013540e82741"
               },
               {
                  "eventname" : "info",
                  "contract" : "0xccd533440d0317e9f366c50648d0013540e82741",
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
            "vmstate" : "HALT",
            "stack" : [
               {
                  "type" : "Integer",
                  "value" : "1"
               }
            ]
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
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json ccd533440d0317e9f366c50648d0013540e82741 main
```
... введите пароль `qwerty`:
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```
Результат:
```
Network fee: 0.0117652
System fee: 0.0717313
Total fee: 0.0834965
Relay transaction (y|N)> y
Sent invocation transaction 4edf5fcceef8fddd6c145cafaad52fee6ebe83f166d28c1ad862349182f5550d
```
Для проверки значения счетчика, выполните `getapplicaionlog` вызов RPC для вызывающей транзакции:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["4edf5fcceef8fddd6c145cafaad52fee6ebe83f166d28c1ad862349182f5550d"] }' localhost:20331 | json_pp
```
Результат:
```
{
   "id" : 1,
   "result" : {
      "txid" : "0x4edf5fcceef8fddd6c145cafaad52fee6ebe83f166d28c1ad862349182f5550d",
      "executions" : [
         {
            "vmstate" : "HALT",
            "notifications" : [
               {
                  "state" : {
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "VmFsdWUgcmVhZCBmcm9tIHN0b3JhZ2U="
                        }
                     ],
                     "type" : "Array"
                  },
                  "contract" : "0xccd533440d0317e9f366c50648d0013540e82741",
                  "eventname" : "info"
               },
               {
                  "eventname" : "info",
                  "contract" : "0xccd533440d0317e9f366c50648d0013540e82741",
                  "state" : {
                     "value" : [
                        {
                           "value" : "U3RvcmFnZSBrZXkgYWxyZWFkeSBzZXQuIEluY3JlbWVudGluZyBieSAx",
                           "type" : "ByteString"
                        }
                     ],
                     "type" : "Array"
                  }
               },
               {
                  "eventname" : "info",
                  "contract" : "0xccd533440d0317e9f366c50648d0013540e82741",
                  "state" : {
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "TmV3IHZhbHVlIHdyaXR0ZW4gaW50byBzdG9yYWdl"
                        }
                     ],
                     "type" : "Array"
                  }
               }
            ],
            "stack" : [
               {
                  "value" : "2",
                  "type" : "Integer"
               }
            ],
            "trigger" : "Application",
            "gasconsumed" : "7173130"
         }
      ]
   },
   "jsonrpc" : "2.0"
}
```

Теперь поле `stack` содержит значение `2` - счетчик был увеличен, как мы и ожидали.

## Воркшоп. Часть 3
В этой части мы узнаем о стандарте токена NEP-17 и попробуем написать, задеплоить и вызвать более сложный смарт-контракт. Начнем!

### NEP-17
[NEP-17](https://github.com/neo-project/proposals/blob/master/nep-17.mediawiki) - это стандарт токена блокчейна Neo, обеспечивающий системы обобщенным механизмом взаимодействия для токенизированных смарт-контрактов.
Пример с реализацией всех требуемых стандартом методов вы можете найти в [nep17.go](https://github.com/nspcc-dev/neo-go/blob/master/examples/token/nep17/nep17.go)
 
Давайте посмотрим на пример смарт-контракта NEP-17: [token.go](https://github.com/nspcc-dev/neo-go/blob/master/examples/token/token.go)
 
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
Network fee: 0.0308452
System fee: 10.0107577
Total fee: 10.0416029
Relay transaction (y|N)> y
Sent invocation transaction ab6f934a5e2137d008613977b41b0a791e5497c2e97a2a84aed0bb684af2c5c3
Contract: c36534b6b81621178980438c18796f23a463441a
```   

Что означает, что наш контракт был развернут, и теперь мы можем вызывать его.

#### Шаг #2
Давайте вызовем контракт для осуществления операций с NEP-17.

Для начала, запросите символ созданного токена NEP-17:

```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json c36534b6b81621178980438c18796f23a463441a symbol
```                                                                   
Где
- `c36534b6b81621178980438c18796f23a463441a` - хеш нашего контракта, полученный на шаге #1.
- `name` - строка операции, описанная ранее и возвращающая символ токена.

... не забудьте пароль от аккаунта `qwerty`.

Результат:
```
Network fee: 0.0117852
System fee: 0.0141954
Total fee: 0.0259806
Relay transaction (y|N)> y
Sent invocation transaction 1e9d018b4ea8fc3442229ae437bc8451e876216dd04a9c071672753437c9ada5
```                                                                                         
Теперь давайте подробнее посмотрим на полученную вызывающую транзакцию с помощью `getapplicationlog` RPC-вызова:

```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["1e9d018b4ea8fc3442229ae437bc8451e876216dd04a9c071672753437c9ada5"] }' localhost:20331 | json_pp
```               

Результат:
```
{
   "id" : 1,
   "result" : {
      "txid" : "0x1e9d018b4ea8fc3442229ae437bc8451e876216dd04a9c071672753437c9ada5",
      "executions" : [
         {
            "vmstate" : "HALT",
            "stack" : [
               {
                  "value" : "QU5U",
                  "type" : "ByteString"
               }
            ],
            "notifications" : [],
            "trigger" : "Application",
            "gasconsumed" : "1419540"
         }
      ]
   },
   "jsonrpc" : "2.0"
}
```

Поле `stack` полученного JSON-сообщения содержит массив байтов в base64 со значением символа токена.

Следующие команды позволят получить вам дополнительную информацию о токене:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json c36534b6b81621178980438c18796f23a463441a decimals
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json c36534b6b81621178980438c18796f23a463441a totalSupply
```

#### Шаг #3
Настало время для более интересных вещей. Для начала проверим баланс NEP-17 токенов на нашем счету с помощью метода `balanceOf`:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json c36534b6b81621178980438c18796f23a463441a balanceOf NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB
```                             
... с паролем `qwerty`. Результат:
```
Network fee: 0.0120452
System fee: 0.0249927
Total fee: 0.0370379
Relay transaction (y|N)> y
Sent invocation transaction 2f9d830cfc52747c6d0658aeb25c75e334afd1e4badd1dc946c706a22a7d1e10
```
Для более детального рассмотрения транзакции используем `getapplicationlog` RPC-вызов:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["2f9d830cfc52747c6d0658aeb25c75e334afd1e4badd1dc946c706a22a7d1e10"] }' localhost:20331 | json_pp
```
Результат:
```
{
   "id" : 1,
   "jsonrpc" : "2.0",
   "result" : {
      "txid" : "0x2f9d830cfc52747c6d0658aeb25c75e334afd1e4badd1dc946c706a22a7d1e10",
      "executions" : [
         {
            "vmstate" : "HALT",
            "stack" : [
               {
                  "type" : "Integer",
                  "value" : "0"
               }
            ],
            "notifications" : [],
            "gasconsumed" : "2499270",
            "trigger" : "Application"
         }
      ]
   }
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
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json c36534b6b81621178980438c18796f23a463441a mint NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB -- NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB:CalledByEntry
```
Где
- `--` специальный разделитель, служащий для обозначения списка подписантов транзакции
- `NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB` сам подписант транзакции (наш аккаунт)

... с паролем `qwerty`. Результат:
``` 
Network fee: 0.0119952
System fee: 0.1371123
Total fee: 0.1491075
Relay transaction (y|N)> y
Sent invocation transaction 9a54e07e54550e57ab9d7d1a1a001516ae2514bae23f6691632f2a3bc1c2d8b7
```
`getapplicationlog` RPC-вызов для этой транзакции дает нам следующий результат:
```
{
   "jsonrpc" : "2.0",
   "result" : {
      "txid" : "0x9a54e07e54550e57ab9d7d1a1a001516ae2514bae23f6691632f2a3bc1c2d8b7",
      "executions" : [
         {
            "stack" : [
               {
                  "value" : true,
                  "type" : "Boolean"
               }
            ],
            "notifications" : [
               {
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
                           "type" : "Integer",
                           "value" : "1100000000000000"
                        }
                     ]
                  },
                  "contract" : "0xc36534b6b81621178980438c18796f23a463441a",
                  "eventname" : "Transfer"
               }
            ],
            "trigger" : "Application",
            "vmstate" : "HALT",
            "gasconsumed" : "13711230"
         }
      ]
   },
   "id" : 1
}
```
Обратите внимание, что поле `stack` содержит значение `true` - токен был успешно выпущен.
Давайте убедимся в этом, еще раз запросив баланс нашего аккаунта с помощью метода `balanceOf`:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json c36534b6b81621178980438c18796f23a463441a balanceOf NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB
```
... пароль `qwerty`. Результат:
``` 
Network fee: 0.0120452
System fee: 0.0274533
Total fee: 0.0394985
Relay transaction (y|N)> y
Sent invocation transaction 870607e5bbffdaef9adb38cf4ca08125481554bf674d6a63e79c3779c924017c
```
... со следующим сообщением от `getapplicationlog` вызова RPC:
```
{
   "id" : 1,
   "result" : {
      "executions" : [
         {
            "trigger" : "Application",
            "gasconsumed" : "2745330",
            "stack" : [
               {
                  "type" : "Integer",
                  "value" : "1100000000000000"
               }
            ],
            "vmstate" : "HALT",
            "notifications" : []
         }
      ],
      "txid" : "0x870607e5bbffdaef9adb38cf4ca08125481554bf674d6a63e79c3779c924017c"
   },
   "jsonrpc" : "2.0"
}
```
Теперь мы видим целое значение в поле `stack`, а именно, `1100000000000000` является значением баланса токена NEP-17 на нашем аккаунте.

Важно, что токен может быть выпущен лишь однажды.

#### Шаг #5

После того, как мы закончили с выпуском токена, мы можем перевести некоторое количество токена кому-нибудь.
Давайте переведем 5 токенов аккаунту с адресом `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm` с помощью функции `transfer`:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json c36534b6b81621178980438c18796f23a463441a transfer NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm 500000000 null -- NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB:CalledByEntry
```
... пароль `qwerty`. Результат:
``` 
Network fee: 0.0123652
System fee: 0.1188695
Total fee: 0.1312347
Relay transaction (y|N)> y
Sent invocation transaction a796fc3d5b75f6c5289aca9d2f77d6e50c5d9fdbd860068fb5771b99ff747e96
```
Наш любимый вызов RPC `getapplicationlog` говорит нам:
```
{
   "id" : 1,
   "jsonrpc" : "2.0",
   "result" : {
      "txid" : "0xa796fc3d5b75f6c5289aca9d2f77d6e50c5d9fdbd860068fb5771b99ff747e96",
      "executions" : [
         {
            "stack" : [
               {
                  "type" : "Boolean",
                  "value" : true
               }
            ],
            "gasconsumed" : "11886950",
            "vmstate" : "HALT",
            "notifications" : [
               {
                  "eventname" : "Transfer",
                  "contract" : "0xc36534b6b81621178980438c18796f23a463441a",
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
                  }
               }
            ],
            "trigger" : "Application"
         }
      ]
   }
}
```
Заметьте, что поле `stack` содержит `true`, что означает, что токен был успешно переведен с нашего аккаунта.
Теперь давайте проверим баланс аккаунта, на который был совершен перевод (`NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm`), чтобы убедиться, что количество токена на нем = 5:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json c36534b6b81621178980438c18796f23a463441a balanceOf NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm
```
Вызов RPC `getapplicationlog` для этой транзакции возвращает следующий результат:
```
{
   "id" : 1,
   "jsonrpc" : "2.0",
   "result" : {
      "txid" : "0x40f3f3c12d3eeba7e282bbaf76af944310b82504dcb3c09db3ea6c6d8418bb6b",
      "executions" : [
         {
            "gasconsumed" : "2745330",
            "trigger" : "Application",
            "vmstate" : "HALT",
            "stack" : [
               {
                  "type" : "Integer",
                  "value" : "500000000"
               }
            ],
            "notifications" : []
         }
      ]
   }
}
```
Как и ожидалось, мы видим ровно 5 токенов в поле `stack`.
Вы можете самостоятельно убедиться, что с нашего аккаунта были списаны 5 токенов, выполнив метод `balanceOf` с аргументом `NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB`.

## Воркшоп. Часть 4
В этой части подытожим наши знания о смарт-контрактах и исследуем смарт-контракт [4-domain.go](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/4-domain/4-domain.go).
Данный контракт описывает операции регистрации, переноса и удаления доменов, а также операцию получения информации о зарегистрированном домене.

Начнем!

#### Шаг #1
Давайте рассмотрим и исследуем [смарт-контракт](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/4-domain/4-domain.go). В качестве первого параметра контракт принимает на вход строку - действие, одно из следующих значений:
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

Скомпилируйте смарт-контракт [4-domain.go](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/4-domain/4-domain.go) с [конфигурацией](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/4-domain/4-domain.yml)
```
$ ./bin/neo-go contract compile -i 4-domain/4-domain.go -c 4-domain/4-domain.yml -m 4-domain/4-domain.manifest.json
```

... и разверните его:
```
$ ./bin/neo-go contract deploy -i 4-domain/4-domain.nef -m 4-domain/4-domain.manifest.json -e http://localhost:20331 -w my_wallet.json
```

... введите пароль `qwerty`:
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```

Результат:
```
Network fee: 0.0303652
System fee: 10.0107577
Total fee: 10.0411229
Relay transaction (y|N)> y
Sent invocation transaction 69548dfecf70c190e2bc872aa210d53ff7faa7956074154c90a27d4c94420562
Contract: 9042814f07d65d2b835fa1f07d21c22c6e1cbdf7
```   
Вы догадываетесь, что это значит :)

#### Шаг #3

Вызовите контракт, чтобы зарегистрировать домен с именем `my_first_domain`: 
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 9042814f07d65d2b835fa1f07d21c22c6e1cbdf7 register my_first_domain NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB -- NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB:CalledByEntry
```
... пароль: `qwerty`
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```
Результат:
```
Network fee: 0.0122052
System fee: 0.0894353
Total fee: 0.1016405
Relay transaction (y|N)> y
Sent invocation transaction d9f05af09ee497ceb523be3274483143f83f7f24038f37799a962c8dee640357
```
Также вы можете увидеть лог-сообщение в консоли, где запускали ноду neo-go:
```
2022-07-05T13:42:36.592+0300	INFO	runtime log	{"tx": "d9f05af09ee497ceb523be3274483143f83f7f24038f37799a962c8dee640357", "script": "9042814f07d65d2b835fa1f07d21c22c6e1cbdf7", "msg": "RegisterDomain: my_first_domain"}
```
Все получилось. Теперь проверим, был ли наш домен действительно зарегистрирован, с помощью вызова RPC `getapplicationlog`:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["d9f05af09ee497ceb523be3274483143f83f7f24038f37799a962c8dee640357"] }' localhost:20331 | json_pp
```
Результат:
```
{
   "jsonrpc" : "2.0",
   "result" : {
      "executions" : [
         {
            "stack" : [
               {
                  "type" : "Boolean",
                  "value" : true
               }
            ],
            "vmstate" : "HALT",
            "notifications" : [
               {
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "value" : "rtjbO11mTKBL3XvtN5kv+VhWC0E=",
                           "type" : "ByteString"
                        },
                        {
                           "value" : "bXlfZmlyc3RfZG9tYWlu",
                           "type" : "ByteString"
                        }
                     ]
                  },
                  "contract" : "0x9042814f07d65d2b835fa1f07d21c22c6e1cbdf7",
                  "eventname" : "registered"
               }
            ],
            "trigger" : "Application",
            "gasconsumed" : "8943530"
         }
      ],
      "txid" : "0xd9f05af09ee497ceb523be3274483143f83f7f24038f37799a962c8dee640357"
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
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 9042814f07d65d2b835fa1f07d21c22c6e1cbdf7 query my_first_domain
```
... любимейший пароль `qwerty`:
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```
Результат:
```
Network fee: 0.0119552
System fee: 0.0412161
Total fee: 0.0531713
Relay transaction (y|N)> y
Sent invocation transaction 4d1b15e3c891b4e5efe6f093e54c1090476dc0ed0069b228f58578c9360ee1f2
```
и лог-сообщение в консоли запущенной ноды neo-go:
```
2022-07-05T13:44:21.681+0300	INFO	runtime log	{"tx": "4d1b15e3c891b4e5efe6f093e54c1090476dc0ed0069b228f58578c9360ee1f2", "script": "9042814f07d65d2b835fa1f07d21c22c6e1cbdf7", "msg": "QueryDomain: my_first_domain"}
```
Проверим транзакцию с помощью вызова RPC `getapplicationlog`:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["4d1b15e3c891b4e5efe6f093e54c1090476dc0ed0069b228f58578c9360ee1f2"] }' localhost:20331 | json_pp
```
... что даст нам следующий результат:
```
{
   "result" : {
      "executions" : [
         {
            "stack" : [
               {
                  "value" : "rtjbO11mTKBL3XvtN5kv+VhWC0E=",
                  "type" : "ByteString"
               }
            ],
            "trigger" : "Application",
            "gasconsumed" : "4121610",
            "notifications" : [],
            "vmstate" : "HALT"
         }
      ],
      "txid" : "0x4d1b15e3c891b4e5efe6f093e54c1090476dc0ed0069b228f58578c9360ee1f2"
   },
   "id" : 1,
   "jsonrpc" : "2.0"
}
```

с base64 представлением адреса аккаунта `NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB` на стеке и в уведомлениях, что означает, что домен `my_first_domain` был зарегистрирован владельцем с полученным адресом аккаунта.

#### Шаг #5

Вызовите контракт для передачи домена другому аккаунту (например, аккаунту с адресом `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm`):
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 9042814f07d65d2b835fa1f07d21c22c6e1cbdf7 transfer my_first_domain NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm -- NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB
```
... пароль: `qwerty`
```
Enter account NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB password >
```
Результат:
```
Network fee: 0.0122052
System fee: 0.0748064
Total fee: 0.0870116
Relay transaction (y|N)> y
Sent invocation transaction 937ec7539a31246ff88eb0bdda74cf7d9613d4a8ad1b7b33f0a785e458d76a14
```
и лог-сообщение:
```
2022-07-05T13:46:06.746+0300	INFO	runtime log	{"tx": "937ec7539a31246ff88eb0bdda74cf7d9613d4a8ad1b7b33f0a785e458d76a14", "script": "9042814f07d65d2b835fa1f07d21c22c6e1cbdf7", "msg": "TransferDomain: my_first_domain"}
```
Отлично. И `getapplicationlog` вызов RPC...
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["937ec7539a31246ff88eb0bdda74cf7d9613d4a8ad1b7b33f0a785e458d76a14"] }' localhost:20331 | json_pp
```
... говорит нам:
```
{
   "jsonrpc" : "2.0",
   "result" : {
      "executions" : [
         {
            "vmstate" : "HALT",
            "gasconsumed" : "7480640",
            "trigger" : "Application",
            "stack" : [
               {
                  "value" : true,
                  "type" : "Boolean"
               }
            ],
            "notifications" : [
               {
                  "contract" : "0x9042814f07d65d2b835fa1f07d21c22c6e1cbdf7",
                  "eventname" : "deleted",
                  "state" : {
                     "type" : "Array",
                     "value" : [
                        {
                           "type" : "ByteString",
                           "value" : "rtjbO11mTKBL3XvtN5kv+VhWC0E="
                        },
                        {
                           "type" : "ByteString",
                           "value" : "bXlfZmlyc3RfZG9tYWlu"
                        }
                     ]
                  }
               },
               {
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
                  },
                  "contract" : "0x9042814f07d65d2b835fa1f07d21c22c6e1cbdf7"
               }
            ]
         }
      ],
      "txid" : "0x937ec7539a31246ff88eb0bdda74cf7d9613d4a8ad1b7b33f0a785e458d76a14"
   },
   "id" : 1
}
```
Поле `notifications` содержит два события:
- первое с именем `deleted` и полями с дополнительной информацией (домен `my_first_domain` был удален с аккаунта с адресом `NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB`),
- второе с именем `registered` и полями с дополнительной информацией (домен `my_first_domain` был зарегистрирован аккаунтом с адресом `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm`).
Поле `stack` содержит `true`, что значит, что домен был успешно перемещен.

#### Шаг #6

Оставшийся вызов - `delete`, вы можете попробовать выполнить его самостоятельно, создав перед этим еще один домен, например, с именем `my_second_domain`, а затем удалить его из хранилища с помощью:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 9042814f07d65d2b835fa1f07d21c22c6e1cbdf7 register my_second_domain NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB -- NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB:CalledByEntry
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 9042814f07d65d2b835fa1f07d21c22c6e1cbdf7 delete my_second_domain -- NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB
```

## Воркшоп. Часть 5

В этой части будет представлен пример децентрализованного приложения на языке Go,
работающего в экосистеме Neo. Пример показывает варианты использования NeoGo RPC
клиента и набора утилит для работы с кошельками, деплоя и вызва смарт контрактов,
получения данных из блокчейна, обработки происходящих в блокчейне событий и
автоматической генерации RPC биндингов для смарт контрактов.

#### Шаг #1

Убедитесь, что локальная сеть из четырех узлов, описанная в [подготовительной части](#workshop-preparation)
настроена и запущена. Нам также понадобится аккаунт `NbrUYaZgyhSkNoRo9ugRyEMdUZxrhkNaWB`
из [кошелька](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/e11949bf5f1dc1ce4e3b6551b6ae22032945c75d/my_wallet.json) с некоторым количеством токена GAS на нём. Откройте и внимательно
просмотрите пример децентрализованного приложения: [dApp.go](./dApp/dApp.go). Пример включает
в себя различные варианты использования базового API, которые необходимы для
начала разработки своего собственного децентрализванного приложения в экосистеме
Neo:
* Создание [клиента JSON-RPC](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient), его инициализация и примеры вызовов базовых
  RPC методов.
* Управление кошельками и аккаунтами Neo с помощью [пакета NeoGo `wallet`](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/wallet).
* [Расширения](https://github.com/nspcc-dev/neo-go/blob/5fc61be5f6c5349d8de8b61967380feee6b51c55/docs/rpc.md#extensions), предлагаемые NeoGo JSON RPC сервером и поддерживаемые NeoGo RPC клиентом.
* [Веб-сокетное расширение](https://github.com/nspcc-dev/neo-go/blob/5fc61be5f6c5349d8de8b61967380feee6b51c55/docs/rpc.md#websocket-server)NeoGo JSON RPC сервера. Создание, инициализация и использование [веб-сокетного JSON RPC клиента](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient#WSClient). 
* [Пакетная утилита `unwrap`](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/unwrap), предназначенная для извлечения результатов тестовых RPC-вызовов.
* Примеры использования [подсистемы уведомлений](https://github.com/nspcc-dev/neo-go/blob/5fc61be5f6c5349d8de8b61967380feee6b51c55/docs/notifications.md) NeoGo JSON RPC. Подписки на происходящие в блокчейне события и получение уведомлений по веб-сокетному каналу.
* Примеры использования пакета NeoGo `invoker`. Выполнение тестовых вызовов смарт контрактов, скриптов или верификационных методов с удобным интерфейсом [Invoker](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/invoker).
* Статья-введение в Neo witnesses и [области их действия](https://neospcc.medium.com/thou-shalt-check-their-witnesses-485d2bf8375d).
* [Исторические вызовы](https://github.com/nspcc-dev/neo-go/blob/5fc61be5f6c5349d8de8b61967380feee6b51c55/docs/rpc.md#invokecontractverifyhistoric-invokefunctionhistoric-and-invokescripthistoric-calls), являющиеся еще одним расширением NeoGo RPC сервера и клиентский API для их использования.
* Примеры использования пакета NeoGo `actor`. Генерация, настройка, подпись, отправка и осуществление ожидания испонения транзаций с гибким интерфейсом [Actor](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/actor).
* Работа с набором NEP-специфичных и контракт-специфичных пакетов:
  * Использование пакетов [`nep17`](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/nep17) и [`nep11`](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/nep11) для вызовов NEP-17 и NEP-11 совместимых контрактов.
  * Использование акторов, специфичных для нативных контрактов: [`gas`](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/gas), [`neo`](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/neo), [`management`](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/management), [`policy`](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/policy), [`oracle`](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/oracle), [`rolemgmt`](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/rolemgmt), [`notary`](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/notary) пакеты.
* Деплой и вызовы [примера контракта Storage](https://github.com/nspcc-dev/neo-go/tree/5fc61be5f6c5349d8de8b61967380feee6b51c55/examples/storage) с помощью пакета-актора, специфичного для [нативного `ContractManagement` контракта](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/management).
* Демонстрация API итераторов контрактного хранилища:
  * Обход итератора в NeoVM-скрипте с использованием [`CallAndExpandIterator` Invoker API](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/invoker#Invoker.CallAndExpandIterator).
  * Обход сессионного итератора с использованием [`TraverseIterator` Invoker API](https://pkg.go.dev/github.com/nspcc-dev/neo-go/pkg/rpcclient/invoker#Invoker.TraverseIterator).
* Использование автогенерированных [контрактных RPC биндингов](https://github.com/nspcc-dev/neo-go/blob/5fc61be5f6c5349d8de8b61967380feee6b51c55/docs/compiler.md#generating-rpc-contract-bindings).

Внимательно разберите пример децентрализованного приложения и прочитайте
комментарии к блокам кода в нём. Перейдите по ссылкам, представленным в
комментариях и прочитайте документацию к использованным API. Перед тем как
запустить приложение склонируйте репозиторий NeoGo в ту же папку, где расположен
репозиторий воркшопа и скомпилируйте демонстрационный пример контракта Storage
из репозитория NeoGo:
```
$ git clone https://github.com/nspcc-dev/neo-go
$ cd neo-go
$ make build
$ ./bin/neo-go contract compile -i ./examples/storage/storage.go -c ./examples/storage/storage.yml -o examples/storage/storage.nef -m ./examples/storage/storage.manifest.json
```  

В качестве последнего подготовительного шага проверьте, что хэш транзакции,
переводящей GAS с мультисигового аккаунта на аккаунт из `my_wallet.json`, верно
указан в переменной `transferTxH` примера `./dApp/dApp.go`

#### Шаг #2

В консоли перейдите в директорию приложения `dApp` и запустите приложение:
```
$ cd ./dApp
$ go run dApp.go 
``` 

Изучите вывод приложения в консоль, просмотрите примеры кода, представленного в
приложении и используйте их для разработки своего собственного децентрализованного
приложения для экосистемы Neo!

Спасибо!

### Полезные ссылки

* [Наш воркшоп на Medium](https://medium.com/@neospcc/%D1%81%D0%BC%D0%B0%D1%80%D1%82-%D0%BA%D0%BE%D0%BD%D1%82%D1%80%D0%B0%D0%BA%D1%82-%D0%B4%D0%BB%D1%8F-neo-769139352b65)
* [Использование Neo Blockchain Toolkit](https://medium.com/@neospcc/neogo-adds-support-for-neo-blockchain-toolkit-673ea914f661)
* [Документация Neo](https://docs.neo.org/)
* [Neo github](https://github.com/neo-project/neo/)
* [NeoGo github](https://github.com/nspcc-dev/neo-go)
