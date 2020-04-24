package controllers

import logging "kubesphere.io/fluentbit-operator/api/v1alpha2"

var (
	ownerKey = ".metadata.controller"
	apiGVStr = logging.GroupVersion.String()
)
