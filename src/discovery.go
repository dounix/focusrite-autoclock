package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	//"golang.org/x/net/html"
)

type Serveraccouncement struct {
	XMLName  xml.Name `xml:"server-announcement"`
	Port     string   `xml:"port,attr"`
	Hostname string   `xml:"hostname,attr"`

	//	Serverannouncement string   `xml:"server-announcement,attr"`
}

func main() {
	hostName := "localhost"
	portNum := "30096"

	service := hostName + ":" + portNum

	RemoteAddr, err := net.ResolveUDPAddr("udp", service)
	p := bluemonday.UGCPolicy()
	//p := bluemonday.NewPolicy()

	p.AllowElements("server-announcement")
	//LocalAddr := nil
	// see https://golang.org/pkg/net/#DialUDP

	conn, err := net.DialUDP("udp", nil, RemoteAddr)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Established connection to %s \n", service)
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

	//xmlbuffer := make([]byte, 42)
	//nxml, xmladdr, err := conn.ReadFromUDP(xmlbuffer)

	//fmt.Println("UDP Server : ", xmladdr)
	//fmt.Println("XML Received from UDP server : ", string(xmlbuffer[:nxml]))
}
