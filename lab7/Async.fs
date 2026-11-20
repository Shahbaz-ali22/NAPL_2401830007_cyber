module Lab7.AsyncDemo

open System

// Task 1
let task1 =
    async {
        printfn "Task 1 Started"
        do! Async.Sleep(3000)
        printfn "Task 1 Completed"
        return "Result 1: Data Downloaded"
    }

// Task 2
let task2 =
    async {
        printfn "Task 2 Started"
        do! Async.Sleep(2000)
        printfn "Task 2 Completed"
        return "Result 2: File Processed"
    }

// Task 3
let task3 =
    async {
        printfn "Task 3 Started"
        do! Async.Sleep(1000)
        printfn "Task 3 Completed"
        return "Result 3: Database Checked"
    }

let results =
    [| task1; task2; task3 |]
    |> Async.Parallel
    |> Async.RunSynchronously


printfn "\nAll Tasks Results:"
results |> Array.iter (printfn "%s")


