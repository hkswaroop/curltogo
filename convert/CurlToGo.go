package convert

import (
	"errors"
	"fmt"
	"strings"
)

type parts struct {
	url     string
	method  string
	headers []string
	data    []string
}

func CurlToGo(input string) {

	fmt.Println(input)

	parsedData, err := parseInput(input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(parsedData)
	fmt.Println()
	fmt.Println(constructCode(parsedData))
}

//curl https://example.org/api -u usern:userpwd -d \"personal information\"
//curl http://www.example.org:1234/
func parseInput(input string) (parts, error) {

	data := parts{}
	if len(input) < 6 || !strings.HasPrefix(input, "curl") {
		return data, errors.New("Invalid input.")
	}
	split := strings.Split(input, " ")
	if len(split) < 2 {
		return data, errors.New("Invalid input.")
	}
	for i := 0; i < len(split); i++ {
		if split[i] == "curl" {
			continue
		}
		if strings.HasPrefix(split[i], "http") {
			data.url = split[i]
		}

	}
	if len(data.method) == 0 {
		data.method = "GET"
	}

	return data, nil
}

var errBlockStr = "if err != nil {\n\t// handle error\n}\n"
var deferStr = "defer resp.Body.Close()\n"

func constructCode(data parts) string {
	var code string
	if data.method == "GET" || data.method == "HEAD" {
		code = fmt.Sprintf("resp, err := http.Get(%s)\n%s%s", data.url, errBlockStr, deferStr)
	}
	return code
}
