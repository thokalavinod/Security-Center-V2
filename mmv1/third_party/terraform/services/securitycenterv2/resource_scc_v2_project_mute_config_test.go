package securitycenterv2_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccSecurityCenterv2ProjectMuteConfig_basic(t *testing.T) {
	t.Parallel()

	contextBasic := map[string]interface{}{
		"project_id":    envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
		"location":      "global",
	}

	contextHighSeverity := map[string]interface{}{
		"project_id":    envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
		"location":      "us_central",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecurityCenterv2ProjectMuteConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterv2ProjectMuteConfig_basic(contextBasic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_scc_v2_project_mute_config.default", "description", "A test project mute config"),
					resource.TestCheckResourceAttr(
						"google_scc_v2_project_mute_config.default", "filter", "severity = \"LOW\""),
					resource.TestCheckResourceAttr(
						"google_scc_v2_project_mute_config.default", "project_mute_config_id", fmt.Sprintf("tf-test-my-config%s", contextBasic["random_suffix"])),
					resource.TestCheckResourceAttr(
						"google_scc_v2_project_mute_config.default", "location", contextBasic["location"].(string)),
					resource.TestCheckResourceAttr(
						"google_scc_v2_project_mute_config.default", "parent", fmt.Sprintf("projects/%s", contextBasic["project_id"])),
				),
			},
			{
				ResourceName:            "google_scc_v2_project_mute_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "project_mute_config_id"},
			},
			{
				Config: testAccSecurityCenterv2ProjectMuteConfig_highSeverity(contextHighSeverity),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_scc_v2_project_mute_config.default", "description", "A test project mute config with high severity"),
					resource.TestCheckResourceAttr(
						"google_scc_v2_project_mute_config.default", "filter", "severity = \"HIGH\""),
					resource.TestCheckResourceAttr(
						"google_scc_v2_project_mute_config.default", "project_mute_config_id", fmt.Sprintf("tf-test-my-config%s", contextHighSeverity["random_suffix"])),
					resource.TestCheckResourceAttr(
						"google_scc_v2_project_mute_config.default", "location", contextHighSeverity["location"].(string)),
					resource.TestCheckResourceAttr(
						"google_scc_v2_project_mute_config.default", "parent", fmt.Sprintf("projects/%s", contextHighSeverity["project_id"])),
				),
			},
			{
				ResourceName:            "google_scc_v2_project_mute_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parent", "project_mute_config_id"},
			},
		},
	})
}

func testAccSecurityCenterv2ProjectMuteConfig_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_v2_project_mute_config" "default" {
  description          = "A test project mute config"
  filter               = "severity = \"LOW\""
  project_mute_config_id = "tf-test-my-config%{random_suffix}"
  location             = "%{location}"
  parent               = "projects/%{project_id}"
}
`, context)
}

func testAccSecurityCenterv2ProjectMuteConfig_highSeverity(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_scc_v2_project_mute_config" "default" {
  description          = "A test project mute config with high severity"
  filter               = "severity = \"HIGH\""
  project_mute_config_id = "tf-test-my-config%{random_suffix}"
  location             = "%{location}"
  parent               = "projects/%{project_id}"
}
`, context)
}
func testAccCheckSecurityCenterv2ProjectMuteConfigDestroyProducer(t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.Provider.Meta().(*acctest.Config)
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_scc_v2_project_mute_config" {
				continue
			}
			// Initialize Security Command Center Service
			sc, err := securitycenter.NewService(context.Background(), config.GoogleClientOptions...)
			if err != nil {
				return fmt.Errorf("Error creating Security Command Center client: %s", err)
			}
			// Get the project mute config by name
			name := rs.Primary.ID
			_, err = sc.Projects.MuteConfigs.Get(name).Do()
			if err == nil {
				return fmt.Errorf("Project mute config %s still exists", name)
			}
		}

		return nil
	}
}
