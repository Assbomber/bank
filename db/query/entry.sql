-- name: CreateEntry :one
INSERT INTO entries (
    account_id,
    amount
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetEntry :one
SELECT *
FROM entries
WHERE id = $1
LIMIT 1;


-- name: ListEntries :many
SELECT * 
FROM entries
WHERE account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;




/*
    INSERT INTO transfers (from_account_id,to_account_id,amount) VALUES (1,2,10);

    INSERT INTO entries (account_id,amount) VALUES (1,-10);
    INSERT INTO entries (account_id,amount) VALUES (2,10);

    SELECT * FROM accounts WHERE id = 1 LIMIT 1 FOR UPDATE;
    SELECT * FROM accounts WHERE id = 2 LIMIT 1 FOR UPDATE;

    UPDATE accounts SET balance = 10 WHERE account_id = 1;
    UPDATE accounts SET balance = 10 WHERE account_id = 2;

*/