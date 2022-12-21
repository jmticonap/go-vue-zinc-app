package parsers

import (
	"bytes"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ListFiles(path string, limit int) []string {
	var result []string
	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			result = append(result, path)
			limit--
		}
		if limit == 0 {
			return io.EOF
		}

		return nil
	})

	return result
}

func TrimValue(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Trim(value, "\r\n")
	value = strings.Trim(value, "\n")
	value = strings.Trim(value, "\t")

	return value
}

func Check_error(err error) bool {
	return err != nil
}

func Check_fatal(err error) bool {
	return err != nil
}

func GetListFiles(dir_path string) []fs.DirEntry {
	files, err := os.ReadDir(dir_path)
	Check_error(err)
	return files
}

func GetEmailFileString(path string) string {
	content, err := os.ReadFile(path)
	Check_fatal(err)
	return string(content)
}

func Send(method, url, body, auth string) (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(method, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("user-agent", "golang application")
	req.Header.Add("Authorization", "Basic "+b64.StdEncoding.EncodeToString([]byte(auth)))

	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("got error %s", err.Error())
	}
	//defer response.Body.Close()
	return response, nil
}

type ParserFunc func(base_path string, limit int) string

func LoadNewData(path, auth string, limit int, parser ParserFunc) {
	if path == "" {
		log.Fatal(errors.New("must pass path argument"))
	} else {
		str := parser(path, limit)

		resp, e := Send(http.MethodPost, os.Getenv("ZINC_SEARCH_HOST")+"/api/_bulk", str, auth)

		if e != nil {
			log.Print(e)
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			bodyString := string(bodyBytes)
			log.Println(bodyString)
		}
	}
}
