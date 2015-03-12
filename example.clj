(import '(java.io File FileInputStream FileOutputStream))
(import '(org.fressian FressianReader FressianWriter))

(def f (File. "example.fressian"))
;(def r (FressianReader. (FileInputStream. f)))
(def w (FressianWriter. (FileOutputStream. f)))

;(.writeObject w [1 2 3 4 5])
(.writeObject w [1 2 3 4 5 "hello"])

(.close w)
