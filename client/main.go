package main

import (
	"log"
	"os"
	"time"

	flags "github.com/jessevdk/go-flags"
	"github.com/linkernetworks/network-controller/client/common"
	pb "github.com/linkernetworks/network-controller/messages"
	"github.com/linkernetworks/network-controller/utils"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type podOptions struct {
	Name string `long:"podName" description:"The Pod Name, can set by environement variable" env:"POD_NAME" required:"true"`
	NS   string `long:"podNS" description:"The namespace of the Pod, can set by environement variable" env:"POD_NAMESPACE" required:"true"`
	UUID string `long:"podUUID" description:"The UUID of the Pod, can set by environement variable" env:"POD_UUID" required:"true"`
}

type interfaceOptions struct {
	IP string `short:"i" long:"ip" description:"The ip address of the interface, should be CIDR form"`
	//FIXME we will support the static route in the furture
	//Gateway string `short:"g" long:"gw" description:"The gateway of the inteface subnet"`
	VLANTag *int32 `short:"v" long:"vlan" description:"The Vlan Tag of the interface"`
}

type connectOptions struct {
	Bridge    string `short:"b" long:"bridge" description:"Target bridge name" required:"true"`
	Interface string `short:"n" long:"nic" description:"The interface name in the container" required:"true"`
}

type clientOptions struct {
	Server    string           `short:"s" long:"server " description:"target server address, [ip:port] for TCP or unix://[path] for UNIX" required:"true"`
	Connect   connectOptions   `group:"connectOptions"`
	Interface interfaceOptions `group:"interfaceOptions" `
	Pod       podOptions       `group:"podOptions" `
}

var options clientOptions
var parser = flags.NewParser(&options, flags.Default)

func main() {
	var setIP bool
	var setVLANAccessLink bool
	if _, err := parser.Parse(); err != nil {
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	// Verify IP address
	if options.Interface.IP != "" {
		setIP = true
	} else {
		log.Println("We don't have valid IP address from the arguments, we won't set the IP for", options.Connect.Interface)
	}

	if setIP {
		if !utils.IsValidCIDR(options.Interface.IP) {
			log.Fatalf("IP address is not correct: %s", options.Interface.IP)
		}
	}

	if options.Interface.VLANTag != nil {
		setVLANAccessLink = true
	}

	if setVLANAccessLink {
		if !utils.IsValidVLANTag(*options.Interface.VLANTag) {
			log.Fatalf("VLAN Tag is not correct: %d", *options.Interface.VLANTag)
		}
	}

	log.Println("Start to connect to ", options.Server)
	// Set up a connection to the server.
	conn, err := grpc.Dial(options.Server, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	ncClient := pb.NewNetworkControlClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	log.Println(options.Pod.Name, options.Pod.NS, options.Pod.UUID)
	// Find Network Namespace Path
	log.Println("Try to find the network namespace path")
	findNetworkNamespacePathResp, err := ncClient.FindNetworkNamespacePath(ctx,
		&pb.FindNetworkNamespacePathRequest{
			PodName:   options.Pod.Name,
			Namespace: options.Pod.NS,
			PodUUID:   options.Pod.UUID,
		},
	)
	if err != nil {
		log.Fatalf("There is something wrong with find network namespace pathpart.\n %v", err)
	}
	common.CheckFatal(
		findNetworkNamespacePathResp.Success,
		findNetworkNamespacePathResp.Reason,
		"Find network namesace path",
	)

	log.Printf(
		"The path is %s.",
		findNetworkNamespacePathResp.Path,
	)
	// Let's connect bridge
	log.Println(
		"Try to connect bridge",
		findNetworkNamespacePathResp.Path,
		options.Connect.Interface,
		options.Connect.Bridge,
	)
	connectBridgeResp, err := ncClient.ConnectBridge(ctx,
		&pb.ConnectBridgeRequest{
			Path:              findNetworkNamespacePathResp.Path,
			PodUUID:           options.Pod.UUID,
			ContainerVethName: options.Connect.Interface,
			BridgeName:        options.Connect.Bridge,
		},
	)
	if err != nil {
		log.Fatalf("There is something wrong with connect bridge: %v", err)
	}
	common.CheckFatal(
		connectBridgeResp.Success,
		connectBridgeResp.Reason,
		"Connect bridge",
	)

	if setIP {
		configureIfaceResp, err := ncClient.ConfigureIface(ctx,
			&pb.ConfigureIfaceRequest{
				Path:              findNetworkNamespacePathResp.Path,
				IP:                options.Interface.IP,
				ContainerVethName: options.Connect.Interface,
			},
		)
		if err != nil {
			log.Fatalf("There is something wrong with setting configure interface: %v", err)
		}
		common.CheckFatal(
			connectBridgeResp.Success,
			configureIfaceResp.Reason,
			"Configure interface",
		)
	}

	if setVLANAccessLink {
		setPortResp, err := ncClient.SetPort(ctx,
			&pb.SetPortRequest{
				IfaceName: utils.GenerateVethName(options.Pod.UUID, options.Connect.Interface),
				Options: &pb.PortOptions{
					Tag: *options.Interface.VLANTag,
				},
			},
		)
		if err != nil {
			log.Fatalf("There is something wrong with setting configure interface: %v", err)
		}
		common.CheckFatal(
			setPortResp.Success,
			setPortResp.Reason,
			"Set Port with VLAN",
		)
	}
	log.Printf("network-controller client has completed all tasks")
}
