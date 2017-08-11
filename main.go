package main

import (
	"github.com/hkswaroop/curltogo/convert"
)

// "curl http://www.example.org:1234/"
// "curl https://example.org/api -u usern:userpwd -d \"personal information\""
func main() {
	//var url = "curl http://www.example.org:1234/"
	var url = "curl https://example.org/api -u usern:userpwd -d \"personal information\""
	convert.CurlToGo(url)
}
