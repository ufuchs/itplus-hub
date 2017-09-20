package app

import "testing"

func TestConfigService(t *testing.T) {

	BaseDir = "../"

	dataPort := 8080

	svc := NewConfigService().RetrieveAll()
	if svc.LastErr != nil {
		t.Errorf("NewConfigService() throws: %v", svc.LastErr)
	}

	actual := DataPort

	if actual != dataPort {
		t.Errorf("ConfigFactory(): expected '%v', actual '%v'",
			dataPort, actual)
	}

}
