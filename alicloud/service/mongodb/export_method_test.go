package mongodb

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
)

// The two methods below are reflected by name from the acctest framework via
// ResourceCheckInitWithDescribeMethod, so they must be real methods on
// *MongoDBService. They therefore live in the internal test package (package
// mongodb) rather than the external mongodb_test package, which cannot add
// methods to a type from another package. They exist only in the test binary.

// TestDescribeMongoDBPublicNetworkAddress should be visible, which can be reflected by the test framework.
func (s *MongoDBService) TestDescribeMongoDBPublicNetworkAddress(id string) (interface{}, error) {
	object, err := s.DescribeReplicaSetRole(id)
	if err != nil {
		return nil, errs.WrapError(err)
	}

	if replicaSetsMap, ok := object["ReplicaSets"].(map[string]interface{}); ok && replicaSetsMap != nil {
		if replicaSetsList, ok := replicaSetsMap["ReplicaSet"]; ok && replicaSetsList != nil {
			for _, replicaSets := range replicaSetsList.([]interface{}) {
				replicaSetsArg := replicaSets.(map[string]interface{})
				networkType, ok := replicaSetsArg["NetworkType"]
				if ok && networkType == "Public" {
					return object, nil
				}
			}
		}
	}

	// simulate NotFoundError.
	return nil, &errs.ComplexError{
		Err: fmt.Errorf("%s: public network address of %s", errs.ResourceNotfound, id),
	}
}

// TestCheckResourceMongoDBReplicaSetRoleExist should be prefixed by upper case, in other words, should be visible
// outside, otherwise it cannot be reflected by the test framework.
func (s *MongoDBService) TestCheckResourceMongoDBReplicaSetRoleExist(_ string) (interface{}, error) {
	return nil, nil
}
