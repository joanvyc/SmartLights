#include <ESP8266WiFi.h>
#include <WString.h>

#define LDR_PIN A0 // A0
#define LIGHT_SIZE 20

//server config
const char* ssid = "marenostrum";
const char* password = "0123456789";
const char * host = "192.168.1.118";
const uint16_t port = 3000;

struct {
    uint32_t crc32;
    int light_counter;
    int light_values [LIGHT_SIZE];
    int flow_counter;
} data;
/*
   int light_values [LIGHT_SIZE];
   int light_counter = 0;
   int flow_counter = 0;
 */
uint32_t calculateCRC32(const uint8_t *data, size_t length) {
    uint32_t crc = 0xffffffff;
    while (length--) {
        uint8_t c = *data++;
        for (uint32_t i = 0x80; i > 0; i >>= 1) {
            bool bit = crc & 0x80000000;
            if (c & i) {
                bit = !bit;
            }
            crc <<= 1;
            if (bit) {
                crc ^= 0x04c11db7;
            }
        }
    }
    return crc;
}

//prints all rtcData, including the leading crc32
void printMemory() {
    char buf[3];
    uint8_t *ptr = (uint8_t *)&data;
    for (size_t i = 0; i < sizeof(data); i++) {
        sprintf(buf, "%02X", ptr[i]);
        Serial.print(buf);
        if ((i + 1) % 32 == 0) {
            Serial.println();
        } else {
            Serial.print(" ");
        }
    }
    Serial.println();
}


String json_parser () {

    String json_data = String("{ \"flow_label\": ") + data.flow_counter + ", \"samples\": [ " + data.light_values[0];
    for (int i = 1 ; i < LIGHT_SIZE; i++) json_data += String(", ") + data.light_values[i];
    json_data += "] } \0";
    Serial.println( "JSON");
    Serial.println(json_data);
    return json_data;
}

void init_wifi() {
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
}

bool post_petition (WiFiClient &cl,  String path, String data) {
    Serial.println("POST");
    cl.print("POST " + path + " HTTP/1.1\r\n" + "Host: " + host + "\r\n" + "Accept: application/json\r\n" + "Content-Type: application/json\r\n" + "Content-Length: " + data.length() + "\r\n" + "Connection: close\r\n" + "\r\n" + data + "\r\n");
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

void setup() {
    Serial.begin(115200);
    Serial.setTimeout(2000);
    // Read struct from RTC memory
    if (ESP.rtcUserMemoryRead(0, (uint32_t*) &data, sizeof(data))) {
        Serial.println("Read: ");
        printMemory();
        Serial.println();
        uint32_t crcOfData = calculateCRC32((uint8_t*) &data.light_counter, sizeof(data) - sizeof(data.crc32));
        Serial.print("CRC32 of data: ");
        Serial.println(crcOfData, HEX);
        Serial.print("CRC32 read from RTC: ");
        Serial.println(data.crc32, HEX);
        if (crcOfData != data.crc32) {
            Serial.println("CRC32 in RTC memory doesn't match CRC32 of data. Data is probably invalid!");
            data.crc32 = 0;
            data.light_counter = 0;
            for (int i = 0; i < LIGHT_SIZE; i++) data.light_values[i] = 0;
            data.flow_counter = 0;
        } else {
            Serial.println("CRC32 check ok, data is probably valid.");
        }
    }

    // Generate new data set for the struct
    int tmp = analogRead(LDR_PIN);
    data.light_values[data.light_counter] = tmp;
    data.light_counter++;
    Serial.println(tmp);
    Serial.println(data.light_counter);

    if (data.light_counter == LIGHT_SIZE) {
        init_wifi();
        send_data(json_parser());
        data.light_counter = 0;
        data.flow_counter++;
        delay(1000);
    }

    data.crc32 = calculateCRC32((uint8_t*) &data.light_counter, sizeof(data) - sizeof(data.crc32));

    // Write struct to RTC memory
    if (ESP.rtcUserMemoryWrite(0, (uint32_t*) &data, sizeof(data))) {
        Serial.println("Write: ");
        printMemory();
        Serial.println();
    }


    // Deep sleep mode for 30 seconds, the ESP8266 wakes up by itself when GPIO 16 (D0 in NodeMCU board) is connected to the RESET pin
    Serial.println("I'm awake, but I'm going into deep sleep mode for 30 seconds");
    ESP.deepSleep(4e6);

}

void loop() {

}

