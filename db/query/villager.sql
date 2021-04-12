-- name: createVillager :one
INSERT INTO villagers (villager) 
VALUES ($1) 
RETURNING *;

-- name: getVillager :one
SELECT * FROM villagers
WHERE villager = $1 
LIMIT 1;
