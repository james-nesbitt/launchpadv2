# Mirantis Launchpad Design

## Common terms 

### Component

a Launchpad functional module which can participate in Launchpad by:

1. Participating in building commands by adding phases (and/or steps)
2. Provide/Require dependencies to/from other components
3. Provide Implementation specific code

Components are core functionality, collected into a project and used 
for building all other functionality.

Components work by implementing various Launchpad interfaces, indicating
what funcitonality they offer.

### Implementation 

Deliver known functionality to other code, such as APIs or interfaces.

Kept distinct from Components in that two Components may deliver that same
functionality, and so they can share a code base.

example: docker-swarm, kubernetes, the MKE api

### Command 

A single executable operation, built by collecting Phases from various 
providers (typically Components)

A Phase is an isolated set of instructions, typically defined by a 
component. Phases typically advertise their dependencies so that they
can be ordered.

### Dependency / Requirement 

A Requirement is in indication that a Dependency is Required, and the
Dependency is a fullfillment of that requirement.

This allows code to advertise what it needs, without having to know 
how the functionality is provided.  For example, a Component may rely
on Kubernetes, so it advertises what Kubernetes it needs by providing a
Requirement; then another component will see the Requirement and fullfill
it if it can.  Kubernetes could be provided by another Component which 
discovers a Kubernetes project, or by one who will create a Kubernetes 
project.

### Project 

A set of Components, from which Dependencies are pulled, and Commands 
are built.

