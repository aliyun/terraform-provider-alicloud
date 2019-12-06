package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/denverdino/aliyungo/cs"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type CsService struct {
	client *connectivity.AliyunClient
}

func (s *CsService) GetContainerClusterByName(name string) (cluster cs.ClusterType, err error) {
	name = Trim(name)
	invoker := NewInvoker()
	var clusters []cs.ClusterType
	err = invoker.Run(func() error {
		raw, e := s.client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return csClient.DescribeClusters(name)
		})
		if e != nil {
			return e
		}
		clusters, _ = raw.([]cs.ClusterType)
		return nil
	})

	if err != nil {
		return cluster, fmt.Errorf("Describe cluster failed by name %s: %#v.", name, err)
	}

	if len(clusters) < 1 {
		return cluster, GetNotFoundErrorFromString(GetNotFoundMessage("Container Cluster", name))
	}

	for _, c := range clusters {
		if c.Name == name {
			return c, nil
		}
	}
	return cluster, GetNotFoundErrorFromString(GetNotFoundMessage("Container Cluster", name))
}

func (s *CsService) GetContainerClusterAndCertsByName(name string) (*cs.ClusterType, *cs.ClusterCerts, error) {
	cluster, err := s.GetContainerClusterByName(name)
	if err != nil {
		return nil, nil, err
	}
	var certs cs.ClusterCerts
	invoker := NewInvoker()
	err = invoker.Run(func() error {
		raw, e := s.client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return csClient.GetClusterCerts(cluster.ClusterID)
		})
		if e != nil {
			return e
		}
		certs, _ = raw.(cs.ClusterCerts)
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return &cluster, &certs, nil
}

func (s *CsService) DescribeContainerApplication(clusterName, appName string) (app cs.GetProjectResponse, err error) {
	appName = Trim(appName)
	cluster, certs, err := s.GetContainerClusterAndCertsByName(clusterName)
	if err != nil {
		return app, err
	}
	raw, err := s.client.WithCsProjectClient(cluster.ClusterID, cluster.MasterURL, *certs, func(csProjectClient *cs.ProjectClient) (interface{}, error) {
		return csProjectClient.GetProject(appName)
	})
	app, _ = raw.(cs.GetProjectResponse)
	if err != nil {
		if IsExceptedError(err, ApplicationNotFound) {
			return app, GetNotFoundErrorFromString(GetNotFoundMessage("Container Application", appName))
		}
		return app, fmt.Errorf("Getting Application failed by name %s: %#v.", appName, err)
	}
	if app.Name != appName {
		return app, GetNotFoundErrorFromString(GetNotFoundMessage("Container Application", appName))
	}
	return
}

func (s *CsService) WaitForContainerApplication(clusterName, appName string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		app, err := s.DescribeContainerApplication(clusterName, appName)
		if err != nil {
			return err
		}

		if strings.ToLower(app.CurrentState) == strings.ToLower(string(status)) {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(fmt.Sprintf("Waitting for container application %s is timeout and current status is %s.", string(status), app.CurrentState))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *CsService) DescribeCsKubernetes(id string) (cluster cs.KubernetesCluster, err error) {
	invoker := NewInvoker()
	var requestInfo *cs.Client
	var response interface{}

	if err := invoker.Run(func() error {
		raw, err := s.client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			requestInfo = csClient
			return csClient.DescribeKubernetesCluster(id)
		})
		response = raw
		return err
	}); err != nil {
		if NotFoundError(err) || IsExceptedError(err, ErrorClusterNotFound) {
			return cluster, WrapErrorf(err, NotFoundMsg, DenverdinoAliyungo)
		}
		return cluster, WrapErrorf(err, DefaultErrorMsg, id, "DescribeKubernetesCluster", DenverdinoAliyungo)
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["ClusterId"] = id
		addDebug("DescribeKubernetesCluster", response, requestInfo, requestMap)
	}
	cluster, _ = response.(cs.KubernetesCluster)
	if cluster.ClusterID != id {
		return cluster, WrapErrorf(Error(GetNotFoundMessage("CsKubernetes", id)), NotFoundMsg, ProviderERROR)
	}
	return
}

func (s *CsService) WaitForCsKubernetes(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeCsKubernetes(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.ClusterID == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.ClusterID, id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)

	}
}

func (s *CsService) DescribeCsManagedKubernetes(id string) (cluster cs.KubernetesCluster, err error) {
	var requestInfo *cs.Client
	invoker := NewInvoker()
	var response interface{}

	if err := invoker.Run(func() error {
		raw, err := s.client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			requestInfo = csClient
			return csClient.DescribeKubernetesCluster(id)
		})
		response = raw
		return err
	}); err != nil {
		if NotFoundError(err) || IsExceptedError(err, ErrorClusterNotFound) {
			return cluster, WrapErrorf(err, NotFoundMsg, AlibabaCloudSdkGoERROR)
		}
		return cluster, WrapErrorf(err, DefaultErrorMsg, id, "DescribeKubernetesCluster", DenverdinoAliyungo)
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["Id"] = id
		addDebug("DescribeKubernetesCluster", response, requestInfo, requestMap, map[string]interface{}{"Id": id})
	}
	cluster, _ = response.(cs.KubernetesCluster)
	if cluster.ClusterID != id {
		return cluster, WrapErrorf(Error(GetNotFoundMessage("CSManagedKubernetes", id)), NotFoundMsg, ProviderERROR)
	}
	return

}

func (s *CsService) WaitForCSManagedKubernetes(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeCsManagedKubernetes(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.ClusterID == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.ClusterID, id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)

	}
}

func (s *CsService) CsKubernetesInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCsKubernetes(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if string(object.State) == failState {
				return object, string(object.State), WrapError(Error(FailedToReachTargetStatus, string(object.State)))
			}
		}
		return object, string(object.State), nil
	}
}

func (s *CsService) CsManagedKubernetesInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCsManagedKubernetes(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if string(object.State) == failState {
				return object, string(object.State), WrapError(Error(FailedToReachTargetStatus, string(object.State)))
			}
		}
		return object, string(object.State), nil
	}
}

func (s *CsService) CsServerlessKubernetesInstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeCsServerlessKubernetes(id)
		if err != nil {
			if NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		for _, failState := range failStates {
			if string(object.State) == failState {
				return object, string(object.State), WrapError(Error(FailedToReachTargetStatus, string(object.State)))
			}
		}
		return object, string(object.State), nil
	}
}

func (s *CsService) DescribeCsServerlessKubernetes(id string) (*cs.ServerlessClusterResponse, error) {
	cluster := &cs.ServerlessClusterResponse{}
	var requestInfo *cs.Client
	invoker := NewInvoker()
	var response interface{}

	if err := invoker.Run(func() error {
		raw, err := s.client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			requestInfo = csClient
			return csClient.DescribeServerlessKubernetesCluster(id)
		})
		response = raw
		return err
	}); err != nil {
		if IsExceptedError(err, ErrorClusterNotFound) {
			return cluster, WrapErrorf(err, NotFoundMsg, DenverdinoAliyungo)
		}
		return cluster, WrapErrorf(err, DefaultErrorMsg, id, "DescribeServerlessKubernetesCluster", DenverdinoAliyungo)
	}
	if debugOn() {
		requestMap := make(map[string]interface{})
		requestMap["Id"] = id
		addDebug("DescribeServerlessKubernetesCluster", response, requestInfo, requestMap, map[string]interface{}{"Id": id})
	}
	cluster, _ = response.(*cs.ServerlessClusterResponse)
	if cluster != nil && cluster.ClusterId != id {
		return cluster, WrapErrorf(Error(GetNotFoundMessage("CSServerlessKubernetes", id)), NotFoundMsg, ProviderERROR)
	}
	return cluster, nil

}

func (s *CsService) WaitForCSServerlessKubernetes(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for {
		object, err := s.DescribeCsServerlessKubernetes(id)
		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return WrapError(err)
			}
		}
		if object.ClusterId == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.ClusterId, id, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)

	}
}

func (s *CsService) tagsToMap(tags []cs.Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		if !s.ignoreTag(t) {
			result[t.Key] = t.Value
		}
	}
	return result
}

func (s *CsService) ignoreTag(t cs.Tag) bool {
	filter := []string{"^http://", "^https://"}
	for _, v := range filter {
		log.Printf("[DEBUG] Matching prefix %v with %v\n", v, t.Key)
		ok, _ := regexp.MatchString(v, t.Key)
		if ok {
			log.Printf("[DEBUG] Found Alibaba Cloud specific t %s (val: %s), ignoring.\n", t.Key, t.Value)
			return true
		}
	}
	return false
}
