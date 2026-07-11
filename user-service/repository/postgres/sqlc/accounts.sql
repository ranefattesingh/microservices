-- name: CreateAccount :one
INSERT INTO accounts (
    public_id, 
    email, 
    password_hash
) VALUES (
    $1, $2, $3
)
RETURNING id;

-- name: GetAccountWithID :one
SELECT * FROM accounts WHERE public_id = $1;

-- name: GetAccountWithEmail :one
SELECT * FROM accounts WHERE email = $1;

-- name: UpdateAccountEmailWithID :exec
UPDATE accounts
    SET email = $1
WHERE public_id = $2;

-- name: UpdatePasswordWithID :exec
UPDATE accounts
    SET password_hash = $1
WHERE public_id = $2;