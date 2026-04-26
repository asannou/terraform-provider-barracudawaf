# barracudawaf_network_vlan

`barracudawaf_network_vlan` manages `Network VLAN` on the Barracuda Web Application Firewall.

## Example Usage

```hcl
resource "barracudawaf_network_vlan" "example" {
  name      = "vlan10"
  interface = "WAN"
  vlan_id   = "10"
  vsite     = "default"
}
```

## Argument Reference

* `name` (Required) - VLAN Name.
* `interface` (Required) - VLAN Interface.
* `vlan_id` (Required) - VLAN ID.
* `vsite` (Required) - Network Group.
* `comments` (Optional) - Comments.
