# barracudawaf_interface_route

`barracudawaf_interface_route` manages `Interface Route` on the Barracuda Web Application Firewall.

## Example Usage

```hcl
resource "barracudawaf_interface_route" "example" {
  ip_address = "192.168.100.0"
  netmask    = "255.255.255.0"
  interface  = "WAN"
  vsite      = "default"
}
```

## Argument Reference

* `ip_address` (Required) - IP/Network Address.
* `interface` (Required) - Network Interface.
* `netmask` (Required) - Netmask.
* `vsite` (Required) - Network Group.
* `ip_version` (Optional) - IP Version.
* `comments` (Optional) - Comments.
