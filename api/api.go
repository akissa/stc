package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

var (
	c      = resty.New()
	apiKey string
	target string
)

type StConfig struct {
	Folders []struct {
		ID     string `json:"id"`
		Label  string `json:"label"`
		Paused bool   `json:"paused"`
	} `json:"folders"`

	Devices []struct {
		DeviceID string `json:"deviceID"`
		Name     string `json:"name"`
		Paused   bool   `json:"paused"`
	}
}
type SysConn map[string]struct {
	Connected     bool   `json:"connected"`
	InBytesTotal  uint64 `json:"inBytesTotal"`
	OutBytesTotal uint64 `json:"outBytesTotal"`
}

type SysConnections struct {
	Connections SysConn
}

type SysStatus struct {
	CpuPercent float64 `json:"cpuPercent"`
	MyID       string  `json:"myID"`
	Uptime     int64   `json:"uptime"`
}

type SysVersion struct {
	Version string `json:"version"`
}

type DbStatus struct {
	GlobalBytes uint64 `json:"globalBytes"`
	GlobalFiles uint64 `json:"globalFiles"`
	LocalBytes  uint64 `json:"localBytes"`
	LocalFiles  uint64 `json:"localFiles"`
	NeedBytes   uint64 `json:"needBytes"`
	State       string `json:"state"`
}

type DbCompletion struct {
	Completion float64 `json:"completion"`
	NeedBytes  uint64  `json:"needBytes"`
}

func IgnoreCertErrors() {
	c.SetTransport(&http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}})
}

func SetApiKeyTarget(a, t string) error {
	if a == "" || t == "" {
		return fmt.Errorf("apikey and target must be specified")
	}
	apiKey = a
	target = t
	return nil
}

func GetConfig() (StConfig, error) {
	r, err := c.R().SetHeader("X-API-Key", apiKey).Get(target + "/rest/config")
	if err != nil {
		return StConfig{}, err
	}
	if r.IsError() {
		return StConfig{}, fmt.Errorf(r.Status())
	}

	cfg := StConfig{}
	err = json.Unmarshal(r.Body(), &cfg)
	if err != nil {
		return StConfig{}, err
	}

	return cfg, nil
}

func GetFolderStatus(f string) (DbStatus, error) {
	r, err := c.R().SetHeader("X-API-Key", apiKey).SetQueryString("folder=" + f).Get(target + "/rest/db/status")
	if err != nil {
		return DbStatus{}, err
	}
	if r.IsError() {
		return DbStatus{}, fmt.Errorf(r.Status())
	}

	dbs := DbStatus{}
	err = json.Unmarshal(r.Body(), &dbs)
	if err != nil {
		return DbStatus{}, err
	}

	return dbs, nil
}

func GetCompletion(d string) (DbCompletion, error) {
	r, err := c.R().SetHeader("X-API-Key", apiKey).SetQueryString("device=" + d).Get(target + "/rest/db/completion")
	if err != nil {
		return DbCompletion{}, err
	}
	if r.IsError() {
		return DbCompletion{}, fmt.Errorf(r.Status())
	}

	dbc := DbCompletion{}
	err = json.Unmarshal(r.Body(), &dbc)
	if err != nil {
		return DbCompletion{}, err
	}

	return dbc, nil
}

func GetConnection() (SysConn, error) {
	r, err := c.R().SetHeader("X-API-Key", apiKey).Get(target + "/rest/system/connections")
	if err != nil {
		return nil, err
	}
	if r.IsError() {
		return nil, fmt.Errorf(r.Status())
	}

	co := SysConnections{}
	err = json.Unmarshal(r.Body(), &co)
	if err != nil {
		return nil, err
	}

	return co.Connections, nil
}

func GetSysStatus() (SysStatus, error) {
	r, err := c.R().SetHeader("X-API-Key", apiKey).Get(target + "/rest/system/status")
	if err != nil {
		return SysStatus{}, err
	}
	if r.IsError() {
		return SysStatus{}, fmt.Errorf(r.Status())
	}

	st := SysStatus{}
	err = json.Unmarshal(r.Body(), &st)
	if err != nil {
		return SysStatus{}, err
	}

	return st, nil
}

func GetSysVersion() (SysVersion, error) {
	r, err := c.R().SetHeader("X-API-Key", apiKey).Get(target + "/rest/system/version")
	if err != nil {
		return SysVersion{}, err
	}
	if r.IsError() {
		return SysVersion{}, fmt.Errorf(r.Status())
	}

	ve := SysVersion{}
	err = json.Unmarshal(r.Body(), &ve)
	if err != nil {
		return SysVersion{}, err
	}

	return ve, nil
}