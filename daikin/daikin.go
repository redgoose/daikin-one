package daikin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

const (
	EquipmentStatusCool     = 1
	EquipmentStatusOvercool = 2
	EquipmentStatusHeat     = 3
	EquipmentStatusFan      = 4
	EquipmentStatusIdle     = 5
)

var httpClient = &http.Client{Timeout: 10 * time.Second}
var urlBase string = "https://integrator-api.daikinskyport.com"

func New(apiKey string, integratorToken string, email string) *Daikin {
	d := Daikin{
		ApiKey:          apiKey,
		IntegratorToken: integratorToken,
		Email:           email,
	}
	return &d
}

func (d *Daikin) getToken() (string, error) {
	body := []byte(`{
		"email": "` + d.Email + `",
		"integratorToken": "` + d.IntegratorToken + `"
	}`)

	r, err := http.NewRequest("POST", urlBase+"/v1/token", bytes.NewBuffer(body))
	if err != nil {
		return "", errors.New("http.NewRequest failed")
	}

	r.Header.Add("content-type", "application/json")
	r.Header.Add("x-api-key", d.ApiKey)

	res, err := httpClient.Do(r)
	if err != nil {
		return "", errors.New("http request failed")
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token request returned a non-success response: %s", res.Status)
	}

	token := &Token{}
	derr := json.NewDecoder(res.Body).Decode(token)
	if derr != nil {
		return "", errors.New("json decode failed")
	}

	return token.AccessToken, nil
}

func (d *Daikin) ListDevices() (*Locations, error) {
	r, err := http.NewRequest("GET", urlBase+"/v1/devices", nil)
	if err != nil {
		return nil, errors.New("http.NewRequest failed")
	}

	r.Header.Add("content-type", "application/json")
	r.Header.Add("x-api-key", d.ApiKey)

	token, err := d.getToken()
	if err != nil {
		return nil, errors.New("getToken did not return a token")
	}

	r.Header.Add("Authorization", "Bearer "+token)

	res, err := httpClient.Do(r)
	if err != nil {
		return nil, errors.New("http request failed")
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list devices request returned a non-success response: %s", res.Status)
	}

	var locations Locations
	derr := json.NewDecoder(res.Body).Decode(&locations)
	if derr != nil {
		return nil, errors.New("json decode failed")
	}

	return &locations, nil
}

func (d *Daikin) GetDeviceInfo(deviceId string) (*DeviceInfo, error) {
	r, err := http.NewRequest("GET", urlBase+"/v1/devices/"+deviceId, nil)
	if err != nil {
		return nil, errors.New("http.NewRequest failed")
	}

	r.Header.Add("content-type", "application/json")
	r.Header.Add("x-api-key", d.ApiKey)

	token, err := d.getToken()
	if err != nil {
		return nil, errors.New("getToken did not return a token")
	}

	r.Header.Add("Authorization", "Bearer "+token)

	res, err := httpClient.Do(r)
	if err != nil {
		return nil, errors.New("http request failed")
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get device info request returned a non-success response: %s", res.Status)
	}

	var deviceInfo DeviceInfo
	derr := json.NewDecoder(res.Body).Decode(&deviceInfo)
	if derr != nil {
		return nil, errors.New("json decode failed")
	}

	return &deviceInfo, nil
}

func (d *Daikin) UpdateModeSetpoint(deviceId string, options ModeSetpointOptions) error {

	body, err := json.Marshal(options)
	if err != nil {
		return errors.New("json.Marshal failed")
	}

	r, err := http.NewRequest("PUT", urlBase+"/v1/devices/"+deviceId+"/msp", bytes.NewBuffer([]byte(body)))
	if err != nil {
		return errors.New("http.NewRequest failed")
	}

	r.Header.Add("content-type", "application/json")
	r.Header.Add("x-api-key", d.ApiKey)

	token, err := d.getToken()
	if err != nil {
		return errors.New("getToken did not return a token")
	}

	r.Header.Add("Authorization", "Bearer "+token)

	res, err := httpClient.Do(r)
	if err != nil {
		return errors.New("http request failed")
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("update mode setpoint request returned a non-success response: %s", res.Status)
	}

	return nil
}
