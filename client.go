// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/31

package gomodule

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	modzip "golang.org/x/mod/zip"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func client(c HTTPClient) HTTPClient {
	if c != nil {
		return c
	}
	return http.DefaultClient
}

func sentRequest(ctx context.Context, c HTTPClient, method string, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("UserAgent", defaultUA)
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = checkResponse(resp); err != nil {
		return nil, err
	}
	return io.ReadAll(resp.Body)
}

func checkResponse(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		bf, err := io.ReadAll(io.LimitReader(resp.Body, 1024))
		if err != nil {
			return err
		}
		if len(bf) > 0 {
			return errors.New(string(bytes.TrimSpace(bf)))
		}
		return fmt.Errorf("invalid resp status %q", resp.Status)
	}

	if resp.ContentLength > modzip.MaxZipFile {
		return fmt.Errorf("response body is too large (%d bytes; limit is %d bytes)", resp.ContentLength, modzip.MaxZipFile)
	}

	if resp.ContentLength == -1 {
		resp.Body = io.NopCloser(io.LimitReader(resp.Body, modzip.MaxZipFile))
	}
	return nil
}
