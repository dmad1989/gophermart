SELECT O."number"::text as number,
	S."OrderStatus",
	O.ACCRUAL,
	O."uploadDate"
FROM PUBLIC.ORDERS O
JOIN PUBLIC.STATUS S ON S."ID" = O."statusId"
WHERE O.USERID = $1
order by w."uploadDate";