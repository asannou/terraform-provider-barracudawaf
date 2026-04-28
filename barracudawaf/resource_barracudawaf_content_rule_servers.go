package barracudawaf

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	subResourceContentRuleServersParams = map[string][]string{
		"ssl_policy": {
			"client_certificate",
			"enable_ssl_compatibility_mode",
			"validate_certificate",
			"enable_https",
			"enable_sni",
			"enable_ssl_3",
			"enable_tls_1",
			"enable_tls_1_1",
			"enable_tls_1_2",
			"enable_tls_1_3",
		},
		"connection_pooling": {
			"keepalive_timeout",
			"enable_connection_pooling",
		},
	}
)

func resourceCudaWAFContentRuleServers() *schema.Resource {
	return &schema.Resource{
		Create: resourceCudaWAFContentRuleServersCreate,
		Read:   resourceCudaWAFContentRuleServersRead,
		Update: resourceCudaWAFContentRuleServersUpdate,
		Delete: resourceCudaWAFContentRuleServersDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), "/")
				if len(parts) != 3 {
					return nil, fmt.Errorf("invalid ID specified. Supposed to be <service_name>/<rule_name>/<server_name>")
				}
				d.Set("parent", []string{parts[0], parts[1]})
				d.SetId(parts[2])
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"comments":        {Type: schema.TypeString, Optional: true, Description: "Comments"},
			"name":            {Type: schema.TypeString, Optional: true, Description: "Web Server Name"},
			"hostname":        {Type: schema.TypeString, Optional: true, Description: "Hostname"},
			"identifier":      {Type: schema.TypeString, Optional: true, Description: "Identifier:"},
			"ip_address":      {Type: schema.TypeString, Optional: true, Description: "IP Address"},
			"address_version": {Type: schema.TypeString, Optional: true, Description: "Version"},
			"port":            {Type: schema.TypeString, Optional: true, Description: "Port"},
			"status":          {Type: schema.TypeString, Optional: true, Description: "Status"},
			"ssl_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_certificate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Client Certificate",
						},
						"enable_ssl_compatibility_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enable SSL Compatibility Mode",
						},
						"validate_certificate": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Validate Server Certificate",
						},
						"enable_https": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Server uses SSL",
						},
						"enable_sni": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enable SNI",
						},
						"enable_ssl_3": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SSL 3.0 (Insecure)",
						},
						"enable_tls_1": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "TLS 1.0 (Insecure)",
						},
						"enable_tls_1_1": {Type: schema.TypeString, Optional: true, Description: "TLS 1.1"},
						"enable_tls_1_2": {Type: schema.TypeString, Optional: true, Description: "TLS 1.2"},
						"enable_tls_1_3": {Type: schema.TypeString, Optional: true, Description: "TLS 1.3"},
					},
				},
			},
			"connection_pooling": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"keepalive_timeout": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Keepalive Timeout",
						},
						"enable_connection_pooling": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enable Connection Pooling",
						},
					},
				},
			},
			"parent": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},

		Description: "`barracudawaf_content_rule_servers` manages `Content Rule Servers` on the Barracuda Web Application Firewall.",
	}
}

func resourceCudaWAFContentRuleServersCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)
	contentRuleName := parent[1].(string)

	resourceEndpoint := "/services/" + serviceName + "/content-rules/" + contentRuleName + "/content-rule-servers"
	err := client.CreateBarracudaWAFResource(
		name,
		hydrateBarracudaWAFContentRuleServersResource(d, "post", resourceEndpoint),
	)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF resource (%s) (%v) ", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFContentRuleServersSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to create Barracuda WAF sub resource (%s) (%v) ", name, err)
		return err
	}

	d.SetId(name)
	return resourceCudaWAFContentRuleServersRead(d, m)
}

func resourceCudaWAFContentRuleServersRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()
	log.Println("[INFO] Fetching Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)
	contentRuleName := parent[1].(string)

	resourceEndpoint := "/services/" + serviceName + "/content-rules/" + contentRuleName + "/content-rule-servers"
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
		return fmt.Errorf("Barracuda WAF resource (%s) not found on the system", name)
	}

	d.Set("name", name)

	payload := map[string]string{
		"comments":        "comments",
		"hostname":        "hostname",
		"identifier":      "identifier",
		"ip_address":      "ip-address",
		"address_version": "address-version",
		"port":            "port",
		"status":          "status",
	}

	for tfKey, apiKey := range payload {
		if val, ok := dataItems[apiKey]; ok && val != nil {
			if reflect.TypeOf(val).Kind() == reflect.Slice {
				d.Set(tfKey, sortFileList(val.([]interface{}), ""))
			} else {
				d.Set(tfKey, fmt.Sprintf("%v", val))
			}
		} else {
			d.Set(tfKey, nil)
		}
	}

	// Read sub-resources
	for subResource, subResourceParams := range subResourceContentRuleServersParams {
		subResourceEndpoint := fmt.Sprintf("%s/%s/%s", resourceEndpoint, name, strings.Replace(subResource, "_", "-", -1))
		subRequest := &APIRequest{
			Method: "get",
			URL:    subResourceEndpoint,
		}

		subResources, err := client.GetBarracudaWAFResource(name, subRequest)
		if err != nil {
			log.Printf("[ERROR] Unable to Retrieve Barracuda WAF sub-resource (%s) (%v) ", subResource, err)
			continue
		}

		if subResources.Data == nil {
			continue
		}

		var subResourceList []interface{}
		for _, subDataItems := range subResources.Data {
			subMap := make(map[string]interface{})
			for _, param := range subResourceParams {
				apiParam := strings.Replace(param, "_", "-", -1)
				if val, ok := subDataItems[apiParam]; ok && val != nil {
					if reflect.TypeOf(val).Kind() == reflect.Slice {
						subMap[param] = sortFileList(val.([]interface{}), "")
					} else {
						subMap[param] = fmt.Sprintf("%v", val)
					}
				}
			}
			subResourceList = append(subResourceList, subMap)
		}

		sortKey := ""
		for _, param := range subResourceParams {
			if param == "name" {
				sortKey = "name"
				break
			}
		}

		if len(subResourceList) > 0 {
			d.Set(subResource, sortFileList(subResourceList, sortKey))
		} else {
			d.Set(subResource, nil)
		}
	}

	return nil
}

func resourceCudaWAFContentRuleServersUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Updating Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)
	contentRuleName := parent[1].(string)

	resourceEndpoint := "/services/" + serviceName + "/content-rules/" + contentRuleName + "/content-rule-servers"
	err := client.UpdateBarracudaWAFResource(name, hydrateBarracudaWAFContentRuleServersResource(d, "put", resourceEndpoint))

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF resource (%s) (%v)", name, err)
		return err
	}

	err = client.hydrateBarracudaWAFContentRuleServersSubResource(d, name, resourceEndpoint)

	if err != nil {
		log.Printf("[ERROR] Unable to update the Barracuda WAF sub resource (%s) (%v)", name, err)
		return err
	}

	return resourceCudaWAFContentRuleServersRead(d, m)
}

func resourceCudaWAFContentRuleServersDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*BarracudaWAF)

	name := d.Id()

	log.Println("[INFO] Deleting Barracuda WAF resource " + name)

	parent := d.Get("parent").([]interface{})
	serviceName := parent[0].(string)
	contentRuleName := parent[1].(string)

	resourceEndpoint := "/services/" + serviceName + "/content-rules/" + contentRuleName + "/content-rule-servers"
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

func hydrateBarracudaWAFContentRuleServersResource(d *schema.ResourceData, method string, endpoint string) *APIRequest {

	//resourcePayload : payload for the resource
	resourcePayload := map[string]string{
		"comments":        d.Get("comments").(string),
		"name":            d.Get("name").(string),
		"hostname":        d.Get("hostname").(string),
		"identifier":      d.Get("identifier").(string),
		"ip-address":      d.Get("ip_address").(string),
		"address-version": d.Get("address_version").(string),
		"port":            d.Get("port").(string),
		"status":          d.Get("status").(string),
	}

	// parameters not supported for updates
	if method == "put" {
		updatePayloadExceptions := [...]string{"address-version"}
		for _, param := range updatePayloadExceptions {
			delete(resourcePayload, param)
		}
	}

	// remove empty parameters from resource payload
	for key, val := range resourcePayload {
		if len(val) == 0 {
			delete(resourcePayload, key)
		}
	}

	return &APIRequest{
		URL:  endpoint,
		Body: resourcePayload,
	}
}

func (b *BarracudaWAF) hydrateBarracudaWAFContentRuleServersSubResource(
	d *schema.ResourceData,
	name string,
	endpoint string,
) error {

	for subResource, subResourceParams := range subResourceContentRuleServersParams {
		subResourceParamsLength := d.Get(subResource + ".#").(int)

		log.Printf("[INFO] Updating Barracuda WAF sub resource (%s) (%s)", name, subResource)

		for i := 0; i < subResourceParamsLength; i++ {
			subResourcePayload := map[string]string{}
			suffix := fmt.Sprintf(".%d", i)

			for _, param := range subResourceParams {
				paramSuffix := fmt.Sprintf(".%s", param)
				paramVaule := d.Get(subResource + suffix + paramSuffix).(string)

				if len(paramVaule) > 0 {
					param = strings.Replace(param, "_", "-", -1)
					subResourcePayload[param] = paramVaule
				}
			}

			err := b.UpdateBarracudaWAFSubResource(name, endpoint, &APIRequest{
				URL:  strings.Replace(subResource, "_", "-", -1),
				Body: subResourcePayload,
			})

			if err != nil {
				return err
			}
		}
	}

	return nil
}
