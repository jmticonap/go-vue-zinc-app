package parsers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"jmtp/indexer/commons"

	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/fatih/color"
)

var yellow = color.New(color.Bold, color.BgHiWhite).SprintFunc()
var (
	chnMailContent chan string            = make(chan string, 10_000)
	chnMailStruct  chan map[string]string = make(chan map[string]string, 10_000)
	// readingMails     bool                   = false
	// parsingMails     bool                   = false
	noReadFilesCount int             = 0
	countMails       int             = 0
	strBulk          strings.Builder = strings.Builder{}
)

func ReadMails(wg *sync.WaitGroup) {
	err := filepath.WalkDir(commons.DbPath, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			content, err := os.ReadFile(path)
			if err == nil { // if is error pass the next file
				chnMailContent <- string(content)
				countMails += 1
			}
		}
		return nil
	})
	wg.Done()
	if err != nil {
		red := color.New(color.Bold, color.BgRed).SprintFunc()
		fmt.Printf("%s\n", red(" Mails reads: ERROR "))
	} else {
		green := color.New(color.Bold, color.BgGreen).SprintFunc()
		fmt.Printf("%s\n", green(" Mails reads: COMPLETE "))
	}
}

func ParseMails(wg *sync.WaitGroup, w http.ResponseWriter) {
	parsed := 1
	for parsed <= commons.FilesQuantity-noReadFilesCount {
		if len(chnMailContent) > 0 {
			mailContent, ok := <-chnMailContent
			if ok {
				mail, _ := DataRegexParser(mailContent)
				chnMailStruct <- mail
				parsed += 1
				// fmt.Printf("%v/%v\n", parsed, commons.FilesQuantity-noReadFilesCount)
			}
		}
	}
	wg.Done()
	green := color.New(color.Bold, color.BgGreen).SprintFunc()
	fmt.Printf("%s\n", green(" Mails parsing: COMPLETE "))
}

func SendMail() {
	addedMailCount := 0
	indexBulkSend := 0
	for {
		if len(chnMailStruct) > 0 {
			mail, ok := <-chnMailStruct
			if ok {
				jsonStr, _ := json.Marshal(mail)
				strBulk.WriteString(`{ "index" : { "_index" : "mails" } }` + "\n" + string(jsonStr) + "\n")
				if len(strBulk.String()) >= 100_000_000 || addedMailCount == (commons.FilesQuantity-noReadFilesCount-1) {
					LoadBulkV2(&strBulk) //Sending information to Zinc Search
					indexBulkSend += 1
					fmt.Printf("%s\n", yellow(fmt.Sprintf("%v - Bulk block: COMPLETE (%v/%v))", indexBulkSend, addedMailCount, commons.FilesQuantity-noReadFilesCount)))
				}
				addedMailCount += 1
			}
		}
		if len(chnMailStruct) == 0 && addedMailCount == (commons.FilesQuantity-noReadFilesCount) {
			green := color.New(color.Bold, color.BgGreen).SprintFunc()
			fmt.Printf("%s\n", green(" Mails sending: COMPLETE "))
			addedMailCount = 0
		}
	}
}

func DataLoader() {

	var currentSize int64 = 0
	var totalConter int64 = 0
	var batchCounter int64 = 1
	var files_path []string
	var body string

	cyan := color.New(color.Bold, color.BgCyan).SprintFunc()

	filepath.WalkDir(commons.DbPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if totalConter >= commons.Limit && commons.Limit >= 0 {
			return errors.New("set limit for files achive")
		}

		if currentSize >= commons.MaxSizeBatch {
			batchCounter += 1
			fmt.Printf(
				"%s - %s\n",
				cyan(strconv.FormatInt(batchCounter, 10)),
				cyan(strconv.FormatInt(totalConter, 10)))
			sendBody(&files_path, &body)
			currentSize = 0
		}

		if !d.IsDir() {
			totalConter += 1
			files_path = append(files_path, path)
			info, _ := os.Stat(path)

			currentSize += info.Size()
			email_content := GetEmailFileString(path)
			mail, _ := DataRegexParser(email_content)
			jsonStr, err := json.Marshal(mail)

			if err != nil {
				log.Fatal(err)
			}
			body += `{ "index" : { "_index" : "mails" } }` + "\n" + string(jsonStr) + "\n"
		}

		return nil
	})
	if len(files_path) > 0 {
		sendBody(&files_path, &body)
		batchCounter += 1
		fmt.Printf(
			"%s - %s\n",
			cyan(strconv.FormatInt(batchCounter, 10)),
			cyan(strconv.FormatInt(totalConter, 10)))
	}
}

func sendBody(files *[]string, body *string) {
	LoadBulk(*body)
	*body = ""
	*files = []string{}
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
