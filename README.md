# Arken

<img src="https://avatars.githubusercontent.com/u/66809416?s=200&v=4">

A Distributed Digital Archive Built for the World's Open Source and Scientific Data.

[![Go Report Card](https://goreportcard.com/badge/github.com/arken/arken)](https://goreportcard.com/report/github.com/arken/arken)

## Table of Contents

- [A Bit of Backstory](#a-bit-of-backstory)
- [What is Arken?](#what-is-arken)
  - [What's a Keyset](#whats-a-keyset)
    - [Keyset Security](#keyset-security)
    - [Rebalancing Data Across the Community](#rebalancing-data-across-the-community)
- [Getting Started](#getting-started)
- [What's the process as someone who want's to backup important data?](#what's-the-process-as-someone-who-want's-to-backup-important-data?)
- [What's the process as someone donating their extra storage space?](#what's-the-process-as-someone-donating-their-extra-storage-space?)

## A Bit of Backstory

Many researchers, museums, and archivists are struggling to host and protect a vast amount of important public data. 
On the other hand, there are many of us developers, tinkerers, and general computer enthusiasts who have extra storage 
space on our home servers.

The goal of Arken is to build an autonomous system for organizing, balancing, and distributing this data among users who 
can donate their extra space. 

```
+------------GitHub/GitLab/GitTea------------+
|    +--------+   +--------+    +--------+   |
|    | Keyset |   | Keyset |    | Keyset |   |
|    +----|---+   +----|\--+    +----|---+   |
+---------|------------|-\-----------|-------+
          |            |  \          /\
          |            |   \        /  \
          |            |    \      /    \
          |            |     \    /      \
          |            |      \  /        \
          |            |       \/          \
          v            v        v           \
       [Arken]     [Arken]<-->[Arken]<--->[Arken]
```

# What is Arken?

Arken is a management engine that runs on top of the IPFS (Interplanetary File System) protocol. Each instance of Arken 
calculates which important files are hosted by the fewest number of other nodes on the network and should thus be 
locally backed up to reduce the risk of data loss. Arken also knows how much space it's using on your system and will 
respect limits you set by locally deleting data that is backed up by more than 10% of the cluster. 

### What's a Keyset?

Arken uses Keysets to transparently keep track of which files are important to the network and should be
monitored and backed up if needed. Unlike a Pinset in an IPFS cluster, a Keyset is simply a plain text git repository
made of up file identifiers. Additionally, Keysets are easy to audit so you can actually know what data you're helping
preserve. Keyset repositories can contain an arbitrary number of directories used to organize keyset files as long as 
they also contain a `keyset.config` YAML file. This config file provides both a lighthouse file identifier used to 
measure the total number of nodes subscribed to that Keyset, and a replication factor that is the percentage of the
total network that should be storing a file at any given time.

While Keysets tell Arken which files should be stored on the subscribed nodes, they don't contain any of the
data to be backed up onto the network. To import data to a Keyset, users add files to IPFS and record the File 
Identifiers (IPFS CID) to a keyset file. From there, nodes will begin pulling data directly from the user to the cluster.

##### Keyset Security

Since Keysets are openly available through Git repositories, they can be easily audited but can only be changed by 
users who have access to those Git repositories or through pull requests.

##### Rebalancing Data Across the Community

Arken instances will periodically query IPFS for the number of other nodes hosting a particular file and attempt to 
replace one well backed up file on the system with files below the optimal threshold.

## Getting Started

#### Tutorials:
[Getting Started with Arken on a Raspberry Pi](https://github.com/arken/arken/blob/master/docs/raspberry-pi-setup.md)

To start running a node, you can download Arken as a Golang program or as a Docker container. 
**It's recommended to run Arken as a Docker container for simplicity and ease of updating.** 

##### Docker:

```
docker run -d --name arken \
 -v STORAGE:/data/storage \
 -v DATABASE:/data/database \
 -v REPOSITORIES:/data/repositories \
 -v CONFIG:/data/config \
 -e ARKEN_GENERAL_POOLSIZE=2TB \
 -e ARKEN_DB_PATH=/data/database/keys.db \
 -e ARKEN_SOURCES_CONFIG=/data/config/keysets.yaml \
 -e ARKEN_SOURCES_REPOSITORIES=/data/repositories \
 -e ARKEN_SOURCES_STORAGE=/data/storage \
 -p 4001:4001 \
 --restart=always arken/arken
```

##### Go Package:

```
go get github.com/arken/arken
go run arken
```

### What's the process as someone who wants to back up important data?

Let's say that you are a scholar who wants to preserve some important works of humanity, or a researcher who wants 
to back up the DNA of an extinct animal/plant. How would you go about adding your data to the distributed file system? 
First, you would download & run the [Arken Import Tool](https://github.com/arken/ait). Using the Arken Import tool you can create 
a Keyset file of the IPFS identifiers for your data. At this point you can either upload the Keyset to your own Git 
repository (this is best if you want to run your own pool of workers) or make an application to put your data in the
Core Keyset repository. The Core Keyset repository consists of extremely important data to preserve and is what the 
community donating their extra disk space uses by default.

### What's the process as someone donating their extra storage space?

Old computers or servers with some empty storage space make excellent Arken nodes. Check out our guide for configuring a Raspberry Pi with Docker and External Storage Arken [here](https://github.com/arken/arken/blob/master/docs/raspberry-pi-setup.md). After installing the 
Arken program, you can configure it either through environment variables or the Arken configuration file located at `~/.arken/`. You can check out an example of an Arken Docker-Compose file [here](https://github.com/arken/arken/blob/master/docs/examples/docker-compose.yml). The core Keyset will be available by default, but because Keysets are just Git repositories, you can add and use 
any Keyset you'd like. For example, you can donate space to the core community pool but also sync a custom Keyset of 
some vacation pictures amongst yours and a few friends' machines.

After the configuration, that's it! Arken will continue to run in the background, determining files with the fewest 
number of other nodes hosting them and rebalancing as necessary.

## License

Copyright 2020-2021 Alec Scott & Arken Team <team@arken.io>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
