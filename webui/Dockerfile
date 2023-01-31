FROM node:16 as builder

WORKDIR /app

ENV PATH /app/node_modules/.bin:$PATH

COPY package.json ./package.json
COPY yarn.lock ./yarn.lock
COPY .npmrc ./.npmrc

RUN yarn install --quiet

COPY . .

RUN yarn build

FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html

COPY ./config/nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]