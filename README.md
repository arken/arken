~# arken
A Tool to Create a Distributed Filesystem to Securely Backup the World's Important Data.

## What is the idea behind this project?

Many researchers, museums, and archivalists are struggling to host and protect a vast amount of important public data. On the other hand, there are many of us developers, tinkerers, and general computer enthusists who have extra storage space on our home servers.

My idea for Arken is to build an autonomous system for organizing, balancing, and distributing this data to users who can donate their extra space. 

## So what is Arken?

Arken is (for now) a command line application that runs on each client donating space to the project. Arken determines which files from the Keyset are most at risk and should be backed up to the client while respecting the amount of disk space and network usage specified by the user. Under the hood, Arken is simply a manger for IPFS (the Inter Planetary File System) which is an awesome distributed network that syncs files between nodes.

### What's a Keyset?

A Keyset is a list of files that a pool of Arken workers should keep track of and backup if needed. Keysets are git repositories of file identifiers used by IPFS no actual data is included in the keyset. 

### What's the process as someone who want's to backup important data?

Let's say that you are a scholar who want's to perseve some important works of humanity, or a researcher who want's to backup the DNA of an extinct animal/plant. How would you go about adding your data to the distributed file system? First, you would load the Arken Import Tool, and point it at the directory of important data. The tool will now create a keyset file of the IPFS identifiers for your data. At this point you can either upload the keyset to your own git repository (best if you want to run your own pool of workers) or make an application to put your data in the core keyset repository. This core keyset repository will be made up of extremely important data to preserve and is what the community donating their extra disk space will use by default.

### What's the process as someone donating their extra storage space?

Let's say that you have an old computer or server that has some storage space sitting empty on it. After installing the Arken program, you'll be asked how much storage space and network bandwidth you're willing to donate to the community pool. The Core Keyset will be available by default but because Keysets are just github repositories you can add and use any keysets you'd like. (For example, you're donating space to the core community pool but also want to sync a custom Keyset of some vacation pictures between you and a few friends' machines.) After the configuration that's it! Arken will continue to run in the background, determining files with the fewest number of other nodes hosting them and priorize backing them up to your node.

### Rebalancing Data Across the Comunity

On a regular interval an instance of Arken will query IPFS for the number of other nodes hosting a particular file and if any files are below the optimal threshhold Arken will look to replace one well backed up file on the system with an at risk file.