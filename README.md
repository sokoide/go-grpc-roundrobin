# go-grpc-roundrobin

## How to build

```
make
```

## How to run

* Register a DNS with 2 A records (e.g. case2.sokoide.com -> 192.168.1.101 and 102)
* Run the server on both 101 and 102 server
* Run client from the 3rd host

```
./client -host case2.sokoide.com -roundrobin
```

* Check the result -> with `-roundrobin`, it's roundrobin
