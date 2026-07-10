FROM nginx:alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY . /usr/share/nginx/html
RUN sed -i "s|var API_BASE = '[^']*'|var API_BASE = ''|" /usr/share/nginx/html/index.html
