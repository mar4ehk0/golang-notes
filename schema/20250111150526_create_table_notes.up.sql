CREATE TABLE IF NOT EXISTS public.notes (
    "id" serial PRIMARY KEY,
    "title" varchar(255) NOT NULL,
    "body" text NOT NULL,
    "user_id" int
);
ALTER TABLE public.notes ADD CONSTRAINT notes_users_fk FOREIGN KEY (user_id) REFERENCES public.users("id");