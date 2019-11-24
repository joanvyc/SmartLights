#include <ESP8266WiFi.h>
#include <WString.h>

#define LDR_PIN A0 // A0
#define LIGHT_SIZE 20

//server config
const char* ssid = "marenostrum";
const char* password = "0123456789";
const char * host = "192.168.1.118";
const uint16_t port = 3000;

int light_values [LIGHT_SIZE];
int light_counter = 0;
int flow_counter = 0;

void setup () {
    //Setup serial for log
    Serial.begin(115200);
    Serial.println();
    Serial.println();

    Serial.print("Connecting to ");
    Serial.println(ssid);

    //connect to network
    WiFi.mode(WIFI_STA);
    WiFi.begin(ssid, password);

    //wait for wifi connection
    while(WiFi.status() != WL_CONNECTED) {
        delay(500);
        Serial.print(".");
    }

    //Log the newtwork info
    Serial.println("");
    Serial.println("WiFi connnected");
    Serial.println("IP address: ");
    Serial.println(WiFi.localIP());

    // LDR
    for (int i = 0; i < LIGHT_SIZE; i++) light_values[i] = 0;
}


bool getPetition(WiFiClient &cl,  String path){
    cl.print("GET " + path + " HTTP/1.1\r\n" + "Host: " + host + "\r\n" + "Connection: close\r\n" + "\r\n");
}

bool post_petition (WiFiClient &cl,  String path, String data) {
    cl.print("POST " + path + " HTTP/1.1\r\n" + "Host: " + host + "\r\n" + "Accept: application/json\r\n" + "Content-Type: application/json\r\n" + "Content-Length: " + data.length() + "\r\n" + "Connection: close\r\n" + "\r\n" + data + "\r\n");
    //cl.print("POST " + path + " HTTP/1.1\r\n" + "Host: " + host + "\r\n" + "Content-Type: application/json\r\n" + "Content-Lenght: " + data.length() + "\r\n" + "Connection: close\r\n" + "\r\n" + data + "\r\n");
}

void send_data (String data) {
    Serial.println("Sending data...");
    WiFiClient cl;
    if (!cl.connect(host, port)){
        Serial.println("Connection failed");
        return;
    }

    if (cl.connected()){
        post_petition(cl, "/api/postData", data);
    }
}
/*
   bool send_key (String key){
   WiFiClient cl;
   if (!cl.connect(host, port)) {
   Serial.println("Connection failed");
   }

   if (cl.connected()) {
   put_petition(cl, "/api/user_enter", key);
//        getPetition(cl, "/api/user_enter");
Serial.println("Connection sended ");
Serial.print(key);
return true;
}
return false;
}

 */
void delay_s(int s){
    for(int i = 0; i<s; i++) delay(1000);
}

String json_parser () {
    
    String json_data = String("{ \"flow_label\": ") + flow_counter + ", \"samples\": [ " + light_values[0];
    for (int i = 1 ; i < LIGHT_SIZE; i++) json_data += String(", ") + light_values[i];
    json_data += "] } \0";
    Serial.println( "JSON");
    Serial.println(json_data);
    return json_data;
}


void loop() {
    /*
       if(doorOpen()){
       door_time = millis();

       if (in_house){
    //user was in house -> exiting
    sendOutHouse();
    }else{
    //user was outside -> is entering
    sendInHouse();
    }
    }
     */

    int tmp = analogRead(LDR_PIN);
    Serial.println(tmp);
    light_values[light_counter] = tmp; 

    delay(1000);
    if (light_counter == LIGHT_SIZE-1) {
        String json_data = json_parser();
        flow_counter++;
        send_data(json_data);
    }

    //counter = (++counter)%LIGHT_SIZE;
    Serial.println(light_counter);
    light_counter = (++light_counter)%LIGHT_SIZE;

}

