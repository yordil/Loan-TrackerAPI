package repository

import (
	"context"
	"fmt"
	"loantracker/domain"
	"loantracker/mongo"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository struct {
	database   mongo.Database
	collection string
}

func NewUserRepository(database mongo.Database, collection string) domain.UserRepository {
	return &UserRepository{
		database:   database,
		collection: collection}

}


func (u *UserRepository) CreateUser(c context.Context , user domain.User) (domain.User, error) {
	collection := u.database.Collection(u.collection)
	_, err := collection.InsertOne(c, user)

	if err != nil {
		return domain.User{} , err
	}

	return user , nil
}


func (u *UserRepository) FindUserByEmail(c context.Context  , email string) (domain.User, error) {
	collection := u.database.Collection(u.collection)
	var user domain.User
	fmt.Println(email , "**********")
	err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return domain.User{} , err
	}

	return user , nil
}


func (u*UserRepository) FindUserByResetToken(c context.Context , token string) (domain.User, error) {
	collection := u.database.Collection(u.collection)
	var user domain.User
	
	// Filter to find the user by reset_token
	filter := bson.M{"reset_password_token": token}

	err := collection.FindOne(c, filter).Decode(&user)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil

	
}

func (u *UserRepository) FindUserByVerificationToken(c context.Context , token string) (domain.User, error) {
	collection := u.database.Collection(u.collection)
	var user domain.User
	filter := bson.M{"verificationtoken": token}

	err := collection.FindOne(c, filter).Decode(&user)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}


func (u *UserRepository) UpdateUser(c context.Context , user domain.User) (domain.User, error) {
	collection := u.database.Collection(u.collection)
	update := bson.M{"$set": user}
	_, err := collection.UpdateOne(c, bson.M{"email": user.Email}, update)
	if err != nil {
		return domain.User{} , err
	}

	return user , nil
}



func (u *UserRepository) SetResetToken(c context.Context, email domain.ForgotPasswordRequest, token string , expiration time.Time) (domain.User, error) {
    // Reference to the collection
    collection := u.database.Collection(u.collection)

    // Filter to find the user by email
    filter := bson.M{"email": email.Email}

    // Update to set the reset_token
    update := bson.M{"$set": bson.M{"reset_password_token": token , "reset_password_expires": expiration}}

    // Execute the update
    _ , err := collection.UpdateOne(c, filter, update)
    if err != nil {
        return domain.User{}, err // Return an empty User and the error
    }

    var updatedUser domain.User
    err = collection.FindOne(c, filter).Decode(&updatedUser)
    if err != nil {
        return domain.User{}, err
    }

    return updatedUser, nil 

}

// find user By id

func (u *UserRepository) FindUserByID(c context.Context , id string) (domain.User, error) {
	collection := u.database.Collection(u.collection)
	var user domain.User
	// Convert the id to an ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.User{} , err
	}
	err = collection.FindOne(c, bson.M{"_id":objID}).Decode(&user)
	if err != nil {
		return domain.User{} , err
	}

	return user , nil
}


//  get all user 

func (u *UserRepository) GetAllUsers(c context.Context) ([]domain.User, error) {
	collection := u.database.Collection(u.collection)
	cursor, err := collection.Find(c, bson.M{})
	if err != nil {
		return nil, err
	}
	var users []domain.User
	if err = cursor.All(c, &users); err != nil {
		return nil, err
	}
	return users, nil

}


//  delete user

func (u *UserRepository) DeleteUser(c context.Context , id string) (error) {
	collection := u.database.Collection(u.collection)

	// Convert the id to an ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	
	if err != nil {
		return err
	}

	// Filter to find the user by id

	filter := bson.M{"_id": objID}

	// Execute the delete operation
	_, err = collection.DeleteOne(c, filter)

	if err != nil {
		return err
	}

	return nil

}

//  find user by refresh Token 

func (u *UserRepository) FindUserByRefreshToken(c context.Context , token string) (domain.User, error) {
	collection := u.database.Collection(u.collection)
	var user domain.User
	filter := bson.M{"refreshtoken": token}

	err := collection.FindOne(c , filter).Decode(&user) 


	if err != nil {
		return domain.User{} , err
	}

	return user , nil

	}