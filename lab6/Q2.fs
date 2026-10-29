printf "Enter a number: "
let n = int (System.Console.ReadLine())

let rec isPrime n d =
    if n < 2 then false
    elif d * d > n then true
    elif n % d = 0 then false
    else isPrime n (d + 1)

if isPrime n 2 then
    printfn "Prime Number"
else
    printfn "Not a Prime Number"
