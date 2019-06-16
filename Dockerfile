FROM nginx:mainline
ADD dist /srv/www
ADD infra/nginx.conf /etc/nginx/conf.d/default.conf

