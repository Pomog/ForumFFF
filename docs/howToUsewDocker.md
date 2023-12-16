## Docker Usage: 

### FORGOT WSL SUDO PASSWORD? !!!
If you've forgotten your WSL SUDO password, use Windows CMD to reset it:

```bash
wsl --user root
passwd <username>
exit
```
### Docker Desktop uses WSL 2 as the default backend for running Linux containers. Ensure that WSL 2 is enabled on your Windows machine.

```bash
sudo apt-get update
sudo apt-get install -y docker.io
docker --version
```
got: Docker version 24.0.5, build 24.0.5-0ubuntu1~22.04.1

restart
```bash
docker info
```
Client:
Version:    24.0.5
Context:    default
Debug Mode: false

Server:
ERROR: Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?
errors pretty printing info

### INSTALL DOCKER WINDOWS

1. download and install https://docs.docker.com/desktop/install/windows-install/
2. sing up
3. Start Docker Desktop

use WSL terminal
```bash
docker info
```
### now it works

### to build Container
```bash
docker build -t ffforum .
```
  install requierments

### to run Container
```bash
docker run --name=ffforum -p 80:8080 ffforum
```

## type in browser for Testing
http://localhost/docker
or
localhost:80

### to stop Container
```bash
docker stop ffforum
```

### to prune Container
```bash
docker builder prune -a
```
--- Clean Build Context !!!

### comments 
sudo systemctl start docker 	
	-- When you run this command, it starts the Docker service,
	which allows you to use Docker commands to create and manage containers.

sudo systemctl enable docker
	--Running this command ensures that Docker starts automatically whenever you restart your system. 