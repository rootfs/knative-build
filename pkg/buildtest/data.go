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

package buildtest

import (
	"encoding/json"
	"io/ioutil"
)

func Data(relpath string) ([]byte, error) {
	b, err := ioutil.ReadFile(relpath)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func DataAs(relpath string, obj interface{}) error {
	b, err := Data(relpath)
	if err != nil {
		return err
	}
	jb, err := YAMLToJSON(b)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(jb, &obj); err != nil {
		return err
	}
	return nil
}
