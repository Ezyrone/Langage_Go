# DESIGN.md — Justification architecturale

## Découpage en packages

Le projet suit le principe d'inversion de dépendance :

- **`domain`** : types métier (`CheckResult`, `Batch`, `Summary`), interfaces (`Checker`, `Store`) et erreurs. Aucune dépendance extérieure. Tous les autres packages dépendent de `domain`, jamais l'inverse.
- **`checker`** : implémentation concrète du `Checker`. `HTTPChecker` pour la production, `MockChecker` déterministe pour les tests. L'interface dans `domain` permet de substituer l'un par l'autre sans modifier le code appelant.
- **`pool`** : moteur concurrent. Reçoit un `Checker` (interface) et une liste d'URLs, retourne un `Batch`. Découplé du transport HTTP.
- **`store`** : persistance. `MemoryStore` implémente `Store` avec une `map` protégée par `sync.RWMutex`. Le `Get` wrappe `ErrBatchNotFound` via `fmt.Errorf("%w")` pour permettre `errors.Is` dans la couche API.
- **`api`** : couche HTTP. Handlers, routage, validation, DTOs, middlewares. Traduit les erreurs métier en codes HTTP via `errors.Is` / `errors.As`.
- **`cmd/urlwatch`** : assemblage uniquement. Crée les dépendances, les câble et lance le serveur.

Ce découpage garantit que chaque package a une responsabilité unique et est testable isolément.

## Modèle de concurrence

### Worker pool borné

Le nombre de goroutines est strictement borné par `opts.Concurrency`. Les URLs sont distribuées via un channel bufferisé `jobs` (taille = nombre d'URLs) et les résultats collectés via un channel bufferisé `results`.

**Pourquoi des channels bufferisés ?**
- `jobs` bufferisé à `len(urls)` : permet d'envoyer toutes les URLs sans bloquer l'émetteur, puis de fermer le channel. Les workers consomment à leur rythme.
- `results` bufferisé à `len(urls)` : permet aux workers d'écrire leur résultat sans attendre que le collecteur lise, évitant tout deadlock.

### Fan-out / fan-in

- **Fan-out** : les N workers lisent depuis le même channel `jobs`. La distribution est naturelle et équitable.
- **Fan-in** : tous les workers écrivent dans le même channel `results`. Une goroutine séparée attend via `sync.WaitGroup` puis ferme `results`, signalant la fin de la collecte.

### Context et timeouts

- **Timeout global (batch)** : calculé dynamiquement comme `perURLTimeout × (⌈nURLs/concurrency⌉ + 1)`. Ce calcul laisse assez de temps pour traiter toutes les vagues d'URLs, plus une marge. Un timeout identique au per-URL causerait des expirations prématurées avec beaucoup d'URLs et peu de concurrence.
- **Timeout per-URL** : chaque URL reçoit un `context.WithTimeout` enfant du context batch, borné par `timeout_ms` configuré par le client.
- Le `MockChecker` écoute `ctx.Done()` pour simuler une annulation, ce qui est testé dans `TestRun_RespectsContext_Cancellation`.

### Fuites de goroutines

Risques identifiés et mitigations :
1. **Workers qui ne terminent pas** : le channel `jobs` est fermé après l'envoi, donc le `range jobs` termine toujours. Le context enfant empêche un `Check` de bloquer indéfiniment.
2. **Channel `results` jamais fermé** : une goroutine dédiée appelle `wg.Wait()` puis `close(results)`, garantissant que `range results` se termine.
3. **Goroutine de fermeture bloquée** : impossible car chaque worker fait `defer wg.Done()`, et la boucle `range jobs` termine grâce au `close(jobs)`.

### Échecs partiels

Les échecs d'URLs individuelles n'interrompent pas le lot. Chaque `CheckResult` porte son propre champ `Error`. Le résumé agrégé (`Summary.Down`) compte les échecs. Le client obtient toujours un résultat complet.

## Gestion des erreurs

- **Sentinelle** : `ErrBatchNotFound` dans `domain` — utilisé avec `errors.Is` dans le handler pour retourner 404.
- **Type personnalisé** : `ValidationError` avec le champ fautif — utilisé avec `errors.As` dans le handler pour retourner 400.
- **Wrapping** : `store.Get` wrappe l'erreur sentinelle avec `fmt.Errorf("batch %q: %w", id, ErrBatchNotFound)`, enrichissant le message tout en préservant la chaîne `errors.Is`.

Traduction vers HTTP :
- `errors.Is(err, ErrBatchNotFound)` → 404
- `errors.As(err, &ValidationError{})` → 400
- Tout le reste → 500

## Choix stdlib (net/http) vs Gin

J'ai choisi `net/http` (stdlib) pour plusieurs raisons :
- Aucune dépendance externe requise (pas de `go.sum` complexe).
- Le routage est simple (3 endpoints) et ne justifie pas un framework.
- Les middlewares sont implémentés comme des `func(http.Handler) http.Handler`, pattern standard et composable.
- L'examen demande de démontrer la maîtrise du langage, pas d'un framework.

## Pourquoi Go est un bon choix ici (3 arguments + 1 limite)

1. **Goroutines et channels natifs** : le worker pool se code en quelques lignes sans bibliothèque de concurrence externe. En Java, il faudrait un `ExecutorService` + `Future` + `CountDownLatch`. En Python, `asyncio` ou `ThreadPoolExecutor` avec le GIL.

2. **Compilation statique et déploiement** : un seul binaire sans dépendance runtime. Idéal pour un microservice conteneurisé (image Docker < 15 Mo). En comparaison, une JVM ou un interpréteur Python alourdit l'image.

3. **`context.Context` intégré à la stdlib** : la propagation de timeout et d'annulation est native dans `net/http`, `database/sql`, etc. Pas besoin d'un framework pour gérer les timeouts proprement.

**Limite ressentie** : l'absence de génériques matures complique certains patterns (ex. écrire une fonction utilitaire de map/filter sur des slices typées demande de la verbosité ou des dépendances sur `slices`/`maps`). Le typage des erreurs est aussi moins ergonomique qu'un `Result<T, E>` en Rust.
