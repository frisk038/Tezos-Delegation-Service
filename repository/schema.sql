CREATE TABLE delegations (
id bigint PRIMARY KEY, -- unique identifier of the delegation
ts TIMESTAMP NOT NULL, -- Timestamp of the delegation
amount text, -- amount delegated
delegator text, -- delegator address
block text, -- block identifier
);