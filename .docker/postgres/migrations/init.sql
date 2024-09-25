-- CREATE DATABASE playground_app;

-- config `work_mem` for `playground_app` db
alter database playground_app set work_mem = '64MB';


-- construct generic_status enum
CREATE TYPE generic_status AS ENUM ('pending','completed','failed');

-- create balances table
CREATE TABLE public.balances (
    id CHAR(27) PRIMARY KEY,
    amount NUMERIC(15, 2),
    created_at TIMESTAMP DEFAULT current_timestamp,
    updated_at TIMESTAMP DEFAULT current_timestamp
);

-- setup balance seeder
INSERT INTO public.balances(id, amount)
VALUES
    ('2TeSprhp2cN6nEIcayZsjjvnlsK', 0),
    ('2TeSppzLaJaxldlzMOkqYO37vqw', 0),
    ('2TeSprCLB0tF6HsJT9eCb1IBht0', 0),
    ('2TeSps1ECSPx2IRrWMKEd6oHvSJ', 0),
    ('2TeSpqL5Cq3t0rNJ4RDcR4ky029', 0),
    ('2TeSpo6Okj6BedA62PUO3nRYADm', 0),
    ('2TeSprwUG0oLfgGbIIObipgf5be', 0),
    ('2TeSpqi0vPnYFf82cdBCE4Gaxjx', 0),
    ('2TeSppFijPxu2bIKpRLhLy3ye4j', 0),
    ('2TeSpnNHCn0Yq0vJddsv0dcSviy', 0),
    ('2TeSpsBbTYE5ZBe3zpO1HDRf44o', 0),
    ('2TeSprmZUc7HGpqlm7bUi6RB9lu', 0),
    ('2TeSpqKkgW9pz15PaqR8ZIS0Wjh', 0),
    ('2TeSpqfWmsmNLffPB4fhErgYRiC', 0),
    ('2TeSpnH3d9Hr9zPyNpoBrzdd8fx', 0),
    ('2TeSpsz6jeis3vJHCoSfFrnGPVl', 0),
    ('2TeSptEF7KEwZ9cdlowOH4HSmbR', 0),
    ('2TeSpsvjaufyDq78BwcFFZ6ibfi', 0),
    ('2TeSpsFLvMYhPU3CpC21wgyOcS5', 0),
    ('2TeSpmhKLpqgNXTNOvh3JPLkbjy', 0),
    ('2TWlPQ2AhstX9PtJ5UTOE6xQ7Ga', 0),
    ('2TWlPVmPjhonQ2DpOFt09O990th', 0),
    ('2TWlPdYWFbP3iXPMIKVRmdZ3ozC', 0),
    ('2TWlPkjd1YcRDCBDQk11nygVDpe', 0),
    ('2TWlPw3U64nhEtzeazP5ELd7q4c', 0),
    ('2TWlQ2JWBetv9MUkFsLPd4zhLa4', 0),
    ('2TWlQ83iEHKIOtvMCfSnBwM2sEB', 0),
    ('2TWlQJsEhRIW8XpeGGz2u75phWN', 0),
    ('2TWlQRnYNr5ViFi0wLTatYImXz5', 0),
    ('2TWlQYYZVIC566T3XFVldQckPsB', 0),
    ('2TWlQcFLvqE67A2qbZpWAdSJiZL', 0),
    ('2TWlQjcMepYBrRJRRzwgXWLA0gX', 0),
    ('2TWlQvGoGcHxUa4iod28bMsfW7e', 0),
    ('2TWlR5WFe7VUVMSxhgpEmfqlAAX', 0),
    ('2TWlR8ifwogFmktTwET0Eb2s4PE', 0),
    ('2TWlRFfCwsZo903aTO7xSRCYQIU', 0),
    ('2TWlRRtrVofuy7C1ZzcWIICVEME', 0),
    ('2TWlRZBGQlbQe34dVNU3GhQKshe', 0),
    ('2TWlRedfruxmFYvYLJur7oGesXY', 0),
    ('2TWlRnR3C0hopS7NkcwyjIOq5Kd', 0)
;

-- construct ledger_entry_type enum
CREATE TYPE ledger_entry_type AS ENUM ('in','out');

-- create ledgers table
CREATE TABLE public.ledgers (
    id CHAR(27) PRIMARY KEY,
    balance_id CHAR(27) NOT NULL REFERENCES balances,
    "type" ledger_entry_type,
    amount NUMERIC(15, 2),
    balance_after NUMERIC(15, 2),
    balance_before NUMERIC(15, 2),
    created_at TIMESTAMP DEFAULT current_timestamp,
    updated_at TIMESTAMP DEFAULT current_timestamp
);


-- create deposits table
CREATE TABLE deposits (
    id CHAR(27) PRIMARY KEY,
    balance_id CHAR(27) NOT NULL REFERENCES balances,
    "status" generic_status,
    amount NUMERIC(15, 2),
    created_at TIMESTAMP DEFAULT current_timestamp,
    updated_at TIMESTAMP DEFAULT current_timestamp
);

-- create withdrawals table
CREATE TABLE withdrawals (
    id CHAR(27) PRIMARY KEY,
    balance_id CHAR(27) NOT NULL REFERENCES balances,
    "status" generic_status,
    amount NUMERIC(15, 2),
    created_at TIMESTAMP DEFAULT current_timestamp,
    updated_at TIMESTAMP DEFAULT current_timestamp
);

-- create transfers table
CREATE TABLE transfers (
    id CHAR(27) PRIMARY KEY,
    balance_id_from CHAR(27) NOT NULL REFERENCES balances,
    balance_id_to CHAR(27) NOT NULL REFERENCES balances,
    "status" generic_status,
    amount NUMERIC(15, 2),
    created_at TIMESTAMP DEFAULT current_timestamp,
    updated_at TIMESTAMP DEFAULT current_timestamp
);

-- index transactions' status
CREATE INDEX idx_transfers_status ON public.transfers ("status");

-- index ledgers' type
CREATE INDEX idx_ledgers_type ON public.ledgers (type);