package main

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/xml"
	"flag"
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/microcosm-cc/bluemonday"
)

func main() {
	log.SetLevel(log.InfoLevel) //TODO change to info

	udphost := flag.String("h", "localhost", "hostname for initial UDP discovery")
	debugPtr := flag.Bool("d", false, "debug enable")
	tracePtr := flag.Bool("t", false, "trace enable")
	flag.Parse()

	if *debugPtr {
		log.SetLevel(log.DebugLevel)
	}

	if *tracePtr {
		log.SetLevel(log.TraceLevel)
	}

	log.Debug("hostname set to ", *udphost)
	valueMap := make(map[int]string)
	var deviceArrivalMsg FocusriteMessage //save the arrival struct to this global!@!

	conn := connectTcp(discoverTcpService(*udphost))
	clientInit(conn)

	go bgKeepAlive(conn) //send keep alives in background, can't image conn is thread safe..

	go bgWatchClock(conn, valueMap, &deviceArrivalMsg) //watch the clock and make sure it's what we want..

	for {
		rootMesssageRouter(conn, valueMap, &deviceArrivalMsg, decodeFocusriteMessage(readMsg(conn)))
		time.Sleep(5 * time.Millisecond)
	}
} //end main

func bgWatchClock(conn *net.TCPConn, valueMap map[int]string, deviceArrivalMsg *FocusriteMessage) {
	// log.Println("running watchclock")
	watchTicker := time.NewTicker(5 * time.Second) //need to be slightly longer than it takes to lock..
	for t := range watchTicker.C {
		log.Debug("running watchclock", t)
		log.Debug("source input spdif meter ID: ", valueMap[deviceArrivalMsg.DeviceArrival.Device.Inputs.SpdifRca[0].Meter.ID])
		log.Debug("clock locked: ", valueMap[deviceArrivalMsg.DeviceArrival.Device.Clocking.Locked.ID])
		log.Debug("clock source: ", valueMap[deviceArrivalMsg.DeviceArrival.Device.Clocking.ClockSource.ID])
		if valueMap[deviceArrivalMsg.DeviceArrival.Device.Clocking.Locked.ID] != "true" &&
			valueMap[deviceArrivalMsg.DeviceArrival.Device.Clocking.ClockSource.ID] == "S/PDIF" {
			log.Info("Setting clock to Internal")
			setInternal := fmt.Sprintf("<set devid=\"1\"><item id=\"%d\" value=\"Internal\"/></set>", deviceArrivalMsg.DeviceArrival.Device.Clocking.ClockSource.ID)
			writeMsg(conn, setInternal)
		}
		spdifLevel, _ := strconv.ParseInt(valueMap[deviceArrivalMsg.DeviceArrival.Device.Inputs.SpdifRca[0].Meter.ID], 10, 8)
		log.Debug("spdiflevel:", spdifLevel)
		if valueMap[deviceArrivalMsg.DeviceArrival.Device.Clocking.ClockSource.ID] == "Internal" &&
			spdifLevel > -100 {
			log.Info("Setting clock to S/PDIF")
			setSpdif := fmt.Sprintf("<set devid=\"1\"><item id=\"%d\" value=\"S/PDIF\"/></set>", deviceArrivalMsg.DeviceArrival.Device.Clocking.ClockSource.ID)
			writeMsg(conn, setSpdif)
		}
	}
}

func rootMesssageRouter(conn *net.TCPConn, valueMap map[int]string, deviceArrivalMsg *FocusriteMessage, m FocusriteMessage) {
	// a array of the tags we will route
	focusriteStructs := []string{m.ClientDetails.XMLName.Local, m.DeviceArrival.XMLName.Local,
		m.DeviceSet.XMLName.Local, m.Approval.XMLName.Local, m.KeepAlive.XMLName.Local}
	sort.Strings(focusriteStructs)
	handler := focusriteStructs[len(focusriteStructs)-1]
	log.Trace("target handler: ", handler)

	switch handler {
	case "client-details":
		log.Debug("get the client deets")
	case "device-arrival":
		log.Debug("processing arrival")
		*deviceArrivalMsg = m
		log.Debug("sending subscribe message after receiving device arrival")
		//conn.Write([]byte(`Length=00002e <device-subscribe devid="1" subscribe="true"/>`))
		writeMsg(conn, `<device-subscribe devid="1" subscribe="true"/>`)
	case "keep-alive":
		log.Debug("received the keepalive")
		return
	case "set":
		log.Trace("updating the value map")
		for i := range m.DeviceSet.Item {
			valueMap[m.DeviceSet.Item[i].ID] = m.DeviceSet.Item[i].Value
			// log.Printf("updating map with device id %d, with value %s", m.DeviceSet.Item[i].ID, m.DeviceSet.Item[i].Value)
		}

	case "approval":
		log.Debug("check approval")
		log.Debug("approval message: ", m.Approval)
		if m.Approval.Authorised == "false" {
			log.Error("Please authorize in Focusrite app")
		}

	default:
		log.Warn("unknown handler: ", handler)
	}

	return

}

func decodeFocusriteMessage(payload string) FocusriteMessage {
	var m FocusriteMessage
	xml.Unmarshal([]byte(payload), &m)
	return m
}

func discoverTcpService(host string) string {
	DiscoveryService := host + ":30096"
	const DiscoveryTag = "server-announcement"

	type Serveraccouncement struct {
		XMLName  xml.Name `xml:DiscoveryTag`
		Port     string   `xml:"port,attr"`
		Hostname string   `xml:"hostname,attr"`
	}

	RemoteAddr, err := net.ResolveUDPAddr("udp", DiscoveryService)
	if err != nil {
		log.Fatal(err)
	}
	p := bluemonday.UGCPolicy()
	p.AllowElements(DiscoveryTag)

	conn, err := net.DialUDP("udp", nil, RemoteAddr)

	if err != nil {
		log.Fatal(err)
	}

	log.Info("Sending UDP discovery to: ", DiscoveryService)
	log.Debug("Remote UDP address: ", conn.RemoteAddr().String())
	log.Debug("Local UDP client address: ", conn.LocalAddr().String())

	defer conn.Close()

	command := `<client-discovery app="SAFFIRE-CONTROL" version="4" device="iOS"/>`
	log.Trace("Length is: ", len(command))
	fullmessage := fmt.Sprintf("Length=%06x %s", len(command), command)
	log.Trace("before sanitize")
	log.Trace(p.Sanitize(fullmessage))
	log.Trace("after sanitize")
	log.Trace(fullmessage)

	_, err = conn.Write([]byte(fullmessage))
	if err != nil {
		log.Error("failed to send UDP message, not sure how this would ever fail, UPD is connectionless", err)
		os.Exit(1)
	}

	// receive udp message from server
	buffer := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	n, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Error("failed to read UDP message, is focusrite control running on the discovery host?")
		log.Error("Error message: ", err)
		os.Exit(1)
	}

	log.Info("UDP Server : ", addr)
	log.Debug("Received from UDP server : ", string(buffer[:n]))
	xmlmessage := strings.SplitN(string(buffer[:n]), " ", 2)[1]
	log.Debug("udp xml message: ", xmlmessage)

	var disco Serveraccouncement

	xml.Unmarshal([]byte(xmlmessage), &disco)
	log.Debug("port: ", disco.Port)
	log.Debug("hostname: ", disco.Hostname)
	return disco.Hostname + ":" + disco.Port
}

func connectTcp(serviceName string) *net.TCPConn {

	tcpAddr, err := net.ResolveTCPAddr("tcp", serviceName)
	if err != nil {
		log.Error("resolution failed: ", err.Error())
		os.Exit(1)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Error("Dial failed:", err.Error())
		os.Exit(1)
	}
	return conn
}

func clientInit(conn *net.TCPConn) {
	myHostname, _ := os.Hostname()
	log.Debug("hostname: ", myHostname)
	clientMsg := fmt.Sprintf(`<client-details hostname="autoclock-%s" client-key="11111111-2042-4050-8FED-B3CA5BABB11D"/>`, myHostname)
	writeMsg(conn, clientMsg)
}

func bgKeepAlive(conn *net.TCPConn) {
	keepAliveTicker := time.NewTicker(3 * time.Second)
	for t := range keepAliveTicker.C {
		log.Debug("sending keep alive", t)
		// sendKeepAlive(conn)
		writeMsg(conn, `<keep-alive/>`)
		//conn.Write([]byte(`Length=00000d <keep-alive/>`))
	}
}

func writeMsg(conn *net.TCPConn, msg string) {
	fullmsg := fmt.Sprintf("Length=%06x %s", len(msg), msg)
	log.Debug("fullmessage: ", fullmsg)
	conn.Write([]byte(fullmsg))
}

func readMsg(conn *net.TCPConn) string {
	reply := make([]byte, 7)
	var response string
	_, err := conn.Read(reply)
	if err != nil {
		log.Error("Read from server failed:", err.Error())
		os.Exit(1)
	}
	if string(reply[:]) == "Length=" {
		length := make([]byte, 6)
		_, err = conn.Read(length)
		if err != nil {
			log.Error("Read length from server failed:", err.Error())
			os.Exit(1)
		}

		space := make([]byte, 1)
		_, err = conn.Read(space)
		if err != nil {
			log.Error("Read space from server failed:", err.Error())
			os.Exit(1)
		}

		blength, err := hex.DecodeString("00" + string(length)) // Length= is a 6 digit hex value, padding to 8 digit string so it can cast to an int32
		if err != nil {
			log.Error("hex length decode failed:", err.Error())
			os.Exit(1)
		}

		blength2 := binary.BigEndian.Uint32(blength)
		log.Trace("blength2 bytes: ", blength2)

		payload := make([]byte, blength2)
		_, err = conn.Read(payload)
		if err != nil {
			log.Error("Read length failed:", err.Error())
			os.Exit(1)
		}

		response = fmt.Sprintf("<focusritemessage>%s</focusritemessage>", string(payload))

	}
	return response
}
