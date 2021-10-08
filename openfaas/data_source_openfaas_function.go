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
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "openfaas-fn",
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
			"env_vars": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"constraints": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"secrets": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"labels": {
				Type:             schema.TypeMap,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: labelsDiffFunc,
			},
			"annotations": {
				Type:             schema.TypeMap,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: annotationsDiffFunc,
			},
			"limits": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"memory": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cpu": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"requests": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"memory": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cpu": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"read_only_root_file_system": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func dataSourceOpenFaaSFunctionRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	namespace := d.Get("namespace").(string)
	config := meta.(Config)

	log.Printf("[DEBUG] Reading function Balancer: %s", name)
	function, err := config.Client.GetFunctionInfo(context.Background(), name, namespace)
	if err != nil {
		return fmt.Errorf("error retrieving function: %s", err)
	}

	d.SetId(makeID(function.Name, function.Namespace))

	return flattenOpenFaaSFunctionResource(d, function)
}
