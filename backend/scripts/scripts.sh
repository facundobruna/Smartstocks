#!/bin/bash

# Script de configuración inicial para Smart Stocks Backend

echo "🚀 Configurando Smart Stocks Backend..."

# Colores
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Verificar Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go no está instalado. Por favor instala Go 1.21+${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Go encontrado: $(go version)${NC}"

# Verificar Docker
if ! command -v docker &> /dev/null; then
    echo -e "${YELLOW}⚠️  Docker no encontrado. Se recomienda instalar Docker para desarrollo${NC}"
else
    echo -e "${GREEN}✅ Docker encontrado${NC}"
fi

# Crear archivo .env si no existe
if [ ! -f .env ]; then
    echo -e "${YELLOW}📝 Creando archivo .env...${NC}"
    cp .env.example .env
    echo -e "${GREEN}✅ Archivo .env creado. Por favor configura tus variables.${NC}"
else
    echo -e "${GREEN}✅ Archivo .env ya existe${NC}"
fi

# Instalar dependencias
echo -e "${YELLOW}📦 Instalando dependencias de Go...${NC}"
go mod download
go mod tidy
echo -e "${GREEN}✅ Dependencias instaladas${NC}"

# Crear directorios necesarios
echo -e "${YELLOW}📁 Creando directorios...${NC}"
mkdir -p bin
mkdir -p logs
mkdir -p tmp
echo -e "${GREEN}✅ Directorios creados${NC}"

# Preguntar si quiere usar Docker
echo -e "\n${YELLOW}¿Deseas levantar los servicios con Docker Compose? (y/n)${NC}"
read -r response

if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
    if command -v docker-compose &> /dev/null || command -v docker &> /dev/null; then
        echo -e "${YELLOW}🐳 Levantando servicios con Docker...${NC}"
        docker-compose up -d
        echo -e "${GREEN}✅ Servicios levantados${NC}"
        echo -e "${GREEN}📊 La API estará disponible en http://localhost:8080${NC}"
        echo -e "${YELLOW}Ver logs: docker-compose logs -f${NC}"
    else
        echo -e "${RED}❌ Docker Compose no está disponible${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  Recuerda configurar MySQL y Redis manualmente${NC}"
    echo -e "${YELLOW}⚠️  Ejecuta las migraciones: make migrate-up${NC}"
    echo -e "${YELLOW}⚠️  Luego ejecuta: make run${NC}"
fi

echo -e "\n${GREEN}🎉 Configuración completada!${NC}"
echo -e "\n${YELLOW}Comandos útiles:${NC}"
echo -e "  make help          - Ver todos los comandos"
echo -e "  make run           - Ejecutar servidor"
echo -e "  make docker-up     - Levantar con Docker"
echo -e "  make docker-logs   - Ver logs"
echo -e "  make test          - Ejecutar tests"