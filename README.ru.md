## Требования

Debian 10, Docker, docker-compose, go
https://docs.docker.com/install/linux/docker-ce/debian/
https://golang.org/dl/

## Развёртывание

```
$ git clone https://github.com/CityOfZion/neo-local.git
$ cd neo-local
$ git checkout -b 4nodes 0.12
$ make start
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


# Ссылки

https://github.com/neo-project/neo/
https://github.com/nspcc-dev/neo-go
https://medium.com/@neospcc/%D1%81%D0%BC%D0%B0%D1%80%D1%82-%D0%BA%D0%BE%D0%BD%D1%82%D1%80%D0%B0%D0%BA%D1%82-%D0%B4%D0%BB%D1%8F-neo-769139352b65
