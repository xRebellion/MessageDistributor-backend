package recipient

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/xRebellion/MessageDistributor/integration"
)

type Recipient struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name        string             `bson:"name" json:"name"`
	ContactInfo []Contact          `bson:"contactInfo" json:"contactInfo"`
}

type Contact struct {
	MediaName   string `bson:"mediaName" json:"mediaName"`     // Name of the communication media / platform
	MediaUserID string `bson:"mediaUserId" json:"mediaUserId"` // The identifier that is acknowledged by the target platform, such as phone number, tag, or something similar
	Preferred   bool   `bson:"preferred" json:"preferred"`     // Whether this contact info is preferrable to the recipient or not to receive messages
	IsSendable  bool   `bson:"-" json:"isSendable"`            // Is this contact info available to be sent
}

type Service struct {
	integrations map[string]integration.Integration
	collection   *mongo.Collection
}

func (s *Service) CreateRecipient(ctx context.Context, name string, contactInfo []Contact) (*Recipient, error) {
	res, err := s.collection.InsertOne(ctx, Recipient{
		Name:        name,
		ContactInfo: contactInfo,
	})
	if err != nil {
		return nil, err
	}
	return &Recipient{
		ID:          res.InsertedID.(primitive.ObjectID),
		Name:        name,
		ContactInfo: contactInfo,
	}, nil
}

func (s *Service) ListRecipients(ctx context.Context) ([]Recipient, error) {
	cursor, err := s.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	result := []Recipient{}
	for cursor.Next(ctx) {
		var bsonResult Recipient
		if err := cursor.Decode(&bsonResult); err != nil {
			return nil, err
		}
		for i, ci := range bsonResult.ContactInfo {
			if _, isSendable := s.integrations[ci.MediaName]; isSendable {
				ci.IsSendable = isSendable
				bsonResult.ContactInfo[i] = ci
			}
		}
		result = append(result, bsonResult)
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) GetRecipientByID(ctx context.Context, id string) (*Recipient, error) {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	res := s.collection.FindOne(ctx, bson.M{"_id": objectID})

	var recipient *Recipient
	if err := res.Decode(&recipient); err != nil {
		return nil, err
	}
	return recipient, nil
}

func NewService(integrations map[string]integration.Integration, collection *mongo.Collection) *Service {
	return &Service{
		integrations: integrations,
		collection:   collection,
	}
}
