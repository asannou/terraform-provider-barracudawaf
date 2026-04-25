# barracudawaf_destination_nat

`barracudawaf_destination_nat` manages `Destination NAT` on the Barracuda Web Application Firewall.

## Example Usage

```hcl
resource "barracudawaf_destination_nat" "example" {
  pre_destination_address  = "203.0.113.10"
  pre_destination_netmask  = "255.255.255.255"
  post_destination_address = "192.168.1.10"
  incoming_interface       = "WAN"
  protocol                 = "TCP"
  vsite                    = "default"
}
```

## Argument Reference

* `pre_destination_address` (Required) - Pre-DNAT Destination.
* `incoming_interface` (Required) - Incoming Interface.
* `post_destination_address` (Required) - Post-DNAT Destination.
* `pre_destination_netmask` (Required) - Pre-DNAT Destination Mask.
* `protocol` (Required) - Protocol.
* `vsite` (Required) - Network Group.
* `pre_destination_port` (Optional) - Destination Port.
* `comments` (Optional) - Comments.
