# Daikin One CLI

## Overview

Daikin One is a lightweight CLI wrapper for the [Daikin Open API](https://www.daikinone.com/openapi/index.html) to manage Daikin One devices/thermostats. 

It is augmented with logging and chart generation functionality to provide historical usage insights which is missing from Daikin One devices.

## Quick Start

1. Install [Go](https://golang.org/doc/install)
2. Install `daikin-one`:
```
go install github.com/redgoose/daikin-one@latest
```

3. Copy `.daikin.yml` to your home directory and populate with your integrator token, api key, and email. Refer to the [getting started](https://www.daikinone.com/openapi/documentation/index.html#gettingstarted) section of the Daikin Open API documentation to get those values. 


4. Get a device id for the device you want to manage:

	```
	daikin-one device ls
	```

	This will return a list of devices associated with your account:

		```
		[
		        {
		                "locationName": "Home",
		                "devices": [
		                        {
		                                "id": "0a000000-0000-00aa-00a0-0a00aa0a00aa",
		                                "name": "Main Room",
		                                "model": "ONEPLUS",
		                                "firmwareVersion": "3.2.19"
		                        }
		                ]
		        }
		]
		```
		
		

5. With a device id, you can start managing your device. For example, to retrieve device configuration and state values use:

	```
	daikin-one device info --device-id <device-id>
	```
	
	This will output something like:
	
	```
	{
	        "coolSetpoint": 22,
	        "heatSetpoint": 20,
	        "fanCirculateSpeed": 0,
	        "equipmentStatus": 5,
	        "humOutdoor": 73,
	        "tempIndoor": 21.8,
	        "setpointDelta": 2,
	        "equipmentCommunication": 0,
	        "modeEmHeatAvailable": true,
	        "geofencingEnabled": true,
	        "scheduleEnabled": true,
	        "humIndoor": 57,
	        "modeLimit": 1,
	        "setpointMinimum": 10,
	        "fan": 0,
	        "tempOutdoor": 21,
	        "mode": 2,
	        "setpointMaximum": 32
	}
	```

	Refer to the Daikin Open API [data definitions](https://www.daikinone.com/openapi/documentation/index.html#datadefinitions) for information on fields and possible values.

Run `daikin-one -h` and `daikin-one device -h` for a full list of available commands.

## License

MIT © redgoose, see [LICENSE](https://github.com/regdoose/daikin-one/blob/master/LICENSE) for details.
