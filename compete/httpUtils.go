package compete

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func get(
	url string,
	header http.Header,
) ([]byte, http.Header, error) {

	// if url == `https://www.puzzle-shingoki.com/` {
	// 	return ioutil.ReadFile("./compete/exampleRequests/samplePuzzleResponse.html")
	// }

	return doRequest(`GET`, url, header, nil)
}

func post(
	url string,
	header http.Header,
	data io.Reader,
) ([]byte, http.Header, error) {

	// if url == `https://www.puzzle-shingoki.com/` {
	// 	return ioutil.ReadFile("./compete/exampleRequests/samplePostCompletePuzzleResponse.html")
	// }

	return doRequest(`POST`, url, header, nil)
}

func doRequest(
	method string,
	url string,
	header http.Header,
	data io.Reader,
) ([]byte, http.Header, error) {

	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, nil, err
	}
	if header != nil {
		req.Header = header
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer response.Body.Close()

	writeToFile(`./temp/lastRequestHeaders.txt`, []byte(fmt.Sprintf("%+v\n", response.Header)))

	if response.StatusCode != http.StatusOK {
		resBytes, err := ioutil.ReadAll(response.Body)
		log.Printf("full response: %+v\n%s\n%s\n",
			response,
			response.Body,
			string(resBytes),
		)
		if err != nil {
			return nil, nil, err
		}

		contentType := response.Header.Get(`Content-Type`)
		if strings.Contains(contentType, `text/plain`) {
			return nil, nil, fmt.Errorf("bad response: \"%s\"", string(resBytes))
		}

		return nil, nil, fmt.Errorf(`bad response from server`)
	}

	resBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, nil, err
	}
	return resBytes, response.Header, nil
}
