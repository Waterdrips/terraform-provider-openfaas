package openfaas

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceOpenFaaSSecret() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOpenFaaSSecretRead,
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
			"value": {
				Type:             schema.TypeString,
				Optional:         true,
				Sensitive:        true,
				DiffSuppressFunc: secretValueDiffFunc,
			},
		},
	}
}

func secretValueDiffFunc(_, _, _ string, _ *schema.ResourceData) bool {
	return true
}
func dataSourceOpenFaaSSecretRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	namespace := d.Get("namespace").(string)
	config := meta.(Config)

	log.Printf("[DEBUG] Reading Secret: %s", name)
	secrets, err := config.Client.GetSecretList(context.Background(), namespace)
	if err != nil {
		return fmt.Errorf("error retrieving function: %s", err)
	}

	for _, s := range secrets {
		if s.Name == name {
			d.SetId(s.Name + s.Namespace)
			return flattenOpenFaaSSecretResource(d, s)
		}
	}

	return fmt.Errorf("unable to find a matching secret. name: [%s] namespace [%s]", name, namespace)
}
