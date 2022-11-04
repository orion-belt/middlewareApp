package common
import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func SendHttpRequest(
	method, url,
	payload string,
	client *http.Client,
	// optional headers, each header is a slice of strings in the form:
	//   [<header name>, value1, value2...]
	headers ...[]string,
) (int, string, error) {

	var body io.Reader = nil
	if len(payload) > 0 {
		body = strings.NewReader(payload)
	}
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return 0, "", err
	}
	for _, h := range headers {
		if l := len(h); l > 0 {
			k := h[0]
			if l > 1 {
				request.Header.Set(k, h[1])
				if l > 2 {
					for _, v := range h[2:] {
						request.Header.Add(k, v)
					}
				}
			} else {
				request.Header.Set(k, "")
			}
		}
		request.Header.Set(h[0], h[1])
	}
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
		// request.Header.Set("radius-packet-encoding", "base64/binary")
	}
	// var client = GetHttpClient()

	response, error := client.Do(request)
	if error != nil {
		return response.StatusCode, "", err
	}

	defer response.Body.Close()
	contents, _ := ioutil.ReadAll(response.Body)

	//return response.StatusCode, string(contents), nil
	return response.StatusCode, string(contents), nil
}
