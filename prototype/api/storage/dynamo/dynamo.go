package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pbivrell/lights/api/storage"
)

const (
	Region       = "us-east-1"
	CredsProfile = "lights-dev"
	LightTable   = "light-lights"
	UserTable    = "light-users"
	SessionTable = "light-sessions"
)

type Dynamo struct {
	client *dynamodb.DynamoDB
}

func New() *Dynamo {
	return &Dynamo{
		client: dynamodb.New(session.New(), aws.NewConfig().WithRegion(Region)),
	}
}

func NewFromProfile(profile string) (*Dynamo, error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewSharedCredentials("", profile),
		Region:      aws.String(Region),
	})

	return &Dynamo{
		client: dynamodb.New(sess),
	}, err
}

func (d *Dynamo) write(table string, datum interface{}) error {
	av, err := dynamodbattribute.MarshalMap(datum)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(table),
	}

	_, err = d.client.PutItem(input)
	return err
}

func (d *Dynamo) read(table string, key string, data interface{}) error {
	result, err := d.client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(key),
			},
		},
	})
	if err != nil {
		return err
	}

	if result.Item == nil {
		return storage.ErrorNotFound
	}

	return dynamodbattribute.UnmarshalMap(result.Item, data)

}

func (d *Dynamo) delete(table string, key string) error {
	_, err := d.client.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String(table),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(key),
			},
		},
	})

	return err
}

func (d *Dynamo) DeleteSession(key string) error {
	return d.delete(SessionTable, key)
}

func (d *Dynamo) WriteSession(s *storage.Session) error {
	return d.write(SessionTable, s)
}

func (d *Dynamo) ReadSession(key string) (*storage.Session, error) {

	sess := &storage.Session{}

	return sess, d.read(SessionTable, key, sess)
}

func (d *Dynamo) DeleteUser(key string) error {
	return d.delete(UserTable, key)
}

func (d *Dynamo) WriteUser(u *storage.User) error {
	return d.write(UserTable, u)
}

func (d *Dynamo) ReadUser(key string) (*storage.User, error) {

	user := &storage.User{}

	return user, d.read(UserTable, key, user)
}

func (d *Dynamo) WriteLight(l *storage.Light) error {
	return d.write(LightTable, l)
}

func (d *Dynamo) ReadLight(key string) (*storage.Light, error) {

	light := &storage.Light{}

	return light, d.read(LightTable, key, light)
}
