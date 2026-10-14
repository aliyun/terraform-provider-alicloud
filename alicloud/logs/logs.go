// Package logs hosts the provider's debug logging. AddDebug is the
// entry point every resource and data source calls after an API
// response, gated on the DEBUG environment variable. It moved out of
// package alicloud so service packages can reach it without importing
// the monolith. The trace walk starts at runtime.Caller(1) relative to
// AddDebug's own frame, so calling it through a var alias in package
// alicloud preserves the reported call sites exactly.
package logs

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/fc-go-sdk"
	"github.com/denverdino/aliyungo/cs"
)

const DefaultDebugMsg = "\n*************** %s Response *************** \n%s\n%s******************************\n\n"

func DebugOn() bool {
	for _, part := range strings.Split(os.Getenv("DEBUG"), ",") {
		if strings.TrimSpace(part) == "terraform" {
			return true
		}
	}
	return false
}

func AddDebug(action, content interface{}, requestInfo ...interface{}) {
	if DebugOn() {
		trace := "[DEBUG TRACE]:\n"
		for skip := 1; skip < 5; skip++ {
			_, filepath, line, _ := runtime.Caller(skip)
			trace += fmt.Sprintf("%s:%d\n", filepath, line)
		}

		if len(requestInfo) > 0 {
			var request = struct {
				Domain     string
				Version    string
				UserAgent  string
				ActionName string
				Method     string
				Product    string
				Region     string
				AK         string
			}{}
			switch requestInfo[0].(type) {
			case *requests.RpcRequest:
				tmp := requestInfo[0].(*requests.RpcRequest)
				request.Domain = tmp.GetDomain()
				request.Version = tmp.GetVersion()
				request.ActionName = tmp.GetActionName()
				request.Method = tmp.GetMethod()
				request.Product = tmp.GetProduct()
				request.Region = tmp.GetRegionId()
			case *requests.RoaRequest:
				tmp := requestInfo[0].(*requests.RoaRequest)
				request.Domain = tmp.GetDomain()
				request.Version = tmp.GetVersion()
				request.ActionName = tmp.GetActionName()
				request.Method = tmp.GetMethod()
				request.Product = tmp.GetProduct()
				request.Region = tmp.GetRegionId()
			case *requests.CommonRequest:
				tmp := requestInfo[0].(*requests.CommonRequest)
				request.Domain = tmp.GetDomain()
				request.Version = tmp.GetVersion()
				request.ActionName = tmp.GetActionName()
				request.Method = tmp.GetMethod()
				request.Product = tmp.GetProduct()
				request.Region = tmp.GetRegionId()
			case *fc.Client:
				client := requestInfo[0].(*fc.Client)
				request.Version = client.Config.APIVersion
				request.Product = "FC"
				request.ActionName = fmt.Sprintf("%s", action)
			case *sls.Client:
				request.Product = "LOG"
				request.ActionName = fmt.Sprintf("%s", action)
			case *tablestore.TableStoreClient:
				request.Product = "OTS"
				request.ActionName = fmt.Sprintf("%s", action)
			case *oss.Client:
				request.Product = "OSS"
				request.ActionName = fmt.Sprintf("%s", action)
			case *datahub.DataHub:
				request.Product = "DataHub"
				request.ActionName = fmt.Sprintf("%s", action)
			case *cs.Client:
				request.Product = "CS"
				request.ActionName = fmt.Sprintf("%s", action)
			}

			requestContent := ""
			if len(requestInfo) > 1 {
				switch requestInfo[1].(type) {
				case *tea.SDKError:
					requestContent = fmt.Sprintf("%#v", requestInfo[1].(*tea.SDKError).Error())
				default:
					requestContent = fmt.Sprintf("%#v", requestInfo[1])
				}
			}

			if len(requestInfo) == 1 {
				if v, ok := requestInfo[0].(map[string]interface{}); ok {
					if res, err := json.Marshal(&v); err == nil {
						requestContent = string(res)
					}
					if res, err := json.Marshal(&content); err == nil {
						content = string(res)
					}
				}
			}

			content = fmt.Sprintf("%vDomain:%v, Version:%v, ActionName:%v, Method:%v, Product:%v, Region:%v\n\n"+
				"*************** %s Request ***************\n%#v\n",
				content, request.Domain, request.Version, request.ActionName,
				request.Method, request.Product, request.Region, request.ActionName, requestContent)
		}

		log.Printf(DefaultDebugMsg, action, content, trace)
	}
}
