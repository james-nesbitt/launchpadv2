package host

import (
	"context"
	"fmt"
)

// ClusterHostsRoleRequirementStandard a requirement for hosts with matching roles
type ClusterHostsRoleRequirementStandard struct {
	Description string
	// Role of hosts needed
	Role string
	// If 0 then no limit - but describe minimum needed hosts, and maximum consumed hosts
	MinimumCount, MaximumCount int

	// Hosts (the type) registered as a source that will be used when the requirement id fullfilled
	hostsSource Hosts
}

func (chr ClusterHostsRoleRequirementStandard) Describe() string {
	return chr.Description //fmt.Sprintf("%s: %+v",chr.Description, chr)
}

func (chr ClusterHostsRoleRequirementStandard) Validate(_ context.Context) error {
	_, err := chr.clusterHostsRoleRequirementFullfill()
	return err
}

func (chr *ClusterHostsRoleRequirementStandard) AssignSourceHosts(hs Hosts) {
	chr.hostsSource = hs
}

func (chr ClusterHostsRoleRequirementStandard) RetrieveMatchingHosts() (Hosts, error) {
	return chr.clusterHostsRoleRequirementFullfill()
}

func (chr ClusterHostsRoleRequirementStandard) clusterHostsRoleRequirementFullfill() (Hosts, error) {
	rhs := Hosts{}

	for _, hs := range chr.hostsSource {
		if chr.MaximumCount > 0 && len(rhs) > chr.MaximumCount {
			break
		}

		if chr.Role != "" && !hs.HasRole(chr.Role) {
			continue
		}

		rhs = append(rhs, hs)
	}

	if chr.MinimumCount > 0 && len(rhs) < chr.MinimumCount {
		return rhs, fmt.Errorf("%w; not enough matching hosts: %s : %+v", ErrRoleHostsDependencyUnMet, chr.Describe(), chr.hostsSource)
	}

	return rhs, nil
}
