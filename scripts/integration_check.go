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
	ossObjectPath := strings.TrimSpace(os.Args[1])
	invocationId := strings.Replace(ossObjectPath, "/", "_", -1)
	log.Println("trace id:", invocationId)
	urlPrefix := "https://terraform-provider-alicloud-ct.oss-ap-southeast-1.aliyuncs.com"
	runLogFileName := "terraform.run.log"
	runResultFileName := "terraform.run.result"
	runLogUrl := urlPrefix + "/" + ossObjectPath + "/" + runLogFileName
	runResultUrl := urlPrefix + "/" + ossObjectPath + "/" + runResultFileName
	lastLineNum := 0
	deadline := time.Now().Add(time.Duration(24) * time.Hour)
	finish := false
	exitCode := 0
	log.Printf("see integration test log: %s", runLogUrl)
	for !time.Now().After(deadline) {
		runLogResponse, err := http.Get(runLogUrl)
		if err != nil || runLogResponse.StatusCode != 200 {
			log.Println("waiting for job running...")
			time.Sleep(10 * time.Second)
			continue
			//os.Exit(1)
		}
		defer runLogResponse.Body.Close()

		runLogContent := make([]byte, 100000000)
		lineNum, er := runLogResponse.Body.Read(runLogContent)
		if er != nil && fmt.Sprint(er) != "EOF" {
			log.Println("[ERROR] reading run log response failed:", err)
		}
		if runLogResponse.StatusCode == 200 {
			fmt.Println(string(runLogContent[lastLineNum:lineNum]))
			lastLineNum = lineNum
		}
		if finish {
			os.Exit(exitCode)
		}
		runResultResponse, err := http.Get(runResultUrl)
		if err != nil || runResultResponse.StatusCode != 200 {
			//log.Println("waiting for job finished...")
			time.Sleep(5 * time.Second)
			continue
		}
		defer runResultResponse.Body.Close()
		runResultContent := make([]byte, 100000)
		_, err = runResultResponse.Body.Read(runResultContent)
		if err != nil && fmt.Sprint(er) != "EOF" {
			log.Println("[ERROR] reading run result response failed:", err)
		}
		finish = true
		if !strings.HasPrefix(string(runResultContent), "PASS") {
			exitCode = 1
		}
	}
	log.Println("[ERROR] Timeout: waiting for job finished timeout after 24 hours.")
}
