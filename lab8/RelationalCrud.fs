open System
open Microsoft.Data.Sqlite

// Connect to SQLite Database
let conn = new SqliteConnection("Data Source=student.db")
conn.Open()

// Helper function to run SQL command
let runSql (sql: string) =
    let cmd = conn.CreateCommand()
    cmd.CommandText <- sql
    cmd.ExecuteNonQuery() |> ignore

// Helper function to print table
let showStudents () =
    printfn "--- Student Table ---"
    let cmd = conn.CreateCommand()
    cmd.CommandText <- "SELECT * FROM Students"
    use r = cmd.ExecuteReader()
    while r.Read() do
        printfn "Id: %O | Name: %O | Marks: %O" r.[0] r.[1] r.[2]

// Create Table
runSql "CREATE TABLE IF NOT EXISTS Students (Id INT, Name TEXT, Marks INT)"
runSql "DELETE FROM Students" // Clear table for fresh run

// 1. CREATE
printfn "\n1. CREATE"
runSql "INSERT INTO Students VALUES (1, 'Alice', 85)"
runSql "INSERT INTO Students VALUES (2, 'Bob', 90)"
showStudents ()

// 2. READ
printfn "\n2. READ (Student 1)"
let cmd = conn.CreateCommand()
cmd.CommandText <- "SELECT * FROM Students WHERE Id = 1"
use r = cmd.ExecuteReader()
if r.Read() then
    printfn "Found Student: %O, Marks: %O" r.[1] r.[2]

// 3. UPDATE
printfn "\n3. UPDATE (Bob Marks to 95)"
runSql "UPDATE Students SET Marks = 95 WHERE Id = 2"
showStudents ()

// 4. DELETE
printfn "\n4. DELETE (Delete Alice)"
runSql "DELETE FROM Students WHERE Id = 1"
showStudents ()

conn.Close()
