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
	"github.com/enkhalifapro/go-fast/viewModels"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type IUserService interface {
	FindById(id string) (error, *models.User)
	FindOne(query *bson.M) (error, *models.User)
	Find(query *models.User) (users []*models.User)
	Insert(user *models.User) error
	Login(loginViewModel *viewModels.LoginViewModel) bool
	CurrentUser(sessionToken string) (error, *models.User)
	VerifyEmail(token string) error
	ResendVerifyEmail(user *models.User) error
	UpdateById(updaterId string, userId string, newUser *viewModels.UpdateUserViewModel) error
	ChangePassword(updaterId string, userId string, newPassword string) error
	SendPasswordResetEmail(user *models.User) error
	ValidatePasswordResetToken(token string) (err error, userId string)
	Logout(token string) error
}

type UserService struct {
	dbName            string
	uri               string
	dialMongoWithInfo string
	collectionName    string
	sessionService    *SessionService
	cryptUtil         utilities.ICryptUtil
	configUtil        utilities.IConfigUtil
	mailUtil          utilities.IMailUtil
	slugUtil          utilities.ISlugUtil
}

func NewUserService(configUtil utilities.IConfigUtil, cryptUtil utilities.ICryptUtil) *UserService {
	r := UserService{}
	r.uri = configUtil.GetConfig("dbUri")
	r.dbName = configUtil.GetConfig("dbName")
	r.dialMongoWithInfo = configUtil.GetConfig("dialMongoWithInfo")
	r.collectionName = "users"
	r.cryptUtil = cryptUtil
	r.configUtil = configUtil
	r.sessionService = NewSessionService(configUtil, cryptUtil)
	r.mailUtil = utilities.NewMailUtil(configUtil)
	r.slugUtil = utilities.NewSlugUtil()
	return &r
}

func (r UserService) populateRole(user *models.User, session *mgo.Session) {
	roleCollection := session.DB(r.dbName).C("roles")
	roleCollection.FindId(user.RoleId).One(&user.Role)
}

func (r UserService) FindById(id string) (error, *models.User) {
	session, _ := r.newSession()
	defer session.Close()
	session.SetSafe(&mgo.Safe{})
	collection := session.DB(r.dbName).C(r.collectionName)
	user := models.User{}
	err := collection.FindId(bson.ObjectIdHex(id)).One(&user)
	if err != nil {
		return err, nil
	}
	r.populateRole(&user, session)
	return nil, &user
}

func (r UserService) FindOne(query *bson.M) (error, *models.User) {
	session, _ := r.newSession()
	defer session.Close()
	session.SetSafe(&mgo.Safe{})
	collection := session.DB(r.dbName).C(r.collectionName)
	user := models.User{}
	err := collection.Find(query).One(&user)
	return err, &user
}

func (r UserService) Find(query *models.User) (users []*models.User) {
	session, _ := r.newSession()
	defer session.Close()
	session.SetSafe(&mgo.Safe{})
	collection := session.DB(r.dbName).C(r.collectionName)
	bsonQuery, _ := bson.Marshal(query)
	find := collection.Find(bsonQuery).Iter()
	user := models.User{}
	for find.Next(&user) {
		r.populateRole(&user, session)
		users = append(users, &user)
	}
	return users
}

func (r UserService) newSession() (*mgo.Session, error) {
	fmt.Println("connection is ")
	fmt.Println(r.uri)
	fmt.Println("is dial with info")
	fmt.Println(r.dialMongoWithInfo)
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
	fmt.Println("without ssl")
	return mgo.Dial(r.uri)
}

func (r UserService) Insert(user *models.User) error {
	user.Id = bson.NewObjectId()
	user.Slug = r.slugUtil.GetSlug(user.UserName)
	user.Password = r.cryptUtil.Bcrypt(user.Password)
	user.EmailVerified = false
	user.VerifyToken = r.cryptUtil.NewEncryptedToken()

	session, _ := r.newSession()
	defer session.Close()
	session.SetSafe(&mgo.Safe{})
	collection := session.DB(r.dbName).C(r.collectionName)
	err := collection.Insert(user)
	// send verify email
	globalVars := map[string]interface{}{"FNAME": user.FirstName, "ACTIVATE_ACCOUNT": "google.com/verify/registration?token=" + user.VerifyToken + "&email=" + user.Email}

	_, mailErr := r.mailUtil.SendTemplate(user.Email, "john@curtisdigital.com", "John Curtis", "Verify your e-mail", globalVars, 606061)
	if mailErr != nil {
		fmt.Println("verify email send error")
		fmt.Println(mailErr)
	}
	return err
}

func (r UserService) ResendVerifyEmail(user *models.User) error {
	// send verify email
	globalVars := map[string]interface{}{"FNAME": user.FirstName, "ACTIVATE_ACCOUNT": "google.com/verify/registration?token=" + user.VerifyToken + "&email=" + user.Email}

	_, mailErr := r.mailUtil.SendTemplate(user.Email, "john@curtisdigital.com", "John Curtis", "Verify your e-mail", globalVars, 606061)
	if mailErr != nil {
		fmt.Println("verify email send error")
		fmt.Println(mailErr)
	}
	return mailErr
}

func (r UserService) NewPasswordResetToken(user *models.User) (error, *models.PasswordResetToken) {
	passwordResetToken := models.PasswordResetToken{}
	passwordResetToken.Id = bson.NewObjectId()
	passwordResetToken.UserId = user.Id.Hex()
	passwordResetToken.Token = r.cryptUtil.NewEncryptedToken()
	passwordResetToken.IsValid = true
	passwordResetToken.CreatedAt = time.Now().UTC()
	passwordResetToken.UpdatedAt = time.Now().UTC()
	passwordResetToken.UpdaterId = user.Id.Hex()
	session, _ := r.newSession()
	defer session.Close()
	session.SetSafe(&mgo.Safe{})
	collection := session.DB(r.dbName).C("passwordresettokens")

	// expire all old token related to user
	_, err := collection.UpdateAll(bson.M{"userid": user.Id.Hex()}, bson.M{"$set": bson.M{
		"isvalid": false}})

	if err != nil && err.Error() != "not found" {
		return err, nil
	}

	// insert new token
	err = collection.Insert(passwordResetToken)
	return err, &passwordResetToken
}

func (r UserService) SendPasswordResetEmail(user *models.User) error {
	err, token := r.NewPasswordResetToken(user)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// send verify email
	globalVars := map[string]interface{}{"FNAME": user.FirstName, "RESET_PASSWORD": "http://google.com/password/reset?token=" + token.Token + "&email=" + user.Email}

	_, mailErr := r.mailUtil.SendTemplate(user.Email, "john@curtisdigital.com", "John Curtis", "Reset your password", globalVars, 613062)
	if mailErr != nil {
		fmt.Println("verify email send error")
		fmt.Println(mailErr)
	}
	return mailErr
}

func (r UserService) ValidatePasswordResetToken(token string) (err error, userId string) {
	session, _ := r.newSession()
	defer session.Close()
	session.SetSafe(&mgo.Safe{})
	collection := session.DB(r.dbName).C("passwordresettokens")

	passwordResetToken := models.PasswordResetToken{}
	err = collection.Find(bson.M{"token": token, "isvalid": true}).One(&passwordResetToken)
	if err != nil {
		return err, ""
	}
	err = collection.RemoveId(passwordResetToken.Id)
	if err != nil {
		return err, ""
	}
	return err, passwordResetToken.UserId
}

func (r UserService) Login(loginViewModel *viewModels.LoginViewModel) bool {
	err, user := r.FindOne(&bson.M{
		"email": loginViewModel.Email})
	if err != nil {
		fmt.Println(err)
		return false
	}
	isValidPassowrd := r.cryptUtil.CompareHashAndPassword(user.Password, loginViewModel.Password)
	if isValidPassowrd != true {
		fmt.Println("it is not valid pass")
		return false
	}
	userSession := models.Session{}
	userSession.UserId = user.Id.Hex()
	userSession.ExpiryDate = time.Now().UTC().Add(time.Hour * 24)
	r.sessionService.Insert(&userSession)
	loginViewModel.Token = userSession.Token
	loginViewModel.UserId = user.Id.Hex()
	loginViewModel.UserName = user.UserName
	loginViewModel.FirstName = user.FirstName
	loginViewModel.LastName = user.LastName
	loginViewModel.Image = user.Image
	return true
}

func (r UserService) CurrentUser(sessionToken string) (error, *models.User) {
	err, userSession := r.sessionService.Find(sessionToken)
	if err != nil {
		return err, nil
	}
	return r.FindById(userSession.UserId)
}

func (r UserService) UpdateById(updaterId string, userId string, newUser *viewModels.UpdateUserViewModel) error {
	session, _ := r.newSession()
	defer session.Close()
	session.SetSafe(&mgo.Safe{})
	collection := session.DB(r.dbName).C(r.collectionName)
	err := collection.UpdateId(bson.ObjectIdHex(userId), bson.M{"$set": bson.M{
		"username":  newUser.UserName,
		"firstname": newUser.FirstName,
		"lastname":  newUser.LastName,
		"email":     newUser.Email,
		"image":     newUser.Image,
		"roleid":    newUser.RoleId,
		"updaterid": updaterId,
		"updatedat": time.Now().UTC()}})
	return err
}

func (r UserService) ChangePassword(updaterId string, userId string, newPassword string) error {
	session, _ := r.newSession()
	defer session.Close()
	session.SetSafe(&mgo.Safe{})
	collection := session.DB(r.dbName).C(r.collectionName)
	err := collection.UpdateId(bson.ObjectIdHex(userId), bson.M{"$set": bson.M{
		"password":  r.cryptUtil.Bcrypt(newPassword),
		"updaterid": updaterId,
		"updatedat": time.Now().UTC()}})
	return err
}

func (r UserService) VerifyEmail(token string) error {
	session, _ := r.newSession()
	defer session.Close()
	session.SetSafe(&mgo.Safe{})
	collection := session.DB(r.dbName).C(r.collectionName)
	err := collection.Update(bson.M{"verifytoken": token}, bson.M{"$set": bson.M{"emailverified": true}})
	return err
}

func (r UserService) Logout(token string) error {
	return r.sessionService.Logout(token)
}
