package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestTerraformAwsLogsCloudtrail(t *testing.T) {
	// Note: do not run this test in t.Parallel() mode.
	// Running this test in parallel with other tests in the module
	// often causes issues when attempting to empty and delete the bucket.

	testName := fmt.Sprintf("terratest-aws-logs-%s", strings.ToLower(random.UniqueId()))
	awsRegion := "us-west-2"

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/cloudtrail/",
		Vars: map[string]interface{}{
			"region":    awsRegion,
			"test_name": testName,
		},
		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": awsRegion,
		},
	}

	defer terraform.Destroy(t, terraformOptions)
	// Empty and delete logs_bucket before terraform destroy
	defer aws.DeleteS3Bucket(t, awsRegion, testName)
	defer aws.EmptyS3Bucket(t, awsRegion, testName)
	terraform.InitAndApply(t, terraformOptions)
}
