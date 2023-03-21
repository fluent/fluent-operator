//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package output

import (
	"github.com/fluent/fluent-operator/apis/fluentbit/v1alpha2/plugins"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AzureBlob) DeepCopyInto(out *AzureBlob) {
	*out = *in
	if in.SharedKey != nil {
		in, out := &in.SharedKey, &out.SharedKey
		*out = new(plugins.Secret)
		(*in).DeepCopyInto(*out)
	}
	if in.AutoCreateContainer != nil {
		in, out := &in.AutoCreateContainer, &out.AutoCreateContainer
		*out = new(bool)
		**out = **in
	}
	if in.EmulatorMode != nil {
		in, out := &in.EmulatorMode, &out.EmulatorMode
		*out = new(bool)
		**out = **in
	}
	if in.TLS != nil {
		in, out := &in.TLS, &out.TLS
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AzureBlob.
func (in *AzureBlob) DeepCopy() *AzureBlob {
	if in == nil {
		return nil
	}
	out := new(AzureBlob)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AzureLogAnalytics) DeepCopyInto(out *AzureLogAnalytics) {
	*out = *in
	if in.CustomerID != nil {
		in, out := &in.CustomerID, &out.CustomerID
		*out = new(plugins.Secret)
		(*in).DeepCopyInto(*out)
	}
	if in.SharedKey != nil {
		in, out := &in.SharedKey, &out.SharedKey
		*out = new(plugins.Secret)
		(*in).DeepCopyInto(*out)
	}
	if in.TimeGenerated != nil {
		in, out := &in.TimeGenerated, &out.TimeGenerated
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AzureLogAnalytics.
func (in *AzureLogAnalytics) DeepCopy() *AzureLogAnalytics {
	if in == nil {
		return nil
	}
	out := new(AzureLogAnalytics)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DataDog) DeepCopyInto(out *DataDog) {
	*out = *in
	if in.TLS != nil {
		in, out := &in.TLS, &out.TLS
		*out = new(bool)
		**out = **in
	}
	if in.IncludeTagKey != nil {
		in, out := &in.IncludeTagKey, &out.IncludeTagKey
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataDog.
func (in *DataDog) DeepCopy() *DataDog {
	if in == nil {
		return nil
	}
	out := new(DataDog)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Elasticsearch) DeepCopyInto(out *Elasticsearch) {
	*out = *in
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(int32)
		**out = **in
	}
	if in.HTTPUser != nil {
		in, out := &in.HTTPUser, &out.HTTPUser
		*out = new(plugins.Secret)
		(*in).DeepCopyInto(*out)
	}
	if in.HTTPPasswd != nil {
		in, out := &in.HTTPPasswd, &out.HTTPPasswd
		*out = new(plugins.Secret)
		(*in).DeepCopyInto(*out)
	}
	if in.LogstashFormat != nil {
		in, out := &in.LogstashFormat, &out.LogstashFormat
		*out = new(bool)
		**out = **in
	}
	if in.TimeKeyNanos != nil {
		in, out := &in.TimeKeyNanos, &out.TimeKeyNanos
		*out = new(bool)
		**out = **in
	}
	if in.IncludeTagKey != nil {
		in, out := &in.IncludeTagKey, &out.IncludeTagKey
		*out = new(bool)
		**out = **in
	}
	if in.GenerateID != nil {
		in, out := &in.GenerateID, &out.GenerateID
		*out = new(bool)
		**out = **in
	}
	if in.ReplaceDots != nil {
		in, out := &in.ReplaceDots, &out.ReplaceDots
		*out = new(bool)
		**out = **in
	}
	if in.TraceOutput != nil {
		in, out := &in.TraceOutput, &out.TraceOutput
		*out = new(bool)
		**out = **in
	}
	if in.TraceError != nil {
		in, out := &in.TraceError, &out.TraceError
		*out = new(bool)
		**out = **in
	}
	if in.CurrentTimeIndex != nil {
		in, out := &in.CurrentTimeIndex, &out.CurrentTimeIndex
		*out = new(bool)
		**out = **in
	}
	if in.TLS != nil {
		in, out := &in.TLS, &out.TLS
		*out = new(plugins.TLS)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Elasticsearch.
func (in *Elasticsearch) DeepCopy() *Elasticsearch {
	if in == nil {
		return nil
	}
	out := new(Elasticsearch)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *File) DeepCopyInto(out *File) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new File.
func (in *File) DeepCopy() *File {
	if in == nil {
		return nil
	}
	out := new(File)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Firehose) DeepCopyInto(out *Firehose) {
	*out = *in
	if in.TimeKey != nil {
		in, out := &in.TimeKey, &out.TimeKey
		*out = new(string)
		**out = **in
	}
	if in.TimeKeyFormat != nil {
		in, out := &in.TimeKeyFormat, &out.TimeKeyFormat
		*out = new(string)
		**out = **in
	}
	if in.DataKeys != nil {
		in, out := &in.DataKeys, &out.DataKeys
		*out = new(string)
		**out = **in
	}
	if in.LogKey != nil {
		in, out := &in.LogKey, &out.LogKey
		*out = new(string)
		**out = **in
	}
	if in.RoleARN != nil {
		in, out := &in.RoleARN, &out.RoleARN
		*out = new(string)
		**out = **in
	}
	if in.Endpoint != nil {
		in, out := &in.Endpoint, &out.Endpoint
		*out = new(string)
		**out = **in
	}
	if in.STSEndpoint != nil {
		in, out := &in.STSEndpoint, &out.STSEndpoint
		*out = new(string)
		**out = **in
	}
	if in.AutoRetryRequests != nil {
		in, out := &in.AutoRetryRequests, &out.AutoRetryRequests
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Firehose.
func (in *Firehose) DeepCopy() *Firehose {
	if in == nil {
		return nil
	}
	out := new(Firehose)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Forward) DeepCopyInto(out *Forward) {
	*out = *in
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(int32)
		**out = **in
	}
	if in.TimeAsInteger != nil {
		in, out := &in.TimeAsInteger, &out.TimeAsInteger
		*out = new(bool)
		**out = **in
	}
	if in.SendOptions != nil {
		in, out := &in.SendOptions, &out.SendOptions
		*out = new(bool)
		**out = **in
	}
	if in.RequireAckResponse != nil {
		in, out := &in.RequireAckResponse, &out.RequireAckResponse
		*out = new(bool)
		**out = **in
	}
	if in.EmptySharedKey != nil {
		in, out := &in.EmptySharedKey, &out.EmptySharedKey
		*out = new(bool)
		**out = **in
	}
	if in.Username != nil {
		in, out := &in.Username, &out.Username
		*out = new(plugins.Secret)
		(*in).DeepCopyInto(*out)
	}
	if in.Password != nil {
		in, out := &in.Password, &out.Password
		*out = new(plugins.Secret)
		(*in).DeepCopyInto(*out)
	}
	if in.TLS != nil {
		in, out := &in.TLS, &out.TLS
		*out = new(plugins.TLS)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Forward.
func (in *Forward) DeepCopy() *Forward {
	if in == nil {
		return nil
	}
	out := new(Forward)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HTTP) DeepCopyInto(out *HTTP) {
	*out = *in
	if in.HTTPUser != nil {
		in, out := &in.HTTPUser, &out.HTTPUser
		*out = new(plugins.Secret)
		(*in).DeepCopyInto(*out)
	}
	if in.HTTPPasswd != nil {
		in, out := &in.HTTPPasswd, &out.HTTPPasswd
		*out = new(plugins.Secret)
		(*in).DeepCopyInto(*out)
	}
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(int32)
		**out = **in
	}
	if in.AllowDuplicatedHeaders != nil {
		in, out := &in.AllowDuplicatedHeaders, &out.AllowDuplicatedHeaders
		*out = new(bool)
		**out = **in
	}
	if in.Headers != nil {
		in, out := &in.Headers, &out.Headers
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.TLS != nil {
		in, out := &in.TLS, &out.TLS
		*out = new(plugins.TLS)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HTTP.
func (in *HTTP) DeepCopy() *HTTP {
	if in == nil {
		return nil
	}
	out := new(HTTP)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Kafka) DeepCopyInto(out *Kafka) {
	*out = *in
	if in.Rdkafka != nil {
		in, out := &in.Rdkafka, &out.Rdkafka
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.DynamicTopic != nil {
		in, out := &in.DynamicTopic, &out.DynamicTopic
		*out = new(bool)
		**out = **in
	}
	if in.QueueFullRetries != nil {
		in, out := &in.QueueFullRetries, &out.QueueFullRetries
		*out = new(int64)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Kafka.
func (in *Kafka) DeepCopy() *Kafka {
	if in == nil {
		return nil
	}
	out := new(Kafka)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Loki) DeepCopyInto(out *Loki) {
	*out = *in
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(int32)
		**out = **in
	}
	if in.HTTPUser != nil {
		in, out := &in.HTTPUser, &out.HTTPUser
		*out = new(plugins.Secret)
		(*in).DeepCopyInto(*out)
	}
	if in.HTTPPasswd != nil {
		in, out := &in.HTTPPasswd, &out.HTTPPasswd
		*out = new(plugins.Secret)
		(*in).DeepCopyInto(*out)
	}
	if in.TenantID != nil {
		in, out := &in.TenantID, &out.TenantID
		*out = new(plugins.Secret)
		(*in).DeepCopyInto(*out)
	}
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.LabelKeys != nil {
		in, out := &in.LabelKeys, &out.LabelKeys
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.TLS != nil {
		in, out := &in.TLS, &out.TLS
		*out = new(plugins.TLS)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Loki.
func (in *Loki) DeepCopy() *Loki {
	if in == nil {
		return nil
	}
	out := new(Loki)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Null) DeepCopyInto(out *Null) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Null.
func (in *Null) DeepCopy() *Null {
	if in == nil {
		return nil
	}
	out := new(Null)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OpenSearch) DeepCopyInto(out *OpenSearch) {
	*out = *in
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(int32)
		**out = **in
	}
	if in.HTTPUser != nil {
		in, out := &in.HTTPUser, &out.HTTPUser
		*out = new(plugins.Secret)
		(*in).DeepCopyInto(*out)
	}
	if in.HTTPPasswd != nil {
		in, out := &in.HTTPPasswd, &out.HTTPPasswd
		*out = new(plugins.Secret)
		(*in).DeepCopyInto(*out)
	}
	if in.LogstashFormat != nil {
		in, out := &in.LogstashFormat, &out.LogstashFormat
		*out = new(bool)
		**out = **in
	}
	if in.TimeKeyNanos != nil {
		in, out := &in.TimeKeyNanos, &out.TimeKeyNanos
		*out = new(bool)
		**out = **in
	}
	if in.IncludeTagKey != nil {
		in, out := &in.IncludeTagKey, &out.IncludeTagKey
		*out = new(bool)
		**out = **in
	}
	if in.GenerateID != nil {
		in, out := &in.GenerateID, &out.GenerateID
		*out = new(bool)
		**out = **in
	}
	if in.ReplaceDots != nil {
		in, out := &in.ReplaceDots, &out.ReplaceDots
		*out = new(bool)
		**out = **in
	}
	if in.TraceOutput != nil {
		in, out := &in.TraceOutput, &out.TraceOutput
		*out = new(bool)
		**out = **in
	}
	if in.TraceError != nil {
		in, out := &in.TraceError, &out.TraceError
		*out = new(bool)
		**out = **in
	}
	if in.CurrentTimeIndex != nil {
		in, out := &in.CurrentTimeIndex, &out.CurrentTimeIndex
		*out = new(bool)
		**out = **in
	}
	if in.SuppressTypeName != nil {
		in, out := &in.SuppressTypeName, &out.SuppressTypeName
		*out = new(bool)
		**out = **in
	}
	if in.Workers != nil {
		in, out := &in.Workers, &out.Workers
		*out = new(int32)
		**out = **in
	}
	if in.TLS != nil {
		in, out := &in.TLS, &out.TLS
		*out = new(plugins.TLS)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OpenSearch.
func (in *OpenSearch) DeepCopy() *OpenSearch {
	if in == nil {
		return nil
	}
	out := new(OpenSearch)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OpenTelemetry) DeepCopyInto(out *OpenTelemetry) {
	*out = *in
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(int32)
		**out = **in
	}
	if in.HTTPUser != nil {
		in, out := &in.HTTPUser, &out.HTTPUser
		*out = new(plugins.Secret)
		(*in).DeepCopyInto(*out)
	}
	if in.HTTPPasswd != nil {
		in, out := &in.HTTPPasswd, &out.HTTPPasswd
		*out = new(plugins.Secret)
		(*in).DeepCopyInto(*out)
	}
	if in.Header != nil {
		in, out := &in.Header, &out.Header
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.LogResponsePayload != nil {
		in, out := &in.LogResponsePayload, &out.LogResponsePayload
		*out = new(bool)
		**out = **in
	}
	if in.AddLabel != nil {
		in, out := &in.AddLabel, &out.AddLabel
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.TLS != nil {
		in, out := &in.TLS, &out.TLS
		*out = new(plugins.TLS)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OpenTelemetry.
func (in *OpenTelemetry) DeepCopy() *OpenTelemetry {
	if in == nil {
		return nil
	}
	out := new(OpenTelemetry)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PrometheusRemoteWrite) DeepCopyInto(out *PrometheusRemoteWrite) {
	*out = *in
	if in.HTTPUser != nil {
		in, out := &in.HTTPUser, &out.HTTPUser
		*out = new(plugins.Secret)
		(*in).DeepCopyInto(*out)
	}
	if in.HTTPPasswd != nil {
		in, out := &in.HTTPPasswd, &out.HTTPPasswd
		*out = new(plugins.Secret)
		(*in).DeepCopyInto(*out)
	}
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(int32)
		**out = **in
	}
	if in.Headers != nil {
		in, out := &in.Headers, &out.Headers
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.LogResponsePayload != nil {
		in, out := &in.LogResponsePayload, &out.LogResponsePayload
		*out = new(bool)
		**out = **in
	}
	if in.AddLabels != nil {
		in, out := &in.AddLabels, &out.AddLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Workers != nil {
		in, out := &in.Workers, &out.Workers
		*out = new(int32)
		**out = **in
	}
	if in.TLS != nil {
		in, out := &in.TLS, &out.TLS
		*out = new(plugins.TLS)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PrometheusRemoteWrite.
func (in *PrometheusRemoteWrite) DeepCopy() *PrometheusRemoteWrite {
	if in == nil {
		return nil
	}
	out := new(PrometheusRemoteWrite)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Splunk) DeepCopyInto(out *Splunk) {
	*out = *in
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(int32)
		**out = **in
	}
	if in.SplunkToken != nil {
		in, out := &in.SplunkToken, &out.SplunkToken
		*out = new(plugins.Secret)
		(*in).DeepCopyInto(*out)
	}
	if in.HTTPUser != nil {
		in, out := &in.HTTPUser, &out.HTTPUser
		*out = new(plugins.Secret)
		(*in).DeepCopyInto(*out)
	}
	if in.HTTPPasswd != nil {
		in, out := &in.HTTPPasswd, &out.HTTPPasswd
		*out = new(plugins.Secret)
		(*in).DeepCopyInto(*out)
	}
	if in.HTTPDebugBadRequest != nil {
		in, out := &in.HTTPDebugBadRequest, &out.HTTPDebugBadRequest
		*out = new(bool)
		**out = **in
	}
	if in.SplunkSendRaw != nil {
		in, out := &in.SplunkSendRaw, &out.SplunkSendRaw
		*out = new(bool)
		**out = **in
	}
	if in.EventFields != nil {
		in, out := &in.EventFields, &out.EventFields
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Workers != nil {
		in, out := &in.Workers, &out.Workers
		*out = new(int32)
		**out = **in
	}
	if in.TLS != nil {
		in, out := &in.TLS, &out.TLS
		*out = new(plugins.TLS)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Splunk.
func (in *Splunk) DeepCopy() *Splunk {
	if in == nil {
		return nil
	}
	out := new(Splunk)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Stackdriver) DeepCopyInto(out *Stackdriver) {
	*out = *in
	if in.ServiceAccountEmail != nil {
		in, out := &in.ServiceAccountEmail, &out.ServiceAccountEmail
		*out = new(plugins.Secret)
		(*in).DeepCopyInto(*out)
	}
	if in.ServiceAccountSecret != nil {
		in, out := &in.ServiceAccountSecret, &out.ServiceAccountSecret
		*out = new(plugins.Secret)
		(*in).DeepCopyInto(*out)
	}
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.AutoformatStackdriverTrace != nil {
		in, out := &in.AutoformatStackdriverTrace, &out.AutoformatStackdriverTrace
		*out = new(bool)
		**out = **in
	}
	if in.Workers != nil {
		in, out := &in.Workers, &out.Workers
		*out = new(int32)
		**out = **in
	}
	if in.ResourceLabels != nil {
		in, out := &in.ResourceLabels, &out.ResourceLabels
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Stackdriver.
func (in *Stackdriver) DeepCopy() *Stackdriver {
	if in == nil {
		return nil
	}
	out := new(Stackdriver)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Stdout) DeepCopyInto(out *Stdout) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Stdout.
func (in *Stdout) DeepCopy() *Stdout {
	if in == nil {
		return nil
	}
	out := new(Stdout)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Syslog) DeepCopyInto(out *Syslog) {
	*out = *in
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(int32)
		**out = **in
	}
	if in.SyslogMaxSize != nil {
		in, out := &in.SyslogMaxSize, &out.SyslogMaxSize
		*out = new(int32)
		**out = **in
	}
	if in.TLS != nil {
		in, out := &in.TLS, &out.TLS
		*out = new(plugins.TLS)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Syslog.
func (in *Syslog) DeepCopy() *Syslog {
	if in == nil {
		return nil
	}
	out := new(Syslog)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TCP) DeepCopyInto(out *TCP) {
	*out = *in
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(int32)
		**out = **in
	}
	if in.TLS != nil {
		in, out := &in.TLS, &out.TLS
		*out = new(plugins.TLS)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TCP.
func (in *TCP) DeepCopy() *TCP {
	if in == nil {
		return nil
	}
	out := new(TCP)
	in.DeepCopyInto(out)
	return out
}
