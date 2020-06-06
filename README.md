# FsX: Distributed File System

FsX is designed as robust, distributed, scalable file system. It is perfect for
infrequent-updated, large volume file system, e.g., photo repository.

Each file repository is fully managed by a process `fsx`. At each time point,
there could be one of the following cases:

1. There is no `fsx` online. All of them are offline.
1. There is some `fsx` online. But only one of them is `master`; others are
   `replica`s.
1. There should be one file system which is allowed be changed out of
   management. This is the `master`.
1. All other file systems should be updated in passive mode, only via `fsx`.
