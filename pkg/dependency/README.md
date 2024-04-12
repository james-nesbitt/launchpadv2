# Dependency management

This package provides the interfaces and some helper functionality
for managing interdependency between components, and their derivatives
such as phases and actions.

## Requirements versus Dependencies 

### Requirements 

Requirements are a mechanism for a component to describe functionality that it needs 
in the form af an identifiable struct.

### Dependency 

A Dependency is a fullfiller of one or more requirements.

## The Mechanisms

Any component which requires functionality should advertise this by preoviding a set
of Requirement objects.  Each Requirement must then be matched with a dependency by
another component.
When a component acceptts a Requirement, it produces a Dependency. The Requirer can 
then rely on the Dependency to produce the functionality it needs.  The Component that
produces the Dependency can use one Dependency to fullfill multiple Requirements.

The Requirer can convert the Dependency to an interface that it expects to deliver 
the expected functionality.

## The Process 

1. Component A declares it needs functionality X by providing Requirement A-Rx
2. Component B declares it also need functionality X by providing Requiremnt B-Rx
3. Component F received Requirement A-Rx and decides that it can fullfill it, so it
   returns Dependency Dx
4. Component F also receives Requirement B-Rx which it can also fullfill, so it
   returns again Dependency Dx, or returns a new Dx_2
5. Dependency Dx can be asked if it is Valid: "Do you think that you can provide 
   functionality X"
6. Component A asks Dependency Dx to fullfill functionality X by first converting
   it to an expected Interface 
   (if you know what functionality you need then you should know what interface you want)
7. Dependency Dx can be asked is functionality X is ready to be fullfilled

Now a 3rd party can connect all dependencies, and order them by coordinating the 
collection of requirements, and remembering who agreed to fullfill them. The 
Dependency objects can be used to determine if a system of deendencies can be met 
(if a configuration is missing something) but also other systems can identify 
that they deliver for ordering.
