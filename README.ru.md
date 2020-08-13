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
	"github.com/nspcc-dev/neo-go/pkg/interop/runtime"
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
Переведите немного GAS с мультисигового аккаунта на аккаунт, который мы будем использовать в дальнейшем.

1. Создадим транзакцию перевода GAS токенов:
    ```
        $ ./bin/neo-go wallet nep5 transfer -w .docker/wallets/wallet1.json --out my_tx.json -r http://localhost:20331 --from NUVPACMnKFhpuHjsRjhUvXz1XhqfGZYVtY --to NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt --token gas --amount 29999999
    ``` 
    Где
    - `./bin/neo-go` запускает neo-go
    - `wallet nep5 transfer` - команда с аргументами в [neo-go](https://github.com/nspcc-dev/neo-go/blob/master/cli/wallet/nep5.go#L108)
    - `-w .docker/wallets/wallet1.json` - путь к [кошельку](https://github.com/nspcc-dev/neo-go/blob/master/.docker/wallets/wallet1.json) для первой ноды в созданной частной сети
    - `--out my_tx.json` - файл для записи подписанной транзакции
    - `-r http://localhost:20331` - RPC-эндпоинт ноды
    - `--from NUVPACMnKFhpuHjsRjhUvXz1XhqfGZYVtY` - мультисиговый аккаунт, являющийся владельцем GAS
    - `--to NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt` - наш аккаунт из [кошелька](https://github.com/nspcc-dev/neo-go-sc-wrkshp/blob/master/my_wallet.json)
    - `--token gas` - имя переводимого токена (в данном случае это GAS)
    - `--amount 29999999` - количество GAS для перевода
    
    Введите пароль `one`:
    ```
    Password >
    ```
    Результатом является транзакция, подписанная первой нодой и сохраненная в `my_tx.json`.

2. Подпишите созданную транзакцию, используя адрес второй ноды:

    ```
    $ ./bin/neo-go wallet multisig sign -w .docker/wallets/wallet2.json --in my_tx.json --out my_tx2.json --addr NUVPACMnKFhpuHjsRjhUvXz1XhqfGZYVtY
    ```
    Где
    - `-w .docker/wallets/wallet2.json` - путь к [кошельку](https://github.com/nspcc-dev/neo-go/blob/master/.docker/wallets/wallet2.json) для второй ноды в частной сети
    - `--in my_tx.json` - транзакция перевода, созданная на предыдущем шаге
    - `--out my_tx2.json` - выходной файл для записи подписанной транзакции
    - `--addr NUVPACMnKFhpuHjsRjhUvXz1XhqfGZYVtY` - мультисиговый аккаунт для подписи транзакции
    
    Введите пароль `two`:
    ```
    Enter password to unlock wallet and sign the transaction
    Password >
    ```
    Результатом является транзакция, подписанная первой и второй нодой частной сети.

3. Подпишите транзакцию, использую адрес третьей ноды и отправьте ее в цепочку:
    ```
    $ ./bin/neo-go wallet multisig sign -w ./.docker/wallets/wallet3.json --in my_tx2.json --out my_tx3.json --addr NUVPACMnKFhpuHjsRjhUvXz1XhqfGZYVtY -r http://localhost:20331
    ```
    Введите пароль `three`:
    ```
    Enter password to unlock wallet and sign the transaction
    Password >
    ```
    Результатом является транзакция, подписанная первой, второй и третьей нодами частной сети, отправленная в цепочку.

4. Проверьте баланс:

    На данный момент на балансе аккаунта `NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt` должно находиться 29999999 GAS.
    Чтобы проверить, что трансфер прошел успешно, воспользуйтесь `getnep5transfers` RPC-вызовом:
    ```
    curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getnep5transfers", "params": ["NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt"] }' localhost:20331 | json_pp
    ```
    Результат должен выглядеть следующим образом:
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
Sent deployment transaction 8fc1dee25649d4292ac2cb19e7be2e65288ef5fab0add933b3f54c31477852dd for contract 28dbf93dc07a3d9b84ce6499132b874844784f9c
```

На данном этапе ваш контракт ‘Hello World’ развернут и может быть вызван. В следующем шаге вызовем этот контракт.

#### Шаг 4
Вызовите контракт.
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 28dbf93dc07a3d9b84ce6499132b874844784f9c main
```

Где
- `contract invokefunction` запускает вызов контракта с заданными параметрами
- `-r http://localhost:20331` определяет эндпоинт RPC, используемый для вызова функции
- `-w my_wallet.json` - кошелек
- `28dbf93dc07a3d9b84ce6499132b874844784f9c` хеш контракта, полученный в результате выполнения предыдущей команды (развертывание из шага 6)
- `Main` - вызываемый метод контракта

Введите пароль `qwerty` для аккаунта:
```
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```

Результат:
```
Sent invocation transaction 5025daa84c6c41b672444e5f26fc125811900309c84295b35b2db3d70bab18b8
```
В консоли, где была запущена нода (шаг 5), вы увидите:
```
2020-07-02T17:01:46.175+0300	INFO	runtime log	{"script": "9c4f784448872b139964ce849b3d7ac03df9db28", "logs": "\"Hello, world!\""}
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
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getrawtransaction", "params": ["8fc1dee25649d4292ac2cb19e7be2e65288ef5fab0add933b3f54c31477852dd", 1] }' localhost:20331 | json_pp
```

Где:
- `"jsonrpc": "2.0"` - версия протокола
- `"id": 1` - id текущего запроса
- `"method": "getrawtransaction"` - запрашиваемый метод
- `"params": ["8fc1dee25649d4292ac2cb19e7be2e65288ef5fab0add933b3f54c31477852dd", 1]` массив параметров запроса, где
   - `8fc1dee25649d4292ac2cb19e7be2e65288ef5fab0add933b3f54c31477852dd` - хеш разворачивающеей транзакции, полученный после выполнения шага 6
   - `1` это `verbose` параметр для получения детального ответа в формате json-строки 

Результат:
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
[GetApplicationLog](https://docs.neo.org/v3/docs/en-us/reference/rpc/latest-version/api/getapplicationlog.html) возвращает лог контракта по соответствующему хешу транзакции.

Запросите информацию о контракте для нашей вызывающей транзакции, полученной на шаге 4:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["5025daa84c6c41b672444e5f26fc125811900309c84295b35b2db3d70bab18b8"] }' localhost:20331 | json_pp
```

Где в качестве параметра:
- `5025daa84c6c41b672444e5f26fc125811900309c84295b35b2db3d70bab18b8` - хеш вызывающей транзакции из шага 4

Результат:
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

#### Другие полезные вызовы RPC
```
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getversion", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getblockcount", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getconnectioncount", "params": [] }' localhost:20331
curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getcontractstate", "params": ["0x9c33bbf2f5afbbc8fe271dd37508acd93573cffc"] }' localhost:20331
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
Sent deployment transaction e3e61ca9d713bddda110a002504f9d5f94f871a15675c2e514d08efb1521ea1b for contract 79edf20dc8ee247981756787a638e9679026c16c
```   

Что означает, что наш контракт развернут, и теперь мы можем вызывать его.

#### Шаг #3
Поскольку мы не вызывали наш смарт-контракт раньше, в его хранилище нет никаких значений, поэтому при первом вызове он должен создать новое значение (равное `1`) и положить его в хранилище.
Давайте проверим:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 79edf20dc8ee247981756787a638e9679026c16c main
```
... введите пароль `qwerty`:
```
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```
Результат:
```
Sent invocation transaction 7f658d0008f7c9bb3f05d655450ead36b9e5d7a2ad0f9e5f7895a43579d2616d
```
Для проверки значения счетчика вызовем `getapplicaionlog` вызов RPC для вызывающей транзакции:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["7f658d0008f7c9bb3f05d655450ead36b9e5d7a2ad0f9e5f7895a43579d2616d"] }' localhost:20331 | json_pp
```
Результат:
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
  - `Storage key not yet set. Setting to 1`, которое было вызвано после того, как мы поняли, что полученное значение = 0
  - `New value written into storage`, которое было вызвано после того, как мы записали новое значение в хранилище
  
И последняя часть - поле `stack`. Данное поле содержит все возвращенные контрактом значения, поэтому здесь вы можете увидеть целое `1`,
которое является значением счетчика, определяющего количество вызовов смарт-контракта.

#### Шаг #4
Для того чтобы убедиться, что все работает как надо, давайте вызовем наш контракт еще раз и проверим, что счетчик будет увеличен: 
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 79edf20dc8ee247981756787a638e9679026c16c main
```
... введите пароль `qwerty`:
```
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```
Результат:
```
Sent invocation transaction bbe338b20228e7067c71da2137d14ac7cbd01158469979d882982aee7f41e316
```
Для проверки значения счетчика, выполните `getapplicaionlog` вызов RPC для вызывающей транзакции:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["bbe338b20228e7067c71da2137d14ac7cbd01158469979d882982aee7f41e316"] }' localhost:20331 | json_pp
```
Результат:
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
Для компиляции [token.go](https://github.com/nspcc-dev/neo-go/blob/master/examples/token/token.go)
можно использовать файл [конфигурации](https://github.com/nspcc-dev/neo-go/blob/master/examples/token/token.yml).
Поскольку данный контракт использует хранилище, необходимо установить флаг 
```
hasstorage: true
```
Скомпилируйте смарт-контракт:
```
$ ./bin/neo-go contract compile -i examples/token/token.go -c 3-token.yml -m examples/token/token.manifest.json
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
Sent deployment transaction d275df152996bdab279efb0942b15fd26c6a2b69965cdc77540f401f6a91aedc for contract 4f84a3b4a32125056b5b2bf5b6c7addc3616985f
```   

Что означает, что наш контракт был развернут, и теперь мы можем вызывать его.

#### Шаг #2
Давайте вызовем контракт для осуществления операций с nep5.

Для начала, запросите имя созданного токена nep5:

```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 4f84a3b4a32125056b5b2bf5b6c7addc3616985f name
```                                                                   
Где
- `4f84a3b4a32125056b5b2bf5b6c7addc3616985f` - хеш нашего контракта, полученный на шаге #1.
- `name` - строка операции, описанная ранее и возвращающая имя токена.

... не забудьте пароль от аккаунта `qwerty`.

Результат:
```
Sent invocation transaction de56bd31d785e6656a5eb20f48fd94c98d19700a9de4e721aebf68ba585490f2
```                                                                                         
Теперь давайте подробнее посмотрим на полученную вызывающую транзакцию с помощью `getapplicationlog` RPC-вызова:

```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["de56bd31d785e6656a5eb20f48fd94c98d19700a9de4e721aebf68ba585490f2"] }' localhost:20331 | json_pp
```               

Результат:
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

Поле `stack` полученного JSON-сообщения содержит массив байтов в base64 со значением имени токена.

Следующие команды позволят получить вам дополнительную информацию о токене:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 4f84a3b4a32125056b5b2bf5b6c7addc3616985f symbol
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 4f84a3b4a32125056b5b2bf5b6c7addc3616985f decimals
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 4f84a3b4a32125056b5b2bf5b6c7addc3616985f totalSupply
```

#### Шаг #3
Настало время для более интересных вещей. Для начала проверим баланс nep5 токенов на нашем счету с помощью метода `balanceOf`:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 4f84a3b4a32125056b5b2bf5b6c7addc3616985f balanceOf NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt
```                             
... с паролем `qwerty`. Результат:
```
Sent invocation transaction de1a69a6c0b5389458f35c44e28d1bea0a0ab4fe7b4ad0013e34b36d0995069a
```
Для более детального рассмотрения транзакции используем `getapplicationlog` RPC-вызов:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["de1a69a6c0b5389458f35c44e28d1bea0a0ab4fe7b4ad0013e34b36d0995069a"] }' localhost:20331 | json_pp
```
Результат:
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
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 4f84a3b4a32125056b5b2bf5b6c7addc3616985f mint NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt -- f8d79b436e9d3a54a38cc877d182b71e381d7a5c:CalledByEntry
```
Где
- `--` специальный разделитель, служащий для обозначения списка подписантов транзакции
- `f8d79b436e9d3a54a38cc877d182b71e381d7a5c` сам подписант транзакции (шестнадцатиричное LE представление нашего аккаунта `NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt`)

... с паролем `qwerty`. Результат:
``` 
Sent invocation transaction 2f4eff5532808080ae4c6a1a844c32a53015ca3cad5c308c5bcf3d8261399d06
```
`getapplicationlog` RPC-вызов для этой транзакции дает нам следующий результат:
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
Обратите внимание, что поле `stack` содержит значение `1` - токен был успешно выпущен.
Давайте убедимся в этом, еще раз запросив баланс нашего аккаунта с помощью метода `balanceOf`:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 4f84a3b4a32125056b5b2bf5b6c7addc3616985f balanceOf NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt
```
... пароль `qwerty`. Результат:
``` 
Sent invocation transaction cb3a28b2670836d1034655c19b5fbea62593f5f834977e72b1599f28c0ea79d5
```
... со следующим сообщением от `getapplicationlog` вызова RPC:
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
Теперь мы видим целое значение в поле `stack`, а именно, `1100000000000000` является значением баланса токена nep5 на нашем аккаунте.

Важно, что токен может быть выпущен лишь однажды.

#### Шаг #5

После того, как мы закончили с выпуском токена, мы можем перевести некоторое количество токена кому-нибудь.
Давайте переведем 5 токенов аккаунту с адресом `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm` с помощью функции `transfer`:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 4f84a3b4a32125056b5b2bf5b6c7addc3616985f transfer NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm 500000000 -- f8d79b436e9d3a54a38cc877d182b71e381d7a5c:CalledByEntry
```
... пароль `qwerty`. Результат:
``` 
Sent invocation transaction 65864bd07566318141fd781ece65a72d3c185f3f7190ac7096f7b9b7b62583ee
```
Наш любимый вызов RPC `getapplicationlog` говорит нам:
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
Заметьте, что поле `stack` содержит `1`, что означает, что токен был успешно переведен с нашего аккаунта.
Теперь давайте проверим баланс аккаунта, на который был совершен перевод (`NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm`), чтобы убедиться, что количество токена на нем = 5:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 4f84a3b4a32125056b5b2bf5b6c7addc3616985f balanceOf NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm
```
Вызов RPC `getapplicationlog` для этой транзакции возвращает следующий результат:
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
Sent deployment transaction 2796d6d641d3e7528f7cb7c30cc8b0099e1291e2af3a14ef4a8d89f3dca14c44 for contract 459f3383e8b087a67454da141eeba34521eaccfc
```   
Вы догадываетесь, что это значит :)

#### Шаг #3

Вызовите контракт, чтобы зарегистрировать домен с именем `my_first_domain`: 
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 459f3383e8b087a67454da141eeba34521eaccfc register my_first_domain NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt -- f8d79b436e9d3a54a38cc877d182b71e381d7a5c:CalledByEntry
```
... пароль: `qwerty`
```
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```
Результат:
```
Sent invocation transaction f6995b104a03fdb3dc45b236c371523276b8c17036631e3df9fdc5575d608988
```
Также вы можете увидеть лог-сообщение в консоли, где запускали ноду neo-go:
```
2020-07-06T16:50:37.713+0300	INFO	runtime log	{"script": "fcccea2145a3eb1e14da5474a687b0e883339f45", "logs": "\"RegisterDomain: my_first_domain\""}
```
Все получилось. Теперь проверим, был ли наш домен действительно зарегистрирован, с помощью вызова RPC `getapplicationlog`:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["f6995b104a03fdb3dc45b236c371523276b8c17036631e3df9fdc5575d608988"] }' localhost:20331 | json_pp
```
Результат:
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
Здесь мы в особенности заинтересованы в двух полях полученного json:

Первое поле - `notifications`, оно содержит одно уведомление с именем `registered`:
- `bXlfZmlyc3RfZG9tYWlu` - строка в base64, которая может быть декодирована в имя нашего домена - `my_first_domain`,
- `E9RIbx59YF1EGCw9f0E19Ms6o/A=` - строка, которая декодируется в адрес нашего аккаунта `NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt`.

Второе поле - `stack`, в котором лежит `1` - значение, возвращенное смарт-контрактом.

Все эти значения дают нам понять, что наш домен был успешно зарегистрирован.    

#### Шаг #4

Вызовите контракт, чтобы запросить информацию об адресе аккаунта, зарегистрировавшего домен `my_first_domain`:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 459f3383e8b087a67454da141eeba34521eaccfc query my_first_domain
```
... любимейший пароль `qwerty`:
```
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```
Результат:
```
Sent invocation transaction 7890938840734c26a66ed642ab50ce9d5f1c075aa090a901c87fccb638d92012
```
и лог-сообщение в консоли запущенной ноды neo-go:
```
2020-07-06T17:02:10.782+0300	INFO	runtime log	{"script": "fcccea2145a3eb1e14da5474a687b0e883339f45", "logs": "\"QueryDomain: my_first_domain\""}
```
Проверим транзакцию с помощью вызова RPC `getapplicationlog`:
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["7890938840734c26a66ed642ab50ce9d5f1c075aa090a901c87fccb638d92012"] }' localhost:20331 | json_pp
```
... что даст нам следующий результат:
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

с base64 представлением адреса аккаунта `NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt` на стеке и в уведомлениях, что означает, что домен `my_first_domain` был зарегистрирован владельцем с полученным адресом аккаунта.

#### Шаг #5

Вызовите контракт для передачи домена другому аккаунту (например, аккаунту с адресом `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm`):
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 459f3383e8b087a67454da141eeba34521eaccfc transfer my_first_domain NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm -- f8d79b436e9d3a54a38cc877d182b71e381d7a5c
```
... пароль: `qwerty`
```
Enter account NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt password >
```
Результат:
```
Sent invocation transaction 7df4b7d4f318e1e2c3007ec919bbcbff620685df71add4834e0f1a35f57ef893
```
и лог-сообщение:
```
2020-07-06T17:10:08.521+0300	INFO	runtime log	{"script": "fcccea2145a3eb1e14da5474a687b0e883339f45", "logs": "\"TransferDomain: my_first_domain\""}
```
Отлично. И `getapplicationlog` вызов RPC...
```
curl -d '{ "jsonrpc": "2.0", "id": 1, "method": "getapplicationlog", "params": ["7df4b7d4f318e1e2c3007ec919bbcbff620685df71add4834e0f1a35f57ef893"] }' localhost:20331 | json_pp
```
... говорит нам:
```
{
   "result" : {
      "gas_consumed" : "6390910",
      "stack" : [
         {
            "value" : "1",
            "type" : "Integer"
         }
      ],
      "trigger" : "Application",
      "vmstate" : "HALT",
      "txid" : "0x5c1da38f000abec53dc135dc404e768a1613749c3377b73d960782fac464d941",
      "notifications" : [
         {
            "state" : {
               "type" : "Array",
               "value" : [
                  {
                     "value" : "ZGVsZXRlZA==",
                     "type" : "ByteString"
                  },
                  {
                     "type" : "ByteString",
                     "value" : "E9RIbx59YF1EGCw9f0E19Ms6o/A="
                  },
                  {
                     "type" : "ByteString",
                     "value" : "bXlfZmlyc3RfZG9tYWlu"
                  }
               ]
            },
            "contract" : "0x459f3383e8b087a67454da141eeba34521eaccfc"
         },
         {
            "state" : {
               "value" : [
                  {
                     "value" : "cmVnaXN0ZXJlZA==",
                     "type" : "ByteString"
                  },
                  {
                     "type" : "ByteString",
                     "value" : "50l6vFaauRKm8hPVkr3Aw2CeHQs="
                  },
                  {
                     "type" : "ByteString",
                     "value" : "bXlfZmlyc3RfZG9tYWlu"
                  }
               ],
               "type" : "Array"
            },
            "contract" : "0x459f3383e8b087a67454da141eeba34521eaccfc"
         }
      ]
   },
   "jsonrpc" : "2.0",
   "id" : 1
}
```
Поле `notifications` содержит два события:
- первое с именем `deleted` и полями с дополнительной информацией (домен `my_first_domain` был удален с аккаунта с адресом `NULwe3UAHckN2fzNdcVg31tDiaYtMDwANt`),
- второе с именем `registered` и полями с дополнительной информацией (домен `my_first_domain` был зарегистрирован аккаунтом с адресом `NgzuJWWGVEwFGsRrgzj8knswEYRJrTe7sm`).
Поле `stack` содержит `1`, что значит, что домен был успешно перемещен.

#### Шаг #6

Оставшийся вызов - `delete`, вы можете попробовать выполнить его самостоятельно, создав перед этим еще один домен, например, с именем `my_second_domain`, а затем удалить его из хранилища с помощью:
```
$ ./bin/neo-go contract invokefunction -r http://localhost:20331 -w my_wallet.json 459f3383e8b087a67454da141eeba34521eaccfc delete my_second_domain -- f8d79b436e9d3a54a38cc877d182b71e381d7a5c
```

Спасибо!

### Полезные ссылки

* [Наш воркшоп на Medium](https://medium.com/@neospcc/%D1%81%D0%BC%D0%B0%D1%80%D1%82-%D0%BA%D0%BE%D0%BD%D1%82%D1%80%D0%B0%D0%BA%D1%82-%D0%B4%D0%BB%D1%8F-neo-769139352b65)
* [Использование NEO Blockchain Toolkit](https://medium.com/@neospcc/neogo-adds-support-for-neo-blockchain-toolkit-673ea914f661)
* [Документация NEO](https://docs.neo.org/)
* [NEO github](https://github.com/neo-project/neo/)
* [NEO-GO github](https://github.com/nspcc-dev/neo-go)
