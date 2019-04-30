package alicloud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

type SlbService struct {
	client *connectivity.AliyunClient
}

type SlbTag struct {
	TagKey   string
	TagValue string
}

const max_num_per_time = 50
const tags_max_num_per_time = 5
const tags_max_page_size = 50

func (s *SlbService) BuildSlbCommonRequest() (*requests.CommonRequest, error) {
	// Get product code from the built request
	slbReq := slb.CreateCreateLoadBalancerRequest()
	req, err := s.client.NewCommonRequest(slbReq.GetProduct(), slbReq.GetLocationServiceCode(), strings.ToUpper(string(Https)), connectivity.ApiVersion20140515)
	if err != nil {
		err = WrapError(err)
	}
	return req, err
}

func (s *SlbService) DescribeSLB(id string) (response *slb.DescribeLoadBalancerAttributeResponse, err error) {

	request := slb.CreateDescribeLoadBalancerAttributeRequest()
	request.LoadBalancerId = id
	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeLoadBalancerAttribute(request)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{LoadBalancerNotFound}) {
			err = WrapErrorf(Error(GetNotFoundMessage("SLB", id)), NotFoundMsg, AlibabaCloudSdkGoERROR)
		} else {
			err = WrapErrorf(err, DefaultErrorMsg, id, request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		return
	}
	addDebug(request.GetActionName(), raw)
	response, _ = raw.(*slb.DescribeLoadBalancerAttributeResponse)
	if response.LoadBalancerId == "" {
		err = WrapErrorf(Error(GetNotFoundMessage("SLB", id)), NotFoundMsg, ProviderERROR)
	}
	return
}

func (s *SlbService) DescribeLoadBalancerRuleId(slbId string, port int, domain, url string) (string, error) {
	req := slb.CreateDescribeRulesRequest()
	req.LoadBalancerId = slbId
	req.ListenerPort = requests.NewInteger(port)
	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeRules(req)
	})
	if err != nil {
		return "", WrapErrorf(err, DefaultErrorMsg, slbId, req.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(req.GetActionName(), raw)
	rules, _ := raw.(*slb.DescribeRulesResponse)
	for _, rule := range rules.Rules.Rule {
		if rule.Domain == domain && rule.Url == url {
			return rule.RuleId, nil
		}
	}

	return "", GetNotFoundErrorFromString(fmt.Sprintf("Rule is not found based on domain %s and url %s.", domain, url))
}

func (s *SlbService) DescribeLoadBalancerRuleAttribute(ruleId string) (*slb.DescribeRuleAttributeResponse, error) {
	req := slb.CreateDescribeRuleAttributeRequest()
	req.RuleId = ruleId
	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeRuleAttribute(req)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{InvalidRuleIdNotFound}) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("SLB Rule", ruleId))
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, ruleId, req.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(req.GetActionName(), raw)
	rule, _ := raw.(*slb.DescribeRuleAttributeResponse)
	if rule == nil || rule.LoadBalancerId == "" {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("SLB Rule", ruleId))
	}
	return rule, err
}

func (s *SlbService) DescribeSlbVServerGroupAttribute(groupId string) (*slb.DescribeVServerGroupAttributeResponse, error) {
	req := slb.CreateDescribeVServerGroupAttributeRequest()
	req.VServerGroupId = groupId
	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeVServerGroupAttribute(req)
	})
	if err != nil {
		if IsExceptedErrors(err, []string{VServerGroupNotFoundMessage, InvalidParameter}) {
			return nil, GetNotFoundErrorFromString(GetNotFoundMessage("SLB VServer Group", groupId))
		}
		return nil, fmt.Errorf("DescribeSlbVServerGroupAttribute got an error: %#v", err)
	}
	group, _ := raw.(*slb.DescribeVServerGroupAttributeResponse)
	if group == nil || group.VServerGroupId == "" {
		return nil, GetNotFoundErrorFromString(GetNotFoundMessage("SLB VServer Group", groupId))
	}
	return group, err
}

func (s *SlbService) DescribeLoadBalancerListenerAttribute(loadBalancerId string, port int, protocol Protocol) (listener map[string]interface{}, err error) {
	req, err := s.BuildSlbCommonRequest()
	if err != nil {
		err = WrapError(err)
		return
	}
	req.ApiName = fmt.Sprintf("DescribeLoadBalancer%sListenerAttribute", strings.ToUpper(string(protocol)))
	req.QueryParams["LoadBalancerId"] = loadBalancerId
	req.QueryParams["ListenerPort"] = string(requests.NewInteger(port))
	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.ProcessCommonRequest(req)
	})
	if err != nil {
		return
	}
	resp, _ := raw.(*responses.CommonResponse)
	if err = json.Unmarshal(resp.GetHttpContentBytes(), &listener); err != nil {
		err = fmt.Errorf("Unmarshalling body got an error: %#v.", err)
	}

	return

}

func (s *SlbService) WaitForSLB(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeSLB(id)

		if err != nil {
			if NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return WrapError(err)
		} else if strings.ToLower(object.LoadBalancerStatus) == strings.ToLower(string(status)) {
			//TODO
			break
		}
		if time.Now().After(deadline) {
			return WrapErrorf(err, WaitTimeoutMsg, id, GetFunc(1), timeout, object.LoadBalancerStatus, status, ProviderERROR)
		}
		time.Sleep(DefaultIntervalShort * time.Second)
	}
	return nil
}

func (s *SlbService) WaitForListener(loadBalancerId string, port int, protocol Protocol, status Status, timeout int) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	for {
		listener, err := s.DescribeLoadBalancerListenerAttribute(loadBalancerId, port, protocol)
		if err != nil && !IsExceptedErrors(err, []string{LoadBalancerNotFound}) {
			return err
		}

		if value, ok := listener["Status"]; ok && strings.ToLower(value.(string)) == strings.ToLower(string(status)) {
			//TODO
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("LoadBalancer Listener", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)

	}
	return nil
}

func (s *SlbService) slbRemoveAccessControlListEntryPerTime(list []interface{}, aclId string) error {
	req := slb.CreateRemoveAccessControlListEntryRequest()
	req.AclId = aclId
	b, _ := json.Marshal(list)
	req.AclEntrys = string(b)
	_, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.RemoveAccessControlListEntry(req)
	})
	if err != nil {
		if !IsExceptedError(err, SlbAclEntryEmpty) {
			return fmt.Errorf("RemoveAccessControlListEntry got an error: %#v", err)
		}
	}

	return nil
}

func (s *SlbService) SlbRemoveAccessControlListEntry(list []interface{}, aclId string) error {
	num := len(list)

	if num <= 0 {
		return nil
	}

	t := (num + max_num_per_time - 1) / max_num_per_time
	for i := 0; i < t; i++ {
		start := i * max_num_per_time
		end := (i + 1) * max_num_per_time

		if end > num {
			end = num
		}

		slice := list[start:end]
		if err := s.slbRemoveAccessControlListEntryPerTime(slice, aclId); err != nil {
			return err
		}
	}

	return nil
}

func (s *SlbService) slbAddAccessControlListEntryPerTime(list []interface{}, aclId string) error {
	req := slb.CreateAddAccessControlListEntryRequest()
	req.AclId = aclId
	b, _ := json.Marshal(list)
	req.AclEntrys = string(b)
	_, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.AddAccessControlListEntry(req)
	})
	if err != nil {
		return fmt.Errorf("AddAccessControlListEntry got an error: %#v", err)
	}

	return nil
}

func (s *SlbService) SlbAddAccessControlListEntry(list []interface{}, aclId string) error {
	num := len(list)

	if num <= 0 {
		return nil
	}

	t := (num + max_num_per_time - 1) / max_num_per_time
	for i := 0; i < t; i++ {
		start := i * max_num_per_time
		end := (i + 1) * max_num_per_time

		if end > num {
			end = num
		}
		slice := list[start:end]
		if err := s.slbAddAccessControlListEntryPerTime(slice, aclId); err != nil {
			return err
		}
	}

	return nil
}

// Flattens an array of slb.AclEntry into a []map[string]string
func (s *SlbService) FlattenSlbAclEntryMappings(list []slb.AclEntry) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))

	for _, i := range list {
		l := map[string]interface{}{
			"entry":   i.AclEntryIP,
			"comment": i.AclEntryComment,
		}
		result = append(result, l)
	}

	return result
}

// Flattens an array of slb.AclEntry into a []map[string]string
func (s *SlbService) flattenSlbRelatedListenerMappings(list []slb.RelatedListener) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))

	for _, i := range list {
		l := map[string]interface{}{
			"load_balancer_id": i.LoadBalancerId,
			"protocol":         i.Protocol,
			"frontend_port":    i.ListenerPort,
			"acl_type":         i.AclType,
		}
		result = append(result, l)
	}

	return result
}

func (s *SlbService) describeSlbCACertificate(caCertificateId string) (*slb.CACertificate, error) {
	request := slb.CreateDescribeCACertificatesRequest()
	request.CACertificateId = caCertificateId
	raw, error := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeCACertificates(request)
	})
	if error != nil {
		return nil, error
	}
	caCertificates, _ := raw.(*slb.DescribeCACertificatesResponse)

	if len(caCertificates.CACertificates.CACertificate) != 1 {
		msg := fmt.Sprintf("DescribeCACertificates id %s got an error %s",
			caCertificateId, SlbCACertificateIdNotFound)
		var err = GetNotFoundErrorFromString(msg)
		return nil, err
	}

	serverCertificate := caCertificates.CACertificates.CACertificate[0]
	return &serverCertificate, nil
}

func (s *SlbService) describeSlbServerCertificate(serverCertificateId string) (*slb.ServerCertificate, error) {
	request := slb.CreateDescribeServerCertificatesRequest()
	request.ServerCertificateId = serverCertificateId

	raw, error := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeServerCertificates(request)
	})
	if error != nil {
		return nil, error
	}
	serverCertificates, _ := raw.(*slb.DescribeServerCertificatesResponse)

	if len(serverCertificates.ServerCertificates.ServerCertificate) != 1 {
		msg := fmt.Sprintf("DescribeServerCertificates id %s got an error %s",
			serverCertificateId, SlbServerCertificateIdNotFound)
		err := GetNotFoundErrorFromString(msg)
		return nil, err
	}

	serverCertificate := serverCertificates.ServerCertificates.ServerCertificate[0]

	return &serverCertificate, nil
}

func (s *SlbService) readFileContent(file_name string) (string, error) {
	b, err := ioutil.ReadFile(file_name)
	if err != nil {
		return "", err
	}
	return string(b), err
}

// setTags is a helper to set the tags for a resource. It expects the
// tags field to be named "tags"
func (s *SlbService) setSlbInstanceTags(d *schema.ResourceData) error {

	if d.HasChange("tags") {
		oraw, nraw := d.GetChange("tags")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		create, remove := diffTags(tagsFromMap(o), tagsFromMap(n))

		// Set tags
		if len(remove) > 0 {
			if err := s.slbRemoveTags(remove, d.Id()); err != nil {
				return err
			}
		}

		if len(create) > 0 {
			if err := s.slbAddTags(create, d.Id()); err != nil {
				return err
			}
		}

		d.SetPartial("tags")
	}

	return nil
}

func toSlbTagsString(tags []Tag) string {
	slbTags := make([]SlbTag, 0, len(tags))

	for _, tag := range tags {
		slbTag := SlbTag{
			TagKey:   tag.Key,
			TagValue: tag.Value,
		}
		slbTags = append(slbTags, slbTag)
	}

	b, _ := json.Marshal(slbTags)

	return string(b)
}

func (s *SlbService) slbAddTagsPerTime(tags []Tag, loadBalancerId string) error {
	request := slb.CreateAddTagsRequest()
	request.LoadBalancerId = loadBalancerId
	request.Tags = toSlbTagsString(tags)

	_, error := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.AddTags(request)
	})

	if error != nil {
		return fmt.Errorf("AddTags got an error: %#v", error)
	}

	return nil
}

func (s *SlbService) slbRemoveTagsPerTime(tags []Tag, loadBalancerId string) error {
	request := slb.CreateRemoveTagsRequest()
	request.LoadBalancerId = loadBalancerId
	request.Tags = toSlbTagsString(tags)

	_, error := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.RemoveTags(request)
	})

	if error != nil {
		return fmt.Errorf("RemoveTags got an error: %#v", error)
	}

	return nil
}

func (s *SlbService) slbAddTags(tags []Tag, loadBalancderId string) error {
	num := len(tags)

	if num <= 0 {
		return nil
	}

	t := (num + tags_max_num_per_time - 1) / tags_max_num_per_time
	for i := 0; i < t; i++ {
		start := i * tags_max_num_per_time
		end := (i + 1) * tags_max_num_per_time

		if end > num {
			end = num
		}
		slice := tags[start:end]
		if err := s.slbAddTagsPerTime(slice, loadBalancderId); err != nil {
			return err
		}
	}

	return nil
}

func (s *SlbService) slbRemoveTags(tags []Tag, loadBalancderId string) error {
	num := len(tags)

	if num <= 0 {
		return nil
	}

	t := (num + tags_max_num_per_time - 1) / tags_max_num_per_time
	for i := 0; i < t; i++ {
		start := i * tags_max_num_per_time
		end := (i + 1) * tags_max_num_per_time

		if end > num {
			end = num
		}
		slice := tags[start:end]
		if err := s.slbRemoveTagsPerTime(slice, loadBalancderId); err != nil {
			return err
		}
	}

	return nil
}

func (s *SlbService) toTags(tagSet []slb.TagSet) (tags []Tag) {
	result := make([]Tag, 0, len(tagSet))
	for _, t := range tagSet {
		tag := Tag{
			Key:   t.TagKey,
			Value: t.TagValue,
		}
		result = append(result, tag)
	}

	return result
}

func (s *SlbService) describeTagsPerTime(loadBalancerId string, pageNumber, pageSize int) (tags []Tag, err error) {
	request := slb.CreateDescribeTagsRequest()
	request.LoadBalancerId = loadBalancerId
	request.PageNumber = requests.NewInteger(pageNumber)
	request.PageSize = requests.NewInteger(pageSize)

	raw, err := s.client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeTags(request)
	})

	if err != nil {
		tmp := make([]Tag, 0)
		return tmp, err
	}
	resp, _ := raw.(*slb.DescribeTagsResponse)

	return s.toTags(resp.TagSets.TagSet), nil
}

func (s *SlbService) describeTags(loadBalancerId string) (tags []Tag, err error) {
	result := make([]Tag, 0, 50)

	for i := 1; ; i++ {
		tagList, err := s.describeTagsPerTime(loadBalancerId, i, tags_max_page_size)
		if err != nil {
			return result, err
		}

		if len(tagList) == 0 {
			break
		}
		result = append(result, tagList...)
	}

	return result, nil
}

func (s *SlbService) slbTagsToMap(tags []Tag) map[string]string {
	result := make(map[string]string)
	for _, t := range tags {
		result[t.Key] = t.Value
	}

	return result
}
