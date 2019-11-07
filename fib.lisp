; syntax: definition
; tokens: keywords, identifiers, numbers, parens
; types: number
(def max 10)
(print max)

; syntax: list
; tokens: quotes
; types: list
(def cases '(0 1 2 max))
(print cases)

(print (< 5 (+ 2 1)))
(print (= (+ 3 7) max))

(print
  (if (< 7 max) "zwei seelen wohnen ach" "in meiner brust"))

(if (= 1 2)
  (print "foo")
  (seq
    (print "bar")
    (print "baz")))

; syntax: lambda, call, if
; tokens: operators
; types: function
(def sum (fn (a b) (+ a b)))
(print (sum 5 2))

; closures
(def x 17)
(def testX
  (fn (x)
    (fn () x)))
(print ((testX 13)))

; syntax: string
; tokens: double-quotes
; types: string
(def greeting "Some fibonacci numbers:")
(print greeting)

; syntax: defun
(defun fib (n)
  (if (< n 2)
    1
    (+ (fib (- n 1)) (fib (- n 2)))))

(print (first '(1 2 3)))
(print (rest '(1 2 3)))

; syntax: seq
; types: null
(print "defining foreach")
(defun foreach (list f)
  (if (empty list)
    null
    (seq
      (f (first list))
      (foreach (rest list) f))))

(print "printing lines")
(foreach cases print)
