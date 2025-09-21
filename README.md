# My Java App

## Overview
This project is a Java application that utilizes the Telegram Bots API to create a bot using long polling. It includes a main application class and unit tests to ensure functionality.

## Project Structure
```
my-java-app
├── src
│   ├── main
│   │   ├── java
│   │   │   └── App.java
│   │   └── resources
│   └── test
│       └── java
│           └── AppTest.java
├── build.gradle
├── settings.gradle
└── README.md
```

## Prerequisites
- Java Development Kit (JDK) 8 or higher
- Gradle 6.0 or higher

## Setup
1. Clone the repository:
   ```
   git clone <repository-url>
   cd my-java-app
   ```

2. Build the project using Gradle:
   ```
   ./gradlew build
   ```

## Running the Application
To run the application, execute the following command:
```
./gradlew run
```

## Running Tests
To run the unit tests, use the command:
```
./gradlew test
```

## Dependencies
This project uses the following dependencies:
- `org.telegram:telegrambots-longpolling:9.1.0`
- `org.telegram:telegrambots-client:9.1.0`

## License
This project is licensed under the MIT License. See the LICENSE file for more details.