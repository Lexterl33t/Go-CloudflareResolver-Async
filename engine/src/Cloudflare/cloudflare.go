package Cloudflare

import (
	"CldResolver/engine/src/Http"
	"bufio"
	"fmt"
	"net"
	"os"
)

func Readlines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func SendRequest(ch *chan any, data string, url string) {

	_, err := Http.GetRequest(fmt.Sprintf("https://%v.%v", data, url))
	fmt.Println(err)
	if err != nil {
		*ch <- map[string]any{
			"Domain":     fmt.Sprintf("https://%v.%v", data, url),
			"IP":         "No host",
			"Cloudflare": false,
			"Error":      err,
		}
		return
	}

	addr, err := net.LookupIP(fmt.Sprintf("%v.%v", data, url))
	if err != nil {
		*ch <- map[string]any{
			"Domain":     fmt.Sprintf("https://%v.%v", data, url),
			"IP":         "No host",
			"Cloudflare": false,
			"Error":      err,
		}
		return
	}
	*ch <- map[string]any{
		"Domain":     fmt.Sprintf("https://%v.%v", data, url),
		"IP":         addr,
		"Cloudflare": true,
		"Error":      nil,
	}
	return
}

func ExtractMap(data chan map[string]any) bool {
	select {
	case msg := <-data:
		if msg["Error"] == nil {
			return true
		} else {
			return false
		}
	}
}

func Resolve(url string, wordlist string) ([]any, error) {
	content, err := Readlines(wordlist)
	if err != nil {
		return nil, err
	}

	ch := make(chan any, len(content))

	for i := 0; i < len(content); i++ {
		go SendRequest(&ch, content[i], url)
	}

	var res []any

	for i := 0; i < len(content); i++ {
		select {
		case msg := <-ch:
			if msg.(map[string]any)["Error"] == nil {
				//reg := regexp.MustCompile(`103\\.21\\.24[4-7]\\.[0-9]{1,3}|103\\.22\\.20[0-3]\\.[0-9]{1,3}|103\\.31\\.[4-7]\\.[0-9]{1,3}|104\\.(1[6-9]|2[0-9]|3[0-1])\\.[0-9]{1,3}\\.[0-9]{1,3}|108\\.162\\.(19[2-9]|2[0-5][0-9])\\.[0-9]{1,3}|131\\.0\\.7[2-5]\\.[0-9]{1,3}|141\\.101\\.(6[4-9]|[7-9][0-9]|1[0-1][0-9]|12[0-7])\\.[0-9]{1,3}|162\\.15[8-9]\\.[0-9]{1,3}\\.[0-9]{1,3}|172\\.(6[4-9]|7[0-1])\\.[0-9]{1,3}\\.[0-9]{1,3}|173\\.245\\.(4[8-9]|5[0-9]|6[0-3])\\.[0-9]{1,3}|188\\.114\\.(9[6-9]|10[0-9]|11[0-1])\\.[0-9]{1,3}|190\\.93\\.2[4-5][0-9]\\.[0-9]{1,3}|197\\.234\\.24[0-3]\\.[0-9]{1,3}|198\\.41\\.(12[8-9]|1[3-9][0-9]|2[0-5][0-9])\\.[0-9]{1,3}|199\\.27\\.(12[8-9]|13[0-5])\\.[0-9]{1,3}`)
				res = append(res, msg)
			}
		}
	}

	return res, nil
}
