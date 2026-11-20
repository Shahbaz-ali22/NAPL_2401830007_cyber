open System
open MongoDB.Bson
open MongoDB.Driver

// Connect to MongoDB
let client = MongoClient("mongodb://localhost:27017")
let db = client.GetDatabase("College")
let col = db.GetCollection<BsonDocument>("students")

// Clear collection for fresh run
col.DeleteMany(BsonDocument()) |> ignore

// Helper function to print collection
let showStudents () =
    printfn "--- MongoDB Students ---"
    let docs = col.Find(BsonDocument()).ToList()
    for d in docs do
        printfn "Id: %O | Name: %O | Marks: %O" d.["_id"] d.["name"] d.["marks"]

// 1. CREATE
printfn "\n1. CREATE"
let s1 = BsonDocument([ ("_id", BsonInt32(1)); ("name", BsonString("Alice")); ("marks", BsonInt32(85)) ])
let s2 = BsonDocument([ ("_id", BsonInt32(2)); ("name", BsonString("Bob")); ("marks", BsonInt32(90)) ])
col.InsertOne(s1)
col.InsertOne(s2)
showStudents ()

// 2. READ
printfn "\n2. READ (Student 1)"
let filter1 = Builders<BsonDocument>.Filter.Eq("_id", 1)
let found = col.Find(filter1).FirstOrDefault()
if found <> null then
    printfn "Found Student: %O, Marks: %O" found.["name"] found.["marks"]

// 3. UPDATE
printfn "\n3. UPDATE (Bob Marks to 95)"
let filter2 = Builders<BsonDocument>.Filter.Eq("_id", 2)
let update = Builders<BsonDocument>.Update.Set("marks", 95)
col.UpdateOne(filter2, update) |> ignore
showStudents ()

// 4. DELETE
printfn "\n4. DELETE (Delete Alice)"
col.DeleteOne(filter1) |> ignore
showStudents ()
