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

type IRoleService interface {
	FindOne(query *bson.M) (error, *models.Role)
	Find(query *bson.M) (err error, roles []*models.Role)
	Insert(updaterId string, role *models.Role) error
	UpdateByName(updaterId string, roleName string, newRole *models.Role) error
	DeleteByName(name string) error
}

type RoleService struct {
	dbName            string
	uri               string
	dialMongoWithInfo string
	collectionName    string
	sessionService    *SessionService
	configUtil        utilities.IConfigUtil
	slugUtil          utilities.ISlugUtil
}

func NewRoleService(configUtil utilities.IConfigUtil, slugUtil utilities.ISlugUtil) *RoleService {
	r := RoleService{}
	r.uri = configUtil.GetConfig("dbUri")
	r.dbName = configUtil.GetConfig("dbName")
	r.dialMongoWithInfo = configUtil.GetConfig("dialMongoWithInfo")
	r.collectionName = "roles"
	r.configUtil = configUtil
	r.slugUtil = slugUtil
	return &r
}

func (r RoleService) newSession() (*mgo.Session, error) {
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

func (r RoleService) FindOne(query *bson.M) (error, *models.Role) {
	session, _ := r.newSession()
	defer session.Close()
	session.SetSafe(&mgo.Safe{})
	collection := session.DB(r.dbName).C(r.collectionName)
	role := models.Role{}
	err := collection.Find(query).One(&role)
	return err, &role
}

func (r RoleService) Find(query *bson.M) (err error, roles []*models.Role) {
	session, _ := r.newSession()
	defer session.Close()
	session.SetSafe(&mgo.Safe{})
	collection := session.DB(r.dbName).C(r.collectionName)
	roles = make([]*models.Role, 0)
	err = collection.Find(query).Sort("order").All(&roles)
	return err, roles
}

func (r RoleService) Insert(updaterId string, role *models.Role) error {
	role.Id = bson.NewObjectId()
	role.CreatedAt = time.Now().UTC()
	role.UpdatedAt = time.Now().UTC()
	role.UpdaterId = updaterId
	role.Slug = r.slugUtil.GetSlug(role.Name)
	session, _ := r.newSession()
	defer session.Close()
	session.SetSafe(&mgo.Safe{})
	collection := session.DB(r.dbName).C(r.collectionName)
	err := collection.Insert(role)
	return err
}

func (r RoleService) UpdateByName(updaterId string, roleName string, newRole *models.Role) error {
	session, _ := r.newSession()
	defer session.Close()
	session.SetSafe(&mgo.Safe{})
	collection := session.DB(r.dbName).C(r.collectionName)
	err := collection.Update(bson.M{"slug": roleName}, bson.M{"$set": bson.M{
		"name":      newRole.Name,
		"slug":      r.slugUtil.GetSlug(newRole.Name),
		"updaterid": updaterId,
		"updatedat": time.Now().UTC()}})
	return err
}

func (r RoleService) DeleteByName(name string) error {
	session, _ := r.newSession()
	defer session.Close()
	session.SetSafe(&mgo.Safe{})
	collection := session.DB(r.dbName).C(r.collectionName)
	err := collection.Remove(bson.M{"slug": name})
	return err
}
