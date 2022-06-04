package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	cConfig "gocalcharger/client/config"
	"gocalcharger/gui/client"
	"gocalcharger/gui/server"
	"gocalcharger/gui/tabs"
	sConfig "gocalcharger/server/config"
	"log"
	"os"
	"time"
)

const (
	serverConfigPath = `./tests/data/config/server.toml`
	clientConfigPath = `./tests/data/config/client.toml`
)

var (
	serverConfig sConfig.ServerConfig
	clientConfig cConfig.ClientConfig
)

func updateTime() {
	if len(tabs.Items) > 0 {
		//tabs.Items[len(tabs.Items)-1].RemainTime = time.Now().Format("Time: 03:04:05")
		//tabs.Update()
	}
}

func testLoadConfig() {
	if _, err := os.Stat(serverConfigPath); err != nil {
		log.Fatalf("error loading server config: %v", err)
	}
	err := sConfig.LoadConfigFile(serverConfigPath, &serverConfig)
	if err != nil {
		log.Fatalf("error loading server config: %v", err)
	}

	if _, err = os.Stat(clientConfigPath); err != nil {
		log.Fatalf("error loading client config: %v", err)
	}
	err = cConfig.LoadConfigFile(clientConfigPath, &clientConfig)
	if err != nil {
		log.Fatalf("error loading client config: %v", err)
	}
}

func testApplyConfig() {
	tabs.ApplyConfigs(serverConfig, clientConfig)
}

func main() {
	// Test loading configs.
	testLoadConfig()
	testApplyConfig()

	a := app.New()
	w := a.NewWindow("GocalChargerGui")
	w.Resize(fyne.NewSize(800, 600))
	downloadTab := tabs.NewDownloadTab()
	networkTab := tabs.NewNetworkTab()
	settingsTab := tabs.NewSettingsTab()
	tab := container.NewAppTabs(downloadTab, networkTab, settingsTab)
	w.SetContent(tab)

	w.Show()

	go func() {
		for range time.Tick(time.Second) {
			updateTime()
		}
	}()
	client.StartReceivingChannels()
	go server.StartServer()
	a.Run()
}
