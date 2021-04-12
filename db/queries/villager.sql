-- name: createVillager :one
INSERT INTO villagers (name, image, species, personality, birthday, quote) 
VALUES ($1, $2, $3, $4, $5, $6) 
RETURNING *;

-- name: getVillager :one
SELECT * FROM villagers
WHERE name = $1 
LIMIT 1;

-- name: getVillagers :many
SELECT * FROM villagers
ORDER BY name
LIMIT $1;