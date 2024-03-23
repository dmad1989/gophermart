select
	U."ID",
	U.LOGIN as login,
	U.PASSWORD as password
from
	users u
where
	u.login = $1