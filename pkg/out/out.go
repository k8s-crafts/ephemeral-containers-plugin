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

package out

import (
	"fmt"
	"io"

	klog "k8s.io/klog/v2"
)

var (
	// destination where output is sent
	outFile io.Writer

	// destination where error output is sent
	errFile io.Writer
)

// Set outFile to an io.Writer
func SetOutFile(file io.Writer) {
	outFile = file
}

// Get outFile
func GetOutFile() io.Writer {
	return outFile
}

// Set errFile to an io.Writer
func SetErrFile(file io.Writer) {
	errFile = file
}

// Get errFile
func GetErrFile() io.Writer {
	return errFile
}

// Write a formatted string with a newline to stdout
func Stringf(format string, a ...interface{}) {
	// Flush log to ensure correct log order
	klog.Flush()
	defer klog.Flush()

	if outFile == nil {
		klog.Errorf("[unset outFile]: %s", fmt.Sprintf(format, a...))
		return
	}

	klog.Infof(format, a...)

	if _, err := fmt.Fprintf(outFile, format, a...); err != nil {
		klog.Errorf("FPrint failed: %v", err)
	}
}

// Write a formatted string with a newline to stderr
func Errf(format string, a ...interface{}) {
	// Flush log to ensure correct log order
	klog.Flush()
	defer klog.Flush()

	if errFile == nil {
		klog.Errorf("[unset errFile]: %s", fmt.Sprintf(format, a...))
		return
	}

	klog.Warningf(format, a...)

	if _, err := fmt.Fprintf(errFile, format, a...); err != nil {
		klog.Errorf("FPrint failed: %v", err)
	}
}

// Write a formatted string with a newline to stdout
func Ln(format string, a ...interface{}) {
	Stringf(format+"\n", a...)
}

// Write a formatted string with a newline to stderr
func ErrLn(format string, a ...interface{}) {
	Errf(format+"\n", a...)
}
