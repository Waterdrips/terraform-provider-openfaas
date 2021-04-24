package openfaas

import (
	"context"
	"errors"
	"fmt"
	"github.com/openfaas/faas-cli/proxy"
	"github.com/openfaas/faas-provider/types"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// TestAccResourceOpenFaaSFunction_basic requires an anonymous OpenFaaS
// deployment running on localhost:8080, with a secret foo. i.e. `faas secret create foo --from-literal baz`
func TestAccResourceOpenFaaSFunction_basic(t *testing.T) {
	var conf types.FunctionStatus
	functionName := fmt.Sprintf("testaccopenfaasfunction-basic-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "openfaas_function.function_test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckOpenFaaSFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenFaaSFunctionConfig_basic(functionName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckOpenFaaSFunctionExists("openfaas_function.function_test", &conf),
					resource.TestCheckResourceAttr("openfaas_function.function_test", "name", functionName),
					resource.TestCheckResourceAttr("openfaas_function.function_test", "image", "functions/alpine:latest"),
					resource.TestCheckResourceAttr("openfaas_function.function_test", "f_process", "sha512sum"),
					resource.TestCheckResourceAttr("openfaas_function.function_test", "labels.%", "2"),
					resource.TestCheckResourceAttr("openfaas_function.function_test", "labels.Name", "TestAccOpenFaaSFunction_basic"),
					resource.TestCheckResourceAttr("openfaas_function.function_test", "labels.Environment", "Test"),
					resource.TestCheckResourceAttr("openfaas_function.function_test", "annotations.%", "1"),
					resource.TestCheckResourceAttr("openfaas_function.function_test", "annotations.CreatedDate", "Mon Sep  3 07:15:55 BST 2018"),
					resource.TestCheckResourceAttr("openfaas_function.function_test", "requests.#", "1"),
					resource.TestCheckResourceAttr("openfaas_function.function_test", "requests.2082038905.memory", "10m"),
					resource.TestCheckResourceAttr("openfaas_function.function_test", "requests.2082038905.cpu", "100m"),
					resource.TestCheckResourceAttr("openfaas_function.function_test", "limits.#", "1"),
					resource.TestCheckResourceAttr("openfaas_function.function_test", "limits.1197768549.memory", "20m"),
					resource.TestCheckResourceAttr("openfaas_function.function_test", "limits.1197768549.cpu", "200m"),
					resource.TestCheckResourceAttr("openfaas_function.function_test", "secrets.#", "1"),
					resource.TestCheckResourceAttr("openfaas_function.function_test", "secrets.2356372769", "foo"),
				),
			},
		},
	})
}

func testAccCheckOpenFaaSFunctionDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openfaas_function" {
			continue
		}

		_, err := config.Client.GetFunctionInfo(context.Background(), rs.Primary.ID, "")

		if err == nil {
			return fmt.Errorf("function %q still exists", rs.Primary.ID)
		}

		// Verify the error
		if isFunctionNotFound(err) {
			return nil
		} else {
			return fmt.Errorf("unexpected error checking function destroyed: %s", err.Error())
		}
	}

	return nil
}

func testAccCheckOpenFaaSFunctionExists(n string, res *types.FunctionStatus) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no function ID is set")
		}

		config := testAccProvider.Meta().(Config)

		function, err := config.Client.GetFunctionInfo(context.Background(), rs.Primary.ID, "")

		if err != nil {
			return err
		}

		*res = function
		return nil
	}
}

func testAccOpenFaaSFunctionConfig_basic(functionName string) string {
	timeout := 1 * time.Second
	ofClient, _ := proxy.NewClient(newAuthChain("", "", "", "http://localhost:8080"), "http://localhot:8080", GetDefaultCLITransport(true, &timeout), &timeout)
	ofClient.CreateSecret(context.Background(), types.Secret{
		Name:      functionName,
		Namespace: "openfaas-fn",
		Value:     "foo",
	})
	return fmt.Sprintf(`resource "openfaas_function" "function_test" {
  name      = "%s"
  image     = "functions/alpine:latest"
  f_process = "sha512sum"
  labels = {
    Name        = "TestAccOpenFaaSFunction_basic"
    Environment = "Test"
  }

  annotations = {
    CreatedDate = "Mon Sep  3 07:15:55 BST 2018"
  }

  requests {
    memory = "10m"
    cpu    = "100m"
  }

  limits {
    memory = "20m"
    cpu    = "200m"
  }

  //secrets = ["%s"]
}`, functionName, functionName)
}
