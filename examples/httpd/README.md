# Apache HTTP Server (httpd) Setup on AWS EC2 Instance for POC

This guide provides step-by-step instructions to set up an Apache HTTP Server (`httpd`) on an AWS EC2 instance for a POC. The goal is to have a fully functional web server that starts automatically on instance reboot and serves a JSON response with the current time and server IP.

## Prerequisites

1. An AWS account.
2. An EC2 instance running Ubuntu 22.04.
3. SSH access to the EC2 instance.

### Step 1 - Install the HTTP server 
```
sudo apt update
sudo apt install apache2 -y
sudo apt install libapache2-mod-wsgi-py3 -y
```

### Step 2 - Enable Apahce to Start on Boot
```sudo systemctl enable apache2``` 

### Step 3 - Copy scripts and configuration 
```
scp -i /path/to/your-key-pair.pem time.py ubuntu@your-ec2-public-dns:/home/ubuntu/
scp -i /path/to/your-key-pair.pem 000-default.conf ubuntu@your-ec2-public-dns:/home/ubuntu/
scp -i /path/to/your-key-pair.pem ../scripts/awscliv2.zip ubuntu@your-ec2-public-dns:/home/ubuntu/example.zip

ssh -i /path/to/your-key-pair.pem ubuntu@your-ec2-public-dns

sudo mv /home/ubuntu/time.py /var/www/html/time.py  
sudo mv /home/ubuntu/000-default.conf /etc/apache2/sites-available/000-default.conf
sudo mv /home/ubuntu/example.zip /var/www/html/example.zip
```

restart server:
```sudo systemctl restart apache2```  

Try locally: 
```
curl http://localhost/time

{"current_time": "2024-05-18T09:41:11.156297", "server_ip": "10.0.14.215"}

```


To download zipped file (65MB~): 
```
wget http://localhost/example.zip
```