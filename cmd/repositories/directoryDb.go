package repositories

import (
    "CENT_Notes/cmd/models"
    "CENT_Notes/cmd/storage"
    "fmt"
    "database/sql"
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
    
    fmt.Printf("Executing query for user ID: %d\n", userId)
    
    rows, err := db.Query(`
        SELECT id, name, user_id, parent_id, created_at, updated_at
        FROM directories
        WHERE user_id = $1
        ORDER BY name`, userId)
    if err != nil {
        fmt.Printf("Query error: %v\n", err)
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var dir models.Directory
        err := rows.Scan(&dir.Id, &dir.Name, &dir.UserId, &dir.ParentId, 
            &dir.CreatedAt, &dir.UpdatedAt)
        if err != nil {
            fmt.Printf("Scan error: %v\n", err)
            return nil, err
        }
        directories = append(directories, dir)
    }
    
    fmt.Printf("Found %d directories in database\n", len(directories))
    return directories, nil
}

func GetDirectory(directoryId int) (models.Directory, error) {
    db := storage.GetDB()
    var directory models.Directory
    
    sqlStatement := `
        SELECT id, name, user_id, parent_id, created_at, updated_at
        FROM directories
        WHERE id = $1`
    
    err := db.QueryRow(sqlStatement, directoryId).Scan(
        &directory.Id,
        &directory.Name,
        &directory.UserId,
        &directory.ParentId,
        &directory.CreatedAt,
        &directory.UpdatedAt,
    )
    
    if err == sql.ErrNoRows {
        return directory, fmt.Errorf("directory not found")
    }
    
    if err != nil {
        return directory, fmt.Errorf("error fetching directory: %v", err)
    }
    
    return directory, nil
} 