# barracudawaf_network_acl

`barracudawaf_network_acl` manages `Network ACL` on the Barracuda Web Application Firewall.

## Example Usage

```hcl
resource "barracudawaf_network_acl" "example" {
  name     = "allow-ssh"
  priority = "10"
  action   = "ACCEPT"
  protocol = "TCP(6)"
  destination_port = "22"
}
```

## Argument Reference

* `name` (Required) - Name.
* `priority` (Required) - Priority.
* `action` (Optional) - Action.
* `comments` (Optional) - Comments.
* `destination_address` (Optional) - Destination IP Address.
* `destination_netmask` (Optional) - Destination Netmask.
* `destination_port` (Optional) - Destination Port Range.
* `enable_logging` (Optional) - Log Status.
* `icmp_response` (Optional) - ICMP Response.
* `interface` (Optional) - Interface.
* `ip_version` (Optional) - IP Protocol Version.
* `ipv6_destination_address` (Optional) - Destination IP Address (IPv6).
* `ipv6_destination_netmask` (Optional) - Destination Netmask (IPv6).
* `ipv6_source_address` (Optional) - Source IP Address (IPv6).
* `ipv6_source_netmask` (Optional) - Source Netmask (IPv6).
* `max_connections` (Optional) - Max Number of Connections.
* `max_half_open_connections` (Optional) - Max Connection Rate.
* `protocol` (Optional) - Protocol.
* `source_address` (Optional) - Source IP Address.
* `source_netmask` (Optional) - Source Netmask.
* `source_port` (Optional) - Source Port Range.
* `status` (Optional) - Enabled.
* `vsite` (Optional) - Network Group.
