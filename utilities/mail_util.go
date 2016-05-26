package utilities

import (
	"github.com/hjr265/postmark.go/postmark"
	"net/mail"
	"strings"
	"fmt"
)

type IMailUtil interface {
	Send(toEmail string, fromEmail string, fromName string, subject string, body string) (*postmark.Result, error)
	SendTemplate(toEmail string, fromEmail string, fromName string, subject string, globalVars map[string]interface{}, templateId int) (*postmark.Result, error)
}

type MailUtil struct {
	postmarkClient postmark.Client
}

func NewMailUtil(configUtil IConfigUtil) IMailUtil {
	util := MailUtil{}
	postmarkKey := configUtil.GetConfig("postmarkKey")
	util.postmarkClient = postmark.Client{
		ApiKey: postmarkKey,
		Secure: true,
	}
	return &util
}

func (util MailUtil) Send(toEmail string, fromEmail string, fromName string, subject string, body string) (*postmark.Result, error) {
	return util.postmarkClient.Send(&postmark.Message{
		From: &mail.Address{
			Name:    fromName,
			Address: fromEmail,
		},
		To: []*mail.Address{
			{
				Name:    "",
				Address: toEmail,
			},
		},
		Subject:  subject,
		TextBody: strings.NewReader(body),
	})
}

func (util MailUtil) SendTemplate(toEmail string, fromEmail string, fromName string, subject string, globalVars map[string]interface{}, templateId int) (*postmark.Result, error) {
	result, err := util.postmarkClient.Send(&postmark.Message{
		From: &mail.Address{
			Name:    fromName,
			Address: fromEmail,
		},
		To: []*mail.Address{
			{
				Name:    "",
				Address: toEmail,
			},
		},
		TemplateId: templateId,
		TemplateModel:globalVars})
	if err != nil {
		fmt.Println("error sending template e-mail")
	}
	fmt.Println(result)
	return result, err
}