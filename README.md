# Arken

A Tool to Create a Distributed Filesystem to Securely Backup the World's Important Data.

## Table of Contents

- [A Bit of Backstory](#a-bit-of-backstory)
- [What is Arken?](#what-is-arken)
  - [What's a Keyset](#what's-a-keyset)
    - [Keyset Security](#keyset-security)
    - [Rebalancing Data Across the Comunity](#rebalancing-data-across-the-comunity)
- [Getting Started](#getting-started)
- [What's the process as someone who want's to backup important data?](#what's-the-process-as-someone-who-want's-to-backup-important-data?)
- [What's the process as someone donating their extra storage space?](#what's-the-process-as-someone-donating-their-extra-storage-space?)

## A Bit of Backstory

Many Researchers, Museums, and Archivists are struggling to host and protect a vast amount of important public data. On the other hand, there are many of us developers, tinkerers, and general computer enthusists who have extra storage space on our home servers.

The idea for Arken is to build an autonomous system for organizing, balancing, and distributing this data to users who can donate their extra space. 

```
+------------GitHub/GitLab/GitTea------------+
|    +--------+   +--------+   +--------+    |
|    | KeySet |   | KeySet |   | KeySet |    |
|    +----|---+   +----|\---+   +----|---+    |
+---------|------------|-\-----------|--------+
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

Arken is management engine that runs on top of the IPFS (Interplanetary File System) protocol. Each instance of Arken calculates which important files are hosted by the fewest number of other nodes on the network and should be locally backed up to reduce the risk of data loss. Arken also knows how much space it's using on your system and will respect limits you set by looking to locally delete data that is backed up by more than 10% of the cluster. 

### What's a Keyset?

A Keyset is how Arken transparently keeps track of knowing which files are important to the network and thus should be monitored and backed up if needed. Unlike a Pinset in an IPFS cluster,  a Keyset is simply a plain text git repository made of up file identifiers which makes it easy to audit so you actually know what data you're helping to preserve. Keyset repositories can contain an arbitrary number of directories used to organize keyset files as long as they also contain a keyset.config YAML file. This config file provides both a lighthouse file identifier used to measure the total number of nodes subscribed to that keyset, and a replication factor that is the percentage of the total network that should be storing a file at any one time.

While Keysets tell Arken which files should be stored on the subscribed nodes. A Keyset doesn't contain any of the data to backup onto the network. To import data to a Keyset users add files to IPFS and record the File Identifiers (IPFS CID) to a keyset file. From there nodes will begin pulling data directly from the user to the cluster.

##### Keyset Security

Because Keysets are openly available through Git repositories they can be easily audited and can only be changed by users who have access to those git repositories or through pull requests.

##### Rebalancing Data Across the Comunity

On a regular interval an instance of Arken will query IPFS for the number of other nodes hosting a particular file and if any files are below the optimal threshhold Arken will look to replace one well backed up file on the system with an at risk file.

## Getting Started

To start running an Arken Node you can download it as a Golang program or as a docker container. **For simplicity and easy updates it is recommended to run Arken as a docker container.** 

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
 --restart=always arkenproject/arken
```

##### Go Package:

```
go get github.com/arkenproject/arken
go run arken
```

### What's the process as someone who want's to backup important data?

Let's say that you are a scholar who want's to perseve some important works of humanity, or a researcher who want's to backup the DNA of an extinct animal/plant. How would you go about adding your data to the distributed file system? First, you would load the Arken Import Tool, and point it at the directory of important data. The tool will now create a keyset file of the IPFS identifiers for your data. At this point you can either upload the keyset to your own git repository (best if you want to run your own pool of workers) or make an application to put your data in the core keyset repository. This core keyset repository will be made up of extremely important data to preserve and is what the community donating their extra disk space will use by default.

### What's the process as someone donating their extra storage space?

Let's say that you have an old computer or server that has some storage space sitting empty on it. After installing the Arken program, you'll be asked how much storage space and network bandwidth you're willing to donate to the community pool. The Core Keyset will be available by default but because Keysets are just github repositories you can add and use any keysets you'd like. (For example, you're donating space to the core community pool but also want to sync a custom Keyset of some vacation pictures between you and a few friends' machines.) After the configuration that's it! Arken will continue to run in the background, determining files with the fewest number of other nodes hosting them and priorize backing them up to your node.
