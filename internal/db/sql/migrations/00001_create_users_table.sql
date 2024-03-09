-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.users
(
    "ID" bigint  NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 1000000 ),
    login text NOT NULL,
    password text NOT NULL,
    pointsCurrent bigint,
    pointWithdrawn bigint,
    CONSTRAINT "PK_USER" PRIMARY KEY ("ID"),
    CONSTRAINT "UK_LOGIN" UNIQUE (login)
);

ALTER TABLE IF EXISTS public.users
    OWNER to postgres;
	
	
CREATE INDEX IF NOT EXISTS  "IDX_USER_LOGIN"
    ON public.users USING btree
    (login ASC NULLS LAST)
    WITH (deduplicate_items=True)
    TABLESPACE pg_default;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.users;
-- +goose StatementEnd
