package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "translate"
	app.Usage = "translate is a cli tools for translation written in go and cli"

	//default action
	app.Action = func(c *cli.Context) {
		text := c.Args().First()

		u, _ := url.Parse("http://translate.google.com/translate_a/t?client=t&ie=UTF-8&oe=UTF-8&hl=zh-CN&sl=en&tl=zh-CN")
		q := u.Query()
		q.Add("text", text)

		//匹配中文
		re, _ := regexp.Compile("[\u4e00-\u9fa5]")
		if re.MatchString(text) {
			q.Set("sl", "zh-CN")
			q.Set("tl", "en")
		}

		u.RawQuery = q.Encode()

		response, err := http.Get(u.String())
		if err != nil {
			fmt.Println("http get error", err)
		}
		defer response.Body.Close()

		body, _ := ioutil.ReadAll(response.Body)
		bodystr := string(body)
		fmt.Println(bodystr[strings.Index(bodystr, "[[[\"")+4 : strings.Index(bodystr, "\",")])
	}
	app.Run(os.Args)
}
