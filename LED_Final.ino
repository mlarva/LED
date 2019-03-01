#include <Ethernet.h>
#include <SPI.h>
#include <PubSubClient.h>
#include <Adafruit_NeoPixel.h>

#define PIN 14

Adafruit_NeoPixel strip = Adafruit_NeoPixel(60, PIN, NEO_RGBW);

byte mac[] = { 0xDE, 0xAD, 0xBE, 0xEF, 0xFE, 0xED };
const char* mqttServer = "192.168.0.18";
const int mqttPort = 1883;
EthernetClient ethClient;
PubSubClient client(ethClient);

void setup() {
  Serial.begin(9600);
  Serial.println("Initializing");
  Serial.println(Ethernet.begin(mac)); 
  Serial.println(Ethernet.localIP());
  
  client.setServer(mqttServer, mqttPort);
  client.setCallback(callback);
  while (!client.connected()) {
    Serial.println("Connecting to MQTT...");
    if (client.connect("ESP8266Client")) {
      Serial.println("connected");  
    } else {
      Serial.print("failed with state ");
      Serial.print(client.state());
      delay(2000);
    }
  }
  client.publish("rgb", "Hello from Arduino");
  client.subscribe("rgb");
  
  strip.setBrightness(100);
  for(uint16_t i=0; i<strip.numPixels(); i++) {
    strip.setPixelColor(i, strip.Color(255, 0, 50));
  }
  strip.begin();
  strip.show(); // Initialize all pixels to 'off'
}

void loop() {
  client.loop();
}

void callback(char* topic, byte* payload, unsigned int length) {
  Serial.print("Message arrived in topic: ");
  Serial.println(topic);
  Serial.print("Message:");
  for (int i = 0; i < length; i++) {
    Serial.print((char)payload[i]);
  }
  Serial.println();
  Serial.println("-----------------------");
  parseCommand(payload, length);
}

void parseCommand(byte* payload, unsigned int length) {
  String command;
  String sRed;
  String sGreen;
  String sBlue;
  int nRed;
  int nGreen;
  int nBlue;
  int commaLocation;
  for (int i = 0; i < length; i++) {
    command.concat((char)payload[i]);
  }
  command.remove(0,1);
  command.remove(command.length()-1);
  commaLocation = command.indexOf(',');  
  sGreen = command.substring(0, commaLocation);  
  command.remove(0,commaLocation+1);
  commaLocation = command.indexOf(',');  
  sRed = command.substring(0, commaLocation);
  command.remove(0,commaLocation+1);
  sBlue = command;   //captures first data String
  nRed = sRed.toInt();
  nGreen = sGreen.toInt();
  nBlue = sBlue.toInt();
  for(uint16_t i=0; i<strip.numPixels(); i++) {
    strip.setPixelColor(i, strip.Color(nRed, nGreen, nBlue));
  }
  strip.show();
}
