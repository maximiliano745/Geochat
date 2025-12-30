# Base stage: Usa una imagen ligera de Node.js para construir el proyecto.
FROM node:20-alpine AS build

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos de manifiesto del proyecto
COPY package*.json ./

# Instala las dependencias del proyecto
RUN npm install

# Copia el resto de los archivos de la aplicación
COPY . .

# Construye la aplicación Vue/Vite
RUN npm run build

# Final stage: Usa una imagen de NGINX para servir los archivos estáticos
FROM nginx:alpine

# Copia los archivos estáticos de la aplicación construida al directorio de NGINX
COPY --from=build /app/dist /usr/share/nginx/html

# Puerto por defecto para NGINX
EXPOSE 80

# Comando para iniciar NGINX
CMD ["nginx", "-g", "daemon off;"]
