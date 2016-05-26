package services

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"time"

	"github.com/enkhalifapro/go-fast/models"
	"github.com/enkhalifapro/go-fast/utilities"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type SessionService struct {
	dbName            string
	uri               string
	dialMongoWithInfo string
	collectionName    string
	userService       *UserService
	configUtil        utilities.IConfigUtil
	cryptUtil         utilities.ICryptUtil
}

func NewSessionService(configUtil utilities.IConfigUtil, cryptUtil utilities.ICryptUtil) *SessionService {
	r := SessionService{}
	r.uri = configUtil.GetConfig("dbUri")
	r.dbName = configUtil.GetConfig("dbName")
	r.dialMongoWithInfo = configUtil.GetConfig("dialMongoWithInfo")
	r.collectionName = "sessions"
	r.cryptUtil = cryptUtil
	r.configUtil = configUtil
	return &r
}

func (r SessionService) newSession() (*mgo.Session, error) {
	if r.dialMongoWithInfo == "true" {
		tlsConfig := &tls.Config{}
		roots := x509.NewCertPool()
		path := r.configUtil.GetConfig("path")
		if ca, err := ioutil.ReadFile(path + "/ssh/mongo.pem"); err == nil {
			roots.AppendCertsFromPEM(ca)
		}
		tlsConfig.RootCAs = roots

		dialInfo, _ := mgo.ParseURL(r.uri)
		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			fmt.Println("try connect")
			conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
			return conn, err
		}
		//Here is the session you are looking for. Up to you from here ;)
		return mgo.DialWithInfo(dialInfo)
	}
	return mgo.Dial(r.uri)
}

func (r SessionService) Insert(userSession *models.Session) error {
	userSession.Id = bson.NewObjectId()
	randomStr := r.cryptUtil.RandomString(100)
	userSession.Token = r.cryptUtil.Encrypt(randomStr)
	session, _ := r.newSession()
	defer session.Close()
	session.SetSafe(&mgo.Safe{})
	collection := session.DB(r.dbName).C(r.collectionName)
	return collection.Insert(userSession)
}

func (r SessionService) Valid(sessionToken string) bool {
	session, _ := r.newSession()
	defer session.Close()
	session.SetSafe(&mgo.Safe{})
	collection := session.DB(r.dbName).C(r.collectionName)
	count, _ := collection.Find(&bson.M{"token": sessionToken,
		"expirydate": bson.M{"$gte": time.Now().UTC()}}).Count()
	return count > 0
}

func (r SessionService) Find(sessionToken string) (error, *models.Session) {
	session, _ := r.newSession()
	defer session.Close()
	session.SetSafe(&mgo.Safe{})
	collection := session.DB(r.dbName).C(r.collectionName)
	userSession := models.Session{}
	err := collection.Find(&bson.M{"token": sessionToken,
		"expirydate": bson.M{"$gte": time.Now().UTC()}}).One(&userSession)
	return err, &userSession
}

func (r SessionService) Logout(token string) error {
	session, _ := r.newSession()
	defer session.Close()
	session.SetSafe(&mgo.Safe{})
	collection := session.DB(r.dbName).C(r.collectionName)
	err := collection.Update(bson.M{"token": token}, bson.M{"$set": bson.M{"expirydate": time.Now().UTC()}})
	return err
}
