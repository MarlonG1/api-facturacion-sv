# API Facturación El Salvador

[![Go Version](https://img.shields.io/badge/Go-1.23-blue.svg)](https://golang.org/)

Una API para la gestión, emisión y transmisión de Documentos Tributarios Electrónicos (DTE) que cumple con los requisitos establecidos por la autoridad fiscal.

## 📋 Características

- Emisión de facturas electrónicas
- Emisión de comprobantes de crédito fiscal (CCF)
- Invalidación de documentos
- Manejo de contingencias
- Transmisión de documentos a Hacienda
- Monitoreo de métricas
- Autenticación JWT
- Firmado electrónico de documentos

## 🏗️ Arquitectura

Este proyecto está implementado siguiendo principios de:

- **Arquitectura Hexagonal (Ports & Adapters)**: Separación clara entre lógica de dominio y acceso a recursos externos
- **Domain Driven Design (DDD)**: Modelado de dominio basado en el negocio

### Capas de la arquitectura:

- **Dominio**: Modelos, reglas de negocio y puertos (interfaces)
- **Aplicación**: Casos de uso y orquestación
- **Infraestructura**: Adaptadores para bases de datos, API, comunicación externa, etc.
- **Bootstrap**: Configuración e inicialización de dependencias

## 🛠️ Tecnologías

- **Go 1.23**: Lenguaje de programación principal
- **Base de datos**: Actualmente soporta MySQL y PostgreSQL
- **Redis**: Caché y almacenamiento de tokens
- **Docker y Docker Compose**: Contenerización y orquestación de servicios
- **Gorilla Mux**: Router HTTP
- **GORM**: ORM para acceso a base de datos
- **JWT**: Autenticación basada en tokens

## 🔧 Requisitos previos

- Docker y Docker Compose
- Go 1.23+ (sólo para desarrollo)
- Certificados de firma digital (para ambiente de producción o pruebas)

## 📦 Instalación

### Con Docker (Recomendado)

1. Clonar el repositorio:
```bash
git clone https://github.com/MarlonG1/api-facturacion-sv.git
cd dte-microservice
```

2. Colocar certificados de firma digital en la carpeta `scripts/temp` (crear si no existe)

3. Iniciar los servicios:
```bash
docker-compose up -d
```

### Para desarrollo

```bash
docker-compose --profile dev up -d
```

### Configuración

Las variables de entorno están predefinidas en el archivo `docker-compose.yml`. Modifícalo según tus necesidades.

## 🚀 Uso

### API Endpoints

El servicio expone las siguientes APIs:

#### Endpoints

- `POST /api/v1/auth/login`: Autenticación de usuarios
- `POST /api/v1/auth/register`: Registro de nuevos clientes
- `POST /api/v1/dte/invoices`: Crear factura electrónica
- `POST /api/v1/dte/ccf`: Crear comprobante de crédito fiscal
- `POST /api/v1/dte/invalidation`: Invalidar documento
- `GET /api/v1/dte`: Listar todos los documentos emitidos por el usuario
- `GET /api/v1/test`: Prueba los componentes del sistema
- `GET /api/v1/metrics`: Obtener métricas de los endpoints
- `GET /api/v1/health`: Estado de salud del servicio

## 🚧 Gestión de contingencias

El sistema maneja automáticamente contingencias cuando:

1. Hay fallas de conexión con el sistema de Hacienda
2. Hay problemas de conectividad a internet
3. Hay fallas en el firmado digital de documentos
4. Sistema de Hacienda no está disponible

Los documentos se almacenan y retransmiten según las reglas configuradas.

## 🔐 Seguridad

- Autenticación basada en tokens JWT
- Validación estricta de entradas
- Firmado digital de documentos

## 🔄 Integración Continua (CI)

Este proyecto utiliza un pipeline de integración continua con dos ramas específicas para la generación de builds:

- **release-amd64**: Compilación y despliegue de la versión para arquitectura `amd64`
- **release-arm64**: Compilación y despliegue de la versión para arquitectura `arm64`

Cada rama se encarga de generar imágenes optimizadas para su respectiva arquitectura, asegurando compatibilidad en distintos entornos de ejecución.

## 👥 Contribución

Para contribuir a este proyecto:

1. Analizar documentación antes de sugerir implementaciones
2. Respetar la arquitectura establecida
3. Mantener consistencia con implementaciones existentes
4. Validar contra JSON Schema
5. No asumir comportamientos no documentados
6. Justificar cualquier complejidad adicional

## 📚 Documentación adicional

- [Guía de referencias JSON Schema y catálogos oficiales](https://factura.gob.sv/informacion-tecnica-y-funcional/)