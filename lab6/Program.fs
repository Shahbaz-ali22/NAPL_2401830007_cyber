open System
open System.Text.RegularExpressions

// 1. Reverse string and check palindrome
printf "Enter a string: "
let str = Console.ReadLine()
let rev = new string (Array.rev (str.ToCharArray()))
printfn "Reversed: %s" rev

if str = rev then
    printfn "Palindrome"
else
    printfn "Not a Palindrome"

// 2. Transform string
printf "\nEnter a string to transform: "
let text = Console.ReadLine()
let isNum, num = Int32.TryParse(text)

let transformed =
    match isNum with
    | true -> string (num * 2)
    | false when text |> Seq.forall Char.IsLetter -> text.ToUpper()
    | false -> text

printfn "Transformed: %s" transformed

// 3. Email validation
printf "\nEnter an email: "
let email = Console.ReadLine()

match Regex.IsMatch(email, @"^[^@\s]+@[^@\s]+\.[^@\s]+$") with
| true  -> printfn "Valid Email"
| false -> printfn "Invalid Email"
