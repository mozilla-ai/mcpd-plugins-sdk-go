# `mcpd` plugins SDK for Go

Go SDK for developing mcpd middleware plugins.

## Overview

This SDK provides Go types and gRPC interfaces for building middleware plugins that integrate with `mcpd`.
Plugins can process HTTP requests and responses, implementing capabilities like authentication, rate limiting,
content transformation, and observability.

## Installation

```bash
go get github.com/mozilla-ai/mcpd-plugins-sdk-go@latest
```

Then in your project:
```bash
go mod tidy
```

## Usage

The SDK provides two approaches: **with helpers** (recommended) for convenience, and **explicit** for full control.

### Option 1: Using SDK Helpers (Recommended)

The SDK provides the `Serve()` helper and `BasePlugin` struct to minimize boilerplate:

```go
package main

import (
	"context"
	"log"

	"github.com/mozilla-ai/mcpd-plugins-sdk-go/pkg/plugins/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MyPlugin struct {
	mcpdpluginsv1.BasePlugin // Provides sensible defaults for all methods.
}

func (p *MyPlugin) GetMetadata(ctx context.Context, _ *emptypb.Empty) (*mcpdpluginsv1.Metadata, error) {
	return &mcpdpluginsv1.Metadata{
		Name:        "my-plugin",
		Version:     "1.0.0",
		Description: "Example plugin that does something useful",
	}, nil
}

func (p *MyPlugin) GetCapabilities(ctx context.Context, _ *emptypb.Empty) (*mcpdpluginsv1.Capabilities, error) {
	return &mcpdpluginsv1.Capabilities{
		Flows: []mcpdpluginsv1.Flow{mcpdpluginsv1.FlowRequest},
	}, nil
}

func (p *MyPlugin) HandleRequest(ctx context.Context, req *mcpdpluginsv1.HTTPRequest) (*mcpdpluginsv1.HTTPResponse, error) {
	// Custom request processing logic here.

	return &mcpdpluginsv1.HTTPResponse{
		Continue:   true,
		StatusCode: 0,
		Headers:    req.Headers,
		Body:       req.Body,
	}, nil
}

func main() {
	if err := mcpdpluginsv1.Serve(&MyPlugin{}); err != nil {
		log.Fatal(err)
	}
}
```

**What `BasePlugin` provides:**
- `CheckHealth()` - returns OK
- `CheckReady()` - returns OK
- `Configure()`, `Stop()` - no-ops
- `HandleRequest()`, `HandleResponse()` - pass through unchanged

Override only the methods you need!

### Option 2: Explicit Implementation

For full control over the server lifecycle:

```go
package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"

	"github.com/mozilla-ai/mcpd-plugins-sdk-go/pkg/plugins/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MyPlugin struct {
	mcpdpluginsv1.UnimplementedPluginServer
}

func (p *MyPlugin) GetMetadata(ctx context.Context, _ *emptypb.Empty) (*mcpdpluginsv1.Metadata, error) {
	return &mcpdpluginsv1.Metadata{
		Name:        "my-plugin",
		Version:     "1.0.0",
		Description: "Example plugin",
	}, nil
}

func (p *MyPlugin) GetCapabilities(ctx context.Context, _ *emptypb.Empty) (*mcpdpluginsv1.Capabilities, error) {
	return &mcpdpluginsv1.Capabilities{
		Flows: []mcpdpluginsv1.Flow{mcpdpluginsv1.FlowRequest},
	}, nil
}

func (p *MyPlugin) CheckHealth(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (p *MyPlugin) CheckReady(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (p *MyPlugin) Configure(ctx context.Context, cfg *mcpdpluginsv1.PluginConfig) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (p *MyPlugin) Stop(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (p *MyPlugin) HandleRequest(ctx context.Context, req *mcpdpluginsv1.HTTPRequest) (*mcpdpluginsv1.HTTPResponse, error) {
	return &mcpdpluginsv1.HTTPResponse{
		Continue:   true,
		Headers:    req.Headers,
		Body:       req.Body,
	}, nil
}

func (p *MyPlugin) HandleResponse(ctx context.Context, resp *mcpdpluginsv1.HTTPResponse) (*mcpdpluginsv1.HTTPResponse, error) {
	return &mcpdpluginsv1.HTTPResponse{
		Continue:   true,
		StatusCode: resp.StatusCode,
		Headers:    resp.Headers,
		Body:       resp.Body,
	}, nil
}

func main() {
	var address, network string
	flag.StringVar(&address, "address", "", "gRPC address (socket path for unix, host:port for tcp)")
	flag.StringVar(&network, "network", "unix", "Network type (unix or tcp)")
	flag.Parse()

	if address == "" {
		log.Fatal("--address flag is required")
	}

	lis, err := net.Listen(network, address)
	if err != nil {
		log.Fatalf("failed to listen on %s %s: %v", network, address, err)
	}

	if network == "unix" {
		defer func() { _ = os.Remove(address) }()
	}

	grpcServer := grpc.NewServer()
	mcpdpluginsv1.RegisterPluginServer(grpcServer, &MyPlugin{})

	log.Printf("Plugin server listening on %s %s", network, address)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
```

## Import Path

The Go package name is `mcpdpluginsv1`, following Kubernetes-style versioned naming (e.g., `corev1`, `appsv1`):

```go
import "github.com/mozilla-ai/mcpd-plugins-sdk-go/pkg/plugins/v1"
```

## Proto Versioning

The SDK follows the versioning of [mcpd-proto](https://github.com/mozilla-ai/mcpd-proto):

- **API Version**: `plugins/v1/` (in proto repo) maps to `pkg/plugins/v1/` (in SDK)
- **Release Version**: Proto repo tags like `v0.0.1`, `v0.0.2`, etc.
- **SDK Version**: This repo's tags track SDK releases and may differ from proto versions

Current proto version: **v0.1.0**

## Repository Structure

```
mcpd-plugins-sdk-go/
├── README.md           # This file.
├── LICENSE             # Apache 2.0 license.
├── Makefile            # Proto fetching and code generation.
├── go.mod              # Go module definition.
├── go.sum              # Dependency checksums.
├── .gitignore          # Ignores tmp/ directory.
├── tmp/                # Downloaded protos (gitignored).
└── pkg/
    └── plugins/
        └── v1/
            ├── base.go            # BasePlugin helper.
            ├── constants.go       # Flow constant aliases.
            ├── server.go          # Serve() helper.
            ├── plugin.pb.go       # Generated protobuf types.
            └── plugin_grpc.pb.go  # Generated gRPC service.
```

## For SDK Maintainers

### Prerequisites

- Go 1.25.1 or later
- protoc (Protocol Buffer Compiler)
- protoc-gen-go and protoc-gen-go-grpc plugins

Install protoc plugins:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### Regenerating Code

The SDK uses a Makefile-based approach to fetch proto definitions and generate Go code:

**Fetch protos and generate:**
```bash
make all
go mod tidy
```

**Update to a new proto version:**
1. Edit `PROTO_VERSION` in the Makefile (e.g., `v0.0.2`)
2. Run `make clean all`
3. Run `go mod tidy`
4. Commit the updated generated files

**Run linter:**
```bash
make lint
```

**Clean generated files:**
```bash
make clean
```

Note: The `clean` target only removes generated `.pb.go` files, preserving helper files like `base.go` and `server.go`.

## Health Checking

The SDK follows [gRPC Health Checking Protocol](https://grpc.github.io/grpc/core/md_doc_health-checking.html) conventions:

```protobuf
rpc CheckHealth(google.protobuf.Empty) returns (google.protobuf.Empty);
rpc CheckReady(google.protobuf.Empty) returns (google.protobuf.Empty);
```

## License

Apache 2.0 - See LICENSE file for details.

## Contributing

This is an early PoC. Contribution guidelines coming soon.
