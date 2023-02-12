package parsers

import (
	"bytes"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"jmtp/indexer/commons"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func ListFiles(path string, limit int) []string {
	var result []string
	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}
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

func CheckFatal(err error) bool {
	return err != nil
}

func GetListFiles(dir_path string) []fs.DirEntry {
	files, err := os.ReadDir(dir_path)
	Check_error(err)
	return files
}

func GetEmailFileString(path string) string {
	content, err := os.ReadFile(path)
	CheckFatal(err)
	return string(content)
}

func Send(method, url, body, auth string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBufferString(body))
	if err != nil {
		return nil, fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Set("user-agent", "golang application")
	req.Header.Add("Authorization", "Basic "+b64.StdEncoding.EncodeToString([]byte(auth)))

	response, err := commons.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("got error %s", err.Error())
	}

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

func LoadBulk(body string) {
	resp, err := Send(http.MethodPost, os.Getenv("ZINC_SEARCH_HOST")+"/api/_bulk", body, commons.Auth)

	if err != nil {
		redF := color.New(color.FgRed).SprintFunc()
		redB := color.New(color.Bold, color.BgRed).SprintFunc()
		log.Fatalf("%s %s\n", redB(" ERROR "), redF(err.Error()))
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

func LoadBulkV2(body *strings.Builder) {
	resp, err := Send(http.MethodPost, os.Getenv("ZINC_SEARCH_HOST")+"/api/_bulk", body.String(), commons.Auth)

	if err != nil {
		redF := color.New(color.FgRed).SprintFunc()
		redB := color.New(color.Bold, color.BgRed).SprintFunc()
		log.Fatalf("%s %s\n", redB(" ERROR "), redF(err.Error()))
	}

	defer resp.Body.Close()
	defer body.Reset()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		log.Println(bodyString)
	}
}
