package main

import (
	"fmt"
	"net"
	"sync"
)

type TapSender struct{}

func (t *TapSender) Serve() {
	var wg sync.WaitGroup

	for i := 0; i < SenderRoutine; i++ {
		wg.Add(1)
		go t.SendTraffic(&wg)
	}

	wg.Wait()

}

func (t *TapSender) SendTraffic(wg *sync.WaitGroup) {
	defer wg.Done()

	// fmt.Println("Sender started")
	destAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", "192.168.52.2", TapPort))
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := net.DialUDP("udp", nil, destAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	packet := make([]byte, PacketSize)

	for i := 0; i < PacketSize; i++ {
		packet[i] = byte(i)
	}

	for {
		// fmt.Println("Sent")
		_, err := conn.Write(packet)

		if err != nil {
			fmt.Println(err)
			return
		}

		// time.Sleep(1 * time.Microsecond)
	}
}
