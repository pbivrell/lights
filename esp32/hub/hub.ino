#include <WiFi.h>
#include <WebServer.h>
#include <Update.h>
#include <EEPROM.h>
#include <FS.h>


#define MAX_INSTRUCTIONS 20000
#define EEPROM_DATA_SIZE 32
#define EEPROM_TOTAL_SIZE (2 * EEPROM_DATA_SIZE) + 2
#define LED_PIN     15
#define LED_DEFAULT 1000

const String Version = "0.0.1";
const String DeviceKey = "theoriginallightskey";
const String DeviceID = "Hub1";
const char* ssidAP     = "LightHub";
const char* passwordAP = "defaultpassword";

IPAddress localIP(192, 168, 1, 115);
IPAddress gateway(192, 168, 1, 254);
IPAddress subnet(255, 255, 255, 0);

bool uploading = false, uploaded = false;

WebServer server(80);


// readCreds retrieves credentials from EEPROM, returning a boolean indicating if the credentials exist.
bool readCreds(char ssid[EEPROM_DATA_SIZE], char password[EEPROM_DATA_SIZE]) {

  char ok[2];
  ok[0] = char(EEPROM.read(0 + (EEPROM_DATA_SIZE * 2)));
  ok[1] = char(EEPROM.read(1 + (EEPROM_DATA_SIZE * 2)));
  
  if (ok[0] != 'O' || ok[1] != 'K') {
    return false;
  }
  
  for (int i = 0; i < EEPROM_DATA_SIZE; ++i){
    ssid[i] = char(EEPROM.read(i));
    if (ssid[i] == '\0') {
      break;
    }
  }

  for (int i = 0; i < EEPROM_DATA_SIZE; ++i){
    password[i] = char(EEPROM.read(i + EEPROM_DATA_SIZE));
    if (password[i] == '\0') {
      break;
    }
  }

  return true;
}

bool writeCreds(String ssidS, String passwordS) {

  char ssid[EEPROM_DATA_SIZE];
  char password[EEPROM_DATA_SIZE];

  if (ssidS.length() > EEPROM_DATA_SIZE || passwordS.length() > EEPROM_DATA_SIZE) {
    return false;
  }

  ssidS.toCharArray(ssid, 32);
  passwordS.toCharArray(password, 32);
  
  for (int i = 0; i < ssidS.length(); ++i){
    EEPROM.write(i, ssid[i]);
  }
  EEPROM.write(ssidS.length(), '\0');
  
  for (int i = 0; i < passwordS.length() ; ++i){
    EEPROM.write(i + EEPROM_DATA_SIZE, password[i]);
  }
  EEPROM.write(passwordS.length() + EEPROM_DATA_SIZE, '\0');


  char ok[2];
  EEPROM.write(0 + (EEPROM_DATA_SIZE * 2), 'O');
  EEPROM.write(1 + (EEPROM_DATA_SIZE * 2), 'K');

  EEPROM.commit();
  return true;
}

// startAP creates an AP with the provided ssid and password at the globally configured network.
// startAP then also configures and starts a webserver with routes for connecting to a local network.
void startAP(const char *ssid, const char *password) {

  WiFi.softAPConfig(localIP, gateway, subnet);
  WiFi.softAP(ssid, password);

  IPAddress IP = WiFi.softAPIP();
  Serial.print("AP IP address: ");
  Serial.println(IP);

  server.enableCORS();
  server.on("/login", handleLogin);
  server.on("/ip", handleIP);
  server.on("/networks", handleNetwork);

  server.begin();
}

void handleIP() {
    Serial.println("Handler: /ip");
    server.send(200, "text/html", WiFi.localIP().toString()); 
}

void handleNetwork() {
  Serial.println("Handler: /networks");

  int n = WiFi.scanNetworks();
   
  String networkJSON = "[";

  for (int i = 0; i < n; ++i) {      
    networkJSON += '"' + WiFi.SSID(i) + '"';
    if (i < n-1) {
      networkJSON += ",";
    }
    delay(10);
  }

  networkJSON += "]";

  server.send(200, "application/json", networkJSON);
  
}

void startServer() {

  server.stop();
  
  server.on("/status", handleStatus);
  server.on("/upload",HTTP_POST, []() {
    server.send(200, "text/plain", "");
  }, handleFileUpload);
  server.enableCORS();

  server.begin();

}

String statusJSON(boolean secret) {
    String key = secret ? "" : DeviceKey;
    String json =  "{" 
      "\"ip\":\""  + WiFi.localIP().toString() + "\"," +
      "\"mac\":\""  + WiFi.macAddress() + "\"," + 
      "\"key\":\""  + key + "\"," +
      "\"id\":\"" + DeviceID + "\"," +
      "\"version\":\"" + Version + "\"" +
    "}";
    return json;
}

void handleStatus() {
    Serial.println("Handler: /status");

    String json = statusJSON(true);
    
    server.send(200, "application/json", json);
}

// startWifi takes the provided ssid and password and attempts to join the network. It will return true once the status is connected or
// after timeout * 500ms it will return false
bool startWiFi(const char * ssid, const char * password, int timeout) {
  Serial.println("Connecting to WIFI");
  
  WiFi.begin(ssid, password);
  Serial.println("");

  int n = 0;
  // Wait for connection
  while (WiFi.status() != WL_CONNECTED && n++ < timeout) {
    delay(500);
    Serial.print(".");
  }
  Serial.println("");
  Serial.print("Connected to ");
  Serial.println(ssid);
  Serial.print("IP address: ");
  Serial.println(WiFi.localIP());
  return WiFi.status() == WL_CONNECTED;
}

// handleLogin is a response handler that reads the query paramters ssid and password then attempts to
// join that network. It will return 200 OK if connected and 500 Not Connected otherwise.
void handleLogin() {
  Serial.println("Handler: /login");

  String ssid = server.arg("ssid");
  String password = server.arg("password");

  if(startWiFi(ssid.c_str(), password.c_str(), 15)) {
    if (writeCreds(ssid, password)) {
      server.send(200, "text/html", "OK");
      ESP.restart();
      return;
    }
  }
  
  server.send(500, "text/html", "Not connected");
  
}

size_t instructions = 0;
size_t upload_size = 0;

static uint8_t dataBuffer[MAX_INSTRUCTIONS];

void handleFileUpload(){ 
  Serial.println("Handler: /upload");

  HTTPUpload& upload = server.upload();

  if(upload.status == UPLOAD_FILE_START){
    uploading = true;
    upload_size = 0;
    //Serial.println("Beginning upload");
  } else if(upload.status == UPLOAD_FILE_WRITE){
    //Serial.println("Copying upload");
    //Serial.println(upload.totalSize);
    //Serial.println(upload.currentSize);
    memcpy(&dataBuffer[upload_size], &upload.buf[0], upload.currentSize);
    upload_size += upload.currentSize;
  } else if(upload.status == UPLOAD_FILE_END) {
    if (upload_size != upload.totalSize) {
      Serial.println("entire file not uploaded");
    }else {
      //Serial.println("success");
      //Serial.println(upload_size);
      instructions = upload.totalSize / 8;
      uploaded = true;    
    }
  } else {
    Serial.println("upload failed");
  }
}

void setup() {
  
   char ssid[32];
   char password[32]; 

   Serial.begin(115200);

   EEPROM.begin(EEPROM_TOTAL_SIZE);
   
   delay(5000);
   
   bool hasCredentials = readCreds(ssid, password);
   if (hasCredentials) {
       if (startWiFi(ssid, password, 60)) {
           startServer();
       }
   } else {
     WiFi.mode(WIFI_MODE_APSTA);  
     startAP(ssidAP, passwordAP);
   }
}


void loop() {
  server.handleClient();
}
