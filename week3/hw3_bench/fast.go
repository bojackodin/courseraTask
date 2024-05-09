package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	easyjson "github.com/mailru/easyjson"
	// "log"
)

type UserInfo struct {
	Browsers []string `json:"browsers"`
	// Company  string   `json:"-"`
	// Country  string   `json:"-"`
	Email string `json:"email"`
	// Job      string   `json:"-"`
	Name string `json:"name"`
	// Phone    string   `json:"-"`
}

type Users []UserInfo

const filePath string = "./data/users.txt"

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	users := make([]UserInfo, 0, 1024)
	for scanner.Scan() {
		user := UserInfo{}
		err := easyjson.Unmarshal(scanner.Bytes(), &user)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}

	seenBrowsers := make(map[string]struct{})
	// uniqueBrowsers := 0
	foundUsers := ""
	var builder strings.Builder

	for i, user := range users {
		isAndroid := false
		isMSIE := false

		// browsers, ok := user["browsers"].([]interface{})
		// if !ok {
		// 	// log.Println("cant cast browsers")
		// 	continue
		// }

		for _, browser := range user.Browsers {
			// r, _ := regexp.Compile("Android")

			if strings.Contains(browser, "Android") {
				isAndroid = true
				// notSeenBefore := true
				// for _, item := range seenBrowsers {
				// 	if item == browser {
				// 		notSeenBefore = false
				// 	}
				// }
				if _, ok := seenBrowsers[browser]; !ok {
					seenBrowsers[browser] = struct{}{}
				}
				// if notSeenBefore {
				// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
				// seenBrowsers = append(seenBrowsers, browser)
				// uniqueBrowsers++
				// }
			}

			if strings.Contains(browser, "MSIE") {
				isMSIE = true
				if _, ok := seenBrowsers[browser]; !ok {
					seenBrowsers[browser] = struct{}{}
				}
			}
		}

		// for _, browser := range user.Browsers {
		// 	// r, _ := regexp.Compile("MSIE")
		// 	if ok := strings.Contains(browser, "MSIE"); ok {
		// 		isMSIE = true
		// 		// notSeenBefore := true
		// 		// for _, item := range seenBrowsers {
		// 		// 	if item == browser {
		// 		// 		notSeenBefore = false
		// 		// 	}
		// 		// }
		// 		if _, ok := seenBrowsers[browser]; !ok {
		// 			seenBrowsers[browser] = struct{}{}
		// 		}
		// 		// if notSeenBefore {
		// 		// 	// log.Printf("SLOW New browser: %s, first seen: %s", browser, user["name"])
		// 		// 	// seenBrowsers = append(seenBrowsers, browser)
		// 		// 	// uniqueBrowsers++
		// 		// 	seenBrowsers[browser] = struct{}{}
		// 		// }
		// 	}
		// }

		if !(isAndroid && isMSIE) {
			continue
		}

		// log.Println("Android and MSIE user:", user["name"], user["email"])
		email := strings.ReplaceAll(user.Email, "@", " [at] ")
		// foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, email)
		builder.WriteString(fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, email))
	}
	foundUsers = builder.String()

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
