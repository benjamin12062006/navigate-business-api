# 04 — Contrato neutral para verificar la identidad

Fecha de creación: 2026-10-01
Estado: Implementado y verificado
Ámbito: `navigate-business-api`, frontera entre HTTP y proveedores de identidad
Relacionada con: [módulo auth y migraciones](03-modulo-auth-y-migraciones-estandar.md), [composición modular](../config/01-composicion-modular-y-dependencias-compartidas.md), [arquitectura basada en Tiago](../../docs/fonts/go-backend/arquitectura/01-tiago-arquitectura-y-repositorios.md)
Revisión de partida: `2cfc3ec5f9eb061967618a202bdc5d9499538e5e`; árbol de trabajo con cambios previos conservados

## Necesidad y punto de partida

Corregir un acoplamiento puntual sin cambiar las rutas ni el comportamiento observable de autenticación. El middleware HTTP aceptaba una interfaz `TokenVerifier`, pero importaba el adaptador Firebase para reconocer `firebaseauth.ErrInvalidToken`. La interfaz también vivía en el paquete HTTP. Por ello, otro verificador podía satisfacer el método `Verify` y aun así necesitar devolver un error perteneciente a Firebase para obtener la respuesta `401` correcta.

El servicio y el repositorio ya estaban separados de Firebase. Este incremento se limita al contrato transversal de verificación y no modifica Firebase Admin, la persistencia, Expo ni los endpoints.

## Evidencia e interpretación

La lección 25 distingue transporte, reglas y almacenamiento, y propone inyectar abstracciones en lugar de hacer que una capa dependa de una implementación intercambiable.[^F1] La reconstrucción guiada recupera ese mecanismo mediante contratos para almacenamiento y autenticación que se inyectan en la composición.[^R1]

La fuente no prescribe dónde ubicar `TokenVerifier`, cómo nombrar errores ni una estructura exacta de carpetas. La adaptación para este proyecto consiste en hacer que `internal/auth` sea dueño del puerto y de las categorías de error independientes del proveedor. Firebase traduce errores de su SDK al contrato neutral; HTTP traduce ese contrato a códigos de estado.

Relación aplicada:

```text
necesidad: sustituir o probar el verificador sin conocer Firebase
→ fuente: depender de abstracciones e inyectar implementaciones
→ reconstrucción: contrato Authenticator intercambiable
→ interpretación: el puerto y sus errores no pertenecen a un adaptador concreto
→ decisión: mover contrato y errores a internal/auth
```

## Alternativas y decisión

| Enfoque | Ventaja | Límite | Decisión |
| --- | --- | --- | --- |
| Conservar el error en Firebase | No requiere cambios y funciona con el proveedor actual. | HTTP y sus pruebas deben importar Firebase; otro verificador puede convertir un token inválido en `503`. | Descartado porque contradice el propósito del puerto existente. |
| Crear un paquete `ports` adicional | Hace explícita la categoría arquitectónica. | Añade una carpeta para un contrato pequeño y separa conceptos que ya pertenecen a auth. | Pospuesto; no aporta valor en este incremento. |
| Alojar puerto y errores en `internal/auth` | Mantiene un contrato neutral con el mínimo de estructura y respeta la dirección hacia el núcleo. | El núcleo de la funcionalidad incorpora el concepto de token de identidad. | Elegido; es la corrección mínima completa. |

## Implementación y funcionamiento

El nuevo [contrato neutral](../../internal/auth/verifier.go) declara `TokenVerifier`, `ErrInvalidIdentityToken` y `ErrIdentityProviderUnavailable` dentro de `internal/auth`. El contrato expresa lo que la autenticación necesita sin mencionar un SDK concreto.

El [adaptador Firebase](../../internal/auth/firebase/verifier.go) conserva la llamada a Firebase Admin y la regla `google.com`, pero ahora traduce sus resultados al contrato neutral. Una aserción de compilación confirma que `*firebaseauth.Verifier` implementa `auth.TokenVerifier`.

El [middleware HTTP](../../internal/auth/http/middleware.go) depende de `auth.TokenVerifier` y clasifica únicamente errores de `internal/auth`: token inválido produce `401`; proveedor indisponible o un fallo inesperado conservan `503`. Ya no importa el paquete Firebase. El [módulo](../../internal/auth/module/module.go) también recibe el puerto neutral. Las [pruebas HTTP](../../internal/auth/http/handler_test.go) dejaron de fabricar errores de Firebase y comprueban las dos categorías comunes; la [prueba del adaptador](../../internal/auth/firebase/verifier_test.go) verifica su traducción al error neutral.

| Decisión respaldada | Archivo y símbolo | Adaptación realizada |
| --- | --- | --- |
| El consumidor depende de una abstracción intercambiable.[^F1][^R1] | [`auth.TokenVerifier`](../../internal/auth/verifier.go) | El puerto salió del adaptador HTTP y pasó al núcleo de la funcionalidad. |
| Los adaptadores traducen hacia el contrato de la aplicación. | [`firebaseauth.Verifier.Verify`](../../internal/auth/firebase/verifier.go) | Los errores del SDK se convierten en categorías independientes de Firebase. |
| HTTP interpreta el contrato, no el proveedor. | [`Authenticator.RequireGoogle`](../../internal/auth/http/middleware.go) | La respuesta conserva `401`/`503` sin importar `internal/auth/firebase`. |

## Comprobación

| Comando o prueba ejecutada | Resultado observado |
| --- | --- |
| `GOCACHE=/private/tmp/expo-go-build-cache go test ./...` | Pasó en todos los paquetes; las pruebas HTTP conservaron `401`, `503`, sincronización y autenticador compartido. |
| `GOCACHE=/private/tmp/expo-go-build-cache go vet ./...` | Pasó sin diagnósticos. |
| Inspección de importaciones de `./internal/auth/http` con `go list` | No apareció `internal/auth/firebase`; HTTP depende del contrato `internal/auth`. |
| `git diff --check` | Pasó sin errores de espacios. |

La primera invocación de Go no pudo leer el caché global del host por las restricciones del entorno. Se repitieron las comprobaciones con un caché temporal permitido y finalizaron correctamente. No se llamó a Firebase real ni se modificó su configuración; las pruebas del contrato usan dobles y la prueba del adaptador usa tokens decodificados construidos en memoria.

## Consecuencias y continuidad

El contrato HTTP permanece igual: token inválido devuelve `401` y proveedor no disponible devuelve `503`. La corrección no añade otro proveedor; evita que el proveedor actual forme parte del contrato entre capas y permite que un futuro adaptador satisfaga el mismo puerto sin devolver errores definidos por Firebase.

No se cambió la bitácora `03`: conserva la organización histórica que originó el módulo. Esta entrada registra la corrección posterior y se relaciona con ella sin reescribir ese contexto.

## Fuentes

[^F1]: [25. Clean Layered Architecture](../../docs/reconstrucciones-cursos-go/APIs%20de%20Producción%20con%20Go/Sección%205%20Scaffolding%20our%20API%20Server/25.%20Clean%20Layered%20Architecture.md), Tiago Taquelim, *The complete guide for building and deploying production web APIs in Golang*, sección 5. Localizador: líneas 17–31; buscar `Transportschicht`, `Schnittstelle` y `HTTP-Handler`. Consulta: 2026-10-01.
[^R1]: [25. Clean Layered Architecture — Reconstrucción guiada](../../docs/reconstrucciones-cursos-go/APIs%20de%20Producción%20con%20Go/Sección%205%20Scaffolding%20our%20API%20Server/25.%20Clean%20Layered%20Architecture%20-%20Reconstrucción%20guiada.md), fases 1–3 y contrato `Authenticator`; evidencia derivada con certeza alta para responsabilidades y dirección de dependencias. Consulta: 2026-10-01.
