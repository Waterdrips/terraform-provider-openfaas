package openfaas

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/openfaas/faas-provider/types"
	"net/http"
)

func resourceOpenFaaSSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpenFaaSSecretCreate,
		Delete: resourceOpenFaaSSecretDelete,
		Read:   resourceOpenFaaSSecretRead,
		Update: resourceOpenFaaSSecretUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Default:  "openfaas-fn",
				Required: false,
				ForceNew: true,
				Optional: true,
			},
		},
	}
}

func resourceOpenFaaSSecretUpdate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	namespace := d.Get("namespace").(string)
	value := d.Get("value").(string)

	secret := types.Secret{
		Name:      name,
		Namespace: namespace,
		Value:     value,
	}
	config := meta.(Config)
	statusCode, resp := config.Client.UpdateSecret(context.Background(), secret)

	if statusCode != http.StatusAccepted {
		return fmt.Errorf(resp)
	}
	return nil
}

func resourceOpenFaaSSecretCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	value := d.Get("value").(string)
	namespace := d.Get("namespace").(string)

	secret := types.Secret{
		Name:      name,
		Namespace: namespace,
		Value:     value,
	}

	config := meta.(Config)

	statusCode, errString := config.Client.CreateSecret(context.Background(), secret)
	if statusCode != http.StatusAccepted {
		return fmt.Errorf(errString)
	}

	d.SetId(name + namespace)
	return nil
}

func resourceOpenFaaSSecretRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	namespace := d.Get("namespace").(string)
	config := meta.(Config)

	secret, err := config.Client.GetSecretList(context.Background(), namespace)
	if err != nil {
		return err
	}
	for _, s := range secret {
		if s.Name == name {
			return flattenOpenFaaSSecretResource(d, s)
		}
	}

	d.SetId("")
	return nil
}

func resourceOpenFaaSSecretDelete(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	namespace := d.Get("namespace").(string)

	secret := types.Secret{
		Name:      name,
		Namespace: namespace,
	}

	config := meta.(Config)

	err := config.Client.RemoveSecret(context.Background(), secret)
	return err
}
