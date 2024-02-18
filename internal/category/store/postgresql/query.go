package postgresql

const queryGetCategory = `
	SELECT
		c.id,
		c.name,
		c.description,
		c.create_time,
		c.update_time
	FROM
		category c
	%s
`
