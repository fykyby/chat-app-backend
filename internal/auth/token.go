package auth

import (
	"context"
	"errors"

	"github.com/fykyby/chat-app-backend/internal/database"
	"github.com/fykyby/chat-app-backend/internal/model"
	"github.com/go-chi/jwtauth/v5"
)

func GetClaimedUser(ctx context.Context, db *database.Queries) (model.ClaimedUser, error) {
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

	if name, ok := claims["name"]; ok {
		claimedUser.Name = name.(string)
	} else {
		return claimedUser, errors.New("Name not found in claims")
	}

	if email, ok := claims["email"]; ok {
		claimedUser.Email = email.(string)
	} else {
		return claimedUser, errors.New("Email not found in claims")
	}

	_, err = db.GetUserByData(ctx, database.GetUserByDataParams{
		ID:    claimedUser.ID,
		Name:  claimedUser.Name,
		Email: claimedUser.Email,
	})
	if err != nil {
		return claimedUser, errors.New("Invalid JWT user")
	}

	return claimedUser, nil
}
