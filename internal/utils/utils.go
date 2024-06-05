package utils

import (
	"math"
	"strconv"
	"time"

	"github.com/axyut/dairygo/internal/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetMongoTimeFromHTMLDate(date string, defaultTime time.Time) (dateTime primitive.DateTime) {

	dateTime = primitive.NewDateTimeFromTime(defaultTime)
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

func SetToPayToRecieve(trans_type string, payment bool, Price float64, ToPay float64, ToReceive float64) (toPay float64, toReceive float64) {

	// update audience to pay and to receive
	if trans_type == string(types.Sold) {
		if !payment {
			// case when user has remaining ToPay to audience
			if ToPay > 0 {
				ToPay -= Price
				// fmt.Println("here sell 1", ToPay)
				if ToPay < 0 {
					ToReceive += math.Abs(ToPay) // convert to positive
					ToPay = 0
					// fmt.Println("here sell 2", ToPay, ToReceive)
				}
			} else {
				ToReceive += Price
			}
			// fmt.Println("here sell", ToReceive)
		}
		if payment {
			if ToPay > 0 {
				ToPay -= Price
				if ToPay < 0 {
					ToReceive += math.Abs(ToPay) // convert to positive
					ToPay = 0
				}
			}
		}
	} else if trans_type == string(types.Bought) {
		if !payment {

			// case when user has remaining ToReceive from audience
			if ToReceive > 0 {
				ToReceive -= Price
				if ToReceive < 0 {
					ToPay += math.Abs(ToReceive) // convert to positive
					ToReceive = 0
				}
			} else {
				ToPay += Price
			}
			// fmt.Println("here buy", ToPay)
		}
		if payment {
			if ToReceive > 0 {
				ToReceive -= Price
				if ToReceive < 0 {
					ToPay += math.Abs(ToReceive) // convert to positive
					ToReceive = 0
				}
			}
		}
	}
	// fmt.Println("final ", ToPay, ToReceive)

	return ToPay, ToReceive
}

func Str(f float64) string {
	if strconv.FormatFloat(f, 'f', -1, 64) == "0" {
		return ""
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}
