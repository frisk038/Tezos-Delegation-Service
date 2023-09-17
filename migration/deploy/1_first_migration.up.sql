-- Create a table named 'delegations' to store information about delegations on the Tezos network.
CREATE TABLE delegations (
    id bigint PRIMARY KEY,       -- Unique identifier of the delegation
    ts TIMESTAMP NOT NULL,      -- Timestamp of the delegation (required and not null)
    amount bigint NOT NULL,     -- Amount delegated (required and not null)
    delegator text NOT NULL,    -- Delegator's address (required and not null)
    block text NOT NULL         -- Block identifier (required and not null)
);
