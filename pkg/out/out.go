// MIT License

// Copyright (c) 2024 k8s-crafts Authors

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
