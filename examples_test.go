package main

import (
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
		{
			name: "testing resource webplayer",
			dir:  "./examples/resources/corellium_webplayer",
		},
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
				TerraformDir: tt.dir,
			}

			defer terraform.Destroy(t, terraformOptions)
			terraform.InitAndApply(t, terraformOptions)
		})
	}
}
