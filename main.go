package main

import (
	"fmt"
	"os"

	"github.com/deranjer/tinyMonitor-Agent/config"

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

func connectServer(clientSettings config.ClientConfig, Logger zerolog.Logger, serverSendCh chan string) {
	sock, err = pair.NewSocket()
	if err != nil {
		Logger.Fatal().Err(err).Msg("Can't setup new pair socket")
	}
	err = sock.Dial(clientSettings.DialAddr)
	if err != nil {
		Logger.Fatal().Err(err).Str("Connection String", clientSettings.DialAddr).Msg("Failed to Dial Socket address of server")
	}
	Logger.Info().Str("Address", clientSettings.DialAddr).Msg("Dial to server opened with no errors")

	err = sock.Send([]byte("Client Connecting to server"))
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
