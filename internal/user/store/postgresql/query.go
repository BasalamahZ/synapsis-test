package postgresql

const queryCreateUser = `
	INSERT INTO
		user_info
	(
		email,
		name,
		password,
		phone_number,
		create_time
	) VALUES (
		:email,
		:name,
		:password,
		:phone_number,
		:create_time
	) RETURNING
 		id
`

const queryGetUser = `
	SELECT 
		u.id,
		u.email,
		u.name,
		u.password,
		u.phone_number,
		u.create_time,
		u.update_time
	FROM
		user_info u
	WHERE
		%s
`

const queryUpdateUser = `
	UPDATE
		user_info
	SET
		password = :password,
		update_time = :update_time
	WHERE
		id = :id
`
