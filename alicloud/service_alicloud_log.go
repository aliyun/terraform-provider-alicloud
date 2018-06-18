package alicloud

import (
	"fmt"

	"github.com/aliyun/aliyun-log-go-sdk"
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
