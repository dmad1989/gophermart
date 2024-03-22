SELECT O."number",
	S."CalcStatus"
FROM PUBLIC.ORDERS O
JOIN STATUS S ON S."ID" = O."statusId"
WHERE O."statusId" IN (1,2,3);