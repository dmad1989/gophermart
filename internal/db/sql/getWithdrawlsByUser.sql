SELECT W."orderNum"::text AS ORDERNUM,
	W."pointsSum",
	W."processedDate"
FROM PUBLIC.WITHDRAWLS W
WHERE W.USERID = $1
ORDER BY W."processedDate";