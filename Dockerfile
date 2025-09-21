# ---- Build Stage ----
FROM gradle:9.1.0-jdk17 AS build
WORKDIR /app
COPY . .
RUN gradle clean shadowJar --no-daemon

# ---- Run Stage ----
FROM eclipse-temurin:17-jre
WORKDIR /app
COPY --from=build /app/build/libs/churkabot-tg-*-all.jar app.jar

# Set environment variable for bot token (override at runtime)
ENV CHURKA_BOT_TOKEN=changeme

ENTRYPOINT ["java", "-jar", "app.jar"]
