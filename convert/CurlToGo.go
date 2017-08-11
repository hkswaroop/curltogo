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
	auth 	auth
}

type auth struct {
	user string
	pass string
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
		if strings.HasPrefix(split[i], "-") {
			data, i, _ = data.handleCommand(split, i)
		}
	}

	return data, nil
}

func (p parts) handleCommand(split []string, i int) (parts, int, error) {
	if len(split) < i +1 {
		return p, i, errors.New("Wrong curl")
	}
	if split[i] == "-u"{
		spl := strings.Split(split[i+1], ":")
		if len(spl) == 2 {
			p.auth = auth{}
			p.auth.user = spl[0]
			p.auth.pass = spl[1]
			i+=1
			if len(p.method) == 0 {
				p.method += "POST"
			}
		}
	} else if split[i] == "-d" {
		p.data = append(p.data, split[i+1])
		i+=1
		if len(p.method) == 0 {
			p.method += "POST"
		}
	}
	return p, i, nil
}

var errBlockStr = "if err != nil {\n\t// handle error\n}\n"
var deferStr = "defer resp.Body.Close()\n"

func constructCode(data parts) string {
	if len(data.method) == 0 {
		data.method = "GET"
	}
	var code string
	if data.method == "GET" || data.method == "HEAD" {
		code = fmt.Sprintf("resp, err := http.Get(%s)\n%s%s", data.url, errBlockStr, deferStr)
	} else {
		var body string
		for  i:=0;i<len(data.data);i++ {
			if i != 0 {
				body +="&"
			}
			body += data.data[i]
		}
		code = fmt.Sprintf("body := strings.NewReader(`%s`)\n",body)
		code += fmt.Sprintf("req, err := http.NewRequest(\"%s\", \"%s\", body)\n", data.method, data.url)
		code += errBlockStr
		if len(data.auth.user) != 0 || len(data.auth.pass) != 0 {
			code += fmt.Sprintf("req.SetBasicAuth(\"%s\", \"%s\")\n", data.auth.user, data.auth.pass)
		}
		if len(data.headers) == 0 {
			code += "req.Header.Set(\"Content-Type\", \"application/x-www-form-urlencoded\")\n"
		}
		code += "resp, err := http.DefaultClient.Do(req)\n" + errBlockStr + deferStr
	}
	return code
}