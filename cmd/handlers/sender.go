package handlers

import (
	"fmt"
	"net/url"
	"os"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/manubidegain/piggy-api/cmd/api/configuration"
	"github.com/manubidegain/piggy-api/cmd/entities"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

const TWILIO_PHONE_NUMBER = "+19382531274"

func EmailSignup(client *auth.Client,
	ctx *gin.Context,
	config *configuration.Config,
	user *entities.User,
	senderName string,
	userID string) {

	url := fmt.Sprintf("%smember-login?email=%s&first_name=%s&last_name=%s&invited_by=%s&&user_id=%s",
		config.Sender.CallbackURL,
		url.QueryEscape(user.Email),
		user.FirstName,
		user.LastName,
		senderName,
		userID,
	)

	actionCodeSettings := &auth.ActionCodeSettings{
		URL:             url,
		HandleCodeInApp: false,
	}
	emailLink, err := client.EmailSignInLink(ctx, user.Email, actionCodeSettings)
	if err != nil {
		fmt.Println("An error has ocurred creating the email link" + err.Error())
	}
	fmt.Println(emailLink)
	from := mail.NewEmail("Piggy", "piggy@gmail.com")
	apikey := os.Getenv("SENDGRID_API_KEY")
	senderClient := sendgrid.NewSendClient(apikey)
	m := mail.NewV3Mail()
	m.SetFrom(from)

	p := mail.NewPersonalization()
	tos := []*mail.Email{
		mail.NewEmail(user.FirstName, user.Email),
	}
	p.AddTos(tos...)
	m.SetTemplateID("d-d2ae83297aac42f1921e101c0b89820d")
	p.SetDynamicTemplateData("inviteName", senderName)
	p.SetDynamicTemplateData("buttonUrl", emailLink)

	m.AddPersonalizations(p)
	response, err := senderClient.Send(m)
	if err != nil {
		fmt.Println("An error has ocurred sending the email" + err.Error())
	}
	fmt.Println(response)

}

func PasswordResetEmail(client *auth.Client, ctx *gin.Context, config *configuration.Config, email string) error {

	url := fmt.Sprintf("%sreset-password",
		config.Sender.CallbackURL,
	)

	actionCodeSettings := &auth.ActionCodeSettings{
		URL:             url,
		HandleCodeInApp: false,
	}
	emailLink, err := client.PasswordResetLinkWithSettings(ctx, email, actionCodeSettings)
	if err != nil {
		return err
	}
	fmt.Println(emailLink)
	from := mail.NewEmail("piggy", "nacho@piggybanking.com")
	apikey := os.Getenv("SENDGRID_API_KEY")
	senderClient := sendgrid.NewSendClient(apikey)
	m := mail.NewV3Mail()
	m.SetFrom(from)

	p := mail.NewPersonalization()
	tos := []*mail.Email{
		mail.NewEmail("piggy user", email),
	}
	p.AddTos(tos...)
	m.SetTemplateID("d-6886e1422d2244e3aadd1c38dffcc2fe")
	p.SetDynamicTemplateData("buttonUrl", emailLink)

	m.AddPersonalizations(p)
	response, err := senderClient.Send(m)
	fmt.Println(response)
	if err != nil {
		return err
	}
	return nil
}

func SendEmailRandomNumber(client *auth.Client, ctx *gin.Context, config *configuration.Config, email string, random string) error {

	from := mail.NewEmail("Roookie", "nacho@piggybanking.com")
	apikey := os.Getenv("SENDGRID_API_KEY")
	senderClient := sendgrid.NewSendClient(apikey)
	m := mail.NewV3Mail()
	m.SetFrom(from)

	p := mail.NewPersonalization()
	tos := []*mail.Email{
		mail.NewEmail("piggy user", email),
	}
	p.AddTos(tos...)
	m.SetTemplateID("d-93683e1c14d44bf3b3b75dbc0ea08caa")
	p.SetDynamicTemplateData("randomNumber", random)

	m.AddPersonalizations(p)
	response, err := senderClient.Send(m)
	fmt.Println(response)
	if err != nil {
		return err
	}
	return nil
}

func SendSMSRandomNumber(ctx *gin.Context, number string, random string) error {

	client := twilio.NewRestClient()

	params := &openapi.CreateMessageParams{}
	params.SetTo(number)
	params.SetFrom(TWILIO_PHONE_NUMBER)
	params.SetBody(fmt.Sprintf("Your piggy verification code from piggy is : %v", random))

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		return err
	} else {
		return nil
	}
}
