package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type PaiWorkspaceServiceV2 struct {
	client *connectivity.AliyunClient
}

// DescribePaiWorkspaceWorkspace <<< Encapsulated get interface for PaiWorkspace Workspace.

func (s *PaiWorkspaceServiceV2) DescribePaiWorkspaceWorkspace(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	WorkspaceId := id
	action := fmt.Sprintf("/api/v1/workspaces/%s", WorkspaceId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	request["WorkspaceId"] = id

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("AIWorkSpace", "2021-02-04", action, query, nil, nil)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"100400008", "100700008", "100400027"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Workspace", id)), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (s *PaiWorkspaceServiceV2) PaiWorkspaceWorkspaceStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePaiWorkspaceWorkspace(id)
		if err != nil {
			if NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribePaiWorkspaceWorkspace >>> Encapsulated.
// DescribePaiWorkspaceDataset <<< Encapsulated get interface for PaiWorkspace Dataset.

func (s *PaiWorkspaceServiceV2) DescribePaiWorkspaceDataset(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	DatasetId := id
	action := fmt.Sprintf("/api/v1/datasets/%s", DatasetId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	request["DatasetId"] = id

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("AIWorkSpace", "2021-02-04", action, query, nil, nil)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"201300003"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Dataset", id)), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (s *PaiWorkspaceServiceV2) PaiWorkspaceDatasetStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePaiWorkspaceDataset(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribePaiWorkspaceDataset >>> Encapsulated.
// DescribePaiWorkspaceExperiment <<< Encapsulated get interface for PaiWorkspace Experiment.

func (s *PaiWorkspaceServiceV2) DescribePaiWorkspaceExperiment(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	ExperimentId := id
	action := fmt.Sprintf("/api/v1/experiments/%s", ExperimentId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	request["ExperimentId"] = id

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("AIWorkSpace", "2021-02-04", action, query, nil, nil)

		if err != nil {
			if IsExpectedErrors(err, []string{"NotFoundErrorProblem"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"InternalServerErrorProblem"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Experiment", id)), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (s *PaiWorkspaceServiceV2) PaiWorkspaceExperimentStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePaiWorkspaceExperiment(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribePaiWorkspaceExperiment >>> Encapsulated.
// DescribePaiWorkspaceDatasetversion <<< Encapsulated get interface for PaiWorkspace Datasetversion.

func (s *PaiWorkspaceServiceV2) DescribePaiWorkspaceDatasetversion(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		err = WrapError(fmt.Errorf("invalid Resource Id %s. Expected parts' length %d, got %d", id, 2, len(parts)))
	}
	DatasetId := parts[0]
	VersionName := parts[1]
	action := fmt.Sprintf("/api/v1/datasets/%s/versions/%s", DatasetId, VersionName)
	request = make(map[string]interface{})
	query = make(map[string]*string)

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("AIWorkSpace", "2021-02-04", action, query, nil, nil)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"201300003"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("Datasetversion", id)), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (s *PaiWorkspaceServiceV2) PaiWorkspaceDatasetversionStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePaiWorkspaceDatasetversion(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribePaiWorkspaceDatasetversion >>> Encapsulated.
// DescribePaiWorkspaceRun <<< Encapsulated get interface for PaiWorkspace Run.

func (s *PaiWorkspaceServiceV2) DescribePaiWorkspaceRun(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	RunId := id
	action := fmt.Sprintf("/api/v1/runs/%s", RunId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	request["RunId"] = id

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("AIWorkSpace", "2021-02-04", action, query, nil, nil)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	currentStatus := response["Name"]
	if currentStatus == nil {
		return object, WrapErrorf(Error(GetNotFoundMessage("Run", id)), NotFoundMsg, response)
	}

	return response, nil
}

func (s *PaiWorkspaceServiceV2) PaiWorkspaceRunStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePaiWorkspaceRun(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribePaiWorkspaceRun >>> Encapsulated.

// DescribePaiWorkspaceCodeSource <<< Encapsulated get interface for PaiWorkspace CodeSource.

func (s *PaiWorkspaceServiceV2) DescribePaiWorkspaceCodeSource(id string) (object map[string]interface{}, err error) {
	client := s.client
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	CodeSourceId := id
	action := fmt.Sprintf("/api/v1/codesources/%s", CodeSourceId)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	request["CodeSourceId"] = id

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		response, err = client.RoaGet("AIWorkSpace", "2021-02-04", action, query, nil, nil)

		if err != nil {
			if IsExpectedErrors(err, []string{"201400004"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"201400002"}) {
			return object, WrapErrorf(Error(GetNotFoundMessage("CodeSource", id)), NotFoundMsg, response)
		}
		return object, WrapErrorf(err, DefaultErrorMsg, id, action, AlibabaCloudSdkGoERROR)
	}
	return response, nil
}

func (s *PaiWorkspaceServiceV2) PaiWorkspaceCodeSourceStateRefreshFunc(id string, field string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribePaiWorkspaceCodeSource(id)
		if err != nil {
			if NotFoundError(err) {
				return object, "", nil
			}
			return nil, "", WrapError(err)
		}

		v, err := jsonpath.Get(field, object)
		currentStatus := fmt.Sprint(v)

		if strings.HasPrefix(field, "#") {
			v, _ := jsonpath.Get(strings.TrimPrefix(field, "#"), object)
			if v != nil {
				currentStatus = "#CHECKSET"
			}
		}

		for _, failState := range failStates {
			if currentStatus == failState {
				return object, currentStatus, WrapError(Error(FailedToReachTargetStatus, currentStatus))
			}
		}
		return object, currentStatus, nil
	}
}

// DescribePaiWorkspaceCodeSource >>> Encapsulated.
