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

; syntax: seq
; types: null
(defun foreach (list f)
  (if (= (first list) null)
    null
    (seq
      (f (first list))
      (foreach (rest list)))))

(def lines (foreach cases print))

; all together
(foreach lines (fn (line) (print '(line "\n"))))
