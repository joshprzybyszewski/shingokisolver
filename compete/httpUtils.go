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
) ([]byte, error) {

	// if url == `https://www.puzzle-shingoki.com/` {
	// 	return ioutil.ReadFile("./compete/samplePuzzleResponse.html")
	// }

	var data io.Reader

	req, err := http.NewRequest(`GET`, url, data)
	if err != nil {
		return nil, err
	}
	if header != nil {
		req.Header = header
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	resBytes, err := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		log.Printf("full response: %+v\n%s\n%s\n",
			response,
			response.Body,
			string(resBytes),
		)

		contentType := response.Header.Get(`Content-Type`)
		if strings.Contains(contentType, `text/plain`) {
			return nil, fmt.Errorf("bad response: \"%s\"", string(resBytes))
		}

		return nil, fmt.Errorf(`bad response from server`)
	} else if err != nil {
		return nil, err
	}

	return resBytes, nil
}
