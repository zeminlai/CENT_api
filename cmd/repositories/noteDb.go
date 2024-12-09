package repositories

import (
    "CENT_Notes/cmd/models"
    "CENT_Notes/cmd/storage"
)

// CreateNote inserts a new note into the database
func CreateNote(note models.Note) (models.Note, error) {
    db := storage.GetDB()
    sqlStatement := `
        INSERT INTO notes (title, content, user_id, created_at, updated_at) 
        VALUES ($1, $2, $3, NOW(), NOW()) 
        RETURNING id, created_at, updated_at`
    
    err := db.QueryRow(sqlStatement, note.Title, note.Content, note.UserId).
        Scan(&note.Id, &note.CreatedAt, &note.UpdatedAt)
    if err != nil {
        return note, err
    }
    return note, nil
}

// GetNote retrieves a note from the database by ID
func GetNote(id int) (models.Note, error) {
    db := storage.GetDB()
    var note models.Note
    sqlStatement := `
        SELECT id, title, content, user_id, created_at, updated_at 
        FROM notes 
        WHERE id = $1`
    
    err := db.QueryRow(sqlStatement, id).
        Scan(&note.Id, &note.Title, &note.Content, &note.UserId, &note.CreatedAt, &note.UpdatedAt)
    if err != nil {
        return note, err
    }
    return note, nil
}

// GetUserNotes retrieves all notes for a specific user
func GetUserNotes(userId int) ([]models.Note, error) {
    db := storage.GetDB()
    var notes []models.Note
    
    sqlStatement := `
        SELECT id, title, content, user_id, created_at, updated_at 
        FROM notes 
        WHERE user_id = $1 
        ORDER BY created_at DESC`
    
    rows, err := db.Query(sqlStatement, userId)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var note models.Note
        err := rows.Scan(&note.Id, &note.Title, &note.Content, &note.UserId, 
            &note.CreatedAt, &note.UpdatedAt)
        if err != nil {
            return nil, err
        }
        notes = append(notes, note)
    }
    return notes, nil
}

// UpdateNote updates an existing note in the database
func UpdateNote(note models.Note) (models.Note, error) {
    db := storage.GetDB()
    sqlStatement := `
        UPDATE notes 
        SET title = $1, content = $2, updated_at = NOW() 
        WHERE id = $3 AND user_id = $4 
        RETURNING updated_at`
    
    err := db.QueryRow(sqlStatement, note.Title, note.Content, note.Id, note.UserId).
        Scan(&note.UpdatedAt)
    if err != nil {
        return note, err
    }
    return note, nil
}

// DeleteNote deletes a note from the database
func DeleteNote(id int, userId int) error {
    db := storage.GetDB()
    sqlStatement := `DELETE FROM notes WHERE id = $1 AND user_id = $2`
    result, err := db.Exec(sqlStatement, id, userId)
    if err != nil {
        return err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if rowsAffected == 0 {
        return sql.ErrNoRows
    }
    return nil
} 