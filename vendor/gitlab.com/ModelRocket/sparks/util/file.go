/*************************************************************************
 * MIT License
 * Copyright (c) 2018 Model Rocket
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package util

import (
	"os"
	"path/filepath"
)

// FindFileReverse will find an absolute path for a file, from cwd
func FindFileReverse(name string) (string, error) {
	var rval string

	cwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return rval, err
	}

	if err := ReverseWalk(cwd, func(path string, info os.FileInfo, err error) error {
		// look for the file here
		path = filepath.Join(path, name)
		if _, err := os.Stat(path); err == nil {
			rval = path
			return os.ErrExist
		}
		return nil
	}); err != os.ErrExist {
		return rval, os.ErrNotExist
	}

	return rval, nil
}

// OpenFileReverse find and open a file in reverse
func OpenFileReverse(name string) (*os.File, error) {
	var rval *os.File

	cwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return rval, err
	}

	if err := ReverseWalk(cwd, func(path string, info os.FileInfo, err error) error {
		// look for the file here
		path = filepath.Join(path, name)
		if fd, err := os.Open(path); err == nil {
			rval = fd
			return os.ErrExist
		}
		return nil
	}); err != os.ErrExist {
		return rval, os.ErrNotExist
	}

	return rval, nil
}

func FileBaseName(name string) string {
	name = filepath.Base(name)
	var extension = filepath.Ext(name)
	return name[0 : len(name)-len(extension)]
}
