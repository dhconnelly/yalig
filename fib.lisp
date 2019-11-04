; syntax: definition
; tokens: keywords, identifiers, numbers, parens
; types: number
(def max 10)
max

; syntax: list
; tokens: quotes
; types: list
(def cases '(0 1 2 max))
cases

; syntax: string
; tokens: double-quotes
; types: string
(def greeting "Some fibonacci numbers:\n")
(print greeting)

; syntax: lambda, call, if
; tokens: operators
; types: function
(def fib
  (fn (n)
    (if (< n 2)
      1
      (+ (fib (- n 1)) (fib (- n 2))))))

; syntax: seq
; types: null
(def foreach
  (fn (list f)
    (if (= (first list) null)
      null
      (seq
        (f (first list))
        (printall (rest list))))))

; all together
(foreach cases
  (fn (num)
    (seq
      (print '("fib(" num ") = " (fib num) "\n")))))
