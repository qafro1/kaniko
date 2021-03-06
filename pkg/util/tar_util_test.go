/*
Copyright 2018 Google LLC

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

package util

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/GoogleContainerTools/kaniko/testutil"
)

var regularFiles = []string{"file", "file.tar", "file.tar.gz"}
var uncompressedTars = []string{"uncompressed", "uncompressed.tar"}
var compressedTars = []string{"compressed", "compressed.tar.gz"}

func Test_IsLocalTarArchive(t *testing.T) {
	testDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("err setting up temp dir: %v", err)
	}
	defer os.RemoveAll(testDir)
	if err := setUpFilesAndTars(testDir); err != nil {
		t.Fatal(err)
	}
	// Test we get the correct result for regular files
	for _, regularFile := range regularFiles {
		isTarArchive := IsFileLocalTarArchive(filepath.Join(testDir, regularFile))
		testutil.CheckErrorAndDeepEqual(t, false, nil, false, isTarArchive)
	}
	// Test we get the correct result for uncompressed tars
	for _, uncompressedTar := range uncompressedTars {
		isTarArchive := IsFileLocalTarArchive(filepath.Join(testDir, uncompressedTar))
		testutil.CheckErrorAndDeepEqual(t, false, nil, true, isTarArchive)
	}
	// Test we get the correct result for compressed tars
	for _, compressedTar := range compressedTars {
		isTarArchive := IsFileLocalTarArchive(filepath.Join(testDir, compressedTar))
		testutil.CheckErrorAndDeepEqual(t, false, nil, true, isTarArchive)
	}
}

func setUpFilesAndTars(testDir string) error {
	regularFilesAndContents := map[string]string{
		regularFiles[0]: "",
		regularFiles[1]: "something",
		regularFiles[2]: "here",
	}
	if err := testutil.SetupFiles(testDir, regularFilesAndContents); err != nil {
		return err
	}

	for _, uncompressedTar := range uncompressedTars {
		tarFile, err := os.Create(filepath.Join(testDir, uncompressedTar))
		if err != nil {
			return err
		}
		if err := createTar(testDir, tarFile); err != nil {
			return err
		}
	}

	for _, compressedTar := range compressedTars {
		tarFile, err := os.Create(filepath.Join(testDir, compressedTar))
		if err != nil {
			return err
		}
		gzr := gzip.NewWriter(tarFile)
		if err := createTar(testDir, gzr); err != nil {
			return err
		}
	}
	return nil
}

func createTar(testdir string, writer io.Writer) error {

	w := tar.NewWriter(writer)
	defer w.Close()
	for _, regFile := range regularFiles {
		filePath := filepath.Join(testdir, regFile)
		fi, err := os.Stat(filePath)
		if err != nil {
			return err
		}
		if err := AddToTar(filePath, fi, map[uint64]string{}, w); err != nil {
			return err
		}
	}
	return nil
}
