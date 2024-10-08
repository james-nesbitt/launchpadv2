# TODO

## Logging

- [ ] (log-sanitize) Sanitize output to prevent logging of sensitive information. Probably components should be able to add sensitive values to a list as they run

## cli 

- [X] (allow-empty-project) we should be able to run some cli commands without a project/config, such as version.

## Action 

- [X] (phase-overwrite-warn) Phase accessors should probably warn when an existing phase is dropped because a new one with the same name is added.
- [ ] (phase-order-debug) The Phase sort graph library only returns a boolean on failure which means we cannot give the user details about any failure (look for a new library?)
- [ ] (allow-dep-match-fail) Allow dependency matching to fail, to allow separation of actions for things like the cli. This places the responsibility of Dependency matching on the user

## Dependency 

- [ ] (unique-event-ids) Events and Event keys need to be unique, but the code is very reled about enforcing that. We should probably make the Unique part immutable and unique.
- [ ] (dependency-selection-direction) Design a system that allows components to be more selective about which depenendcy they use to fullfill their needs. This matters more when there is more that one dependency fulfiller for a requirement.

## Project 

- [X] (rename) Rename "cluster" to "project"
- [ ] (allow-dep-match-fail) Allow dependency matching to fail, to allow separation of actions for things like the cli. This places the responsibility of Dependency matching on the user

## Host 

- [x] (host-network-discover) Hosts need to be able to discover their networking (internal) so that we can activate swarm using the advertise address
- [x] (flexible-host) Allow Hosts to have component specific configuration/functionalily so that we can configure host specific things like docker-sudo, mcr-daemon-config. Maybe a plugin system, or allow multiple handlers per host? Currently I am looking at having a system for allowing hosts to be decorated on demand
- [ ] (host-hooks) Allow host hooks. Perhaps rely on action.Events or action.Phases as markers.
- [x] (host-as-a-component) should we stop treating host components separately from product components.
- [ ] (log-command-identifier) Log entries for a command should have some identifier to show that they are connected so that you can correlate the various output lines to one execution.

## Component

- [ ] (component-state) develop a better state strategy/pattern for components. Requirements would include: locking, debug/output, config comparison, io/storage?

## Docker:DockerExec

- [ ] (dockerexec-cache) we call docker info and swarm inspect repeatedly for leaders. This host should likely cache this information. The whole dockerexec client could use a caching mechanism.
- [ ] (docker-exec-injection) confirm that we are not allowing shell injection with user data being used in docker commands

## Products 

- [x] (product-discover-strategy) develop a strategy for the Discover operations that is more clear about state, and allows for expecting to fail (like checking that a product is uninstalled after running uninstall.

## Product:MKE3

- [X] (mke3-config-separate) separate install and upgrade config so that it is obvious to the user what values are supported for each operation.
- [ ] (mke3-implementation) MKE3 implementation is needed to allow the client-bundle download

## Product:K0s

- [X] (k0s-initial) start K0s component 
- [x] (k0s-multiple-controller) currently, adding controllers fails with the nextgen terraform aws example
- [x] (k0s-kubernetes-implementation) k0s can produce a kubernetes implementation
- [ ] (k0s-project) K0s project operations should install/uninstall k0s
- [ ] (k0s-config-dig) K0sConfig should just be a dig, but it wasn't working well (ignorance more than technical issues)

## Product:MKEx 

- [X] (mkex-initial) start mkex component 
- [X] (mkex-swarm-initialize) get mkex to activate swarm

## 
