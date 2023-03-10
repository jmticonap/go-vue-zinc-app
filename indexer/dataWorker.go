package main

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func check_error(err error) bool {
	return err != nil
}

func check_fatal(err error) bool {
	return err != nil
}

func getListFiles(dir_path string) []fs.DirEntry {
	files, err := os.ReadDir(dir_path)
	check_error(err)
	return files
}

func getEmailFileString(path string) string {
	content, err := os.ReadFile(path)
	check_fatal(err)
	return string(content)
}

func trimValue(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Trim(value, "\r\n")
	value = strings.Trim(value, "\n")
	value = strings.Trim(value, "\t")

	return value
}

func DataParser(base_path string, limit int) string {
	var result string
	files, err := os.ReadDir(base_path) //dirs of all users
	check_error(err)

	//========================================
	var i int
	for _, file := range files {
		if file.IsDir() {
			path := base_path + "/" + file.Name()

			email_dirs := getListFiles(path)

			for _, email_dir := range email_dirs {
				emails_files_path := path + "/" + email_dir.Name()
				emails_files := getListFiles(emails_files_path)

				for _, email_file := range emails_files {
					//===============================================================================
					email_content := getEmailFileString(emails_files_path + "/" + email_file.Name())

					keys := []string{
						"Message-ID:", "Date:", "From:", "To:",
						"Subject:", "Mime-Version:", "Content-Type:",
						"Content-Transfer-Encoding:", "X-From:",
						"X-To:", "X-cc:", "X-bcc:", "X-Folder:",
						"X-Origin:", "X-FileName:",
					}

					mail := make(map[string]string)
					mail["user"] = file.Name()
					mail["category"] = email_dir.Name()
					for i := 0; i < len(keys); i++ {
						base_index := strings.Index(email_content, keys[i])
						if i == len(keys)-1 {
							last_values := strings.Split(email_content[base_index+len(keys[i]):], "\n")

							mail[strings.Trim(keys[i], ":")] = trimValue(last_values[0])
							content := last_values[1:]

							mail["Content"] = strings.TrimLeft(trimValue(strings.Join(content, "\n")), "\n\t")
						} else {
							from := base_index + len(keys[i])
							next_base_index := strings.Index(email_content, keys[i+1])
							if from < next_base_index {
								mail[strings.Trim(keys[i], ":")] = trimValue(email_content[from:next_base_index])
							} else {
								mail[strings.Trim(keys[i], ":")] = "none"
							}
						}
					}
					jsonStr, err := json.Marshal(mail)
					if err != nil {
						log.Fatal(err)
					}
					result += `{ "index" : { "_index" : "mails" } }` + "\n" + string(jsonStr) + "\n"
					if i == limit-1 {
						return result
					}
					i++
				}
			}
		}
	}
	return result
}

func DataParserV2(base_path string, limit int) string {
	var result string
	files := ListFiles(base_path, limit)

	for _, file := range files {
		email_content := getEmailFileString(file)
		mail, _ := DataRegexParser(email_content)
		jsonStr, err := json.Marshal(mail)

		if err != nil {
			log.Fatal(err)
		}
		result += `{ "index" : { "_index" : "mails" } }` + "\n" + string(jsonStr) + "\n"
	}

	return result
}

func DataRegexParser(str string) (map[string]string, error) {
	str = strings.Replace(str, "\r\n\r\n", "\n\n", 1)
	result := make(map[string]string)
	mailParts := strings.Split(str, "\n\n")
	header := mailParts[0]
	content := strings.Join(mailParts[1:], "\n\n")

	var re_key = regexp.MustCompile(`(?m)^(Message-ID:|Date:|From:|To:|Subject:|Cc:|Mime-Version:|Content-Type:|Content-Transfer-Encoding:|Bcc:|X-From:|X-To:|X-cc:|X-bcc:|X-Folder:|X-Origin:|X-FileName:)`)
	var keys = re_key.FindAllString(header, -1)

	var re_value = regexp.MustCompile(`(?m): .*$`)
	var values = re_value.FindAllString(header, -1)

	if len(keys) != len(values) {
		return nil, errors.New("problems with keys and values")
	}

	for i, match := range keys {
		var key = strings.ToLower(match[:len(match)-1])
		if len(values[i]) > 2 {
			result[key] = values[i][1:]
		} else {
			result[key] = ""
		}
	}
	result["content"] = content

	return result, nil
}

func send(method, url, body, auth string) (*http.Response, error) {
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

// ServeHTTP calls f(w, r).
func (f ParserFunc) ParserMail(base_path string, limit int) string {
	return f(base_path, limit)
}

func loadNewData(path, auth string, limit int, parser ParserFunc) {
	if path == "" {
		log.Fatal(errors.New("must pass path argument"))
	} else {
		//base_path := "/home/juancho/Programming/GoProjects/enron_mail_20110402/maildir"
		str := parser(path, limit)
		//err := os.WriteFile("mails.ndjson", []byte(str), 0644)
		//check_error(err)

		resp, e := send(http.MethodPost, "http://localhost:4080/api/_bulk", str, auth)

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
