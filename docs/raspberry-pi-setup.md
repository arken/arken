# How to Set Up a Raspberry Pi to Run Arken

(This guide is going to mainly focus on installing Arken on a Raspberry Pi running Raspbian/RaspberryPiOS)

## Install Docker

1. Install docker's prerequisites

```bash
sudo apt-get install apt-transport-https ca-certificates software-properties-common -y
```

2. Download & Install Docker

```bash
curl -fsSL get.docker.com -o get-docker.sh && sh get-docker.sh
```

3. Set up the Apt Repository for Docker

Add the following line to your apt configuration.

###### /etc/apt/sources.list

```bash
deb https://download.docker.com/linux/raspbian/ stretch stable
```

4. Update the RPi

```bash
sudo apt update && sudo apt upgrade
```

5. Start & Enable the Docker Daemon

```bash
systemctl start --now docker.service
```

## Install Docker-Compose

1. Download and install pip3

```bash
sudo apt install python3-pip
```

2. Allow docker-compose to be executed as a program

```bash
sudo pip3 install docker-compose
```

## Install USBMount & Prepare External Drive

1. Make a directory to use as a mount point. For this tutorial I'm going to use `/mnt/data` but you can replace that with any other path you'd like.

```bash
sudo mkdir -p /mnt/data
```

**(If you're using the internal SD card of the Raspberry Pi you can move on to installing Arken from here.)**

2. Make the mount point immutable unless a drive is mounted.

```bash
sudo chattr +i /mnt/data
```

3. Install USBMount

```bash
sudo apt install usbmount
```

4. Configure USBMount by changing the configuration file to match the following.

###### /etc/usbmount/usbmount.conf

```roboconf
# Configuration file for the usbmount package, which mounts removable
# storage devices when they are plugged in and unmounts them when they
# are removed.

# Change to zero to disable usbmount
ENABLED=1

# Mountpoints: These directories are eligible as mointpoints for
# removable storage devices.  A newly plugged in device is mounted on
# the first directory in this list that exists and on which nothing is
# mounted yet.
MOUNTPOINTS="/mnt/data"

# Filesystem types: removable storage devices are only mounted if they
# contain a filesystem type which is in this list.
FILESYSTEMS="vfat ext2 ext3 ext4 hfsplus"
```

5. Disable UDev Private Mount to fix Auto Mount Not Working

```bash
sudo nano /lib/systemd/system/systemd-udevd.service
```

Change `PrivateMount=yes`  --> `PrivateMount=no`



6. Reboot the system. When logging back in with a drive attached you should see it as mounted to the expected location by typing.

```bash
lsblk
```

 Expected Output:

```
NAME        MAJ:MIN RM  SIZE RO TYPE MOUNTPOINT
sda           8:0    0  1.8T  0 disk 
`-sda1        8:1    0  1.8T  0 part /mnt/data
mmcblk0     179:0    0 29.8G  0 disk 
|-mmcblk0p1 179:1    0  256M  0 part /boot
`-mmcblk0p2 179:2    0 29.6G  0 part /
```

## ---

## Install Arken

1. Download the Arken Docker Compose Configuration

```bash
wget https://raw.githubusercontent.com/arkenproject/arken/master/docs/examples/docker-compose.yml
```

2. Open the `docker-compose.yml` you just downloaded and edit the configuration values to reflect your system. 
   
   ```bash
   nano docker-compose.yml
   ```
   
   - Replace: `</PATH/TO/MOUNTED/DRIVE/HERE>` with your mount point. In this tutorial we used `/mnt/data`.
   
   - Replace `<YOUR-DRIVE-SIZE>` with the size of your external drive, or the size of your SD card minus 5GB for the operating system.
     
     - Use the format `1TB` or  `1GB`
   
   - (Optional) Add the following line to the "environment" section under Arken. This email will only be used to send you an alert if your node doesn't check in for 24 hours or if your hard drive fails. After that alert is sent all your data will be scrubbed from our systems.
     
     ```yaml
          - ARKEN_STATS_EMAIL="you@example.com"
     ```

3. Start the Arken Application!

```bash
sudo docker-compose up -d --remove-orphans
```

4. Check the stats of the application

```bash
sudo docker-compose logs -f arken
```

(CTRL+C to escape.)

That's it! You should now have Arken running on your Raspberry pi!
