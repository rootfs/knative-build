/*
Copyright 2018 Google, Inc. All rights reserved.

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
package builder

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"

	"github.com/knative/build/pkg/apis/build/v1alpha1"
)

// ApplyTemplate applies the values in the template to the build, and replaces
// placeholders for declared parameters with the build's matching arguments.
func ApplyTemplate(u *v1alpha1.Build, tmpl *v1alpha1.BuildTemplate) (*v1alpha1.Build, error) {
	build := u.DeepCopy()
	if tmpl == nil {
		return build, nil
	}
	tmpl = tmpl.DeepCopy()
	build.Spec.Steps = tmpl.Spec.Steps
	build.Spec.Volumes = append(build.Spec.Volumes, tmpl.Spec.Volumes...)

	// Apply template arguments or parameter defaults.
	replacements := map[string]string{}
	if tmpl != nil {
		for _, p := range tmpl.Spec.Parameters {
			if p.Default != nil {
				replacements[p.Name] = *p.Default
			}
		}
	}
	if build.Spec.Template != nil {
		for _, a := range build.Spec.Template.Arguments {
			replacements[a.Name] = a.Value
		}
	}

	applyReplacements := func(in string) string {
		for k, v := range replacements {
			in = strings.Replace(in, fmt.Sprintf("${%s}", k), v, -1)
		}
		return in
	}

	// Apply variable expansion to steps fields.
	steps := build.Spec.Steps
	for i := range steps {
		steps[i].Name = applyReplacements(steps[i].Name)
		for ia, a := range steps[i].Args {
			steps[i].Args[ia] = applyReplacements(a)
		}
		for ie, e := range steps[i].Env {
			steps[i].Env[ie].Value = applyReplacements(e.Value)
		}
		steps[i].WorkingDir = applyReplacements(steps[i].WorkingDir)
		for ic, c := range steps[i].Command {
			steps[i].Command[ic] = applyReplacements(c)
		}
		for iv, v := range steps[i].VolumeMounts {
			steps[i].VolumeMounts[iv].Name = applyReplacements(v.Name)
			steps[i].VolumeMounts[iv].MountPath = applyReplacements(v.MountPath)
			steps[i].VolumeMounts[iv].SubPath = applyReplacements(v.SubPath)
		}
	}

	if buildTmpl := build.Spec.Template; buildTmpl != nil && len(buildTmpl.Env) > 0 {
		// Apply variable expansion to the build's overriden
		// environment variables
		for i, e := range buildTmpl.Env {
			buildTmpl.Env[i].Value = applyReplacements(e.Value)
		}

		for i := range steps {
			steps[i].Env = applyEnvOverride(steps[i].Env, buildTmpl.Env)
		}
	}

	return build, nil
}

func applyEnvOverride(src, override []corev1.EnvVar) []corev1.EnvVar {
	result := make([]corev1.EnvVar, 0, len(src)+len(override))
	overrides := make(map[string]bool)

	for _, env := range override {
		overrides[env.Name] = true
	}

	for _, env := range src {
		if _, present := overrides[env.Name]; !present {
			result = append(result, env)
		}
	}

	return append(result, override...)
}
