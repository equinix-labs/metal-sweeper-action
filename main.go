package main

import (
	"fmt"
	"log"
	"os"

	"github.com/packethost/packngo"
)

var (
	version = "dev"

	warn = log.New(os.Stderr, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
)

const (
	tokenEnv       = "METAL_AUTH_TOKEN"
	projectEnv     = "METAL_PROJECT_ID"
	keepProjectEnv = "KEEP_PROJECT"

	consumerToken = "metal-sweeper-action"
	uaFmt         = "metal-sweeper-action/%s %s"
)

func main() {
	authToken := os.Getenv(tokenEnv)

	client := packngo.NewClientWithAuth(consumerToken, authToken, nil)
	client.UserAgent = fmt.Sprintf(uaFmt, version, client.UserAgent)

	projectID := os.Getenv(projectEnv)

	// What to delete?
	// - Devices
	// - ProjectVirtualNetworks
	// - Project
	// - TODO(displague) SpotMarkets
	// - TODO(displague) VolumeAttachments
	// - TODO(displague) Volumes (are these project specific?)

	devices, _, err := client.Devices.List(projectID, nil)

	if err != nil {
		warn.Println("Could not list devices", err)
	}

	for _, device := range devices {
		fmt.Println("Deleting device", device.Hostname)
		if _, err := client.Devices.Delete(device.ID, true); err != nil {
			warn.Println("Could not delete device", err)
		}
	}

	vlans, _, err := client.ProjectVirtualNetworks.List(projectID, nil)

	if err != nil {
		warn.Println("Could not list vlans", err)
	}

	for _, vlan := range vlans.VirtualNetworks {
		fmt.Println("Deleting vlan", vlan.Description)
		if _, err := client.ProjectVirtualNetworks.Delete(vlan.ID); err != nil {
			warn.Println("Could not delete vlan", err)
		}
	}

	if os.Getenv(keepProjectEnv) == "false" {
		_, err := client.Projects.Delete(projectID)
		if err != nil {
			warn.Println("Could not delete project", err)
		}
	}
}
