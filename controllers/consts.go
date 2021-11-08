package controllers

import (
	loggingv1alpha2 "kubesphere.io/fluentbit-operator/apis/kubesphere.io/v1alpha2"
)

var (
	ownerKey = ".metadata.controller"
	apiGVStr = loggingv1alpha2.SchemeGroupVersion.String()
)
