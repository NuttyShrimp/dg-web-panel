package staff

import (
	"degrens/panel/internal/db"
	panel_models "degrens/panel/internal/db/models/panel"
	"degrens/panel/lib/graylogger"

	"gorm.io/gorm"
)

func GetAllNotes() ([]panel_models.Notes, error) {
	notes := []panel_models.Notes{}
	err := db.MariaDB.Client.Preload("User").Order("updated_at DESC").Find(&notes).Error
	if err == gorm.ErrRecordNotFound {
		return notes, nil
	}
	return notes, err
}

func CreateNote(userId uint, note string) error {
	dbNote := panel_models.Notes{
		Note:      note,
		CreatorID: userId,
	}
	err := db.MariaDB.Client.Save(&dbNote).Error
	if err == nil {
		graylogger.Log("staff:notes:create", "userId", userId, "note", note)
	}
	return err
}

func UpdateNote(userId uint, noteId uint, note string) error {
	dbNote := panel_models.Notes{}
	orgNote := string(dbNote.Note)
	err := db.MariaDB.Client.First(&dbNote, noteId).Error
	if err != nil {
		return err
	}
	dbNote.Note = note
	err = db.MariaDB.Client.Save(&dbNote).Error
	if err == nil {
		graylogger.Log("staff:notes:update", "userId", userId, "originalNote", orgNote, "newNote", note, "noteStruct", dbNote)
	}
	return err
}

func DeleteNote(userId uint, noteId uint) error {
	dbNote := panel_models.Notes{}
	err := db.MariaDB.Client.First(&dbNote, noteId).Error
	if err != nil {
		return err
	}
	err = db.MariaDB.Client.Where("id = ?", noteId).Delete(&panel_models.Notes{}).Error
	if err == nil {
		graylogger.Log("staff:notes:delete", "userId", userId, "note", dbNote)
	}
	return err
}
