package addressbooks

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/NomNes/go-sendpulse"
)

const (
	EmailStatusNew = iota
	EmailStatusActive
	EmailStatusActivationRequested
	EmailStatusActivationPending
	EmailStatusUnsubscribed
	EmailStatusRejected
	EmailStatusUnsubscribedFromAll
	EmailStatusActivationSent
	EmailStatusBlockedByUser
	EmailStatusDeliveryErrors
	EmailStatusBlockedByHosts
	EmailStatusBlockedBySendersName
	EmailStatusBlockedByAddress
	EmailStatusDeletedByUser
	EmailStatusTemporarilyUnavailable
)

type Email struct {
	Email         string `json:"email"`
	Status        int    `json:"status"`
	StatusExplain string `json:"status_explain"`
}

type EmailWithExtra struct {
	Email
	Phone     string            `json:"phone"`
	Variables map[string]string `json:"variables"`
}

// GetEmails Retrieving a List of Emails from a Mailing List
func (c *Collection) GetEmails(ctx context.Context, id int, options *go_sendpulse.ListOptions) (list []EmailWithExtra, err error) {
	r := c.client.R().SetContext(ctx).SetResult(&list)
	if options != nil {
		if options.Limit != nil {
			r.SetQueryParam("limit", fmt.Sprintf("%d", *options.Limit))
		}
		if options.Offset != nil {
			r.SetQueryParam("offset", fmt.Sprintf("%d", *options.Offset))
		}
	}
	_, err = r.Get(fmt.Sprintf("/addressbooks/%d/emails", id))
	return list, err
}

func (c *Collection) GetEmailsTotal(ctx context.Context, id int) (total int, err error) {
	var res string
	r := c.client.R().SetContext(ctx).SetResult(&res)
	_, err = r.Get(fmt.Sprintf("/addressbooks/%d/emails/total", id))
	if err == nil {
		total, err = strconv.Atoi(res)
	}
	return total, err
}

// FindByVariable Find All Contacts in Mailing List by Value of Variable
func (c *Collection) FindByVariable(ctx context.Context, id int, variableName, searchValue string) (list []Email, err error) {
	r := c.client.R().SetContext(ctx).SetResult(&list)
	_, err = r.Get(fmt.Sprintf("/addressbooks/%d/variables/%s/%s", id, variableName, searchValue))
	return list, err
}

type EmailWithVariables struct {
	Email     string            `json:"email"`
	Variables map[string]string `json:"variables,omitempty"`
}

type AddEmailsOptions struct {
	Confirmation bool
	SenderEmail  string
	TemplateId   string
	MessageLang  string
}

// AddEmailsWithVariables Add single-opt-in
func (c *Collection) AddEmailsWithVariables(ctx context.Context, id int, emails []EmailWithVariables, options *AddEmailsOptions) error {
	body := map[string]interface{}{"emails": emails}
	if options != nil {
		if options.Confirmation {
			body["confirmation"] = true
		}
		if options.SenderEmail != "" {
			body["sender_email"] = options.SenderEmail
		}
		if options.TemplateId != "" {
			body["template_id"] = options.TemplateId
		}
		if options.MessageLang != "" {
			body["message_lang"] = options.MessageLang
		}
	}
	var res go_sendpulse.ResultResponse
	r := c.client.R().SetContext(ctx).SetHeader("Content-Type", "application/json").SetBody(body).SetResult(&res)
	_, err := r.Post(fmt.Sprintf("/addressbooks/%d/emails", id))
	if err == nil {
		if !res.Result {
			err = errors.New("failed")
		}
	}
	return err
}

// AddEmails Add single-opt-in
func (c *Collection) AddEmails(ctx context.Context, id int, emails []string, options *AddEmailsOptions) error {
	var items []EmailWithVariables
	for _, email := range emails {
		items = append(items, EmailWithVariables{Email: email})
	}
	return c.AddEmailsWithVariables(ctx, id, items, options)
}

// DeleteEmails Deleting Email Addresses from a Mailing List
func (c *Collection) DeleteEmails(ctx context.Context, id int, emails []string) error {
	var res go_sendpulse.ResultResponse
	r := c.client.R().SetContext(ctx).SetResult(res).SetBody(emails)
	_, err := r.Delete(fmt.Sprintf("/addressbooks/%d/emails", id))
	if err == nil {
		if !res.Result {
			err = errors.New("failed")
		}
	}
	return err
}

type EmailInformation struct {
	Email
	AddressBookId int             `json:"abook_id"`
	Variables     []EmailVariable `json:"variables"`
}

type EmailVariable struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// GetEmail Retrieving Information for a Specific Email Address from a Mailing List
func (c *Collection) GetEmail(ctx context.Context, id int, email string) (list []EmailInformation, err error) {
	r := c.client.R().SetContext(ctx)
	_, err = r.Get(fmt.Sprintf("/addressbooks/%d/emails/%s", id, email))
	return list, err
}
