package constants

const (
	SecretVolName    = "config"
	FluentdMountPath = "/fluentd/etc"
	BufferMountPath  = "/buffers"

	MetricsName = "metrics"

	DefaultBind        string = "0.0.0.0"
	DefaultForwardPort int32  = 24424
	DefaultHttpPort    int32  = 9880
	DefaultMetricsPort int32  = 2021
	// 101 is the fsGroup that fluentd runs as in the kubesphere image
	DefaultFsGroup int64 = 101

	DefaultForwardName = "forward"
	DefaultHttpName    = "http"

	InputForwardType = "forward"
	InputHttpType    = "http"
)
