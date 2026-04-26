# barracudawaf_parameter_optimizer

`barracudawaf_parameter_optimizer` manages `Parameter Optimizer` on the Barracuda Web Application Firewall.

## Example Usage

```hcl
resource "barracudawaf_parameter_optimizer" "example" {
  parent      = [barracudawaf_services.example.name]
  name        = "example-optimizer"
  start_token = "some-token"
}
```

## Argument Reference

* `parent` (Required) - The parent service name.
* `name` (Required) - Optimizer Name.
* `start_token` (Required) - Start Token.
