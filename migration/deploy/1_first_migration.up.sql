CREATE TABLE delegations (
id bigint PRIMARY KEY, -- unique identifier of the delegation
ts TIMESTAMP NOT NULL, -- Timestamp of the delegation
amount bigint NOT NULL, -- amount delegated
delegator text NOT NULL, -- delegator address
block text NOT NULL-- block identifier
);
