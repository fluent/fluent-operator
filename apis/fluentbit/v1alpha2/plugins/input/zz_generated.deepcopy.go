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

package input

import ()

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Dummy) DeepCopyInto(out *Dummy) {
	*out = *in
	if in.Rate != nil {
		in, out := &in.Rate, &out.Rate
		*out = new(int32)
		**out = **in
	}
	if in.Samples != nil {
		in, out := &in.Samples, &out.Samples
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Dummy.
func (in *Dummy) DeepCopy() *Dummy {
	if in == nil {
		return nil
	}
	out := new(Dummy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeExporterMetrics) DeepCopyInto(out *NodeExporterMetrics) {
	*out = *in
	if in.ScrapeInterval != nil {
		in, out := &in.ScrapeInterval, &out.ScrapeInterval
		*out = new(int32)
		**out = **in
	}
	if in.Path != nil {
		in, out := &in.Path, &out.Path
		*out = new(Path)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeExporterMetrics.
func (in *NodeExporterMetrics) DeepCopy() *NodeExporterMetrics {
	if in == nil {
		return nil
	}
	out := new(NodeExporterMetrics)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PrometheusScrapeMetrics) DeepCopyInto(out *PrometheusScrapeMetrics) {
	*out = *in
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PrometheusScrapeMetrics.
func (in *PrometheusScrapeMetrics) DeepCopy() *PrometheusScrapeMetrics {
	if in == nil {
		return nil
	}
	out := new(PrometheusScrapeMetrics)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Systemd) DeepCopyInto(out *Systemd) {
	*out = *in
	if in.SystemdFilter != nil {
		in, out := &in.SystemdFilter, &out.SystemdFilter
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Systemd.
func (in *Systemd) DeepCopy() *Systemd {
	if in == nil {
		return nil
	}
	out := new(Systemd)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Tail) DeepCopyInto(out *Tail) {
	*out = *in
	if in.ReadFromHead != nil {
		in, out := &in.ReadFromHead, &out.ReadFromHead
		*out = new(bool)
		**out = **in
	}
	if in.RefreshIntervalSeconds != nil {
		in, out := &in.RefreshIntervalSeconds, &out.RefreshIntervalSeconds
		*out = new(int64)
		**out = **in
	}
	if in.RotateWaitSeconds != nil {
		in, out := &in.RotateWaitSeconds, &out.RotateWaitSeconds
		*out = new(int64)
		**out = **in
	}
	if in.SkipLongLines != nil {
		in, out := &in.SkipLongLines, &out.SkipLongLines
		*out = new(bool)
		**out = **in
	}
	if in.Multiline != nil {
		in, out := &in.Multiline, &out.Multiline
		*out = new(bool)
		**out = **in
	}
	if in.MultilineFlushSeconds != nil {
		in, out := &in.MultilineFlushSeconds, &out.MultilineFlushSeconds
		*out = new(int64)
		**out = **in
	}
	if in.ParserN != nil {
		in, out := &in.ParserN, &out.ParserN
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.DockerMode != nil {
		in, out := &in.DockerMode, &out.DockerMode
		*out = new(bool)
		**out = **in
	}
	if in.DockerModeFlushSeconds != nil {
		in, out := &in.DockerModeFlushSeconds, &out.DockerModeFlushSeconds
		*out = new(int64)
		**out = **in
	}
	if in.DisableInotifyWatcher != nil {
		in, out := &in.DisableInotifyWatcher, &out.DisableInotifyWatcher
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Tail.
func (in *Tail) DeepCopy() *Tail {
	if in == nil {
		return nil
	}
	out := new(Tail)
	in.DeepCopyInto(out)
	return out
}
