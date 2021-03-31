package openfaas

import (
	"fmt"
	"testing"

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
				Config: testAccOpenFaaSFunctionConfig_basic(name),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
