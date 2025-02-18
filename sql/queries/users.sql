-- name: GetUserByEmailOrUsername :one
SELECT * FROM users
WHERE email = $1 OR username = $2;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1;

-- name: GetAllUsers :many
SELECT * FROM users;

-- name: ChangeUsername :one
UPDATE users
SET username = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ChangePassword :exec
UPDATE users
SET password = $2, updated_at = NOW()
WHERE id = $1;

-- name: DeleteUser :one
DELETE FROM users
WHERE id = $1
RETURNING *;
