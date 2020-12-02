# Forward Error Logs to A Remote HTTP Endpoint

This tutorial guides you on how to filter error logs and send them out via HTTP.

# Demo

First, deploy HTTP sample server for receiving logs (see main.go).

```shell
kubectl create deploy log-receiver --image=kubespheredev/example-log-receiver
kubectl expose deploy log-receiver --port=8080 --target-port=8080
```  
Second, setup the logging pipeline. It deploys a log generator and logging agents (fluent bit). The logging agent will forward all error logs containing `ERRO` to the sample server above.

```shell
kubectl apply -f manifests/setup/
kubectl apply -f docs/user-guides/forwarding-logs-via-http/deploy/
```

Note: for KubeSphere users who have enabled logging, you can simply apply the YAML files in the folder `kubesphere`. Don't forget to adapt the http receiver endpoint to your setup.

# Sample Output

The logging agents (fluent bit) forward logs in the following format:

```json
[
    {
        "date":1603182646.218455,
        "log":"16:47:09.634 [peer] GetLocalAddress -> ERRO 033 Auto detected peer address: 9.3.158.178:30303",
        "time":"2020-10-20T08:30:46.218455275Z",
        "kubernetes":{
            "pod_name":"log-generator",
            "namespace_name":"default",
            "container_name":"log-generator",
            "docker_id":"aa740a1b77bf181dc8f26723c848e0692a214e318468b9818fb19611136cd360",
            "container_image":"busybox:latest"
        }
    },
    {
        "date":1603182647.220342,
        "log":"16:47:09.634 [peer] GetLocalAddress -> ERRO 033 Auto detected peer address: 9.3.158.178:30303",
        "time":"2020-10-20T08:30:47.220342086Z",
        "kubernetes":{
            "pod_name":"log-generator",
            "namespace_name":"default",
            "container_name":"log-generator",
            "docker_id":"aa740a1b77bf181dc8f26723c848e0692a214e318468b9818fb19611136cd360",
            "container_image":"busybox:latest"
        }
    },
    {
        "date":1603182648.221025,
        "log":"16:47:09.634 [peer] GetLocalAddress -> ERRO 033 Auto detected peer address: 9.3.158.178:30303",
        "time":"2020-10-20T08:30:48.221024896Z",
        "kubernetes":{
            "pod_name":"log-generator",
            "namespace_name":"default",
            "container_name":"log-generator",
            "docker_id":"aa740a1b77bf181dc8f26723c848e0692a214e318468b9818fb19611136cd360",
            "container_image":"busybox:latest"
        }
    },
    {
        "date":1603182649.222065,
        "log":"16:47:09.634 [peer] GetLocalAddress -> ERRO 033 Auto detected peer address: 9.3.158.178:30303",
        "time":"2020-10-20T08:30:49.222065229Z",
        "kubernetes":{
            "pod_name":"log-generator",
            "namespace_name":"default",
            "container_name":"log-generator",
            "docker_id":"aa740a1b77bf181dc8f26723c848e0692a214e318468b9818fb19611136cd360",
            "container_image":"busybox:latest"
        }
    },
    {
        "date":1603182650.222759,
        "log":"16:47:09.634 [peer] GetLocalAddress -> ERRO 033 Auto detected peer address: 9.3.158.178:30303",
        "time":"2020-10-20T08:30:50.222759014Z",
        "kubernetes":{
            "pod_name":"log-generator",
            "namespace_name":"default",
            "container_name":"log-generator",
            "docker_id":"aa740a1b77bf181dc8f26723c848e0692a214e318468b9818fb19611136cd360",
            "container_image":"busybox:latest"
        }
    }
]
```

You may check out the logs of log-receiver.

# TLS Support

Please read comments on the file `deploy/output-http.yaml` for how to send logs to an HTTPS endpoint.