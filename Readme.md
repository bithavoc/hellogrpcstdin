# GRPC external plugin example

```
go get github.com/bithavoc/hellogrpcstdin
```

An example of GRPC Dialing with Stdin and Stdout of a child process.

`greeter_server` launches `greeter_client` as a process and uses `os.Stdin` and `os.Stdout` to listen for GRPC requests.

**Note:** You don't really need this library to implement GRPC plugins in Go, you can always use [hashicorp/go-plugin](https://github.com/hashicorp/go-plugin) which uses unix sockets.

## Running Example

**Compile Client(Plugin):**

```shell
(cd greeter_client; go build)
```

**Compile Server(Plugin):**

```shell
(cd greeter_server; go build)
```

**Run Server(Plugin):**

```shell
(cd greeter_server; ./greeter_server)
```

You should see an output similar to this:

```
Server starting
accepting, waiting for only conn
agent process running, will ready listener
accepting, only conn ready
accepting, waiting for only conn
2018/07/16 19:13:33 Server replying, check greeter_server/client.log
```

`cat greeter_server/client.log`:

```
2018/07/16 19:15:51 client starting
2018/07/16 19:15:51 addresss 
2018/07/16 19:15:51 Greeting: Hello world from server
```
