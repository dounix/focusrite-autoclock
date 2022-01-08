package main

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

func main() {

	valueMap := make(map[int]string)
	var deviceArrivalMsg FocusriteMessage //save the arrival struct to this global!@!

	conn := connectTcp()
	clientInit(conn)

	go bgKeepAlive(conn) //send keep alives in background, can't image conn is thread safe..

	go bgWatchClock(conn, valueMap, &deviceArrivalMsg) //watch the clock and make sure it's what we want..

	for {
		rootMesssageRouter(conn, valueMap, &deviceArrivalMsg, decodeFocusriteMessage(readMsg(conn)))
		//  These aren't avaialble until we decode the device arrival(field IDs) and device set(values)
		//  log.Printf("source input spdif meter ID %+v\n", valueMap[deviceArrivalMsg.DeviceArrival.Device.Inputs.SpdifRca[0].Meter.ID])
		//	log.Printf("clock locked %+v\n", valueMap[deviceArrivalMsg.DeviceArrival.Device.Clocking.Locked.ID])
		//	log.Printf("clock source %+v\n", valueMap[deviceArrivalMsg.DeviceArrival.Device.Clocking.ClockSource.ID])
		time.Sleep(5 * time.Millisecond)
	}
} //end main

func bgWatchClock(conn *net.TCPConn, valueMap map[int]string, deviceArrivalMsg *FocusriteMessage) {
	// log.Println("running watchclock")
	watchTicker := time.NewTicker(5 * time.Second) //need to be slightly longer than it takes to lock..
	for t := range watchTicker.C {
		log.Println("running watchclock", t)
		log.Printf("source input spdif meter ID %+v\n", valueMap[deviceArrivalMsg.DeviceArrival.Device.Inputs.SpdifRca[0].Meter.ID])
		log.Printf("clock locked %+v\n", valueMap[deviceArrivalMsg.DeviceArrival.Device.Clocking.Locked.ID])
		log.Printf("clock source %+v\n", valueMap[deviceArrivalMsg.DeviceArrival.Device.Clocking.ClockSource.ID])
		if valueMap[deviceArrivalMsg.DeviceArrival.Device.Clocking.Locked.ID] == "false" &&
			valueMap[deviceArrivalMsg.DeviceArrival.Device.Clocking.ClockSource.ID] == "S/PDIF" {
			log.Println("Setting clock to Internal")
			setInternal := fmt.Sprintf("<set devid=\"1\"><item id=\"%d\" value=\"Internal\"/></set>", deviceArrivalMsg.DeviceArrival.Device.Clocking.ClockSource.ID)
			writeMsg(conn, setInternal)
		}
		spdifLevel, _ := strconv.ParseInt(valueMap[deviceArrivalMsg.DeviceArrival.Device.Inputs.SpdifRca[0].Meter.ID], 10, 8)
		//log.Println("spdiflevel:", spdifLevel)
		if valueMap[deviceArrivalMsg.DeviceArrival.Device.Clocking.ClockSource.ID] == "Internal" &&
			spdifLevel > -100 {
			log.Println("Setting clock to S/PDIF")
			setSpdif := fmt.Sprintf("<set devid=\"1\"><item id=\"%d\" value=\"S/PDIF\"/></set>", deviceArrivalMsg.DeviceArrival.Device.Clocking.ClockSource.ID)
			writeMsg(conn, setSpdif)
		}
		// sendKeepAlive(conn)
		// conn.Write([]byte(`Length=00000d <keep-alive/>`))
	}
}

func rootMesssageRouter(conn *net.TCPConn, valueMap map[int]string, deviceArrivalMsg *FocusriteMessage, m FocusriteMessage) {
	// a array of the tags we will route
	focusriteStructs := []string{m.ClientDetails.XMLName.Local, m.DeviceArrival.XMLName.Local,
		m.DeviceSet.XMLName.Local, m.Approval.XMLName.Local, m.KeepAlive.XMLName.Local}
	sort.Strings(focusriteStructs)
	handler := focusriteStructs[len(focusriteStructs)-1]
	// log.Printf("target handler: %+v", handler)
	//log.Printf("root message routing %+v", m)
	//swtich
	// log.Printf("root message routing %+v", m.DeviceArrival.XMLName.Local)
	// log.Printf("root message routing %+v", m)
	switch handler {
	case "client-details":
		log.Println("get the client deets")
	case "device-arrival":
		log.Println("arrival")
		*deviceArrivalMsg = m
		log.Println("sending subscribe message after receiving device arrival")
		//conn.Write([]byte(`Length=00002e <device-subscribe devid="1" subscribe="true"/>`))
		writeMsg(conn, `<device-subscribe devid="1" subscribe="true"/>`)
	case "keep-alive":
		log.Println("received the keepalive")
		return
	case "set":
		// log.Println("updating the value map")
		for i := range m.DeviceSet.Item {
			valueMap[m.DeviceSet.Item[i].ID] = m.DeviceSet.Item[i].Value
			// log.Printf("updating map with device id %d, with value %s", m.DeviceSet.Item[i].ID, m.DeviceSet.Item[i].Value)
		}

	case "approval":
		log.Println("check approval")
		log.Printf("approval message: %+v", m.Approval)

	default:
		log.Printf("unknown handler %s", handler)
	}

	return

}

func decodeFocusriteMessage(payload string) FocusriteMessage {
	var m FocusriteMessage
	xml.Unmarshal([]byte(payload), &m)
	// sort.Slice(deviceset.Item, func(i, j int) bool {
	// })
	return m
}

func discoverTcpService() string {
	const DiscoveryService = "localhost:30096"
	const DiscoveryTag = "server-announcement"

	type Serveraccouncement struct {
		XMLName  xml.Name `xml:DiscoveryTag`
		Port     string   `xml:"port,attr"`
		Hostname string   `xml:"hostname,attr"`
	}

	RemoteAddr, err := net.ResolveUDPAddr("udp", DiscoveryService)
	p := bluemonday.UGCPolicy()
	p.AllowElements(DiscoveryTag)

	conn, err := net.DialUDP("udp", nil, RemoteAddr)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Established connection to %s \n", DiscoveryService)
	log.Printf("Remote UDP address : %s \n", conn.RemoteAddr().String())
	log.Printf("Local UDP client address : %s \n", conn.LocalAddr().String())

	defer conn.Close()

	command := `<client-discovery app="SAFFIRE-CONTROL" version="4" device="iOS"/>`
	log.Printf("Length is %x\n", len(command))
	fullmessage := fmt.Sprintf("Length=%06x %s", len(command), command)
	log.Println("before sanitize")
	log.Println(p.Sanitize(fullmessage))
	log.Println("after sanitize")
	log.Println(fullmessage)

	_, err = conn.Write([]byte(fullmessage))
	if err != nil {
		log.Println(err)
	}

	// receive udp message from server
	buffer := make([]byte, 1024)
	n, addr, err := conn.ReadFromUDP(buffer)

	fmt.Println("UDP Server : ", addr)
	fmt.Println("Received from UDP server : ", string(buffer[:n]))
	//	xmlmessage := strings.SplitN(string(buffer[:n]), " ", 1)
	xmlmessage := strings.SplitN(string(buffer[:n]), " ", 2)[1]
	log.Printf("xml message : %s ", xmlmessage)
	//tokenizer := html.NewTokenizer(xmlmessage)

	var disco Serveraccouncement

	xml.Unmarshal([]byte(xmlmessage), &disco)
	//	fmt.Println(disco)
	//	fmt.Printf("%V\n", disco)
	fmt.Printf("port %s\n", disco.Port)
	fmt.Printf("hostname %s\n", disco.Hostname)
	return disco.Hostname + ":" + disco.Port
}

func connectTcp() *net.TCPConn {
	serviceName := discoverTcpService()
	tcpAddr, err := net.ResolveTCPAddr("tcp", serviceName)
	if err != nil {
		log.Println("resolution failed: ", err.Error())
		os.Exit(1)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Println("Dial failed:", err.Error())
		os.Exit(1)
	}
	return conn
}

func clientInit(conn *net.TCPConn) {
	// conn.Write([]byte("Length=000055 <client-details hostname=\"jimsne\" client-key=\"11111111-2042-4050-8FED-B3CA5BABB11D\"/>"))
	writeMsg(conn, `<client-details hostname="jimsne" client-key="11111111-2042-4050-8FED-B3CA5BABB11D"/>`)

	// _, err := conn.Write([]byte("Length=000055 <client-details hostname=\"iphone\" client-key=\"11111111-2042-4050-8FED-A3CA5BABB11D\"/>"))

}

func bgKeepAlive(conn *net.TCPConn) {
	keepAliveTicker := time.NewTicker(3 * time.Second)
	for t := range keepAliveTicker.C {
		log.Println("sending keep alive", t)
		// sendKeepAlive(conn)
		writeMsg(conn, `<keep-alive/>`)
		//conn.Write([]byte(`Length=00000d <keep-alive/>`))
	}
}

func writeMsg(conn *net.TCPConn, msg string) {
	fullmsg := fmt.Sprintf("Length=%06x %s", len(msg), msg)
	log.Printf("fullmessage: %s", fullmsg)
	conn.Write([]byte(fullmsg))
}

func readMsg(conn *net.TCPConn) string {
	reply := make([]byte, 7)
	var response string
	_, err := conn.Read(reply)
	if err != nil {
		println("Read from server failed:", err.Error())
		os.Exit(1)
	}
	if string(reply[:]) == "Length=" {
		length := make([]byte, 6)
		_, err = conn.Read(length)
		if err != nil {
			println("Read length from server failed:", err.Error())
			os.Exit(1)
		}

		space := make([]byte, 1)
		_, err = conn.Read(space)
		if err != nil {
			println("Read space from server failed:", err.Error())
			os.Exit(1)
		}
		// log.Printf("space is: x%sx", space)

		//log.Println("Length= was matched")
		blength, err := hex.DecodeString("00" + string(length)) // Length= is a 6 digit hex value, padding to 8 digit string so it can cast to an int32
		if err != nil {
			println("hex length decode failed:", err.Error())
			os.Exit(1)
		}
		blength2 := binary.BigEndian.Uint32(blength)

		// log.Printf("blength2 is is %d bytes", blength2)

		payload := make([]byte, blength2)
		_, err = conn.Read(payload)
		if err != nil {
			println("Read length failed:", err.Error())
			os.Exit(1)
		}

		response = fmt.Sprintf("<focusritemessage>%s</focusritemessage>", string(payload))

	}
	return response
}
