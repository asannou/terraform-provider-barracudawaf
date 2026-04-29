package barracudawaf

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCudaWAFAttackTypePattern() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFAttackTypePatternCreate,
		Read:   resourceCudaWAFAttackTypePatternRead,
		Update: resourceCudaWAFAttackTypePatternUpdate,
		Delete: resourceCudaWAFAttackTypePatternDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), "/")
				if len(parts) != 2 {
					return nil, fmt.Errorf("invalid ID specified. Supposed to be <attack_type_name>/<pattern_name>")
				}
				d.Set("parent", []string{parts[0]})
				d.SetId(parts[1])
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Pattern Name",
			},
			"regex": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Pattern Regex",
			},
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v == "" {
						return
					}
					validModes := []string{"Passive", "Active", "Off"}
					found := false
					for _, mode := range validModes {
						if v == mode {
							found = true
							break
						}
					}
					if !found {
						errs = append(errs, fmt.Errorf("%q must be one of %v", key, validModes))
					}
					return
				},
				Description: "Operating Mode",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Pattern Description",
			},
			"algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v == "" {
						return
					}
					validAlgorithms := []string{
						"None",
						"Credit Card Check Digit",
						"Korean Resident Registration Number Check Digit",
						"Poland National Identification Number Check Digit",
						"Czech National Identity Card Number Check Digit",
						"Hungarian Identity Card Number Check Digit",
					}
					found := false
					for _, algorithm := range validAlgorithms {
						if v == algorithm {
							found = true
							break
						}
					}
					if !found {
						errs = append(errs, fmt.Errorf("%q must be one of %v", key, validAlgorithms))
					}
					return
				},
				Description: "Pattern Algorithm",
			},
			"case_sensitive": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "Yes" && v != "No" && v != "" {
						errs = append(errs, fmt.Errorf("%q must be either Yes, No or empty", key))
					}
					return
				},
				Description: "Case Sensitivity",
			},
			"parent": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				ForceNew: true,
			},
		},

		Description: "`barracudawaf_attack_type_pattern` manages `Attack Type Pattern` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFAttackTypePatternCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	attackTypeName := parent[0].(string)

	resourceEndpoint := "/attack-types/" + attackTypeName + "/attack-patterns"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFAttackTypePatternResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFAttackTypePatternRead(d, m)
}

func resourceCudaWAFAttackTypePatternRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	attackTypeName := parent[0].(string)

	resourceEndpoint := "/attack-types/" + attackTypeName + "/attack-patterns"
	request := &APIRequest{
		Method: "get",
		URL:    resourceEndpoint,
	}

	var dataItems map[string]interface{}
	resources, err := client.GetBarracudaWAFResource(name, request)

	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	if resources.Data == nil {
		log.Printf("[WARN] Barracuda WAF resource (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	found := false
	for _, dataItems = range resources.Data {
		if dataItems["name"] == name {
			found = true
			break
		}
	}

	if !found {
		log.Printf("[WARN] Barracuda WAF resource (%s) not found on the system", name)
		d.SetId("")
		return nil
	}

	d.Set("name", name)
	d.Set("regex", dataItems["regex"])
	d.Set("mode", dataItems["mode"])
	d.Set("description", dataItems["description"])
	d.Set("algorithm", dataItems["algorithm"])
	d.Set("case_sensitive", dataItems["case-sensitive"])

	return nil
}

func resourceCudaWAFAttackTypePatternUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	attackTypeName := parent[0].(string)

	resourceEndpoint := "/attack-types/" + attackTypeName + "/attack-patterns"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFAttackTypePatternResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFAttackTypePatternRead(d, m)
}

func resourceCudaWAFAttackTypePatternDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	attackTypeName := parent[0].(string)

	resourceEndpoint := "/attack-types/" + attackTypeName + "/attack-patterns"
	request := &APIRequest{
		Method: "delete",
		URL:    resourceEndpoint,
	}

	err := client.DeleteBarracudaWAFResource(name, request)

	if err != nil {
		return fmt.Errorf("Unable to delete the Barracuda WAF resource (%s) (%v)", name, err)
	}

	return nil
}

func hydrateBarracudaWAFAttackTypePatternResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]interface{}{
		"name":           d.Get("name").(string),
		"regex":          d.Get("regex").(string),
		"mode":           d.Get("mode").(string),
		"description":    d.Get("description").(string),
		"algorithm":      d.Get("algorithm").(string),
		"case-sensitive": d.Get("case_sensitive").(string),
	}

	// parameters not supported for updates
	if method == "put" {
		updatePayloadExceptions := [...]string{"name"}
		for _, param := range updatePayloadExceptions {
			delete(resourcePayload, param)
		}
	}

	// remove empty parameters from resource payload
	for key, val := range resourcePayload {
		if v, ok := val.(string); ok && len(v) == 0 {
			delete(resourcePayload, key)
		}
	}

	return &APIRequest{
		URL:  endpoint,
		Body: resourcePayload,
	}
}
