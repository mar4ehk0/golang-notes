CREATE TABLE IF NOT EXISTS public.tags (
    "id" serial PRIMARY KEY,
    "name" varchar(255) UNIQUE NOT NULL
);