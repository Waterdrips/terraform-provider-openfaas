package openfaas

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceOpenFaaSFunction() *schema.Resource {
	return &schema.Resource{
		Create: resourceOpenFaaSFunctionCreate,
		Read:   resourceOpenFaaSFunctionRead,
		Update: resourceOpenFaaSFunctionUpdate,
		Delete: resourceOpenFaaSFunctionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"image": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"f_process": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"env_vars": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"registry_auth": {
				Type:     schema.TypeString,
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
				DiffSuppressFunc: labelsDiffFunc,
			},
			"annotations": {
				Type:             schema.TypeMap,
				Optional:         true,
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

func resourceOpenFaaSFunctionCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	deploySpec := expandDeploymentSpec(d, name)
	config := meta.(Config)

	statusCode := config.Client.DeployFunction(context.Background(), deploySpec)
	if statusCode >= 300 {
		return fmt.Errorf("error deploying function %s status code %d", name, statusCode)
	}

	d.SetId(name)
	return nil
}

func resourceOpenFaaSFunctionRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Id()
	config := meta.(Config)

	function, err := config.Client.GetFunctionInfo(context.Background(), name, "")

	if err != nil {
		if isFunctionNotFound(err) {
			d.SetId("")
			return nil
		}

		return err
	}

	return flattenOpenFaaSFunctionResource(d, function)
}

func resourceOpenFaaSFunctionUpdate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	deploySpec := expandDeploymentSpec(d, name)
	config := meta.(Config)
	statusCode := config.Client.DeployFunction(context.Background(), deploySpec)
	if statusCode >= 300 {
		return fmt.Errorf("error deploying function %s status code %d", name, statusCode)
	}

	return nil
}

func resourceOpenFaaSFunctionDelete(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	config := meta.(Config)

	err := config.Client.DeleteFunction(context.Background(), name, "")
	return err
}

func isFunctionNotFound(err error) bool {
	return strings.Contains(err.Error(), "404") ||
		strings.Contains(err.Error(), "No such function")
}

var defaultLabels = map[string]struct{}{
	"labels.faas_function": {},
	"labels.function":      {},
	"labels.uid":           {},
}

var defaultAnnotations = map[string]struct{}{
	"annotations.prometheus.io.scrape": {},
}

const extraProviderLabelsCount = 3

func labelsDiffFunc(k, old, new string, d *schema.ResourceData) bool {
	return defaultFieldsDiffFunc(k, old, new, defaultLabels)
}

func annotationsDiffFunc(k, old, new string, d *schema.ResourceData) bool {
	return defaultFieldsDiffFunc(k, old, new, defaultAnnotations)
}

func defaultFieldsDiffFunc(k, old, new string, defaults map[string]struct{}) bool {
	if _, ok := defaults[k]; ok {
		if new != "" {
			return old == new
		}
		return true
	}

	o, err := strconv.Atoi(old)
	if err != nil {
		return old == new
	}

	n, err := strconv.Atoi(new)
	if err != nil {
		return old == new
	}

	return o == n
}
