package psql

const (
	createUserQuery = "INSERT INTO users(id, name, email, password, is_admin, create_date, update_date) VALUES ($1,$2,$3,$4,$5,NOW(),NOW())"
)
