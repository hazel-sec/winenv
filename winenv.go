// Package winenv collect Windows environment informations at local
package winenv

import (
	"fmt"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// WindowsEnvInfo to get
type WindowsEnvInfo struct {
	Product      string `json:"product"`
	Version      string `json:"version"`
	Release      string `json:"release"`
	Build        string `json:"build"`
	ServicePack  string `json:"service_pack"`
	Architecture string `json:"architecture"`
}

// NewWindowsEnvInfo from Environment variable and Registry
func NewWindowsEnvInfo() (envinfo WindowsEnvInfo, err error) {

	k, err := registry.OpenKey(
		registry.LOCAL_MACHINE,
		`SOFTWARE\Microsoft\Windows NT\CurrentVersion`,
		registry.QUERY_VALUE|registry.WOW64_64KEY,
	)
	if err != nil {
		return WindowsEnvInfo{}, err
	}
	defer k.Close()

	product, _, err := k.GetStringValue("ProductName")
	if err != nil {
		return WindowsEnvInfo{}, err
	}

	// keys must be exists
	version, _, err := k.GetStringValue("CurrentVersion")
	if err != nil {
		return WindowsEnvInfo{}, err
	}

	build, _, err := k.GetStringValue("CurrentBuildNumber")
	if err != nil {
		return WindowsEnvInfo{}, err
	}

	var arch string

	buildLabEx, _, err := k.GetStringValue("BuildLabEx")
	if err != nil {
		return WindowsEnvInfo{}, err
	}
	switch {
	case strings.Contains(buildLabEx, "amd64"):
		arch = "amd64"
	case strings.Contains(buildLabEx, "x86"):
		arch = "x86"
	case len(buildLabEx) == 0:
		arch = "unknown"
	default:
		items := strings.Split(buildLabEx, ".")
		arch = items[2]
	}

	// keys maybe not exisits
	release, _, _ := k.GetStringValue("ReleaseId")
	servicepack, _, _ := k.GetStringValue("CSDVersion")

	return WindowsEnvInfo{
		Version:      version,
		Product:      product,
		Release:      release,
		Build:        build,
		ServicePack:  servicepack,
		Architecture: arch,
	}, nil
}

// String for WindowsEnvInfo
func (e WindowsEnvInfo) String() (str string) {

	str = e.Product

	if e.ServicePack != "" {
		sp := strings.Replace(e.ServicePack, "Service Pack ", "SP", 1)
		str = fmt.Sprintf("%s %s", str, sp)
	}

	str = fmt.Sprintf("%s (%s) [ver. %s] [build %s]", str, e.Architecture, e.Version, e.Build)

	if e.Release != "" {
		str = fmt.Sprintf("%s [release %s]", str, e.Release)
	}

	return
}
