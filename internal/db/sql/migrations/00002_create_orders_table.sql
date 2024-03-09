-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.orders
(
    "ID" bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 100000000 ),
    "number" bigint NOT NULL,
    userid bigint NOT NULL,
    "uploadDate" date NOT NULL,
    status text,
    accrual bigint,
    CONSTRAINT "PK_ID" PRIMARY KEY ("ID"),
    CONSTRAINT "UK_NUMBER" UNIQUE ("number"),
    CONSTRAINT "FK_USERID" FOREIGN KEY (userid)
        REFERENCES public.users ("ID") MATCH SIMPLE
        ON UPDATE RESTRICT
        ON DELETE NO ACTION
        NOT VALID
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.orders
    OWNER to postgres;


CREATE INDEX IF NOT EXISTS "IDX_ORDER_USER"
    ON public.orders USING btree
    (userid ASC NULLS LAST)
    WITH (deduplicate_items=False)
    TABLESPACE pg_default;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.orders;
-- +goose StatementEnd
