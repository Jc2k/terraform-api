package command

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof" // Needed for debugging
	"strings"

	"github.com/xanzy/terraform-api/api"
	"github.com/xanzy/terraform-api/api/pb"
	"google.golang.org/grpc"
)

// APICommand is a special Command implementation that enables Terraform
// to run as a service providing a gRPC API endpoint.
type APICommand struct {
	Meta

	// When this channel is closed, the gRPC server will be stopped.
	ShutdownCh <-chan struct{}
}

// Run parses the command line arguments and runs Terraform as a service.
func (c *APICommand) Run(args []string) int {
	var ip string
	var port int

	cmdFlags := flag.NewFlagSet("api", flag.ContinueOnError)
	cmdFlags.StringVar(&ip, "ip", "127.0.0.1", "ip")
	cmdFlags.IntVar(&port, "port", 8080, "port")
	cmdFlags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := cmdFlags.Parse(args); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing flags: %v", err))
		return 1
	}

	s := api.NewServer(
		c.ContextOpts.Providers,
		c.ContextOpts.Provisioners,
	)

	// Create and register the gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterTerraformServer(grpcServer, s)

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Create a new listener
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating listener: %v", err))
		return 1
	}

	go func() {
		// Wait until we receive a shutdown command
		// and then stop the gRPC server
		<-c.ShutdownCh
		grpcServer.Stop()
	}()

	// Start serving incoming gRPC connections
	err = grpcServer.Serve(listener)
	if err != nil && err != grpc.ErrServerStopped {
		c.Ui.Error(fmt.Sprintf("Error serving incoming connections: %v", err))
		return 1
	}

	return 0
}

// Help displays the usage info.
func (c *APICommand) Help() string {
	helpText := `
Usage: terraform api [options]

  Run Terraform as a service providing a gRPC API endpoint.

Options:

  -ip=127.0.0.1           The IP address the service will bind to. Defaults
                          to 127.0.0.1.

  -port=8080              The port the service will bind to. Defaults to 8080.
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a description of the API command.
func (c *APICommand) Synopsis() string {
	return "Run Terraform as a service providing a gRPC API endpoint"
}
