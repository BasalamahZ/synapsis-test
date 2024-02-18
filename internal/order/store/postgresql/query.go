package postgresql

const queryCreateOrder = `
	INSERT INTO
		transaction
	(
		user_id,
		product_id,
		total_amount,
		quantity,
		response_midtrans,
		create_time
	) VALUES (
		:user_id,
		:product_id,
		:total_amount,
		:quantity,
		:response_midtrans,
		:create_time
	) RETURNING
		id
`


const queryGetOrder = `
	SELECT
		t.id,
		t.user_id,
		ui.name AS user_name,
		ui.email AS user_email,
		t.product_id,
		p.name AS product_name,
		p.price AS product_price,
		t.quantity,
		t.total_amount,
		t.status,
		t.response_midtrans,
		t.create_time,
		t.update_time
	FROM
		transaction t
	LEFT JOIN
		product_cart pc
	ON
		t.user_id = pc.user_id
	AND
		t.product_id = pc.product_id
	LEFT JOIN
		product p
	ON
		p.id = t.product_id
	LEFT JOIN
		user_info ui
	ON
		ui.id = t.user_id
	%s
`
