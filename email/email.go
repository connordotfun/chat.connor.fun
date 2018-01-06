package email

import (
	mailjet "github.com/mailjet/mailjet-apiv3-go"
	"github.com/aaronaaeng/chat.connor.fun/config"
	"fmt"
)

const (
	sender = "no-reply@connor.fun"

	accountVerificationEmailTemplate = `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title>Account Verification</title>
		</head>
		<body style="margin: 0; padding: 0;">
			<table align="center" border="1" cellpadding="0" cellspacing="0" width="600" style="border-collapse: collapse;
					border: 2px solid #dbdbdb; border-right-color: #535353; border-bottom-color: #535353; background-color: #bfbfbf;">
				<tr style="padding-bottom: 50px;">
					<td style="padding-left: 4px;">
						<h1>Thanks for signing up!</h1>
						<em>Please verify your email to activate your account</em>
					</td>
				</tr>
				<tr>
					<td style="padding-left: 4px;">
						<strong> Hey {{var:username}},</strong>
						<br/>
						It looks like you've signed up for a new chat.connor.fun account! Before you can start chatting
						on the coolest hangout place on the world wide web, please verify your account.
					</td>
				</tr>
				<tr>
					<td align="center" style="padding: 10px 10px 10px 10px">
						<a id="verify" href="{{var:link}}" target="_blank" style="padding: 4px 3px 4px 3px; text-decoration: none;
							border: 2px solid #dbdbdb; border-right-color: #535353; border-bottom-color: #535353; background-color: #bfbfbf;">
							Verify Email
						</a>
					</td>
				</tr>
			</table>
		</body>
		</html>
	`
)

type verificationEmailVars struct {
	Username string
	Link string
}

func SendAccountVerificationEmail(toEmail string, username string, veriLink string) error {
	mailjetClient := mailjet.NewMailjetClient(config.MailjetPubKey, config.MailjetPrivKey)

	mail := &mailjet.InfoSendMail{
		FromEmail: sender,
		Recipients: []mailjet.Recipient{
			{
				Email: toEmail,
			},
		},
		Subject:  "Verify you chat.connor.fun account",
		HTMLPart: accountVerificationEmailTemplate,
		MjTemplateLanguage: "true",
		Vars: verificationEmailVars{
			Username: username,
			Link: veriLink,
		},
	}
	res, err := mailjetClient.SendMail(mail)
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}