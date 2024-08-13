package kubernetes

import (
	kubeclientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type Config struct {
	KubeClientConfig kubeclientcmdapi.Config
}
