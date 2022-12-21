package parsers

import (
	"encoding/json"
	"errors"
	"log"
	"regexp"
	"strings"
)

func DataParserV2(base_path string, limit int) string {
	var result string
	files := ListFiles(base_path, limit)

	for _, file := range files {
		email_content := GetEmailFileString(file)
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
