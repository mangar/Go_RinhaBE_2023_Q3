CREATE UNLOGGED TABLE IF NOT EXISTS rinha23_clientes
(
    id          varchar(100) not null Primary key,
    apelido     varchar(100),
    nome        varchar(100),
    nascimento  varchar(100),
    stack       varchar(1000),
    search_content       varchar(10000),
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

ALTER TABLE rinha23_clientes SET (autovacuum_enabled = false);


-- --
-- select count(*) from rinha23_clientes;

-- --
-- select * from rinha23_clientes;