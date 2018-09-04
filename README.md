# Conch(贝壳)

[![Build Status](https://travis-ci.com/blockchainworkers/conch.svg?branch=master)](https://travis-ci.com/blockchainworkers/conch)
[![GoDoc](https://godoc.org/github.com/blockchainworkers/conch?status.svg)](https://godoc.org/github.com/blockchainworkers/conch)


一个基于tendermint的共识机制引擎实现的简单虚拟货币 

## 一些相关知识的链接

[拜占庭容错](https://en.wikipedia.org/wiki/Byzantine_fault_tolerance)

[状态机副本](https://en.wikipedia.org/wiki/State_machine_replication)

[区块链](https://en.wikipedia.org/wiki/Blockchain_(database))


## 什么是tendermint
Tendermint Core 是一个拜占庭容错中间件, 它实现了一套共识机制引擎. 可以基于此项目使用任意的编程语言实现基于此共识机制的各种项目。如分布式系统, 虚拟货币等. 下面是一些基于此引擎开源的一些明星项目. 我们使用此共识机制更大的目的是在于学习。 感谢[tendermint](https://tendermint.com)这个伟大的开源项目. 让开发区块链工程更easy, 也能从此开源项目中学习到很多有用的知识.

[cosmos network](https://cosmos.network/)

[ethermint 基于POS+BFT以太坊实现](https://github.com/cosmos/ethermint)

[Hyperledger Burrow](https://github.com/hyperledger/burrow)


### 测试链已经启动

#### 使用docker镜像连接到测试链

- 创建配置和数据存储路径 
    ```shell
        mkdir /opt/conch/config
        mkdir /opt/conch/data
    ```
- 将genesis.json, config.toml文件放入/opt/conch/config
    > genesis.json
    ```json
        {
            "genesis_time": "2018-08-31T09:19:29.144335522Z",
            "chain_id": "conch-testnet-wupeaking",
            "consensus_params": {
            "block_size_params": {
                "max_bytes": "22020096",
                "max_txs": "10000",
                "max_gas": "-1"
            },
            "tx_size_params": {
                "max_bytes": "10240",
                "max_gas": "-1"
            },
            "block_gossip_params": {
                "block_part_size_bytes": "65536"
            },
            "evidence_params": {
                "max_age": "100000"
            }
            },
            "validators": [
            {
                "pub_key": {
                "type": "tendermint/PubKeySecp256k1",
                "value": "A7C1pYP/mrQ6Jnp3oQMpAVKpUOnAQjKpLA95e7MbV/eR"
                },
                "power": "10",
                "name": "zhangsan"
            }
            ],
            "app_hash": ""
        }
    ```

    > config.tmol
    ```tmol

        # TCP or UNIX socket address of the ABCI application,
        # or the name of an ABCI application compiled in with the Tendermint binary
        proxy_app = "conchapp"

        # A custom human readable name for this node
        moniker = "conchapp-pc"

        # If this node is many blocks behind the tip of the chain, FastSync
        # allows them to catchup quickly by downloading blocks in parallel
        # and verifying their commits
        fast_sync = true

        # Database backend: leveldb | memdb
        db_backend = "leveldb"

        # Database directory
        db_path = "data"

        # Output level for logging, including package level options
        log_level = "main:info,state:info,*:error"

        ##### additional base config options #####

        # Path to the JSON file containing the initial validator set and other meta data
        genesis_file = "config/genesis.json"

        # Path to the JSON file containing the private key to use as a validator in the consensus protocol
        priv_validator_file = "config/priv_validator.json"

        # Path to the JSON file containing the private key to use for node authentication in the p2p protocol
        node_key_file = "config/node_key.json"

        # Mechanism to connect to the ABCI application: socket | grpc
        abci = "socket"

        # TCP or UNIX socket address for the profiling server to listen on
        prof_laddr = ""

        # If true, query the ABCI app on connecting to a new peer
        # so the app can decide if we should keep the connection or not
        filter_peers = false

        ##### advanced configuration options #####

        ##### rpc server configuration options #####
        [rpc]

        # TCP or UNIX socket address for the RPC server to listen on
        laddr = "tcp://0.0.0.0:26657"

        # TCP or UNIX socket address for the gRPC server to listen on
        # NOTE: This server only supports /broadcast_tx_commit
        grpc_laddr = ""

        # Maximum number of simultaneous connections.
        # Does not include RPC (HTTP&WebSocket) connections. See max_open_connections
        # If you want to accept more significant number than the default, make sure
        # you increase your OS limits.
        # 0 - unlimited.
        grpc_max_open_connections = 900

        # Activate unsafe RPC commands like /dial_seeds and /unsafe_flush_mempool
        unsafe = false

        # Maximum number of simultaneous connections (including WebSocket).
        # Does not include gRPC connections. See grpc_max_open_connections
        # If you want to accept more significant number than the default, make sure
        # you increase your OS limits.
        # 0 - unlimited.
        max_open_connections = 900

        ##### peer to peer configuration options #####
        [p2p]

        # Address to listen for incoming connections
        laddr = "tcp://0.0.0.0:26656"

        # Address to advertise to peers for them to dial
        # If empty, will use the same port as the laddr,
        # and will introspect on the listener or use UPnP
        # to figure out the address.
        external_address = ""

        # Comma separated list of seed nodes to 
        seeds = ""
        persistent_peers = "bf3e958f7d8a935bff9d0d10396d4916009c4ee5@120.55.49.80:26656"

        # UPNP port forwarding
        upnp = false

        # Path to address book
        addr_book_file = "config/addrbook.json"

        # Set true for strict address routability rules
        addr_book_strict = true

        # Time to wait before flushing messages out on the connection, in ms
        flush_throttle_timeout = 100

        # Maximum number of peers to connect to
        max_num_peers = 50

        # Maximum size of a message packet payload, in bytes
        max_packet_msg_payload_size = 1024

        # Rate at which packets can be sent, in bytes/second
        send_rate = 5120000

        # Rate at which packets can be received, in bytes/second
        recv_rate = 5120000

        # Set true to enable the peer-exchange reactor
        pex = true

        # Seed mode, in which node constantly crawls the network and looks for
        # peers. If another node asks it for addresses, it responds and disconnects.
        #
        # Does not work if the peer-exchange reactor is disabled.
        seed_mode = false

        # Comma separated list of peer IDs to keep private (will not be gossiped to other peers)
        private_peer_ids = ""

        ##### mempool configuration options #####
        [mempool]

        recheck = true
        recheck_empty = true
        broadcast = true
        wal_dir = "data/mempool.wal"

        # size of the mempool
        size = 100000

        # size of the cache (used to filter transactions we saw earlier)
        cache_size = 100000

        ##### consensus configuration options #####
        [consensus]

        wal_file = "data/cs.wal/wal"

        # All timeouts are in milliseconds
        timeout_propose = 3000
        timeout_propose_delta = 500
        timeout_prevote = 1000
        timeout_prevote_delta = 500
        timeout_precommit = 1000
        timeout_precommit_delta = 500
        timeout_commit = 1000

        # Make progress as soon as we have all the precommits (as if TimeoutCommit = 0)
        skip_timeout_commit = false

        # EmptyBlocks mode and possible interval between empty blocks in seconds
        create_empty_blocks = true
        create_empty_blocks_interval = 20

        # Reactor sleep duration parameters are in milliseconds
        peer_gossip_sleep_duration = 100
        peer_query_maj23_sleep_duration = 2000

        ##### transactions indexer configuration options #####
        [tx_index]

        # What indexer to use for transactions
        #
        # Options:
        #   1) "null" (default)
        #   2) "kv" - the simplest possible indexer, backed by key-value storage (defaults to levelDB; see DBBackend).
        indexer = "kv"

        # Comma-separated list of tags to index (by default the only tag is tx hash)
        #
        # It's recommended to index only a subset of tags due to possible memory
        # bloat. This is, of course, depends on the indexer's DB and the volume of
        # transactions.
        index_tags = ""

        # When set to true, tells indexer to index all tags. Note this may be not
        # desirable (see the comment above). IndexTags has a precedence over
        # IndexAllTags (i.e. when given both, IndexTags will be indexed).
        index_all_tags = false

        ##### instrumentation configuration options #####
        [instrumentation]

        # When true, Prometheus metrics are served under /metrics on
        # PrometheusListenAddr.
        # Check out the documentation for the list of available metrics.
        prometheus = false

        # Address to listen for Prometheus collector(s) connections
        prometheus_listen_addr = ":26660"

        # Maximum number of simultaneous connections.
        # If you want to accept more significant number than the default, make sure
        # you increase your OS limits.
        # 0 - unlimited.
        max_open_connections = 3

    ```

- 启动docker镜像

    ```shell
    docker run --name=conch-node --net=host -v /opt/conch:/opt/conch/data -d blockchainworkers/conch 
    ```
#### 使用二进制文件加入测试链

- 编译二进制文件
    ```shell
    > mkdir -p /opt/conch/src/github.com/blockchainworkers
    > mkdir -p /opt/conch/data/config &&  mkdir -p /opt/conch/data/data
    > export BUILD_FLAGS=-ldflags "-X github.com/blockchainworkers/conch/version.GitCommit=`git rev-parse --short=8 HEAD`"
    > export GOPATH=/opt/conch
    > cd /opt/conch/src/github.com/blockchainworkers &&  git clone https://github.com/blockchainworkers/conch.git conch
    > cd conch
    > go build $(BUILD_FLAGS) -o conchd ./cmd/conch/main.go
    ```
- 将genesis.json, config.toml文件放入/opt/conch/data/config

- 启动节点

    ```shell
        ./conchd node --home /opt/conch/data
    ```

#### 关于API说明以及转账说明文档即将到来
    未完待续...
