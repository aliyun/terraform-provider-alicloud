package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	stateFilePath := strings.TrimSpace(os.Args[1])
	stateContent, err := os.ReadFile(stateFilePath + "/terraform.tfstate")
	if err != nil {
		log.Println("reading the state file failed, error: ", err)
		os.Exit(1)
	}

	resourceState := new(TerraformState)
	err = json.Unmarshal(stateContent, resourceState)
	if err != nil {
		log.Println("unmarshalling the state content failed, error: ", err)
		os.Exit(1)
	}
	var appendBuf bytes.Buffer
	for _, res := range resourceState.Resources {
		if res.Mode != "managed" || !strings.HasPrefix(res.Type, "alicloud_") {
			continue
		}
		if _, ok := importUnsupportedResources[res.Type]; ok {
			log.Printf("[WARNING] %s does not support import.", res.Type)
			continue
		}
		for _, instance := range res.Instances {
			item := instance.(map[string]interface{})
			id := item["attributes"].(map[string]interface{})["id"]
			to := res.Type + "." + res.Name
			// -generate-config-out does not support count and for_each
			// some example depends on specified region
			to += "-import"
			if v, ok := item["index_key"]; ok {
				//if vv, okk := v.(string) ; okk {
				//	to += fmt.Sprintf("[\"%v\"]", vv)
				//}else {
				//	to += fmt.Sprintf("[%v]", v)
				//}
				if vv, okk := v.(string); okk {
					bufs := make([]byte, 0)
					for _, ch := range vv {
						if (ch > 47 && ch < 58) || (ch > 64 && ch < 91) || (ch > 96 && ch < 123) || ch == 45 || ch == 95 {
							bufs = append(bufs, byte(ch))
							continue
						}
						bufs = append(bufs, '-')
					}
					to += "-" + string(bufs)
				} else {
					to += fmt.Sprintf("-%v", v)
				}
			}
			appendBuf.WriteString(fmt.Sprintf(`
import {
 id = "%s"
 to = %s
}

`, id, to))
		}
	}
	os.WriteFile(stateFilePath+"/import.tf", appendBuf.Bytes(), 0644)
}

type TerraformState struct {
	Version          int         `json:"version"`
	TerraformVersion string      `json:"terraform_version"`
	Serial           int         `json:"serial"`
	Lineage          string      `json:"lineage"`
	Outputs          interface{} `json:"outputs"`
	Resources        []Resource  `json: "resources"`
	CheckResults     interface{} `json:"check_results"`
}

type Resource struct {
	Mode      string        `json:"mode"`
	Type      string        `json: "type"`
	Name      string        `json:"name"`
	Provider  string        `json:"provider"`
	Instances []interface{} `json:"instances"`
}

var importUnsupportedResources = map[string]struct{}{
	"alicloud_ram_access_key":             {},
	"alicloud_ram_role_attachment":        {},
	"alicloud_maxcompute_project":         {},
	"alicloud_api_gateway_app_attachment": {},
	"alicloud_nas_smb_acl_attachment":     {},
	"alicloud_ots_instance_attachment":    {},
}
