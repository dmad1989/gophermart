-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.withdrawls
(
    "ID" bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 100000000 ),
    "orderNum" bigint,
    "pointsSum"  numeric(10,2),
    "processedDate" date,
    userid bigint,
    CONSTRAINT "PK_ID_WITHDRAWLS" PRIMARY KEY ("ID"),
    CONSTRAINT "FK_WITHDRAWLS_USERID" FOREIGN KEY (userid)
        REFERENCES public.users ("ID") MATCH SIMPLE
        ON UPDATE RESTRICT
        ON DELETE NO ACTION
        NOT VALID
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.withdrawls
    OWNER to postgres;


CREATE INDEX IF NOT EXISTS "IDX_WITHDRAWLS_USER"
    ON public.withdrawls USING btree
    (userid ASC NULLS LAST)
    WITH (deduplicate_items=True)
    TABLESPACE pg_default;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.withdrawls;
-- +goose StatementEnd
