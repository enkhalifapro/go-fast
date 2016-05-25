package services

import (
	"testing"

	"github.com/enkhalifapro/go-fast/mocks"
	"github.com/enkhalifapro/go-fast/models"
	"github.com/enkhalifapro/go-fast/utilities"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func setupClearUserCollection(configUtil utilities.IConfigUtil) {
	session, _ := mgo.Dial(configUtil.GetConfig("dbUri"))
	defer session.Close()
	session.SetSafe(&mgo.Safe{})
	collection := session.DB(configUtil.GetConfig("dbName")).C("users")
	// clear efficiencies collection
	_, err := collection.RemoveAll(bson.M{})
	if err != nil {
		panic(err)
	}
}

func setupAddfakeUser(configUtil utilities.IConfigUtil, token string) {
	session, _ := mgo.Dial(configUtil.GetConfig("dbUri"))
	defer session.Close()
	session.SetSafe(&mgo.Safe{})
	collection := session.DB(configUtil.GetConfig("dbName")).C("users")
	// clear efficiencies collection
	err := collection.Insert(models.User{Id: bson.NewObjectId(), UserName: "ayman"})
	if err != nil {
		panic(err)
	}
}

func TestFindById(t *testing.T) {
	Convey("Given I have no users", t, func() {
		configUtil := mocks.ConfigUtilMock{}
		setupClearUserCollection(configUtil)
		cryptUtil := utilities.NewCryptUtil()
		sessionService := NewSessionService(configUtil, cryptUtil)
		Convey("When get unknown session", func() {
			validateResult := sessionService.Valid("unknown")
			Convey("validation result should be false", func() {
				So(validateResult, ShouldEqual, false)
			})
		})
	})
}
