package plugins

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-openapi/errors"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ConfigMapLoader struct {
	client    client.Client
	namespace string
}

func NewConfigMapLoader(c client.Client, ns string) ConfigMapLoader {
	return ConfigMapLoader{
		client:    c,
		namespace: ns,
	}
}

func (cl ConfigMapLoader) LoadConfigMap(selector v1.ConfigMapKeySelector, namespace string) (string, error) {
	var configMap v1.ConfigMap
	ctx := context.Background()
	if err := cl.client.Get(ctx, client.ObjectKey{Name: selector.Name, Namespace: namespace}, &configMap); err != nil {
		return "", err
	}

	if v, ok := configMap.Data[selector.Key]; !ok {
		return "", errors.NotFound(fmt.Sprintf("The key %s is not found.", selector.Key))
	} else {
		return strings.TrimSuffix(v, "\n"), nil
	}
}
