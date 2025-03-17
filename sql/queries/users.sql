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

-- name: SubscribeUser :exec
WITH subscribed_update AS (
    UPDATE users
    SET subscribed_to = array_append(subscribed_to, $1)
    WHERE id = $2
)
UPDATE users
SET subscribers = array_append(subscribers, $2)
WHERE users.id = $1;

-- name: UnsubscribeUser :exec
WITH unsubscribed_update AS (
    UPDATE users
    SET subscribed_to = array_remove(subscribed_to, $1)
    WHERE id = $2
)
UPDATE users
SET subscribers = array_remove(subscribers, $2)
WHERE users.id = $1;

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

-- name: VerifyVerificationCode :exec
UPDATE users 
SET verification_code = 0
WHERE id = $1;