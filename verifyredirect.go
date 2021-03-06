package main

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

func VerifyRedirect(redirect Redirect) error {
	response, err := http.Get(redirect.source)
	if err != nil {
		return err
	}

	defer func() {
		err := response.Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	redirectResponseURL := response.Request.URL.String()
	if redirectResponseURL == redirect.target {
		return nil
	}

	resolveTargetResponse, err := http.Get(redirect.target)
	if err != nil {
		return err
	}

	defer func() {
		err := resolveTargetResponse.Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	resolvedTargetURL := resolveTargetResponse.Request.URL.String()
	if resolvedTargetURL == redirectResponseURL {
		return nil
	}

	return errors.Errorf(
		"source uri %s does resolve to %s,"+
			"not to targetUri %s which resolves to %s",
		redirect.source,
		redirectResponseURL,
		redirect.target,
		resolvedTargetURL)
}
