package postgresql

const queryGetProduct = `
	SELECT
		p.id,
		p.name,
		p.price,
		p.description,
		p.category_id,
		c.name AS category_name,
		p.create_time,
		p.update_time
	FROM
		product p
	LEFT JOIN
		category c
	ON
		c.id = p.category_id
	%s
`

const queryAddProductCart = `
	INSERT INTO
		product_cart
	(
		user_id,
		product_id,
		quantity,
		create_time
	) VALUES (
		:user_id,
		:product_id,
		:quantity,
		:create_time
	)
`

const queryGetProductCart = `
	SELECT
		pc.user_id,
		pc.product_id,
		p.name AS product_name,
		p.price AS product_price,
		pc.quantity,
		pc.create_time,
		pc.update_time
	FROM
		product_cart pc
	LEFT JOIN
		product p
	ON
		p.id = pc.product_id
	%s
`

const queryDeleteProductCart = `
	DELETE FROM
		product_cart
	WHERE
		user_id = :user_id
	AND
		product_id = :product_id
`
