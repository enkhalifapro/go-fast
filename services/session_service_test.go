package services

import (
	"testing"
	"time"

	"github.com/enkhalifapro/go-fast/mocks"
	"github.com/enkhalifapro/go-fast/models"
	"github.com/enkhalifapro/go-fast/utilities"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func setupClearSessionCollection(configUtil utilities.IConfigUtil) {
	session, _ := mgo.Dial(configUtil.GetConfig("dbUri"))
	defer session.Close()
	session.SetSafe(&mgo.Safe{})
	collection := session.DB(configUtil.GetConfig("dbName")).C("sessions")
	// clear efficiencies collection
	_, err := collection.RemoveAll(bson.M{})
	if err != nil {
		panic(err)
	}
}

func setupAddfakeUserSession(configUtil utilities.IConfigUtil, token string) {
	session, _ := mgo.Dial(configUtil.GetConfig("dbUri"))
	defer session.Close()
	session.SetSafe(&mgo.Safe{})
	collection := session.DB(configUtil.GetConfig("dbName")).C("sessions")
	// clear efficiencies collection
	err := collection.Insert(models.Session{Id: bson.NewObjectId(), Token: token, ExpiryDate: time.Now().UTC().Add(time.Hour * 24)})
	if err != nil {
		panic(err)
	}
}

func setupGetSessionsCount(configUtil utilities.IConfigUtil) (int, error) {
	session, _ := mgo.Dial(configUtil.GetConfig("dbUri"))
	defer session.Close()
	session.SetSafe(&mgo.Safe{})
	collection := session.DB(configUtil.GetConfig("dbName")).C("sessions")
	// clear efficiencies collection
	count, err := collection.Count()
	if err != nil {
		panic(err)
	}
	return count, err
}

func TestValid(t *testing.T) {
	Convey("Given I have no user sessions", t, func() {
		configUtil := mocks.ConfigUtilMock{}
		setupClearSessionCollection(configUtil)
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

func TestFind(t *testing.T) {
	Convey("Given I have one user session 'session-1-token'", t, func() {
		configUtil := mocks.ConfigUtilMock{}
		setupClearSessionCollection(configUtil)
		setupAddfakeUserSession(configUtil, "session-1-token")
		cryptUtil := utilities.NewCryptUtil()
		sessionService := NewSessionService(configUtil, cryptUtil)
		Convey("When query 'session-1-token'", func() {
			err, session := sessionService.Find("session-1-token")
			Convey("err should be nil", func() {
				So(err, ShouldEqual, nil)
			})
			Convey("session token should equal 'session-1-token'", func() {
				So(session.Token, ShouldEqual, "session-1-token")
			})
		})
	})
}

func TestLogout(t *testing.T) {
	Convey("Given I have one user session 'session-1-token'", t, func() {
		configUtil := mocks.ConfigUtilMock{}
		setupClearSessionCollection(configUtil)
		setupAddfakeUserSession(configUtil, "session-1-token")
		cryptUtil := utilities.NewCryptUtil()
		sessionService := NewSessionService(configUtil, cryptUtil)
		Convey("When logout 'session-1-token'", func() {
			err := sessionService.Logout("session-1-token")
			Convey("err should be nil", func() {
				So(err, ShouldEqual, nil)
			})
			Convey("when validate 'session-1-token' status should be false", func() {
				status := sessionService.Valid("session-1-token")
				So(status, ShouldEqual, false)
			})
		})
	})
}
