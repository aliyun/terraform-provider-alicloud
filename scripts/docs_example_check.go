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
	log.Println("run log path:", ossObjectPath)
	log.Println("trace id:", invocationId)
	urlPrefix := "https://terraform-provider-alicloud-ct-eu.oss-eu-central-1.aliyuncs.com"
	runLogFileName := "terraform-example.run.result.log"
	runLogUrl := urlPrefix + "/" + ossObjectPath + "/" + runLogFileName
	lastLineNum := 0
	deadline := time.Now().Add(time.Duration(24) * time.Hour)
	finish := false
	exitCode := 0
	for !time.Now().After(deadline) {
		runLogResponse, err := http.Get(runLogUrl)
		if err != nil || runLogResponse.StatusCode != 200 {
			log.Println("waiting for job running...")
			time.Sleep(10 * time.Second)
			continue
		}
		defer runLogResponse.Body.Close()

		runLogContent := make([]byte, 100000000)
		lineNum, er := runLogResponse.Body.Read(runLogContent)
		if er != nil && fmt.Sprint(er) != "EOF" {
			log.Println("[ERROR] reading run log response failed:", err)
		}
		if runLogResponse.StatusCode == 200 {
			if lineNum > lastLineNum {
				fmt.Println(string(runLogContent[lastLineNum:lineNum]))
				lastLineNum = lineNum
			}
		}
		if finish {
			log.Println("run log path:", ossObjectPath)
			log.Println("run log url:", runLogUrl)
			os.Exit(exitCode)
		}
		finish = true
		if strings.Contains(string(runLogContent), "FAIL") {
			exitCode = 1
		}
	}
	log.Println("[ERROR] Timeout: waiting for job finished timeout after 24 hours.")
}
