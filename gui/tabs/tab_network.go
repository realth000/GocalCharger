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
	serverPort        = binding.NewString()
	serverPermitFiles = binding.NewString()
	serverSSLEnabled  = binding.NewBool()
	serverSSLCert     = binding.NewString()
	serverSSLKey      = binding.NewString()
	serverCACert      = binding.NewString()
)

// Client configs
var (
	clientRemoteServerIP      = binding.NewString()
	clientRemoteServerPort    = binding.NewString()
	clientName                = binding.NewString()
	clientSSLEnable           = binding.NewBool()
	clientSSLCert             = binding.NewString()
	clientSSLKey              = binding.NewString()
	clientSSLCACert           = binding.NewString()
	clientSSLMutualAuth       = binding.NewBool()
	clientSSLDownloadFile     = binding.NewBool()
	clientSSLDownloadFilePath = binding.NewString()
)

// SSL control
var (
	serverSSLCheck    *widget.Check
	serverCertEntry   *widget.Entry
	serverKeyEntry    *widget.Entry
	serverCACertEntry *widget.Entry
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
	portEntry := widget.NewEntry()
	portHBox := container.New(layout.NewFormLayout(), portLabel, portEntry)
	return portHBox
}

func makeServerSSLConfigArea() *fyne.Container {
	serverSSLCheck = widget.NewCheck("Enable SSL", func(b bool) { _ = serverSSLEnabled.Set(b); updateServerSSL() })
	certLabel := widget.NewLabel("Certificate path")
	serverCertEntry = widget.NewEntryWithData(serverSSLCert)
	serverCertEntry.SetPlaceHolder("*.pem")
	keyLabel := widget.NewLabel("Private key path")
	serverKeyEntry = widget.NewEntryWithData(serverSSLKey)
	serverKeyEntry.SetPlaceHolder("*.key")
	caCertLabel := widget.NewLabel("CA certificate path")
	serverCACertEntry = widget.NewEntryWithData(serverCACert)
	serverCACertEntry.SetPlaceHolder("*.pem")
	return container.NewVBox(serverSSLCheck, container.New(layout.NewFormLayout(), certLabel, serverCertEntry, keyLabel, serverKeyEntry, caCertLabel, serverCACertEntry))
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
	b, _ := serverSSLEnabled.Get()
	if b {
		if serverCertEntry != nil {
			serverCertEntry.Enable()
		}
		if serverKeyEntry != nil {
			serverKeyEntry.Enable()
		}
		if serverCACertEntry != nil {
			serverCACertEntry.Enable()
		}
	} else {
		if serverCertEntry != nil {
			serverCertEntry.Disable()
		}
		if serverKeyEntry != nil {
			serverKeyEntry.Disable()
		}
		if serverCACertEntry != nil {
			serverCACertEntry.Disable()
		}
	}
}

func makeClientNetworkConfigArea() fyne.CanvasObject {
	serverIPLabel := widget.NewLabel("Remote server IP")
	serverIPEntry := widget.NewEntry()
	serverPortLabel := widget.NewLabel("Remote server port")
	serverPortEntry := widget.NewEntry()
	return container.New(layout.NewFormLayout(), serverIPLabel, serverIPEntry, serverPortLabel, serverPortEntry)
}

func makeClientSSLConfigArea() fyne.CanvasObject {
	clientSSLCheck = widget.NewCheck("Enable SSL", func(b bool) { clientSSLEnable.Set(b); updateClientSSL() })
	certLabel := widget.NewLabel("Certificate path")
	clientCertEntry = widget.NewEntryWithData(clientSSLCert)
	clientCertEntry.SetPlaceHolder("*.pem")
	keyLabel := widget.NewLabel("Private key path")
	clientKeyEntry = widget.NewEntryWithData(clientSSLKey)
	clientKeyEntry.SetPlaceHolder("*.key")
	caCertLabel := widget.NewLabel("CA certificate path")
	clientCACertEntry = widget.NewEntryWithData(clientSSLCACert)
	clientCACertEntry.SetPlaceHolder("*.pem")
	return container.NewVBox(clientSSLCheck, container.New(layout.NewFormLayout(), certLabel, clientCertEntry, keyLabel, clientKeyEntry, caCertLabel, clientCACertEntry))
}

func updateClientSSL() {
	b, _ := clientSSLEnable.Get()
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
	fmt.Println("test connect Server")
}

func ApplyConfigs(s sConfig.ServerConfig, c cConfig.ClientConfig) {
	_ = serverPort.Set(strconv.Itoa(int(s.Port)))
	_ = serverPermitFiles.Set(s.PermitFiles)
	_ = serverSSLEnabled.Set(s.SSL)
	_ = serverSSLCert.Set(s.SSLCert)
	_ = serverSSLKey.Set(s.SSLKey)
	_ = serverCACert.Set(s.SSLCACert)

	_ = clientRemoteServerIP.Set(c.ServerUrl)
	_ = clientRemoteServerPort.Set(strconv.Itoa(int(c.ServerPort)))
	_ = clientName.Set(c.ClientName)
	_ = clientSSLEnable.Set(c.SSL)
	_ = clientSSLCert.Set(c.SSLCert)
	_ = clientSSLKey.Set(c.SSLKey)
	_ = clientSSLCACert.Set(c.SSLCACert)
	_ = clientSSLMutualAuth.Set(c.MutualAuth)
	_ = clientSSLDownloadFile.Set(c.DownloadFile)
	_ = clientSSLDownloadFilePath.Set(c.DownloadFilePath)
}
