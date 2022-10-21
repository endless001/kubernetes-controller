package manager

import (
	"github.com/spf13/pflag"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"kubernetes-controller/internal/annotations"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

type Config struct {
	LogLevel  string
	LogFormat string

	APIServerHost    string
	CacheSyncTimeout time.Duration
	SyncPeriod       time.Duration

	KubeConfigPath           string
	IngressClassName         string
	IngressExtV1beta1Enabled bool
	IngressNetV1beta1Enabled bool
	IngressNetV1Enabled      bool
	IngressClassNetV1Enabled bool
	ServiceEnabled           bool
	LeaderElectionID         string

	UpdateStatus bool

	MetricsAddr string
	ProbeAddr   string

	TermDelay time.Duration
}

func (c *Config) FlagSet() *pflag.FlagSet {
	flagSet := pflag.NewFlagSet("", pflag.ExitOnError)
	flagSet.StringVar(&c.KubeConfigPath, "kubeconfig", "C:\\Users\\longqing\\.kube\\config", "Path to the kubeconfig file.")

	flagSet.StringVar(&c.IngressClassName, "ingress-class", annotations.DefaultIngressClass, `Name of the ingress class to route through this controller.`)
	
	flagSet.BoolVar(&c.IngressNetV1Enabled, "enable-controller-ingress-networkingv1", true, "Enable the networking.k8s.io/v1 Ingress controller.")
	return flagSet
}
func (c *Config) GetKubeConfig() (*rest.Config, error) {
	config, err := clientcmd.BuildConfigFromFlags(c.APIServerHost, c.KubeConfigPath)
	if err != nil {
		return nil, err
	}
	return config, err
}
func (c *Config) GetKubeClient() (client.Client, error) {
	config, err := c.GetKubeConfig()
	if err != nil {
		return nil, err
	}
	return client.New(config, client.Options{})
}
