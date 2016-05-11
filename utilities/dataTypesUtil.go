package utilities

import (
	"strconv"
	"time"
	"log"
)

func ToInt(value string) (error, int) {
	intInput, err := strconv.ParseInt(value, 10, 0)
	if err != nil {
		return err, 0;
	}
	return nil, int(intInput)
}
// Extracts time from string with format "dd/MM/yyyy"
// stdZeroDay => "02"
// stdZeroMonth => "01"
// stdLongYear => "2006"
func GetTimeNowFromString(dateString string) (time.Time, error) {
	uiDate, err := time.Parse("02/01/2006", dateString)
	if err != nil {
		log.Println("Error converting string to date from string :", dateString)
		return uiDate, err
	}
	now := time.Now().UTC()
	uiDateNowTime := time.Date(uiDate.Year(), uiDate.Month(), uiDate.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.UTC)
	return uiDateNowTime, err
}