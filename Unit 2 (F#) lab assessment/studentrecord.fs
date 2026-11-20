open Microsoft.Data.Sqlite

type Student = {
    Id: int
    Name: string
    Marks: int
}

let students = [
    { Id = 1; Name = "Shahbaz"; Marks = 80 }
    { Id = 2; Name = "Rahul"; Marks = 40 }
    { Id = 3; Name = "Harshit"; Marks = 70 }
]

printfn "Students"

for s in students do
    printfn "%s - %d" s.Name s.Marks

let names =
    students |> List.map (fun s -> s.Name)

printfn "\nNames"

for name in names do
    printfn "%s" name

let passed =
    students |> List.filter (fun s -> s.Marks >= 50)

printfn "\nPassed Students"

for s in passed do
    printfn "%s" s.Name

printfn "\nResult"

for s in students do
    match s.Marks with
    | x when x >= 50 -> printfn "%s = PASS" s.Name
    | _ -> printfn "%s = FAIL" s.Name

let db = new SqliteConnection("Data Source=students.db")
db.Open()

let cmd = db.CreateCommand()
cmd.CommandText <-
    "CREATE TABLE IF NOT EXISTS Students
     (Id INTEGER PRIMARY KEY, Name TEXT, Marks INTEGER)"

cmd.ExecuteNonQuery() |> ignore

let addStudent s =
    async {
        let c = db.CreateCommand()
        c.CommandText <- "INSERT INTO Students VALUES ($id,$name,$marks)"

        c.Parameters.AddWithValue("$id", s.Id) |> ignore
        c.Parameters.AddWithValue("$name", s.Name) |> ignore
        c.Parameters.AddWithValue("$marks", s.Marks) |> ignore

        c.ExecuteNonQuery() |> ignore
    }

addStudent students[0]
|> Async.RunSynchronously

let read = db.CreateCommand()
read.CommandText <- "SELECT * FROM Students"

use r = read.ExecuteReader()

printfn "\nDatabase"

while r.Read() do
    printfn "%d %s %d"
        (r.GetInt32(0))
        (r.GetString(1))
        (r.GetInt32(2))

let update = db.CreateCommand()
update.CommandText <- "UPDATE Students SET Marks = 90 WHERE Id = 1"
update.ExecuteNonQuery() |> ignore

let delete = db.CreateCommand()
delete.CommandText <- "DELETE FROM Students WHERE Id = 2"
delete.ExecuteNonQuery() |> ignore

printfn "\nCRUD completed"

db.Close()+