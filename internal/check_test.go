package internal_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"gitlab.com/Oppodelldog/checkredirects/internal"

	"github.com/stretchr/testify/assert"
)

const testSourceURI = "some-Source-uri"
const testTargetURI = "some-Target-uri"
const redirectsFilename = "redirects"

var originals = struct {
	osArgs             []string
	verifyRedirectFunc internal.VerifyRedirectFuncDef
}{
	osArgs:             os.Args,
	verifyRedirectFunc: internal.VerifyRedirectFunc,
}

func restoreOriginals() {
	os.Args = originals.osArgs
	internal.VerifyRedirectFunc = originals.verifyRedirectFunc
}

func TestMainFunc_ProcessesAllLinesOfFile(t *testing.T) {
	defer restoreOriginals()

	prepareTestTempFolder(t)
	writeTestFile(t, 3)

	numberOfVerifyCalls := 0
	internal.VerifyRedirectFunc = func(ctx context.Context, redirect internal.Redirect) error {
		assert.Exactly(t, testSourceURI, redirect.Source)
		assert.Exactly(t, testTargetURI, redirect.Target)
		numberOfVerifyCalls++

		return nil
	}

	internal.Check(redirectsFilename, 1, "\t")

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

	err := ioutil.WriteFile(redirectsFilename, data, 0600)
	if err != nil {
		t.Fatalf("Did not expect os.Chdir to return an error, but got: %v ", err)
	}
}

func TestEnsureVerifyRedirectImplementationIsUsed(t *testing.T) {
	assert.Exactly(t, reflect.ValueOf(internal.VerifyRedirectFunc).Pointer(), reflect.ValueOf(internal.VerifyRedirectFunc).Pointer())
}
