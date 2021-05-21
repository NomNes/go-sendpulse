package addressbooks

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/NomNes/go-sendpulse"
)

const (
	StatusActive = iota
	StatusDeleted
	statusUnknown
	StatusWaiting
	StatusBlockedByService
	StatusBlockedByDaemon
)

type payload struct {
	ID               int                       `json:"id"`
	Name             string                    `json:"name"`
	AllEmailQty      int                       `json:"all_email_qty"`
	ActiveEmailQty   int                       `json:"active_email_qty"`
	InactiveEmailQty int                       `json:"inactive_email_qty"`
	CreationDate     go_sendpulse.JsonDateTime `json:"creationdate"`
	Status           int                       `json:"status"`
	StatusExplain    string                    `json:"status_explain"`
}

type Item struct {
	payload
}

// ID List ID
func (i *Item) ID() int {
	return i.payload.ID
}

// Name List name
func (i *Item) Name() string {
	return i.payload.Name
}

// SetName Set List name
func (i *Item) SetName(name string) {
	i.payload.Name = name
}

// AllEmailQty Total number of emails
func (i *Item) AllEmailQty() int {
	return i.payload.AllEmailQty
}

// ActiveEmailQty Number of active emails
func (i *Item) ActiveEmailQty() int {
	return i.payload.ActiveEmailQty
}

// InactiveEmailQty Number of inactive emails
func (i *Item) InactiveEmailQty() int {
	return i.payload.InactiveEmailQty
}

// CreationDate Date of creation
func (i *Item) CreationDate() time.Time {
	return time.Time(i.payload.CreationDate)
}

// Status Status code
func (i *Item) Status() int {
	return i.payload.Status
}

// StatusExplain Status explanation
func (i *Item) StatusExplain() string {
	return i.payload.StatusExplain
}

// Save Creating/Editing the Name of a Mailing List
func (i *Item) Save(ctx context.Context, client go_sendpulse.Client) (err error) {
	var res string
	r := client.R().SetContext(ctx).SetResult(&res)
	if i.ID() > 0 {
		// Update
		r.SetBody(map[string]string{"name": i.Name()})
		_, err = r.Put(fmt.Sprintf("/addressbooks/%d", i.ID()))
		if err == nil {
			if res != "result = true" {
				err = errors.New("failed")
			}
		}
	} else {
		// Create
		r.SetBody(map[string]string{"bookName": i.Name()})
		_, err = r.Post("/addressbooks")
		var id int
		id, err = strconv.Atoi(res)
		if err == nil {
			if id > 0 {
				i.payload.ID = id
			} else {
				err = errors.New("failed")
			}
		}
	}
	return err
}

// Delete Deleting a Mailing List
func (c *Collection) Delete(ctx context.Context, id int) (err error) {
	var res go_sendpulse.ResultResponse
	r := c.client.R().SetContext(ctx).SetResult(&res)
	_, err = r.Delete(fmt.Sprintf("/addressbooks/%d", id))
	if err == nil {
		if !res.Result {
			err = errors.New("failed")
		}
	}
	return err
}

type Cost struct {
	Currency                  string `json:"cur"`
	SentEmailsQty             int    `json:"sent_emails_qty"`
	OverdraftAllEmailsPrice   int    `json:"overdraftAllEmailsPrice"`
	AddressesDeltaFromBalance int    `json:"addressesDeltaFromBalance"`
	AddressesDeltaFromTariff  int    `json:"addressesDeltaFromTariff"`
	MaxEmailsPerTask          int    `json:"max_emails_per_task"`
	Result                    bool   `json:"result"`
}

// CostOfCampaign Calculating the Cost of a Campaign Sent to a Mailing List
func (c *Collection) CostOfCampaign(ctx context.Context, id int) (cost Cost, err error) {
	r := c.client.R().SetContext(ctx).SetResult(&cost)
	_, err = r.Get(fmt.Sprintf("/addressbooks/%d/cost", id))
	return cost, err
}
