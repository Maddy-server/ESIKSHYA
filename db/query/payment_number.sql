-- name: CreatePaymentNumber :exec
INSERT INTO payment_number(
    number,method,save,parent_id
) VALUES (
    ?,?,?,?
);

-- name: GetPaymentNumberInfo :one
SELECT * FROM payment_number WHERE parent_id=?;

