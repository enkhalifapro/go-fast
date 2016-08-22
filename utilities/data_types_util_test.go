package utilities

import (
	"testing"
)

func TestToInt(t *testing.T) {
	// act
	_, intValue := ToInt("996")
	// assert
	if intValue != 996 {
		t.Error("Error Converting string into Integer")
	}
}

/*
func TestGetTimeNowFromString(t *testing.T) {
	// Only pass t into top-level Convey calls
	Convey("Given I have date string 16/05/2016", t, func() {
		strDate := "16/05/2016"
		Convey("When covert to date", func() {

			date,err:= GetTimeNowFromString(strDate)
			if err != nil {
				panic(err)
			}
			Convey("error should be nil", func() {
				So(err, ShouldEqual, nil)
			})
			Convey("month should be may", func() {
				So(date.Month().String(), ShouldEqual, "May")
			})
			Convey("year should be 2016", func() {
				So(date.Year(), ShouldEqual, 2016)
			})
			Convey("day should be 16", func() {
				So(date.Day(), ShouldEqual, 16)
			})
		})
	})
}*/
