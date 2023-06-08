package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestExamples_resources(t *testing.T) {
	tests := []struct {
		name   string
		before func() error
		after  func() error
		dir    string
	}{
		{
			name: "testing demo",
			dir:  "./examples",
		},
		{
			name: "testing resource image",
			before: func() error {
				_, err := os.Create("/tmp/example.img")

				return err
			},
			after: func() error {
				return os.Remove("/tmp/example.img")
			},
			dir: "./examples/resources/corellium_image",
		},
		{
			name: "testing resource instance",
			dir:  "./examples/resources/corellium_instance",
		},
		{
			name: "testing resource instance with local-exec provisioner",
			dir:  "./examples/resources/corellium_instance/local-exec",
		},
		{
			name: "testing resource instance with remote-exec provisioner",
			dir:  "./examples/resources/corellium_instance/remote-exec",
		},
		{
			name: "testing resource project",
			dir:  "./examples/resources/corellium_project",
		},
		/*{
			name: "testing resource snapshot",
			dir:  "./examples/resources/corellium_snapshot",
		},*/
		{
			name: "testing resource user",
			dir:  "./examples/resources/corellium_user",
		},
		/*{
			name: "testing resource webplayer",
			dir:  "./examples/resources/corellium_webplayer",
		},*/
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil && tt.after != nil {
				if err := tt.before(); err != nil {
					t.Fatal(err)
				}

				defer tt.after() //nolint:errcheck
			}

			terraformOptions := &terraform.Options{
				BackendConfig: map[string]interface{}{
					"bucket":                      os.Getenv("S3_BUCKET"),
					"endpoint":                    os.Getenv("S3_ENDPOINT"),
					"region":                      os.Getenv("S3_REGION"),
					"access_key":                  os.Getenv("S3_ACCESS_KEY"),
					"secret_key":                  os.Getenv("S3_SECRET_KEY"),
					"sse_customer_key":            os.Getenv("S3_SSE_KEY"),
					"key":                         fmt.Sprintf("%s/terraform.tfstate", tt.dir),
					"skip_credentials_validation": true,
					"skip_region_validation":      true,
					"encrypt":                     true,
				},
				TerraformDir: tt.dir,
			}

			defer terraform.Destroy(t, terraformOptions)
			terraform.InitAndApply(t, terraformOptions)
		})
	}
}
