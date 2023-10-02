#!/bin/bash

# Nom de l'image Docker que vous souhaitez créer
IMAGE_NAME="forum_image"

# Nom du conteneur Docker
CONTAINER_NAME="forum_container"

# Répertoire de travail où se trouvent les fichiers de votre forum
WORK_DIR="."

# Port depuis l'hôte vers le conteneur
HOST_PORT=10000
CONTAINER_PORT=80

# Étape 1 : Créez l'image Docker
docker build -t $IMAGE_NAME $WORK_DIR

# Étape 2 : Créez et lancez un conteneur à partir de l'image
docker run -d --name $CONTAINER_NAME -p $HOST_PORT:$CONTAINER_PORT $IMAGE_NAME

# Étape 3 : Liste les fichiers et dossiers à l'intérieur du conteneur
docker exec $CONTAINER_NAME ls -l 

# Étape 4 : Affichez des informations sur le conteneur nouvellement créé
docker ps -a | grep $CONTAINER_NAME


