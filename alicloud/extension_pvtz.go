package alicloud

var PvtzThrottlingUserCatcher = Catcher{Reason: ThrottlingUser, RetryCount: 30, RetryWaitSeconds: 2}
var PvtzSystemBusyCatcher = Catcher{Reason: "System.Busy", RetryCount: 30, RetryWaitSeconds: 5}

func PvtzInvoker() Invoker {
	i := Invoker{}
	i.AddCatcher(PvtzThrottlingUserCatcher)
	i.AddCatcher(ServiceBusyCatcher)
	i.AddCatcher(PvtzSystemBusyCatcher)
	return i
}
