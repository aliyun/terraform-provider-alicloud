package sample

import (
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// BucketReplicationSample shows how to set, get and delete bucket's replication.
func BucketReplicationSample() {
	// New client
	client, err := oss.New(endpoint, accessID, accessKey)
	if err != nil {
		HandleError(err)
	}

	// Create src, dst buckets with default parameters
	err = client.CreateBucket(bucketName)
	if err != nil {
		HandleError(err)
	}
	dstBucketName := fmt.Sprintf("%s-dst", bucketName)
	err = client.CreateBucket(dstBucketName)
	if err != nil {
		HandleError(err)
	}

	// Case 1: Set the replication. The applied objects' prefixSet includes xx/ and the action is all. HistoricalObjectReplication is enabled
	destnation := oss.Destination{
		Bucket:       dstBucketName,
		Location:     "oss-cn-beijing",
		TransferType: "internal",
	}
	rule1 := oss.ReplicationRule{
		PrefixSet:                   &oss.PrefixSet{Prefix: []string{"xx/"}},
		Action:                      "ALL",
		Destination:                 &destnation,
		HistoricalObjectReplication: "enabled",
	}

	err = client.SetBucketReplication(bucketName, rule1)
	if err != nil {
		HandleError(err)
	}

	// Case 2: Get the bucket's replication
	rep, err := client.GetBucketReplication(bucketName)
	if err != nil {
		HandleError(err)
	}
	fmt.Printf("Bucket Replication Rule: %v", rep.Rule)

	// Case 3: Get the bucket's replication
	repLocation, err := client.GetBucketReplicationLocation(bucketName)
	if err != nil {
		HandleError(err)
	}
	fmt.Printf("Bucket Replication Location: %v", repLocation)

	// Case 4: Get the bucket's replication
	repProgress, err := client.GetBucketReplicationProgress(bucketName, oss.RuleId("e8972aa9-c10f-47e3-be04-b*****17d678"))
	if err != nil {
		HandleError(err)
	}
	fmt.Printf("Bucket Replication Progress: %v", repProgress.Rule)

	// Case 5: Set the replication, The applied objects' prefixSet includes test/ and the action is put. HistoricalObjectReplication is disabled
	destnation = oss.Destination{
		Bucket:       dstBucketName,
		Location:     "oss-cn-beijing",
		TransferType: "oss_acc",
	}
	rule2 := oss.ReplicationRule{
		PrefixSet:                   &oss.PrefixSet{Prefix: []string{"test/"}},
		Action:                      "PUT",
		Destination:                 &destnation,
		HistoricalObjectReplication: "disabled",
	}

	err = client.SetBucketReplication(bucketName, rule2)
	if err != nil {
		HandleError(err)
	}

	// Case 6: Set the replication, The rule ID is rule3 and the syncRole is User
	rule3 := oss.ReplicationRule{
		PrefixSet:                   &oss.PrefixSet{Prefix: []string{"sss/"}},
		Action:                      "ALL",
		Destination:                 &destnation,
		HistoricalObjectReplication: "enabled",
		SyncRole:                    "User",
	}

	err = client.SetBucketReplication(bucketName, rule3)
	if err != nil {
		HandleError(err)
	}

	// Case 7: Set the lifecycle. The rule ID is rule4 and the
	SourceSelectionCriteria := oss.SourceSelectionCriteria{SseKmsEncryptedObjects: &oss.SseKmsEncryptedObjects{Status: "Enabled"}}
	EncryptionConfiguration := oss.EncryptionConfiguration{ReplicaKmsKeyID: "1a8d780d-0d34-49e5-ab45-7b9b*******4"} // todo kms-id 需要自行创建
	rule4 := oss.ReplicationRule{
		PrefixSet:                   &oss.PrefixSet{Prefix: []string{"ttt/"}},
		Action:                      "ALL",
		Destination:                 &destnation,
		HistoricalObjectReplication: "enabled",
		SourceSelectionCriteria:     &SourceSelectionCriteria,
		EncryptionConfiguration:     &EncryptionConfiguration,
	}

	err = client.SetBucketReplication(bucketName, rule4)
	if err != nil {
		HandleError(err)
	}

	// Case 8: Delete bucket's replication
	err = client.DeleteBucketReplication(bucketName, rep.Rule.ID)
	if err != nil {
		HandleError(err)
	}

	// Delete bucket
	err = client.DeleteBucket(bucketName)
	if err != nil {
		HandleError(err)
	}
	err = client.DeleteBucket(fmt.Sprintf("%s-dst", bucketName))
	if err != nil {
		HandleError(err)
	}

	fmt.Println("BucketReplicationSample completed")
}
