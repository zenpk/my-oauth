# MyOAuth

Self-hosted OAuth2.0 implementation, with PKCE support.

Note: This is the side-project of other side-projects.

## Configuration

### Back End

Edit `conf-prod.json` file to configure the backend-related settings.

### Front End

Edit the `BASE` constant in `frontend/src/apis/basic.ts` to point to the actual backend service endpoint.

## Build with Docker

```shell
sudo docker build -t myoauth .
sudo docker run -dp 20476:80 myoauth
```

## Screenshots

Coming soon.

## API

Due to being too lazy to write docs, please refer to the `tests` folder.

## Usage

Recommend using `sdks`.
