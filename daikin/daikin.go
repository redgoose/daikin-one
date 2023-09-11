package daikin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type Daikin struct {
	ApiKey          string
	IntegratorToken string
	Email           string
}

type Token struct {
	AccessToken          string `json:"accessToken"`
	AccessTokenExpiresIn int    `json:"accessTokenExpiresIn"`
	TokenType            string `json:"tokenType"`
}

type Locations []Location

type Location struct {
	LocationName string   `json:"locationName"`
	Devices      []Device `json:"devices"`
}

type Device struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	Model           string `json:"model"`
	FirmwareVersion string `json:"firmwareVersion"`
}

type DeviceInfo struct {
	CoolSetpoint           float32 `json:"coolSetpoint"`
	HeatSetpoint           float32 `json:"heatSetpoint"`
	FanCirculateSpeed      int     `json:"fanCirculateSpeed"`
	EquipmentStatus        int     `json:"equipmentStatus"`
	HumOutdoor             int     `json:"humOutdoor"`
	TempIndoor             float32 `json:"tempIndoor"`
	SetpointDelta          float32 `json:"setpointDelta"`
	EquipmentCommunication int     `json:"equipmentCommunication"`
	ModeEmHeatAvailable    bool    `json:"modeEmHeatAvailable"`
	GeofencingEnabled      bool    `json:"geofencingEnabled"`
	ScheduleEnabled        bool    `json:"scheduleEnabled"`
	HumIndoor              int     `json:"humIndoor"`
	ModeLimit              int     `json:"modeLimit"`
	SetpointMinimum        float32 `json:"setpointMinimum"`
	Fan                    int     `json:"fan"`
	TempOutdoor            float32 `json:"tempOutdoor"`
	Mode                   int     `json:"mode"`
	SetpointMaximum        float32 `json:"setpointMaximum"`
}

type ModeSetpointOptions struct {
	Mode         int     `json:"mode"`
	HeatSetpoint float32 `json:"heatSetpoint"`
	CoolSetpoint float32 `json:"coolSetpoint"`
}

var httpClient = &http.Client{Timeout: 10 * time.Second}
var urlBase string = "https://integrator-api.daikinskyport.com"

func New(apiKey string, integratorToken string, email string) *Daikin {
	d := &Daikin{
		ApiKey:          apiKey,
		IntegratorToken: integratorToken,
		Email:           email,
	}
	return d
}

func (d *Daikin) getToken() string {
	body := []byte(`{
		"email": "` + d.Email + `",
		"integratorToken": "` + d.IntegratorToken + `"
	}`)

	r, err := http.NewRequest("POST", urlBase+"/v1/token", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	r.Header.Add("content-type", "application/json")
	r.Header.Add("x-api-key", d.ApiKey)

	res, err := httpClient.Do(r)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		panic(res.Status)
	}

	token := &Token{}
	derr := json.NewDecoder(res.Body).Decode(token)
	if derr != nil {
		panic(derr)
	}

	return token.AccessToken
}

func (d *Daikin) ListDevices() Locations {
	r, err := http.NewRequest("GET", urlBase+"/v1/devices", nil)
	if err != nil {
		panic(err)
	}

	r.Header.Add("content-type", "application/json")
	r.Header.Add("x-api-key", d.ApiKey)
	r.Header.Add("Authorization", "Bearer "+d.getToken())

	res, err := httpClient.Do(r)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		panic(res.Status)
	}

	var locations Locations
	derr := json.NewDecoder(res.Body).Decode(&locations)
	if derr != nil {
		panic(derr)
	}

	return locations
}

func (d *Daikin) GetDeviceInfo(deviceId string) DeviceInfo {
	r, err := http.NewRequest("GET", urlBase+"/v1/devices/"+deviceId, nil)
	if err != nil {
		panic(err)
	}

	r.Header.Add("content-type", "application/json")
	r.Header.Add("x-api-key", d.ApiKey)
	r.Header.Add("Authorization", "Bearer "+d.getToken())

	res, err := httpClient.Do(r)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		panic(res.Status)
	}

	var deviceInfo DeviceInfo
	derr := json.NewDecoder(res.Body).Decode(&deviceInfo)
	if derr != nil {
		panic(derr)
	}

	return deviceInfo
}

func (d *Daikin) UpdateModeSetpoint(deviceId string, options ModeSetpointOptions) {

	body, err := json.Marshal(options)
	if err != nil {
		panic(err)
	}

	r, err := http.NewRequest("PUT", urlBase+"/v1/devices/"+deviceId+"/msp", bytes.NewBuffer([]byte(body)))
	if err != nil {
		panic(err)
	}

	r.Header.Add("content-type", "application/json")
	r.Header.Add("x-api-key", d.ApiKey)
	r.Header.Add("Authorization", "Bearer "+d.getToken())

	res, err := httpClient.Do(r)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		panic(res.Status)
	}
}
