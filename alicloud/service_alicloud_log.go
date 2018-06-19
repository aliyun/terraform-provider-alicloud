package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform/helper/resource"
)

func (client *AliyunClient) DescribeLogProject(name string) (project *sls.LogProject, err error) {
	project, err = client.logconn.GetProject(name)
	if err != nil {
		return project, fmt.Errorf("GetProject %s got an error: %#v.", name, err)
	}
	if project == nil || project.Name == "" {
		return project, GetNotFoundErrorFromString(GetNotFoundMessage("Log Project", name))
	}
	return
}

func (client *AliyunClient) DescribeLogStore(projectName, name string) (store *sls.LogStore, err error) {
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		store, err = client.logconn.GetLogStore(projectName, name)
		if err != nil {
			if IsExceptedErrors(err, []string{ProjectNotExist, LogStoreNotExist}) {
				return resource.NonRetryableError(GetNotFoundErrorFromString(GetNotFoundMessage("Log Store", name)))
			}
			if IsExceptedErrors(err, []string{InternalServerError}) {
				return resource.RetryableError(fmt.Errorf("GetLogStore %s got an error: %#v.", name, err))
			}
			return resource.NonRetryableError(fmt.Errorf("GetLogStore %s got an error: %#v.", name, err))
		}
		return nil
	})

	if err != nil {
		return
	}

	if store == nil || store.Name == "" {
		return store, GetNotFoundErrorFromString(GetNotFoundMessage("Log Store", name))
	}
	return
}
