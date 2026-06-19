# JOURNAL_IA.md — Usage de l'IA

## Utilisation générale

J'ai utilisé l'IA comme assistant pour accélérer l'écriture du code boilerplate et la structuration du projet, ainsi que la rédaction des fichiers.md, pour une mise au propre plus professionnel.
Je fournis les idées, l'IA corrige, modifie mais surtout explique le raisonemment afin de travailler en binôme.

## Ce que j'ai accepté

- **Structure du projet** : l'IA a proposé l'arborescence `cmd/` + `internal/` conforme aux conventions Go. J'ai accepté car c'est le standard de facto pour les projets Go.
- **Worker pool** : le pattern fan-out/fan-in avec channels bufferisés et `sync.WaitGroup` correspond exactement à ce qui est enseigné en cours. J'ai vérifié l'absence de fuite de goroutines et de data race.
- **Middlewares** : le pattern `func(http.Handler) http.Handler` est standard en Go. J'ai accepté l'implémentation après vérification.

## Ce que j'ai vérifié/modifié

- **Gestion du context** : j'ai vérifié que le timeout par URL utilise bien un context enfant du context de batch, et que `ctx.Done()` est bien écouté dans le mock checker.
- **Tests** : j'ai vérifié que `go test -race ./...` passe sans erreur, confirmant l'absence de data race.
- **Contrat JSON** : j'ai vérifié que les struct tags produisent exactement les noms de champs demandés dans le sujet (`batch_id`, `status_code`, `latency_ms`, etc.).
- **Validation** : j'ai ajusté les bornes de validation pour correspondre exactement au sujet (concurrency 1-50, timeout 100-30000).

## Ce que j'ai appris

- Le pattern d'une goroutine dédiée à la fermeture du channel de résultats (`go func() { wg.Wait(); close(results) }()`) est élégant et évite les deadlocks.
- L'utilisation de `errors.As` avec un type personnalisé `ValidationError` pour extraire le champ fautif est plus propre qu'un switch sur des chaînes de caractères.
