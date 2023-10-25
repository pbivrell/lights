#include <WiFi.h>
#include <WebServer.h>
#include <Update.h>
#include <Adafruit_NeoPixel.h>
#include <EEPROM.h>
#include <FS.h>
#include <HTTPClient.h>

#define MAX_INSTRUCTIONS 20000
#define LED_PIN     15
#define LED_DEFAULT 1000

const String Version = "0.3.0";

uint8_t LightCount = 50;

const String DeviceKey = "theoriginallightskey";
const String DeviceID = "OriginalLight1";

const char* ssid     = "LightHub";
const char* password = "defaultpassword";

const String hubURL = "http://8.8.4.4/register";

Adafruit_NeoPixel strip(LED_DEFAULT, LED_PIN, NEO_GRB + NEO_KHZ800);

bool FoundHub = false;
bool LightStatus = false, LightChange = false;
bool uploading = false, uploaded = false;
String runningFile = "default.bin";

WebServer server(80);

void startServer() {

  server.stop();

  server.on("/count", handleCount);
  server.on("/color", handleColor);
  server.on("/toggle", handleToggle);
  server.on("/status", handleStatus);
  server.on("/upload",HTTP_POST, []() {
    server.send(200, "text/plain", "");
  }, handleFileUpload);
  server.enableCORS();

  server.begin();

}

bool reportHub() {
  
  HTTPClient http;   
  http.begin(hubURL); 
  http.addHeader("Content-Type", "application/x-www-form-urlencoded");
  int httpResponseCode = http.POST("id="+DeviceID+"&ip="+WiFi.localIP().toString()); 
  String response = http.getString();  
  return true;
}

String statusJSON(boolean secret) {
    String toggled = LightStatus ? "true" : "false";
    String key = secret ? "" : DeviceKey;
    String json =  "{" 
      "\"status\":"  + toggled + "," +
      "\"ip\":\""  + WiFi.localIP().toString() + "\"," +
      "\"mac\":\""  + WiFi.macAddress() + "\"," + 
      "\"key\":\""  + key + "\"," +
      "\"id\":\"" + DeviceID + "\"," +
      "\"version\":\"" + Version + "\"" +
    "}";
    return json;
}

void handleCount() {
  Serial.println("Handler: /count");
  LightCount  = server.arg("c").toInt();
  server.send(200, "text/html", "OK");
}

void handleStatus() {
  Serial.println("Handler: /status");
  String json = statusJSON(true);  
  server.send(200, "application/json", json);
}

void handleToggle() {

  Serial.println("Handler: /toggle");

  LightStatus = !LightStatus;
  LightChange = true;

  Serial.println("Changed lights status");
  Serial.println(LightStatus);

  server.send(200, "text/html", "OK"); 
}

size_t instructions = 0;
static uint8_t dataBuffer[MAX_INSTRUCTIONS];


void handleColor() {

  Serial.println("handling: /color");
  int r = server.arg("r").toInt();
  int g = server.arg("g").toInt();
  int b = server.arg("b").toInt();
  Serial.println(r);
  Serial.println(g);
  Serial.println(b);
  setColor(r,g,b);
}

void setColor(uint8_t r, uint8_t g, uint8_t b) {
  instructions = 3;
  dataBuffer[0] = 0x03;
  dataBuffer[1] = LightCount;
  dataBuffer[2] = 0x0;
  dataBuffer[3] = 0x0;
  dataBuffer[4] = 0x0;
  dataBuffer[5] = 0x0;
  dataBuffer[6] = 0x0; 
  dataBuffer[7] = 0x0;    
        
  dataBuffer[0 + 8] = 0x04;
  dataBuffer[1 + 8] = r;
  dataBuffer[2 + 8] = g;
  dataBuffer[3 + 8] = b;
  dataBuffer[4 + 8] = 0x0;
  dataBuffer[5 + 8] = 0x0;
  dataBuffer[6 + 8] = LightCount; 
  dataBuffer[7 + 8] = 0x0;

  dataBuffer[0 + 16] = 0x02;
  dataBuffer[1 + 16] = 0xE8;
  dataBuffer[2 + 16] = 0x03;
  dataBuffer[3 + 16] = 0x0;
  dataBuffer[4 + 16] = 0x0;
  dataBuffer[5 + 16] = 0x0;
  dataBuffer[6 + 16] = 0x0; 
  dataBuffer[7 + 16] = 0x0;
  LightStatus = true;
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
  return WiFi.status() == WL_CONNECTED;
}

size_t upload_size = 0;

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

void serveTime(uint32_t delayAmount){ 

  unsigned long start = millis();

  Serial.println("delaying");

  for(;;) {

    if  (WiFi.status() != WL_CONNECTED) {
      FoundHub = false;
      if (startWiFi(ssid, password, 10)) {
        Serial.println("ReConnected to Hub");
        startServer();
        if (reportHub()) {
          FoundHub = true;
        }
      }
    }

    // serve a connection
    server.handleClient();

    // Light status has been changed, the current state of the lights is invalid
    if (LightChange) {
      LightChange = false;
      break;
    }

    // Something new has been uploaded restart the animation
    // Regardless of the wait time
    if (uploaded) {
      uploading = false;
      uploaded = false;
      //Serial.println("uploaded terminating serve");
      break;
    }

    // The amount of delay has expired return to possibly running
    // animiations
    if (millis() > start + delayAmount) {
      //Serial.println("no longer delaying");
      //Serial.println(start);
      //Serial.println(millis());
      //Serial.println(delayAmount);
      //delay(500);
      break;
    }
  }
}

int MissingHubColor = 0;

void runInstructions(  Adafruit_NeoPixel * strip){

  strip->clear();

  if (!LightStatus) {
    strip->show();

    // Lights are off. Don't show anything. Give the server
    // time to server as that is the only way this will change
    serveTime(500);
    return;
  }

  if (uploading) {
     // Don't start a new animation while something is being actively uploaded
     serveTime(500);
     return;
  }

  uint16_t index;
  uint16_t amount;
  uint32_t delayAmount;
  uint16_t count; 
  uint8_t r, g, b;

  Serial.println("Running instructions");

  uint8_t lastInstruction = 0;

  int counter = 0;

  Serial.println(instructions);
  
  while (counter < instructions) {

    if (!LightStatus) {
      strip->clear();
      strip->show();
      break;
    }

    if (uploading) {
      break;
    }
    
    if(dataBuffer[(counter * 8) + 0] == 0x3) {
      count = (dataBuffer[(counter * 8) + 2] << 8) | dataBuffer[(counter * 8) + 1];
      strip->updateLength(count);
      strip->clear();
      //Serial.println("Set Count");
      //Serial.println(count);
    }else if(dataBuffer[(counter * 8) + 0] == 0x2) {
      delayAmount = (dataBuffer[(counter * 8) + 4] << 24) | (dataBuffer[(counter * 8) + 3] << 16) | (dataBuffer[(counter * 8) + 2] << 8) | dataBuffer[(counter * 8) + 1];
      strip->show();
     
      //Serial.println("Delay");
      //Serial.println(delayAmount);

      // While animation is waiting. Run the server
      serveTime(delayAmount);
      
    }else if(dataBuffer[(counter * 8) + 0] == 0x1) {
      
      r = dataBuffer[(counter * 8) + 1];
      g = dataBuffer[(counter * 8) + 2];
      b = dataBuffer[(counter * 8) + 3];
      index = (dataBuffer[(counter * 8) + 5] << 8) | dataBuffer[(counter * 8) + 4];
      
      //Serial.println("color");
      //Serial.println(r);
      //Serial.println(g);
      //Serial.println(b);
      //Serial.println(index);
      strip->setPixelColor(index, strip->Color(g, r, b)); 

    }else if(dataBuffer[(counter * 8) + 0] == 0x4) {
      r = dataBuffer[(counter * 8) + 1];
      g = dataBuffer[(counter * 8) + 2];
      b = dataBuffer[(counter * 8) + 3];
      index = (dataBuffer[(counter * 8) + 5] << 8) | dataBuffer[(counter * 8) + 4];
      amount = (dataBuffer[(counter * 8) + 7] << 8) | dataBuffer[(counter * 8) + 6];

      for (uint16_t i = 0; i < amount; i++) {
        strip->setPixelColor(i+index, strip->Color(g, r, b)); 
      }
      Serial.println("set colors");
    }else {
      Serial.println("Invalid");
      //Serial.println(dataBuffer[(counter * 8) + 0]);
      //Serial.println(dataBuffer[(counter * 8) + 1]);
      //Serial.println(dataBuffer[(counter * 8) + 2]);
      //Serial.println(dataBuffer[(counter * 8) + 3]);
      //Serial.println(dataBuffer[(counter * 8) + 4]);
      //Serial.println(dataBuffer[(counter * 8) + 5]);
      //Serial.println(dataBuffer[(counter * 8) + 6]);
      //Serial.println(dataBuffer[(counter * 8) + 7]);
    }

    if (!FoundHub){
      strip->setPixelColor(0, strip->Color(37, MissingHubColor, 37));
      strip->show();
    }

    MissingHubColor+=20;
    if (MissingHubColor > 255) {
      MissingHubColor = 0;
    }

    lastInstruction = dataBuffer[(counter * 8) + 0];
    counter++;
  }

  if (lastInstruction != 0x2) {
    Serial.println("forced delay");
    serveTime(1000);
  }
}

void setup() {

  Serial.begin(115200);

  strip.begin(); 
   
  delay(500);

  /*if (startWiFi(ssid, password, 60)) {
    Serial.println("Connected to Hub");
    startServer();
    if (reportHub()){
      FoundHub = true;
    }
   
  }else {

  } */ 

  setColor(0xFC, 0xB6, 0x12);      
  
}


void loop() {
  runInstructions(&strip);
}
