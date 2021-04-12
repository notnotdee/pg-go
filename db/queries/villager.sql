-- name: CreateVillager :one
INSERT INTO villagers (name, image, species, personality, birthday, quote) 
VALUES ($1, $2, $3, $4, $5, $6) 
RETURNING *;

-- name: GetVillager :one
SELECT * FROM villagers
WHERE name = $1 
LIMIT 1;

-- name: GetVillagers :many
SELECT * FROM villagers
ORDER BY name
LIMIT $1;