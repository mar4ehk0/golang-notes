CREATE TABLE IF NOT EXISTS public.tags_notes (
    "id" serial PRIMARY KEY,
    "tag_id" int NOT NULL,
    "note_id" int NOT NULL
);
ALTER TABLE public.tags_notes ADD CONSTRAINT tags_notes_tag_fk FOREIGN KEY (tag_id) REFERENCES public.tags("id");
ALTER TABLE public.tags_notes ADD CONSTRAINT tags_notes_note_fk FOREIGN KEY (note_id) REFERENCES public.notes("id");