CREATE UNLOGGED TABLE IF NOT EXISTS rinha23_clientes
(
    id          varchar(100) not null Primary key,
    apelido     varchar(500),
    nome        varchar(500),
    nascimento  varchar(100),
    stack       varchar(10000),
    search_content       varchar(100000),
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

ALTER TABLE rinha23_clientes SET (autovacuum_enabled = false);


-- --
-- select count(*) from rinha23_clientes;

-- --
-- select * from rinha23_clientes;