CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(50) NOT NULL UNIQUE,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    password TEXT NOT NULL,
    role VARCHAR(20) DEFAULT 'CUSTOMER',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);