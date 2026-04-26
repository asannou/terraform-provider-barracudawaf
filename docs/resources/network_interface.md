# barracudawaf_network_interface

`barracudawaf_network_interface` manages `Network Interface` on the Barracuda Web Application Firewall.

## Example Usage

```hcl
resource "barracudawaf_network_interface" "wan" {
  name                    = "WAN"
  duplexity               = "Full Duplex"
  speed                   = "1000 Mbps"
  auto_negotiation_status = "On"
}
```

## Argument Reference

* `name` (Required) - NIC Card Name.
* `auto_negotiation_status` (Optional) - Auto-Negotiation Status.
* `duplexity` (Required) - NIC Cards Duplexity.
* `speed` (Required) - Default System Log Level.
