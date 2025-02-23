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

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: ResetPassword :exec
UPDATE users
SET password = $2, updated_at = NOW()
WHERE id = $1;

-- name: SendResetVerificationCode :exec
UPDATE users
SET verification_code = $2
WHERE id = $1;

-- name: VerifyVrificationCode :exec
UPDATE users 
SET verification_code = 0
WHERE id = $1;