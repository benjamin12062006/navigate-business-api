Air quedó integrado en `navigate-business-api`.

Cambios realizados:

- Air `v1.67.4` registrado en [go.mod](/Users/benji/Documents/EXPO-FASTAPI/navigate-business-api/go.mod) mediante la directiva `tool`.
- Configuración del watcher en [.air.toml](/Users/benji/Documents/EXPO-FASTAPI/navigate-business-api/.air.toml).
- Carpeta temporal `tmp/` agregada a [.gitignore](/Users/benji/Documents/EXPO-FASTAPI/navigate-business-api/.gitignore).
- Instrucciones incorporadas al [README.md](/Users/benji/Documents/EXPO-FASTAPI/navigate-business-api/README.md).

Para iniciar el backend con recarga automática:

```bash
cd /Users/benji/Documents/EXPO-FASTAPI/navigate-business-api
go tool air
```

Comprobé el ciclo completo:

1. Air compiló `./cmd/api`.
2. Inició la API.
3. Detectó una modificación en `router.go`.
4. Recompiló y reinició el servidor.
5. La ruta respondió `200 OK`.
6. Al cerrar Air, eliminó `tmp/`.
7. `go test ./...` pasó correctamente.

Air está registrado en el proyecto, así que no necesitas instalarlo globalmente. Este es el mecanismo recomendado actualmente por su [documentación oficial](https://github.com/air-verse/air).

Los cambios están funcionando y permanecen sin commit ni tag; `v1.1.0` continúa representando el incremento anterior.
