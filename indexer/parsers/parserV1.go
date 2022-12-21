package parsers

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

func DataParserV1(base_path string, limit int) string {
	var result string
	files, err := os.ReadDir(base_path) //dirs of all users
	Check_error(err)

	var i int
	for _, file := range files {
		if file.IsDir() {
			path := base_path + "/" + file.Name()

			email_dirs := GetListFiles(path)

			for _, email_dir := range email_dirs {
				emails_files_path := path + "/" + email_dir.Name()
				emails_files := GetListFiles(emails_files_path)

				for _, email_file := range emails_files {
					//===============================================================================
					email_content := GetEmailFileString(emails_files_path + "/" + email_file.Name())

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

							mail[strings.Trim(keys[i], ":")] = TrimValue(last_values[0])
							content := last_values[1:]

							mail["Content"] = strings.TrimLeft(TrimValue(strings.Join(content, "\n")), "\n\t")
						} else {
							from := base_index + len(keys[i])
							next_base_index := strings.Index(email_content, keys[i+1])
							if from < next_base_index {
								mail[strings.Trim(keys[i], ":")] = TrimValue(email_content[from:next_base_index])
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
