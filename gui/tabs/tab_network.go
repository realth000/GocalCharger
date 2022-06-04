package tabs

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	cConfig "gocalcharger/client/config"
	sConfig "gocalcharger/server/config"
	"strconv"
)

var (
	serverStatus = binding.NewString()
)

// Server configs
var (
	ServerPort        = binding.NewString()
	ServerPermitFiles = binding.NewString()
	ServerSSLEnabled  = binding.NewBool()
	ServerSSLCert     = binding.NewString()
	ServerSSLKey      = binding.NewString()
	ServerCACert      = binding.NewString()
)

// Client configs
var (
	ClientRemoteServerIP      = binding.NewString()
	ClientRemoteServerPort    = binding.NewString()
	ClientName                = binding.NewString()
	ClientSSLEnable           = binding.NewBool()
	ClientSSLCert             = binding.NewString()
	ClientSSLKey              = binding.NewString()
	ClientSSLCACert           = binding.NewString()
	ClientSSLMutualAuth       = binding.NewBool()
	ClientSSLDownloadFile     = binding.NewBool()
	ClientSSLDownloadFilePath = binding.NewString()
)

// SSL control
var (
	serverSSLCheck    *widget.Check
	serverCertEntry   *widget.Entry
	serverKeyEntry    *widget.Entry
	ServerCACertEntry *widget.Entry
	clientSSLCheck    *widget.Check
	clientCertEntry   *widget.Entry
	clientKeyEntry    *widget.Entry
	clientCACertEntry *widget.Entry
)

func init() {
	_ = serverStatus.Set("closed")
}

func newServerControlArea() fyne.CanvasObject {
	netConfig := makeServerNetworkConfigArea()
	sslConfig := makeServerSSLConfigArea()
	updateServerSSL()
	applyButton := widget.NewButton("Apply", reloadServerConfig)
	return widget.NewCard("Server", "Server network configuration", container.NewVBox(netConfig, sslConfig, container.NewHBox(applyButton)))
}

func makeServerNetworkConfigArea() fyne.CanvasObject {
	portLabel := widget.NewLabel("Localhost port")
	portEntry := widget.NewEntryWithData(ServerPort)
	portHBox := container.New(layout.NewFormLayout(), portLabel, portEntry)
	return portHBox
}

func makeServerSSLConfigArea() *fyne.Container {
	serverSSLCheck = widget.NewCheck("Enable SSL", func(b bool) { _ = ServerSSLEnabled.Set(b); updateServerSSL() })
	b, _ := ServerSSLEnabled.Get()
	serverSSLCheck.SetChecked(b)
	certLabel := widget.NewLabel("Certificate path")
	serverCertEntry = widget.NewEntryWithData(ServerSSLCert)
	serverCertEntry.SetPlaceHolder("*.pem")
	keyLabel := widget.NewLabel("Private key path")
	serverKeyEntry = widget.NewEntryWithData(ServerSSLKey)
	serverKeyEntry.SetPlaceHolder("*.key")
	caCertLabel := widget.NewLabel("CA certificate path")
	ServerCACertEntry = widget.NewEntryWithData(ServerCACert)
	ServerCACertEntry.SetPlaceHolder("*.pem")
	return container.NewVBox(serverSSLCheck, container.New(layout.NewFormLayout(), certLabel, serverCertEntry, keyLabel, serverKeyEntry, caCertLabel, ServerCACertEntry))
}

func newClientControlArea() fyne.CanvasObject {
	netConfig := makeClientNetworkConfigArea()
	sslConfig := makeClientSSLConfigArea()
	updateClientSSL()
	applyButton := widget.NewButton("Apply", reloadClientConfig)
	testConnectServerButton := widget.NewButton("Test connect server", testConnectServer)
	buttonBox := container.NewHBox(applyButton, testConnectServerButton)
	return widget.NewCard("Client", "Client network configuration", container.NewVBox(netConfig, sslConfig, buttonBox))
}

func updateServerSSL() {
	b, _ := ServerSSLEnabled.Get()
	if b {
		if serverCertEntry != nil {
			serverCertEntry.Enable()
		}
		if serverKeyEntry != nil {
			serverKeyEntry.Enable()
		}
		if ServerCACertEntry != nil {
			ServerCACertEntry.Enable()
		}
	} else {
		if serverCertEntry != nil {
			serverCertEntry.Disable()
		}
		if serverKeyEntry != nil {
			serverKeyEntry.Disable()
		}
		if ServerCACertEntry != nil {
			ServerCACertEntry.Disable()
		}
	}
}

func makeClientNetworkConfigArea() fyne.CanvasObject {
	serverIPLabel := widget.NewLabel("Remote server IP")
	serverIPEntry := widget.NewEntryWithData(ClientRemoteServerIP)
	serverPortLabel := widget.NewLabel("Remote server port")
	serverPortEntry := widget.NewEntryWithData(ClientRemoteServerPort)
	return container.New(layout.NewFormLayout(), serverIPLabel, serverIPEntry, serverPortLabel, serverPortEntry)
}

func makeClientSSLConfigArea() fyne.CanvasObject {
	clientSSLCheck = widget.NewCheck("Enable SSL", func(b bool) { _ = ClientSSLEnable.Set(b); updateClientSSL() })
	b, _ := ClientSSLEnable.Get()
	clientSSLCheck.SetChecked(b)
	certLabel := widget.NewLabel("Certificate path")
	clientCertEntry = widget.NewEntryWithData(ClientSSLCert)
	clientCertEntry.SetPlaceHolder("*.pem")
	keyLabel := widget.NewLabel("Private key path")
	clientKeyEntry = widget.NewEntryWithData(ClientSSLKey)
	clientKeyEntry.SetPlaceHolder("*.key")
	caCertLabel := widget.NewLabel("CA certificate path")
	clientCACertEntry = widget.NewEntryWithData(ClientSSLCACert)
	clientCACertEntry.SetPlaceHolder("*.pem")
	return container.NewVBox(clientSSLCheck, container.New(layout.NewFormLayout(), certLabel, clientCertEntry, keyLabel, clientKeyEntry, caCertLabel, clientCACertEntry))
}

func updateClientSSL() {
	b, _ := ClientSSLEnable.Get()
	if b {
		if clientCertEntry != nil {
			clientCertEntry.Enable()
		}
		if clientKeyEntry != nil {
			clientKeyEntry.Enable()
		}
		if clientCACertEntry != nil {
			clientCACertEntry.Enable()
		}
	} else {
		if clientCertEntry != nil {
			clientCertEntry.Disable()
		}
		if clientKeyEntry != nil {
			clientKeyEntry.Disable()
		}
		if clientCACertEntry != nil {
			clientCACertEntry.Disable()
		}
	}
}
func NewNetworkTab() *container.TabItem {
	serverArea := newServerControlArea()
	clientArea := newClientControlArea()
	networkTab := container.NewTabItem("Network", container.NewVBox(serverArea, clientArea))
	return networkTab
}

func reloadServerConfig() {
	fmt.Println("reload server network config")
}

func SetServerSSL(b bool) {
	serverSSLCheck.SetChecked(b)
}

func reloadClientConfig() {
	fmt.Println("reload client network config")
}

func testConnectServer() {

}

func ApplyConfigs(s sConfig.ServerConfig, c cConfig.ClientConfig) {
	_ = ServerPort.Set(strconv.Itoa(int(s.Port)))
	_ = ServerPermitFiles.Set(s.PermitFiles)
	_ = ServerSSLEnabled.Set(s.SSL)
	_ = ServerSSLCert.Set(s.SSLCert)
	_ = ServerSSLKey.Set(s.SSLKey)
	_ = ServerCACert.Set(s.SSLCACert)

	_ = ClientRemoteServerIP.Set(c.ServerUrl)
	_ = ClientRemoteServerPort.Set(strconv.Itoa(int(c.ServerPort)))
	_ = ClientName.Set(c.ClientName)
	_ = ClientSSLEnable.Set(c.SSL)
	_ = ClientSSLCert.Set(c.SSLCert)
	_ = ClientSSLKey.Set(c.SSLKey)
	_ = ClientSSLCACert.Set(c.SSLCACert)
	_ = ClientSSLMutualAuth.Set(c.MutualAuth)
	_ = ClientSSLDownloadFile.Set(c.DownloadFile)
	_ = ClientSSLDownloadFilePath.Set(c.DownloadFilePath)
}
