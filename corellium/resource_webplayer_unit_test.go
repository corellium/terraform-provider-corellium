package corellium

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestUnitCorelliumV1WebPlayerResource exercises the corellium_webplayer
// resource's Create/Read/Delete lifecycle against a mock API server.
// Update is a no-op (there's no webplayer update endpoint), so no update
// step is exercised. Read bypasses the generated client and builds its
// own request from the provider's configured host (GetWebPlayerManual in
// resource_webplayer.go), so this test also guards against that
// hand-rolled URL construction breaking again. No TF_ACC, live token, or
// network egress beyond loopback.
func TestUnitCorelliumV1WebPlayerResource(t *testing.T) {
	const sessionID = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	const instanceID = "55555555-5555-5555-5555-555555555555"
	const projectID = "33333333-3333-3333-3333-333333333333"

	srv := newMockAPIServer(t)
	srv.handle("GET /api/v3/instances/"+instanceID, jsonHandler(http.StatusOK, map[string]any{
		"id":           instanceID,
		"name":         "test",
		"key":          "abc123",
		"flavor":       "iphone7plus",
		"type":         "ios",
		"project":      projectID,
		"state":        "on",
		"stateChanged": "2026-01-01T00:00:00Z",
		"startedAt":    "2026-01-01T00:00:00Z",
		"created":      "2026-01-01T00:00:00Z",
		"osVersion":    "15.7.5",
		"bootOptions":  map[string]any{},
	}))
	srv.handle("POST /api/v3/webplayer", jsonHandler(http.StatusOK, map[string]any{
		"identifier": sessionID,
		"token":      "webplayer-token",
		"expiration": time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
	}))
	// Read bypasses the generated client (GetWebPlayerManual), hitting
	// GET {host}/api/v1/webplayer/{sessionId} directly.
	srv.handle("GET /api/v1/webplayer/"+sessionID, jsonHandler(http.StatusOK, []map[string]any{
		{
			"instanceid": instanceID,
			"identifier": sessionID,
			"project":    projectID,
			"token":      "webplayer-token",
			"clientid":   "",
			"features":   map[string]any{},
		},
	}))
	srv.handle("DELETE /api/v3/webplayer/"+sessionID, jsonHandler(http.StatusOK, []string{sessionID}))

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
                resource "corellium_webplayer" "test" {
                    project          = %q
                    instanceid       = %q
                    expiresinseconds = 3600
                    features = {
                        apps           = false
                        console        = false
                        coretrace      = false
                        devicecontrol  = false
                        devicedelete   = false
                        files          = false
                        frida          = false
                        images         = false
                        messaging      = false
                        netmon         = false
                        network        = false
                        portforwarding = false
                        profile        = false
                        sensors        = false
                        settings       = false
                        snapshots      = false
                        strace         = false
                        system         = false
                        connect        = false
                    }
                }`, projectID, instanceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("corellium_webplayer.test", "id"),
					resource.TestCheckResourceAttr("corellium_webplayer.test", "identifier", sessionID),
					resource.TestCheckResourceAttr("corellium_webplayer.test", "project", projectID),
					resource.TestCheckResourceAttr("corellium_webplayer.test", "instanceid", instanceID),
					resource.TestCheckResourceAttrSet("corellium_webplayer.test", "token"),
				),
			},
		},
	})
}
