/*
do request for  formData with file
*/
package form_data_request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/artisanhe/tools/courier/status_error"
)

func NewFileRequest(uri string, fileHeader *multipart.FileHeader) (*http.Request, error) {
	fileContents, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer fileContents.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreatePart(fileHeader.Header)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, fileContents)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, nil
}

func DoClient(req *http.Request, timeOut time.Duration, respResult interface{}) error {
	client := &http.Client{
		Timeout: timeOut,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	if resp.StatusCode/200 == 1 {
		err := json.Unmarshal(data, respResult)
		if err != nil {
			msg := fmt.Sprintf("[%d] %s %s", resp.StatusCode, req.Method, req.URL)
			return status_error.HttpRequestFailed.StatusError().WithDesc(msg)
		}
		return nil
	} else {
		statusError := &status_error.StatusError{}
		err := json.Unmarshal(data, statusError)
		if err != nil {
			msg := fmt.Sprintf("[%d] %s %s", resp.StatusCode, req.Method, req.URL)
			return status_error.HttpRequestFailed.StatusError().WithDesc(msg)
		}
		return nil
	}
}
