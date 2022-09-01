package interfaces

import (
	"fmt"
	loggerInterfaces "github.com/fajarardiyanto/flt-go-logger/interfaces"
	"github.com/shirou/gopsutil/v3/host"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type (
	// ContextKeyType is a private struct that is used for storing values in net.Context
	ContextKeyType struct{}

	// ParamsMapType is a private type that is used to store route params
	ParamsMapType map[string]string
)

var (
	// ContextKey is the key that is used to store values in net.Context for each request
	ContextKey = ContextKeyType{}
	// GoVersion is golang version run
	GoVersion string
	// OSBuildName version os running
	OSBuildName string
	//BuildDate date application running
	BuildDate string
)

// GetAllParams returns all route params stored in http.Request.
func GetAllParams(r *http.Request) ParamsMapType {
	if values, ok := r.Context().Value(ContextKey).(ParamsMapType); ok {
		return values
	}
	return nil
}

func ResolveAddress(addr []string, logger loggerInterfaces.Logger) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			logger.Info("Environment variable PORT=\"%s\"", port)
			return ":" + port
		}
		logger.Info("Environment variable PORT is undefined. Using port :8080 by default")
		return ":8080"
	case 1:
		var port string
		if strings.Contains(addr[0], ":") {
			port = addr[0]
		} else {
			port = ":" + addr[0]
		}

		return port
	default:
		panic("too many parameters")
	}
}

func ShowVersion(ver string) {
	const show = `[INFO] Server is starting
├── Version	         : %s
├── GoVersion	     : %s
├── Compiler	     : %s
├── BuildTime	     : %s
└── Path is register :`
	version := fmt.Sprintf(show, ver, GoVersion, OSBuildName, BuildDate)
	fmt.Println(version)
}

func init() {
	if BuildDate == "" {
		BuildDate = time.Now().Format("Mon Jan 02 15:04:05 MST 2006")
	}

	if GoVersion == "" {
		GoVersion = fmt.Sprintf("%s %s/%s", strings.Replace(runtime.Version(), "go", "", -1), runtime.GOOS, runtime.GOARCH)
	}

	if len(OSBuildName) == 0 {
		_, err := os.Stat("/etc/os-release")
		if err == nil {
			if dat, err := exec.Command("sh", "-c", `awk -F'=' '/PRETTY_NAME/ {print $2}' /etc/os-release | sed "s/\"//g"`).CombinedOutput(); err == nil {
				OSBuildName = strings.TrimSpace(string(dat))
			} else {
				SetOSVersion()
			}
		} else {
			SetOSVersion()
		}
	}
}

func SetOSVersion() {
	if info, err := host.Info(); err == nil {
		OSBuildName = fmt.Sprintf("%s (%s)",
			info.PlatformFamily, info.PlatformVersion)
	}
}
