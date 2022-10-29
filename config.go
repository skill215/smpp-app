package smppapp

type SmppConfig struct {
	SmppApp struct {
		SmppServer struct {
		} `yaml:"smpp-server"`
		SmppClient struct {
		} `yaml:"smpp-client"`
		SmppMessage struct {
		} `yaml:"smpp-ao"`
		Service struct {
		} `yaml:"service"`
		Log struct {
		} `yaml:"log"`
	} `yaml:"smpp-app"`
}
