package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testSourceURI = "some-source-uri"
const testTargetURI = "some-target-uri"

var originals = struct {
	osArgs             []string
	verifyRedirectFunc verifyRedirectFuncDef
}{
	osArgs:             os.Args,
	verifyRedirectFunc: verifyRedirectFunc,
}

func restoreOriginals() {
	os.Args = originals.osArgs
	verifyRedirectFunc = originals.verifyRedirectFunc
}

func TestMainFunc_ProcessesAllLinesOfFile(t *testing.T) {
	defer restoreOriginals()

	prepareTestTempFolder(t)
	writeTestFile(t, 3)

	numberOfVerifyCalls := 0
	verifyRedirectFunc = func(redirect Redirect) error {
		assert.Exactly(t, testSourceURI, redirect.source)
		assert.Exactly(t, testTargetURI, redirect.target)
		numberOfVerifyCalls++

		return nil
	}

	os.Args = []string{"", "-c=1"}

	main()

	assert.Exactly(t, 3, numberOfVerifyCalls)
}

func prepareTestTempFolder(t *testing.T) {
	const testTempFolder = "/tmp/checkredirects"

	err := os.RemoveAll(testTempFolder)
	if err != nil {
		t.Fatalf("Did not expect os.Remove to return an error, but got: %v ", err)
	}

	err = os.MkdirAll(testTempFolder, 0777)
	if err != nil {
		t.Fatalf("Did not expect os.MkdirAll to return an error, but got: %v ", err)
	}

	err = os.Chdir(testTempFolder)
	if err != nil {
		t.Fatalf("Did not expect os.Chdir to return an error, but got: %v ", err)
	}
}

func writeTestFile(t *testing.T, linesToCheck int) {
	var data []byte

	for i := 0; i < linesToCheck; i++ {
		data = append(data, []byte(fmt.Sprintf("%s\t%s\n", testSourceURI, testTargetURI))...)
	}

	err := ioutil.WriteFile(redirectsFileName, data, 0600)
	if err != nil {
		t.Fatalf("Did not expect os.Chdir to return an error, but got: %v ", err)
	}
}

func TestEnsureVerifyRedirectImplementationIsUsed(t *testing.T) {
	assert.Exactly(t, reflect.ValueOf(verifyRedirectFunc).Pointer(), reflect.ValueOf(VerifyRedirect).Pointer())
}
