package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestTerraformAwsLogsCombined(t *testing.T) {
	// Note: do not run this test in t.Parallel() mode.

	testName := fmt.Sprintf("terratest-aws-logs-%s", strings.ToLower(random.UniqueId()))
	// AWS only supports one configuration recorder per region.
	// Each test using aws-config will need to specify a different region.
	awsRegion := "us-east-2"
	vpcAzs := aws.GetAvailabilityZones(t, awsRegion)[:3]
	testRedshift := !testing.Short()

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/combined/",
		Vars: map[string]interface{}{
			"region":        awsRegion,
			"vpc_azs":       vpcAzs,
			"test_name":     testName,
			"test_redshift": testRedshift,
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
