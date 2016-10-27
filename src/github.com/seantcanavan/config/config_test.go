package config

import (
	"testing"

	"github.com/seantcanavan/utils"
)

func TestSimplecfgigLoadFromFile(t *testing.T) {
	err := ConfigFromFile(LOCAL_INTERNAL_PATH)

	if err != nil {
		t.Errorf("generic unmarshal error: %v", err)
	}

	if Cfg.CheckInFrequencySeconds != 3600 {
		t.Errorf("Cfg.CheckInFrequencySeconds did not unmarshal correctly: %v", Cfg.CheckInFrequencySeconds)
	}

	if Cfg.UpdateFrequencySeconds != 3600 {
		t.Errorf("Cfg.UpdateFrequencySeconds did not unmarshal correctly: %v", Cfg.UpdateFrequencySeconds)
	}

	if Cfg.RemoteUpdateURI != "https://github.com/seantcanavan/anon-eth-net.git" {
		t.Errorf("Cfg.RemoteUpdateURI did not unmarshal correctly: %v", Cfg.RemoteUpdateURI)
	}

	if Cfg.RemoteVersionURI != "https://raw.githubusercontent.com/seantcanavan/anon-eth-net/master/src/github.com/seantcanavan/main/version.no" {
		t.Errorf("Cfg.RemoteVersionURI did not unmarshal correctly: %v", Cfg.RemoteVersionURI)
	}

	if Cfg.LocalVersionURI != "version.no" {
		t.Errorf("Cfg.LocalVersionURI did not unmarshal correctly: %v", Cfg.LocalVersionURI)
	}

	if Cfg.MineEther != false {
		t.Errorf("Cfg.MineEther did not unmarshal correctly: %v", Cfg.MineEther)
	}

	if Cfg.GPUMine != false {
		t.Errorf("Cfg.GPUMine did not unmarshal correctly: %v", Cfg.MineEther)
	}

	if Cfg.CPUMine != true {
		t.Errorf("Cfg.CPUMine did not unmarshal correctly: %v", Cfg.CPUMine)
	}

	if Cfg.EtherWallet != "" {
		t.Errorf("Cfg.EtherWallet did not unmarshal correctly: %v", Cfg.EtherWallet)
	}

	if Cfg.DeviceId != "519a2a15-afad-4c1a-94a3-71660c83504b" {
		t.Errorf("Cfg.DeviceId did not unmarshal correctly: %v", Cfg.DeviceId)
	}

	if Cfg.DeviceName != "LG Smart Fridge" {
		t.Errorf("Cfg.DeviceName did not unmarshal correctly: %v", Cfg.DeviceName)
	}

	if Cfg.InitialStartup != false {
		t.Errorf("Cfg.InitialStartup did not unmarshal correctly: %v", Cfg.InitialStartup)
	}

	if Cfg.FirstRunAfterUpdate != false {
		t.Errorf("Cfg.FirstRunAfterUpDate did not unmarshal correctly: %v", Cfg.FirstRunAfterUpdate)
	}

	if Cfg.LogVolatility != 3 {
		t.Errorf("Cfg.LogVolatility did not unmarshal correctly: %v", Cfg.LogVolatility)
	}

	fileLines, fileErr := utils.ReadLines(Cfg.CheckInGmailCredentialsFile)
	if fileErr != nil {
		t.Errorf("issue reading in email credentials file")
	}

	if Cfg.CheckInGmailAddress != fileLines[0] {
		t.Errorf("Cfg.CheckInGmailAddress did not load correctly: %v", Cfg.CheckInGmailAddress)
	}

	if Cfg.CheckInGmailPassword != fileLines[1] {
		t.Errorf("Cfg.CheckInGmailPassword did not load correctly: %v", Cfg.CheckInGmailPassword)
	}

}
