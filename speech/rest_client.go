package speech

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/vito-ai/auth"
	"github.com/vito-ai/auth/option"
)

var ErrNotFinish = errors.New("result is not complete yet")
var ErrFailed = errors.New("result failed")

type restClient struct {
	// endpoint to rtzr api server host
	endpoint string

	//httpClient
	httpClient *http.Client
}

// Make New Client for RESTful STT API
func NewRestClient(cliopts *option.ClientOption) (*restClient, error) {
	if cliopts == nil {
		cliopts = option.DefaultClientOption()
	}
	httpClient, err := auth.NewAuthClient(cliopts)
	if err != nil {
		return nil, err
	}

	c := &restClient{
		endpoint:   cliopts.GetRestEndpoint(),
		httpClient: httpClient,
	}

	return c, nil
}

func (c *restClient) Close() error {
	c.httpClient = nil
	return nil
}

func (c *restClient) Recognize(ctx context.Context, param *RecognizeRequest) (*RecognizeResponse, error) {
	resId, err := c.RecognizeAsync(ctx, param)
	if err != nil {
		return nil, err
	}

	resp, err := c.receiveResultWithPolling(ctx, resId, 4*time.Second)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *restClient) RecognizeAsync(ctx context.Context, param *RecognizeRequest) (ResultId, error) {
	isPipeClose := false

	r, w := io.Pipe()
	writer := multipart.NewWriter(w)
	defer func() {
		if !isPipeClose {
			r.Close()
			w.Close()
			writer.Close()
		}
	}()

	err := param.AudioSource.validate()
	if err != nil {
		return "", err
	}

	errCh := make(chan error, 1)
	defer close(errCh)

	go func() {
		defer w.Close()
		if err := createConfigField(writer, param.Config); err != nil {
			errCh <- err
			return
		}
		if param.AudioSource.FilePath != "" {
			if err := createFileFieldWithLocal(writer, param.AudioSource.FilePath); err != nil {
				errCh <- err
				return
			}
		} else {
			if err := createFileFieldWithData(writer, param.AudioSource.Content); err != nil {
				errCh <- err
				return
			}
		}
		if err := writer.Close(); err != nil {
			errCh <- err
			return
		}

		errCh <- nil
	}()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, r)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	response, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	r.Close()
	isPipeClose = true
	defer response.Body.Close()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case err := <-errCh:
		if err != nil {
			return "", err
		}
	}
	resByte, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	if response.StatusCode != 200 {
		return "", fmt.Errorf("server error : %d\n%s", response.StatusCode, string(resByte))
	}
	result := &RecognizeResponse{}
	if err = json.Unmarshal(resByte, &result); err != nil {
		return "", err
	}

	return result.Id, nil
}

func (c *restClient) ReceiveResult(ctx context.Context, resultId ResultId) (*RecognizeResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpoint+"/"+string(resultId), nil)
	if err != nil {
		return nil, err
	}

	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("server request error")
	}
	defer response.Body.Close()

	result := &RecognizeResponse{}
	resByte, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(resByte, &result); err != nil {
		return nil, err
	}
	switch result.Status {
	case "completed":
		return result, nil
	case "transcribing":
		return nil, ErrNotFinish
	case "failed":
		return nil, ErrFailed
	default:
		return nil, fmt.Errorf("server response error : %s", string(resByte))
	}
}

func (c *restClient) receiveResultWithPolling(ctx context.Context, resultId ResultId, delay time.Duration) (*RecognizeResponse, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(delay):
			res, err := c.ReceiveResult(ctx, resultId)
			if err != nil {
				if errors.Is(err, ErrNotFinish) {
					continue
				}
				return nil, err
			}

			if res != nil {
				return res, nil
			}
			return nil, errors.New("nil response return")
		}
	}
}

func createFileFieldWithLocal(writer *multipart.Writer, filePath string) error {
	audiofile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer audiofile.Close()

	fw, err := writer.CreateFormFile("file", audiofile.Name())
	if err != nil {
		return err
	}

	if _, err = io.Copy(fw, audiofile); err != nil {
		return err
	}
	return nil
}

func createFileFieldWithData(writer *multipart.Writer, contents []byte) error {
	fw, err := writer.CreateFormFile("file", "rtzr-default-audiofile")
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(contents)
	if _, err = io.Copy(fw, buffer); err != nil {
		return err
	}

	return nil
}

func createConfigField(writer *multipart.Writer, config RecognitionConfig) error {
	fw, err := writer.CreateFormField("config")
	if err != nil {
		return err
	}

	j, err := json.Marshal(config)
	if err != nil {
		return err
	}
	if _, err := fw.Write(j); err != nil {
		return err
	}

	return nil
}
