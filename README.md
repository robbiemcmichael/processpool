# processpool

Maintains a pool of processes which can be fetched via HTTP.

Each process is assumed to listen on a TCP port, so `processpool` will
automatically assign a unique port to each process and wait for the TCP port to
be ready before serving it to clients via HTTP. The client is responsible for
shutting down the process (e.g. by killing the PID) after which `processpool`
will start a new process to reuse that port.

## Usage

```
Usage:
  processpool [flags] -- command [args]

Flags:
  -a, --address string   address for processpool to listen on (default "127.0.0.1:8080")
  -h, --help             help for processpool
  -p, --port int         starting port for the first process to use (required)
  -n, --processes int    number of processes (required)
```

## Example

Run `processpool` to maintain a pool of 10 Python HTTP servers:

```bash
bin/processpool --port 33000 --processes 10 -- bash -c 'python -m http.server $PROCESSPOOL_PORT'
```

Use `curl` to fetch a PID and port for the Python process:

```bash
$ curl http://127.0.0.1:8080
{"pid":14159,"port":33000}
```
