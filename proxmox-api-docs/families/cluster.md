# Cluster Family

Scope: `/cluster/*`

Best fit: cluster config, orchestration, metrics, notifications, firewall, HA, SDN, Ceph, and replication.

Main subtrees:

- sdn
- notifications
- firewall
- ha
- mapping
- acme
- config
- metrics
- qemu
- ceph
- jobs
- backup
- replication
- bulk-action
- resources

Type pattern:

- Collection endpoints usually return arrays.
- Config endpoints usually behave like item resources.
- Action endpoints often return `null` or a task/job token.

Shared types to keep:

- `ClusterReplicationJob`
- `ClusterMetricServer`
- `ClusterMetricExportResponse`
- `ClusterNotificationMatcher`
- `ClusterNotificationTarget`
- `ClusterNotificationEndpointBase`
- `ClusterFirewallRule`
- `ClusterFirewallGroup`
- `ClusterHaGroup`
- `ClusterHaResource`
- `ClusterBackupJob`
- `ClusterSdnZone`

Useful path params:

- `{id}`
- `{name}`
- `{sid}`
- `{vmid}`
- `{group}`
- `{pos}`

Representative endpoints:

- `GET /cluster/replication`
- `POST /cluster/replication`
- `GET /cluster/metrics/export`
- `GET /cluster/notifications/endpoints/smtp`
- `POST /cluster/notifications/targets/{name}/test`