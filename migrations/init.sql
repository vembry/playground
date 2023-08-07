-- CREATE DATABASE credit;

-- create users table
CREATE TABLE public.users (
    id CHAR(27) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp
);

-- setup dummy users
INSERT INTO public.users (id, name)
VALUES 
    ('2TWlPQ2AhstX9PtJ5UTOE6xQ7Ga', 'user_1'),
    ('2TWlPVmPjhonQ2DpOFt09O990th', 'user_2'),
    ('2TWlPdYWFbP3iXPMIKVRmdZ3ozC', 'user_3'),
    ('2TWlPkjd1YcRDCBDQk11nygVDpe', 'user_4'),
    ('2TWlPw3U64nhEtzeazP5ELd7q4c', 'user_5'),
    ('2TWlQ2JWBetv9MUkFsLPd4zhLa4', 'user_6'),
    ('2TWlQ83iEHKIOtvMCfSnBwM2sEB', 'user_7'),
    ('2TWlQJsEhRIW8XpeGGz2u75phWN', 'user_8'),
    ('2TWlQRnYNr5ViFi0wLTatYImXz5', 'user_9'),
    ('2TWlQYYZVIC566T3XFVldQckPsB', 'user_10'),
    ('2TWlQcFLvqE67A2qbZpWAdSJiZL', 'user_11'),
    ('2TWlQjcMepYBrRJRRzwgXWLA0gX', 'user_12'),
    ('2TWlQvGoGcHxUa4iod28bMsfW7e', 'user_13'),
    ('2TWlR5WFe7VUVMSxhgpEmfqlAAX', 'user_14'),
    ('2TWlR8ifwogFmktTwET0Eb2s4PE', 'user_15'),
    ('2TWlRFfCwsZo903aTO7xSRCYQIU', 'user_16'),
    ('2TWlRRtrVofuy7C1ZzcWIICVEME', 'user_17'),
    ('2TWlRZBGQlbQe34dVNU3GhQKshe', 'user_18'),
    ('2TWlRedfruxmFYvYLJur7oGesXY', 'user_19'),
    ('2TWlRnR3C0hopS7NkcwyjIOq5Kd', 'user_20')
;

-- construct transaction_status enum
CREATE TYPE transaction_status AS ENUM ('pending','success','failed');

-- create transactions table
CREATE TABLE public.transactions (
    id CHAR(27) PRIMARY KEY,
    user_id CHAR(27),
    status transaction_status,
    description VARCHAR(50) NOT NULL,
    remarks VARCHAR(50) NOT NULL,
    amount NUMERIC(15, 2),
    created_at TIMESTAMP DEFAULT current_timestamp,
    updated_at TIMESTAMP DEFAULT current_timestamp,

    FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT unique_id_userid UNIQUE(id, user_id)
);

-- index transactions' status
CREATE INDEX idx_transactions_status ON public.transactions (status);

-- construct ledger_entry_type enum
CREATE TYPE ledger_entry_type AS ENUM ('in','out');

-- create ledgers table
CREATE TABLE public.ledgers (
    id CHAR(27) PRIMARY KEY,
    user_id CHAR(27),
    type ledger_entry_type,
    description VARCHAR(50) NOT NULL,
    amount NUMERIC(15, 2),
    balance_after NUMERIC(15, 2),
    balance_before NUMERIC(15, 2),
    created_at TIMESTAMP DEFAULT current_timestamp,

    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- index ledgers' type
CREATE INDEX idx_ledgers_type ON public.ledgers (type);

-- setup pre-balance
INSERT INTO public.ledgers(id, user_id, type, description, amount, balance_after, balance_before)
VALUES
    ('2TWoNoIElOA5lEMz2YJ4hBtWKEW', '2TWlPQ2AhstX9PtJ5UTOE6xQ7Ga', 'in', 'init balance', 0, 0, 0),
    ('2TWoNk4xFevVLpBa9CGcOlJMMtD', '2TWlPVmPjhonQ2DpOFt09O990th', 'in', 'init balance', 0, 0, 0),
    ('2TWoNjy0ZwMeCVJiSbUdEGD2rOF', '2TWlPdYWFbP3iXPMIKVRmdZ3ozC', 'in', 'init balance', 0, 0, 0),
    ('2TWoNlkQxwgYo7mK5SmHKhGgn7H', '2TWlPkjd1YcRDCBDQk11nygVDpe', 'in', 'init balance', 0, 0, 0),
    ('2TWoNjNrMz5Tk7KWrUoyzWd9n1x', '2TWlPw3U64nhEtzeazP5ELd7q4c', 'in', 'init balance', 0, 0, 0),
    ('2TWoNlsyvef315nycZp2Fo38ir2', '2TWlQ2JWBetv9MUkFsLPd4zhLa4', 'in', 'init balance', 0, 0, 0),
    ('2TWoNkSFocY94EtT1bzah03YW7M', '2TWlQ83iEHKIOtvMCfSnBwM2sEB', 'in', 'init balance', 0, 0, 0),
    ('2TWoNjXHN5wO1mCkKHWrmeqHXwl', '2TWlQJsEhRIW8XpeGGz2u75phWN', 'in', 'init balance', 0, 0, 0),
    ('2TWoNnfgdYCl6RP8uuId2UxiPV4', '2TWlQRnYNr5ViFi0wLTatYImXz5', 'in', 'init balance', 0, 0, 0),
    ('2TWoNkrfVNVHZJiuRQse59VuKK9', '2TWlQYYZVIC566T3XFVldQckPsB', 'in', 'init balance', 0, 0, 0),
    ('2TWoNkcUykDQ4sBOAy3IwLgGz0t', '2TWlQcFLvqE67A2qbZpWAdSJiZL', 'in', 'init balance', 0, 0, 0),
    ('2TWoNilW0BrGwmP95XZlIpN6zEg', '2TWlQjcMepYBrRJRRzwgXWLA0gX', 'in', 'init balance', 0, 0, 0),
    ('2TWoNmxDzxJvxVwxv7MQleJ5XL3', '2TWlQvGoGcHxUa4iod28bMsfW7e', 'in', 'init balance', 0, 0, 0),
    ('2TWoNju27StnXldwLZwKvhEV53V', '2TWlR5WFe7VUVMSxhgpEmfqlAAX', 'in', 'init balance', 0, 0, 0),
    ('2TWoNm6hVYSFmLA87XBe9UPU8NT', '2TWlR8ifwogFmktTwET0Eb2s4PE', 'in', 'init balance', 0, 0, 0),
    ('2TWoNiySHd5cook6r4qWbpWBmmh', '2TWlRFfCwsZo903aTO7xSRCYQIU', 'in', 'init balance', 0, 0, 0),
    ('2TWoNmYo3Tnc9isXPcL6xKRzvlp', '2TWlRRtrVofuy7C1ZzcWIICVEME', 'in', 'init balance', 0, 0, 0),
    ('2TWoNlmcEqb6pM0niSL1WgViSsD', '2TWlRZBGQlbQe34dVNU3GhQKshe', 'in', 'init balance', 0, 0, 0),
    ('2TWoNiKJghuwkEANCZ1lavtOBMw', '2TWlRedfruxmFYvYLJur7oGesXY', 'in', 'init balance', 0, 0, 0),
    ('2TWoNiinHa0FTsSPIjzMvIww10k', '2TWlRnR3C0hopS7NkcwyjIOq5Kd', 'in', 'init balance', 0, 0, 0)
;

-- create balances table
CREATE TABLE public.balances (
    id CHAR(27) PRIMARY KEY,
    user_id CHAR(27),
    amount NUMERIC(15, 2),
    created_at TIMESTAMP DEFAULT current_timestamp,
    updated_at TIMESTAMP DEFAULT current_timestamp,
    latest_ledger_point_id CHAR(27),

    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (latest_ledger_point_id) REFERENCES ledgers (id)
);

-- setup pre-balance
INSERT INTO public.balances(id, user_id, amount, latest_ledger_point_id)
VALUES
    ('2TeSprhp2cN6nEIcayZsjjvnlsK', '2TWlPQ2AhstX9PtJ5UTOE6xQ7Ga', 0, '2TWoNoIElOA5lEMz2YJ4hBtWKEW'),
    ('2TeSppzLaJaxldlzMOkqYO37vqw', '2TWlPVmPjhonQ2DpOFt09O990th', 0, '2TWoNk4xFevVLpBa9CGcOlJMMtD'),
    ('2TeSprCLB0tF6HsJT9eCb1IBht0', '2TWlPdYWFbP3iXPMIKVRmdZ3ozC', 0, '2TWoNjy0ZwMeCVJiSbUdEGD2rOF'),
    ('2TeSps1ECSPx2IRrWMKEd6oHvSJ', '2TWlPkjd1YcRDCBDQk11nygVDpe', 0, '2TWoNlkQxwgYo7mK5SmHKhGgn7H'),
    ('2TeSpqL5Cq3t0rNJ4RDcR4ky029', '2TWlPw3U64nhEtzeazP5ELd7q4c', 0, '2TWoNjNrMz5Tk7KWrUoyzWd9n1x'),
    ('2TeSpo6Okj6BedA62PUO3nRYADm', '2TWlQ2JWBetv9MUkFsLPd4zhLa4', 0, '2TWoNlsyvef315nycZp2Fo38ir2'),
    ('2TeSprwUG0oLfgGbIIObipgf5be', '2TWlQ83iEHKIOtvMCfSnBwM2sEB', 0, '2TWoNkSFocY94EtT1bzah03YW7M'),
    ('2TeSpqi0vPnYFf82cdBCE4Gaxjx', '2TWlQJsEhRIW8XpeGGz2u75phWN', 0, '2TWoNjXHN5wO1mCkKHWrmeqHXwl'),
    ('2TeSppFijPxu2bIKpRLhLy3ye4j', '2TWlQRnYNr5ViFi0wLTatYImXz5', 0, '2TWoNnfgdYCl6RP8uuId2UxiPV4'),
    ('2TeSpnNHCn0Yq0vJddsv0dcSviy', '2TWlQYYZVIC566T3XFVldQckPsB', 0, '2TWoNkrfVNVHZJiuRQse59VuKK9'),
    ('2TeSpsBbTYE5ZBe3zpO1HDRf44o', '2TWlQcFLvqE67A2qbZpWAdSJiZL', 0, '2TWoNkcUykDQ4sBOAy3IwLgGz0t'),
    ('2TeSprmZUc7HGpqlm7bUi6RB9lu', '2TWlQjcMepYBrRJRRzwgXWLA0gX', 0, '2TWoNilW0BrGwmP95XZlIpN6zEg'),
    ('2TeSpqKkgW9pz15PaqR8ZIS0Wjh', '2TWlQvGoGcHxUa4iod28bMsfW7e', 0, '2TWoNmxDzxJvxVwxv7MQleJ5XL3'),
    ('2TeSpqfWmsmNLffPB4fhErgYRiC', '2TWlR5WFe7VUVMSxhgpEmfqlAAX', 0, '2TWoNju27StnXldwLZwKvhEV53V'),
    ('2TeSpnH3d9Hr9zPyNpoBrzdd8fx', '2TWlR8ifwogFmktTwET0Eb2s4PE', 0, '2TWoNm6hVYSFmLA87XBe9UPU8NT'),
    ('2TeSpsz6jeis3vJHCoSfFrnGPVl', '2TWlRFfCwsZo903aTO7xSRCYQIU', 0, '2TWoNiySHd5cook6r4qWbpWBmmh'),
    ('2TeSptEF7KEwZ9cdlowOH4HSmbR', '2TWlRRtrVofuy7C1ZzcWIICVEME', 0, '2TWoNmYo3Tnc9isXPcL6xKRzvlp'),
    ('2TeSpsvjaufyDq78BwcFFZ6ibfi', '2TWlRZBGQlbQe34dVNU3GhQKshe', 0, '2TWoNlmcEqb6pM0niSL1WgViSsD'),
    ('2TeSpsFLvMYhPU3CpC21wgyOcS5', '2TWlRedfruxmFYvYLJur7oGesXY', 0, '2TWoNiKJghuwkEANCZ1lavtOBMw'),
    ('2TeSpmhKLpqgNXTNOvh3JPLkbjy', '2TWlRnR3C0hopS7NkcwyjIOq5Kd', 0, '2TWoNiinHa0FTsSPIjzMvIww10k')
;