package discovery

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/mdns"
)

// StartMDNS advertises the SyncDrop server on the local network.
func StartMDNS(port int, serverID string) (*mdns.Server, error) {
	host, err := os.Hostname()
	if err != nil {
		host = "syncdrop-host"
	}

	// Service name: _syncdrop._tcp
	// Domain: local.
	info := []string{
		fmt.Sprintf("version=1.0.0"),
		fmt.Sprintf("serverId=%s", serverID),
	}

	service, err := mdns.NewMDNSService(
		host,
		"_syncdrop._tcp",
		"",
		"",
		port,
		nil,
		info,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create mDNS service: %w", err)
	}

	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		return nil, fmt.Errorf("failed to start mDNS server: %w", err)
	}

	log.Printf("[mDNS] Broadcasting SyncDrop on local network (_syncdrop._tcp) on port %d", port)
	return server, nil
}
