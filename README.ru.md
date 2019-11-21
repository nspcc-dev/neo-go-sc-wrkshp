## Требования

Debian 10, Docker, docker-compose, go:
* https://docs.docker.com/install/linux/docker-ce/debian/
* https://golang.org/dl/

## Развёртывание

```
$ git clone https://github.com/CityOfZion/neo-local.git
$ cd neo-local
$ git checkout -b 4nodes 0.12
$ make start
```

### neo-go

```
$ git clone https://github.com/nspcc-dev/neo-go
$ cd neo-go
$ make build
```

## RPC

```
$ curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getversion", "params": [] }' localhost:20331
$ curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getblockcount", "params": [] }' localhost:20331
$ curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getconnectioncount", "params": [] }' localhost:20331
$ curl -d '{ "jsonrpc": "2.0", "id": 5, "method": "getaccountstate", "params": ["AK2nJJpJr6o664CWJKi1QRXjqeic2zRp8y"] }' localhost:20331
```


## Компиляция
```
$ neo-go contract compile -i $FILE
$ neo-go contract inspect -i $FILE
$ neo-go contract inspect -c -i $FILE
```



## Вызов
```
$ neo-go contract testinvoke -e 127.0.0.1:20331 -i $AVM
```
## Пошаговое выполнение запуска контракта

Условия:
- Запущен neo-local
- Скомпилирован neo-go

### Запуск Hello world контракта (1-print.go)
1. `$ ./bin/neo-go contract compile -i '/1-print.go'`

Где 
- `./bin/neo-go` запускает neo-go из корня.
- `contract compile` комада с аргументом клиента neo-go https://github.com/nspcc-dev/neo-go/blob/master/cli/smartcontract/smart_contract.go#L43
- `-i '/1-print.go'` путь к смарт-контракту

Результат:

Скомпилированный контракт: `1-pring.avm`

2. Скопировать скомпилированный контракт в docker neo-python:
`$ sudo docker cp smart-contracts/1-print.avm  neo-python:/neo-python`

3. Деплой контракта из neo-cli
- Открыть neo-cli консоль
- `sc deploy 1-print.avm True False False  07 05 --fee=1.5`

Результат:
```
Please fill out the following contract details:
[Contract Name] >                                                                                                                                                                                                                      
[Contract Version] >                                                                                                                                                                                                                   
[Contract Author] >                                                                                                                                                                                                                    
[Contract Email] >                                                                                                                                                                                                                     
[Contract Description] >                                                                                                                                                                                                               
Creating smart contract....
                 Name:  
              Version: 
               Author:  
                Email:  
          Description:  
        Needs Storage: True 
 Needs Dynamic Invoke: False 
           Is Payable: False 
{
    "hash": "0x4cf87bd149748abc3cf5fe46040bb8b4c2c0b2c7",
    "script": "5ec56b68164e656f2e53746f726167652e476574436f6e74657874616a00527ac410746573742d73746f726167652d6b65796a51527ac46a00c36a51c37c680f4e656f2e53746f726167652e476574616a52527ac41756616c756520726561642066726f6d2073746f726167656a53527ac46a53c368124e656f2e52756e74696d652e4e6f74696679616a52c3009c6447002553746f72616765206b6579206e6f7420796574207365742e2053657474696e6720746f203168124e656f2e52756e74696d652e4e6f7469667961516a52527ac4624d002a53746f72616765206b657920616c7265616479207365742e20496e6372656d656e74696e67206279203168124e656f2e52756e74696d652e4e6f74696679616a52c351936a52527ac46a00c36a51c36a52c35272680f4e656f2e53746f726167652e507574611e4e65772076616c7565207772697474656e20696e746f2073746f7261676568124e656f2e52756e74696d652e4e6f74696679616a52c36c7566",
    "parameters": [
        "String",
        "Array"
    ],
    "returntype": "ByteArray"
}

``` 

4. Запуск контракта

`sc invoke 0x4cf87bd149748abc3cf5fe46040bb8b4c2c0b2c7`

Где `0x4cf87bd149748abc3cf5fe46040bb8b4c2c0b2c7` hashcode контракта, который можно увидеть из логов на предыдущем шаге.

Результат

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

Необходимо ввести пароль `coz`


### Запуск контракта с использованием хранилища (2-storage.go)

Этот контракт использует системные вызовы(interops), которые позволяют сохранить значение переменной и обращаться к ней из другого контракта.

Компиляция и запуск осуществляется аналогичным образом как и с предыдущим контрактом.

Результат:

```
neo> [I 191031 10:00:09 EventHub:62] [SmartContract.Storage.Get][7758] [4cf87bd149748abc3cf5fe46040bb8b4c2c0b2c7] [tx 50961aedc1cafe460ec7fa1048347c01f59dd4e286d08bbd323088d25f81a959] {'type': 'String', 'value': "b'test-storage-key' -> bytearray(b'')"}
[I 191031 10:00:09 EventHub:62] [SmartContract.Storage.Put][7758] [4cf87bd149748abc3cf5fe46040bb8b4c2c0b2c7] [tx 50961aedc1cafe460ec7fa1048347c01f59dd4e286d08bbd323088d25f81a959] {'type': 'String', 'value': "b'test-storage-key' -> bytearray(b'\\x01')"}
[I 191031 10:00:09 EventHub:62] [SmartContract.Runtime.Notify][7758] [4cf87bd149748abc3cf5fe46040bb8b4c2c0b2c7] [tx 50961aedc1cafe460ec7fa1048347c01f59dd4e286d08bbd323088d25f81a959] {'type': 'ByteArray', 'value': b'Value read from storage'}
[I 191031 10:00:09 EventHub:62] [SmartContract.Runtime.Notify][7758] [4cf87bd149748abc3cf5fe46040bb8b4c2c0b2c7] [tx 50961aedc1cafe460ec7fa1048347c01f59dd4e286d08bbd323088d25f81a959] {'type': 'ByteArray', 'value': b'Storage key not yet set. Setting to 1'}                                                                                                                                                                                                                
[I 191031 10:00:09 EventHub:62] [SmartContract.Runtime.Notify][7758] [4cf87bd149748abc3cf5fe46040bb8b4c2c0b2c7] [tx 50961aedc1cafe460ec7fa1048347c01f59dd4e286d08bbd323088d25f81a959] {'type': 'ByteArray', 'value': b'New value written into storage'}
[I 191031 10:00:09 EventHub:62] [SmartContract.Execution.Success][7758] [4cf87bd149748abc3cf5fe46040bb8b4c2c0b2c7] [tx 50961aedc1cafe460ec7fa1048347c01f59dd4e286d08bbd323088d25f81a959] {'type': 'Array', 'value': [{'type': 'Integer', 'value': '1'}]}
```

При повторном вызове контракта результирующая величина будет инкременетирована до двух `'value': '2'}...`:
```
neo> [I 191031 10:01:11 EventHub:62] [SmartContract.Storage.Get][7770] [4cf87bd149748abc3cf5fe46040bb8b4c2c0b2c7] [tx bcba05164519dce40ce4ab1875f61e2e0cc7854bf7e690ce66f381a8567e8236] {'type': 'String', 'value': "b'test-storage-key' -> bytearray(b'\\x01')"}
[I 191031 10:01:11 EventHub:62] [SmartContract.Storage.Put][7770] [4cf87bd149748abc3cf5fe46040bb8b4c2c0b2c7] [tx bcba05164519dce40ce4ab1875f61e2e0cc7854bf7e690ce66f381a8567e8236] {'type': 'String', 'value': "b'test-storage-key' -> bytearray(b'\\x02')"}
[I 191031 10:01:11 EventHub:62] [SmartContract.Runtime.Notify][7770] [4cf87bd149748abc3cf5fe46040bb8b4c2c0b2c7] [tx bcba05164519dce40ce4ab1875f61e2e0cc7854bf7e690ce66f381a8567e8236] {'type': 'ByteArray', 'value': b'Value read from storage'}
[I 191031 10:01:11 EventHub:62] [SmartContract.Runtime.Notify][7770] [4cf87bd149748abc3cf5fe46040bb8b4c2c0b2c7] [tx bcba05164519dce40ce4ab1875f61e2e0cc7854bf7e690ce66f381a8567e8236] {'type': 'ByteArray', 'value': b'Storage key already set. Incrementing by 1'}                                                                                                                                                                                                           
[I 191031 10:01:11 EventHub:62] [SmartContract.Runtime.Notify][7770] [4cf87bd149748abc3cf5fe46040bb8b4c2c0b2c7] [tx bcba05164519dce40ce4ab1875f61e2e0cc7854bf7e690ce66f381a8567e8236] {'type': 'ByteArray', 'value': b'New value written into storage'}
[I 191031 10:01:11 EventHub:62] [SmartContract.Execution.Success][7770] [4cf87bd149748abc3cf5fe46040bb8b4c2c0b2c7] [tx bcba05164519dce40ce4ab1875f61e2e0cc7854bf7e690ce66f381a8567e8236] {'type': 'Array', 'value': [{'type': 'Integer', 'value': '2'}]}
```


# Ссылки

* https://medium.com/@neospcc/%D1%81%D0%BC%D0%B0%D1%80%D1%82-%D0%BA%D0%BE%D0%BD%D1%82%D1%80%D0%B0%D0%BA%D1%82-%D0%B4%D0%BB%D1%8F-neo-769139352b65
* https://docs.neo.org/
* https://github.com/neo-project/neo/
* https://github.com/nspcc-dev/neo-go
