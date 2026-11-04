package privatezone

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
)

var throttlingUserCatcher = helper.Catcher{Reason: errs.ThrottlingUser, RetryCount: 30, RetryWaitSeconds: 2}
var systemBusyCatcher = helper.Catcher{Reason: "System.Busy", RetryCount: 30, RetryWaitSeconds: 5}

func invoker() helper.Invoker {
	i := helper.Invoker{}
	i.AddCatcher(throttlingUserCatcher)
	i.AddCatcher(helper.ServiceBusyCatcher)
	i.AddCatcher(systemBusyCatcher)
	return i
}
