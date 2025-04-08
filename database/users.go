package database

func GetUsernameByID(userID int) string {
    db := GetDatabase()
    var username string
    err := db.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&username)
    if err != nil {
        return ""
    }
    return username
}