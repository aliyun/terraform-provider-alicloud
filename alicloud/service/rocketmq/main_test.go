package rocketmq_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestMain lets the resource test harness register the -sweep flag so the
// package's instance sweeper is dispatched for this external test package.
func TestMain(m *testing.M) { resource.TestMain(m) }
