package openfaas

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceOpenFaaSFunction() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOpenFaaSFunctionRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"image": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"f_process": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
			"annotations": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceOpenFaaSFunctionRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	config := meta.(Config)

	log.Printf("[DEBUG] Reading function Balancer: %s", name)
	function, err := config.Client.GetFunctionInfo(context.Background(), name, "")
	if err != nil {
		return fmt.Errorf("error retrieving function: %s", err)
	}

	d.SetId(function.Name)

	return flattenOpenFaaSFunctionResource(d, function)
}
