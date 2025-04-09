# gNMI Path Converter

This is a Go project that provides utilities to:

- Convert YANG-style XPATH keys (commonly stored in Redis) into gNMI-encoded `Path` proto messages for use at the **system side (Target)**.
- Convert gNMI `Path` proto messages back into YANG-style XPATH strings for the ** Broker e.g. collector** (Server side).

The tool is helpful for systems using model-driven telemetry based on [gNMI](https://github.com/openconfig/reference/blob/master/rpc/gnmi/gnmi-specification.md) and [OpenConfig YANG models](https://github.com/openconfig/public).

## Features

- Parse YANG-style XPATH into gNMI Path proto
- Convert gNMI Path proto into readable OpenConfig XPATH strings
- Supports multi-key path segments (`foo[bar=val][baz=val2]`)
- Ready for integration with Redis pipelines and telemetry backends

## Example 
```bash

Input XPATH: /switch/clusters/cluster[name=cluster-1]/ports/port[id=eth0]/interfaces/interface[name=eth0.100]/state/counters

Converted to gNMI Path: elem:  {
  name:  "switch"
}
elem:  {
  name:  "clusters"
}
elem:  {
  name:  "cluster"
  key:  {
    key:  "name"
    value:  "cluster-1"
  }
}
elem:  {
  name:  "ports"
}
elem:  {
  name:  "port"
  key:  {
    key:  "id"
    value:  "eth0"
  }
}
elem:  {
  name:  "interfaces"
}
elem:  {
  name:  "interface"
  key:  {
    key:  "name"
    value:  "eth0.100"
  }
}
elem:  {
  name:  "state"
}
elem:  {
  name:  "counters"
}

Converted back to XPATH: /switch/clusters/cluster[name=cluster-1]/ports/port[id=eth0]/interfaces/interface[name=eth0.100]/state/counters
```

