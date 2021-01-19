package stressgen

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"
)

// Registry is a container for data to help us to randomly generate valid
// references between objects, and to ensure those references will
// remain valid under randomly-generated modifications.
//
// While we try to keep the ConfigObjects in a randomly-generated configuration
// as self-contained as possible, a lot of Terraform behaviors only emerge
// as a result of references between objects and so we need to be able to
// randomly generate those too. This type is here to coordinate that.
//
// Each randomly-generated module has its own Registry, because each Terraform
// module has its own separate namespace. Mirroring the usual configuration
// structure, each randomly-generated configuration has one root Registry and
// then an additional Registry for each of the child modules it calls.
type Registry struct {
	// Parent and Children together represent the tree of registries, which
	// mirrors the tree of modules in the generated configuration.
	Parent   *Registry
	Children map[string]*Registry

	// ExprRefTargets is a set (represented as a slice whose order is irrelevant)
	// of candidates for references in expressions generated in later
	// configuration objects.
	//
	// DefunctExprRefTargets tracks references that were available in the
	// previous incarnation of a module that we're currently "editing", but
	// that have now been removed. The value in this map is a suitable
	// expression to use instead of the original reference in any expression
	// that previously referred to it and will continue existing in the new
	// version.
	//
	// Typically an ExprRefTarget is a ConfigExprRef and a DefunctExprRefTarget
	// is a ConfigExprConst with the same value as the original target, but
	// generators consuming these should not typically worry about that and
	// should instead just treat the expression generically, calling
	// ConfigExpr.ExpectedValue to determine what value the expression ought
	// to have after apply.
	ExprRefTargets        []ConfigExpr
	DefunctExprRefTargets map[string]ConfigExpr

	// InputVariableValues tracks values thet the caller must set for some or
	// all of the input variables defined inside the corresponding generated
	// module.
	//
	// Part of generating an ConfigVariable is deciding (randomly) whether to
	// explicitly set it or to let its default (if any) take effect, and this
	// collection only tracks the case where the variable should be explicitly
	// set by the caller. The caller here might either be the test harness
	// itself, for a root module, or in the calling module block for a child
	// module. The config generator is responsible for ensuring that there's
	// always an entry in here for any variable that doesn't define a default,
	// and is therefore required.
	InputVariables map[string]cty.Value
}

// NewRootRegistry creates and returns an empty registry that has no parent
// registry.
func NewRootRegistry() *Registry {
	return &Registry{
		Children: make(map[string]*Registry),
	}
}

// NewChild creates and returns an empty registry that is registered as a child
// of the reciever.
//
// The given name must be unique within the space of child registry names in
// the reciever, or this function will panic. The name of a child registry
// should match the name of the module call that implied its existence.
func (r *Registry) NewChild(name string) *Registry {
	if _, exists := r.Children[name]; exists {
		panic(fmt.Sprintf("registry already has a child module %q", name))
	}
	ret := NewRootRegistry()
	ret.Parent = r
	return r
}
