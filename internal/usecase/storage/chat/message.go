package chat

import (
	"errors"

	"github.com/coffee/support/internal/entity"
	"github.com/gosuit/e"
	"github.com/gosuit/lec"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MessageStorage struct {
	db *mongo.Database
}

func New(db *mongo.Database) *MessageStorage {
	return &MessageStorage{
		db: db,
	}
}

func (s *MessageStorage) WriteMessage(ctx lec.Context, msg entity.Message) e.Error {
	collection := s.db.Collection(MessageCollection)
	msg.MessageId = bson.NewObjectID().Hex()

	_, err := collection.InsertOne(ctx, msg)
	if err != nil {
		return e.E(err).WithMessage("Failed to send message")
	}

	return nil
}

func (s *MessageStorage) GetChat(ctx lec.Context, userId uint64) ([]*entity.Message, e.Error) {
	var res []*entity.Message

	collection := s.db.Collection(MessageCollection)

	filter := bson.M{"userId": userId}

	c, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, e.E(err).WithMessage("Failed to find chat")
	}

	err = c.All(ctx, &res)
	if err != nil {
		return nil, e.E(err).WithMessage(InternalErrorMessage)
	}

	return res, nil
}

func (s *MessageStorage) GetSupportChats(ctx lec.Context, supportId uint64) ([]uint64, e.Error) {
	var SupportChats []entity.SupportToChat

	collection := s.db.Collection(ChatToSupportCollection)

	filter := bson.M{"supportId": supportId}

	c, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, e.E(err).WithMessage("Failed to find chats")
	}

	err = c.All(ctx, &SupportChats)
	if err != nil {
		return nil, e.E(err).WithMessage(InternalErrorMessage)
	}

	res := getUserIdFromEntity(SupportChats)

	return res, nil
}

func (s *MessageStorage) GetSupportIdFromUser(ctx lec.Context, userId uint64) (uint64, e.Error) {
	var userToSupport entity.SupportToChat

	collection := s.db.Collection(ChatToSupportCollection)

	filter := bson.M{"userId": userId}

	err := collection.FindOne(ctx, filter).Decode(&userToSupport)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, nil
		}

		return 0, e.E(err).WithMessage("Failed to get support")
	}

	return userToSupport.SupportId, nil
}

func (s *MessageStorage) SetSupportToChat(ctx lec.Context, sc entity.SupportToChat) e.Error {
	collection := s.db.Collection(ChatToSupportCollection)

	collection.InsertOne(ctx, sc)

	collection = s.db.Collection(SupportCollection)

	update := bson.M{"$inc": bson.M{"countChat": 1}}

	filter := bson.M{"supportId": sc.SupportId}
	_, err := collection.UpdateOne(ctx, filter, update, options.UpdateOne().SetUpsert(true))
	if err != nil {
		return e.E(err).WithMessage("Failed to set support for chat")
	}

	return nil
}

func (s *MessageStorage) ChooseSupport(ctx lec.Context) (entity.Support, e.Error) {
	var res entity.Support

	collection := s.db.Collection(SupportCollection)

	pipeline := mongo.Pipeline{
		{
			{"$group", bson.D{
				{"_id", nil},
				{"countChat", bson.D{{"$min", "$countChat"}}},
			}},
		},
	}

	ctx.Logger().Debug("Running aggregation pipeline")

	c, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return res, e.E(err).WithMessage(InternalErrorMessage)
	}

	ctx.Logger().Debug("Aggregation executed successfully")

	var result struct {
		CountChat int `bson:"countChat"`
	}

	if c.Next(ctx) {
		if err := c.Decode(&result); err != nil {
			return res, e.E(err).WithMessage(InternalErrorMessage)
		}
	} else if err := c.Err(); err != nil {
		return res, e.E(err).WithMessage(InternalErrorMessage)
	}

	filter := bson.M{"countChat": result.CountChat}

	err = collection.FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return res, e.E(err).WithMessage(InternalErrorMessage)
	}

	return res, nil

}

func (s *MessageStorage) ReadChat(ctx lec.Context, userId, readerId uint64) e.Error {
	collection := s.db.Collection(MessageCollection)

	filter := bson.M{"senderId": bson.M{"$ne": readerId}, "userId": userId}

	update := bson.M{"$set": bson.M{"isRead": true}}

	_, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return e.E(err).WithMessage("Failed to read messages")
	}

	return nil
}

func getUserIdFromEntity(in []entity.SupportToChat) []uint64 {
	res := make([]uint64, 0, len(in))

	for _, i := range in {
		res = append(res, i.UserId)
	}

	return res
}
