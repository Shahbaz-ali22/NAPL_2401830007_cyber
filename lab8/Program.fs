open System

// 1. Base Class with a virtual method
type Shape() =
    abstract member Area: unit -> float
    default this.Area() = 0.0

// 2. Derived Class 1: Circle
type Circle(radius: float) =
    inherit Shape()

    override this.Area() =
        3.14 * radius * radius

// Derived Class 2: Rectangle
type Rectangle(length: float, width: float) =
    inherit Shape()

    override this.Area() =
        length * width

// 3. Creating objects
let c = Circle(5.0)
let r = Rectangle(4.0, 6.0)

// Call methods directly
printfn "Circle Area: %.2f" (c.Area())
printfn "Rectangle Area: %.2f" (r.Area())

// Demonstrating Polymorphism using a list of base class (Shape)
let shapes: Shape list = [ c :> Shape; r :> Shape ]

printfn "\n Polymorphism Output "
for s in shapes do
    printfn "Area: %.2f" (s.Area())
