package k0s

import "github.com/k0sproject/dig"

// K0sConfig is an untyped passthrough for k0s cluster configuration.
// Using dig.Mapping (map[string]interface{}) allows forward-compatibility:
// any k0s config fields pass through without requiring struct updates.
//
// Reference: https://docs.k0sproject.io/head/configuration/
type K0sConfig = dig.Mapping
