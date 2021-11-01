package controllers

import (
	clusterfluentbit "kubesphere.io/fluentbit-operator/api/clusterfluentbitoperator/v1alpha2"
	logging "kubesphere.io/fluentbit-operator/api/fluentbitoperator/v1alpha2"
)

var (
	ownerKey = ".metadata.controller"
	apiGVStr = logging.SchemeGroupVersion.String()
	clusterapiGVStr = clusterfluentbit.GroupVersion.String()
)
