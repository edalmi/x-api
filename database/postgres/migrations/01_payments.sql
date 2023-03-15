CREATE TABLE accounts (
       id BIGSERIAL NOT NULL,
       PRIMARY KEY(id)
);

CREATE TABLE subscription_type AS ENUM('SENDER', 'RECEIVER')

CREATE TABLE account_subscriptions(
       id BIGSERIAL NOT NULL,
       account_id BIGINT NOT NULL,
       subscription_type subscription_type NOT NULL,
       PRIMARY KEY(id),
       FOREIGN KEY(account_id) REFERENCES accounts(id),
       UNIQUE(account_id, subscription_type)
);

CREATE TYPE payment_status AS ENUM ('CREATED', 'FAILED');

CREATE TABLE payments (
           id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
           third_party_id VARCHAR(128) NOT NULL,
           amount BIGINT NOT NULL,
           currency VARCHAR(3) NOT NULL,
           status payment_status NOT NULL DEFAULT 'CREATED',
           account_id BIGINT NOT NULL,
           created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
           updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
           UNIQUE(third_party_id, account_id),
           FOREIGN KEY (account_id) REFERENCES accounts(id)
);

