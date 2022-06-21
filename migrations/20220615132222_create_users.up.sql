CREATE TABLE IF NOT EXISTS public.users(
    id bigserial not null primary key,
    email varchar not null unique,
    encryptedPassword varchar not null,
    userRole int not null,
    Active boolean not null
);