package commons

import (
	"io/fs"
	"net/http"
)

var HttpClient *http.Client

var Limit, MaxSizeBatch int64
var DbPath, Auth string
var NewData bool
var MailInfo []fs.FileInfo
