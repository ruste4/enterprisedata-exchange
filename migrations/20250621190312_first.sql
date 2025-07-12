-- +goose Up
CREATE TABLE IF NOT EXISTS exchange_nodes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,                                    
    c_description TEXT NOT NULL,                                             
    node_code VARCHAR(255) NOT NULL,                                         
    this_node_code VARCHAR(255),                                             
    prefix VARCHAR(2),                                                       
    this_prefix VARCHAR(2),                                                  
    c_state VARCHAR(50) CHECK(c_state IN ('active', 'inactive', 'deleted')), 
    user VARCHAR(50),
    pass VARCHAR(255),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,                  
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP                   
);

CREATE TABLE IF NOT EXISTS exchange_sessions (
    session_uuid VARCHAR(50) NOT NULL,
    exchange_nodes_ref INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,                  
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP 
);

-- +goose Down
DROP TABLE exchange_nodes;
DROP TABLE exchange_sessions;