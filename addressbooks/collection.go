package addressbooks

import (
	"context"
	"fmt"

	"github.com/NomNes/go-sendpulse"
)

type Collection struct {
	client go_sendpulse.Client
}

func New(client go_sendpulse.Client) Collection {
	return Collection{client: client}
}

// GetOne Retrieving Mailing List Information
func (c *Collection) GetOne(ctx context.Context, id int) (item *Item, err error) {
	var items []Item
	r := c.client.R().SetContext(ctx).SetResult(&items)
	_, err = r.Get(fmt.Sprintf("/addressbooks/%d", id))
	if len(items) > 0 {
		item = &items[0]
	}
	return item, err
}

// GetList Retrieving a List of Mailing Lists
func (c *Collection) GetList(ctx context.Context, options *go_sendpulse.ListOptions) (list []Item, err error) {
	var l []payload
	r := c.client.R().SetContext(ctx).SetResult(&l)
	if options != nil {
		if options.Limit != nil {
			r.SetQueryParam("limit", fmt.Sprintf("%d", *options.Limit))
		}
		if options.Offset != nil {
			r.SetQueryParam("offset", fmt.Sprintf("%d", *options.Offset))
		}
	}
	_, err = r.Get("/addressbooks")
	for _, item := range l {
		list = append(list, Item{item})
	}
	return list, err
}

// GetVariables Get a List of Variables for a Mailing List
func (c *Collection) GetVariables(ctx context.Context, id int) (list map[string]string, err error) {
	r := c.client.R().SetContext(ctx).SetResult(list)
	_, err = r.Get(fmt.Sprintf("/addressbooks/%d/variables", id))
	return list, err
}
