package alicloud

import (
	"github.com/denverdino/aliyungo/slb"
)

type Listener struct {
	slb.HTTPListenerType

	InstancePort     int
	LoadBalancerPort int
	Protocol         string
	//tcp & udp
	PersistenceTimeout int

	//https
	SSLCertificateId string

	//tcp
	HealthCheckType slb.HealthCheckType

	//api interface: http & https is HealthCheckTimeout, tcp & udp is HealthCheckConnectTimeout
	HealthCheckConnectTimeout int
}

type ListenerErr struct {
	ErrType string
	Err     error
}

func (e *ListenerErr) Error() string {
	return e.ErrType + " " + e.Err.Error()

}

func expandBackendServers(list []interface{}) []slb.BackendServerType {
	result := make([]slb.BackendServerType, 0, len(list))
	for _, i := range list {
		if i.(string) != "" {
			result = append(result, slb.BackendServerType{ServerId: i.(string), Weight: 100})
		}
	}
	return result
}
