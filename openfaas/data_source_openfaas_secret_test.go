package openfaas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TODO NS
// TODO changes?
func TestAccDataSourceOpenFaaSSecret_basic(t *testing.T) {
	name := fmt.Sprintf("testaccopenfaassecret-basic-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOpenFaaSSecretConfigBasic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.openfaas_secret.secret_test", "name", name),
					resource.TestCheckResourceAttr("data.openfaas_secret.secret_test", "namespace", "openfaas-fn"),
				),
			},
		},
	})
}

func testAccDataSourceOpenFaaSSecretConfigBasic(secretName string) string {
	return fmt.Sprintf(`resource "openfaas_secret" "secret_test" {
  name      = "%s"
  value     = "Something"
}

data "openfaas_secret" "secret_test" {
  name = openfaas_secret.secret_test.name
  namespace = "openfaas-fn"
}`, secretName)
}
