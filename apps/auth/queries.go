package auth

import (
	"github.com/rahulsoibam/koubru-prod-api/authutils/googlejwt"
)

func (a *App) dbGetUserIDUsingUsername(username string) (int64, error) {
	var userID int64
	var err error
	err = a.DB.QueryRow("SELECT user_id from KUser WHERE username = $1", username).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (a *App) dbGetUserIDUsingEmail(email string) (int64, error) {
	var userID int64
	var err error
	err = a.DB.QueryRow("SELECT user_id from KUser WHERE email = $1", email).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (a *App) dbGetEncodedHashUsingUserID(userID int64) (string, error) {
	var encodedHash string
	if err := a.AuthDB.QueryRow("SELECT password_hash from credential WHERE user_id = $1", userID).Scan(&encodedHash); err != nil {
		return "", err
	}
	return encodedHash, nil
}

func (a *App) dbRegisterUser(nu *NewUser) (int64, error) {
	var userID int64
	err := a.DB.QueryRow("INSERT INTO KUser (username, full_name, email) VALUES ($1, $2, $3) RETURNING user_id", nu.Username, nu.FullName, nu.Email).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (a *App) dbStorePassword(userID int64, encodedHash string) error {
	_, err := a.AuthDB.Exec("INSERT INTO credential (user_id, password_hash) VALUES ($1, $2)", userID, encodedHash)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) dbGetUserIDUsingFacebook(facebookID string) (int64, error) {
	var err error
	var userID int64
	err = a.DB.QueryRow("SELECT user_id FROM KUser WHERE facebook=$1", facebookID).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (a *App) dbGetUserIDUsingGoogle(googleID string) (int64, error) {
	var err error
	var userID int64
	err = a.DB.QueryRow("SELECT user_id FROM KUser WHERE google=$1", googleID).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (a *App) dbRegisterUserUsingFacebook(fu FacebookUser, username string) (int64, error) {
	var err error
	var userID int64
	err = a.DB.QueryRow("INSERT INTO KUser (username, full_name, photo_url, facebook) VALUES ($1, $2, $3, $4) RETURNING user_id", username, fu.Name, fu.Picture.Data.URL, fu.ID).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (a *App) dbRegisterUserUsingGoogle(cs *googlejwt.ClaimSet, username string) (int64, error) {
	var err error
	var userID int64
	err = a.DB.QueryRow("INSERT INTO KUser (username, full_name, photo_url, google) VALUES ($1, $2, $3, $4) RETURNING user_id", username, cs.Name, cs.Picture, cs.Sub).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
