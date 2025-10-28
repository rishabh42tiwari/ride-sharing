package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type RideFareModel struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	UserID          string             `bson:"userId"`
	PackageSlug     string             `bson:"packageSlug"`
	TotalPriceCents float64            `bson:"totalPriceCents"`
}