/*
 * Copyright 2020 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package web provides convenience routines to access files in the data
// directory.
package web

import (
	"path/filepath"
	"runtime"
)

// basePath is the root directory of this package.
var basePath string

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	basePath = filepath.Dir(currentFile)
}

// Path returns the absolute path the given relative file or directory path,
// relative to the google.golang.org/grpc/examples/data directory in the
// user's GOPATH.  If rel is already absolute, it is returned unmodified.
func Path(rel string) string {
	if filepath.IsAbs(rel) {
		return rel
	}

	return filepath.Join(basePath, rel)
}
