package stressgen

import (
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// ConfigExpr is an interface implemented by types that represent various
// kinds of expression that are relevant to our testing.
//
// Since stresstest is focused mainly on testing graph building and graph
// traversal behaviors, and not on expression evaluation details, we don't
// aim to cover every possible kind of expression here but should aim to model
// all kinds of expression that can contribute in some way to the graph shape.
type ConfigExpr interface {
	// BuildExpr builds the hclwrite representation of the recieving expression,
	// for inclusion in the generated configuration files.
	BuildExpr() *hclwrite.Expression

	// ExpectedValue returns the value this expression ought to return if
	// Terraform behaves correctly. This must be the specific, fully-known
	// value we expect to find in the final state, not any placeholder value
	// that might show up during planning if we were faking a computed resource
	// argument.
	ExpectedValue() cty.Value
}

// ConfigExprConst is an implementation of ConfigExpr representing static,
// constant values.
type ConfigExprConst struct {
	Value cty.Value
}

var _ ConfigExpr = (*ConfigExprConst)(nil)

// BuildExpr implements ConfigExpr.BuildExpr
func (e *ConfigExprConst) BuildExpr() *hclwrite.Expression {
	return hclwrite.NewExpressionLiteral(e.Value)
}

// ExpectedValue implements ConfigExpr.ExpectedValue
func (e *ConfigExprConst) ExpectedValue() cty.Value {
	return e.Value
}
