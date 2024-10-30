// Copyright 2024 k8s-crafts Authors

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testutils

import (
	"math/rand"
	"path"
	"strings"
)

// Generate a random string with a specific length
// NOTE: k8s requires object names to be lowercase
func generateRandom(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"

	b := make([]byte, length)

	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}

	return string(b)
}

// Working directory: <project-root>/e2e
func getTestdataDir() string {
	return path.Join("testdata")
}

func IsErrorNotFound(err error) bool {
	return err != nil && strings.Contains(err.Error(), "not found")
}
