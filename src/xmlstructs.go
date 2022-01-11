package main

import "encoding/xml"

type FocusriteMessage struct {
	XMLName xml.Name `xml:"focusritemessage"`
	//These are all really top level elements, but the response is wrapped in a dummy focusritemessage xml for unmarshalling

	KeepAlive struct {
		XMLName xml.Name `xml:"keep-alive"`
	}

	ClientDetails struct {
		XMLName xml.Name `xml:"client-details"`
		Text    string   `xml:",chardata"`
		ID      string   `xml:"id,attr"`
	}

	Approval struct {
		XMLName    xml.Name `xml:"approval"`
		Text       string   `xml:",chardata"`
		Hostname   string   `xml:"hostname,attr"`
		ID         string   `xml:"id,attr"`
		Type       string   `xml:"type,attr"`
		Authorised string   `xml:"authorised,attr"`
	}

	DeviceSet struct {
		XMLName xml.Name `xml:"set"`
		Text    string   `xml:",chardata"`
		Devid   string   `xml:"devid,attr"`
		Item    []struct {
			Text  string `xml:",chardata"`
			ID    int    `xml:"id,attr"`
			Value string `xml:"value,attr"`
		} `xml:"item"`
	}

	DeviceArrival struct {
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
					ID               int    `xml:"id,attr"`
					SupportsTalkback string `xml:"supports-talkback,attr"`
					Hidden           string `xml:"hidden,attr"`
					Name             string `xml:"name,attr"`
					StereoName       string `xml:"stereo-name,attr"`
					Available        struct {
						Text string `xml:",chardata"`
						ID   int    `xml:"id,attr"`
					} `xml:"available"`
					Meter struct {
						Text string `xml:",chardata"`
						ID   int    `xml:"id,attr"`
					} `xml:"meter"`
					Nickname struct {
						Text string `xml:",chardata"`
						ID   int    `xml:"id,attr"`
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
					ID   int    `xml:"id,attr"`
				} `xml:"locked"`
				ClockSource struct {
					Text string `xml:",chardata"`
					ID   int    `xml:"id,attr"`
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
}
