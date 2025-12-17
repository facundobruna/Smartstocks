# 🚀 Guía de Inicio Rápido - Smart Stocks Backend

Esta guía te ayudará a tener el backend funcionando en menos de 5 minutos.

## ✅ Prerrequisitos

- Docker y Docker Compose instalados
- O bien: Go 1.21+, MySQL 8.0 y Redis 7

## 🎯 Opción 1: Docker (Más Rápido)

### 1. Clonar el repositorio
```bash
git clone <tu-repo>
cd smartstocks-backend
```

### 2. Crear archivo .env
```bash
cp .env.example .env
```

### 3. Levantar servicios
```bash
docker-compose up -d
```

### 4. Verificar que funciona
```bash
curl http://localhost:8080/health
```

Deberías ver:
```json
{
  "status": "healthy",
  "service": "smart-stocks-api"
}
```

**¡Listo!** La API está corriendo en `http://localhost:8080`

---

## 🛠️ Opción 2: Desarrollo Local

### 1. Instalar dependencias
```bash
go mod download
```

### 2. Configurar MySQL
```bash
mysql -u root -p
CREATE DATABASE smartstocks;
```

### 3. Ejecutar migraciones
```bash
mysql -u root -p smartstocks < database/migrations/001_initial_schema.sql
```

### 4. Configurar .env
```bash
cp .env.example .env
# Edita los valores de conexión a MySQL y Redis
```

### 5. Iniciar Redis
```bash
redis-server
```

### 6. Ejecutar servidor
```bash
go run cmd/api/main.go
```

---

## 🧪 Probar la API

### 1. Registrar usuario
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "Password123"
  }'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Password123"
  }'
```

Guarda el `access_token` de la respuesta.

### 3. Obtener perfil (protegido)
```bash
curl -X GET http://localhost:8080/api/v1/user/profile \
  -H "Authorization: Bearer TU_ACCESS_TOKEN"
```

### 4. Listar colegios
```bash
curl http://localhost:8080/api/v1/schools
```

---

## 📦 Importar Colección de Postman

1. Abre Postman
2. Importa el archivo `postman/SmartStocks.postman_collection.json`
3. Las variables se configurarán automáticamente después del login

---

## 🔍 Verificar Servicios

### Verificar MySQL
```bash
# Con Docker
docker exec -it smartstocks-mysql mysql -u smartstocks -psmartst stocks123 -e "SHOW DATABASES;"

# Local
mysql -u root -p -e "SHOW DATABASES;"
```

### Verificar Redis
```bash
# Con Docker
docker exec -it smartstocks-redis redis-cli ping

# Local
redis-cli ping
```

### Ver logs
```bash
# Con Docker
docker-compose logs -f api

# Local
# Los logs aparecerán en la consola donde ejecutaste el servidor
```

---

## 🐛 Solución de Problemas

### Error: "connection refused" en MySQL
```bash
# Verifica que MySQL esté corriendo
docker-compose ps
# O localmente:
sudo systemctl status mysql
```

### Error: "connection refused" en Redis
```bash
# Verifica que Redis esté corriendo
docker-compose ps
# O localmente:
redis-cli ping
```

### Error: "bind: address already in use"
El puerto 8080 ya está en uso. Cambia el puerto en `.env`:
```
PORT=8081
```

### Limpiar y reiniciar todo (Docker)
```bash
docker-compose down -v
docker-compose up -d
```

---

## 📚 Siguientes Pasos

1. Lee el [README.md](README.md) completo
2. Revisa los endpoints en la colección de Postman
3. Explora el código en `internal/`
4. Prepárate para la Fase 2: Sistema de Quizzes

---

## 🎓 Estructura de la Base de Datos

La migración inicial crea:

- ✅ **users** - Usuarios del sistema
- ✅ **user_stats** - Estadísticas y puntos
- ✅ **schools** - Colegios asociados
- ✅ **refresh_tokens** - Tokens de sesión

Datos de prueba incluidos:
- 5 colegios de Argentina

---

## 💡 Tips

- Usa `make help` para ver todos los comandos disponibles
- El token JWT expira en 24 horas
- El refresh token expira en 30 días
- Rate limit: 100 requests por minuto por IP

---

## 📧 ¿Necesitas ayuda?

- Revisa los logs: `docker-compose logs -f`
- Verifica el health check: `curl http://localhost:8080/health`
- Contacto: smartstocksarg@gmail.com