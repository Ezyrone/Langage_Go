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

## Réponses aux questions du TP 4A (Goroutines et Synchronisation)

**Exercice 1 — Que constatez-vous dans la sortie ?**
Sans mécanisme de synchronisation, `main()` se termine immédiatement après avoir lancé les goroutines. Le programme s'arrête avant que les goroutines n'aient le temps de finir leur travail. On observe peu ou pas de messages "Tâche terminée", car quand la fonction `main` retourne, toutes les goroutines encore en cours sont tuées.

**Exercice 2 — Le comportement a-t-il changé ?**
Oui. Grâce à `sync.WaitGroup`, `main()` attend que chaque goroutine appelle `wg.Done()` avant de continuer. Toutes les goroutines terminent désormais leur travail, et le message "Toutes les goroutines ont terminé leur exécution." s'affiche uniquement après la fin de toutes les tâches.

**Exercice 3 — L'ordre des résultats correspond-il à l'ordre des IDs ?**
Non. Les goroutines s'exécutent en concurrence avec des durées aléatoires, donc celle qui finit en premier envoie son résultat en premier dans le canal. L'ordre de lecture des résultats reflète l'ordre de terminaison des goroutines, pas l'ordre de lancement.

**Exercice 4 — Comment le nombre de travailleurs affecte-t-il le temps total ?**
Les tâches sont réparties entre les 3 travailleurs. L'ordre de traitement dépend de la disponibilité de chaque travailleur. Avec plus de travailleurs, le temps total diminue car davantage de tâches sont traitées en parallèle. Avec 3 travailleurs pour 10 tâches, chacun traite environ 3 à 4 tâches.

## Autres

- `training_exo_1/` — Brouillons et tests perso (hello world, module greetings, etc.)
