package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("[ERROR] missing head sha argument: usage: acube_check.go <headSha>")
	}
	sha := strings.TrimSpace(os.Args[1])
	log.Println("acube acctest head sha:", sha)
	urlPrefix := "https://terraform-acctest-code-bucket.oss-cn-beijing.aliyuncs.com"
	objectPath := "dst-pr-" + sha
	runLogUrl := urlPrefix + "/" + objectPath + "/run.log"
	runResultUrl := urlPrefix + "/" + objectPath + "/run.result"
	log.Println("run log url:", runLogUrl)
	log.Println("run result url:", runResultUrl)
	lastLineNum := 0
	deadline := time.Now().Add(30 * time.Minute)
	finish := false
	exitCode := 0
	for !time.Now().After(deadline) {
		runLogResponse, err := http.Get(runLogUrl)
		if err != nil || runLogResponse.StatusCode != 200 {
			if runLogResponse != nil {
				runLogResponse.Body.Close()
			}
			log.Println("waiting for acube job running...")
			time.Sleep(10 * time.Second)
			continue
		}
		runLogContent := make([]byte, 100000000)
		lineNum, er := runLogResponse.Body.Read(runLogContent)
		if er != nil && fmt.Sprint(er) != "EOF" {
			log.Println("[ERROR] reading run log response failed:", er)
		}
		runLogResponse.Body.Close()
		if lineNum > lastLineNum {
			fmt.Println(string(runLogContent[lastLineNum:lineNum]))
			lastLineNum = lineNum
		}
		if finish {
			log.Println("run log url:", runLogUrl)
			os.Exit(exitCode)
		}
		runResultResponse, err := http.Get(runResultUrl)
		if err != nil || runResultResponse.StatusCode != 200 {
			if runResultResponse != nil {
				runResultResponse.Body.Close()
			}
			time.Sleep(10 * time.Second)
			continue
		}
		runResultContent := make([]byte, 100000)
		_, err = runResultResponse.Body.Read(runResultContent)
		if err != nil && fmt.Sprint(err) != "EOF" {
			log.Println("[ERROR] reading run result response failed:", err)
		}
		runResultResponse.Body.Close()
		finish = true
		if !strings.HasPrefix(string(runResultContent), "SUCCESS") {
			exitCode = 1
		}
	}
	log.Println("[ERROR] Timeout: waiting for acube acctest finished timeout after 30 minutes.")
	os.Exit(1)
}
