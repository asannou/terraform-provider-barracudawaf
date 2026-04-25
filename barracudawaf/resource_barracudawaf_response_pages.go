package barracudawaf

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCudaWAFResponsePages() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFResponsePagesCreate,
		Read:   resourceCudaWAFResponsePagesRead,
		Update: resourceCudaWAFResponsePagesUpdate,
		Delete: resourceCudaWAFResponsePagesDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Response Page Name",
			},
			"status_code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Status Code",
			},
			"body": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Body",
			},
			"headers": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Headers",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Session Token Type",
			},
		},

		Description: "`barracudawaf_response_pages` manages `Response Page` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFResponsePagesCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	resourceEndpoint := "/response-pages"
	err := client.CreateBarracudaWAFResource(name, hydrateBarracudaWAFResponsePagesResource(d, "post", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFResponsePagesRead(d, m)
}

func resourceCudaWAFResponsePagesRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	resourceEndpoint := "/response-pages"
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

	for _, dataItems = range resources.Data {
		if dataItems["name"] == name {
			break
		}
	}

	if dataItems["name"] != name {
		return fmt.Errorf("Barracuda WAF resource (%s) not found on the system", name)
	}

	d.Set("name", name)
	d.Set("status_code", dataItems["status-code"])
	d.Set("body", dataItems["body"])
	d.Set("headers", dataItems["headers"])
	d.Set("type", dataItems["type"])

	return nil
}

func resourceCudaWAFResponsePagesUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	resourceEndpoint := "/response-pages"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFResponsePagesResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFResponsePagesRead(d, m)
}

func resourceCudaWAFResponsePagesDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	resourceEndpoint := "/response-pages"
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

func hydrateBarracudaWAFResponsePagesResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"name":        d.Get("name").(string),
		"status-code": d.Get("status_code").(string),
		"body":        d.Get("body").(string),
		"headers":     d.Get("headers").(string),
		"type":        d.Get("type").(string),
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
		if reflect.ValueOf(val).Len() == 0 {
			delete(resourcePayload, key)
		}
	}

	return &APIRequest{
		URL:  endpoint,
		Body: resourcePayload,
	}
}
