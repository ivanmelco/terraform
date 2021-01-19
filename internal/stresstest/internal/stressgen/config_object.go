package stressgen

import (
	"github.com/hashicorp/hcl/v2/hclwrite"

	"github.com/hashicorp/terraform/states"
)

// ConfigObject is an interface implemented by types representing items that
// can be included in a generated test configuration.
//
// Each ConfigObject typically represents one configuration block in a module,
// and has a few different responsibilities. The most important is to generate
// the actual configuration block for the object, but a ConfigObject can be
// made more useful by providing a verifier that checks whether the final
// state matches the goal of the configuration, and by registering objects that
// it makes available in the symbol table that later-constructed objects might
// potentially refer to in order to create a proper dependency graph that is
// more likely to detect race conditions.
//
// Some ConfigObject implementations can contain other nested ConfigObject
// implementations. For example, ModuleCall contains another whole module
// which its parent module will call.
type ConfigObject interface {
	// DisplayName returns a reasonable identifier for this object to use in
	// UI output. It's not necessarily unique across a whole configuration, but
	// should be as unique as possible. For objects that already have
	// conventional absolute address syntax, like resources and modules,
	// the string serialization of those addresses are a good choice.
	DisplayName() string

	// AppendConfig appends one or more configuration blocks to the given
	// body, which represents the top-level body of a .tf file.
	AppendConfig(to *hclwrite.Body)

	// CheckState compares the relevant parts of the given state to the
	// original configuration for itself and returns one or more errors if
	// anything doesn't match expectations.
	CheckState(in *states.State) []error
}
