-- +goose Up
-- +goose StatementBegin
INSERT INTO public.status ("ID", "OrderStatus", "CalcStatus") 
                   VALUES (1, 'NEW', null)  on conflict do nothing;
INSERT INTO public.status ("ID", "OrderStatus", "CalcStatus") 
                   VALUES (2, 'PROCESSING', 'REGISTERED') on conflict do nothing;
INSERT INTO public.status ("ID", "OrderStatus", "CalcStatus") 
                   VALUES (3, 'PROCESSING', 'PROCESSING') on conflict do nothing;
INSERT INTO public.status ("ID", "OrderStatus", "CalcStatus") 
                   VALUES (4, 'INVALID', 'INVALID') on conflict do nothing;
INSERT INTO public.status ("ID", "OrderStatus", "CalcStatus") 
                   VALUES (5, 'PROCESSED', 'PROCESSED') on conflict do nothing;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM public.status;
-- +goose StatementEnd
