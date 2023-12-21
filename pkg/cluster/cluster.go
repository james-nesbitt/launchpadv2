package cluster

import (
	"github.com/Mirantis/launchpad/pkg/host"
	"github.com/Mirantis/launchpad/pkg/product"
)

// Cluster function handler for the complete cluster
type Cluster struct {
	Hosts    host.Hosts
	Products product.Products
}
