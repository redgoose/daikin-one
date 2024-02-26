package charts

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/redgoose/daikin-one/internal/db"
)

func TestConvertTempsCtoF(t *testing.T) {
	// Define a test case with known input and expected output
	testCases := []struct {
		input    []db.PeriodData
		expected []db.PeriodData
	}{
		{
			input: []db.PeriodData{
				// NOTE
				{Period: "Morning", TempIndoor: 20.0, TempOutdoor: 15.0, HumidityIndoor: 50.0, HumidityOutdoor: 60.0, CoolSetpoint: 22.0, HeatSetpoint: 18.0, RunTime: 120},
				{Period: "Afternoon", TempIndoor: 25.0, TempOutdoor: 30.0, HumidityIndoor: 55.0, HumidityOutdoor: 65.0, CoolSetpoint: 24.0, HeatSetpoint: 19.0, RunTime: 180},
			},
			expected: []db.PeriodData{
				{Period: "Morning", TempIndoor: 68.0, TempOutdoor: 59.0, HumidityIndoor: 50.0, HumidityOutdoor: 60.0, CoolSetpoint: 71.6, HeatSetpoint: 64.4, RunTime: 120},
				{Period: "Afternoon", TempIndoor: 77.0, TempOutdoor: 86.0, HumidityIndoor: 55.0, HumidityOutdoor: 65.0, CoolSetpoint: 75.2, HeatSetpoint: 66.2, RunTime: 180},
			},
		},
	}

	for _, tc := range testCases {
		// Convert temperatures from C to F
		converted := convertTempsCtoF(tc.input)

		// Check each field in the converted struct against the expected output
		for i, cv := range converted {

			// Expect the temperature fields to change
			if cv.TempIndoor != tc.expected[i].TempIndoor {
				t.Errorf("Test failed for Period %s: expected TempIndoor = %.2f F, got TempIndoor = %.2f F",
					cv.Period, tc.expected[i].TempIndoor, cv.TempIndoor)
			}

			if cv.TempOutdoor != tc.expected[i].TempOutdoor {
				t.Errorf("Test failed for Period %s: expected TempOutdoor = %.2f F, got TempOutdoor = %.2f F",
					cv.Period, tc.expected[i].TempOutdoor, cv.TempOutdoor)
			}

			if cv.CoolSetpoint != tc.expected[i].CoolSetpoint {
				t.Errorf("Test failed for Period %s: expected CoolSetpoint = %.2f, got CoolSetpoint = %.2f",
					cv.Period, tc.expected[i].CoolSetpoint, cv.CoolSetpoint)
			}

			if cv.HeatSetpoint != tc.expected[i].HeatSetpoint {
				t.Errorf("Test failed for Period %s: expected HeatSetpoint = %.2f, got HeatSetpoint = %.2f",
					cv.Period, tc.expected[i].HeatSetpoint, cv.HeatSetpoint)
			}

			// Expect the non-temperature fields to not change.
			if cv.HumidityIndoor != tc.expected[i].HumidityIndoor {
				t.Errorf("Test failed for Period %s: expected HumidityIndoor = %.2f%%, got HumidityIndoor = %.2f%%",
					cv.Period, tc.expected[i].HumidityIndoor, cv.HumidityIndoor)
			}

			if cv.HumidityOutdoor != tc.expected[i].HumidityOutdoor {
				t.Errorf("Test failed for Period %s: expected HumidityOutdoor = %.2f%%, got HumidityOutdoor = %.2f%%",
					cv.Period, tc.expected[i].HumidityOutdoor, cv.HumidityOutdoor)
			}

			if cv.RunTime != tc.expected[i].RunTime {
				t.Errorf("Test failed for Period %s: expected RunTime = %d, got RunTime = %d",
					cv.Period, tc.expected[i].RunTime, cv.RunTime)
			}

		}
	}
}
