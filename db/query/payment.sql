-- name: CreatePayment :exec
INSERT INTO payment(
    transaction_id,transaction_token,method,parent_id,child_id,amount,pay_at,expire_at
) VALUES (
    ?,?,?,?,?,?,?,?
);

-- name: GetPaymentList :many
SELECT * from payment WHERE parent_id=? AND expire_at >= ? ;

-- name: GetPayment :one
SELECT * from payment WHERE child_id=? AND expire_at >= ? ;