package openfaas

import (
	"context"
	"fmt"
	"github.com/openfaas/faas-cli/proxy"
	"github.com/openfaas/faas-provider/types"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOpenFaaSFunction_importBasic(t *testing.T) {
	resourceName := "openfaas_function.function_test"
	name := fmt.Sprintf("testaccopenfaasfunction-basic-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOpenFaaSFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenFaaSFunctionImport_basic(name),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccOpenFaaSFunctionImport_basic(functionName string) string {
	timeout := 1 * time.Second
	ofClient, _ := proxy.NewClient(newAuthChain("", "", "", "http://127.0.0.1:8080"), "http://127.0.0.1:8080", GetDefaultCLITransport(true, &timeout), &timeout)
	ofClient.CreateSecret(context.Background(), types.Secret{
		Name:      "test-import",
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

  secrets = ["test-import"]
}`, functionName)
}
