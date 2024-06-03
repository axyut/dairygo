package utils

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetMongoTimeFromHTMLDate(date string) (dateTime primitive.DateTime) {

	dateTime = primitive.NewDateTimeFromTime(time.Now())
	if date != "" {
		timeP, err := time.Parse("2006-01-02", date)
		if err != nil {
			return
		}
		timeP = timeP.Add(time.Hour * time.Duration(time.Now().UTC().Hour()))
		timeP = timeP.Add(time.Minute * time.Duration(time.Now().UTC().Minute()))
		timeP = timeP.Add(time.Second * time.Duration(time.Now().UTC().Second()))
		timeP = timeP.Local()
		dateTime = primitive.NewDateTimeFromTime(timeP)
	}
	return
}
