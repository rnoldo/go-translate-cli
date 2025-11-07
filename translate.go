package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/codegangsta/cli"
)

const (
	defaultSource = "en"
	defaultTarget = "zh-CN"
	apiEndpoint   = "http://translate.google.com/translate_a/t"
)

var chineseRegex = regexp.MustCompile(`[\p{Han}]`)

func main() {
	app := cli.NewApp()
	app.Name = "translate"
	app.Usage = "translate is a cli tools for translation written in go and cli"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "from, f",
			Usage: "source language code (default auto-detect)",
		},
		cli.StringFlag{
			Name:  "to, t",
			Usage: "target language code (default auto-detect counterpart)",
		},
	}

	app.Action = func(c *cli.Context) error {
		text := strings.TrimSpace(strings.Join(c.Args(), " "))
		if text == "" {
			return cli.NewExitError("please provide text to translate", 1)
		}

		source := c.String("from")
		target := c.String("to")

		if source == "" && target == "" {
			// Fallback to the old auto direction if the user did not specify languages.
			if containsChinese(text) {
				source = "zh-CN"
				target = "en"
			} else {
				source = defaultSource
				target = defaultTarget
			}
		} else {
			if source == "" {
				source = "auto"
			}
			if target == "" {
				target = defaultTarget
			}
		}

		result, err := translate(text, source, target)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		fmt.Println(result)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func containsChinese(text string) bool {
	return chineseRegex.MatchString(text)
}

func translate(text, source, target string) (string, error) {
	params := url.Values{
		"client": {"t"},
		"ie":     {"UTF-8"},
		"oe":     {"UTF-8"},
		"hl":     {"zh-CN"},
		"sl":     {source},
		"tl":     {target},
		"text":   {text},
	}

	reqURL := apiEndpoint + "?" + params.Encode()
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(reqURL)
	if err != nil {
		return "", fmt.Errorf("request translation failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("translation api returned %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read translation response: %w", err)
	}

	return extractTranslation(body)
}

func extractTranslation(body []byte) (string, error) {
	var payload interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", fmt.Errorf("parse translation response: %w", err)
	}

	root, ok := payload.([]interface{})
	if !ok || len(root) == 0 {
		return "", fmt.Errorf("unexpected translation payload")
	}

	sentences, ok := root[0].([]interface{})
	if !ok || len(sentences) == 0 {
		return "", fmt.Errorf("translation payload missing sentences")
	}

	var builder strings.Builder
	for _, sentence := range sentences {
		parts, ok := sentence.([]interface{})
		if !ok || len(parts) == 0 {
			continue
		}

		chunk, _ := parts[0].(string)
		builder.WriteString(chunk)
	}

	if builder.Len() == 0 {
		return "", fmt.Errorf("translation payload did not contain text")
	}

	return builder.String(), nil
}
