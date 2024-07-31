package main

import (
	"fmt"
	"net"
	"os/exec"
	"sync"
	"time"

	"github.com/songgao/water"
	"github.com/vishvananda/netlink"
)

type TapReceiver struct {
	tap  *water.Interface
	stat chan int
}

func NewTapReceiver() *TapReceiver {
	// Create Tap device
	iface, err := water.New(water.Config{
		DeviceType: water.TAP,
		PlatformSpecificParams: water.PlatformSpecificParams{
			Name:       "tap-test",
			MultiQueue: true,
		},
	})
	if err != nil {
		panic(err)
	}

	link, err := netlink.LinkByName(iface.Name())
	if err != nil {
		panic(err)
	}

	ipv4 := net.ParseIP(TapIP)

	// Set IP address
	netlink.AddrAdd(link, &netlink.Addr{
		IPNet: &net.IPNet{
			IP:   ipv4,
			Mask: net.CIDRMask(TapSubnet, 32),
		},
	})

	// Set UP
	netlink.LinkSetUp(link)

	time.Sleep(2000 * time.Millisecond)

	// Change neighbor table
	exec.Command("sudo", "ip", "neigh", "change", RemoteIP, "lladdr", RemoteMac, "dev", "tap-test").Run()

	return &TapReceiver{
		tap:  iface,
		stat: make(chan int, 1000),
	}
}

func (t *TapReceiver) Serve() {

	var wg sync.WaitGroup

	for i := 0; i < ReceiverRoutine; i++ {
		wg.Add(1)
		go t.HandleTraffic(&wg)
	}

	go t.Stat()

	wg.Wait()

}

func (t *TapReceiver) Stat() {
	var sum int
	tick := time.Tick(1 * time.Second)
	for {
		select {
		case <-tick:
			PrintByteRate(float64(sum))
			sum = 0
		case n := <-t.stat:
			sum += n
		}
	}
}

func (t *TapReceiver) HandleTraffic(wg *sync.WaitGroup) {
	defer wg.Done()

	frame := make([]byte, 5000)
	sum := 0
	tick := time.Tick(100 * time.Millisecond)

	for {
		select {
		case <-tick:
			t.stat <- sum
			sum = 0
		default:
			n, err := t.tap.Read([]byte(frame))
			if err != nil {
				panic(err)
			}
			sum += n
		}

	}
}

func PrintByteRate(byteRate float64) {
	bitRate := byteRate * 8
	var throughput string
	if bitRate < 1024 {
		throughput = fmt.Sprintf("%.2f bps", bitRate)
	} else if bitRate < 1024*1024 {
		throughput = fmt.Sprintf("%.2f Kbps", bitRate/1024)
	} else if bitRate < 1024*1024*1024 {
		throughput = fmt.Sprintf("%.2f Mbps", bitRate/1024/1024)
	} else {
		throughput = fmt.Sprintf("%.2f Gbps", bitRate/1024/1024/1024)
	}
	fmt.Printf("Throughput: %s\n", throughput)
}

func (t *TapReceiver) Close() {
	// Destroy Tap device
	t.tap.Close()
}
