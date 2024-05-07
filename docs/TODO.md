# TODO

## Action 

- [X] (phase-overwrite-warn) Phase accessors should probably warn when an existing phase is dropped because a new one with the same name is added.
- [ ] (phase-order-debug) The Phase sort graph library only returns a boolean on failure which means we cannot give the user details about any failure (look for a new library?)

## Dependency 

- [ ] (unique-event-ids) Events and Event keys need to be unique, but the code is very reled about enforcing that. We should probably make the Unique part immutable and unique.
