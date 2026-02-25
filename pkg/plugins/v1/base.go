package mcpdpluginsv1

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

// BasePlugin provides sensible default implementations for all plugin methods.
// Plugin developers can embed this struct and override only the methods they need.
//
// Default behaviors:
//   - Configure: no-op
//   - Stop: no-op
//   - GetMetadata: returns empty metadata (should be overridden)
//   - GetCapabilities: returns no flows (should be overridden)
//   - CheckHealth: returns OK
//   - CheckReady: returns OK
//   - HandleRequest: passes through unchanged (continue=true)
//   - HandleResponse: passes through unchanged (continue=true)
//
// Usage:
//
//	import (
//	    "context"
//
//	    "github.com/mozilla-ai/mcpd-plugins-sdk-go/pkg/plugins/v1"
//	    "google.golang.org/protobuf/types/known/emptypb"
//	)
//
//	type MyPlugin struct {
//	    mcpdpluginsv1.BasePlugin
//	}
//
//	func (p *MyPlugin) GetMetadata(ctx context.Context, _ *emptypb.Empty) (*mcpdpluginsv1.Metadata, error) {
//	    return &mcpdpluginsv1.Metadata{
//	        Name: "my-plugin",
//	        Version: "1.0.0",
//	    }, nil
//	}
//
//	func (p *MyPlugin) HandleRequest(ctx context.Context, req *mcpdpluginsv1.HTTPRequest) (*mcpdpluginsv1.HTTPResponse, error) {
//	    // Custom logic here.
//	    return &mcpdpluginsv1.HTTPResponse{Continue: true}, nil
//	}
type BasePlugin struct {
	UnimplementedPluginServer
}

// Configure is a no-op by default.
func (b *BasePlugin) Configure(ctx context.Context, cfg *PluginConfig) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

// Stop is a no-op by default.
func (b *BasePlugin) Stop(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

// GetMetadata returns empty metadata by default. Plugins should override this.
func (b *BasePlugin) GetMetadata(ctx context.Context, _ *emptypb.Empty) (*Metadata, error) {
	return &Metadata{}, nil
}

// GetCapabilities returns no flows by default. Plugins should override this.
func (b *BasePlugin) GetCapabilities(ctx context.Context, _ *emptypb.Empty) (*Capabilities, error) {
	return &Capabilities{}, nil
}

// CheckHealth returns OK by default.
func (b *BasePlugin) CheckHealth(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

// CheckReady returns OK by default.
func (b *BasePlugin) CheckReady(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

// HandleRequest passes through the request unchanged with continue=true.
func (b *BasePlugin) HandleRequest(ctx context.Context, req *HTTPRequest) (*HTTPResponse, error) {
	return &HTTPResponse{
		Continue:   true,
		StatusCode: 0,
		Headers:    req.Headers,
		Body:       req.Body,
	}, nil
}

// HandleResponse passes through the response unchanged with continue=true.
func (b *BasePlugin) HandleResponse(ctx context.Context, resp *HTTPResponse) (*HTTPResponse, error) {
	return &HTTPResponse{
		Continue:   true,
		StatusCode: resp.StatusCode,
		Headers:    resp.Headers,
		Body:       resp.Body,
	}, nil
}
