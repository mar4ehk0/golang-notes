package dto

type TagNoteDto struct {
	TagID  int
	NoteID int
}

func NewTagsNotesForNote(tagsID []int, noteID int) []TagNoteDto {
	result := make([]TagNoteDto, 0, len(tagsID))

	for _, tagID := range tagsID {
		dto := TagNoteDto{TagID: tagID, NoteID: noteID}
		result = append(result, dto)
	}

	return result
}
