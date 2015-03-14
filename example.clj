(import '(java.io File FileInputStream FileOutputStream))
(import '(org.fressian FressianReader FressianWriter))

(def f (File. "example.fressian"))
;(def r (FressianReader. (FileInputStream. f)))
(def w (FressianWriter. (FileOutputStream. f)))

;(.writeObject w 89)
;(.writeObject w 33554431)
(.writeObject w (float 1.2345))
;(.writeObject w 0.0)
;(.writeObject w 3.257329852835)

;(.writeObject w (java.util.Date.))

;(.writeObject w [1 2 3 4 5])
;(.writeObject w [1 2 3 4 5 "hello"])

;(.writeObject w {"hey" 3, "ho" 2, "answer" 42})

(.close w)
