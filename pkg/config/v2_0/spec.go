package v20

/**
* V2.0 Launchpad Spec
*
* This spec is similar to the v1 spec, which had a section for hosts
* followed by a number of product configurations. In v2 we treat the
* products as a section sitting next to the hosts, and both the hosts
* and the products are keyed maps.
*
* ------------------------------
* hosts:
*   one:
*     roles:
*     - manager
*   two:
*     handler: rig
*     roles:
*     - bastion
*     ssh:
*       address: 127.0.0.1
*       key_path: ./key.pem
*       user: sshuser
*
* products:
*   k0s:
*     version: 1.24
*   registry:
*     handler: msr3
*     version: 3.1.1
* ------------------------------
*
* NOTES:
*  - hosts are turned into a HostComponent in the cluster
*
 */

// Spec defines cluster spec.
type Spec struct {
	Hosts    SpecHosts    `yaml:"hosts" validate:"required,min=1"`
	Products SpecProducts `yaml:"products" validate:"required,min=1"`
}
