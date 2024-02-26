-- Add new fields:
    -- OutdoorHeatDemand
    -- OutdoorCoolDemand
    -- IndoorFanActual
    -- IndoorHeatActual

ALTER TABLE daikin
ADD outdoor_heat REAL;

ALTER TABLE daikin
ADD outdoor_cool REAL;

ALTER TABLE daikin
ADD indoor_fan REAL;

ALTER TABLE daikin
ADD indoor_heat REAL;