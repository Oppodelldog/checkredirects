package internal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

var ErrInvalidSourceRequest = errors.New("Invalid source request")
var ErrInvalidTargetRequest = errors.New("Invalid target request")
var ErrSourceRequest = errors.New("error in requesting source")
var ErrTargetRequest = errors.New("error in requesting target")
var ErrRedirectCheckFailed = errors.New("error redirect check failed")

func VerifyRedirect(ctx context.Context, redirect Redirect) error {
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, redirect.Source, nil)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidSourceRequest, err)
	}

	response, err := http.DefaultClient.Do(r)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSourceRequest, err)
	}

	defer func() {
		err := response.Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	redirectResponseURL := response.Request.URL.String()
	if redirectResponseURL == redirect.Target {
		return nil
	}

	r, err = http.NewRequestWithContext(ctx, http.MethodGet, redirect.Target, nil)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidTargetRequest, err)
	}

	resolveTargetResponse, err := http.DefaultClient.Do(r)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrTargetRequest, err)
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

	return fmt.Errorf(
		"%w: Source uri %s does resolve to %s,"+
			"not to targetUri %s which resolves to %s",
		ErrRedirectCheckFailed,
		redirect.Source,
		redirectResponseURL,
		redirect.Target,
		resolvedTargetURL)
}
