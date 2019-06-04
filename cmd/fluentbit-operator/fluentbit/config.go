package fluentbit

var fluentBitConfigTemplate = `
[SERVICE]
    Flush        1
    Daemon       Off
    Log_Level    info
    Parsers_File parsers.conf

[INPUT]
    Name             tail
    Path             /var/log/containers/*.log
    Parser           docker
    Tag              kube.*
    Refresh_Interval 5
    Mem_Buf_Limit    5MB
    Skip_Long_Lines  On
    DB               /tail-db/tail-containers-state.db
    DB.Sync          Normal

    [FILTER]
    Name                kubernetes
    Match               kube.*
    Kube_URL            https://kubernetes.default.svc:443
    Kube_CA_File        /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
    Kube_Token_File     /var/run/secrets/kubernetes.io/serviceaccount/token

[FILTER]
    Name                nest
    Match               kube.*
    Operation           lift
    Nested_under        kubernetes
    Prefix_with         kubernetes_

[FILTER]
    Name                modify
    Match               kube.*
    Remove              stream

[FILTER]
    Name                modify
    Match               kube.*
    Remove              kubernetes_labels

[FILTER]
    Name                modify
    Match               kube.*
    Remove              kubernetes_annotations

[FILTER]
    Name                modify
    Match               kube.*
    Remove              kubernetes_pod_id

[FILTER]
    Name                modify
    Match               kube.*
    Remove              kubernetes_docker_id

[FILTER]
    Name                nest
    Match               kube.*
    Operation           nest
    Wildcard            kubernetes_*
    Nested_under        kubernetes
    Remove_prefix       kubernetes_

[OUTPUT]
    Name  null
    Match *
`

var fluentBitSettingsTemplate = `
Enable     1
`
