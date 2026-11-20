(def coll1 [1 2 3 4 5 6 7 8 9 10])
(def coll2 '(5 10 15 20 25 30))

(println "--- Collection 1 (Vector) ---")
(println "Original:" coll1)
(println "1. map (squares):" (map #(* % %) coll1))
(println "2. filter (evens):" (filter even? coll1))
(println "3. reduce (sum):" (reduce + coll1))
(println "3. reduce (max):" (reduce max coll1))
(println "4. take (first 3):" (take 3 coll1))
(println "4. drop (first 3):" (drop 3 coll1))

(def chained1
  (->> coll1
       (filter even?)
       (map #(* % %))
       (take 3)))
(println "Chained (filter evens -> square -> take 3):" chained1)

(println "\n--- Collection 2 (List) ---")
(println "Original:" coll2)
(println "1. map (squares):" (map #(* % %) coll2))
(println "2. filter (evens):" (filter even? coll2))
(println "3. reduce (sum):" (reduce + coll2))
(println "3. reduce (max):" (reduce max coll2))
(println "4. take (first 2):" (take 2 coll2))
(println "4. drop (first 2):" (drop 2 coll2))

(def chained2
  (->> coll2
       (filter even?)
       (map #(* % 2))
       (reduce +)))
(println "Chained (filter evens -> double -> sum):" chained2)

