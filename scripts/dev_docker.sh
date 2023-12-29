#!/bin/bash

# To use this locally you'll need to create a IAM service account key,
# then copy that file into your image and use that to auth as the service account 
# also need to set ADC to have it work properly
# https://stackoverflow.com/questions/45472882/how-to-authenticate-google-cloud-sdk-on-a-docker-ubuntu-image

# Create a local bridge network for both of our containers
docker network create -d bridge rkc-bot-network  

# Start the ngrok container on our newly created network and retrieve the publically accessible URL for the tunnel
docker run --rm -it -e NGROK_AUTHTOKEN=<YOUR_NGROK_AUTHTOKEN> --network rkc-bot-network --name ngrok ngrok/ngrok:alpine http <BOT_CONTAINER_TAG>:<WEBHOOK PORT>

# Set the appropriate env variables in Dockerfile and then build the Telegram bot image
# WEBHOOK_DOMAIN will need to change the public URL for the tunnel will change everytime you start the ngrok container
# Also need to change the `ListenAddr` on all `webhookOpts` to `<BOT_CONTAINER_TAG>/<WEBHOOK PORT>`
docker build -t aptl06/rkc-bot .

# Run the ngrok container on the same network as ngrok
docker run -it --name <BOT_CONTAINER_TAG> --network rkc-bot-network aptl06/rkc-bot
