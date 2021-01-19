package stressgen

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/states"
	"github.com/zclconf/go-cty/cty"
)

// ConfigOutput is an implementation of ConfigObject representing the
// declaration of an output value.
type ConfigOutput struct {
	Addr      addrs.AbsOutputValue
	Value     ConfigExpr
	Sensitive bool
}

var _ ConfigObject = (*ConfigOutput)(nil)

// DisplayName implements ConfigObject.DisplayName.
func (o *ConfigOutput) DisplayName() string {
	return o.Addr.String()
}

// AppendConfig implements ConfigObject.AppendConfig.
func (o *ConfigOutput) AppendConfig(to *hclwrite.Body) {
	block := hclwrite.NewBlock("output", []string{o.Addr.OutputValue.Name})
	body := block.Body()
	body.SetAttributeRaw("value", o.Value.BuildExpr().BuildTokens(nil))
	if o.Sensitive {
		body.SetAttributeValue("sensitive", cty.True)
	}
	to.AppendBlock(block)
}

func (o *ConfigOutput) CheckState(in *states.State) []error {
	// Only root module output values are recorded in a final state.
	// Other output values are only there as an intermediary to pass
	// values between modules.
	// (Technically we _could_ test others if we wanted, by relying on
	// the fact that an in-memory State has a transient cache of them,
	// but the goal of stresstest is to verify externally-visible behavior,
	// not implementation details.)
	if !o.Addr.Module.IsRoot() {
		return nil
	}

	var errs []error
	os := in.OutputValue(o.Addr)
	if os == nil {
		errs = append(errs, fmt.Errorf("root module output value %s is not tracked in the state", o.Addr.OutputValue.Name))
		// This problem prevents us from checking any others.
		return errs
	}
	wantV := o.Value.ExpectedValue()
	gotV := os.Value
	if !wantV.RawEquals(gotV) {
		errs = append(errs, ErrUnexpected{
			Message: fmt.Sprintf("wrong value for root module output value %s", o.Addr.OutputValue.Name),
			Got:     gotV,
			Want:    wantV,
		})
	}
	if got, want := os.Sensitive, o.Sensitive; got != want {
		errs = append(errs, ErrUnexpected{
			Message: fmt.Sprintf("wrong sensitive flag for root module output value %s", o.Addr.OutputValue.Name),
			Got:     got,
			Want:    want,
		})
	}
	return errs
}
