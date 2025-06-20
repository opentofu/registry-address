// Copyright (c) The OpenTofu Authors
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package regaddr

import (
	"fmt"
	"log"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/opentofu/svchost"
)

func TestProviderString(t *testing.T) {
	tests := []struct {
		Input Provider
		Want  string
	}{
		{
			Provider{
				Type:      "test",
				Hostname:  DefaultProviderRegistryHost,
				Namespace: "hashicorp",
			},
			NewProvider(DefaultProviderRegistryHost, "hashicorp", "test").String(),
		},
		{
			Provider{
				Type:      "test-beta",
				Hostname:  DefaultProviderRegistryHost,
				Namespace: "hashicorp",
			},
			NewProvider(DefaultProviderRegistryHost, "hashicorp", "test-beta").String(),
		},
		{
			Provider{
				Type:      "test",
				Hostname:  "registry.opentofu.com",
				Namespace: "hashicorp",
			},
			"registry.opentofu.com/hashicorp/test",
		},
		{
			Provider{
				Type:      "test",
				Hostname:  DefaultProviderRegistryHost,
				Namespace: "othercorp",
			},
			DefaultProviderRegistryHost.ForDisplay() + "/othercorp/test",
		},
	}

	for _, test := range tests {
		got := test.Input.String()
		if got != test.Want {
			t.Errorf("wrong result for %s\n", test.Input.String())
		}
	}
}

func TestProviderLegacyString(t *testing.T) {
	tests := []struct {
		Input Provider
		Want  string
	}{
		{
			Provider{
				Type:      "test",
				Hostname:  DefaultProviderRegistryHost,
				Namespace: LegacyProviderNamespace,
			},
			"test",
		},
		{
			Provider{
				Type:      "opentf",
				Hostname:  TransitionalBuiltInProviderHost,
				Namespace: BuiltInProviderNamespace,
			},
			"opentf",
		},
	}

	for _, test := range tests {
		got := test.Input.LegacyString()
		if got != test.Want {
			t.Errorf("wrong result for %s\ngot:  %s\nwant: %s", test.Input.String(), got, test.Want)
		}
	}
}

func TestProviderDisplay(t *testing.T) {
	tests := []struct {
		Input Provider
		Want  string
	}{
		{
			Provider{
				Type:      "test",
				Hostname:  DefaultProviderRegistryHost,
				Namespace: "hashicorp",
			},
			"hashicorp/test",
		},
		{
			Provider{
				Type:      "test",
				Hostname:  "registry.opentofu.com",
				Namespace: "hashicorp",
			},
			"registry.opentofu.com/hashicorp/test",
		},
		{
			Provider{
				Type:      "test",
				Hostname:  DefaultProviderRegistryHost,
				Namespace: "othercorp",
			},
			"othercorp/test",
		},
		{
			Provider{
				Type:      "terraform",
				Namespace: BuiltInProviderNamespace,
				Hostname:  TransitionalBuiltInProviderHost,
			},
			"terraform.io/builtin/terraform",
		},
	}

	for _, test := range tests {
		got := test.Input.ForDisplay()
		if got != test.Want {
			t.Errorf("wrong result for %s: %q\n", test.Input.String(), got)
		}
	}
}

func TestProviderIsBuiltIn(t *testing.T) {
	tests := []struct {
		Input Provider
		Want  bool
	}{
		{
			Provider{
				Type:      "test",
				Hostname:  TransitionalBuiltInProviderHost,
				Namespace: BuiltInProviderNamespace,
			},
			true,
		},
		{
			Provider{
				Type:      "opentf",
				Hostname:  TransitionalBuiltInProviderHost,
				Namespace: BuiltInProviderNamespace,
			},
			true,
		},
		{
			Provider{
				Type:      "test",
				Hostname:  TransitionalBuiltInProviderHost,
				Namespace: "boop",
			},
			false,
		},
		{
			Provider{
				Type:      "test",
				Hostname:  DefaultProviderRegistryHost,
				Namespace: BuiltInProviderNamespace,
			},
			false,
		},
		{
			Provider{
				Type:      "test",
				Hostname:  DefaultProviderRegistryHost,
				Namespace: "hashicorp",
			},
			false,
		},
		{
			Provider{
				Type:      "test",
				Hostname:  "registry.opentofu.com",
				Namespace: "hashicorp",
			},
			false,
		},
		{
			Provider{
				Type:      "test",
				Hostname:  DefaultProviderRegistryHost,
				Namespace: "othercorp",
			},
			false,
		},
	}

	for _, test := range tests {
		got := test.Input.IsBuiltIn()
		if got != test.Want {
			t.Errorf("wrong result for %s\ngot:  %#v\nwant: %#v", test.Input.String(), got, test.Want)
		}
	}
}

func TestProviderIsLegacy(t *testing.T) {
	tests := []struct {
		Input Provider
		Want  bool
	}{
		{
			Provider{
				Type:      "test",
				Hostname:  DefaultProviderRegistryHost,
				Namespace: LegacyProviderNamespace,
			},
			true,
		},
		{
			Provider{
				Type:      "test",
				Hostname:  "registry.opentofu.com",
				Namespace: LegacyProviderNamespace,
			},
			false,
		},
		{
			Provider{
				Type:      "test",
				Hostname:  DefaultProviderRegistryHost,
				Namespace: "hashicorp",
			},
			false,
		},
	}

	for _, test := range tests {
		got := test.Input.IsLegacy()
		if got != test.Want {
			t.Errorf("wrong result for %s\n", test.Input.String())
		}
	}
}

func ExampleParseProviderSource() {
	pAddr, err := ParseProviderSource("hashicorp/aws")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v", pAddr)
	// Output: regaddr.Provider{Type:"aws", Namespace:"hashicorp", Hostname:svchost.Hostname("registry.opentofu.org")}
}

func TestParseProviderSource(t *testing.T) {
	tests := map[string]struct {
		Want Provider
		Err  bool
	}{
		"registry.opentofu.org/hashicorp/aws": {
			Provider{
				Type:      "aws",
				Namespace: "hashicorp",
				Hostname:  DefaultProviderRegistryHost,
			},
			false,
		},
		"registry.opentofu.org/HashiCorp/AWS": {
			Provider{
				Type:      "aws",
				Namespace: "hashicorp",
				Hostname:  DefaultProviderRegistryHost,
			},
			false,
		},
		"opentofu.org/builtin/opentofu": {
			Provider{
				Type:      "opentofu",
				Namespace: BuiltInProviderNamespace,
				Hostname:  BuiltInProviderHost,
			},
			false,
		},
		"terraform.io/builtin/terraform": {
			Provider{
				Type:      "terraform",
				Namespace: BuiltInProviderNamespace,
				Hostname:  TransitionalBuiltInProviderHost,
			},
			false,
		},
		// v0.12 representation
		// In most cases this would *likely* be the same provider
		// we otherwise represent as being in the default namespace,
		// but we cannot be sure in the context of the source string alone.
		"terraform": {
			Provider{
				Type:      "terraform",
				Namespace: UnknownProviderNamespace,
				Hostname:  DefaultProviderRegistryHost,
			},
			false,
		},
		"hashicorp/aws": {
			Provider{
				Type:      "aws",
				Namespace: "hashicorp",
				Hostname:  DefaultProviderRegistryHost,
			},
			false,
		},
		"HashiCorp/AWS": {
			Provider{
				Type:      "aws",
				Namespace: "hashicorp",
				Hostname:  DefaultProviderRegistryHost,
			},
			false,
		},
		"aws": {
			Provider{
				Type:      "aws",
				Namespace: UnknownProviderNamespace,
				Hostname:  DefaultProviderRegistryHost,
			},
			false,
		},
		"AWS": {
			Provider{
				Type:      "aws",
				Namespace: UnknownProviderNamespace,
				Hostname:  DefaultProviderRegistryHost,
			},
			false,
		},
		"example.com/foo-bar/baz-boop": {
			Provider{
				Type:      "baz-boop",
				Namespace: "foo-bar",
				Hostname:  svchost.Hostname("example.com"),
			},
			false,
		},
		"foo-bar/baz-boop": {
			Provider{
				Type:      "baz-boop",
				Namespace: "foo-bar",
				Hostname:  DefaultProviderRegistryHost,
			},
			false,
		},
		"localhost:8080/foo/bar": {
			Provider{
				Type:      "bar",
				Namespace: "foo",
				Hostname:  svchost.Hostname("localhost:8080"),
			},
			false,
		},
		"example.com/too/many/parts/here": {
			Provider{},
			true,
		},
		"/too///many//slashes": {
			Provider{},
			true,
		},
		"///": {
			Provider{},
			true,
		},
		"/ / /": { // empty strings
			Provider{},
			true,
		},
		"badhost!/hashicorp/aws": {
			Provider{},
			true,
		},
		"example.com/badnamespace!/aws": {
			Provider{},
			true,
		},
		"example.com/bad--namespace/aws": {
			Provider{},
			true,
		},
		"example.com/-badnamespace/aws": {
			Provider{},
			true,
		},
		"example.com/badnamespace-/aws": {
			Provider{},
			true,
		},
		"example.com/bad.namespace/aws": {
			Provider{},
			true,
		},
		"example.com/hashicorp/badtype!": {
			Provider{},
			true,
		},
		"example.com/hashicorp/bad--type": {
			Provider{},
			true,
		},
		"example.com/hashicorp/-badtype": {
			Provider{},
			true,
		},
		"example.com/hashicorp/badtype-": {
			Provider{},
			true,
		},
		"example.com/hashicorp/bad.type": {
			Provider{},
			true,
		},

		// We forbid the opentofu- prefix both because it's redundant to
		// include "opentofu" in a provider name and because we use
		// the longer prefix opentofu-provider- to hint for users who might be
		// accidentally using the git repository name or executable file name
		// instead of the provider type. We also continue our predecessor
		// project's tradition of forbidding similar prefixes for its own
		// project name.
		"example.com/opentofu/opentofu-provider-bad": {
			Provider{},
			true,
		},
		"example.com/opentofu/opentofu-bad": {
			Provider{},
			true,
		},
		"example.com/opentofu/terraform-provider-bad": {
			Provider{},
			true,
		},
		"example.com/opentofu/terraform-bad": {
			Provider{},
			true,
		},
	}

	for name, test := range tests {
		got, err := ParseProviderSource(name)
		if diff := cmp.Diff(test.Want, got); diff != "" {
			t.Errorf("mismatch (%q): %s", name, diff)
		}
		if err != nil {
			if test.Err == false {
				t.Errorf("got error: %s, expected success", err)
			}
		} else {
			if test.Err {
				t.Errorf("got success, expected error")
			}
		}
	}
}

func TestParseProviderPart(t *testing.T) {
	tests := map[string]struct {
		Want  string
		Error string
	}{
		`foo`: {
			`foo`,
			``,
		},
		`FOO`: {
			`foo`,
			``,
		},
		`Foo`: {
			`foo`,
			``,
		},
		`abc-123`: {
			`abc-123`,
			``,
		},
		`Испытание`: {
			`испытание`,
			``,
		},
		`münchen`: { // this is a precomposed u with diaeresis
			`münchen`, // this is a precomposed u with diaeresis
			``,
		},
		`münchen`: { // this is a separate u and combining diaeresis
			`münchen`, // this is a precomposed u with diaeresis
			``,
		},
		`abc--123`: {
			``,
			`cannot use multiple consecutive dashes`,
		},
		`xn--80akhbyknj4f`: { // this is the punycode form of "испытание", but we don't accept punycode here
			``,
			`cannot use multiple consecutive dashes`,
		},
		`abc.123`: {
			``,
			`dots are not allowed`,
		},
		`-abc123`: {
			``,
			`must contain only letters, digits, and dashes, and may not use leading or trailing dashes`,
		},
		`abc123-`: {
			``,
			`must contain only letters, digits, and dashes, and may not use leading or trailing dashes`,
		},
		``: {
			``,
			`must have at least one character`,
		},
	}

	for given, test := range tests {
		t.Run(given, func(t *testing.T) {
			got, err := ParseProviderPart(given)
			if test.Error != "" {
				if err == nil {
					t.Errorf("unexpected success\ngot:  %s\nwant: %s", err, test.Error)
				} else if got := err.Error(); got != test.Error {
					t.Errorf("wrong error\ngot:  %s\nwant: %s", got, test.Error)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error\ngot:  %s\nwant: <nil>", err)
				} else if got != test.Want {
					t.Errorf("wrong result\ngot:  %s\nwant: %s", got, test.Want)
				}
			}
		})
	}
}

func TestProviderEquals(t *testing.T) {
	tests := []struct {
		InputP Provider
		OtherP Provider
		Want   bool
	}{
		{
			NewProvider(DefaultProviderRegistryHost, "foo", "test"),
			NewProvider(DefaultProviderRegistryHost, "foo", "test"),
			true,
		},
		{
			NewProvider(DefaultProviderRegistryHost, "foo", "test"),
			NewProvider(DefaultProviderRegistryHost, "bar", "test"),
			false,
		},
		{
			NewProvider(DefaultProviderRegistryHost, "foo", "test"),
			NewProvider(DefaultProviderRegistryHost, "foo", "my-test"),
			false,
		},
		{
			NewProvider(DefaultProviderRegistryHost, "foo", "test"),
			NewProvider("example.com", "foo", "test"),
			false,
		},
	}
	for _, test := range tests {
		t.Run(test.InputP.String(), func(t *testing.T) {
			got := test.InputP.Equals(test.OtherP)
			if got != test.Want {
				t.Errorf("wrong result\ngot:  %v\nwant: %v", got, test.Want)
			}
		})
	}
}
