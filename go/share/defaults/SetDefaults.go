package defaults

/*
func LoadDefaultImplementations() {}

const (
	DEFAULT_MAX_DATA_SIZE     = 1024 * 1024
	DEFAULT_EDGE_QUEUE_SIZE   = 10000
	DEFAULT_SWITCH_QUEUE_SIZE = 50000
	DEFAULT_SWITCH_PORT       = 50000
)

var DefaultSecurityProvider interfaces.ISecurityProvider

func init() {
	initLogger()
	initEdgeConfig()
	initSecurityProvider()
}

func initLogger() {
	interfaces.SetLogger(logger.NewLoggerImpl(&logger.FmtLogMethod{}))
}

func initEdgeConfig() {
	interfaces.SetEdgeConfig(interfaces.NewMessageConfig(DEFAULT_MAX_DATA_SIZE, DEFAULT_EDGE_QUEUE_SIZE, DEFAULT_EDGE_QUEUE_SIZE, DEFAULT_SWITCH_PORT, true, 30))
	interfaces.SetEdgeSwitchConfig(interfaces.NewMessageConfig(DEFAULT_MAX_DATA_SIZE, DEFAULT_EDGE_QUEUE_SIZE, DEFAULT_EDGE_QUEUE_SIZE, DEFAULT_SWITCH_PORT, false, 0))
	interfaces.SetSwitchConfig(interfaces.NewMessageConfig(DEFAULT_MAX_DATA_SIZE, DEFAULT_SWITCH_QUEUE_SIZE, DEFAULT_SWITCH_QUEUE_SIZE, DEFAULT_SWITCH_PORT, true, 30))
}

func initSecurityProvider() {
	hash := md5.New()
	secret := "Default Security Provider"
	hash.Write([]byte(secret))
	kHash := hash.Sum(nil)
	k := base64.StdEncoding.EncodeToString(kHash)
	DefaultSecurityProvider = shallow_security.NewShallowSecurityProvider(k, secret)
}
*/
