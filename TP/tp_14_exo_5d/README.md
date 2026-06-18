# TP 14 - Accès aux Bases de Données en Go

## Comparaison des trois approches

### database/sql (package standard)

**Avantages :**
- Fait partie de la bibliothèque standard Go, aucune dépendance supplémentaire (hormis le driver).
- Contrôle total sur les requêtes SQL écrites.
- Léger et performant, pas de couche d'abstraction.
- Interface stable et bien documentée.

**Inconvénients :**
- Mapping manuel des colonnes vers les champs de struct via `rows.Scan()` — verbeux et source d'erreurs si l'ordre des colonnes change.
- Le code de lecture (`Query` + boucle `rows.Next()` + `Scan`) est répétitif.
- Pas de vérification à la compilation du mapping colonnes/struct.

**Quand l'utiliser :** Projets simples, microservices avec peu de requêtes, ou quand on veut un contrôle maximal sur le SQL sans dépendance externe.

### sqlx

**Avantages :**
- Extension directe de `database/sql` (compatible à 100%).
- `db.Select()` et `db.Get()` mappent automatiquement les résultats vers des structs grâce aux tags `db:"..."`.
- Élimine le code de scan répétitif — le code est plus concis et lisible.
- Pas de magie : on écrit toujours du SQL pur, donc on garde le contrôle.

**Inconvénients :**
- Dépendance externe (bien que très légère et stable).
- Toujours besoin d'écrire les requêtes SQL à la main.
- Pas de génération automatique de schéma ou de migrations.

**Quand l'utiliser :** La majorité des projets Go en production. C'est le meilleur compromis entre contrôle SQL et ergonomie du code.

### GORM

**Avantages :**
- ORM complet : pas besoin d'écrire de SQL pour les opérations courantes.
- `AutoMigrate` crée et met à jour les tables automatiquement à partir des structs.
- API chainable intuitive (`db.Where(...).Find(...)`, etc.).
- Gestion intégrée des relations, hooks, transactions, soft delete.

**Inconvénients :**
- Abstraction lourde qui peut masquer les requêtes SQL réellement exécutées.
- Moins performant pour les requêtes complexes (jointures multiples, sous-requêtes).
- Courbe d'apprentissage pour les fonctionnalités avancées.
- Peut générer des requêtes SQL inefficaces si mal utilisé.

**Quand l'utiliser :** Prototypage rapide, applications CRUD classiques, projets où la vitesse de développement prime sur l'optimisation fine des requêtes.

## Différences clés observées dans le code

| Aspect | database/sql | sqlx | GORM |
|---|---|---|---|
| Connexion | `sql.Open()` | `sqlx.Open()` | `gorm.Open()` |
| Création table | SQL manuel | SQL manuel | `AutoMigrate()` |
| Insert | `db.Exec()` | `db.Exec()` | `db.Create()` |
| Select all | `Query()` + boucle `Scan()` | `db.Select()` | `db.Find()` |
| Select one | `QueryRow().Scan()` | `db.Get()` | `db.First()` |
| Lignes de code pour GetAll | ~12 lignes | ~3 lignes | ~3 lignes |

## Exécution

```bash
go run main.go
```

Pour repartir d'une base vierge, supprimer les fichiers `.db` avant de relancer :

```bash
rm -f test.db test_sqlx.db gorm_test.db
go run main.go
```
