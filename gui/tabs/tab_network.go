package tabs

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var (
	serverStatus     = binding.NewString()
	serverSSLEnabled = true
	clientSSLEnable  = true
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
	return widget.NewCard("Server", "Server network configuration", container.NewVBox(netConfig, sslConfig))
}

func makeServerNetworkConfigArea() fyne.CanvasObject {
	portLabel := widget.NewLabel("Localhost port")
	portEntry := widget.NewEntry()
	portHBox := container.New(layout.NewFormLayout(), portLabel, portEntry)
	return portHBox
}

func makeServerSSLConfigArea() *fyne.Container {
	serverSSLCheck = widget.NewCheck("Enable SSL", func(b bool) { serverSSLEnabled = b; updateServerSSL() })
	serverSSLEnabled = serverSSLCheck.Checked
	certLabel := widget.NewLabel("Certificate path")
	serverCertEntry = widget.NewEntry()
	serverCertEntry.SetPlaceHolder("*.pem")
	keyLabel := widget.NewLabel("Private key path")
	serverKeyEntry = widget.NewEntry()
	serverKeyEntry.SetPlaceHolder("*.key")
	caCertLabel := widget.NewLabel("CA certificate path")
	serverCACertEntry = widget.NewEntry()
	serverCACertEntry.SetPlaceHolder("*.pem")
	return container.NewVBox(serverSSLCheck, container.New(layout.NewFormLayout(), certLabel, serverCertEntry, keyLabel, serverKeyEntry, caCertLabel, serverCACertEntry))
}

func newClientControlArea() fyne.CanvasObject {
	netConfig := makeClientNetworkConfigArea()
	sslConfig := makeClientSSLConfigArea()
	updateClientSSL()
	return widget.NewCard("Client", "Client network configuration", container.NewVBox(netConfig, sslConfig))
}

func updateServerSSL() {
	if serverSSLEnabled {
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
	clientSSLCheck = widget.NewCheck("Enable SSL", func(b bool) { clientSSLEnable = b; updateClientSSL() })
	clientSSLEnable = clientSSLCheck.Checked
	certLabel := widget.NewLabel("Certificate path")
	clientCertEntry = widget.NewEntry()
	clientCertEntry.SetPlaceHolder("*.pem")
	keyLabel := widget.NewLabel("Private key path")
	clientKeyEntry = widget.NewEntry()
	clientKeyEntry.SetPlaceHolder("*.key")
	caCertLabel := widget.NewLabel("CA certificate path")
	clientCACertEntry = widget.NewEntry()
	clientCACertEntry.SetPlaceHolder("*.pem")
	return container.NewVBox(clientSSLCheck, container.New(layout.NewFormLayout(), certLabel, clientCertEntry, keyLabel, clientKeyEntry, caCertLabel, clientCACertEntry))
}

func updateClientSSL() {
	if clientSSLEnable {
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

func SetServerSSL(b bool) {
	serverSSLCheck.SetChecked(b)
}
