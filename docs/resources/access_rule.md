# barracudawaf_access_rule

`barracudawaf_access_rule` manages `Access Rule` on the Barracuda Web Application Firewall.

## Example Usage

```hcl
resource "barracudawaf_access_rule" "example" {
  parent          = [barracudawaf_services.example.name]
  name            = "example-access-rule"
  attribute_names = ["src_ip"]
  attribute_values = ["10.0.0.1"]
}
```

## Argument Reference

* `parent` (Required) - The parent service name.
* `name` (Required) - Access Rule Name.
* `attribute_names` (Optional) - Access Rule Attribute Names.
* `attribute_values` (Optional) - Access Rule Attribute Values.
