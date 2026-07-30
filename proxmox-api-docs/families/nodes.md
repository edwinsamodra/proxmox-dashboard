# Nodes Family

Scope: `/nodes/{node}/*`

Best fit: node-local compute, guest CRUD, guest actions, storage, network, disks, services, telemetry, and shell helpers.

Main subtrees:

- qemu
- lxc
- ceph
- storage
- disks
- sdn
- firewall
- scan
- apt
- certificates
- services
- network
- status
- version
- time
- dns
- hosts
- tasks
- syslog
- journal
- rrd
- rrddata

Type pattern:

- `qemu` and `lxc` are the canonical collection + item + action trees.
- Create/update payloads are usually large and optional-heavy.
- Read endpoints usually return resource objects.
- Console and agent helpers usually return session wrappers or strings.

Shared types to keep:

- `NodeInventory`
- `NodeStatus`
- `NodeVersion`
- `NodeNetworkInterface`
- `NodeStorageEntry`
- `NodeDisk`
- `NodeService`
- `NodeCephStatus`
- `NodeTask`
- `NodeVm`
- `NodeContainer`
- `VmCreateRequest`
- `VmUpdateRequest`
- `ContainerCreateRequest`
- `ContainerUpdateRequest`
- `TaskResponse`

Useful path params:

- `{node}`
- `{vmid}`
- `{upid}`
- `{pos}`
- `{name}`

Representative endpoints:

- `GET /nodes/{node}`
- `GET /nodes/{node}/qemu`
- `POST /nodes/{node}/qemu`
- `GET /nodes/{node}/lxc/{vmid}`
- `POST /nodes/{node}/qemu/{vmid}/status/start`