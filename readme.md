# Account Credential Manager (Go gRPC)

This project provides a gRPC-based account credential management service implemented in Go.

## Features


## Service Overview

This project provides two backend implementations for account management:

- **Database-backed service (`AccountsServiceDB`)**: Stores accounts persistently using a database. Recommended for production use where data durability is required.
- **In-memory service (`AccountsServiceMem`)**: Stores accounts in memory for fast, ephemeral operations. Useful for testing or lightweight deployments.

Both services implement the same gRPC API, supporting:
  - Account creation
  - Token retrieval and regeneration
  - Token expiry checks
  - Listing all accounts

Error handling uses gRPC status codes, and all operations are logged for transparency.

You can easily switch between backends depending on your needs.

## API Definition

The gRPC service and message definitions are located in [`./api/proto/v1/accounts.proto`](./api/proto/v1/accounts.proto).

## Getting Started

1. **Install dependencies**
    ```sh
    make build
    ```

2. **Run the server**
    ```sh
    make run
    ```

## Usage

Refer to the [`accounts.proto`](./api/proto/v1/accounts.proto) file for service methods and message formats.