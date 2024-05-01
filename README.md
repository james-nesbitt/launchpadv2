# Mirantis Launchpad

A tool for installing Mirantis Containers products.

Originally developped as an SSH/WinRM tool, now adapted with a more abstract
implementation to handle a more diverse set of products.

# Building 

Try the goreleaser based Make targets to build the tool:

Builds release targets:

``` make clean dist```

Build a single binary for your platform:

```make local```

## Docs

See ./docs/README.md for more.
