package main

import (
	"context"
	"fmt"
	"log"
	"os"

	metal "github.com/equinix-labs/metal-go/metal/v1"
)

var (
	version = "dev"
	warn    = log.New(os.Stderr, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	client  *metal.APIClient
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

	config := metal.NewConfiguration()
	config.AddDefaultHeader("X-Auth-Token", authToken)
	config.UserAgent = fmt.Sprintf(uaFmt, version, config.UserAgent)
	client = metal.NewAPIClient(config)

	projectID := os.Getenv(projectEnv)

	// What to delete?
	// - Devices
	// - ProjectVirtualNetworks
	// - Project
	// - TODO(displague) SpotMarkets
	// - TODO(displague) VolumeAttachments
	// - TODO(displague) Volumes (are these project specific?)

	devices, err := getAllProjectDevices(projectID)

	if err != nil {
		warn.Println("Could not list devices", err)
	}

	for _, device := range devices {
		fmt.Println("Deleting device", device.Hostname)
		if _, err := client.DevicesApi.DeleteDevice(context.Background(), device.GetId()).Execute(); err != nil {
			warn.Println("Could not delete device", err)
		}
	}

	vlans, _, err := client.VLANsApi.FindVirtualNetworks(context.Background(), projectID).Execute()

	if err != nil {
		warn.Println("Could not list vlans", err)
	}

	for _, vlan := range vlans.VirtualNetworks {
		fmt.Println("Deleting vlan", vlan.Description)
		if _, _, err := client.VLANsApi.DeleteVirtualNetwork(context.Background(), vlan.GetId()).Execute(); err != nil {
			warn.Println("Could not delete vlan", err)
		}
	}

	if os.Getenv(keepProjectEnv) == "false" {
		_, err := client.ProjectsApi.DeleteProject(context.Background(), projectID).Execute()
		if err != nil {
			warn.Println("Could not delete project", err)
		}
	}
}

func getAllProjectDevices(projectID string) ([]metal.Device, error) {
	var devices []metal.Device
	var page int32 = 1

	for {
		devicePage, _, err := client.DevicesApi.FindProjectDevices(context.Background(), projectID).Page(page).Execute()

		if err != nil {
			return nil, err
		}

		devices = append(devices, devicePage.Devices...)
		if devicePage.Meta.GetLastPage() > devicePage.Meta.GetCurrentPage() {
			page = page + 1
			continue
		}

		return devices, nil
	}
}
