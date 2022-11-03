package system

import (
	"net"
	"strings"

	"github.com/jaypipes/ghw"
)

func GetWirelessMacAddr() (result string, err error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, ifa := range ifas {
		if strings.HasPrefix(ifa.Name, "wl") {
			return ifa.HardwareAddr.String(), nil
		}
	}

	return "", nil
}

func GetFullInformation() (interface{}, error) {
	result := map[string]interface{}{}
	cpu, err := ghw.CPU()
	if err != nil {
		return nil, err
	}
	memory, err := ghw.Memory()
	if err != nil {
		return nil, err
	}
	storage, err := ghw.Block()
	if err != nil {
		return nil, err
	}
	topology, err := ghw.Topology()
	if err != nil {
		return nil, err
	}
	network, err := ghw.Network()
	if err != nil {
		return nil, err
	}
	pci, err := ghw.PCI()
	if err != nil {
		return nil, err
	}
	gpu, err := ghw.GPU()
	if err != nil {
		return nil, err
	}
	chassis, err := ghw.Chassis()
	if err != nil {
		return nil, err
	}
	bios, err := ghw.BIOS()
	if err != nil {
		return nil, err
	}
	baseboard, err := ghw.Baseboard()
	if err != nil {
		return nil, err
	}
	product, err := ghw.Product()
	if err != nil {
		return nil, err
	}

	result["cpu"] = cpu.JSONString(false)
	result["memory"] = memory.JSONString(false)
	result["storage"] = storage.JSONString(false)
	result["topology"] = topology.JSONString(false)
	result["network"] = network.JSONString(false)
	result["pci"] = pci.JSONString(false)
	result["gpu"] = gpu.JSONString(false)
	result["chassis"] = chassis.JSONString(false)
	result["bios"] = bios.JSONString(false)
	result["baseboard"] = baseboard.JSONString(false)
	result["product"] = product.JSONString(false)

	return result, nil
}
