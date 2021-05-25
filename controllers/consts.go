package controllers

import logging "kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2"

var (
	ownerKey = ".metadata.controller"
	apiGVStr = logging.SchemeGroupVersion.String()
)
