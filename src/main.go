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
	//	service := discoverTcpService()
	//	log.Printf("discoverTcpService returned %s", service)
	conn := connectTcp()
	clientInit(conn)
	log.Printf("first read: %s", readMsg(conn))
	devicemsg := readMsg(conn)
	//log.Printf("second read: %s", devicemsg)
	decodeDeviceArrival(devicemsg)

	devicesettings := readMsg(conn)

	decodeDeviceSettings(devicesettings)
	// log.Printf("second read: %s", readMsg(conn))
	// log.Printf("third read: %s", readMsg(conn))

	// for i := range myconfig {
	// if myconfig[i].Key == "key1" {
	// Found!
	// }
	// }

}

func decodeDeviceSettings(payload string) string {
	//generate at  https://www.onlinetool.io/xmltogo/
	//log.Printf("decodeDeviceArrival payload id %s", payload)
	type DeviceSet struct {
		XMLName xml.Name `xml:"set"`
		Text    string   `xml:",chardata"`
		Devid   string   `xml:"devid,attr"`
		Item    []struct {
			Text  string `xml:",chardata"`
			ID    string `xml:"id,attr"`
			Value string `xml:"value,attr"`
		} `xml:"item"`
	}

	var deviceset DeviceSet

	xml.Unmarshal([]byte(payload), &deviceset)
	log.Printf("device set struct %+v\n\n", deviceset)
	log.Printf("device set struct %+v\n\n", deviceset.Item[3])
	log.Printf("device set struct %+v\n\n", deviceset.Item[33])
	return "set things are complicated"

}

func decodeDeviceArrival(payload string) string {
	//generate at  https://www.onlinetool.io/xmltogo/
	//log.Printf("decodeDeviceArrival payload id %s", payload)
	type DeviceArrival struct {
		XMLName xml.Name `xml:"device-arrival"`
		Text    string   `xml:",chardata"`
		Device  struct {
			Text         string `xml:",chardata"`
			ID           string `xml:"id,attr"`
			Protocol     string `xml:"protocol,attr"`
			Model        string `xml:"model,attr"`
			Class        string `xml:"class,attr"`
			BusID        string `xml:"bus-id,attr"`
			SerialNumber string `xml:"serial-number,attr"`
			Version      string `xml:"version,attr"`
			Nickname     struct {
				Text string `xml:",chardata"`
				ID   string `xml:"id,attr"`
			} `xml:"nickname"`
			SealBroken struct {
				Text string `xml:",chardata"`
				ID   string `xml:"id,attr"`
			} `xml:"seal-broken"`
			Snapshot struct {
				Text string `xml:",chardata"`
				ID   string `xml:"id,attr"`
			} `xml:"snapshot"`
			SaveSnapshot struct {
				Text string `xml:",chardata"`
				ID   string `xml:"id,attr"`
			} `xml:"save-snapshot"`
			ResetDevice struct {
				Text string `xml:",chardata"`
				ID   string `xml:"id,attr"`
			} `xml:"reset-device"`
			Preset struct {
				Text string `xml:",chardata"`
				ID   string `xml:"id,attr"`
				Enum []struct {
					Text  string `xml:",chardata"`
					Value string `xml:"value,attr"`
				} `xml:"enum"`
			} `xml:"preset"`
			Firmware struct {
				Text    string `xml:",chardata"`
				Version struct {
					Text string `xml:",chardata"`
					ID   string `xml:"id,attr"`
				} `xml:"version"`
				NeedsUpdate struct {
					Text string `xml:",chardata"`
					ID   string `xml:"id,attr"`
				} `xml:"needs-update"`
				FirmwareProgress struct {
					Text string `xml:",chardata"`
					ID   string `xml:"id,attr"`
				} `xml:"firmware-progress"`
				UpdateFirmware struct {
					Text string `xml:",chardata"`
					ID   string `xml:"id,attr"`
				} `xml:"update-firmware"`
				RestoreFactory struct {
					Text string `xml:",chardata"`
					ID   string `xml:"id,attr"`
				} `xml:"restore-factory"`
			} `xml:"firmware"`
			Mixer struct {
				Text      string `xml:",chardata"`
				Available struct {
					Text string `xml:",chardata"`
					ID   string `xml:"id,attr"`
				} `xml:"available"`
				Inputs struct {
					Text     string `xml:",chardata"`
					AddInput struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"add-input"`
					AddInputWithoutReset struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"add-input-without-reset"`
					AddStereoInput struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"add-stereo-input"`
					RemoveInput struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"remove-input"`
					FreeInputs struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"free-inputs"`
					Input []struct {
						Text   string `xml:",chardata"`
						Source struct {
							Text string `xml:",chardata"`
							ID   string `xml:"id,attr"`
						} `xml:"source"`
						Stereo struct {
							Text string `xml:",chardata"`
							ID   string `xml:"id,attr"`
						} `xml:"stereo"`
					} `xml:"input"`
				} `xml:"inputs"`
				Mixes struct {
					Text string `xml:",chardata"`
					Mix  []struct {
						Text       string `xml:",chardata"`
						ID         string `xml:"id,attr"`
						Name       string `xml:"name,attr"`
						StereoName string `xml:"stereo-name,attr"`
						Meter      struct {
							Text string `xml:",chardata"`
							ID   string `xml:"id,attr"`
						} `xml:"meter"`
						Input []struct {
							Text string `xml:",chardata"`
							Gain struct {
								Text string `xml:",chardata"`
								ID   string `xml:"id,attr"`
							} `xml:"gain"`
							Pan struct {
								Text string `xml:",chardata"`
								ID   string `xml:"id,attr"`
							} `xml:"pan"`
							Mute struct {
								Text string `xml:",chardata"`
								ID   string `xml:"id,attr"`
							} `xml:"mute"`
							Solo struct {
								Text string `xml:",chardata"`
								ID   string `xml:"id,attr"`
							} `xml:"solo"`
						} `xml:"input"`
					} `xml:"mix"`
				} `xml:"mixes"`
			} `xml:"mixer"`
			Inputs struct {
				Text     string `xml:",chardata"`
				Analogue []struct {
					Text             string `xml:",chardata"`
					ID               string `xml:"id,attr"`
					SupportsTalkback string `xml:"supports-talkback,attr"`
					Hidden           string `xml:"hidden,attr"`
					Name             string `xml:"name,attr"`
					StereoName       string `xml:"stereo-name,attr"`
					Available        struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"available"`
					Meter struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"meter"`
					Nickname struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"nickname"`
					Mode struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
						Enum []struct {
							Text  string `xml:",chardata"`
							Value string `xml:"value,attr"`
						} `xml:"enum"`
					} `xml:"mode"`
					Air struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"air"`
					Pad struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"pad"`
				} `xml:"analogue"`
				SpdifRca []struct {
					Text             string `xml:",chardata"`
					ID               string `xml:"id,attr"`
					SupportsTalkback string `xml:"supports-talkback,attr"`
					Hidden           string `xml:"hidden,attr"`
					Name             string `xml:"name,attr"`
					StereoName       string `xml:"stereo-name,attr"`
					Available        struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"available"`
					Meter struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"meter"`
					Nickname struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"nickname"`
				} `xml:"spdif-rca"`
				Playback []struct {
					Text             string `xml:",chardata"`
					ID               string `xml:"id,attr"`
					SupportsTalkback string `xml:"supports-talkback,attr"`
					Hidden           string `xml:"hidden,attr"`
					Name             string `xml:"name,attr"`
					StereoName       string `xml:"stereo-name,attr"`
					Available        struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"available"`
					Meter struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"meter"`
					Nickname struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"nickname"`
				} `xml:"playback"`
			} `xml:"inputs"`
			Outputs struct {
				Text     string `xml:",chardata"`
				Analogue []struct {
					Text             string `xml:",chardata"`
					Name             string `xml:"name,attr"`
					StereoName       string `xml:"stereo-name,attr"`
					Headphone        string `xml:"headphone,attr"`
					Independent      string `xml:"independent,attr"`
					Monitor          string `xml:"monitor,attr"`
					HardwareControls string `xml:"hardware-controls,attr"`
					Available        struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"available"`
					Meter struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"meter"`
					AssignMix struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"assign-mix"`
					AssignTalkbackMix struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"assign-talkback-mix"`
					Mute struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"mute"`
					Source struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"source"`
					Stereo struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"stereo"`
					Nickname struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"nickname"`
					HardwareControl struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"hardware-control"`
					Gain struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"gain"`
				} `xml:"analogue"`
				SpdifRca []struct {
					Text       string `xml:",chardata"`
					Name       string `xml:"name,attr"`
					StereoName string `xml:"stereo-name,attr"`
					Available  struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"available"`
					Meter struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"meter"`
					AssignMix struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"assign-mix"`
					AssignTalkbackMix struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"assign-talkback-mix"`
					Mute struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"mute"`
					Source struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"source"`
					Stereo struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"stereo"`
					Nickname struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"nickname"`
				} `xml:"spdif-rca"`
				Loopback []struct {
					Text       string `xml:",chardata"`
					Name       string `xml:"name,attr"`
					StereoName string `xml:"stereo-name,attr"`
					Available  struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"available"`
					Meter struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"meter"`
					AssignMix struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"assign-mix"`
					AssignTalkbackMix struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"assign-talkback-mix"`
					Mute struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"mute"`
					Source struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"source"`
					Stereo struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"stereo"`
					Nickname struct {
						Text string `xml:",chardata"`
						ID   string `xml:"id,attr"`
					} `xml:"nickname"`
				} `xml:"loopback"`
			} `xml:"outputs"`
			RecordOutputs string `xml:"record-outputs"`
			Monitoring    struct {
				Text              string `xml:",chardata"`
				MonitorGroupPairs string `xml:"monitor-group-pairs"`
			} `xml:"monitoring"`
			Clocking struct {
				Text   string `xml:",chardata"`
				Locked struct {
					Text string `xml:",chardata"`
					ID   string `xml:"id,attr"`
				} `xml:"locked"`
				ClockSource struct {
					Text string `xml:",chardata"`
					ID   string `xml:"id,attr"`
					Enum []struct {
						Text  string `xml:",chardata"`
						Value string `xml:"value,attr"`
					} `xml:"enum"`
				} `xml:"clock-source"`
				SampleRate struct {
					Text string `xml:",chardata"`
					ID   string `xml:"id,attr"`
					Enum []struct {
						Text  string `xml:",chardata"`
						Value string `xml:"value,attr"`
					} `xml:"enum"`
				} `xml:"sample-rate"`
				ClockMaster struct {
					Text string `xml:",chardata"`
					ID   string `xml:"id,attr"`
				} `xml:"clock-master"`
			} `xml:"clocking"`
			Settings struct {
				Text       string `xml:",chardata"`
				BufferSize struct {
					Text string `xml:",chardata"`
					ID   string `xml:"id,attr"`
				} `xml:"buffer-size"`
				PhantomPersistence struct {
					Text string `xml:",chardata"`
					ID   string `xml:"id,attr"`
				} `xml:"phantom-persistence"`
				DelayCompensation string `xml:"delay-compensation"`
			} `xml:"settings"`
			Dante      string `xml:"dante"`
			State      string `xml:"state"`
			QuickStart struct {
				Text    string `xml:",chardata"`
				URL     string `xml:"url,attr"`
				MsdMode struct {
					Text string `xml:",chardata"`
					ID   string `xml:"id,attr"`
				} `xml:"msd-mode"`
			} `xml:"quick-start"`
			HaloSettings struct {
				Text             string `xml:",chardata"`
				AvailableColours struct {
					Text string `xml:",chardata"`
					Enum []struct {
						Text  string `xml:",chardata"`
						Value string `xml:"value,attr"`
					} `xml:"enum"`
				} `xml:"available-colours"`
				GoodMeterColour struct {
					Text string `xml:",chardata"`
					ID   string `xml:"id,attr"`
				} `xml:"good-meter-colour"`
				PreClipMeterColour struct {
					Text string `xml:",chardata"`
					ID   string `xml:"id,attr"`
				} `xml:"pre-clip-meter-colour"`
				ClippingMeterColour struct {
					Text string `xml:",chardata"`
					ID   string `xml:"id,attr"`
				} `xml:"clipping-meter-colour"`
				EnablePreviewMode struct {
					Text string `xml:",chardata"`
					ID   string `xml:"id,attr"`
				} `xml:"enable-preview-mode"`
				Halos struct {
					Text string `xml:",chardata"`
					ID   string `xml:"id,attr"`
				} `xml:"halos"`
			} `xml:"halo-settings"`
		} `xml:"device"`
	}

	var devicearrival DeviceArrival

	xml.Unmarshal([]byte(payload), &devicearrival)
	log.Printf("device arrival struct %+v\n\n", devicearrival)
	log.Printf("device arrival struct mixes %+v\n\n", devicearrival.Device.Mixer)
	log.Printf("device arrival struct device ID %+v\n\n", devicearrival.Device.ID)
	log.Printf("this is the protocol of the device %s", devicearrival.Device.Protocol)
	return "things are complicated"

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
		log.Printf("space is: x%sx", space)

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

		response = string(payload)

	}
	return response
}
