# Forward Error Logs to A Remote HTTP Endpoint

This tutorial guides you on how to filter error logs and send them out via HTTP.

# Demo

First, deploy HTTP sample server for receiving logs (see main.go).

```shell
kubectl create deploy log-receiver --image=vasth/web0:latest --port=8080
kubectl expose deploy log-receiver --port=8080 --target-port=8080
```  
Second, setup the logging pipeline. It deploys a log generator and logging agents (fluent bit). The logging agent will forward all error logs containing `ERRO` to the sample server above.

```shell
kubectl apply -f manifests/setup/
kubectl apply -f docs/user-guides/forwarding-logs-via-http/deploy/
```

# Sample Output

The logging agents (fluent bit) forward logs in the following format:

```json
[
    {
        "date":1603093215.693762,
        "log":"16:47:09.634 [peer] GetLocalAddress -> ERRO 033 Auto detected peer address: 9.3.158.178:30303",
        "stream":"stdout",
        "time":"2020-10-19T07:40:15.693761554Z"
    },
    {
        "date":1603093216.694311,
        "log":"16:47:09.634 [peer] GetLocalAddress -> ERRO 033 Auto detected peer address: 9.3.158.178:30303",
        "stream":"stdout",
        "time":"2020-10-19T07:40:16.694310592Z"
    },
    {
        "date":1603093217.695296,
        "log":"16:47:09.634 [peer] GetLocalAddress -> ERRO 033 Auto detected peer address: 9.3.158.178:30303",
        "stream":"stdout",
        "time":"2020-10-19T07:40:17.695295635Z"
    },
    {
        "date":1603093218.696355,
        "log":"16:47:09.634 [peer] GetLocalAddress -> ERRO 033 Auto detected peer address: 9.3.158.178:30303",
        "stream":"stdout",
        "time":"2020-10-19T07:40:18.696354877Z"
    },
    {
        "date":1603093219.697264,
        "log":"16:47:09.634 [peer] GetLocalAddress -> ERRO 033 Auto detected peer address: 9.3.158.178:30303",
        "stream":"stdout",
        "time":"2020-10-19T07:40:19.69726438Z"
    }
]
```