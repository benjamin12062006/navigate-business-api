# Reconstrucciones de cursos Go

Fecha de incorporación al proyecto: 2026-09-27

Esta carpeta conserva dentro del proyecto las 16 reconstrucciones guiadas que antes estaban en `/private/tmp/codex-evidence-guides`. Fueron movidas, no regeneradas. La comparación SHA-256 confirmó antes del traslado que cada archivo temporal era idéntico a su guía correspondiente en la biblioteca local.

También se copiaron las 16 transcripciones originales correspondientes junto a sus guías. Esto conserva los enlaces relativos de `Fuente primaria` y `Contexto previo`, y permite auditar cada reconstrucción sin salir de la documentación del proyecto. Las copias no reemplazan la biblioteca maestra.

## Alcance real

Las guías actuales proceden únicamente de estos dos cursos principales:

1. `APIs de Producción con Go`.
2. `Building, scaling and deploying distributed systems in Golang and Kubernetes`.

Cubren arquitectura por capas, repositorios, transacciones SQL, autenticación, autorización y una introducción a la estructura de sistemas distribuidos. Una reconstrucción es material derivado para estudiar el código; la transcripción original continúa siendo la fuente primaria.

## Lo que todavía no está aquí

Los cursos de Trevor Sawler, Maximilian Schwarzmüller y otros autores fueron localizados como fuentes complementarias, pero todavía no tienen reconstrucciones guardadas en esta carpeta. No debe interpretarse el inventario de cursos como si todos ya hubieran sido leídos, reconstruidos o aplicados.

Boris Paskhaver fue propuesto como referencia pedagógica, pero no aparece como autor de un curso dentro de la biblioteca local de Go. Por esa razón no se le atribuyen decisiones técnicas de este backend.

## Regla para añadir nuevas guías

Solo se incorporará una reconstrucción cuando una decisión o incremento del proyecto necesite comprender código, contratos u orden de cambios de una lección concreta. Cada guía nueva debe conservar su curso y sección, enlazar la transcripción original y distinguir lo enseñado por el instructor de la adaptación realizada en `navigate-business-api`.
