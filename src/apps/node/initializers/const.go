package initializers

import (
	"os"
	"time"

	"github.com/vany-egorov/ha-eta/lib/environment"
)

const (
	EnvEnvironment string = "HA_ETA_NODE_ENVIRONMENT"
	EnvConfig      string = "HA_ETA_NODE_CONFIG"
	EnvName        string = "HA_ETA_NODE_NAME"
	EnvLog         string = "HA_ETA_NODE_LOG"
)

const (
	DefaultEnvironment environment.Environment = environment.EnvProd

	DefaultPathConfig string = "/etc/ha-eta-node/config.yml"
	DefaultPathLog    string = "/var/log/ha-eta-node"

	DefaultDaemonPidfileMode os.FileMode = 0644
	DefaultDaemonWorkdir     string      = ""

	DefaultServerINETHost string = "0.0.0.0"
	DefaultServerINETPort int    = 80
	DefaultServerUNIXAddr string = "/run/ha-eta-node/ha-eta-node.sock"

	DefaultPeriodMemstats time.Duration = 60 * time.Second
)

var (
	DefaultCORSAllowMethods []string = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"}
	DefaultCORSAllowHeaders []string = []string{
		"Authorization", "Content-Type", "Accept",
		"Origin", "User-Agent", "DNT",
		"Cache-Control", "X-Mx-ReqToken", "Keep-Alive",
		"X-Requested-With", "If-Modified-Since", "Content-Disposition",
		"X-Token",
	}
)

const (
	ServiceName string = "HA-ETA-Node"
	Company     string = "©WHEELY"
	Logo        string = `
 ___   ___  ________           ______  _________ ________           ___   __   ______  ______  ______
/__/\ /__/\/_______/\         /_____/\/________//_______/\         /__/\ /__/\/_____/\/_____/\/_____/\
\::\ \\  \ \::: _  \ \  ______\::::_\/\__.::.__\\::: _  \ \  ______\::\_\\  \ \:::_ \ \:::_ \ \::::_\/_
 \::\/_\ .\ \::(_)  \ \/______/\:\/___/\ \::\ \  \::(_)  \ \/______/\:.  -\  \ \:\ \ \ \:\ \ \ \:\/___/\
  \:: ___::\ \:: __  \ \__::::\/\::___\/_ \::\ \  \:: __  \ \__::::\/\:. _    \ \:\ \ \ \:\ \ \ \::___\/_
   \: \ \\::\ \:.\ \  \ \        \:\____/\ \::\ \  \:.\ \  \ \        \. \ -\  \ \:\_\ \ \:\/.:| \:\____/\
    \__\/ \::\/\__\/\__\/         \_____\/  \__\/   \__\/\__\/         \__\/ \__\/\_____\/\____/_/\_____\/
`
)
