package main

import (
	"github.com/hkswaroop/curltogo/convert"
)

// "curl http://www.example.org:1234/"
// "curl https://example.org/api -u usern:userpwd -d \"personal information\""
func main() {
	var TYPE = "curl http://www.example.org:1234/"

	convert.CurlToGo(TYPE)
}
