package controllers

import (
	fluentbitv1alpha2 "github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2"
	fluentdv1alpha1 "github.com/fluent/fluent-operator/v2/apis/fluentd/v1alpha1"
)

var (
	fluentbitOwnerKey = ".fluentbit.metadata.controller"
	collectorOwnerKey = ".collector.metadata.controller"
	fluentdOwnerKey   = ".fluentd.metadata.controller"
	fluentbitApiGVStr = fluentbitv1alpha2.SchemeGroupVersion.String()
	fluentdApiGVStr   = fluentdv1alpha1.SchemeGroupVersion.String()
	fluentdAgentMode  = "agent"
)
