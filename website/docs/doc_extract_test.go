package docs

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

func TestDocumentCheck(t *testing.T) {
	DocType := "datasource"
	fileRegex := regexp.MustCompile("([a-zA-Z0-9_]*).html.markdown")
	if DocType == "resource" {
		files, _ := ioutil.ReadDir("./r")
		for _, file := range files {
			log.Infof("The File %s being processed", file.Name())
			fieldNameMatched := fileRegex.FindAllStringSubmatch(file.Name(), -1)
			fileName := fieldNameMatched[0][1]
			// extract Terraform template
			cmd := exec.Command("terrafmt", "blocks", fmt.Sprintf("./r/%s", file.Name()), "-j")
			out, err := cmd.CombinedOutput()
			if err != nil {
				log.Fatalf("cmd.Run() failed with %s\n", err)
			}
			res := make(map[string]interface{}, 0)
			err = json.Unmarshal(out, &res)
			assert.Nil(t, err)
			// output the target path
			for i := 0; i < int(res["block_count"].(float64)); i++ {
				err := os.Mkdir(fmt.Sprintf("./target/r/%s_%d", fileName, i), os.ModePerm)
				if err != nil {
					log.Errorf("Output The FIle failed,%s", err.Error())
					if !strings.Contains(err.Error(), "file exists") {
						t.Fatal()
					}
				}
				codeObj := res["blocks"].([]interface{})[i]
				err = OutputTemplate(fmt.Sprintf("./target/r/%s_%d/%s_%d.tf", fileName, i, fileName, i), codeObj.(map[string]interface{})["text"].(string))
				assert.Nil(t, err)
			}
		}
	}
	if DocType == "datasource" {
		files, _ := ioutil.ReadDir("./d")
		for _, file := range files {
			log.Infof("The File %s being processed", file.Name())
			fieldNameMatched := fileRegex.FindAllStringSubmatch(file.Name(), -1)
			fileName := fieldNameMatched[0][1]
			// extract Terraform template
			cmd := exec.Command("terrafmt", "blocks", fmt.Sprintf("./d/%s", file.Name()), "-j")
			out, err := cmd.CombinedOutput()
			if err != nil {
				log.Fatalf("cmd.Run() failed with %s\n", err)
			}
			res := make(map[string]interface{}, 0)
			err = json.Unmarshal(out, &res)
			assert.Nil(t, err)
			// output the target path
			for i := 0; i < int(res["block_count"].(float64)); i++ {
				err := os.Mkdir(fmt.Sprintf("./target/d/%s_%d", fileName, i), os.ModePerm)
				if err != nil {
					log.Errorf("Output The FIle failed,%s", err.Error())
					if !strings.Contains(err.Error(), "file exists") {
						t.Fatal()
					}
				}
				codeObj := res["blocks"].([]interface{})[i]
				err = OutputTemplate(fmt.Sprintf("./target/d/%s_%d/%s_%d.tf", fileName, i, fileName, i), codeObj.(map[string]interface{})["text"].(string))
				assert.Nil(t, err)
			}
		}
	}
}

func OutputTemplate(targetPath, res string) error {
	f, err := os.Create(targetPath)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		_, err = f.Write([]byte(res))
	}
	return nil
}
