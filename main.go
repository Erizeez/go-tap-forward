package main

const (
	TapIP           = "192.168.52.1"
	RemoteIP        = "192.168.52.2"
	RemoteMac       = "00:11:22:33:44:55"
	TapSubnet       = 24
	TapPort         = 8080
	PacketSize      = 1400
	ReceiverRoutine = 4
	SenderRoutine   = 4
)

func main() {
	tapReceiver := NewTapReceiver()
	defer tapReceiver.Close()

	go tapReceiver.Serve()

	tapSender := &TapSender{}
	go tapSender.Serve()

	select {}
}
