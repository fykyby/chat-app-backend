package auth

import (
	"context"
	"errors"

	"github.com/fykyby/chat-app-backend/internal/model"
	"github.com/go-chi/jwtauth/v5"
)

func GetClaimedUser(ctx context.Context) (model.ClaimedUser, error) {
	claimedUser := model.ClaimedUser{}

	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return claimedUser, err
	}

	if id, ok := claims["id"]; ok {
		claimedUser.ID = int32(id.(float64))
	} else {
		return claimedUser, errors.New("ID not found in claims")
	}

	if _, ok := claims["name"]; ok {
		claimedUser.Name = claims["name"].(string)
	} else {
		return claimedUser, errors.New("Name not found in claims")
	}

	if _, ok := claims["email"]; ok {
		claimedUser.Email = claims["email"].(string)
	} else {
		return claimedUser, errors.New("Email not found in claims")
	}

	// TODO?: Check if user with combined data above exists in database

	return claimedUser, nil
}
