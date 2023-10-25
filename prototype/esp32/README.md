# API documentation

The following page documents the various API's that make up the light project. There are various API's and consumers continated inside of the light project, each documented below by it's serving platform with the intended consumer noted.

* ESP-Light
* ESP-Hub
* Backend API

## ESP Hub APIs

The Hub is the interconnection point between the dumb household light devices and the webserver. The Hub's responability is to track the lights connected to it via it's WiFi network. Then it must control it's collected lights using 2 mechanisms. Polling a webserver over the internet for changes to this it's lights states. Or push changes via direct HTTP API calls from clients on the same network as the device.

### Startup

The start up process executes the following steps.

1. Read EEPROM memory for WiFi credentials. if credentials goto 2 else goto 3
2. Attempt to connect to WiFi network for 30s
3. Start access point
4. Start DNS server

### State

The following is documentation about the internal state that the hub devices holds.

#### Lights List

The Light List maintains a set of data about the lists which have been registered with it.

Foreach Light it holds the following information:

```
  uint8_t Count;
  String IP;
  String ID;
  String Pattern;
  bool Status;
  String Updated;
```

This state is not presisted on shutdown. It can be reconstructed from a combination of a light re-registering and an update from an external source (push or pull). This state is meant to represent data about what lights are "online" and what state they are currently in.

#### Hub Last Update

The hub maintants status about who and when it was last updated. These can be either "cloud" or an identifier for the controlling device on the same network as the hub. When a change comes in from a device. The rest of updators will have invalid current state of the lights. 


### Soft AP API (AKA webserver)
 The hub upon start will create an Access Point. This access point can be indentified with the SSID **LightHub** using the intentionally static password **defaultpassword**

> This password is intentionally static because we want the Light devices to require 0 configuration. See design decision "Why a hub". The hub network is hidden but can be connected to by those who are instructed how (or savvy enough).

The softAP API enables device network configuration and simiplistic tools for debugging and light controlling when a network connection is unavialable.

####/register

Client: Light Device

Purpose: To track a light device as it connects to the hub

State: This handler updates the [Light List State](#Lights-List) to include the new light.

### Polling "API"

### DNS Server

The hub device has a DNS server for the sole purpose of functioniting as a Captive Portal. The TLDR on capative portal is that when devices connect to a new network the device will attempt a request on the newly acquired network. To do this the device firsst has to make a DNS request (port 53). The hub being the provider of this network intercepts this DNS request. This request is rerouted to the IP of the hub's webserver, the webserver then has handlers to enable configuring the WiFi. 


## Design Decision

### Why a hub

### Polling and Pushing
