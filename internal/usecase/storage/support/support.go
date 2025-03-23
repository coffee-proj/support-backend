package support

import (
	"github.com/coffee/support/internal/entity"
	"github.com/gosuit/e"
	"github.com/gosuit/lec"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type SupportStorage struct {
	db *mongo.Database
}

func New(db *mongo.Database) *SupportStorage {
	return &SupportStorage{
		db: db,
	}
}

func (s *SupportStorage) AddSupport(ctx lec.Context, supportId uint64) e.Error {
	collection := s.db.Collection(SupportCollection)

	r := collection.FindOne(ctx, bson.M{"supportId": supportId})

	if r.Decode(&entity.Support{}) == nil {
		return e.New("Support with this id already exists", e.BadInput)
	}

	support := entity.Support{
		SupportId: supportId,
	}

	_, err := collection.InsertOne(ctx, support)
	if err != nil {
		return e.E(err).WithMessage("Failed to add support")
	}

	return nil
}

func (s *SupportStorage) ReplaceSupport(ctx lec.Context, targetId, replacementId uint64) e.Error {
	collection := s.db.Collection(ChatToSupportCollection)

	filter := bson.M{"supportId": targetId}

	update := bson.M{"$set": bson.M{"supportId": replacementId}}

	res, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return e.E(err).WithMessage("Failed to set support for chat")
	}

	collection = s.db.Collection(SupportCollection)

	filter = bson.M{"supportId": replacementId}

	update = bson.M{"$inc": bson.M{"countChat": res.ModifiedCount}}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return e.E(err).WithMessage("Failed to set support for chat")
	}

	return nil
}

func (s *SupportStorage) RemoveSupport(ctx lec.Context, supportId uint64) e.Error {
	collection := s.db.Collection(SupportCollection)

	filter := bson.M{"supportId": supportId}

	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return e.E(err).WithMessage("Failed to remove support")
	}

	return nil
}

func (s *SupportStorage) GetAllSupports(ctx lec.Context, offset, limit int64) ([]entity.Support, e.Error) {
	result := []entity.Support{}
	collection := s.db.Collection(SupportCollection)

	cur, err := collection.Find(ctx, bson.M{}, options.Find().SetSkip(offset).SetLimit(limit))
	if err != nil {
		return nil, e.E(err).WithMessage("Failed to find supports")
	}

	if err := cur.All(ctx, &result); err != nil {
		return nil, e.E(err).WithMessage("Failed to parse supports")
	}

	return result, nil
}
