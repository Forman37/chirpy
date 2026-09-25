-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password, is_chirpy_red)
VALUES(
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    false
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: FetchUser :one
SELECT * from users
WHERE email = $1;

-- name: GetUserFromRefreshToken :one
SELECT * FROM users
WHERE id = $1;

-- name: UpdateUser :exec
UPDATE users
SET email = $1, hashed_password = $2
WHERE id = $3;

-- name: UpdateToChirpyRed :exec
UPDATE users
SET is_chirpy_red = true
WHERE id = $1;
