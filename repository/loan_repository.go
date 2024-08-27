package repository

import (
	"context"
	"loantracker/domain"
	"loantracker/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LoanRepository struct {
	database   mongo.Database
	collection string
}

func NewLoanRepository(database mongo.Database, collection string) domain.LoanRepository {
	return &LoanRepository{
		database:   database,
		collection: collection}

}

func (l *LoanRepository) CreateLoan(ctx context.Context , loan domain.Loan) (domain.Loan, error) {
	collection := l.database.Collection(l.collection)
	_, err := collection.InsertOne(ctx, loan)

	if err != nil {
		return domain.Loan{} , err	
	}

	return loan , nil

}


func (l *LoanRepository) GetAllLoans(ctx context.Context, status string, order string) ([]domain.Loan, error) {
    collection := l.database.Collection(l.collection)

    // Build the filter based on the status parameter
    filter := bson.M{}
    if status != "" && status != "all" {
        filter["status"] = status
    }

    // Determine the sort order
    sortOrder := 1 // Default is ascending
    if order == "desc" {
        sortOrder = -1
    }

    // Build the sort criteria
    sort := bson.D{{"status", sortOrder}}

    // Execute the query
    cursor, err := collection.Find(ctx, filter, options.Find().SetSort(sort))
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var loans []domain.Loan

    // Decode the documents
    for cursor.Next(ctx) {
        var loan domain.Loan
        if err := cursor.Decode(&loan); err != nil {
            return nil, err
        }
        loans = append(loans, loan)
    }



    return loans, nil
}


func (l *LoanRepository) GetLoanByID(ctx context.Context , id string) (domain.Loan , error) {

	collection := l.database.Collection(l.collection)
	var Loan domain.Loan
	// Convert the id to an ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Loan{} , err
	}
	err = collection.FindOne(ctx, bson.M{"_id":objID}).Decode(&Loan)
	if err != nil {
		return domain.Loan{} , err
	}

	return Loan , nil

}
//  update Loan by id 
func (l *LoanRepository) UpdateLoan(ctx context.Context , id string , loan domain.Loan) (domain.Loan , error) {
	collection := l.database.Collection(l.collection)
	// Convert the id to an ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Loan{} , err
	}
	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{
			"condition": loan.Condition,
		},
	}
	_ , err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return domain.Loan{} , err
	}


	var updatedLoan domain.Loan
	err = collection.FindOne(ctx, filter).Decode(&updatedLoan)
	if err != nil {
		return domain.Loan{} , err
	}

	return updatedLoan , nil
}

// delete loan by id

func (l *LoanRepository) DeleteLoan(ctx context.Context , id string) error {
	collection := l.database.Collection(l.collection)
	// Convert the id to an ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	return nil
}