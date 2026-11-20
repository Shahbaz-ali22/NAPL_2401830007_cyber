// ============================================================================
// Lab 8 (Exp 15): Simple Database CRUD Operations (Relational & NoSQL)
// Zero-dependency pure F# script (runs directly via: dotnet fsi SimpleDatabaseCrud.fsx)
// ============================================================================

open System

// ============================================================================
// PART 1: Relational Database (Table with Fixed Schema: Id, Name, Marks)
// ============================================================================
printfn "=========================================================="
printfn " PART 1: RELATIONAL DATABASE (SQL / SQLite Model)"
printfn "=========================================================="

// Define fixed entity schema
type Student = { Id: int; Name: string; Marks: int }

// Relational Table represented as a mutable list of rows
let mutable studentTable: Student list = []

let showTable (step: string) =
    printfn "\n--- Table State after: %s ---" step
    if studentTable.IsEmpty then
        printfn "(Table is empty)"
    else
        for s in studentTable do
            printfn "Row -> ID: %d | Name: %-8s | Marks: %d" s.Id s.Name s.Marks
    printfn "----------------------------------------------------"

// 1. CREATE
studentTable <- studentTable @ [ { Id = 1; Name = "Alice"; Marks = 85 } ]
studentTable <- studentTable @ [ { Id = 2; Name = "Bob";   Marks = 90 } ]
showTable "CREATE (Inserted Alice & Bob)"

// 2. READ
printfn "\n[READ] Searching for Student with ID = 1:"
match studentTable |> List.tryFind (fun s -> s.Id = 1) with
| Some s -> printfn "Found: Name = %s, Marks = %d" s.Name s.Marks
| None   -> printfn "Student not found."

// 3. UPDATE
studentTable <- studentTable |> List.map (fun s -> if s.Id = 2 then { s with Marks = 95 } else s)
showTable "UPDATE (Updated Bob's Marks to 95)"

// 4. DELETE
studentTable <- studentTable |> List.filter (fun s -> s.Id <> 1)
showTable "DELETE (Deleted Alice with ID = 1)"


// ============================================================================
// PART 2: NoSQL Database (Document / MongoDB Model with Flexible Key-Values)
// ============================================================================
printfn "\n\n=========================================================="
printfn " PART 2: NoSQL DATABASE (MongoDB / Document Model)"
printfn "=========================================================="

// A document is a flexible map of Key -> Value (Schema-less)
type Document = Map<string, obj>

let mutable mongoCollection: Document list = []

let showCollection (step: string) =
    printfn "\n--- Collection State after: %s ---" step
    if mongoCollection.IsEmpty then
        printfn "(Collection is empty)"
    else
        for doc in mongoCollection do
            let fields = doc |> Map.toList |> List.map (fun (k, v) -> sprintf "%s: %O" k v) |> String.concat ", "
            printfn "{ %s }" fields
    printfn "----------------------------------------------------"

// 1. CREATE (Inserting documents with flexible fields)
let doc1 = Map [ ("_id", box 1); ("name", box "Alice"); ("marks", box 85) ]
let doc2 = Map [ ("_id", box 2); ("name", box "Bob");   ("marks", box 90); ("club", box "Robotics") ] // extra field!

mongoCollection <- mongoCollection @ [ doc1; doc2 ]
showCollection "CREATE (Inserted 2 Documents)"

// 2. READ (Find by _id)
printfn "\n[READ] Query document where _id = 2:"
match mongoCollection |> List.tryFind (fun d -> d.["_id"] = box 2) with
| Some d -> printfn "Found Document: %A" d
| None   -> printfn "Document not found."

// 3. UPDATE (Update marks field)
mongoCollection <- mongoCollection |> List.map (fun d ->
    if d.["_id"] = box 2 then d.Add("marks", box 95) else d
)
showCollection "UPDATE (Changed marks of _id 2 to 95)"

// 4. DELETE (Delete document with _id = 1)
mongoCollection <- mongoCollection |> List.filter (fun d -> d.["_id"] <> box 1)
showCollection "DELETE (Deleted document with _id = 1)"

