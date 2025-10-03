package v1

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

// Serve is a convenience function that handles all the boilerplate for running a plugin server.
// It parses command-line flags, sets up the appropriate network listener, creates a gRPC server,
// and serves the plugin implementation.
//
// Usage:
//
//	import (
//	    "log"
//
//	    pluginv1 "github.com/mozilla-ai/mcpd-plugins-sdk-go/pkg/plugins/v1/plugins"
//	)
//
//	func main() {
//	    if err := pluginv1.Serve(&MyPlugin{}); err != nil {
//	        log.Fatal(err)
//	    }
//	}
func Serve(impl PluginServer) error {
	var address, network string
	flag.StringVar(&address, "address", "", "gRPC address (socket path for unix, host:port for tcp)")
	flag.StringVar(&network, "network", "unix", "Network type (unix or tcp)")
	flag.Parse()

	if address == "" {
		return fmt.Errorf("--address flag is required")
	}

	lis, err := net.Listen(network, address)
	if err != nil {
		return fmt.Errorf("failed to listen on %s %s: %w", network, address, err)
	}

	// Clean up unix socket file when done.
	if network == "unix" {
		defer func() { _ = os.Remove(address) }()
	}

	grpcServer := grpc.NewServer()
	RegisterPluginServer(grpcServer, impl)

	// Handle graceful shutdown.
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		log.Println("Shutting down gracefully...")
		grpcServer.GracefulStop()
	}()

	log.Printf("Plugin server listening on %s %s", network, address)
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}
