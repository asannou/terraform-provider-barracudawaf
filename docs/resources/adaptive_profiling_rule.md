# barracudawaf_adaptive_profiling_rule

`barracudawaf_adaptive_profiling_rule` manages `Adaptive Profiling Rule` on the Barracuda Web Application Firewall.

## Example Usage

```hcl
resource "barracudawaf_adaptive_profiling_rule" "example" {
  parent = [barracudawaf_services.example.name]
  name   = "example-profiling"
  host   = "*"
  url    = "/*"
  status = "On"
}
```

## Argument Reference

* `parent` (Required) - The parent service name.
* `name` (Required) - Learn Rule Name.
* `host` (Required) - Host Match.
* `url` (Required) - URL Match.
* `status` (Optional) - Status.
* `learn_from_request` (Optional) - Learn From Request.
* `learn_from_response` (Optional) - Learn From Response.
