package openfaas

import (
	"context"
	"errors"
	"fmt"
	"github.com/openfaas/faas-provider/types"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Test NS

// Test requires an anonymous OpenFaaS
// deployment running on localhost:8080, with a secret foo. i.e. `faas secret create foo --from-literal baz`
func TestAccResourceOpenFaaSSecret_basic(t *testing.T) {
	var conf types.Secret
	secretName := fmt.Sprintf("testaccopenfaassecret-basic-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "openfaas_secret.secret_test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOpenFaaSSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceOpenFaaSSecretConfigBasic(secretName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckOpenFaaSSecretExists("openfaas_secret.secret_test", &conf),
					resource.TestCheckResourceAttr("openfaas_secret.secret_test", "name", secretName),
					resource.TestCheckResourceAttr("openfaas_secret.secret_test", "namespace", "openfaas-fn"),
				),
			},
		},
	})
}

func testAccCheckOpenFaaSSecretDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openfaas_secret" {
			continue
		}

		name, namespace, _ := decodeID(rs.Primary.ID)

		// eugh... eventual consistency
		var err error
		var secrets []types.Secret
		for i := 0; i < 15; i++ {
			secrets, err = config.Client.GetSecretList(context.Background(), namespace)
			if err != nil {
				break
			}
			time.Sleep(time.Second * 1)
		}

		if err != nil {
			return fmt.Errorf("error getting secrets list %v", err)
		}

		if _, ok := findSecret(secrets, name); !ok {
			return nil
		} else {
			return fmt.Errorf("Secret exists [%s] in namespace [%s] ", name, namespace)
		}
	}

	return nil
}

func findSecret(secrets []types.Secret, name string) (*types.Secret, bool) {
	for _, secret := range secrets {
		if secret.Name == name {
			return &secret, true
		}
	}
	return nil, false
}

func testAccCheckOpenFaaSSecretExists(n string, res *types.Secret) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no secret id set")
		}

		config := testAccProvider.Meta().(Config)
		name, namespace, _ := decodeID(rs.Primary.ID)
		function, err := config.Client.GetSecretList(context.Background(), namespace)

		if err != nil {
			return err
		}

		if _, ok := findSecret(function, name); !ok {
			return err
		}

		return nil
	}
}

func testAccResourceOpenFaaSSecretConfigBasic(secretName string) string {
	return fmt.Sprintf(`resource "openfaas_secret" "secret_test" {
  name      = "%s"
  value     = "Something"
}
`, secretName)
}
