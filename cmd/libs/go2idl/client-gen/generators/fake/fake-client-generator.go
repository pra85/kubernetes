/*
Copyright 2015 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fake

import (
	"path/filepath"
	"strings"

	"k8s.io/kubernetes/cmd/libs/go2idl/generator"
	"k8s.io/kubernetes/cmd/libs/go2idl/types"
)

func PackageForGroup(group string, version string, typeList []*types.Type, packageBasePath string, srcTreePath string, boilerplate []byte) generator.Package {
	outputPackagePath := filepath.Join(packageBasePath, group, version, "fake")
	// TODO: should make this a function, called by here and in client-generator.go
	realClientPath := filepath.Join(packageBasePath, group, version)
	return &generator.DefaultPackage{
		PackageName: "fake",
		PackagePath: outputPackagePath,
		HeaderText:  boilerplate,
		PackageDocumentation: []byte(
			`// Package fake has the automatically generated clients.
`),
		// GeneratorFunc returns a list of generators. Each generator makes a
		// single file.
		GeneratorFunc: func(c *generator.Context) (generators []generator.Generator) {
			generators = []generator.Generator{
				// Always generate a "doc.go" file.
				generator.DefaultGen{OptionalName: "doc"},
			}
			// Since we want a file per type that we generate a client for, we
			// have to provide a function for this.
			for _, t := range typeList {
				generators = append(generators, &genFakeForType{
					DefaultGen: generator.DefaultGen{
						OptionalName: "fake_" + strings.ToLower(c.Namers["private"].Name(t)),
					},
					outputPackage: outputPackagePath,
					group:         group,
					typeToMatch:   t,
					imports:       generator.NewImportTracker(),
				})
			}

			generators = append(generators, &genFakeForGroup{
				DefaultGen: generator.DefaultGen{
					OptionalName: "fake_" + group + "_client",
				},
				outputPackage:  outputPackagePath,
				realClientPath: realClientPath,
				group:          group,
				types:          typeList,
				imports:        generator.NewImportTracker(),
			})
			return generators
		},
		FilterFunc: func(c *generator.Context, t *types.Type) bool {
			return types.ExtractCommentTags("+", t.SecondClosestCommentLines)["genclient"] == "true"
		},
	}
}
