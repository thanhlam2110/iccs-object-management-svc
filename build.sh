#!/bin/bash
# version 2020.05.06
# /bin/bash build.sh <project> <branch> [Dockerfile]
project=""				# project-name
tag=""			 		# docker image tag (branch)
Dockerfile="Dockerfile" # Dockerfile

# Check parameter project and branch
if [ ! -z "$2" ]; then
	project="$1"
	if [ $2 == "master" ]; then
	tag="latest"
	elif [ $2 == "pro" ] || [ $2 == "prod" ]  || [ $2 == "release" ]; then
			tag="stable"
	else
		tag=$2
	fi
else
	echo "Can not build empty tag for project $2"
	exit 1
fi
# Check parameter Dockerfile
if [ ! -z "$3" ]; then
	Dockerfile="$3"
fi
# Check Dockerfile exits
if [ -f "./Dockerfile" ]; then
	echo "$(date +'%d-%m-%Y %H:%M:%S') [$HOSTNAME] Build $Dockerfile"
	DOCKER_BUILDKIT=1 docker build -f $Dockerfile -t $project:$tag . 
else
	echo "$Dockerfile not found!"
	exit 1
fi
retVal=$?
if [ $retVal -gt 0 ]; then
    echo "Build failed - exit code: $retVal"
    exit 1
fi
# Push to local registry
if [ ! -z "$localRegistry" ]; then
    localRegistryServer=$(echo $localRegistry | cut -d'|' -f1)
    localRegistryUser=$(echo $localRegistry | cut -d'|' -f2)
    localRegistryPassword=$(echo $localRegistry | cut -d'|' -f3)
	if [ ! -z "$localRegistryServer" ]; then
		echo "Push $project:$tag to local Docker Registry: $localRegistryServer"
		docker login -u $localRegistryUser -p $localRegistryPassword $localRegistryServer
		docker tag $project:$tag $localRegistryServer/$project:$tag
		docker push $localRegistryServer/$project:$tag
	else
		echo "Local Registry not set. Skip push to local Registry!"
	fi
else
	echo "Local Registry not set. Skip push to local Registry!"
fi
# Push to remote registry (dockerhub or 3rd provider)
if [ ! -z "$remoteRegistry" ]; then
    remoteRegistryServer=$(echo $remoteRegistry | cut -d'|' -f1)
    remoteRegistryUser=$(echo $remoteRegistry | cut -d'|' -f2)
    remoteRegistryPassword=$(echo $remoteRegistry | cut -d'|' -f3)
	if [ -z "$remoteRegistryServer" ]; then # push to dockerhub
		echo "Push $remoteRegistryUser/$project:$tag to DockerHub"
		docker login -u $remoteRegistryUser -p $remoteRegistryPassword
		docker tag $project:$tag $remoteRegistryUser/$project:$tag
		docker push $remoteRegistryUser/$project:$tag
	else # push to 3rd provider
		echo "Push $remoteRegistryServer/$project:$tag to $remoteRegistryServer"
		docker login -u $remoteRegistryUser -p $remoteRegistryPassword $remoteRegistryServer
		docker tag $project:$tag $remoteRegistryServer/$project:$tag
		docker push $remoteRegistryServer/$project:$tag
	fi
else
	echo "Remote Registry not set. Skip push to remote Registry!"
fi

echo "$(date +'%d-%m-%Y %H:%M:%S') [$HOSTNAME] End build"
