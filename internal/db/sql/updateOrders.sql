UPDATE PUBLIC.ORDERS O
SET "statusId" = S."ID",
	ACCRUAL = :accrual
FROM PUBLIC.STATUS S
WHERE S."CalcStatus" = :calcstatus
	AND number = :number