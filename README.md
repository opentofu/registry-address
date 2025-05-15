# OpenTofu Registry Addresses

This Go library contains types to represent module and provider registry
addresses as used in OpenTofu, along with the canonical implementations of
parsing and comparing those addresses.

**Provider** addresses can be found in

 - [`tofu show -json <FILE>`](https://opentofu.org/docs/internals/json-format/#configuration-representation) (`full_name`)
 - [`tofu version -json`](https://opentofu.org/docs/cli/commands/version/) (`provider_selections`)
 - [`tofu providers schema -json`](https://opentofu.org/docs/cli/commands/providers/schema/#providers-schema-representation) (keys of `provider_schemas`)
 - within `required_providers` block in OpenTofu modules
 - OpenTofu [CLI configuration file](https://opentofu.org/docs/cli/config/config-file/#provider-installation)
 - Plugin [reattach configurations](https://www.terraform.io/plugin/debugging#running-terraform-with-a-provider-in-debug-mode)

**Module** addresses can be found within `source` argument
of `module` blocks in OpenTofu modules (`*.tf`)
and parts of the address (namespace and name) in the Module Registry API.

## Compatibility

The module aims for compatibility with OpenTofu v1.5 and later.

## Usage

### Provider

```go
pAddr, err := regaddr.ParseProviderSource("foo/bar")
if err != nil {
	// deal with error
}

// pAddr == regaddr.Provider{
//   Type:      "foo",
//   Namespace: "bar",
//   Hostname:  regaddr.DefaultProviderRegistryHost,
// }
```

### Module

```go
mAddr, err := regaddr.ParseModuleSource("foo/bar/baz//modules/example")
if err != nil {
	// deal with error
}

// mAddr == Module{
//   Package: ModulePackage{
//     Host:         regaddr.DefaultModuleRegistryHost,
//     Namespace:    "foo",
//     Name:         "bar",
//     TargetSystem: "baz",
//   },
//   Subdir: "modules/example",
// },
```

## Other Module Address Formats

OpenTofu supports [various other module source address types](https://opentofu.org/docs/language/modules/sources/)
that are not handled by this library. This library focuses only on the
syntax used for OpenTofu module registries.
