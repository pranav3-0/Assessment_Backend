package repository

import "gorm.io/gorm"

type AuthRepositoryImpl struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &AuthRepositoryImpl{db: db}
}

type AuthRepository interface {
	SaveOTP(email string, otp string) error
	ValidateOTP(email string, otp string) (bool, error)
}

func (r *AuthRepositoryImpl) SaveOTP(email string, otp string) error {

	return r.db.Exec(`
		INSERT INTO email_otp_verification
		(email, otp, created_on, is_used)
		VALUES (?, ?, NOW(), false)
	`, email, otp).Error
}

func (r *AuthRepositoryImpl) ValidateOTP(email string, otp string) (bool, error) {

	var count int64

	err := r.db.Raw(`
		SELECT COUNT(*)
		FROM email_otp_verification
		WHERE email = ?
		AND otp = ?
		AND is_used = false
		AND created_on > NOW() - INTERVAL '5 minutes'
	`, email, otp).Scan(&count).Error

	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	r.db.Exec(`
		UPDATE email_otp_verification
		SET is_used = true
		WHERE email = ?
		AND otp = ?
	`, email, otp)

	return true, nil
}