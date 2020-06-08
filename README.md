# FsX: Distributed File Backup System

FsX is designed as a robust, distributed, scalable,
unidirectional-synchronization file backup system. It is perfect for
infrequent-modification, large volume file system with replicas, e.g., photo
repository.

It is designed with the following features in mind:
1. full history in logs. 
2. scalability
3. simple diff cross replicas.

Each file repository is fully managed by a process `fsx`, with exception for
`master` node.

1. There are some `fsx`s online. But only one of them is `master` node; others
   are `replica` nodes.
1. There should be one file system which is allowed be changed out of
   management. This is the `master` node. This is ensured by the user, which
   should be trivial.
1. All other file systems should be updated in passive mode, only via `fsx`.

## Cmd Logs

Each node has a Cmd Logs, which is a sequence of the Cmds. Following the Cmds,
it allows any node creating the same file system.

```go
type NodeState struct {
  NodeName    string  // Unique in the cluster.
  NextVersion uint64
  CmdLogs     CmdLogs // Len(CmdLogs) == NextVersion
}

type CmdLogs []CmdLog

type CmdLog struct {
  Version     uint64
  MasterName  string  // Name of master node for this Cmd.
  Cmd         Cmd
}

type Cmd struct {

}
```



## Master Node

During start up, `master` node detects the file system change. It checks file
path and size as sanity check.
