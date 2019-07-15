package main

import (
	"fmt"
	"os"

	"github.com/deranjer/tinyMonitor-Agent/config"

	"github.com/rs/zerolog"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/sub"

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

func main() {
	clientSettings, Logger := config.SetupClient() //setup logging and all server settings
	Logger.Info().Msg("Server and Logger configuration complete")
	sock, err = sub.NewSocket()
	if err != nil {
		Logger.Fatal().Err(err).Msg("Can't get new subscribe socket")
	}
	err = sock.Dial(clientSettings.DialAddr)
	if err != nil {
		Logger.Fatal().Err(err).Str("Connection String", clientSettings.DialAddr).Msg("Failed to Dial Socket address")
	}
	Logger.Info().Str("Address", clientSettings.DialAddr).Msg("Dial to server opened with no errors")
	err = sock.SetOption(mangos.OptionSubscribe, []byte(""))
	if err != nil {
		die("cannot subscribe: %s", err.Error())
	}
	for {
		if msg, err = sock.Recv(); err != nil {
			die("Cannot recv: %s", err.Error())
		}
		fmt.Printf("CLIENT(%s): RECEIVED %s\n", "client1", string(msg))
	}

}
