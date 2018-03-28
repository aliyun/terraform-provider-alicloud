package alicloud

import (
	"fmt"
	"strings"
	"time"

	"bytes"
	"log"

	"github.com/denverdino/aliyungo/cs"
	"gopkg.in/yaml.v2"
)

func ContainerApplicationTemplateAreEquivalent(tem1, tem2 string) (bool, error) {
	var obj1 interface{}
	err := yaml.Unmarshal([]byte(tem1), &obj1)
	if err != nil {
		return false, err
	}

	canonicalYaml1, _ := yaml.Marshal(obj1)

	var obj2 interface{}
	err = yaml.Unmarshal([]byte(tem2), &obj2)
	if err != nil {
		return false, err
	}

	canonicalYaml2, _ := yaml.Marshal(obj2)

	equal := bytes.Compare(canonicalYaml1, canonicalYaml2) == 0
	if !equal {
		log.Printf("[DEBUG] Canonical template are not equal.\nFirst: %s\nSecond: %s\n",
			canonicalYaml1, canonicalYaml2)
	}
	return equal, nil
}

func (client *AliyunClient) GetContainerClusterByName(name string) (cluster cs.ClusterType, err error) {
	name = Trim(name)
	clusters, err := client.csconn.DescribeClusters(name)
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

func (client *AliyunClient) GetApplicationClientByClusterName(name string) (c *cs.ProjectClient, err error) {
	cluster, err := client.GetContainerClusterByName(name)
	if err != nil {
		return nil, err
	}

	certs, err := client.csconn.GetClusterCerts(cluster.ClusterID)
	if err != nil {
		return
	}

	c, err = cs.NewProjectClient(cluster.ClusterID, cluster.MasterURL, certs)

	if err != nil {
		return nil, fmt.Errorf("Getting Application Client failed by cluster id %s: %#v.", cluster.ClusterID, err)
	}
	c.SetDebug(false)
	c.SetUserAgent(getUserAgent())

	return
}

func (client *AliyunClient) DescribeContainerApplication(clusterName, appName string) (app cs.GetProjectResponse, err error) {
	appName = Trim(appName)
	conn, err := client.GetApplicationClientByClusterName(clusterName)
	if err != nil {
		return app, err
	}
	app, err = conn.GetProject(appName)
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

func (client *AliyunClient) WaitForContainerApplication(clusterName, appName string, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		app, err := client.DescribeContainerApplication(clusterName, appName)
		if err != nil {
			return err
		}

		if strings.ToLower(app.CurrentState) == strings.ToLower(string(status)) {
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("Container Application", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}
