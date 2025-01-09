package repositories

import (
    "CENT_Notes/cmd/models"
    "CENT_Notes/cmd/storage"
)

func CreateDirectory(dir models.Directory) (models.Directory, error) {
    db := storage.GetDB()
    sqlStatement := `
        INSERT INTO directories (name, user_id, parent_id, created_at, updated_at)
        VALUES ($1, $2, $3, NOW(), NOW())
        RETURNING id, created_at, updated_at`
    
    err := db.QueryRow(sqlStatement, dir.Name, dir.UserId, dir.ParentId).
        Scan(&dir.Id, &dir.CreatedAt, &dir.UpdatedAt)
    return dir, err
}

func GetUserDirectories(userId int) ([]models.Directory, error) {
    db := storage.GetDB()
    var directories []models.Directory
    
    rows, err := db.Query(`
        SELECT id, name, user_id, parent_id, created_at, updated_at
        FROM directories
        WHERE user_id = $1
        ORDER BY name`, userId)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var dir models.Directory
        err := rows.Scan(&dir.Id, &dir.Name, &dir.UserId, &dir.ParentId, 
            &dir.CreatedAt, &dir.UpdatedAt)
        if err != nil {
            return nil, err
        }
        directories = append(directories, dir)
    }
    return directories, nil
} 