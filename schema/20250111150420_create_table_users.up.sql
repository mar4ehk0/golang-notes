CREATE TABLE IF NOT EXISTS public.users (
    "id" serial PRIMARY KEY,
    "email" varchar(255) UNIQUE NOT NULL,
    "password" varchar(255) NOT NULL
);