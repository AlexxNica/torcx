// Copyright 2017 CoreOS Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ftests

import (
	"os"
	"os/exec"
	"testing"
)

func TestGeneratorEmpty(t *testing.T) {
	if !IsInContainer() {
		cfg := RktConfig{
			imageName: EmptyImage,
		}
		RunTestInContainer(t, cfg)
		return
	}

	// Prepare an empty vendor profile
	_ = os.MkdirAll("/usr/share/torcx/profiles", 0755)
	fp, err := os.Create("/usr/share/torcx/profiles/vendor.json")
	if err != nil && !os.IsExist(err) {
		t.Fatal(err)
	}
	fp.Close()

	cmd := exec.Command("torcx-generator")
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		t.Error(string(bytes))
	}
}