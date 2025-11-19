# Distributed File Storage System

A peer-to-peer (P2P) distributed file storage system built in Go that enables multiple nodes to store, retrieve, and replicate files across a decentralized network.

## Overview

This project implements a distributed file storage system where multiple server nodes can connect to each other, share files, and maintain data redundancy. Each node acts as both a client and server, allowing files to be stored locally and replicated across connected peers in the network.

## Features

### Core Functionality
- **P2P Network Communication**: TCP-based transport layer with custom message encoding using Gob serialization
- **Distributed Storage**: Files are stored across multiple nodes with automatic replication
- **Content-Addressed Storage (CAS)**: Uses SHA-1 hashing for deterministic file paths and deduplication
- **File Replication**: Automatic broadcasting of stored files to all connected peers
- **Network Bootstrap**: Nodes can join existing networks by connecting to bootstrap nodes
- **Concurrent Peer Management**: Thread-safe peer tracking and management with RWMutex locks

### Technical Highlights
- **Custom Transport Layer**: Implements TCP-based peer-to-peer transport with handshaking
- **Message-Based Protocol**: Supports `MessageStoreFile` and `MessageGetFile` for distributed operations
- **Pluggable Path Transformer**: Flexible file path generation (supports both flat and hierarchical storage)
- **Graceful Error Handling**: Comprehensive error handling for network and storage operations
- **Stream-Based I/O**: Efficient streaming of files between nodes using Go's io.Reader/io.Writer interfaces
- **AES Encryption Support**: Files can be encrypted using AES-CTR cipher mode with 256-bit keys
- **Unicast & Broadcast Messaging**: Supports both unicast and broadcast communication patterns
- **Binary Protocol**: Efficient file size transmission using binary encoding (Little Endian)

## Architecture

### Components

**1. P2P Transport Layer** (`p2p/`)
- `tcp_transport.go`: TCP server that listens for incoming connections and manages peer connections
- `handshaker.go`: Handshake protocol for establishing peer connections
- `decoder.go`: Message decoding logic using Gob serialization
- `message.go`: Message type definitions for network communication

**2. File Server** (`server.go`)
- Orchestrates file storage and retrieval operations
- Manages connected peers
- Broadcasts messages to peers
- Handles incoming RPC messages
- Bootstrap node connection management

**3. Storage Layer** (`store.go`)
- `Store` interface for file operations (Read, Write, Delete, Has, Clear)
- `CASPathTransformer`: Hierarchical directory structure based on SHA-1 hash of file keys
- `DefaultPathTransformer`: Simple flat directory structure

**4. Encryption Module** (`crypto.go`)
- `copyEncrypt()`: Encrypts data using AES-CTR with random IV prepended
- `copyDcrypt()`: Decrypts AES-CTR encrypted data with embedded IV
- `newEncrptionKey()`: Generates random 256-bit encryption keys

## File Structure

```
├── p2p/                    # P2P networking layer
│   ├── tcp_transport.go
│   ├── tcp_transport_test.go
│   ├── handshaker.go
│   ├── decoder.go
│   ├── message.go
│   └── transport.go
├── main.go                 # Entry point and server setup
├── server.go               # FileServer implementation
├── store.go                # Storage layer implementation
├── crypto.go               # Encryption/Decryption utilities
├── store_test.go           # Storage tests
├── crypto_test.go          # Encryption tests
├── go.mod                  # Go module definition
└── Taskfile.yml            # Task automation
```

## Getting Started

### Prerequisites
- Go 1.16 or higher
- TCP ports available for server listening (default: 4001, 5001)

### Building
```bash
go build -o bin/fs
```

### Running

**Start the first node (bootstrap node):**
```bash
go run main.go
```

**Start additional nodes connected to the bootstrap node:**
```bash
# In the main.go, configure bootstrap nodes
fs2 := makeServer(":4001", ":5001")  // Connects to node on port 5001
```

### Example Usage
```go
// Store a file across the network
data := bytes.NewReader([]byte("File content"))
fs2.StoreData("my-file-key", data)

// Retrieve a file (from local storage or network)
file, err := fs.Get("my-file-key")
```

## Storage Strategy

### Content-Addressed Storage (CAS)
Files are organized using SHA-1 hashing with hierarchical directory structure:
- Hash: `b300938bfa8fe3c2a10eccd85aeac38fc9f15d50`
- Path: `5001_DFA/b3009/38bfa/8fe3c/2a10e/ccd85/aeac3/8fc9f/15d50/b300938bfa8fe3c2a10eccd85aeac38fc9f15d50`

This approach provides:
- Automatic deduplication
- Deterministic file locations
- Distributed load balancing

## Network Protocol

### Message Types
1. **MessageStoreFile**: Notifies peers of a new file to store
   - `Key`: File identifier
   - `Size`: File size in bytes

2. **MessageGetFile**: Requests a file from peers
   - `Key`: File identifier to retrieve

### Communication Patterns
- **Broadcast**: Sends messages to all connected peers simultaneously
- **Unicast**: Sends messages to specific peer nodes
- **Stream Protocol**: Uses binary encoding for file size transmission with `IncomingStream` message type

### Communication Flow
1. Node stores a file locally via `StoreData(key, reader)`
2. Broadcasts `MessageStoreFile` metadata to all connected peers
3. Sends actual file data via stream to each peer
4. Peers receive and store the file using `handleMessageStoreFile()`
5. When retrieving via `Get(key)`, node checks local storage first
6. If not found locally, broadcasts `MessageGetFile` request to network
7. Peers respond with file size (binary encoded) followed by file data
8. Node writes received data to local storage and returns reader

## Current Status

### Implemented
- ✅ P2P TCP transport with handshaking
- ✅ File storage and retrieval with CAS
- ✅ Peer connection management
- ✅ Message encoding/decoding (Gob)
- ✅ File replication across peers
- ✅ Bootstrap node connectivity
- ✅ Concurrent request handling
- ✅ Unicast and broadcast messaging
- ✅ Binary protocol for efficient size transmission
- ✅ AES-CTR encryption functions

### In Progress / Planned
- 🔄 Integrate encryption into file storage/retrieval pipeline
- 🔄 Persistence of peer metadata
- 🔄 File versioning and conflict resolution
- 🔄 Data compression support
- 🔄 Performance optimization and benchmarking
- 🔄 Comprehensive test coverage for encryption
- 🔄 REST API layer
- 🔄 Multi-stream parallel file transfer

## Testing
```bash
go test ./...
```

## Technologies & Concepts

- **Go**: Concurrency with goroutines and channels
- **Networking**: TCP sockets, custom protocol design, binary protocol
- **Distributed Systems**: P2P architecture, replication, eventual consistency
- **Serialization**: Gob encoding for binary message format
- **File I/O**: Stream-based I/O with io.Reader/io.Writer interfaces
- **Cryptography**: SHA-1 hashing for content addressing, AES-CTR encryption
- **Encryption**: AES cipher in CTR mode with random initialization vectors
- **Concurrency Patterns**: RWMutex for thread-safe operations, channel-based communication

## Future Enhancements

1. **Transparent Encryption**: Automatically encrypt/decrypt files during storage operations
2. **Consensus Protocol**: Implement RAFT or similar for consistency guarantees
3. **Delete Propagation**: Ensure deleted files are removed from all peers
4. **Network Topology**: Implement Kademlia DHT for better peer discovery
5. **Monitoring & Metrics**: Add Prometheus metrics for network health
6. **REST API**: HTTP interface for easier interaction
7. **Docker Support**: Containerization for easier deployment
8. **TLS/SSL**: Secure peer communication with certificate-based authentication
9. **Key Management**: Distributed key management for encryption
10. **Data Compression**: GZIP or Snappy compression for bandwidth optimization

## Author
Apurv