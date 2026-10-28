package corellium

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestUnitCorelliumV1ImageResource exercises the corellium_image resource's
// Create/Read/Delete lifecycle against a mock API server. Update is a
// no-op in the resource implementation (images can't be changed in
// place), so no update step is exercised here. No TF_ACC, live token, or
// network egress beyond loopback.
func TestUnitCorelliumV1ImageResource(t *testing.T) {
	const imageID = "88888888-8888-8888-8888-888888888888"
	const projectID = "33333333-3333-3333-3333-333333333333"

	// Create opens the image file before uploading it, so it must exist.
	imagePath := filepath.Join(t.TempDir(), "example.img")
	if err := os.WriteFile(imagePath, []byte("image"), 0o600); err != nil {
		t.Fatal(err)
	}

	srv := newMockAPIServer(t)
	imageBody := map[string]any{
		"status":    "active",
		"id":        imageID,
		"name":      "test",
		"type":      "backup",
		"project":   projectID,
		"createdAt": "2026-01-01T00:00:00Z",
	}
	srv.handle("POST /api/v3/images", jsonHandler(http.StatusOK, imageBody))
	srv.handle("GET /api/v3/images/"+imageID, jsonHandler(http.StatusOK, imageBody))
	srv.handle("DELETE /api/v3/images/"+imageID, jsonHandler(http.StatusOK, map[string]any{}))

	providerConfig := fmt.Sprintf(`
		provider "corellium" {
			token = "unit-test-token"
			host  = %q
		}`, srv.hostAttr())

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
                resource "corellium_image" "test" {
                    name         = "test"
                    type         = "backup"
                    filename     = %q
                    encapsulated = false
                    project      = %q
                }`, imagePath, projectID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("corellium_image.test", "id"),
					resource.TestCheckResourceAttr("corellium_image.test", "name", "test"),
					resource.TestCheckResourceAttr("corellium_image.test", "type", "backup"),
					resource.TestCheckResourceAttr("corellium_image.test", "filename", imagePath),
					resource.TestCheckResourceAttr("corellium_image.test", "project", projectID),
					resource.TestCheckResourceAttr("corellium_image.test", "status", "active"),
				),
			},
		},
	})
}
