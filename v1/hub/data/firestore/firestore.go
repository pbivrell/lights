package firestore

import (
	"context"
	"fmt"
	"hub/data"

	"cloud.google.com/go/firestore"
)

const (
	projectID       = "lights-403101"
	lightCollection = "devices/lights/"
)

type Connection struct {
	client *firestore.Client
}

func Connect(ctx context.Context) (*Connection, error) {

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return &Connection{
		client: client,
	}, nil
}

func (c *Connection) Close() {
	c.client.Close()
}

func (c *Connection) WriteLight(ctx context.Context, u data.Light) error {

	doc := c.client.Doc(fmt.Sprintf("%s%d", lightCollection, u.ID))

	_, err := doc.Create(ctx, u)
	return err
}

func (c *Connection) ReadLight(ctx context.Context, id int64) (data.Light, error) {

	var user data.Light

	doc := c.client.Doc(fmt.Sprintf("%s%d", lightCollection, id))
	docsnap, err := doc.Get(ctx)
	if err != nil {
		return user, err
	}

	return user, docsnap.DataTo(&user)
}
