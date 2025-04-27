-- Create the "todos" table
CREATE TABLE todos (
                       id BIGSERIAL PRIMARY KEY,              -- Auto-incrementing primary key
                       description VARCHAR(255),  -- Description with a maximum length of 255 characters
                       status VARCHAR(255) DEFAULT 'TO DO',       -- Status with a maximum length of 255 characters
                       created BIGINT, -- Created timestamp
                       updated BIGINT -- Updated timestamp
);

-- Create an index on the "updated" column if you plan to sort/filter by it often
CREATE INDEX idx_id ON todos(id);
CREATE INDEX idx_status ON todos(status);
CREATE INDEX idx_created ON todos(created);
CREATE INDEX idx_updated ON todos(updated);