create unique index if not exists client_exposure_idx on client_exposure (client_id);
create unique index if not exists client_exposure_detail_idx on client_exposure_detail (client_id, exposure_currency);
create unique index if not exists currency_rate_idx on currency_rate (base_currency, quoted_currency);