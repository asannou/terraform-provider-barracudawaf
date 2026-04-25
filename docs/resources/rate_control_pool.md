# barracudawaf_rate_control_pool

`barracudawaf_rate_control_pool` manages `Rate Control Pool` on the Barracuda Web Application Firewall.

## Example Usage

```hcl
resource "barracudawaf_rate_control_pool" "example" {
  name                     = "example-pool"
  max_active_requests      = "100"
  max_per_client_backlog   = "10"
}
```

## Argument Reference

* `name` (Required) - Rate Control Pool Name.
* `max_active_requests` (Optional) - Maximum Active Requests.
* `max_per_client_backlog` (Optional) - Maximum Per Client Backlog.
* `max_unconfigured_clients` (Optional) - Maximum Unconfigured Clients.
