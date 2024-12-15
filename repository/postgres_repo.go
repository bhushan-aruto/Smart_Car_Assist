package repository

import "database/sql"

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{
		db,
	}
}

func (repo *PostgresRepo) GetUserIdByEmail(email string) (string, error) {
	query := `select user_id from users where email=$1`
	var userId string
	err := repo.db.QueryRow(query, email).Scan(&userId)
	return userId, err
}

func (repo *PostgresRepo) CheckUserBookingExists(userId string) (bool, error) {
	query := `select exists(select 1 from bookings where user_id = $1)`
	var bookingExists bool

	err := repo.db.QueryRow(query, userId).Scan(&bookingExists)

	return bookingExists, err
}

func (repo *PostgresRepo) DeleteUserBooking(userId string) error {
	query := `delete from bookings where user_id = $1`

	_, err := repo.db.Exec(query, userId)
	return err
}
