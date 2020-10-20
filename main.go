package main

import (
	"fmt"
	"os"

	"github.com/packethost/packngo"
)

func main() {
	// PACKET_AUTH_TOKEN should be already set
	// the client will use it by default (?)
	client, err := packngo.NewClient()
	if err != nil {
		panic(err)
	}

	projectID := os.Getenv("PROJECT_ID")
	// TODO(displague) lookup project by name
	// , _, err := c.Projects.List(nil)

	// What to delete?
	// - Devices
	// - TODO(displague) SpotMarkets
	// - TODO(displague) VolumeAttachments
	// - TODO(displague) Volumes (are these project specific?)
	// - TODO(displague) ProjectVirtualNetworks

	devices, _, err := client.Devices.List(projectID, nil)

	if err != nil {
		panic(err)
	}

	for _, device := range devices {
		fmt.Println("Deleting", device.Hostname)
		if _, err := client.Devices.Delete(device.ID, true); err != nil {
			panic(err)
		}
	}

}
