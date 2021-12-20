package main

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"log"
	"net"
	"os"
	"reflect"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

func main() {
	service := discoverTcpService()
	log.Printf("discoverTcpService returned %s", service)
	connectTcp()

}

func discoverTcpService() string {
	const DiscoveryService = "localhost:30096"
	const DiscoveryTag = "server-announcement"

	type Response struct {
	}

	type Serveraccouncement struct {
		XMLName  xml.Name `xml:DiscoveryTag`
		Port     string   `xml:"port,attr"`
		Hostname string   `xml:"hostname,attr"`
	}

	RemoteAddr, err := net.ResolveUDPAddr("udp", DiscoveryService)
	p := bluemonday.UGCPolicy()
	p.AllowElements(DiscoveryTag)
	//LocalAddr := nil
	// see https://golang.org/pkg/net/#DialUDP

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

	//_, err = conn.Write(message)
	_, err = conn.Write([]byte(fullmessage))
	if err != nil {
		log.Println(err)
	}

	// receive message from server
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

func connectTcp() {
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

	_, err = conn.Write([]byte("Length=000055 <client-details hostname=\"iphone\" client-key=\"11111111-2042-4050-8FED-A3CA5BABB11D\"/>"))
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}

	reply := make([]byte, 7)
	_, err = conn.Read(reply)
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
		log.Println("Length= was matched")
		blength, err := hex.DecodeString("00" + string(length))
		if err != nil {
			println("hex decode failed:", err.Error())
			os.Exit(1)
		}
		//	var blength2 = new(int64) + blength
		blength2 := binary.BigEndian.Uint32(blength)

		log.Printf("length in string is %s", length)
		log.Printf("length in hex is %+v", length)
		log.Printf("blength in hex is %+v", blength)
		log.Printf("blength2 is is %d", blength2)
		log.Printf("type of blength %s\n", reflect.TypeOf(blength))

		// payload := make([]byte, blength[2]+1) //why do I have to add one, because it's a buffer size, and not an offset..  how to cast 3 []bytes into int32

		payload := make([]byte, blength2+1) //why do I have to add one, because it's a buffer size, and not an offset..  how to cast 3 []bytes into int32
		_, err = conn.Read(payload)
		if err != nil {
			println("Read length failed:", err.Error())
			os.Exit(1)
		}

		println("payload: ", string(payload))

		// println("length bytes", string(length))
		// println("length bytes2", string(blength))

		// reply2 := make([]byte, int(4096))
		// _, err = conn.Read(reply2)
		// if err != nil {
		// 	println("Read from server failed:", err.Error())
		// 	os.Exit(1)
		// }

		// println("reply from server length", string(reply2))

	}

	// reply2 := make([]byte, int(blength))
	// _, err = conn.Read(reply2)
	// if err != nil {
	// 	println("Read from server failed:", err.Error())
	// 	os.Exit(1)
	// }

	println("reply from server length", string(reply))

	//	println("reply from server=", string(reply2))

	conn.Close()

}
