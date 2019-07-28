package main

import (
	"fmt"
	"os"
	"time"

	"github.com/deranjer/tinyMonitor-Agent/config"
	"github.com/deranjer/tinyMonitor/messaging"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/net"

	"github.com/rs/zerolog"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/pair"

	// register transports
	_ "nanomsg.org/go/mangos/v2/transport/tcp" //TODO change to /all to register all transports
)

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

var (
	err  error
	sock mangos.Socket
	msg  []byte
	//Logger is the global Logger variable
	Logger zerolog.Logger
)

func guessDefaultConnection(interfaces []net.InterfaceStat) net.InterfaceAddr { //TODO actually fix this
	var netGuess net.InterfaceAddr
	for _, netInterface := range interfaces {
		if netInterface.Name == "lo" {
			continue
		} else if netInterface.Name == "ens18" {
			netGuess = netInterface.Addrs[0]
		}

	}
	return netGuess
}

func connectServer(clientSettings config.ClientConfig, Logger zerolog.Logger, serverSendCh chan string) {
	systemInfo, err := host.Info()
	if err != nil {
		Logger.Fatal().Err(err).Msg("Unable to get information from host!")
	}
	ipAddr, err := net.Interfaces()
	if err != nil {
		Logger.Fatal().Err(err).Msg("Unable to get network information from host!")
	}
	ipAddrGuess := guessDefaultConnection(ipAddr)
	sock, err = pair.NewSocket()
	if err != nil {
		Logger.Fatal().Err(err).Msg("Can't setup new pair socket")
	}
	err = sock.Dial(clientSettings.DialAddr)
	if err != nil {
		Logger.Fatal().Err(err).Str("Connection String", clientSettings.DialAddr).Msg("Failed to Dial Socket address of server")
	}
	Logger.Info().Str("Address", clientSettings.DialAddr).Msg("Dial to server opened with no errors")

	msgStruct := messaging.BaseMessage{
		MessageType: "RegisterAgent",
		MessageBody: messaging.RegisterAgent{
			AgentHostName: systemInfo.Hostname,
			AgentIPAddr:   ipAddrGuess.Addr,
			AgentJoinDate: time.Now(),
		},
	}
	b := messaging.MessageEncode(msgStruct)
	err = sock.Send(b)
	if err != nil {
		Logger.Fatal().Err(err).Str("Connection String", clientSettings.DialAddr).Msg("Failed to send message to server")
	}
	Logger.Info().Str("Address", clientSettings.DialAddr).Msg("Message Sent")
	msg, err = sock.Recv()
	if err != nil {
		Logger.Error().Err(err).Msg("Agent failed to receive pair message from Server")
	} else {
		Logger.Debug().Str("Message Body", string(msg)).Msg("Message Received from Server")
		return
	}

}

func main() {
	clientSettings, Logger := config.SetupClient() //setup logging and all server settings
	Logger.Info().Msg("Server and Logger configuration complete")
	serverSendCh := make(chan string)
	connectServer(clientSettings, Logger, serverSendCh)

}
