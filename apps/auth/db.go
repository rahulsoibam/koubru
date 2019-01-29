package auth

import (
	"database/sql"

	"github.com/rahulsoibam/koubru-prod-api/authutils/googlejwt"
)

func dbGetUserIDUsingUsername(db *sql.DB, username string) (int64, error) {
	var userID int64
	var err error
	err = db.QueryRow("SELECT user_id from KUser WHERE username = $1", username).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func dbGetUserIDUsingEmail(db *sql.DB, email string) (int64, error) {
	var userID int64
	var err error
	err = db.QueryRow("SELECT user_id from KUser WHERE email = $1", email).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func dbGetEncodedHashUsingUserID(authDB *sql.DB, userID int64) (string, error) {
	var encodedHash string
	if err := authDB.QueryRow("SELECT password_hash from credential WHERE user_id = $1", userID).Scan(&encodedHash); err != nil {
		return "", err
	}
	return encodedHash, nil
}

func dbRegisterUser(db *sql.DB, nu *NewUser) (int64, error) {
	var userID int64
	err := db.QueryRow("INSERT INTO KUser (username, full_name, email) VALUES ($1, $2, $3) RETURNING user_id", nu.Username, nu.FullName, nu.Email).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func dbStorePassword(authDB *sql.DB, userID int64, encodedHash string) error {
	_, err := authDB.Exec("INSERT INTO credential (user_id, password_hash) VALUES ($1, $2)", userID, encodedHash)
	if err != nil {
		return err
	}
	return nil
}

func dbGetUserIDUsingFacebook(db *sql.DB, facebookID string) (int64, error) {
	var err error
	var userID int64
	err = db.QueryRow("SELECT user_id FROM KUser WHERE facebook=$1", facebookID).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func dbGetUserIDUsingGoogle(db *sql.DB, googleID string) (int64, error) {
	var err error
	var userID int64
	err = db.QueryRow("SELECT user_id FROM KUser WHERE google=$1", googleID).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func dbRegisterUserUsingFacebook(db *sql.DB, fu FacebookUser, username string) (int64, error) {
	var err error
	var userID int64
	err = db.QueryRow("INSERT INTO KUser (username, full_name, facebook) VALUES ($1, $2, $3) RETURNING user_id", username, fu.Name, fu.ID).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func dbRegisterUserUsingGoogle(db *sql.DB, cs *googlejwt.ClaimSet, username string) (int64, error) {
	var err error
	var userID int64
	err = db.QueryRow("INSERT INTO KUser (username, full_name, google) VALUES ($1, $2, $3) RETURNING user_id", username, cs.Name, cs.Sub).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
