# Monitoring

Fluent Bit comes with a built-in HTTP Server. According to the official [documentation](https://docs.fluentbit.io/manual/administration/monitoring) of fluentbit You can enable this by enabling the HTTP server from the fluent bit configuration file:

```conf
[SERVICE]
    HTTP_Server  On
    HTTP_Listen  0.0.0.0
    HTTP_PORT    2020
```

When you use the fluent-operator, You can enable this from `FluentBitConfig` manifest. Example is below:

```yaml
apiVersion: fluentbit.fluent.io/v1alpha2
kind: ClusterFluentBitConfig
metadata:
  name: fluent-bit-config
  labels:
    app.kubernetes.io/name: fluent-bit
spec:
  filterSelector:
    matchLabels:
      fluentbit.fluent.io/enabled: 'true'
  inputSelector:
    matchLabels:
      fluentbit.fluent.io/enabled: 'true'
  outputSelector:
    matchLabels:
      fluentbit.fluent.io/enabled: 'true'
  service:
    httpListen: 0.0.0.0
    httpPort: 2020
    httpServer: true
    parsersFile: parsers.conf

```

Once HTTP server is enabled, you should be able to get the information:

```bash
curl <podIP>:2020 | jq .

{
  "fluent-bit": {
    "version": "1.8.3",
    "edition": "Community",
    "flags": [
      "FLB_HAVE_PARSER",
      "FLB_HAVE_RECORD_ACCESSOR",
      "FLB_HAVE_STREAM_PROCESSOR",
      "FLB_HAVE_TLS",
      "FLB_HAVE_OPENSSL",
      "FLB_HAVE_AWS",
      "FLB_HAVE_SIGNV4",
      "FLB_HAVE_SQLDB",
      "FLB_HAVE_METRICS",
      "FLB_HAVE_HTTP_SERVER",
      "FLB_HAVE_SYSTEMD",
      "FLB_HAVE_FORK",
      "FLB_HAVE_TIMESPEC_GET",
      "FLB_HAVE_GMTOFF",
      "FLB_HAVE_UNIX_SOCKET",
      "FLB_HAVE_PROXY_GO",
      "FLB_HAVE_JEMALLOC",
      "FLB_HAVE_LIBBACKTRACE",
      "FLB_HAVE_REGEX",
      "FLB_HAVE_UTF8_ENCODER",
      "FLB_HAVE_LUAJIT",
      "FLB_HAVE_C_TLS",
      "FLB_HAVE_ACCEPT4",
      "FLB_HAVE_INOTIFY"
    ]
  }
}
```

