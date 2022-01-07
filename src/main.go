package main

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

func main() {

	// mainmap := make(map[int]string)

	// log.Println(mainmap)
	//	service := discoverTcpService()
	//	log.Printf("discoverTcpService returned %s", service)
	valueMap := make(map[int]string)

	conn := connectTcp()
	clientInit(conn)
	// log.Printf("first read: %s\n\n", readMsg(conn))
	// msgText := readMsg(conn)
	// log.Printf("2nd read read before decode: %+v\n\n", msgText)
	// msg := decodeFocusriteMessage(msgText)
	// rootMesssageRouter(valueMap, msg)
	// msg2 := decodeFocusriteMessage(msgText)
	// rootMesssageRouter(valueMap, msg2)

	msg := decodeFocusriteMessage(readMsg(conn))
	rootMesssageRouter(valueMap, msg)

	msg2 := decodeFocusriteMessage(readMsg(conn))
	rootMesssageRouter(valueMap, msg2)

	// read2 := decodeFocusriteMessage(msg)
	// log.Printf("2nd read decoded: %+v\n\n", read2)
	// log.Printf("testing some fields clock source ID %+v", read2.DeviceArrival.Device.Clocking.ClockSource)

	//log.Printf("second read: %s", devicemsg)
	// maindevicearrival := decodeDeviceArrival(devicemsg)
	// //log.Printf("device arrival struct %+v\n\n", maindevicearrival)
	// // log.Printf("device arrival struct mixes %+v\n\n", maindevicearrival.Device.Mixer)

	// //	clockingLockedId, _ := strconv.ParseUint(maindevicearrival.Device.Clocking.Locked.ID, 10, 16)
	// //	log.Printf("int for locked id  %+v\n\n", clockingLockedId)
	// log.Printf("int for locked id  %+v\n\n", maindevicearrival.Device.Clocking.Locked)
	// log.Printf("type for locked id is %T\n\n", maindevicearrival.Device.Clocking.Locked.ID)
	// log.Printf("locked id int is %d\n\n", maindevicearrival.Device.Clocking.Locked.ID)

	// log.Printf("locked source ID int %d\n\n", maindevicearrival.Device.Clocking.ClockSource.ID)

	// log.Printf("source input spdif meter ID %d\n\n", maindevicearrival.Device.Inputs.SpdifRca[len(maindevicearrival.Device.Inputs.SpdifRca)-1].Meter.ID)

	// devicesettings := readMsg(conn)

	// maindeviceset := decodeDeviceSettings(devicesettings)
	// //	log.Printf("main device set struct item 33 %+v\n\n", maindeviceset.Item[33])

	// //	log.Printf("main device set struct item 33 id is %d\n\n", maindeviceset.Item[33].ID)

	// log.Printf("get control id is %d", findControlId(maindevicearrival.Device.Clocking.Locked.ID, maindeviceset))

	// log.Printf("get control value is %s", getControlValue(maindevicearrival.Device.Clocking.Locked.ID, maindeviceset))
	// log.Printf("get meter for spdif value is %s", getControlValue(maindevicearrival.Device.Inputs.SpdifRca[0].Meter.ID, maindeviceset))

	// log.Printf("get clocksource value is %s", getControlValue(maindevicearrival.Device.Clocking.ClockSource.ID, maindeviceset))

	// log.Printf("mainmap 441 %s", mainmap[441])
	// updateMap(mainmap, maindeviceset)
	// log.Printf("mainmap 441 %s", mainmap[441])
	// const Length=00002e <device-subscribe devid="1" subscribe="true"/>

	// func findControlId(controlID int, deviceSet DeviceSet ) int {

	// for i := range myconfig {
	// if myconfig[i].Key == "key1" {
	// Found!
	// }
	// }

	//intVar, err := strconv.Atoi(strVar)

	//todo find the
	//	log.Printf("is locked clock fuck xxx %+v\n\n", maindeviceset.Item[clockingLockedId])

	// 	log.Printf("is locked clock fuck xxx %+v\n\n", maindeviceset.Item[strconv.Atoi(maindevicearrival.Device.Clocking.Locked.ID)])
	// log.Printf("second read: %s", readMsg(conn))
	// log.Printf("third read: %s", readMsg(conn))

	// for i := range myconfig {
	// if myconfig[i].Key == "key1" {
	// Found!
	// }
	// }

}

func rootMesssageRouter(mainmap map[int]string, m FocusriteMessage) {
	log.Printf("root message routing %+v", m)
	//swtich
	log.Printf("root message routing %+v", m.DeviceArrival.XMLName.Local)
	// log.Printf("root message routing %+v", m)

	if m.DeviceArrival.XMLName.Local == "" {
		log.Printf("wasn't dev arrival")
	}

	// for i := range m.DeviceSet.Item {
	// 	mainmap[m.DeviceSet.Item[i].ID] = m.DeviceSet.Item[i].Value
	// 	log.Printf("updating map with device id %d, with value %s", m.DeviceSet.Item[i].ID, m.DeviceSet.Item[i].Value)
	// }
	return

}

func updateMap(mainmap map[int]string, m FocusriteMessage) {
	for i := range m.DeviceSet.Item {
		mainmap[m.DeviceSet.Item[i].ID] = m.DeviceSet.Item[i].Value
		log.Printf("updating map with device id %d, with value %s", m.DeviceSet.Item[i].ID, m.DeviceSet.Item[i].Value)
	}
	return

}

func getControlValue(controlID int, m FocusriteMessage) string {
	// return 123
	//	return 123
	for i := range m.DeviceSet.Item {
		if m.DeviceSet.Item[i].ID == controlID {
			// log.Printf("found a matching thing at index %d", i)
			return m.DeviceSet.Item[i].Value
		}
	}
	return "9999"
}

func findControlId(controlID int, m FocusriteMessage) int {
	// return 123
	//	return 123
	for i := range m.DeviceSet.Item {
		if m.DeviceSet.Item[i].ID == controlID {
			log.Printf("found a matching thing at index %d", i)
			return i
		}
	}
	return 9999
}

// for i := range myconfig {
//     if myconfig[i].Key == "key1" {
//         // Found!
//     }
// }

func decodeFocusriteMessage(payload string) FocusriteMessage {
	var m FocusriteMessage
	xml.Unmarshal([]byte(payload), &m)
	// sort.Slice(deviceset.Item, func(i, j int) bool {
	// })
	return m
}

// func decodeDeviceSettings(payload string) DeviceSet {
// 	var deviceset DeviceSet

// 	xml.Unmarshal([]byte(payload), &deviceset)
// 	// sort.Slice(deviceset.Item, func(i, j int) bool {
// 	// })
// 	return deviceset
// }

// //returns a device arrival struct
// func decodeDeviceArrival(payload string) DeviceArrival {
// 	//generate at  https://www.onlinetool.io/xmltogo/

// 	var devicearrival DeviceArrival

// 	xml.Unmarshal([]byte(payload), &devicearrival)
// 	return devicearrival

// }

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
	conn.Write([]byte("Length=000055 <client-details hostname=\"iphone\" client-key=\"11111111-2042-4050-8FED-A3CA5BABB11D\"/>"))
	// _, err := conn.Write([]byte("Length=000055 <client-details hostname=\"iphone\" client-key=\"11111111-2042-4050-8FED-A3CA5BABB11D\"/>"))

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

		log.Println("Length= was matched")
		blength, err := hex.DecodeString("00" + string(length)) // Length= is a 6 digit hex value, padding to 8 digit string so it can cast to an int32
		if err != nil {
			println("hex decode failed:", err.Error())
			os.Exit(1)
		}
		blength2 := binary.BigEndian.Uint32(blength)

		log.Printf("blength2 is is %d bytes", blength2)

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
