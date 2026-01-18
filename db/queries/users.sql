
-- name: CreateUser :one
INSERT INTO users (
    email,
    first_name,
    last_name,
    password,
    role
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id, first_name, last_name, email, password, role, created_at;

-- name: GetUserByEmail :one
SELECT id, email, password, role, created_at 
FROM users 
WHERE email = $1 
LIMIT 1;

-- name: GetUserByID :one
SELECT id, email, first_name, last_name, password, role, created_at 
FROM users 
WHERE id = $1 
LIMIT 1;