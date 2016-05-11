package mocks

type ConfigUtilMock struct {

}

func (r ConfigUtilMock) GetConfig(key string) string {
	if key == "dbName" {
		return "kn-test"
	}
	if key == "dbUri" {
		return "mongodb://root:123456@162.243.1.78:19860"
	}
	if key == "postmarkKey" {
		return "232f02b1-a979-43ed-9e1e-e7bbe560cd6e"
	}
	if key == "domain" {
		return ""
	}
	return ""
}

func (r ConfigUtilMock) GetConfigAsInt(key string) int {
	return 0
}