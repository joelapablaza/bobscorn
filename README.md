# Corn Test Project

Este proyecto contiene dos carpetas principales:

- `Front/`: Aplicación frontend hecha con Next.js
- `Back/`: Backend desarrollado en Go y dockerizado

---

## 🚀 Levantar el backend

No es necesario tener Go instalado localmente. El backend corre dentro de un contenedor Docker usando una imagen de Go.

Para levantarlo, simplemente ejecutá:

```bash
docker compose up --build
```

Esto construirá la imagen (si es necesario) y levantará el servicio del backend automáticamente.

---

## 🖥️ Levantar el frontend

1. Ir a la carpeta del frontend:

```bash
cd Front
```

2. Instalar dependencias usando Bun o PNPM:

```bash
bun install
# o
pnpm install
```

3. Copiar o renombrar el archivo de entorno:

```bash
cp .env.example .env.local
```

4. Ejecutar la aplicación en modo desarrollo:

```bash
bun run dev
# o
pnpm run dev
```

---

Con estos pasos vas a tener el frontend corriendo en modo desarrollo y el backend levantado dentro de Docker, sin necesidad de instalar Go en tu máquina local. ✅
