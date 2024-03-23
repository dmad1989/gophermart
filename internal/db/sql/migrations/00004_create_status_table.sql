-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.status
(
    "ID" bigint NOT NULL,
    "OrderStatus" text COLLATE pg_catalog."default" NOT NULL,
    "CalcStatus" text COLLATE pg_catalog."default",
    CONSTRAINT status_pkey PRIMARY KEY ("ID")
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.status
    OWNER to postgres;

CREATE INDEX IF NOT EXISTS "IDX_ID"
    ON public.status USING btree
    ("ID" ASC NULLS LAST)
    WITH (deduplicate_items=True)
    TABLESPACE pg_default;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.status;
-- +goose StatementEnd
