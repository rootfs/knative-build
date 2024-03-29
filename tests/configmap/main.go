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

package main

import (
	"io/ioutil"
	"os"

	"github.com/golang/glog"
)

func main() {
	configData := os.Getenv("TEST_DATA")
	b, err := ioutil.ReadFile("/config/test.data")
	if err != nil {
		glog.Fatalf("Unexpected error reading /config/test.data: %v", err)
	}
	otherConfigData := string(b)
	if configData != otherConfigData {
		glog.Fatalf("Unexpected mismatch in config data, want %q, but got %q", configData, otherConfigData)
	}
}
