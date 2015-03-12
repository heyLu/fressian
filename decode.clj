(import '(org.fressian FressianReader))

(def r (FressianReader. System/in))

(defn indent [prefix]
  (str "  " prefix))

(defn ppf [prefix f o]
  (println (str prefix (pr-str (f o)))))

(defn pp [prefix o]
  (cond
    (nil? o) (ppf prefix identity o)

    (or (instance? java.util.List o)
        (.isArray (.getClass o)))
    (do
      (ppf prefix type o)
      (doseq [v o]
        (pp (indent prefix) v)))

    (instance? java.util.Map o)
    (do
      (ppf prefix type o)
      (doseq [[k v] o]
        (pp (indent prefix) k)
        (pp (indent (indent prefix)) v)
        (println)))

    (instance? org.fressian.TaggedObject o)
    (do
      (ppf prefix type o)
      (ppf (str (indent prefix) :tag " ") #(.getTag %) o)
      (ppf (indent prefix) identity :value)
      (pp (indent (indent prefix)) (.getValue o))
      (ppf (str (indent prefix) :meta " ") #(.getMeta %) o))

    :else (ppf prefix identity o)))

(defn read-all []
  (try
    (pp "" (.readObject r))
    (println)
    (catch java.io.EOFException e)))

(read-all)
