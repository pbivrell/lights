#include <WiFi.h>
#include <WebServer.h>
#include <Update.h>
#include <EEPROM.h>
#include <FS.h>
#include <DNSServer.h>
#include <HTTPClient.h>


#define MAX_INSTRUCTIONS 20000
#define EEPROM_DATA_SIZE 32
#define EEPROM_TOTAL_SIZE (2 * EEPROM_DATA_SIZE) + 2
#define LED_PIN     15
#define LED_DEFAULT 1000
#define MAX_LIGHTS 10

const String Version = "0.0.1";
const String DeviceKey = "theoriginallightskey";
const String DeviceID = "Hub1";
const char* ssidAP     = "LightHub";
const char* passwordAP = "defaultpassword";

IPAddress apIP(8, 8, 4, 4);
IPAddress subnet(255, 255, 255, 0);

DNSServer dnsServer;

struct Light{
  uint8_t Count;
  String IP;
  String ID;
  String Pattern;
  bool Status;
};


struct Light Lights[MAX_LIGHTS];
int LightCount = 0;

bool WifiConnected = false;

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

  WiFi.softAPConfig(apIP, apIP, subnet);
  WiFi.softAP(ssid, password);

  IPAddress IP = WiFi.softAPIP();
  Serial.print("AP IP address: ");
  Serial.println(IP);

  server.enableCORS();
  server.on("/login", handleLogin);
  server.on("/ip", handleIP);
  server.on("/networks", handleNetwork);
  server.on("/register", handleRegister);
  server.on("/", handleIndex);
  server.on("/lights", handleLights);
  server.on("/cloudaddr", handleCloudAddr);
  server.onNotFound ( handleNotFound );

  server.begin();
}

String cloudAddr = "https://lights.paulbivrell.com";

void handleCloudAddr() {
  cloudAddr = server.arg("addr");
  server.send(200, "text/html", "OK");
}

bool reportCloud() {

  if (cloudAddr == "") {
    return false;
  }
  
  HTTPClient http;   
  http.begin(cloudAddr + "/hub"); 
  http.addHeader("Content-Type", "application/json");
  int httpResponseCode = http.POST(cloudReportJSON()); 
  String response = http.getString();  
  http.end();

  int index = 0;
  uint8_t count = 0;
  String item = "";
  String id = "";
  bool lightStatus = false;
  String pattern = "";
  
  for (auto c : response) {

    if (c == '\n'){
      pattern = item;
      item = "";
      count = 0;
      index = 0;
      LightChanged(lightStatus, id, pattern, count);
    }else if (c == ',') {
      if (index == 0) {
        id = item;
        ++index;
      }else if (c == 1){
        lightStatus = item == "true" ? true : false;
        ++index;
      }else {
        count = item.toInt();
        index = 0;
      }
      item = "";
    }else {
      item += c;
    }
  }
  
  return true;
}

void LightChanged(bool lightStatus, String ID, String pattern, uint8_t count) {
  //Serial.println("light changed");
  //Serial.println(lightStatus);
  //Serial.println(ID);
 //Serial.println(pattern);

  for (int i = 0; i < MAX_LIGHTS; i++) {
    if ( ID == Lights[i].ID ) {
      if (Lights[i].Status != lightStatus) {
        ToggleLight(Lights[i].IP);
        Lights[i].Status = lightStatus;
      }

      if (Lights[i].Pattern != pattern){
        UpdatePattern(Lights[i].IP, pattern);
        Lights[i].Pattern = pattern;
      }

      if (Lights[i].Count != count) {
        UpdateCount(Lights[i].IP, count);
        Lights[i].Count = count;
      }
      
    }
  }
}

void UpdatePattern(String IP, String pattern) {
  if (pattern[0] == '~') {
    ParsePattern(IP, pattern);
  }
}

void ParsePattern(String IP, String pattern) {
  Serial.println(pattern);
  int r,g,b, bright = 0;
  String item = "";
  int count = 0;
  for (auto c : pattern) {
    Serial.println("Parsing");
    Serial.println(c);
    if (c != '-' && c != '.' && c != '~') {
      Serial.println("appending");
      Serial.println(c);
      Serial.println(item);
      item += c;
    } else if (c == '-' || c == '.') {
      Serial.println("found item");
      Serial.println(item);
      if (count == 0){
        r = item.toInt();
      }else if(count == 1) {
        g = item.toInt();
      }else if(count == 2) {
        b = item.toInt();
      }else {
        bright = item.toInt();
      }
      count++;  
      item = "";
    }

    if (c == '.') {
      break;
    }
  }
  ColorLights(IP, r, g,b, bright);
}

void ColorLights(String IP, int r, int g, int b, int bright) {

  if (bright == 0) {

    bright = 1;
  }
  Serial.println("coloring lights");
  int R = (int)(r * ((float)bright/100.0));
  int G = (int)(g * ((float)bright/100.0));
  int B = (int)(b * ((float)bright/100.0));
  HTTPClient http;   
  http.begin("http://" + IP + "/color?r=" + String(R) +"&g=" + String(G) + "&b=" + String(B)); 
  int httpResponseCode = http.GET(); 
  http.end();
}

void ToggleLight(String IP) {
  HTTPClient http;   
  http.begin("http://" + IP + "/toggle"); 
  int httpResponseCode = http.GET(); 
  http.end();

  //Serial.println("toggled");
  //Serial.println(IP);
  //Serial.println(httpResponseCode);
}

void UpdateCount(String IP, uint8_t count) {
  HTTPClient http;   
  http.begin("http://" + IP + "/count?c="+String(count)); 
  int httpResponseCode = http.GET(); 
  http.end();
}

String cloudReportJSON() {
  String json = "{";
  json += "\"hubID\":\"" + DeviceID + "\",";
  json += "\"hubKey\":\"" + DeviceKey + "\",";
  json += "\"lights\":" + currentLightsJSON();
  json += "}";
  //Serial.println(json);
  return json;
  
}

void handleLights() {
  String lightsJSON = currentLightsJSON();
  server.send(200, "application/json", lightsJSON);
}

String currentLightsJSON() {
  String json = "[";
  for (int i = 0; i < MAX_LIGHTS; i++) {

    if (Lights[i].IP == "") {
      break;
    }

    if (i > 0 && i < MAX_LIGHTS - 1 && Lights[i].ID != "") {
      json += ",";
    }
    
    json += "\"" + Lights[i].ID + "\"";
  }

  json += "]";
  return json;
}

void handleRegister() {
  String id = server.arg("id");
  String ip = server.arg("ip");
  Serial.println("handling: /register");
  Serial.println(id);
  Serial.println(ip);

  int light = LightCount;

  for (int i = 0; i < MAX_LIGHTS; i++) {
    if ( id == Lights[i].ID ) {
      light = i;
    }
  }

  Lights[light].IP = ip;
  Lights[light].ID = id;
  Lights[light].Status = false;

  LightCount++;

  if (LightCount > MAX_LIGHTS) {
    LightCount = 0;
  }

  server.send(200, "text/html", "OK");  
}

void handleIP() {
    Serial.println("Handler: /ip");
    server.send(200, "text/html", WiFi.localIP().toString()); 
}

void handleNotFound() {
  if (captivePortal()) { // If caprive portal redirect instead of displaying the error page.
    return;
  }
  String message = "File Not Found\n\n";
  message += "URI: ";
  message += server.uri();
  message += "\nMethod: ";
  message += ( server.method() == HTTP_GET ) ? "GET" : "POST";
  message += "\nArguments: ";
  message += server.args();
  message += "\n";

  for ( uint8_t i = 0; i < server.args(); i++ ) {
    message += " " + server.argName ( i ) + ": " + server.arg ( i ) + "\n";
  }
  server.sendHeader("Cache-Control", "no-cache, no-store, must-revalidate");
  server.sendHeader("Pragma", "no-cache");
  server.sendHeader("Expires", "-1");
  server.send ( 404, "text/plain", message );
}

void handleIndex() {
    Serial.println("Handler: /");

    if (captivePortal()) { // If caprive portal redirect instead of displaying the page.
      return;
    }
  
    server.send(200, "text/html", WifiSetup()); 
}

String WifiSetup() {

  int n = WiFi.scanNetworks();
   
  String ssids[n];

  for (int i = 0; i < n; ++i) {      
    ssids[i] = WiFi.SSID(i) ;
    delay(10);
  }

  String options = "";
  for (int i = 0; i < n; ++i) {
    options += "<option value=\"" +  ssids[i] + "\">" + ssids[i] + "</option>\n";
  }
  
  return 
"<html>"
"  <head>"
"    <meta content=\"width=device-width, initial-scale=1\" name=\"viewport\" />"
""
"    <style>"
"      .text{"
"        font-size: 2em;"
"      }"
""
"      input[type=password], select {"
"        width: 100%;"
"        padding: 12px;"
"        border: 1px solid #ccc;"
"        border-radius: 4px;"
"        resize: vertical;"
"        font-size: 25px;"
"      }"
""
"      #container{"
"        background: #f2f2f2;"
"        padding-left: 1%;"
"        padding-right: 1%;"
"        font-size: 25px;"
"        margin: 0 auto;"
"        width: 75%;"
"      }"
""
"      label{"
"        margin-bottom: 10px;"
"        display: block;"
"      }"
""
"      .subtext{"
"        color: darkgrey;"
"      }"
""
"      input[type=submit] {"
"        float: right;"
"        background-color: green;"
"        color: white;"
"        border: none;"
"        border-radius: 4px;"
"        margin-top: 10px;"
"        padding: 12px 20px;"
"        font-size: 25px;"
"      }"
""
"      @media only screen and (max-width: 600px) {"
"        #container{"
"          margin: 0;"
"          width: 98%;"
"          font-size: 1.4em;"
"        }"
"      }"
"    </style>"
"  </head>"
"  <body>"
"    <p class=\"text\">Welcome to LightHub!</p>"
"    <div id=\"container\">"
"      <p class=\"subtext\">Please configure your hub by connecting to your WiFi</p>"
"      <form action=\"/login\">"
"        <label for=\"wifi\">WiFi</label>"
"        <select name=\"ssid\">" + options +
"        </select><br><br>"
"        <label for=\"password\">Password</label>"
"        <input type=\"password\" id=\"password\" name=\"password\"><br><br>"
"        <input type=\"submit\" value=\"Submit\">"
"      </form>"
"    </div>"
"  </body>"
"</html>";
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


boolean captivePortal() {
  if (!isIp(server.hostHeader())) {
    Serial.print("Request redirected to captive portal");
    server.sendHeader("Location", String("http://") + toStringIp(server.client().localIP()), true);
    server.send ( 302, "text/plain", ""); // Empty content inhibits Content-length header so we have to close the socket ourselves.
    server.client().stop(); // Stop is needed because we sent no content length
    return true;
  }
  return false;
}


boolean isIp(String str) {
  for (int i = 0; i < str.length(); i++) {
    int c = str.charAt(i);
    if (c != '.' && (c < '0' || c > '9')) {
      return false;
    }
  }
  return true;
}

/** IP to String? */
String toStringIp(IPAddress ip) {
  String res = "";
  for (int i = 0; i < 3; i++) {
    res += String((ip >> (8 * i)) & 0xFF) + ".";
  }
  res += String(((ip >> 8 * 3)) & 0xFF);
  return res;
}

void setup() {
  
   char ssid[32];
   char password[32]; 

   Serial.begin(115200);

   EEPROM.begin(EEPROM_TOTAL_SIZE);
   
   delay(500);
   
   bool hasCredentials = readCreds(ssid, password);
   if (hasCredentials) {
       WifiConnected = startWiFi(ssid, password, 60);
   }
   
   WiFi.mode(WIFI_MODE_APSTA);  
   startAP(ssidAP, passwordAP);
   dnsServer.start(53, "*", WiFi.softAPIP());

}

unsigned long last = millis();

void loop() {
  dnsServer.processNextRequest();
  server.handleClient();

  if (millis() > last + 1000) {
    Serial.println("Report");
    reportCloud();
    last = millis();
  }
}
