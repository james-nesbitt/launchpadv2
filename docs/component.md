# Component 

Components in Launchpad are functional units which provide funcitonality that 
can be used by a cluster to build commands, provide cli options, and assist other 
components by fullfilling dependencies.

## Building a Component

A component is any struct that implements the component interface, which is a small
interface that requires only that it can identify itself.

Components provide functionaly by mutation in that they implement other interfaces 
which allow them to provide functionality.

## Registering Components

Any Component instance can be added to a project.Project object, and that project
will actively use the Component for building action.Commands.

Typically, Components are defined and configured using the config process which 
decodes Components defined in configuration. This usually requires that the 
Component has a decoder which is registered with the config builder used.
An example, is that there is a public global array of ProductDecoder functions 
registered in product, which are typically populated in Component init() methods.

## Component interfaces

### Commands 

If the Component implements the action.CommandBuild interface, then Launchpad will 
know to ask the Component to build execution Phases during command runs.

The Component is expected to recognize what Command will be run, based on a string
key, and should add Phases to the Command as needed to procedd.

#### Commands and Dependencies

Any Command Phase added can include dependency.Dependency action.Events that are 
needed to run, and Launchpad will order the action.Phases in order to ensure 
that dependencies are met.

### Dependency interfaces

Components can require dependencies, and can fullfill dependencies.

Launchpad will collect all dependency.Requirements, and try to match them with 
dependency.Dependencies during Validation, to make sure that all needed 
functionality is available.

#### dependency.RequiresDependencies

This Component has dependencies that need to be met by another Component.

The Component should be ready to provide dependency.Requirement objects of the
appropriate type for its needs (e.g. A host.HostRolesRequuirement to indicate
that hosts of a certain role are needed.)


#### dependency.ProvidesDependencies: 

This Component can fullfill dependency.Requirements from other Components.

The Component will be handed requirements, and if the Component can fulfill 
the Requirement, it should return an appropriate dependency.Dependency object.

The Dependency could be called upon at any time, and so should be able to 
be called upon repeatedly.

#### Commands and Dependencies

If dependency fullfillment requires a Command phase, then the Component should 
add the Phase to any Command calls, and Dependency Events for if the Dependency 
will be made available, or removed during the Command execution

## Implementations instead of Components 

Implementations are also a way of providing functionality, but are never used 
as Components.

Implementations provided functionality by defining dependency.Requirement and 
matching dependency.Dependency types which other components can use.

An example is "a kubernetes client", which is a standard, but can be provided 
by either MKE3, or K0S.
