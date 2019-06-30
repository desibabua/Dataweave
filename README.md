# Problem statement

Design Problem: DWRedis
Design an in memory <key, value> store such as Redis. Let's call it DWRedis.

Requirements:
- It should support the major Redis operations such as SET, GET, and SAVE.
- There should be a service that supports start, restart, and stop of DWRedis.
- When the service is started it should load data from a file at a standard location.
- When the service is stopped, it should dump data to the above file.

Some questions you need to consider
- How are you going to implement TTL?
- How will you manage a keyspace that is larger than the memory allocated for
DWRedis?
- How will you implement versioning and checkpointing?
- How are you going to handle crash recovery?
- Assume that you can run your DWRedis in a clustered mode. If a new node is added,
how are you going to redistribute the keyspace?
- How fast is your DWRedis going to be? Can you estimate the throughput--with proper
reasoning--without implementing and testing this?

What you need to deliver?
- A design doc that has the following:
- choice of data structures
 - design details for each operation
- schematics for different workflows (load, dump, checkpointing, crash recovery)
- design considerations to ensure performance, reliability, correctness, and
scalability
- Answers to the questions in the previous section

# Solution

#### Data structure to be used
- Map<Key, pointer to the value object>
- Value object will have following properties
```
Class Value {
    private T key;
    private T value;
    private Date expiryDate; // Will be 0 if TTL is not set
    private String idSequence; // Will be used in case of distributed system to achieve consistency
}
```

- Id sequence is a unique id to the record across the cluster of nodes
- It can be timestamp in ns + node id + unique id per node

#### Insertion algorithm (SET Operation)
```
func insert (T key, T value) {
   if(key is persent in map) {
      return
   }
  node = create_node(key, value)
  map.insert(key, node)
}
```

#### Fetch Algorithm (GET Operation)
```
func fetch (T key) T {
    boolean node = map[key]
    if (node == null) {
        // return 403 or key not found exception
    }
    if (node.expiration_time < curr_time) {
        // return 403 or key not found exception
    }
   
    return node.val
}
```

#### SET Expiration Algorithm (EXPIRE <key> <TTL>)
```
func expire (T key, Date expiry) {
    boolean node = map[key]
    if (node == null) {
        // return 403 or key not found exception
    }
    if (node.expiration_time < curr_time) {
        // return 403 or key not found exception
    }
   
    node.setExpiration(expiry)
}
```

#### Expiry based eviction algorithm (TTL Implementation)
For achieving expiry based eviction we can have a cron job for removing the expired records(as below) and
also while serving the record we need to check the expiry (already done in GET operation)
```
// Cron Job
func evictAllExpired() {
    for key, valNode := range(map) {
        if(valNode..expiration_time < curr_time) {
            map.remove(key)
            free(valNode)
        }
    }
}
```

#### Save Operation - dump workflow
Save operation will just serialize the key, value in a file
```
func save() {
    file = newFile(name: `curr-date`+"_checkpoint")
    for key, valNode := range(map) {
        file.addLine(key, valNode.val, valNode.expire_time)
    }
    file.save()
}
```

#### Restore Operation - save workflow
Restore operation will just de-serialize the file and store it in memory
```
func restore(fileName) {
    file = openFile(fileName)
    for line as Node file.next {
        cmd.SET(line.key, line.value, line.expire_time)
    }
    file.close()
}
```

#### Systemctl for service
```diff
 - Note: Not sure what exactly the ask is assuming we need systemd file(attached below)</red>
```

```
[Unit]
Description=DWRedis
Requires=network.target
After=network.target
[Service]
ExecStart=/bin/sh -c '$install_root_dir/DWRedis-linux start $install_root_dir/DWRedis-linux.yaml >> $install_root_dir/logs/DWRedis.log 2>&1'
ExecReload='/bin/kill -HUP \$MAINPID'
KillSignal=TERM
Restart=on-failure
User=$service_name
RestartSec=30
WorkingDirectory=$install_root_dir
```

For this to work we need to have a linux exec for that something like a cli application which is bundled as linux executable.

Also in case of start/stop/restart we need to save/read the files from a fixed location. The workflow for that is mentioned below.
```
func start(backupFilePath) {
    restore(backupFilePath)
}
func stop(backupFilePath) {
    save(backupFilePath)
    closeAllConnections()
}
func restart(backupFilePath) {
    stop(backupFilePath)
    start(backupFilePath)
}
```

#### How to make the above system distributed?

Things to consider for designing distributed DWRedis
- key space should be distributed between nodes as all the keys can't fit in a single machine
- The checkpoints and persisted file should be distributed(somewhere in HDFS)
- Each key set should be partitioned and replicated
- Also to come up with the consistency pattern I am assuming that Read to Write ratio to be 1: 1000 i.e. system is read heavy.
So in that case I will go be strong consistency pattern.

Q. How does the distributed system look like?

A. 
![cached image](http://www.plantuml.com/plantuml/proxy?src=https://raw.githubusercontent.com/ravik-karn/Dataweave/master/architecture.puml)

Following is the usage of all the components used
- Load Balancer will be used to distributed the workload between multiple Request manager.
 We can use Network layer load balancing for routing the request to nearest server.
 Within in the same node it can use LRU/ Round Robin/ Sticky session etc. to distribute the load. 
- Metadata manger is to keep track about the keyset and it's primary/secondary nodes(ISRs - In Sync Replicas).
Metadata manager is also responsible for electing the primary node for a particular partition. 
For electing the primary node for a particular partition we can use PAXOS algorithm.
- Controller keeps polling about the health of the system to all the storage node. 
Once a node is unhealthy then it sends a message to Metadata manger and metadata manager initiates request for leader election.
Controller is also responsible to spin up new storage node and redistributing keys among the nodes.
- Request manager talks applys a hashing function to find the primary node responsible for the particular keyset(described below). - This info can be cached for performance. 
After getting the node id it requests corrosponding storage node for get/put request. 
- Storage node responsible for storing the key & value map(as described above) 
It also creates backups of it's data at defined times(check-points).
During initialization it loads the data from the HDFS and during shutdown it stores the data back.  

Q. How to distribute large key space?

A. We will divide the data using consistent hashing. 
```
func getNode(T key) Strinng {
    availableNode = zookeeper.getRegisterdNodes()
}
```  

Q. How to implement versioning?

A. For implementing versioning we can keep track of all the write request comming to the Request managers and while rolling to   

