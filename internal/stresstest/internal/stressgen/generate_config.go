package stressgen

import "github.com/hashicorp/hcl/v2/hclwrite"

// GenerateConfigFile generates the potential content of a single configuration
// (.tf) file which declares all of the given configuration objects.
//
// It's the caller's responsibility to make sure that the given objects all
// make sense to be together in a single module, including making sure they all
// together meet any uniqueness constraints and that any objects that refer
// to other objects are given along with the objects they refer to.
func GenerateConfigFile(objs []ConfigObject) []byte {
	f := hclwrite.NewEmptyFile()
	body := f.Body()
	for _, obj := range objs {
		obj.AppendConfig(body)
	}
	return f.Bytes()
}
