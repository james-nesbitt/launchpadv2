//go:build integration
package integration

var (
	IntegrationTestYaml_Legacy_2_0 = `---
apiVersion: launchpad.mirantis.com/v2.0
kind: mke
metadata:
  name: lp2.0-leg
spec:
  hosts:
    ACon0:
      handler: mock
      roles:
      - manager
    AWrk0:
      handler: mock
      roles:
      - worker
    AWrk1:
      handler: mock
      roles:
      - worker
  products:
    mcr:
      version: 23.0.7
      repoURL: https://repos.mirantis.com
      installURLLinux: https://get.mirantis.com/
      installURLWindows: https://get.mirantis.com/install.ps1
      channel: stable
      prune: true
    mke3:
      version: 3.7.3
      imageRepo: docker.io/mirantis
      adminUsername: admin
      installFlags:
      - "--default-node-orchestrator=kubernetes"
      - "--nodeport-range=32768-35535"
      upgradeFlags:
      - "--force-recent-backup"
      - "--force-minimums"
    msr2:
      version: 2.9.30
`
	IntegrationTestYaml_Legacy_2_1 = `---
apiVersion: launchpad.mirantis.com/v2.1
kind: mke
metadata:
  name: lp2.1-leg
spec:
  components:
    hosts:
      ACon0:
        handler: mock
        roles:
        - manager
      AWrk0:
        handler: mock
        roles:
        - worker
      AWrk1:
        handler: mock
        roles:
        - worker
    mcr:
      version: 23.0.7
      repoURL: https://repos.mirantis.com
      installURLLinux: https://get.mirantis.com/
      installURLWindows: https://get.mirantis.com/install.ps1
      channel: stable
      prune: true
    mke3:
      version: 3.7.3
      imageRepo: docker.io/mirantis
      adminUsername: admin
      installFlags:
      - "--default-node-orchestrator=kubernetes"
      - "--nodeport-range=32768-35535"
      upgradeFlags:
      - "--force-recent-backup"
      - "--force-minimums"
    msr2:
      version: 2.9.30
`
	IntegrationTestYaml_NextGen_2_0 = `---
apiVersion: launchpad.mirantis.com/v2.0
kind: mke
metadata:
  name: lp2.1-ng
spec:
  hosts:
    ACon0:
      handler: mock
      roles:
      - controller
    AWrk0:
      handler: mock
      roles:
      - worker
    AWrk1:
      handler: mock
      roles:
      - worker
  products
    k0s:
      version: 1.28-k0s
    mke4:
      version: 4.0.0
    msr4:
      version: 4.0.0
`
	IntegrationTestYaml_NextGen_2_1 = `---
apiVersion: launchpad.mirantis.com/v2.1
kind: mke
metadata:
  name: lp2.1-ng
spec:
  components:
    hosts:
      ACon0:
        handler: mock
        roles:
        - controller
      AWrk0:
        handler: mock
        roles:
        - worker
      AWrk1:
        handler: mock
        roles:
        - worker
    k0s:
      version: 1.28-k0s
    mke4:
      version: 4.0.0
    msr4:
      version: 4.0.0
`

	// used for overriding with custom value for integration tests
	IntegrationTestYamlDoNotCommit = ``

)
