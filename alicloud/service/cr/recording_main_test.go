package cr_test

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/recording"
	"os"
	"testing"
)

// Recording is initialized once for every resource test in this product package.
func TestMain(m *testing.M) { os.Exit(recording.RunMain(m)) }
