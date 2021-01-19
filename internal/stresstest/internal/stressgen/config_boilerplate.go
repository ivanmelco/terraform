package stressgen

import (
	"sort"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/states"
	"github.com/zclconf/go-cty/cty"
)

// ConfigBoilerplate is an implementation of ConfigObject representing some
// items we need to include in our test configurations to create the basis
// of a valid Terraform module.
//
// This particular type is not part of our random generation process, and is
// instead generally instantiated directly as the first object in a module,
// before we start appending randomly-generated items.
type ConfigBoilerplate struct {
	Providers map[string]addrs.Provider
}

var _ ConfigObject = (*ConfigBoilerplate)(nil)

// DisplayName implements ConfigObject.DisplayName.
func (bp *ConfigBoilerplate) DisplayName() string {
	return "module boilerplate"
}

// AppendConfig implements ConfigObject.AppendConfig.
func (bp *ConfigBoilerplate) AppendConfig(to *hclwrite.Body) {
	tfBlock := to.AppendBlock(hclwrite.NewBlock("terraform", nil))
	tfBody := tfBlock.Body()

	providerReqsBlock := tfBody.AppendBlock(hclwrite.NewBlock("provider_requirements", nil))
	providerReqsBody := providerReqsBlock.Body()
	var providerLocalNames = make([]string, 0, len(bp.Providers))
	for localName := range bp.Providers {
		providerLocalNames = append(providerLocalNames, localName)
	}
	sort.Strings(providerLocalNames)
	for _, localName := range providerLocalNames {
		addr := bp.Providers[localName]
		providerReqsBody.SetAttributeValue(localName, cty.ObjectVal(map[string]cty.Value{
			"source": cty.StringVal(addr.String()),
		}))
	}
}

// CheckState implements ConfigObject.CheckState.
func (bp *ConfigBoilerplate) CheckState(in *states.State) []error {
	// Boilerplate doesn't itself produce anything in the state. It's there
	// only in support of other ConfigObjects.
	return nil
}
