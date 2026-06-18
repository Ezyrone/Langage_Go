# Langage Go - M2

Ce dépôt regroupe l'ensemble du cours de Go, les TP et les exercices réalisés dans le cadre du M2 AL à l'ESGI Grenoble

## TPs

| Dossier | Thème | Notions abordées |
|---------|-------|------------------|
| `TP/tp_1_exo_2a/` | Types, variables et constantes | Déclaration explicite, inférence (`:=`), constantes, `iota`, zero values, conversion de type |
| `TP/tp_2_exo_2b/` | Fonctions variadiques et retours multiples | `...type`, retours multiples, gestion d'erreurs avec `error`, filtrage de données |
| `TP/tp_3_exo_2c/` | Slices, Maps et performance | Slices (`append`, `len`, `cap`), Maps, structs, indexation par catégorie, benchmark avec/sans pré-allocation, `sort.Slice` |
| `TP/tp_4_exo_3a/` | Structs et méthodes | Structs, méthodes, receivers (valeur vs pointeur), `String()`, constructeurs avec validation, `math` |
| `TP/tp_5_exo_3b/` | Interfaces et interface vide | Interfaces, implémentation implicite, polymorphisme, `interface{}`, type switch, assertions de type |
| `TP/tp_6_exo_4a/` | Goroutines et synchronisation | Goroutines, `sync.WaitGroup`, canaux (`chan`), pool de travailleurs, concurrence |
| `TP/tp_7_exo_4c/` | Goroutines et synchronisation (DS) | Goroutines, `sync.WaitGroup`, canaux (`chan`), pool de travailleurs, concurrence |
| `TP/tp_8_exo_4d/` | Gestion d'événements avec `select` | `select`, channels multiples, `time.NewTicker`, goroutines productrices, arrêt propre |
| `TP/tp_9_exo_4e/` | Mutex et WaitGroup | `sync.Mutex`, `sync.WaitGroup`, race conditions, `sync/atomic`, sections critiques |
| `TP/tp_10_exo_4f/` | Context, annulation et timeout | `context.WithTimeout`, `context.WithCancel`, `ctx.Done()`, `ctx.Err()`, arrêt propre |
| `TP/tp_11_exo_5a/` | API REST avec `net/http` | `http.ServeMux`, handlers, CRUD en mémoire, `encoding/json`, `github.com/google/uuid` |
| `TP/tp_12_exo_5b/` | API REST avancée avec Gin | `gin-gonic/gin`, middlewares (logger, auth), `binding:"required"`, groupes de routes, versioning API |
| `TP/tp_13_exo_5c/` | Sérialisation JSON et struct tags | `encoding/json`, struct tags, `omitempty`, `json:"-"`, `Marshal`/`Unmarshal`, custom `Marshaler` |
| `TP/tp_14_exo_5d/` | Accès aux bases de données | `database/sql`, `sqlx`, GORM, SQLite, CRUD, ORM, auto-migration |

## Réponses aux questions du TP 4A (Goroutines et Synchronisation)

**Exercice 1 — Que constatez-vous dans la sortie ?**
Sans mécanisme de synchronisation, `main()` se termine immédiatement après avoir lancé les goroutines. Le programme s'arrête avant que les goroutines n'aient le temps de finir leur travail. On observe peu ou pas de messages "Tâche terminée", car quand la fonction `main` retourne, toutes les goroutines encore en cours sont tuées.

**Exercice 2 — Le comportement a-t-il changé ?**
Oui. Grâce à `sync.WaitGroup`, `main()` attend que chaque goroutine appelle `wg.Done()` avant de continuer. Toutes les goroutines terminent désormais leur travail, et le message "Toutes les goroutines ont terminé leur exécution." s'affiche uniquement après la fin de toutes les tâches.

**Exercice 3 — L'ordre des résultats correspond-il à l'ordre des IDs ?**
Non. Les goroutines s'exécutent en concurrence avec des durées aléatoires, donc celle qui finit en premier envoie son résultat en premier dans le canal. L'ordre de lecture des résultats reflète l'ordre de terminaison des goroutines, pas l'ordre de lancement.

**Exercice 4 — Comment le nombre de travailleurs affecte-t-il le temps total ?**
Les tâches sont réparties entre les 3 travailleurs. L'ordre de traitement dépend de la disponibilité de chaque travailleur. Avec plus de travailleurs, le temps total diminue car davantage de tâches sont traitées en parallèle. Avec 3 travailleurs pour 10 tâches, chacun traite environ 3 à 4 tâches.

## Réponses aux questions du TP 5C (Sérialisation JSON et Struct Tags)

**Exercice 1 — Les clés JSON correspondent-elles aux noms des champs ?**
Oui, exactement. Sans struct tags, Go utilise le nom du champ tel quel comme clé JSON (`Nom`, `Age`, `Email`, `Actif` avec majuscule). C'est rarement le format souhaité pour une API (on préfère `snake_case` ou `camelCase`).

**Exercice 2 — Effet de `omitempty` et de `json:"-"` ?**
Le tag `omitempty` fait que `contact_email` est absent du JSON quand `Email == ""`. Le tag `json:"-"` empêche `MotDePasse` d'apparaître dans le JSON, quelle que soit sa valeur, ce qui est essentiel pour ne pas exposer des données sensibles.

**Exercice 3 — Clé inconnue et type incorrect ?**
Une clé JSON sans champ correspondant dans la struct (ex: `description`) est simplement ignorée par `json.Unmarshal`. Si `unit_price` était une string `"79.99"` au lieu d'un nombre, `Unmarshal` retournerait une erreur de type (`cannot unmarshal string into Go struct field ... of type float64`).

**Exercice 4 — Pourquoi toujours vérifier les erreurs ?**
Un JSON malformé ou avec des types incorrects peut provoquer des données corrompues ou des zero values silencieuses, menant à des bugs difficiles à diagnostiquer. Vérifier systématiquement les erreurs de `Marshal`/`Unmarshal` garantit la fiabilité du traitement des données.

**Exercice 5 — Questions de réflexion**

*Question 1 (objet imbriqué) :* Pour un champ `publisher_info` contenant un objet, on crée une struct `Editeur` séparée et on l'utilise comme champ dans `Livre` avec le tag `json:"publisher_info,omitempty"`. Utiliser un pointeur (`*Editeur`) permet de distinguer "absent" (`nil`) de "présent mais vide".

*Question 2 (timestamp Unix) :* Pour sérialiser un `time.Time` en timestamp Unix, on crée un type personnalisé (`UnixTime`) qui implémente les interfaces `json.Marshaler` et `json.Unmarshaler`. `MarshalJSON` retourne `t.Unix()` et `UnmarshalJSON` reconstruit le `Time` depuis le timestamp entier.

## Réponses aux questions du TP 5D (Accès aux Bases de Données)

### Comparaison des trois approches

| Approche | Avantages | Inconvénients | Quand l'utiliser |
|---|---|---|---|
| **database/sql** | Bibliothèque standard, aucune dépendance extra. Contrôle total sur le SQL. Léger et performant. Interface stable et bien documentée. | Mapping manuel via `rows.Scan()` — verbeux et fragile. Code de lecture répétitif (`Query` + `rows.Next()` + `Scan`). Pas de vérification du mapping à la compilation. | Projets simples, microservices avec peu de requêtes, contrôle maximal sur le SQL. |
| **sqlx** | Extension directe de `database/sql` (100% compatible). `db.Select()` et `db.Get()` mappent automatiquement vers des structs via les tags `db:"..."`. Code concis et lisible. On écrit toujours du SQL pur. | Dépendance externe (légère et stable). Requêtes SQL toujours écrites à la main. Pas de génération de schéma ni de migrations. | La majorité des projets Go en production — meilleur compromis contrôle/ergonomie. |
| **GORM** | ORM complet, pas de SQL pour les opérations courantes. `AutoMigrate` crée les tables depuis les structs. API chainable intuitive. Relations, hooks, transactions, soft delete intégrés. | Abstraction lourde qui masque le SQL exécuté. Moins performant sur les requêtes complexes. Courbe d'apprentissage. Peut générer du SQL inefficace. | Prototypage rapide, apps CRUD classiques, vitesse de développement prioritaire. |

### Différences clés observées dans le code

| Aspect | database/sql | sqlx | GORM |
|---|---|---|---|
| Connexion | `sql.Open()` | `sqlx.Open()` | `gorm.Open()` |
| Création table | SQL manuel | SQL manuel | `AutoMigrate()` |
| Insert | `db.Exec()` | `db.Exec()` | `db.Create()` |
| Select all | `Query()` + boucle `Scan()` | `db.Select()` | `db.Find()` |
| Select one | `QueryRow().Scan()` | `db.Get()` | `db.First()` |
| Lignes de code pour GetAll | ~12 lignes | ~3 lignes | ~3 lignes |

## Autres

- `training_exo_1/` — Brouillons et tests perso (hello world, module greetings, etc.)
