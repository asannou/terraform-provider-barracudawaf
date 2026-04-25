# barracudawaf_cluster

`barracudawaf_cluster` manages `Cluster` on the Barracuda Web Application Firewall.

## Example Usage

```hcl
resource "barracudawaf_cluster" "example" {
  cluster_name          = "waf-cluster"
  cluster_shared_secret = "secret123"
}
```

## Argument Reference

* `cluster_name` (Required) - Cluster Name.
* `cluster_shared_secret` (Optional) - Cluster Shared Secret.
* `data_path_failure_action` (Optional) - Data-path failure action.
* `failback_mode` (Optional) - Failback Mode.
* `heartbeat_count_per_interface` (Optional) - Heartbeat Count Per Interface.
* `heartbeat_frequency` (Optional) - Heartbeat Frequency.
* `monitor_link` (Optional) - Monitor Link.
* `transmit_heartbeat_on` (Optional) - Transmit Heartbeat On.
* `vx_aa_enable` (Optional) - Active-Active HA Clustering.
